package tests

import (
	"os"
	"testing"

	"short-url/config"
	"short-url/database"
)

func TestMain(m *testing.M) {
	// 设置测试环境变量
	os.Setenv("DATABASE_DSN", ":memory:")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("SERVER_GIN_MODE", "test")

	// 初始化测试配置
	config.Init()

	// 初始化测试数据库
	database.Init()

	// 运行测试
	code := m.Run()

	// 清理资源
	os.Exit(code)
}