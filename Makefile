# ==========================================
# shorturl - 构建 / 运行入口
# ==========================================
# 后端：
#   make up / down / logs / status   用根目录 compose 起停完整环境
#   make dev / build / test / lint   开发
#   make start / stop                不用 compose 的单容器调试
#
# Raycast 扩展（需要 macOS + Node）：
#   make ext-install / ext-dev / ext-build / ext-lint / ext-test
# ==========================================

EXT_DIR        := raycast-extension

.PHONY: help dev build test lint up down logs status start stop clean \
        ext-install ext-dev ext-build ext-lint ext-test

help:
	@echo "后端："
	@echo "  up           用 docker compose 起完整环境（本地构建镜像）"
	@echo "  down         停止并移除容器（保留数据卷）"
	@echo "  logs         跟踪服务日志"
	@echo "  status       查看当前 compose 环境状态"
	@echo "  dev          热重载（需要 air）"
	@echo "  build        编译到 tmp/short-url"
	@echo "  test         go test ./..."
	@echo "  lint         go fmt + go vet"
	@echo "  start/stop   不用 compose 的单容器调试"
	@echo ""
	@echo "Raycast 扩展："
	@echo "  ext-install  安装依赖"
	@echo "  ext-dev      ray develop"
	@echo "  ext-build    ray build"
	@echo "  ext-lint     eslint"
	@echo "  ext-test     纯函数用例（不需要 Raycast GUI）"

# ---------- 启停 ----------

up:
	docker compose up -d --build
	@echo ""
	@echo "已启动，验证："
	@echo "  curl -s http://localhost:$${SHORTURL_PORT:-8087}/ping"

down:
	docker compose down
	@echo "已停止（数据卷保留；彻底删除用 docker compose down -v）"

logs:
	docker compose logs -f short-url

status:
	@docker compose ps 2>/dev/null || true

# ---------- 开发 ----------

dev:
	go mod tidy
	air

build:
	go mod tidy
	go build -o tmp/short-url main.go
	@echo "构建完成：tmp/short-url"

test:
	go test ./...
	@echo "测试完成"

lint:
	go fmt ./...
	go vet ./...
	@echo "代码检查完成"

# 不依赖 compose 的单容器调试
start:
	docker build -t shorturl:local .
	-docker stop short-url
	-docker rm short-url
	docker run -itd \
		-p 8087:8087 \
		--name short-url \
		--restart unless-stopped \
		-v shorturl-data:/data \
		-e SERVER_DOMAIN="http://localhost:8087/" \
		shorturl:local
	@echo "已启动：http://localhost:8087（数据在命名卷 shorturl-data）"

stop:
	-docker stop short-url
	-docker rm short-url
	@echo "调试容器已停止"

clean:
	-docker stop short-url
	-docker rm short-url
	docker image prune -f
	docker container prune -f
	docker network prune -f
	@echo "本地 Docker 资源已清理"

# ---------- Raycast 扩展 ----------

ext-install:
	cd $(EXT_DIR) && pnpm install

ext-dev:
	cd $(EXT_DIR) && pnpm run dev

ext-build:
	cd $(EXT_DIR) && pnpm run build

ext-lint:
	cd $(EXT_DIR) && pnpm run lint

# 扩展的纯函数用例测试：编译 src 下的纯函数，再用 node 跑断言（不需要 Raycast GUI）
ext-test:
	cd $(EXT_DIR) && ./node_modules/.bin/tsc src/validate-url.ts src/api.ts \
		--target ES2022 --module commonjs --outDir /tmp/shorturl-ext-test --skipLibCheck
	cd $(EXT_DIR) && node tests/validate-url.test.js /tmp/shorturl-ext-test
	cd $(EXT_DIR) && node tests/api.test.js /tmp/shorturl-ext-test
