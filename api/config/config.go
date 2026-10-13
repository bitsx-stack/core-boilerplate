package config

import (
	"os"
)

type Config struct {
	ServerPort string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		ServerPort: GetEnv("PORT", ":8080"),
		JWTSecret:  GetEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
