package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"short-url/config"
	"short-url/models"
)

var DB *gorm.DB

func Init() {
	var err error

	DB, err = gorm.Open(sqlite.Open(config.AppConfig.Database.DSN), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	modelsToMigrate := []interface{}{
		&models.ShortURL{},
		&models.Session{},
		&models.VisitRecord{},
		&models.GeoLocation{},
	}

	if err := DB.AutoMigrate(modelsToMigrate...); err != nil {
		panic("failed to migrate models: " + err.Error())
	}

	fmt.Println("Database initialized successfully")
}
