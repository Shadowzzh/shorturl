package main

import (
	"fmt"
	"net/http"

	"github.com/bwmarrin/snowflake"
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

	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/shorten", handlers.CreateShortURL)
	r.GET("/:id", handlers.RedirectToOrigin)

	r.Run(":" + config.AppConfig.Server.Port)
}
