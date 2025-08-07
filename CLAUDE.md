# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。


使用中文回复

## 项目概述

这是一个使用 Gin 框架构建的 Go Web 应用程序，提供 URL 缩短服务。应用程序使用 SQLite 数据库存储短网址数据，支持通过环境变量进行配置。

## 架构

- **模块化架构**：代码分为 `handlers`、`database`、`services` 等包
- **SQLite 数据库**：使用 GORM 进行数据持久化
- **Gin Web 框架**：RESTful API 端点
- **环境变量配置**：支持通过环境变量配置端口等参数

### 核心组件

- `database` 包：数据库连接和模型定义
- `handlers` 包：HTTP 请求处理器
- `services` 包：业务逻辑处理
- 两个主要端点：
  - `GET /ping`：健康检查端点
  - `POST /shorten`：创建短网址
  - `GET /:id`：短网址重定向

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

### 环境变量配置
项目支持通过环境变量进行配置：

```bash
# 复制示例配置文件
cp .env.example .env

# 编辑配置文件
vim .env
```

主要环境变量：
- `PORT`：服务器端口（默认：3001）

### 运行服务器
服务器默认在 3001 端口运行。可以通过环境变量 `PORT` 修改端口：

```bash
# 使用默认端口 3001
go run main.go

# 或指定端口
PORT=8080 go run main.go
```

访问端点：
- 健康检查：`http://localhost:3001/ping`
- 创建短网址：`POST http://localhost:3001/shorten`
- 访问短网址：`http://localhost:3001/{短网址ID}`

## 开发注意事项

- 应用程序使用 Go 工作区（存在 go.work 文件）
- 数据存储在 SQLite 数据库中（`shorturls.db`）
- 服务器默认监听 3001 端口，可通过 `PORT` 环境变量修改
- 已配置 Air 热重载，修改代码后会自动重启服务器
- Air 配置文件：`.air.toml`，构建输出目录：`tmp/`
- 使用 GORM 进行数据库操作
- 短网址 ID 使用随机生成算法