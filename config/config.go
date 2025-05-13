package config

import (
	"fmt"
	"os"
)

// Global AWS config variables
var (
	AWSRegion             string
	AWSBucketName         string
	ScoutCloudFrontDomain string
	VideoCloudFrontDomain string
	AssetCloudFrontDomain string
)

// InitConfig initializes all config values after LoadEnv is called
func InitConfig() {
	AWSRegion = os.Getenv("AWS_REGION")
	AWSBucketName = os.Getenv("AWS_BUCKET_NAME")
	ScoutCloudFrontDomain = os.Getenv("SCOUT_CLOUDFRONT_DOMAIN")
	VideoCloudFrontDomain = os.Getenv("VIDEO_CLOUDFRONT_DOMAIN")
	AssetCloudFrontDomain = os.Getenv("ASSET_CLOUDFRONT_DOMAIN")

	fmt.Println("DEBUG: Using VIDEO_CLOUDFRONT_DOMAIN =", VideoCloudFrontDomain)
}

var PythonParserURL = os.Getenv("PYTHON_SCOUT_PARSER_URL")
