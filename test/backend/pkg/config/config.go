package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL     string
	JWTSecret string
	Port      string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		DBURL:     getEnv("DATABASE_URL", "postgres://etemademan:password123@localhost:5432/etemademan_db?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		Port:      getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
