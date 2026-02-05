package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBUrl     string
	JWTSecret string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found")
	}

	cfg := &Config{
		Port:      getEnv("PORT", ":5050"),
		DBUrl:     getEnv("DB_DSN", ""),
		JWTSecret: getEnv("JWT_SECRET", "default_secret_key_for_dev"),
	}

	if cfg.DBUrl == "" {
		log.Fatal("DB_URL is not set in .env")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
