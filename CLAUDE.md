# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 项目概述

这是一个使用 Gin 框架构建的简单 Go Web 应用程序，提供基本的 URL 缩短服务和用户管理功能。应用程序由单个 `main.go` 文件组成，包含所有核心功能。

## 架构

- **单文件架构**：所有代码都包含在 `main.go` 中
- **内存存储**：使用简单的 map (`var db = make(map[string]string)`) 进行数据持久化
- **Gin Web 框架**：支持中间件的 RESTful API 端点
- **基本身份验证**：使用 Gin 内置的 BasicAuth 中间件保护路由

### 核心组件

- `setupRouter()`：配置所有路由和中间件
- 内存数据库：使用 Go map 的简单键值存储
- 三个主要端点：
  - `GET /ping`：健康检查端点
  - `GET /user/:name`：获取用户值
  - `POST /admin`：设置用户值的受保护端点（需要基本身份验证）

## 常用命令

### 开发
```bash
# 运行应用程序（热重载）
air

# 运行应用程序（不热重载）
go run main.go

# 构建应用程序
go build -o short-url main.go

# 格式化代码
go fmt main.go

# 运行测试（如果存在）
go test

# 检查潜在问题
go vet main.go
```

### 依赖管理
```bash
# 下载并安装依赖
go mod tidy

# 添加新依赖
go get <包名>

# 更新依赖
go get -u all
```

### 运行服务器
服务器默认在 8080 端口运行。访问端点：
- 健康检查：`http://localhost:8080/ping`
- 用户数据：`http://localhost:8080/user/{用户名}`
- 管理端点（受保护）：`http://localhost:8080/admin`（需要基本身份验证：foo:bar 或 manu:123）

## 开发注意事项

- 应用程序使用 Go 工作区（存在 go.work 文件）
- 所有数据存储在内存中，服务器重启时会丢失
- 基本身份验证凭据是硬编码的：`foo:bar` 和 `manu:123`
- 服务器监听所有接口 (0.0.0.0:8080)
- 已配置 Air 热重载，修改代码后会自动重启服务器
- Air 配置文件：`.air.toml`，构建输出目录：`tmp/`