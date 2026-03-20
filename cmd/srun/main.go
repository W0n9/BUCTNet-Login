package main

import (
	"os"

	"github.com/W0n9/BUCTNet-Login/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const Version = "v1.1.8"

// loginCmd 显式执行登录子命令
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login srun",
	Run:   Login,
}

// logoutCmd 注销当前会话
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout srun",
	Run:   Logout,
}

// infoCmd 查询套餐/流量等在线信息
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get srun info",
	Run:   Info,
}

// configCmd 交互式写入本地账号密码
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config srun",
	Run:   Config,
}

// keepaliveCmd 定时探测在线，掉线则自动重登
var keepaliveCmd = &cobra.Command{
	Use:   "keepalive",
	Short: "Keep login state and auto reconnect when offline",
	Run:   Keepalive,
}

// rootCmd 根命令；无子命令时默认执行登录
var rootCmd = &cobra.Command{
	Use:   "srun [command]",
	Short: "An efficient client for BUCT campus network",
	RunE: func(cmd *cobra.Command, args []string) error {
		return LoginE(cmd, args)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 按全局 --debug 初始化 zap 日志
		logger.InitLogger(debugMode)
		log = logger.GetLogger()
		if debugMode {
			log.Debug("Debug mode enabled")
		}
	},
}

var debugMode bool
var log *zap.SugaredLogger

// main 解析命令行并分发子命令
func main() {

	defer func() {
		if log != nil {
			log.Sync()
		}
	}()
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "debug mode")

	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(VersionString())

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(configCmd)

	// -i/--interval：保活探测周期（秒）
	keepaliveCmd.Flags().IntP("interval", "i", 30, "检查间隔时间（秒）")
	rootCmd.AddCommand(keepaliveCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
