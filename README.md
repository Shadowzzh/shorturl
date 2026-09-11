# shorturl

自托管短链接服务。Go + Gin 后端，SQLite 存储，可选 Redis 缓存。
附带一个 macOS Raycast 客户端，见 [`raycast-extension/`](raycast-extension/)。

## 跑起来

需要 Docker（含 compose）：

```bash
docker compose up -d --build
```

验证：

```bash
curl -s http://localhost:8087/ping
curl -s -X POST http://localhost:8087/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com/very/long/path"}'
```

```json
{"code":200,"data":{"id":2098347898403033088,"short_url":"http://localhost:8087/IXMF3t"},"msg":"Short URL created successfully"}
```

访问 `http://localhost:8087/IXMF3t` 会 **302** 跳转到原始地址。

不用容器直接跑需要 Go 1.24+：`go run .`（读 `config.yaml`，端口 `3001`）。

## 接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/shorten` | 创建短链，body `{"url":"..."}` |
| `GET` | `/<短码>` | 302 跳转到原始地址 |
| `GET` | `/ping` | 存活探针，返回 `pong` |

错误统一返回 `{"code":<非200>,"msg":"..."}`。

## 配置

写进 `.env` 即可（也可直接 `export`）——这两个是最常改的：

```bash
SHORTURL_PORT=8087                      # 对外端口
SHORTURL_DOMAIN=https://s.example.com/  # 返回的 short_url 前缀，结尾必须带 /
```

`SHORTURL_DOMAIN` 不改变实际监听地址，只决定**返回给你的短链长什么样**。

其余用环境变量覆盖（优先级高于配置文件，viper，`.` 换成 `_`）：`SERVER_PORT`、
`SERVER_GIN_MODE`、`DATABASE_DSN`、`REDIS_HOST`、`REDIS_PORT`、`REDIS_PASSWORD`、`REDIS_DB`。
容器内默认值在 `config.docker.yaml`。

**Redis 是可选的**：只做读缓存，连不上只打印一行日志，创建与跳转照常。
不需要缓存就把 `REDIS_HOST` 指向一个不存在的地址。

两个容易踩的点：

- `SERVER_DOMAIN` 必须以 `/` 结尾——返回值是 `domain + 短码` 直接拼接
- 容器里 `REDIS_HOST` 要写 `redis` 而非 `localhost`——Redis 是独立容器

## 开发

```bash
make help       # 全部命令
make test       # go test ./...
make ext-test   # Raycast 扩展的纯函数用例（不需要 GUI）
```
