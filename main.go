package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"short-url/database"
	"short-url/handlers"
)

func main() {
	database.Init()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/shorten", handlers.CreateShortURL)
	r.GET("/:id", handlers.RedirectToOrigin)

	r.Run(":8080")
}
