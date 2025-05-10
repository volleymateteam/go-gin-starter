package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	godotenv.Load(".env")
}

// GetEnvWithDefault returns the value of the environment variable or a default value if not set
func GetEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// GetRequiredEnv returns the value of the environment variable or panics if not set
func GetRequiredEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	panic("Required environment variable not set: " + key)
}
