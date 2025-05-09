package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-gin-starter/dto"
	"go-gin-starter/services"
	"go-gin-starter/utils"
)

// UpdateUserByAdmin updates a user's profile by Admin
func UpdateUserByAdmin(c *gin.Context) {
	var input dto.AdminUpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidInput)
		return
	}

	// Validate UUID
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	// Fetch original user BEFORE updating
	originalUser, err := services.GetUserByID(userID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, utils.ErrUserNotFound)
		return
	}

	// Validate input using helper
	if err := utils.ValidateAdminUpdateInput(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Call service to update user
	updatedUser, err := services.AdminUpdateUser(userID, &input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	metadata := utils.BuildUserUpdateMetadata(originalUser, &input)

	// Prepare audit logging
	adminID := c.MustGet("user_id").(uuid.UUID)
	errLog := services.LogAdminAction(adminID, "update_user", &userID, nil, nil, nil, metadata)
	if errLog != nil {
		fmt.Printf("LogAdminAction failed: %v\n", errLog)
	}

	response := utils.BuildAdminUserResponse(updatedUser)

	utils.RespondSuccess(c, http.StatusOK, response, utils.MsgUserUpdated)
}

// DeleteUserByAdmin deletes any user by ID (Admin only)
func DeleteUserByAdmin(c *gin.Context) {
	// Extract target user ID
	idParam := c.Param("id")
	targetUserID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	// Fetch user before deletion
	targetUser, err := services.GetUserByID(targetUserID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, utils.ErrUserNotFound)
		return
	}

	// Delete user
	err = services.DeleteUserByID(targetUserID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrDatabase)
		return
	}

	metadata := utils.BuildUserDeleteMetadata(targetUser)

	// Add audit logging
	adminID := c.MustGet("user_id").(uuid.UUID)
	_ = services.LogAdminAction(adminID, "delete_user", &targetUserID, nil, nil, nil, metadata)

	utils.RespondSuccess(c, http.StatusOK, nil, utils.MsgUserDeleted)
}

// UpdateUserPermissions updates user permissions by Admin
func UpdateUserPermissions(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	var input dto.UpdatePermissionsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, fmt.Sprintf("%s: %v", utils.ErrInvalidInput, err))
		return
	}

	err = services.UpdateUserPermissions(userID, input.Permissions)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update permissions: %v", err))
		return
	}

	// Get the updated user to include username/email in log
	user, err := services.GetUserByID(userID)
	if err != nil {
		utils.RespondSuccess(c, http.StatusOK, nil, utils.MsgUserPermissionsUpdated)
		return
	}

	// Build metadata with username + email
	metadata := utils.BuildUserPermissionUpdateMetadata(user, input.Permissions)

	adminID := c.MustGet("user_id").(uuid.UUID)

	// Log admin action
	errLog := services.LogAdminAction(adminID, "update_permissions", &userID, nil, nil, nil, metadata)
	if errLog != nil {
		fmt.Printf("LogAdminAction failed: %v\n", errLog)
	}

	response := utils.BuildUserPermissionsResponse(user)
	utils.RespondSuccess(c, http.StatusOK, response, utils.MsgUserPermissionsUpdated)
}

// GetUserPermissions retrieves a user's permissions
func GetUserPermissions(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, utils.ErrUserNotFound)
		return
	}

	response := utils.BuildUserPermissionsResponse(user)
	utils.RespondSuccess(c, http.StatusOK, response, utils.MsgUserPermissionsFetched)
}

// ResetUserPermissions resets a user's extra permissions, keeping only their role-based permissions
func ResetUserPermissions(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, utils.ErrUserNotFound)
		return
	}

	// Reset to an empty array - this only affects extra permissions
	// The role-based permissions will still be retained through the HasPermission function
	emptyPermissions := make([]string, 0)
	err = services.UpdateUserPermissions(userID, emptyPermissions)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to reset permissions: %v", err))
		return
	}

	metadata := utils.BuildUserResetPermissionsMetadata(user)

	// Add audit logging
	adminID := c.MustGet("user_id").(uuid.UUID)
	_ = services.LogAdminAction(adminID, "reset_permissions", &userID, nil, nil, nil, metadata)

	response := utils.BuildUserResetPermissionsResponse(user, emptyPermissions)
	utils.RespondSuccess(c, http.StatusOK, response, utils.MsgUserPermissionsReset)
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
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrFetchAuditFaild)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, logs, utils.MsgAuditLogsFetched)
}
