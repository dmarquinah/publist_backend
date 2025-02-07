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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading the .env file")
	}

	return &Config{
		ServerAddress: ":5000",
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
