// TEMPORARY PLACEHOLDER FILE - DO NOT USE
// This file will be removed once we fully transition to the new architecture
package controllers

import (
	httpPkg "go-gin-starter/pkg/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateTeam handles POST /api/admin/teams
func CreateTeam(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}

// GetAllTeams handles GET /api/admin/teams
func GetAllTeams(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}

// GetTeamByID handles GET /api/admin/teams/:id
func GetTeamByID(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}

// UpdateTeam handles PUT /api/admin/teams/:id
func UpdateTeam(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}

// DeleteTeam handles DELETE /api/admin/teams/:id
func DeleteTeam(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}

// UploadTeamLogo handles PATCH /api/admin/teams/:id/upload-logo
func UploadTeamLogo(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}
