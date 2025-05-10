package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTeam handles POST /api/admin/teams
func CreateTeam(c *gin.Context) {
	var input dto.CreateTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidInput)
		return
	}

	team, err := services.CreateTeamService(&input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer)
		return
	}

	response := utils.BuildTeamResponse(team)

	utils.RespondSuccess(c, http.StatusCreated, response, utils.MsgTeamCreated)
}

// GetAllTeams handles GET /api/admin/teams
func GetAllTeams(c *gin.Context) {
	teams, err := services.GetAllTeamsService()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrDatabase)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, teams, utils.MsgTeamsFetched)
}

// GetTeamByID handles GET /api/admin/teams/:id
func GetTeamByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	team, err := services.GetTeamByIDService(id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, utils.ErrTeamNotFound)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, team, utils.MsgTeamFetched)
}

// UpdateTeam handles PUT /api/admin/teams/:id
func UpdateTeam(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	var input dto.UpdateTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidInput)
		return
	}

	team, err := services.UpdateTeamService(id, &input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, team, utils.MsgTeamUpdated)
}

// DeleteTeam handles DELETE /api/admin/teams/:id
func DeleteTeam(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	// Fetch team info before deletion
	team, _ := services.GetTeamByIDService(id)

	if err := services.DeleteTeamService(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer)
		return
	}

	// Build metadata for audit log
	metadata := models.JSONBMap{}
	if team != nil {
		metadata["team_name"] = team.Name
	}

	// Add audit logging
	adminID := c.MustGet("user_id").(uuid.UUID)
	_ = services.LogAdminAction(adminID, "delete_team", &id, nil, nil, nil, metadata)

	utils.RespondSuccess(c, http.StatusOK, nil, utils.MsgTeamDeleted)
}

// UploadTeamLogo handles PATCH /api/admin/teams/:id/upload-logo
func UploadTeamLogo(c *gin.Context) {
	idParam := c.Param("id")
	teamID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	file, err := c.FormFile("logo")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrFileUploadRequired)
		return
	}

	// delegate everything to service
	newFileName, savePath, err := services.UploadAndSaveTeamLogoService(teamID, file)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrFileSaveFailed)
		return
	}

	// audit log
	adminID := c.MustGet("user_id").(uuid.UUID)
	metadata := models.JSONBMap{"filename": newFileName}
	_ = services.LogAdminAction(adminID, "upload_team_logo", &teamID, nil, nil, nil, metadata)

	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"logo_url": "/uploads/logos/" + newFileName,
	}, utils.MsgLogoUploaded)
}
