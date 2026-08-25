package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func Load() (Config, error) {
	// .env is useful for local development.
	// In Docker/production, environment variables can be provided directly.
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")

	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	if jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	if port == "" {
		port = "8080"
	}

	return Config{
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		Port:        port,
	}, nil
}
