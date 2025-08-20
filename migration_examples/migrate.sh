#!/bin/bash
# 移植BUCTNet-Login新特性到w0n9/srun的自动化脚本
# Migration script to port BUCTNet-Login features to w0n9/srun

set -e

echo "🚀 开始移植BUCTNet-Login新特性到w0n9/srun..."

# 检查当前目录是否为srun项目
if [ ! -f "go.mod" ] || ! grep -q "srun" go.mod; then
    echo "❌ 错误：请在srun项目根目录下运行此脚本"
    exit 1
fi

echo "📋 移植步骤："
echo "1. 备份原始文件"
echo "2. 更新go.mod和依赖"
echo "3. 创建logger包"
echo "4. 更新main.go"
echo "5. 更新cli.go"
echo "6. 测试构建"

# 步骤1: 备份原始文件
echo "📂 1. 备份原始文件..."
mkdir -p backup
cp cmd/srun/main.go backup/main.go.backup
cp cmd/srun/cli.go backup/cli.go.backup
cp go.mod backup/go.mod.backup
echo "✅ 备份完成"

# 步骤2: 更新go.mod
echo "📦 2. 更新go.mod和依赖..."
cat > go.mod << 'EOF'
module github.com/vouv/srun

go 1.24

require (
	github.com/go-resty/resty/v2 v2.16.5
	github.com/spf13/cobra v1.9.1
	go.uber.org/zap v1.27.0
	golang.org/x/term v0.34.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
)
EOF

# 下载依赖
go mod tidy
echo "✅ 依赖更新完成"

# 步骤3: 创建logger包
echo "📝 3. 创建logger包..."
mkdir -p logger
cat > logger/logger.go << 'EOF'
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
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
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
EOF
echo "✅ logger包创建完成"

# 步骤4: 更新main.go (只添加keepalive相关部分)
echo "🔧 4. 更新main.go..."
# 这里需要手动操作，因为需要保留原有代码结构
echo "📝 请手动添加以下内容到cmd/srun/main.go:"
echo "   - 导入 logger 包"
echo "   - 添加 keepaliveCmd 命令"
echo "   - 更新版本到 v1.1.8"
echo "   - 添加 zap logger 初始化"

# 步骤5: 测试构建
echo "🔨 5. 测试构建..."
if go build -o bin/srun ./cmd/srun; then
    echo "✅ 构建成功！"
else
    echo "❌ 构建失败，请检查代码"
    exit 1
fi

echo "🎉 移植完成！"
echo "📋 下一步："
echo "1. 手动完成main.go和cli.go的更新"
echo "2. 测试keepalive功能"
echo "3. 更新README文档"
echo "4. 运行测试确保功能正常"

echo "📁 相关文件："
echo "- backup/ : 原始文件备份"
echo "- logger/ : 新的日志包"
echo "- cmd/srun/ : 需要更新的CLI文件"