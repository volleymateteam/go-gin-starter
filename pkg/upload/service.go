package upload

import (
	"errors"
	"go-gin-starter/pkg/constants"
	"go-gin-starter/pkg/logger"
	"go-gin-starter/pkg/storage"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FileType string

const (
	TeamLogo   FileType = "team_logo"
	SeasonLogo FileType = "season_logo"
	UserAvatar FileType = "user_avatar"
	MatchVideo FileType = "match_video"
	MatchScout FileType = "match_scout"
)

type FileUploadService interface {
	ValidateAndUploadFile(ctx *gin.Context, fileField string, fileType FileType, maxSize int64) (string, error)
}

type FileUploadServiceImpl struct{}

func NewFileUploadService() FileUploadService {
	return &FileUploadServiceImpl{}
}

func (s *FileUploadServiceImpl) ValidateAndUploadFile(ctx *gin.Context, fileField string, fileType FileType, maxSize int64) (string, error) {
	// Set maximum request size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize)

	fileHeader, err := ctx.FormFile(fileField)
	if err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			return "", errors.New(constants.ErrLogoTooLarge)
		}
		return "", errors.New(constants.ErrFileUploadRequired)
	}

	if fileHeader.Size > maxSize {
		return "", errors.New(constants.ErrLogoTooLarge)
	}

	ext := filepath.Ext(fileHeader.Filename)
	allowedExts := s.getAllowedExtensions(fileType)
	if !s.isExtensionAllowed(ext, allowedExts) {
		return "", errors.New(constants.ErrInvalidFileType)
	}

	src, err := fileHeader.Open()
	if err != nil {
		logger.Error("Failed to open uploaded file", zap.Error(err))
		return "", errors.New(constants.ErrUploadFailed)
	}
	defer src.Close()

	// Generate path based on file type
	objectKey := s.generateObjectKey(fileType, ext)

	// Determine content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = s.inferContentType(ext)
	}

	// Upload to S3
	fileURL, err := storage.UploadFileToS3(src, objectKey, contentType)
	if err != nil {
		logger.Error("S3 upload failed", zap.Error(err), zap.String("key", objectKey))
		return "", errors.New(constants.ErrUploadFailed)
	}

	return fileURL, nil
}

func (s *FileUploadServiceImpl) getAllowedExtensions(FileType FileType) []string {
	switch FileType {
	case TeamLogo, SeasonLogo, UserAvatar:
		return []string{".jpg", ".jpeg", ".png"}
	case MatchVideo:
		return []string{".mp4", ".avi", ".mov"}
	case MatchScout:
		return []string{".dvw"}
	default:
		return []string{}
	}
}

func (s *FileUploadServiceImpl) isExtensionAllowed(ext string, allowedExts []string) bool {
	ext = strings.ToLower(ext)
	for _, allowed := range allowedExts {
		if ext == allowed {
			return true
		}
	}
	return false
}

func (s *FileUploadServiceImpl) generateObjectKey(fileType FileType, ext string) string {
	filename := uuid.New().String() + ext

	switch fileType {
	case TeamLogo:
		return filepath.Join("logos/teams", filename)
	case SeasonLogo:
		return filepath.Join("logos/seasons", filename)
	case UserAvatar:
		return filepath.Join("avatars", filename)
	case MatchVideo:
		return filepath.Join("videos", filename)
	case MatchScout:
		return filepath.Join("scouts", filename)
	default:
		return filepath.Join("misc", filename)
	}
}

func (s *FileUploadServiceImpl) inferContentType(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".mkv":
		return "video/x-matroska"
	case ".json":
		return "application/json"
	case ".csv":
		return "text/csv"
	default:
		return "application/octet-stream"
	}
}
