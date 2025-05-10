package config

import (
	"os"
)

// Global AWS config variables
var (
	AWSRegion     string
	AWSBucketName string
)

// InitConfig initializes all config values after LoadEnv is called
func InitConfig() {
	AWSRegion = os.Getenv("AWS_REGION")
	AWSBucketName = os.Getenv("AWS_BUCKET_NAME")
}

var PythonParserURL = os.Getenv("PYTHON_SCOUT_PARSER_URL")
