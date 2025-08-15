package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"short-url/config"
	"short-url/models"
)

var DB *gorm.DB

func Init() {
	var err error

	dialector := postgres.Open(config.AppConfig.Database.DSN)

	DB, err = gorm.Open(dialector, &gorm.Config{})

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
