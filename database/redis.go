package database

import (
	"context"
	"fmt"
	"short-url/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.Host + ":" + config.AppConfig.Redis.Port,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()

	if err != nil {
		fmt.Printf("Error connecting to Redis: %v", err)
	} else {
		fmt.Println("Connected to Redis successfully")
	}

}
