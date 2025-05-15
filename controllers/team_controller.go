package controllers

import (
	"fmt"
	"go-gin-starter/dto"
	"go-gin-starter/models"
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

// TeamController handles team-related HTTP requests
type TeamController struct {
	teamService services.TeamService
}

// NewTeamController creates a new instance of TeamController
func NewTeamController(teamService services.TeamService) *TeamController {
	return &TeamController{
		teamService: teamService,
	}
}

// CreateTeam handles POST /api/admin/teams
func (c *TeamController) CreateTeam(ctx *gin.Context) {
	var input dto.CreateTeamInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	team, err := c.teamService.CreateTeam(&input)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	response := httpPkg.BuildTeamResponse(team)

	httpPkg.RespondSuccess(ctx, http.StatusCreated, response, constants.MsgTeamCreated)
}

// GetAllTeams handles GET /api/admin/teams
func (c *TeamController) GetAllTeams(ctx *gin.Context) {
	teams, err := c.teamService.GetAllTeams()
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, teams, constants.MsgTeamsFetched)
}

// GetTeamByID handles GET /api/admin/teams/:id
func (c *TeamController) GetTeamByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	team, err := c.teamService.GetTeamByID(id)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusNotFound, constants.ErrTeamNotFound)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, team, constants.MsgTeamFetched)
}

// UpdateTeam handles PUT /api/admin/teams/:id
func (c *TeamController) UpdateTeam(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	var input dto.UpdateTeamInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	team, err := c.teamService.UpdateTeam(id, &input)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, team, constants.MsgTeamUpdated)
}

// DeleteTeam handles DELETE /api/admin/teams/:id
func (c *TeamController) DeleteTeam(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	// Fetch team info before deletion
	team, _ := c.teamService.GetTeamByID(id)

	if err := c.teamService.DeleteTeam(id); err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	// Build metadata for audit log
	metadata := models.JSONBMap{}
	if team != nil {
		metadata["team_name"] = team.Name
	}

	// Add audit logging
	adminID := ctx.MustGet("user_id").(uuid.UUID)
	_ = services.LogAdminAction(adminID, "delete_team", &id, nil, nil, nil, metadata)

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgTeamDeleted)
}

// UploadTeamLogo handles PATCH /api/admin/teams/:id/upload-logo
func (c *TeamController) UploadTeamLogo(ctx *gin.Context) {
	// Set maximum request size to prevent memory issues
	maxSize := int64(5 * 1024 * 1024) // 5MB limit
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize)

	teamID, err := uuid.Parse(ctx.Param("id"))
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

	objectKey := fmt.Sprintf("logos/teams/%s%s", uuid.New().String(), ext)
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

	logoURL, err := storage.UploadFileToS3(src, objectKey, contentType)
	if err != nil {
		logger.Error("S3 upload failed", zap.Error(err), zap.String("objectKey", objectKey))
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}

	// Update only the filename in DB
	if err := c.teamService.UpdateTeamLogo(teamID, filepath.Base(objectKey)); err != nil {
		logger.Error("Failed to update team logo in database", zap.Error(err), zap.String("teamID", teamID.String()))
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, gin.H{
		"logo_url": logoURL,
	}, constants.MsgLogoUploaded)
}
