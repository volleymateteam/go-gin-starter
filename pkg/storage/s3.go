package storage

import (
	"fmt"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"

	"go-gin-starter/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadMatchVideoToS3 uploads a video to S3 in the correct folder structure
func UploadMatchVideoToS3(
	uploader *s3manager.Uploader,
	file multipart.File,
	fileHeader *multipart.FileHeader,
	matchID string,
	seasonYear string,
	country string,
	competition string,
	gender string,
) (string, error) {
	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".mp4" && ext != ".mov" && ext != ".mkv" {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	safeCompetition := strings.ReplaceAll(strings.ToLower(competition), " ", "_")
	safeCountry := strings.ToLower(country)
	safeGender := strings.ToLower(gender)
	safeSeason := strings.ReplaceAll(seasonYear, "/", "_")

	key := fmt.Sprintf("videos/%s_%s/%s_%s/%s%s",
		safeSeason,
		safeCountry,
		safeCompetition,
		safeGender,
		matchID,
		ext)

	// Detect MIME type
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "binary/octet-stream" // fallback
	}

	// Define tage if raw video
	tags := "storage=raw" // required for lifecycle transition

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(config.AWSBucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(mimeType),
		Tagging:     aws.String(tags),
	})
	if err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("https://%s/%s", config.VideoCloudFrontDomain, key)
	return publicURL, nil
}

// UploadUserAvatar uploads user avatar image to S3 and returns the public URL
func UploadUserAvatar(uploader *s3manager.Uploader,
	file multipart.File,
	fileHeader *multipart.FileHeader,
	userID string,
) (string, error) {
	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	key := fmt.Sprintf("avatars/%s%s", userID, ext)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.AWSBucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s/%s", config.AssetCloudFrontDomain, key), nil
}
