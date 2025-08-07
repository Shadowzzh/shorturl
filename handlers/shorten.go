package handlers

import (
	"net/http"
	"short-url/config"
	"short-url/models"
	"short-url/services"

	"github.com/gin-gonic/gin"
)

func CreateShortURL(c *gin.Context) {
	var req models.ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid request",
		})
		return
	}

	shortURL, err := services.CreateShortURL(req.URL)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to create short URL",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Short URL created successfully",
		"data": gin.H{
			"id":        shortURL.ID,
			"short_url": config.AppConfig.Server.Domain + shortURL.ID,
		},
	})
}
