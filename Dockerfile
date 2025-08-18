# 构建阶段
FROM golang:1.24-alpine AS builder

# 安装构建依赖
RUN apk --no-cache add gcc musl-dev sqlite-dev

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源码
COPY . .

# 编译应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates sqlite

# 创建工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/main .

# 创建数据库目录
RUN mkdir -p /data

# 设置环境变量
ENV SERVER_GIN_MODE=release
ENV SERVER_PORT=8087

# 暴露端口
EXPOSE 8087

# 运行应用
CMD ["./main"]