package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	Domain  string `mapstructure:"domain"`
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

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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

	// 调试输出
	log.Printf("Final config - Database DSN: %s", AppConfig.Database.DSN)
	log.Printf("Viper get database.dsn: %s", viper.GetString("database.dsn"))
	log.Printf("Viper get DATABASE_DSN: %s", viper.GetString("DATABASE_DSN"))
}

func setDefaults() {
	viper.SetDefault("server.port", "3001")
	viper.SetDefault("server.gin_mode", "debug")

	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.dsn", "host=localhost user=postgres password=postgres dbname=shorturl port=5432 sslmode=disable TimeZone=Asia/Shanghai")

	// redis
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
}
