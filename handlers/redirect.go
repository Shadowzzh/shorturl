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

	if err == nil {
		c.Redirect(http.StatusFound, shortURL.OriginalURL)
		return
	} else {
		shortURL, err = services.GetShortURL(shortID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "Short URL not found",
			})
			return
		}
	}

	go services.UpdateVisitStats(shortID)

	c.Redirect(http.StatusFound, shortURL.OriginalURL)
}
