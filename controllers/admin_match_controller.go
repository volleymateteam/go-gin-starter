package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateMatch handles POST /api/admin/matches
func CreateMatch(c *gin.Context) {
	var input dto.CreateMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	// Validate round
	if !models.IsValidRound(input.Round) {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidMatchRound)
		return
	}

	match, err := services.CreateMatchService(&input)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusCreated, match, constants.MsgMatchCreated)
}

// GetAllMatches handles GET /api/admin/matches
func GetAllMatches(c *gin.Context) {
	matches, err := services.GetAllMatchesService()
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}
	httpPkg.RespondSuccess(c, http.StatusOK, matches, constants.MsgMatchesFetched)
}

// GetMatchByID handles GET /api/admin/matches/:id
func GetMatchByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	match, err := services.GetMatchByIDService(id)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrMatchNotFound)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, match, constants.MsgMatchFetched)
}

// UpdateMatch handles PUT /api/admin/matches/:id
func UpdateMatch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	var input dto.UpdateMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	match, err := services.UpdateMatchService(id, &input)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, match, constants.MsgMatchUpdated)
}

// DeleteMatch handles DELETE /api/admin/matches/:id
func DeleteMatch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	if err := services.DeleteMatchService(id); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgMatchDeleted)
}

// UploadMatchVideo handles PATCH /api/admin/matches/:id/upload-video
func UploadMatchVideo(c *gin.Context) {
	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidMatchID)
		return
	}

	file, err := c.FormFile("video")
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrFileUploadRequired)
		return
	}

	src, err := file.Open()
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}
	defer src.Close()

	videoURL, err := services.UploadMatchVideoService(matchID, src, file)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, gin.H{"video_url": videoURL}, constants.MsgVideoUploaded)
}

// UploadScoutFile handles PATCH /api/admin/matches/:id/upload-scout-file
func UploadMatchScout(c *gin.Context) {
	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidMatchID)
		return
	}

	file, err := c.FormFile("scout_file")
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrFileUploadRequired)
		return
	}

	src, err := file.Open()
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}
	defer src.Close()

	jsonURL, err := services.UploadMatchScoutService(matchID, src, file)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, gin.H{"scout_url": jsonURL}, constants.MsgScoutUploaded)
}
