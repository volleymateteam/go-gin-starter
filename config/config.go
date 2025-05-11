package config

import (
	"os"
)

// Global AWS config variables
var (
	AWSRegion             string
	AWSBucketName         string
	ScoutCloudFrontDomain string
	VideoCloudFrontDomain string
)

// InitConfig initializes all config values after LoadEnv is called
func InitConfig() {
	AWSRegion = os.Getenv("AWS_REGION")
	AWSBucketName = os.Getenv("AWS_BUCKET_NAME")
	ScoutCloudFrontDomain = os.Getenv("SCOUT_CLOUDFRONT_DOMAIN")
	VideoCloudFrontDomain = os.Getenv("VIDEO_CLOUDFRONT_DOMAIN")
}

var PythonParserURL = os.Getenv("PYTHON_SCOUT_PARSER_URL")
