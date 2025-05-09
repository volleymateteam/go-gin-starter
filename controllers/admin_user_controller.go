package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-gin-starter/dto"
	"go-gin-starter/models"
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

	// Validate Gender if provided
	if input.Gender != "" && !models.IsValidGender(input.Gender) {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidGender)
		return
	}

	// Validate Role if provided
	if input.Role != "" && !models.IsValidRole(input.Role) {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidRole)
		return
	}

	// Call service
	updatedUser, err := services.AdminUpdateUser(userID, &input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Audit logging
	adminID := c.MustGet("user_id").(uuid.UUID)

	// Dynamically build list of updated fields
	updatedFields := []string{}
	if input.Username != "" {
		updatedFields = append(updatedFields, "username")
	}
	if input.Email != "" {
		updatedFields = append(updatedFields, "email")
	}
	if input.Gender != "" {
		updatedFields = append(updatedFields, "gender")
	}
	if input.Role != "" {
		updatedFields = append(updatedFields, "role")
	}

	metadata := models.JSONBMap{
		"updated_fields": updatedFields,
	}

	// Get original user data for role change audit
	if input.Role != "" {
		originalUser, err := services.GetUserByID(userID)
		if err == nil && originalUser != nil {
			metadata["old_role"] = originalUser.Role
			metadata["new_role"] = input.Role
		}
	}

	errLog := services.LogAdminAction(adminID, "update_user", &userID, nil, nil, nil, metadata)
	if errLog != nil {
		fmt.Printf("LogAdminAction failed: %v\n", errLog)
	}

	// Prepare response
	response := dto.AdminUserResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Gender:    updatedUser.Gender,
		Role:      updatedUser.Role,
		AvatarURL: "/uploads/avatars/" + updatedUser.Avatar,
		CreatedAt: updatedUser.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: updatedUser.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

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

	// Delete user
	err = services.DeleteUserByID(targetUserID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrDatabase)
		return
	}

	// Add audit logging
	adminID := c.MustGet("user_id").(uuid.UUID)
	_ = services.LogAdminAction(adminID, "delete_user", &targetUserID, nil, nil, nil, models.JSONBMap{})

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

	// add logging
	adminID := c.MustGet("user_id").(uuid.UUID)
	metadata := models.JSONBMap{
		"new_permissions": input.Permissions,
	}
	errLog := services.LogAdminAction(adminID, "update_permissions", &userID, nil, nil, nil, metadata)
	if errLog != nil {
		fmt.Printf("LogAdminAction failed: %v\n", errLog)
	}

	// Get the updated user to return in the response
	user, err := services.GetUserByID(userID)
	if err != nil {
		utils.RespondSuccess(c, http.StatusOK, nil, utils.MsgUserPermissionsUpdated)
		return
	}

	// Get role-based permissions
	rolePerms, exists := models.RolePermissions[user.Role]
	if !exists {
		rolePerms = []string{}
	}

	// Return the updated permissions in the response
	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"user_id":           user.ID,
		"username":          user.Username,
		"role":              user.Role,
		"role_permissions":  rolePerms,
		"extra_permissions": user.ExtraPermissions,
		"all_permissions":   utils.GetAllPermissions(user),
	}, utils.MsgUserPermissionsUpdated)
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

	// Get role-based permissions
	rolePerms, exists := models.RolePermissions[user.Role]
	if !exists {
		rolePerms = []string{}
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"user_id":           user.ID,
		"username":          user.Username,
		"role":              user.Role,
		"role_permissions":  rolePerms,
		"extra_permissions": user.ExtraPermissions,
		// Optionally include effective permissions (combined)
		"all_permissions": utils.GetAllPermissions(user),
	}, utils.MsgUserPermissionsFetched)
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

	// Add audit logging
	adminID := c.MustGet("user_id").(uuid.UUID)
	_ = services.LogAdminAction(adminID, "reset_permissions", &userID, nil, nil, nil, models.JSONBMap{})

	// Get role-based permissions
	rolePerms, exists := models.RolePermissions[user.Role]
	if !exists {
		rolePerms = []string{}
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"user_id":           user.ID,
		"role":              user.Role,
		"role_permissions":  rolePerms,
		"extra_permissions": emptyPermissions,
		"all_permissions":   utils.GetAllPermissions(user),
	}, utils.MsgUserPermissionsReset)
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
