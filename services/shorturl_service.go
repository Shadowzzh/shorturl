package services

import (
	"context"
	"encoding/json"
	"short-url/database"
	"short-url/models"
	"short-url/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

func CreateShortURL(originUrl string) (*models.ShortURL, error) {
	if !strings.HasPrefix(originUrl, "http://") &&
		!strings.HasPrefix(originUrl, "https://") {
		originUrl = "http://" + originUrl
	}

	var existing models.ShortURL
	result := database.DB.Where("original_url=?", originUrl).First(&existing)
	if result.Error == nil {
		return &existing, nil
	}

	shortURL := models.ShortURL{
		ID:          utils.GenerateShortID(),
		VisitCount:  0,
		OriginalURL: originUrl,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&shortURL).Error; err != nil {
		return nil, err
	}

	return &shortURL, nil
}

func GetShortURL(id string) (*models.ShortURL, error) {
	var shortURL models.ShortURL

	if err := database.DB.First(&shortURL, "id=?", id).Error; err != nil {
		return nil, err
	}

	return &shortURL, nil
}

func GetFromCache(id string) (*models.ShortURL, error) {
	ctx := context.Background()

	jsonStr, err := database.RedisClient.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}

	var shortURL models.ShortURL
	err = json.Unmarshal([]byte(jsonStr), &shortURL)

	return &shortURL, err
}

func CacheShortURL(shortURL *models.ShortURL) error {
	ctx := context.Background()

	jsonData, err := json.Marshal(shortURL)
	if err != nil {
		return err
	}

	err = database.RedisClient.Set(ctx, shortURL.ID, jsonData, 24*time.Hour).Err()
	return err
}

func UpdateVisitStats(id string) error {
	now := time.Now()
	return database.DB.Model(&models.ShortURL{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"visit_count": gorm.Expr("visit_count + ?", 1),
			"last_visit":  &now,
		}).Error
}
