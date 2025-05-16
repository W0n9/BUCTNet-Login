package main

import (
	"os"

	"github.com/W0n9/BUCTNet-Login/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const Version = "v1.1.6"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login srun",
	Run:   Login,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout srun",
	Run:   Logout,
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get srun info",
	Run:   Info,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config srun",
	Run:   Config,
}

var keepaliveCmd = &cobra.Command{
	Use:   "keepalive",
	Short: "Keep login state and auto reconnect when offline",
	Run:   Keepalive,
}

var rootCmd = &cobra.Command{
	Use:   "srun [command]",
	Short: "An efficient client for BUCT campus network",
	RunE: func(cmd *cobra.Command, args []string) error {
		return LoginE(cmd, args)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize logger with the correct debugMode value
		logger.InitLogger(debugMode)
		log = logger.GetLogger()
		if debugMode {
			log.Debug("Debug mode enabled")
		}
	},
}

var debugMode bool
var log *zap.SugaredLogger

// main 程序入口
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

	// 添加keepalive命令
	keepaliveCmd.Flags().IntP("interval", "i", 30, "检查间隔时间（秒）")
	rootCmd.AddCommand(keepaliveCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
