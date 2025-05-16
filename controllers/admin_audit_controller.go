package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/services"
)

// AdminAuditController handles audit log operations
type AdminAuditController struct {
	// No service dependency yet, using global service function
}

// NewAdminAuditController creates a new admin audit controller
func NewAdminAuditController() *AdminAuditController {
	return &AdminAuditController{}
}

// GetAuditLogs handles GET /api/admin/audit-logs with optional filters and pagination
func (c *AdminAuditController) GetAuditLogs(ctx *gin.Context) {
	// Read query parameters
	actionType := ctx.Query("action_type")
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "20")

	// Parse page & limit to integers
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Call service
	logs, err := services.GetAuditLogs(actionType, offset, limit)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrFetchAuditFaild)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, logs, constants.MsgAuditLogsFetched)
}
