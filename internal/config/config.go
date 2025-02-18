package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	// Add more configuration options here
}

func New() *Config {
	//Handle load of environment params
	log.Println("Reading env variables...")
	appEnv := getEnv("APP_ENV", "local")

	// Only try to load .env file if not in Docker
	if os.Getenv("DOCKER") != "true" {
		err := godotenv.Load(".env." + appEnv)
		if err != nil {
			log.Printf("Warning: Error loading the .env file: %v", err)
			// Don't fatal here, as env vars might be set through other means
		}
	}

	port := getEnv("APP_PORT", ":5000")

	return &Config{
		ServerAddress: port,
		// Initialize other config values
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
