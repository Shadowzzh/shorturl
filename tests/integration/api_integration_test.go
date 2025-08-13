package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"short-url/config"
	"short-url/database"
	"short-url/handlers"
	"short-url/models"
)

// 集成测试 - 测试完整的请求响应流程，包括真实的数据库操作
func TestMain(m *testing.M) {
	// 设置测试环境
	os.Setenv("DATABASE_DSN", ":memory:")
	os.Setenv("SERVER_GIN_MODE", "test")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")

	// 初始化配置和数据库
	config.Init()
	database.Init()

	// 尝试初始化 Redis，但不要求必须成功（测试环境可能没有 Redis）
	database.InitRedis()

	// 运行测试
	code := m.Run()
	os.Exit(code)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 使用真实的 handlers
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.POST("/shorten", handlers.CreateShortURL)
	router.GET("/:id", handlers.RedirectToOrigin)

	return router
}

func TestPingIntegration(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ping", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestCreateShortURLIntegration(t *testing.T) {
	router := setupRouter()

	t.Run("成功创建短网址", func(t *testing.T) {
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
		assert.Equal(t, "Short URL created successfully", response["msg"])

		// 验证数据已保存到数据库
		data := response["data"].(map[string]interface{})
		shortID := data["id"].(string)

		var shortURL models.ShortURL
		err = database.DB.Where("id = ?", shortID).First(&shortURL).Error
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", shortURL.OriginalURL)
	})

	t.Run("无效请求 - 缺少URL", func(t *testing.T) {
		reqBody := models.ShortenRequest{}
		jsonData, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRedirectIntegration(t *testing.T) {
	router := setupRouter()

	// 先创建一个短网址
	shortURL := models.ShortURL{
		ID:          "test123",
		OriginalURL: "https://example.com",
	}
	database.DB.Create(&shortURL)

	t.Run("成功重定向", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test123", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Location"))
	})

	t.Run("短网址不存在", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/notfound", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// 清理测试数据
	database.DB.Delete(&shortURL)
}
