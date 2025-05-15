package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vouv/srun/core"
	"github.com/vouv/srun/store"
	"golang.org/x/term"
)

// Login 登录函数，处理登录流程
func Login(cmd *cobra.Command, args []string) {
	err := LoginE(cmd, args)
	if err != nil {
		log.Error(err)
	}
}

// LoginE 登录函数，返回错误信息
func LoginE(cmd *cobra.Command, args []string) error {
	account, err := store.ReadAccount()
	if err != nil {
		return err
	}
	log.Info("尝试登录...")

	if err = core.Login(account); err != nil {
		return err
	}
	log.Info("登录成功!")

	return store.WriteAccount(account)
}

// Logout 注销函数，处理注销流程
func Logout(cmd *cobra.Command, args []string) {
	err := LogoutE(cmd, args)
	if err != nil {
		log.Error(err)
	}
}

// LogoutE 注销函数，返回错误信息
func LogoutE(cmd *cobra.Command, args []string) error {
	var err error
	account, err := store.ReadAccount()
	if err != nil {
		return err
	}

	_ = core.Logout(account)

	log.Info("注销成功!")

	return store.WriteAccount(account)
}

// Info 查询信息函数
func Info(cmd *cobra.Command, args []string) {
	err := InfoE(cmd, args)
	if err != nil {
		log.Error(err)
	}
}

// InfoE 查询信息函数，返回错误信息
func InfoE(cmd *cobra.Command, args []string) error {
	info, err := core.Info()
	if err != nil {
		return err
	}
	fmt.Println(info.String())
	return nil
}

// Config 配置账号密码函数
func Config(cmd *cobra.Command, args []string) {
	err := ConfigE(cmd, args)
	if err != nil {
		log.Error(err)
	}
}

// ConfigE 配置账号密码函数，返回错误信息
func ConfigE(cmd *cobra.Command, args []string) error {

	in := os.Stdin
	fmt.Print("设置校园网账号:\n>")
	username := readInput(in)

	// 终端API
	fmt.Print("设置校园网密码(隐私输入):\n>")
	fd := int(os.Stdin.Fd())
	bytePwd, err := term.ReadPassword(fd)
	if err != nil {
		return err
	}
	fmt.Println()
	pwd := string(bytePwd)

	// trim
	username = strings.TrimSpace(username)
	pwd = strings.TrimSpace(pwd)

	if err := store.SetAccount(username, pwd); err != nil {
		return err
	}
	log.Info("账号密码已被保存")
	return nil
}

// readInput 读取输入内容
func readInput(in io.Reader) string {
	reader := bufio.NewReader(in)
	line, _, err := reader.ReadLine()
	if err != nil {
		panic(err)
	}
	return string(line)
}

// VersionString 返回版本信息字符串
func VersionString() string {
	return fmt.Sprintln("System:") +
		fmt.Sprintf("\tOS:%s ARCH:%s GO:%s\n", runtime.GOOS, runtime.GOARCH, runtime.Version()) +
		fmt.Sprintln("About:") +
		fmt.Sprintf("\tVersion: %s\n", Version) +
		fmt.Sprintln("\n\t</> with ❤ By vouv")
}
