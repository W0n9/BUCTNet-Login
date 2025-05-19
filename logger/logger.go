package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

// InitLogger 初始化日志配置
func InitLogger(debug bool) {
	var cfg zap.Config
	if debug {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.CallerKey = "caller"
		cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	} else {
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		// 设置时间格式为人类可读的格式 (ISO8601)
		// cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// 使用控制台编码器替代默认的JSON编码器
		cfg.Encoding = "console"
		cfg.DisableCaller = true
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	log = logger.Sugar()
}

// GetLogger 获取日志实例
func GetLogger() *zap.SugaredLogger {
	if log == nil {
		InitLogger(false)
	}
	return log
}

// Debug 调试日志
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}
