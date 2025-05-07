package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Global AWS config variables
var (
	AWSRegion     string
	AWSBucketName string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AWSRegion = os.Getenv("AWS_REGION")
	AWSBucketName = os.Getenv("AWS_BUCKET_NAME")
}

var PythonParserURL = os.Getenv("PYTHON_SCOUT_PARSER_URL")
