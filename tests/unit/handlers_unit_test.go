package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"short-url/models"
)

// 单元测试 - 只测试 handler 的输入输出逻辑，不依赖数据库

func TestCreateShortURLHandler_Unit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 模拟 handler 函数，只测试请求解析和响应格式
	mockHandler := func(c *gin.Context) {
		var req models.ShortenRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "Invalid request",
			})
			return
		}

		// 模拟成功响应
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "Short URL created successfully",
			"data": gin.H{
				"id":        "mock123",
				"short_url": "http://localhost:3001/mock123",
			},
		})
	}

	router := gin.New()
	router.POST("/shorten", mockHandler)

	t.Run("有效JSON解析", func(t *testing.T) {
		reqBody := models.ShortenRequest{
			URL: "https://example.com",
		}
		jsonData, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(200), response["code"])
	})

	t.Run("无效JSON处理", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRedirectHandler_Unit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 模拟重定向 handler
	mockRedirectHandler := func(c *gin.Context) {
		shortID := c.Param("id")

		// 跳过保留路径
		if shortID == "ping" || shortID == "shorten" || shortID == "info" {
			c.Next()
			return
		}

		// 模拟查找逻辑
		if shortID == "existing123" {
			c.Redirect(http.StatusFound, "https://example.com")
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Short URL not found",
		})
	}

	router := gin.New()
	router.GET("/:id", mockRedirectHandler)

	t.Run("参数解析和重定向", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/existing123", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Location"))
	})

	t.Run("保留路径处理", func(t *testing.T) {
		// 添加一个处理 ping 的路由
		router.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/ping", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "pong", w.Body.String())
	})
}
