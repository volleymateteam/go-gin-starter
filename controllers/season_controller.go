package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/pkg/logger"
	"go-gin-starter/pkg/upload"
	"go-gin-starter/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SeasonController handles season-related HTTP requests
type SeasonController struct {
	seasonService services.SeasonService
	uploadService upload.FileUploadService
}

// NewSeasonController creates a new instance of SeasonController
func NewSeasonController(seasonService services.SeasonService, uploadService upload.FileUploadService) *SeasonController {
	return &SeasonController{
		seasonService: seasonService,
		uploadService: uploadService,
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
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	// Use the file upload service to validate and upload the file
	logoURL, err := c.uploadService.ValidateAndUploadFile(ctx, "logo", upload.SeasonLogo, constants.MaxLogoFileSize)
	if err != nil {
		if err.Error() == constants.ErrLogoTooLarge ||
			err.Error() == constants.ErrFileUploadRequired ||
			err.Error() == constants.ErrInvalidFileType {
			httpPkg.RespondError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		logger.Error("File upload failed", zap.Error(err))
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
