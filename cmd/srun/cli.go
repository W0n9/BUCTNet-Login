package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/vouv/srun/core"
	"github.com/vouv/srun/model"
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

// Keepalive 保持在线功能，定期检查网络状态并自动重连
func Keepalive(cmd *cobra.Command, args []string) {
	err := KeepaliveE(cmd, args)
	if err != nil {
		log.Error(err)
	}
}

// KeepaliveE 保持在线功能，返回错误信息
func KeepaliveE(cmd *cobra.Command, args []string) error {
	interval, err := cmd.Flags().GetInt("interval")
	if err != nil {
		interval = 30 // 默认30秒检查一次
	}

	account, err := store.ReadAccount()
	if err != nil {
		return err
	}

	log.Info("启动保持在线功能...")
	log.Info(fmt.Sprintf("检查间隔: %d秒", interval))
	log.Info(fmt.Sprintf("账号: %s", account.Username))
	log.Info("按 Ctrl+C 退出")

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 定时器，定期检查网络状态
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// 启动时先尝试登录一次
	if err := ensureOnline(account); err != nil {
		log.Info(fmt.Sprintf("初始登录尝试: %s", err.Error()))
	}

	// 主循环
	for {
		select {
		case <-ticker.C:
			// 检查网络状态并确保在线
			if err := ensureOnline(account); err != nil {
				log.Debug(fmt.Sprintf("保持在线出错: %s", err.Error()))
			}
		case <-sigChan:
			// 收到退出信号
			log.Info("正在退出保持在线功能...")
			return nil
		}
	}
}

// ensureOnline 确保网络在线，如果不在线则尝试登录
func ensureOnline(account *model.Account) error {
	// 尝试获取信息，检查是否在线
	info, err := core.Info()
	if info.UserName == "" {
		log.Info("检测到网络离线，尝试重新登录...")
		return core.Login(account)
	}
	if err != nil {
		log.Info("检测到网络离线，尝试重新登录...")
		return core.Login(account)
	}

	log.Debug("网络在线，无需重连")
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
		fmt.Sprintf("\tVersion: %s\n", Version)
}
