package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	GinMode string `mapstructure:"gin_mode"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

var AppConfig *Config

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 设置默认值
	setDefaults()

	// 自动从环境变量读取
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Config file not found, using defaults and environment variables")
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	// 解析配置到结构体
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}
}

func setDefaults() {
	viper.SetDefault("server.port", "3001")
	viper.SetDefault("server.gin_mode", "debug")
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.dsn", "shorturls.db")
}
