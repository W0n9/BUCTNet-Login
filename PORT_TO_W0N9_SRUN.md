# 移植BUCTNet-Login新特性到w0n9/srun的详细指南

## 概述
本文档详细说明了如何将BUCTNet-Login仓库中的新特性移植到w0n9/srun仓库。

## 主要新特性列表

### 1. 增强的日志系统 (logger包)
- **替换**: 将logrus替换为zap logger
- **优势**: 更好的性能、结构化日志、更灵活的配置
- **文件**: `logger/logger.go`

### 2. Keepalive保持在线功能 ⭐
- **功能**: 定时检查网络状态，自动重连
- **命令**: `srun keepalive -i 30`
- **文件**: 在`cmd/srun/main.go`和`cmd/srun/cli.go`中实现

### 3. 改进的错误处理
- **模式**: 实现*E函数返回错误，而不是直接处理
- **好处**: 更好的错误传播和处理

### 4. 版本和依赖更新
- **版本**: v1.1.5 → v1.1.8
- **Go版本**: 1.18 → 1.24
- **依赖**: 更新到最新稳定版本

## 详细移植步骤

### 步骤1: 更新go.mod和依赖
```go
module github.com/vouv/srun  // 或w0n9/srun

go 1.24

require (
    github.com/go-resty/resty/v2 v2.16.5  // 新增
    github.com/spf13/cobra v1.9.1         // 更新
    go.uber.org/zap v1.27.0               // 新增
    golang.org/x/term v0.34.0             // 更新
)
```

### 步骤2: 创建logger包
创建 `logger/logger.go` 文件，实现zap日志系统。

### 步骤3: 更新main.go
在`cmd/srun/main.go`中添加：
- keepalive命令
- zap logger初始化
- 版本更新到v1.1.8

### 步骤4: 更新cli.go
在`cmd/srun/cli.go`中添加：
- KeepaliveE函数
- ensureOnline函数
- 改进的错误处理模式

### 步骤5: 测试和验证
- 构建项目
- 测试keepalive功能
- 验证日志输出
- 确保向后兼容性

## 关键代码差异

### 日志系统差异
**原始 (logrus)**:
```go
import log "github.com/sirupsen/logrus"
log.Info("message")
```

**新版 (zap)**:
```go
import "github.com/W0n9/BUCTNet-Login/logger"
log := logger.GetLogger()
log.Info("message")
```

### Keepalive功能
这是一个全新的功能，需要完整实现：
- 定时器检查网络状态
- 信号处理
- 自动重连逻辑

## 预期影响
1. **性能提升**: zap logger提供更好的性能
2. **用户体验**: keepalive功能解决频繁掉线问题
3. **代码质量**: 更好的错误处理和日志记录
4. **维护性**: 更清晰的代码结构

## 注意事项
- 确保向后兼容性
- 测试所有现有功能
- 更新文档和README
- 考虑渐进式迁移策略