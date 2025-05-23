# BUCTNet-Login


[![Build Status](https://travis-ci.org/vouv/srun.svg?branch=master)](https://travis-ci.org/vouv/srun) [![Go Report Card](https://goreportcard.com/badge/github.com/vouv/srun)](https://goreportcard.com/report/github.com/vouv/srun) ![License](https://img.shields.io/packagist/l/doctrine/orm.svg) [![GoDoc](https://godoc.org/github.com/vouv/srun?status.svg)](https://godoc.org/github.com/vouv/srun/core)

> An efficient client for BUCT campus network

北京化工大学校园网命令行登录工具
- 支持linux、maxOS、windows
- 实验性支持龙芯、部分路由器平台，如MT7621
- 基于Go语言实现

## 如何获取

1. 在[Relesae](https://github.com/W0n9/BUCTNet-Login/releases/latest)处获取
    > 其中MT7621等路由器SOC，选择mipsle架构（于NeWiFi3测试上测试通过）  
    龙芯用户可以尝试选择mips64le架构（因缺少设备，未测试）

2. 自行编译，要求有Golang环境 `go version >= 1.24`   
    1. 先克隆项目
    ```
    $ git clone https://github.com/W0n9/BUCTNet-Login && cd BUCTNet-Login
    ```
    2. 
    ```bash
    $ make all
    ```
    3. 编译好的可执行文件在`bin`文件夹中

## 如何使用
### 查看帮助

```
$ BUCTNet-Login -h
```

### 配置登录程序使用的账户和密码

```
$ BUCTNet-Login config
```

### 登录

```
$ BUCTNet-Login login
```

### 登出

```
$ BUCTNet-Login logout
```

### 查看信息
```
$ BUCTNet-Login info
```

### 保持登录，并每隔5秒检测状态
```
$ BUCTNet-Login keepalive -i 5
```

### About

主要功能与原理

- 本地保存账号到`$HOME/.BUCTNet-Login/account.json`
- 使用账号快速登录校园网，环境支持的情况下也可以一键登录

### Thanks to
- [vouv/srun](https://github.com/vouv/srun)提供的hash算法