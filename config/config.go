package config

import "time"

// Config 定义全局配置
type Config struct {
	// 基础配置
	BaseURL     string
	APITimeout  time.Duration
	DebugMode   bool
	AccountFile string

	// API路径配置
	ChallengeURL string
	PortalURL    string
	SucceedURL   string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		BaseURL:      "http://202.4.130.95",
		APITimeout:   time.Second * 5,
		ChallengeURL: "/cgi-bin/get_challenge",
		PortalURL:    "/cgi-bin/srun_portal",
		SucceedURL:   "/cgi-bin/rad_user_info",
	}
}

var globalConfig = DefaultConfig()

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}

// SetConfig 设置全局配置
func SetConfig(cfg *Config) {
	globalConfig = cfg
}
