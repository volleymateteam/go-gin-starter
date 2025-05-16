package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/models"
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

// TeamController handles team-related HTTP requests
type TeamController struct {
	teamService   services.TeamService
	uploadService upload.FileUploadService
}

// NewTeamController creates a new instance of TeamController
func NewTeamController(teamService services.TeamService, uploadService upload.FileUploadService) *TeamController {
	return &TeamController{
		teamService:   teamService,
		uploadService: uploadService,
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
	teamID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	// Use the file upload service to validate and upload the file
	logoURL, err := c.uploadService.ValidateAndUploadFile(ctx, "logo", upload.TeamLogo, constants.MaxLogoFileSize)
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

	// Update only the filename in DB
	if err := c.teamService.UpdateTeamLogo(teamID, logoURL); err != nil {
		logger.Error("Failed to update team logo in database", zap.Error(err), zap.String("teamID", teamID.String()))
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, gin.H{
		"logo_url": logoURL,
	}, constants.MsgLogoUploaded)
}
