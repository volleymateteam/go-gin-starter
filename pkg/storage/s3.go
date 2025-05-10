package storage

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"go-gin-starter/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadMatchVideoToS3 uploads a video to S3 in the correct folder structure
func UploadMatchVideoToS3(uploader *s3manager.Uploader, file multipart.File, fileHeader *multipart.FileHeader, matchID string, seasonYear string, country string, competition string, gender string) (string, error) {
	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".mp4" && ext != ".mov" && ext != ".mkv" {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	safeCompetition := strings.ReplaceAll(strings.ToLower(competition), " ", "_")
	safeCountry := strings.ToLower(country)
	safeGender := strings.ToLower(gender)
	safeSeason := strings.ReplaceAll(seasonYear, "/", "-") // e.g., 2024-2025

	key := fmt.Sprintf("videos/%s_%s/%s_%s/%s%s", safeSeason, safeCountry, safeCompetition, safeGender, matchID, ext)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.AWSBucketName),
		Key:    aws.String(key),
		Body:   file,
		// ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}
