package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-gin-starter/dto"
	auditPkg "go-gin-starter/pkg/audit"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	validationPkg "go-gin-starter/pkg/validation"
	"go-gin-starter/services"
)

// AdminUserController handles admin-specific user operations
type AdminUserController struct {
	userService services.UserService
}

// NewAdminUserController creates a new admin user controller
func NewAdminUserController(userService services.UserService) *AdminUserController {
	return &AdminUserController{
		userService: userService,
	}
}

// UpdateUserByAdmin updates a user's profile by Admin
func (c *AdminUserController) UpdateUserByAdmin(ctx *gin.Context) {
	var input dto.AdminUpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	// Validate UUID
	idParam := ctx.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	// Fetch original user BEFORE updating
	originalUser, err := c.userService.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	// Validate input using helper
	if err := validationPkg.ValidateAdminUpdateInput(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Call service to update user
	updatedUser, err := c.userService.AdminUpdateUser(userID, &input)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	metadata := auditPkg.BuildUserUpdateMetadata(originalUser, &input)

	// Prepare audit logging
	adminID := ctx.MustGet("user_id").(uuid.UUID)
	errLog := services.LogAdminAction(adminID, "update_user", &userID, nil, nil, nil, metadata)
	if errLog != nil {
		fmt.Printf("LogAdminAction failed: %v\n", errLog)
	}

	response := httpPkg.BuildAdminUserResponse(updatedUser)

	httpPkg.RespondSuccess(ctx, http.StatusOK, response, constants.MsgUserUpdated)
}

// DeleteUserByAdmin deletes any user by ID (Admin only)
func (c *AdminUserController) DeleteUserByAdmin(ctx *gin.Context) {
	// Extract target user ID
	idParam := ctx.Param("id")
	targetUserID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	// Fetch user before deletion
	targetUser, err := c.userService.GetUserByID(targetUserID)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	// Delete user
	err = c.userService.DeleteUserByID(targetUserID)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}

	metadata := auditPkg.BuildUserDeleteMetadata(targetUser)

	// Add audit logging
	adminID := ctx.MustGet("user_id").(uuid.UUID)
	_ = services.LogAdminAction(adminID, "delete_user", &targetUserID, nil, nil, nil, metadata)

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgUserDeleted)
}

// UpdateUserPermissions updates user permissions by Admin
func (c *AdminUserController) UpdateUserPermissions(ctx *gin.Context) {
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
func (c *AdminUserController) GetUserPermissions(ctx *gin.Context) {
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
func (c *AdminUserController) ResetUserPermissions(ctx *gin.Context) {
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

// GetAuditLogs handles GET /api/admin/audit-logs with optional filters and pagination
func (*AdminUserController) GetAuditLogs(ctx *gin.Context) {
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
