// 这些是需要添加到w0n9/srun的cmd/srun/cli.go文件中的新函数
// 主要是keepalive相关功能

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	// 这里需要适配到目标仓库的import路径
	// "github.com/vouv/srun/core"
	// "github.com/vouv/srun/store"
	// "github.com/vouv/srun/model"
)

// Keepalive 保持在线功能，定期检查网络状态并自动重连
func Keepalive(cmd *cobra.Command, args []string) {
	err := KeepaliveE(cmd, args)
	if err != nil {
		log.Error(err)
	}
}

// KeepaliveE 保持在线功能，返回错误信息
// 这是主要的新特性实现
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
// 这是keepalive功能的核心逻辑
func ensureOnline(account *model.Account) error {
	// 尝试获取信息，检查是否在线
	info, err := core.Info()
	if err != nil {
		log.Error(fmt.Sprintf("获取信息失败: %s", err.Error()))
	}
	if info == nil || info.UserName == "" {
		log.Info("检测到网络离线，尝试重新登录...")
		return core.Login(account)
	}

	log.Debug("网络在线，无需重连")
	return nil
}

// 改进的错误处理模式示例
// 原有函数的*E版本，返回错误而不是直接处理

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

func InfoE(cmd *cobra.Command, args []string) error {
	info, err := core.Info()
	if err != nil {
		return err
	}
	fmt.Println(info.String())
	return nil
}

func ConfigE(cmd *cobra.Command, args []string) error {
	// 配置逻辑保持不变，但返回错误
	// ... 原有配置代码 ...
	return nil
}