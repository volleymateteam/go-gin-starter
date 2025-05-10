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

// UpdateUserByAdmin updates a user's profile by Admin
func UpdateUserByAdmin(c *gin.Context) {
	var input dto.AdminUpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	// Validate UUID
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	// Fetch original user BEFORE updating
	originalUser, err := services.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	// Validate input using helper
	if err := validationPkg.ValidateAdminUpdateInput(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Call service to update user
	updatedUser, err := services.AdminUpdateUser(userID, &input)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	metadata := auditPkg.BuildUserUpdateMetadata(originalUser, &input)

	// Prepare audit logging
	adminID := c.MustGet("user_id").(uuid.UUID)
	errLog := services.LogAdminAction(adminID, "update_user", &userID, nil, nil, nil, metadata)
	if errLog != nil {
		fmt.Printf("LogAdminAction failed: %v\n", errLog)
	}

	response := httpPkg.BuildAdminUserResponse(updatedUser)

	httpPkg.RespondSuccess(c, http.StatusOK, response, constants.MsgUserUpdated)
}

// DeleteUserByAdmin deletes any user by ID (Admin only)
func DeleteUserByAdmin(c *gin.Context) {
	// Extract target user ID
	idParam := c.Param("id")
	targetUserID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	// Fetch user before deletion
	targetUser, err := services.GetUserByID(targetUserID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	// Delete user
	err = services.DeleteUserByID(targetUserID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}

	metadata := auditPkg.BuildUserDeleteMetadata(targetUser)

	// Add audit logging
	adminID := c.MustGet("user_id").(uuid.UUID)
	_ = services.LogAdminAction(adminID, "delete_user", &targetUserID, nil, nil, nil, metadata)

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgUserDeleted)
}

// UpdateUserPermissions updates user permissions by Admin
func UpdateUserPermissions(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	var input dto.UpdatePermissionsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, fmt.Sprintf("%s: %v", constants.ErrInvalidInput, err))
		return
	}

	err = services.UpdateUserPermissions(userID, input.Permissions)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update permissions: %v", err))
		return
	}

	// Get the updated user to include username/email in log
	user, err := services.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgUserPermissionsUpdated)
		return
	}

	// Build metadata with username + email
	metadata := auditPkg.BuildUserPermissionUpdateMetadata(user, input.Permissions)

	adminID := c.MustGet("user_id").(uuid.UUID)

	// Log admin action
	errLog := services.LogAdminAction(adminID, "update_permissions", &userID, nil, nil, nil, metadata)
	if errLog != nil {
		fmt.Printf("LogAdminAction failed: %v\n", errLog)
	}

	response := httpPkg.BuildUserPermissionsResponse(user)
	httpPkg.RespondSuccess(c, http.StatusOK, response, constants.MsgUserPermissionsUpdated)
}

// GetUserPermissions retrieves a user's permissions
func GetUserPermissions(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	response := httpPkg.BuildUserPermissionsResponse(user)
	httpPkg.RespondSuccess(c, http.StatusOK, response, constants.MsgUserPermissionsFetched)
}

// ResetUserPermissions resets a user's extra permissions, keeping only their role-based permissions
func ResetUserPermissions(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	// Reset to an empty array - this only affects extra permissions
	// The role-based permissions will still be retained through the HasPermission function
	emptyPermissions := make([]string, 0)
	err = services.UpdateUserPermissions(userID, emptyPermissions)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to reset permissions: %v", err))
		return
	}

	metadata := auditPkg.BuildUserResetPermissionsMetadata(user)

	// Add audit logging
	adminID := c.MustGet("user_id").(uuid.UUID)
	_ = services.LogAdminAction(adminID, "reset_permissions", &userID, nil, nil, nil, metadata)

	response := httpPkg.BuildUserResetPermissionsResponse(user, emptyPermissions)
	httpPkg.RespondSuccess(c, http.StatusOK, response, constants.MsgUserPermissionsReset)
}

// GetAuditLogs handles GET /api/admin/audit-logs with optional filters and pagination
func GetAuditLogs(c *gin.Context) {
	// Read query parameters
	actionType := c.Query("action_type")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

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
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrFetchAuditFaild)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, logs, constants.MsgAuditLogsFetched)
}
