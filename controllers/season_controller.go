package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/pkg/logger"
	"go-gin-starter/pkg/storage"
	"go-gin-starter/services"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SeasonController handles season-related HTTP requests
type SeasonController struct {
	seasonService services.SeasonService
}

// NewSeasonController creates a new instance of SeasonController
func NewSeasonController(seasonService services.SeasonService) *SeasonController {
	return &SeasonController{
		seasonService: seasonService,
	}
}

// CreateSeason handler for creating a new season
func (c *SeasonController) CreateSeason(ctx *gin.Context) {
	var input dto.CreateSeasonInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	season, err := c.seasonService.CreateSeason(&input)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusCreated, season, constants.MsgSeasonCreated)
}

// GetAllSeasons handler for getting all seasons
func (c *SeasonController) GetAllSeasons(ctx *gin.Context) {
	seasons, err := c.seasonService.GetAllSeasons()
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}
	httpPkg.RespondSuccess(ctx, http.StatusOK, seasons, constants.MsgSeasonsFetched)
}

// GetSeasonByID handler for getting a season by ID
func (c *SeasonController) GetSeasonByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	season, err := c.seasonService.GetSeasonByID(id)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusNotFound, constants.ErrSeasonNotFound)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, season, constants.MsgSeasonFetched)
}

// UpdateSeason handler for updating a season
func (c *SeasonController) UpdateSeason(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	var input dto.UpdateSeasonInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	season, err := c.seasonService.UpdateSeason(id, &input)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, season, constants.MsgSeasonUpdated)
}

// DeleteSeason handler for deleting a season
func (c *SeasonController) DeleteSeason(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	if err := c.seasonService.DeleteSeason(id); err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgSeasonDeleted)
}

// UploadSeasonLogo handler for uploading a season logo
func (c *SeasonController) UploadSeasonLogo(ctx *gin.Context) {
	// Set maximum request size to prevent memory issues
	maxSize := int64(5 * 1024 * 1024) // 5MB limit
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize)

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	fileHeader, err := ctx.FormFile("logo")
	if err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrLogoTooLarge)
			return
		}
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrFileUploadRequired)
		return
	}

	if fileHeader.Size > 2*1024*1024 {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrLogoTooLarge)
		return
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidFileType)
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		logger.Error("Failed to open uploaded file", zap.Error(err))
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}
	defer src.Close()

	newFileName := uuid.New().String() + ext
	key := filepath.Join("logos/seasons", newFileName)
	contentType := fileHeader.Header.Get("Content-Type")

	// If content type is missing, try to infer it
	if contentType == "" {
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		default:
			contentType = "application/octet-stream"
		}
	}

	logoURL, err := storage.UploadFileToS3(src, key, contentType)
	if err != nil {
		logger.Error("S3 upload failed", zap.Error(err), zap.String("key", key))
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}

	if err := c.seasonService.UpdateSeasonLogo(id, logoURL); err != nil {
		logger.Error("Failed to update season logo in database", zap.Error(err), zap.String("seasonID", id.String()))
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, gin.H{
		"logo_url": logoURL,
	}, constants.MsgLogoUploaded)
}
