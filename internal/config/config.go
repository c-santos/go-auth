package config

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port: getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET_KEY", "secret"),
	}
}

func getEnv(key string, defaultValue) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
