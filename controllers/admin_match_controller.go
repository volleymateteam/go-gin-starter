package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateMatch handles POST /api/admin/matches
func CreateMatch(c *gin.Context) {
	var input dto.CreateMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidInput)
		return
	}

	match, err := services.CreateMatchService(&input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer)
		return
	}

	utils.RespondSuccess(c, http.StatusCreated, match, utils.MsgMatchCreated)
}

// GetAllMatches handles GET /api/admin/matches
func GetAllMatches(c *gin.Context) {
	matches, err := services.GetAllMatchesService()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrDatabase)
		return
	}
	utils.RespondSuccess(c, http.StatusOK, matches, utils.MsgMatchesFetched)
}

// GetMatchByID handles GET /api/admin/matches/:id
func GetMatchByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	match, err := services.GetMatchByIDService(id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, utils.ErrMatchNotFound)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, match, utils.MsgMatchFetched)
}

// UpdateMatch handles PUT /api/admin/matches/:id
func UpdateMatch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	var input dto.UpdateMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidInput)
		return
	}

	match, err := services.UpdateMatchService(id, &input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, match, utils.MsgMatchUpdated)
}

// DeleteMatch handles DELETE /api/admin/matches/:id
func DeleteMatch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	if err := services.DeleteMatchService(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, nil, utils.MsgMatchDeleted)
}

// UploadMatchVideo handles PATCH /api/admin/matches/:id/upload-video
func UploadMatchVideo(c *gin.Context) {
	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidMatchID)
		return
	}

	file, err := c.FormFile("video")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrFileUploadRequired)
		return
	}

	src, err := file.Open()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrUploadFailed)
		return
	}
	defer src.Close()

	videoURL, err := services.UploadMatchVideoService(matchID, src, file)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{"video_url": videoURL}, utils.MsgVideoUploaded)
}
