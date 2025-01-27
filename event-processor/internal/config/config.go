package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURL   string
	DatabaseDSN   string
	StoreApiHost  string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	Port          int
}

// TODO: check required configs
func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		RabbitMQURL:   getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		DatabaseDSN:   getEnv("DATABASE_DSN", "postgres://user:password@localhost:5432/app?sslmode=disable"),
		StoreApiHost:  getEnv("STORE_API_HOST", "http://localhost:8080"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
		Port:          getEnvAsInt("PORT", 9090),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return fallback
}
