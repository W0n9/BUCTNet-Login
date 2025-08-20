# 移植到w0n9/srun的go.mod更新

# 原始go.mod (vouv/srun):
```
module github.com/vouv/srun

go 1.18

require (
	github.com/moby/moby v20.10.20+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.1.1
)
```

# 新版go.mod (移植后):
```
module github.com/vouv/srun  # 或改为 github.com/w0n9/srun

go 1.24

require (
	github.com/go-resty/resty/v2 v2.16.5  # 新增：HTTP客户端
	github.com/spf13/cobra v1.9.1         # 更新版本
	go.uber.org/zap v1.27.0               # 新增：zap日志库
	golang.org/x/term v0.34.0             # 更新：终端操作
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
)
```

## 主要依赖变更说明：

### 移除的依赖：
- `github.com/sirupsen/logrus` -> 被 `go.uber.org/zap` 替代
- `github.com/moby/moby` -> 被 `golang.org/x/term` 替代

### 新增的依赖：
- `go.uber.org/zap` - 高性能结构化日志库
- `github.com/go-resty/resty/v2` - 现代HTTP客户端
- `golang.org/x/term` - 终端操作库

### 更新的依赖：
- `github.com/spf13/cobra` v1.1.1 -> v1.9.1
- Go版本 1.18 -> 1.24

## 迁移注意事项：

1. **日志库迁移**：
   - 所有 `log "github.com/sirupsen/logrus"` 需要替换为 zap logger
   - 日志调用方式有所不同，需要适配

2. **终端操作**：
   - `github.com/moby/moby/pkg/term` 替换为 `golang.org/x/term`
   - API可能有细微差别

3. **构建要求**：
   - 需要Go 1.24或更高版本
   - 更新构建脚本和CI配置