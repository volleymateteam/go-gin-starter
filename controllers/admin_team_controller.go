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

	utils.RespondSuccess(c, http.StatusCreated, team, utils.MsgTeamCreated)
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

	err = services.DeleteTeamService(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, nil, utils.MsgTeamDeleted)
}

func UploadTeamLogo(c *gin.Context) {
	idParam := c.Param("id")
	teamID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	team, err := services.GetTeamByIDService(teamID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, utils.ErrTeamNotFound)
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
	savePath := filepath.Join("uploads/logos", newFileName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrUploadFailed)
		return
	}

	team.Logo = newFileName
	if err := services.UpdateTeamLogoService(teamID, newFileName); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrUploadFailed)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"logo_url": "/uploads/logos/" + newFileName,
	}, utils.MsgLogoUploaded)
}
