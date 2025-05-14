package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/pkg/storage"
	"go-gin-starter/services"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateSeason handler for creating a new season
func CreateSeason(c *gin.Context) {
	var input dto.CreateSeasonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	season, err := services.CreateSeasonService(&input)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusCreated, season, constants.MsgSeasonCreated)
}

// GetAllSeasons handler for getting all seasons
func GetAllSeasons(c *gin.Context) {
	seasons, err := services.GetAllSeasonsService()
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}
	httpPkg.RespondSuccess(c, http.StatusOK, seasons, constants.MsgSeasonsFetched)
}

// GetSeasonByID handler for getting a season by ID
func GetSeasonByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	season, err := services.GetSeasonByIDService(id)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrSeasonNotFound)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, season, constants.MsgSeasonFetched)
}

// UpdateSeason handler for updating a season
func UpdateSeason(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	var input dto.UpdateSeasonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	season, err := services.UpdateSeasonService(id, &input)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, season, constants.MsgSeasonUpdated)
}

// DeleteSeason handler for deleting a season
func DeleteSeason(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	if err := services.DeleteSeasonService(id); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgSeasonDeleted)
}

// UploadSeasonLogo handler for uploading a season logo
func UploadSeasonLogo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	fileHeader, err := c.FormFile("logo")
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrFileUploadRequired)
		return
	}

	if fileHeader.Size > 2*1024*1024 {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrLogoTooLarge)
		return
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidFileType)
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}
	defer src.Close()

	newFileName := uuid.New().String() + ext
	key := filepath.Join("logos/seasons", newFileName)
	contentType := fileHeader.Header.Get("Content-Type")

	logoURL, err := storage.UploadFileToS3(src, key, contentType)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}

	if err := services.UpdateSeasonLogoService(id, logoURL); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, gin.H{
		"logo_url": logoURL,
	}, constants.MsgLogoUploaded)
}
