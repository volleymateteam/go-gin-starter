package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-gin-starter/dto"
	auditPkg "go-gin-starter/pkg/audit"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/services"
)

// AdminUserPermissionsController handles user permissions operations
type AdminUserPermissionsController struct {
	userService services.UserService
}

// NewAdminUserPermissionsController creates a new admin user permissions controller
func NewAdminUserPermissionsController(userService services.UserService) *AdminUserPermissionsController {
	return &AdminUserPermissionsController{
		userService: userService,
	}
}

// UpdateUserPermissions updates user permissions by Admin
func (c *AdminUserPermissionsController) UpdateUserPermissions(ctx *gin.Context) {
	idParam := ctx.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	var input dto.UpdatePermissionsInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, fmt.Sprintf("%s: %v", constants.ErrInvalidInput, err))
		return
	}

	err = c.userService.UpdateUserPermissions(userID, input.Permissions)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to update permissions: %v", err))
		return
	}

	// Get the updated user to include username/email in log
	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgUserPermissionsUpdated)
		return
	}

	// Build metadata with username + email
	metadata := auditPkg.BuildUserPermissionUpdateMetadata(user, input.Permissions)

	adminID := ctx.MustGet("user_id").(uuid.UUID)

	// Log admin action
	errLog := services.LogAdminAction(adminID, "update_permissions", &userID, nil, nil, nil, metadata)
	if errLog != nil {
		fmt.Printf("LogAdminAction failed: %v\n", errLog)
	}

	response := httpPkg.BuildUserPermissionsResponse(user)
	httpPkg.RespondSuccess(ctx, http.StatusOK, response, constants.MsgUserPermissionsUpdated)
}

// GetUserPermissions retrieves a user's permissions
func (c *AdminUserPermissionsController) GetUserPermissions(ctx *gin.Context) {
	idParam := ctx.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	response := httpPkg.BuildUserPermissionsResponse(user)
	httpPkg.RespondSuccess(ctx, http.StatusOK, response, constants.MsgUserPermissionsFetched)
}

// ResetUserPermissions resets a user's extra permissions, keeping only their role-based permissions
func (c *AdminUserPermissionsController) ResetUserPermissions(ctx *gin.Context) {
	idParam := ctx.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	// Reset to an empty array - this only affects extra permissions
	// The role-based permissions will still be retained through the HasPermission function
	emptyPermissions := make([]string, 0)
	err = c.userService.UpdateUserPermissions(userID, emptyPermissions)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to reset permissions: %v", err))
		return
	}

	metadata := auditPkg.BuildUserResetPermissionsMetadata(user)

	// Add audit logging
	adminID := ctx.MustGet("user_id").(uuid.UUID)
	_ = services.LogAdminAction(adminID, "reset_permissions", &userID, nil, nil, nil, metadata)

	response := httpPkg.BuildUserResetPermissionsResponse(user, emptyPermissions)
	httpPkg.RespondSuccess(ctx, http.StatusOK, response, constants.MsgUserPermissionsReset)
}
