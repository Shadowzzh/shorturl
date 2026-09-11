package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"short-url/database"
	"short-url/models"
	"short-url/utils"
	"time"

	"gorm.io/gorm"
)

type RawGeo struct {
	As          string  `json:"as"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Isp         string  `json:"isp"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Org         string  `json:"org"`
	Query       string  `json:"query"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	Status      string  `json:"status"`
	Timezone    string  `json:"timezone"`
	Zip         string  `json:"zip"`
}

func GetOrCreateGeoLocation(ip string) (*models.GeoLocation, error) {
	resp, err := http.Get("http://ip-api.com/json/" + ip + "?lang=zh-CN")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var raw RawGeo

	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		return nil, err
	}

	if raw.Status != "success" {
		return nil, fmt.Errorf("geo location lookup failed for IP: %s", ip)
	}

	geoLocation := models.GeoLocation{
		ID:          utils.SnowflakeNode.Generate().Int64(),
		Ip:          raw.Query,
		Country:     raw.Country,
		CountryCode: raw.CountryCode,
		Region:      raw.Region,
		RegionName:  raw.RegionName,
		City:        raw.City,
		Lat:         raw.Lat,
		Lon:         raw.Lon,
		Org:         raw.Org,
		Isp:         raw.Isp,
		As:          raw.As,
		Timezone:    raw.Timezone,
	}

	var existing models.GeoLocation
	result := database.DB.Where("ip = ?", geoLocation.Ip).First(&existing)
	if result.Error == nil {
		return &existing, nil
	}

	if err := database.DB.Create(&geoLocation).Error; err != nil {
		return nil, err
	}

	return &geoLocation, nil
}

func GetOrCreateSession(ip string, userAgent string) (*models.Session, error) {
	var existing models.Session
	result := database.DB.Where("ip = ?", ip).First(&existing)
	if result.Error == nil {
		updates := map[string]interface{}{
			"visit_count": gorm.Expr("visit_count + 1"),
			"last_seen":   time.Now(),
		}
		database.DB.Model(&existing).Updates(updates)
		return &existing, nil
	}

	session := &models.Session{
		ID:         utils.SnowflakeNode.Generate().Int64(),
		IP:         ip,
		UserAgent:  userAgent,
		VisitCount: 1,
		LastSeen:   time.Now(),
		FirstSeen:  time.Now(),
	}

	if err := database.DB.Create(&session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func GetOrCreateVisitRecord(sessionID int64, shortURLID string, geoLocationId *int64) (*models.VisitRecord, error) {
	visitRecord := &models.VisitRecord{
		ID:            utils.SnowflakeNode.Generate().Int64(),
		SessionID:     sessionID,
		ShortURLID:    shortURLID,
		GeoLocationId: geoLocationId,
		VisitTime:     time.Now(),
	}

	if err := database.DB.Create(&visitRecord).Error; err != nil {
		return nil, err
	}

	return visitRecord, nil
}
