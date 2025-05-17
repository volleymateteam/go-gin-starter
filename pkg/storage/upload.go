package storage

import (
	"bytes"
	"context"
	"fmt"
	"go-gin-starter/config"
	"go-gin-starter/pkg/logger"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

// UploadFileToS3 uploads a file to AWS S3 and returns the full public URL
func UploadFileToS3(file multipart.File, objectKey string, contentType string) (string, error) {
	// Create a context with timeout for the upload operation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Read file in chunks instead of loading all into memory at once
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, file)
	if err != nil {
		logger.Error("Failed to read file", zap.Error(err))
		return "", err
	}

	// Get S3 configuration
	awsRegion := os.Getenv("AWS_REGION")
	awsBucket := os.Getenv("AWS_BUCKET_NAME")

	if awsRegion == "" || awsBucket == "" {
		logger.Error("Missing AWS configuration",
			zap.String("region", awsRegion),
			zap.String("bucket", awsBucket))
		return "", fmt.Errorf("missing AWS configuration")
	}

	// Create AWS session with configuration
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		logger.Error("Failed to create AWS session", zap.Error(err))
		return "", err
	}

	s3Client := s3.New(sess)

	// Upload with context for timeout handling
	_, err = s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(awsBucket),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		logger.Error("S3 upload failed",
			zap.Error(err),
			zap.String("objectKey", objectKey),
			zap.String("contentType", contentType))
		return "", err
	}

	// Determine CloudFront domain based on object type
	var cloudFront string
	if strings.Contains(objectKey, "avatars") || strings.Contains(objectKey, "logos") || strings.Contains(objectKey, "profile") {
		cloudFront = config.AssetCloudFrontDomain
	} else {
		cloudFront = config.ScoutCloudFrontDomain
	}

	if cloudFront == "" {
		logger.Warn("CloudFront domain is empty", zap.String("objectType", objectKey))
	}

	publicURL := fmt.Sprintf("https://%s/%s", cloudFront, objectKey)
	return publicURL, nil
}

func UploadBytesToS3(data []byte, objectKey, contentType string) (string, error) {
	awsRegion := os.Getenv("AWS_REGION")
	awsBucket := os.Getenv("AWS_BUCKET_NAME")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		return "", err
	}

	s3Client := s3.New(sess)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(awsBucket),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	// Determine CloudFront domain based on object type
	var cloudFront string
	switch {
	case strings.HasPrefix(objectKey, "videos/"):
		cloudFront = config.VideoCloudFrontDomain
	case strings.HasPrefix(objectKey, "scout-files/"):
		cloudFront = config.ScoutCloudFrontDomain
	case strings.HasPrefix(objectKey, "avatars/") || strings.HasPrefix(objectKey, "logos/"):
		cloudFront = config.AssetCloudFrontDomain
	default:
		cloudFront = config.AssetCloudFrontDomain // default to asset domain
	}

	if cloudFront == "" {
		logger.Warn("CloudFront domain is empty", zap.String("objectType", objectKey))
		return "", fmt.Errorf("CloudFront domain not configured for object type: %s", objectKey)
	}

	publicURL := fmt.Sprintf("https://%s/%s", cloudFront, objectKey)
	return publicURL, nil
}
