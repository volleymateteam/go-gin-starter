package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadFileToS3 uploads a file to AWS S3 and returns the full public URL
func UploadFileToS3(file multipart.File, objectKey string, contentType string) (string, error) {
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

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
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
		// ACL:         aws.String("public-read"), // optional, makes it publicly accessible
	})
	if err != nil {
		return "", err
	}

	// publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", awsBucket, awsRegion, objectKey)
	publicURL := fmt.Sprintf("https://%s/%s", os.Getenv("CLOUDFRONT_DOMAIN"), objectKey)

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
		// ACL:         aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	// publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", awsBucket, awsRegion, objectKey)
	publicURL := fmt.Sprintf("https://%s/%s", os.Getenv("CLOUDFRONT_DOMAIN"), objectKey)
	return publicURL, nil
}
