package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const Version = "v0.0.1"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login buct_net",
	Run:   Login,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout buct_net",
	Run:   Logout,
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get buct_net info",
	Run:   Info,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config buct_net_login",
	Run:   Config,
}

var rootCmd = &cobra.Command{
	Use:   "buct_net_login [command]",
	Short: "An efficient client for BUCT campus network",
	RunE: func(cmd *cobra.Command, args []string) error {
		return LoginE(cmd, args)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debugMode {
			log.SetLevel(log.DebugLevel)
		}
	},
}

var debugMode bool

func main() {

	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "debug mode")

	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(VersionString())

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(configCmd)

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		// DisableTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
