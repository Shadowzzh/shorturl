package handlers

import (
	"net/http"
	"short-url/services"

	"github.com/gin-gonic/gin"
)

func RedirectToOrigin(c *gin.Context) {
	shortID := c.Param("id")
	if shortID == "ping" || shortID == "shorten" || shortID == "info" {
		c.Next()
		return
	}

	shortURL, err := services.GetFromCache(shortID)

	if err != nil {
		shortURL, err = services.GetShortURL(shortID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "Short URL not found",
			})
			return
		}
	}

	ip := c.Copy().ClientIP()
	userAgent := c.GetHeader("User-Agent")

	geoLocation, err := services.GetOrCreateGeoLocation(ip)
	if err != nil {
		geoLocation = nil
	}

	session, err := services.GetOrCreateSession(ip, userAgent)
	if err != nil {
		// 会话创建失败，记录日志但不影响主流程
		session = nil
	}

	if geoLocation != nil && session != nil {
		services.GetOrCreateVisitRecord(session.ID, shortID, &geoLocation.ID)
	}

	go services.UpdateVisitStats(shortID)

	c.Redirect(http.StatusFound, shortURL.OriginalURL)
}
