package handlers

import (
	"log"
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
		// 对外只返回通用错误（不泄露内部细节），但把真实原因写进日志便于排障。
		// 历史上这里既没记日志也没带原因，导致 500 只能靠间接手段定位。
		log.Printf("create short url failed: %v (url=%q)", err, req.URL)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to create short URL",
		})
		return
	}

	// 缓存是可选能力：写失败不影响创建结果，但留个记录便于发现 Redis 异常
	if err := services.CacheShortURL(shortURL); err != nil {
		log.Printf("cache short url failed (ignored): %v (code=%q)", err, shortURL.Code)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Short URL created successfully",
		"data": gin.H{
			"id":        shortURL.ID,
			"short_url": config.AppConfig.Server.Domain + shortURL.Code,
		},
	})
}
