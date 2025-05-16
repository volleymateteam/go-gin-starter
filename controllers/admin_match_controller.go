// TEMPORARY PLACEHOLDER FILE - DO NOT USE
// This file will be removed once we fully transition to the new architecture
package controllers

import (
	httpPkg "go-gin-starter/pkg/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateMatch handles POST /api/admin/matches
func CreateMatch(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}

// GetAllMatches handles GET /api/admin/matches
func GetAllMatches(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}

// GetMatchByID handles GET /api/admin/matches/:id
func GetMatchByID(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}

// UpdateMatch handles PUT /api/admin/matches/:id
func UpdateMatch(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}

// DeleteMatch handles DELETE /api/admin/matches/:id
func DeleteMatch(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}

// UploadMatchVideo handles PATCH /api/admin/matches/:id/upload-video
func UploadMatchVideo(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}

// UploadMatchScout handles PATCH /api/admin/matches/:id/upload-scout
func UploadMatchScout(c *gin.Context) {
	httpPkg.RespondError(c, http.StatusNotImplemented, "This endpoint has been moved to the new architecture")
}
