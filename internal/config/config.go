// Конфигурация приложения и загрузка переменных окружения
package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	DB          DatabaseConfig
	Server      ServerConfig
	RateAPI     RateAPIConfig
	Environment string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type RateAPIConfig struct {
	URL   string
	Key   string
	CacheTTL time.Duration
}

func LoadConfig() *Config {
	return &Config{
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "flex_exchange"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  time.Second * 30,
			WriteTimeout: time.Second * 30,
		},
		RateAPI: RateAPIConfig{
			URL:   getEnv("RATE_API_URL", "https://api.exchangerate.host"),
			Key:   os.Getenv("API_KEY"),
			CacheTTL: time.Minute * 5,
		},
		Environment: getEnv("ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}