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

// CreateTeam handles POST /api/admin/teams
func CreateTeam(c *gin.Context) {
	var input dto.CreateTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	team, err := services.CreateTeamService(&input)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	response := httpPkg.BuildTeamResponse(team)

	httpPkg.RespondSuccess(c, http.StatusCreated, response, constants.MsgTeamCreated)
}

// GetAllTeams handles GET /api/admin/teams
func GetAllTeams(c *gin.Context) {
	teams, err := services.GetAllTeamsService()
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, teams, constants.MsgTeamsFetched)
}

// GetTeamByID handles GET /api/admin/teams/:id
func GetTeamByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	team, err := services.GetTeamByIDService(id)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrTeamNotFound)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, team, constants.MsgTeamFetched)
}

// UpdateTeam handles PUT /api/admin/teams/:id
func UpdateTeam(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	var input dto.UpdateTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	team, err := services.UpdateTeamService(id, &input)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, team, constants.MsgTeamUpdated)
}

// DeleteTeam handles DELETE /api/admin/teams/:id
func DeleteTeam(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	// Fetch team info before deletion
	team, _ := services.GetTeamByIDService(id)

	if err := services.DeleteTeamService(id); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInternalServer)
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

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgTeamDeleted)
}

// UploadTeamLogo handles PATCH /api/admin/teams/:id/upload-logo
func UploadTeamLogo(c *gin.Context) {
	// Set maximum request size to prevent memory issues
	maxSize := int64(5 * 1024 * 1024) // 5MB limit
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	fileHeader, err := c.FormFile("logo")
	if err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrLogoTooLarge)
			return
		}
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
		logger.Error("Failed to open uploaded file", zap.Error(err))
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrUploadFailed)
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
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}

	// Update only the filename in DB
	if err := services.UpdateTeamLogoService(teamID, filepath.Base(objectKey)); err != nil {
		logger.Error("Failed to update team logo in database", zap.Error(err), zap.String("teamID", teamID.String()))
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, gin.H{
		"logo_url": logoURL,
	}, constants.MsgLogoUploaded)
}
