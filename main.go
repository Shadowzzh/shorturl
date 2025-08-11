package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"short-url/config"
	"short-url/database"
	"short-url/handlers"
)

func main() {
	// 初始化配置
	config.Init()

	// 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.GinMode)

	database.Init()
	database.InitRedis()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/shorten", handlers.CreateShortURL)
	r.GET("/:id", handlers.RedirectToOrigin)

	r.Run(":" + config.AppConfig.Server.Port)
}
