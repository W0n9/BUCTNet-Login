# BUCTNet-Login

> An efficient client for BUCT campus network

北京化工大学校园网命令行登录工具
- 支持linux、maxOS、windows
- 实验性支持龙芯、部分路由器平台，如MT7621
- 基于Go语言实现

## 如何获取

1. 在[Relesae](https://github.com/W0n9/BUCTNet-Login/releases/latest)处获取

2. 自行编译，要求有Golang环境 `go version >= 1.16`   
    1. 先克隆项目
    ```
    $ git clone https://github.com/W0n9/BUCTNet-Login && cd BUCTNet-Login
    ```
    2. 
    ```bash
    $ make all
    ```
    3. 编译好的可执行文件在bin文件夹中

## 如何使用
### 查看帮助

```
$ BUCTNet-Login -h
```

### 配置登录程序

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

### About

主要功能与原理

- 本地保存账号到`$HOME/.BUCTNet-Login/account.json`
- 使用账号快速登录校园网，环境支持的情况下也可以一键登录