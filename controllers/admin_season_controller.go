package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateSeason handler for creating a new season
func CreateSeason(c *gin.Context) {
	var input dto.CreateSeasonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidInput)
		return
	}

	season, err := services.CreateSeasonService(&input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer)
		return
	}

	utils.RespondSuccess(c, http.StatusCreated, season, utils.MsgSeasonCreated)
}

// GetAllSeasons handler for getting all seasons
func GetAllSeasons(c *gin.Context) {
	seasons, err := services.GetAllSeasonsService()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrDatabase)
		return
	}
	utils.RespondSuccess(c, http.StatusOK, seasons, utils.MsgSeasonsFetched)
}

// GetSeasonByID handler for getting a season by ID
func GetSeasonByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	season, err := services.GetSeasonByIDService(id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, utils.ErrSeasonNotFound)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, season, utils.MsgSeasonFetched)
}

// UpdateSeason handler for updating a season
func UpdateSeason(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	var input dto.UpdateSeasonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidInput)
		return
	}

	season, err := services.UpdateSeasonService(id, &input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, season, utils.MsgSeasonUpdated)
}

// DeleteSeason handler for deleting a season
func DeleteSeason(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	if err := services.DeleteSeasonService(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, nil, utils.MsgSeasonDeleted)
}

// UploadSeasonLogo handler for uploading a season logo
func UploadSeasonLogo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	file, err := c.FormFile("logo")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrFileUploadRequired)
		return
	}

	if file.Size > 2*1024*1024 {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrLogoTooLarge)
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidFileType)
		return
	}

	newFileName := uuid.New().String() + ext
	savePath := filepath.Join("uploads/logos/seasons", newFileName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrUploadFailed)
		return
	}

	// Update in DB
	if err := services.UpdateSeasonLogoService(id, newFileName); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrDatabase)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"logo_url": "/uploads/logos/seasons/" + newFileName,
	}, utils.MsgLogoUploaded)
}
