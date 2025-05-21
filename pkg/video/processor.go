package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go-gin-starter/pkg/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

// VideoProcessor handles video compression and processing
type VideoProcessor struct {
	s3Client   *s3.S3
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
	bucket     string
}

// NewVideoProcessor creates a new video processor instance
func NewVideoProcessor(sess *session.Session, s3Client *s3.S3, bucket string) *VideoProcessor {
	return &VideoProcessor{
		s3Client:   s3Client,
		downloader: s3manager.NewDownloader(sess),
		uploader:   s3manager.NewUploader(sess),
		bucket:     bucket,
	}
}

// ProcessVideo handles the complete video processing pipeline
func (p *VideoProcessor) ProcessVideo(job *VideoProcessingJob) error {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "video-processing-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download raw video
	inputPath := filepath.Join(tempDir, "input"+filepath.Ext(job.InputKey))
	if err := p.downloadVideo(job.InputKey, inputPath); err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}

	// Process for each format
	for format, specs := range DefaultVideoFormats {
		outputPath := filepath.Join(tempDir, fmt.Sprintf("output_%s.mp4", format))

		if err := p.compressVideo(inputPath, outputPath, specs); err != nil {
			logger.Error("Failed to process video format",
				zap.String("format", format),
				zap.Error(err))
			continue
		}

		// Generate output key for this format
		formatKey := strings.Replace(job.OutputKey, "compressed/", fmt.Sprintf("compressed/%s/", format), 1)

		// Upload processed video
		if err := p.uploadVideo(outputPath, formatKey); err != nil {
			logger.Error("Failed to upload processed video",
				zap.String("format", format),
				zap.Error(err))
			continue
		}
	}

	// Generate and upload thumbnail
	thumbnailPath := filepath.Join(tempDir, "thumbnail.jpg")
	if err := p.generateThumbnail(inputPath, thumbnailPath); err != nil {
		logger.Error("Failed to generate thumbnail", zap.Error(err))
	} else {
		thumbnailKey := strings.Replace(job.OutputKey, "compressed/", "thumbnails/", 1)
		thumbnailKey = strings.TrimSuffix(thumbnailKey, filepath.Ext(thumbnailKey)) + ".jpg"

		if err := p.uploadVideo(thumbnailPath, thumbnailKey); err != nil {
			logger.Error("Failed to upload thumbnail", zap.Error(err))
		}
	}

	return nil
}

// downloadVideo downloads a video from S3
func (p *VideoProcessor) downloadVideo(key, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = p.downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	return err
}

// compressVideo compresses a video using ffmpeg
func (p *VideoProcessor) compressVideo(inputPath, outputPath string, format VideoFormat) error {
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c:v", "libx264",
		"-preset", "medium",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "128k",
		"-movflags", "+faststart",
		"-vf", fmt.Sprintf("scale=%s", format.Resolution),
		"-b:v", format.Bitrate,
		"-y",
		outputPath,
	)

	return cmd.Run()
}

// generateThumbnail generates a thumbnail from the video
func (p *VideoProcessor) generateThumbnail(videoPath, thumbnailPath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-vf", "scale=640:-1",
		"-y",
		thumbnailPath,
	)

	return cmd.Run()
}

// uploadVideo uploads a processed video to S3
func (p *VideoProcessor) uploadVideo(filePath, key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	contentType := p.getContentType(filePath)
	_, err = p.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	return err
}

// getContentType determines the content type based on file extension
func (p *VideoProcessor) getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}
