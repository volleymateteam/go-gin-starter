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
