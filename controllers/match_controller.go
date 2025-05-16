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

// MatchController handles match-related HTTP requests
type MatchController struct {
	matchService services.MatchService
}

// NewMatchController creates a new instance of MatchController
func NewMatchController(matchService services.MatchService) *MatchController {
	return &MatchController{
		matchService: matchService,
	}
}

// CreateMatch handles POST /api/admin/matches
func (c *MatchController) CreateMatch(ctx *gin.Context) {
	var input dto.CreateMatchInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	// Validate round
	if !models.IsValidRound(input.Round) {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidMatchRound)
		return
	}

	match, err := c.matchService.CreateMatch(&input)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusCreated, match, constants.MsgMatchCreated)
}

// GetAllMatches handles GET /api/admin/matches
func (c *MatchController) GetAllMatches(ctx *gin.Context) {
	matches, err := c.matchService.GetAllMatches()
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}
	httpPkg.RespondSuccess(ctx, http.StatusOK, matches, constants.MsgMatchesFetched)
}

// GetMatchByID handles GET /api/admin/matches/:id
func (c *MatchController) GetMatchByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	match, err := c.matchService.GetMatchByID(id)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusNotFound, constants.ErrMatchNotFound)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, match, constants.MsgMatchFetched)
}

// UpdateMatch handles PUT /api/admin/matches/:id
func (c *MatchController) UpdateMatch(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	var input dto.UpdateMatchInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	match, err := c.matchService.UpdateMatch(id, &input)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, match, constants.MsgMatchUpdated)
}

// DeleteMatch handles DELETE /api/admin/matches/:id
func (c *MatchController) DeleteMatch(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	if err := c.matchService.DeleteMatch(id); err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgMatchDeleted)
}

// UploadMatchVideo handles PATCH /api/admin/matches/:id/upload-video
func (c *MatchController) UploadMatchVideo(ctx *gin.Context) {
	matchID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidMatchID)
		return
	}

	file, err := ctx.FormFile("video")
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrFileUploadRequired)
		return
	}

	src, err := file.Open()
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}
	defer src.Close()

	videoURL, err := c.matchService.UploadMatchVideo(matchID, src, file)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, gin.H{"video_url": videoURL}, constants.MsgVideoUploaded)
}

// UploadMatchScout handles PATCH /api/admin/matches/:id/upload-scout
func (c *MatchController) UploadMatchScout(ctx *gin.Context) {
	matchID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidMatchID)
		return
	}

	file, err := ctx.FormFile("scout_file")
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrFileUploadRequired)
		return
	}

	src, err := file.Open()
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}
	defer src.Close()

	jsonURL, err := c.matchService.UploadMatchScout(matchID, src, file)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, gin.H{"scout_url": jsonURL}, constants.MsgScoutUploaded)
}
