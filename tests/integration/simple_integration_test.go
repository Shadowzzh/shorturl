package integration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	
	"short-url/config"
	"short-url/database"
	"short-url/models"
	"short-url/services"
)

// 集成测试：测试 services 层与数据库的协作
// 需要真实的数据库连接

func TestShortURLService_Integration(t *testing.T) {
	// 设置测试数据库
	os.Setenv("DATABASE_DSN", ":memory:")
	config.Init()
	database.Init()
	
	t.Run("创建短网址完整流程", func(t *testing.T) {
		// 1. 创建短网址
		originalURL := "https://example.com"
		shortURL, err := services.CreateShortURL(originalURL)
		
		assert.NoError(t, err)
		assert.NotNil(t, shortURL)
		assert.Equal(t, originalURL, shortURL.OriginalURL)
		assert.Len(t, shortURL.ID, 6)
		
		// 2. 验证数据确实保存到数据库了
		var dbRecord models.ShortURL
		err = database.DB.Where("id = ?", shortURL.ID).First(&dbRecord).Error
		assert.NoError(t, err)
		assert.Equal(t, originalURL, dbRecord.OriginalURL)
		
		// 3. 通过 ID 查询短网址
		foundURL, err := services.GetShortURL(shortURL.ID)
		assert.NoError(t, err)
		assert.Equal(t, originalURL, foundURL.OriginalURL)
	})
	
	t.Run("重复URL应该返回相同ID", func(t *testing.T) {
		url := "https://test.com"
		
		// 第一次创建
		first, err := services.CreateShortURL(url)
		assert.NoError(t, err)
		
		// 第二次创建相同URL
		second, err := services.CreateShortURL(url)
		assert.NoError(t, err)
		
		// 应该返回相同的ID
		assert.Equal(t, first.ID, second.ID)
	})
	
	t.Run("查询不存在的ID", func(t *testing.T) {
		_, err := services.GetShortURL("notexist")
		assert.Error(t, err)
	})
}