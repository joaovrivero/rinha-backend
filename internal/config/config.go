package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() (*Config, error) {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "host=db user=rinha password=rinha dbname=rinhadb port=5432 sslmode=disable"),
		Port:        getEnv("PORT", "80"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
