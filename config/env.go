package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Load .env first (fallback if ENV is not set yet)
	_ = godotenv.Load(".env")

	// Get current environment from ENV variable (if already set)
	env := os.Getenv("ENV")
	envFile := ".env"

	// Use .env.dev or .env.prod depending on ENV value
	if env == "dev" {
		envFile = ".env.dev"
	} else if env == "prod" {
		envFile = ".env.prod"
	}

	if err := godotenv.Overload(envFile); err != nil {
		log.Fatalf("❌ Error loading %s file", envFile)
	}
}

// env := os.Getenv("ENV")
// envFile := ".env"

// if env == "dev" {
// 	envFile = ".env.dev"
// } else if env == "prod" {
// 	envFile = ".env.prod"
// }

// if err := godotenv.Load(envFile); err != nil {
// 	log.Fatalf("❌ Error loading %s file", envFile)
// }

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
