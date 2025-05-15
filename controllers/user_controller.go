package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserController handles user-related HTTP requests
type UserController struct {
	userService services.UserService
}

// NewUserController creates a new instance of UserController
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Helper method to extract user ID from context
func (c *UserController) getUserIDFromContext(ctx *gin.Context) (uuid.UUID, bool) {
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		return uuid.UUID{}, false
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, false
	}

	return userID, true
}

// GetProfile handles getting the user's profile
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, ok := c.getUserIDFromContext(ctx)
	if !ok {
		httpPkg.RespondError(ctx, http.StatusUnauthorized, constants.ErrUnauthorized)
		return
	}

	userProfile, err := c.userService.GetUserProfile(userID)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, userProfile, constants.MsgProfileFetched)
}

// UpdateProfile handles updating the user's profile
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	var input dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := c.getUserIDFromContext(ctx)
	if !ok {
		httpPkg.RespondError(ctx, http.StatusUnauthorized, constants.ErrUnauthorized)
		return
	}

	if err := c.userService.UpdateUserProfile(userID, input); err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgUserUpdated)
}

// ChangePassword handles updating the user's password
func (c *UserController) ChangePassword(ctx *gin.Context) {
	var input dto.ChangePasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := c.getUserIDFromContext(ctx)
	if !ok {
		httpPkg.RespondError(ctx, http.StatusUnauthorized, constants.ErrUnauthorized)
		return
	}

	if err := c.userService.ChangeUserPassword(userID, input.OldPassword, input.NewPassword); err != nil {
		if err.Error() == constants.ErrPasswordMismatch {
			httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrPasswordMismatch)
			return
		}
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgPasswordChanged)
}

// DeleteProfile handles deleting the user's account
func (c *UserController) DeleteProfile(ctx *gin.Context) {
	userID, ok := c.getUserIDFromContext(ctx)
	if !ok {
		httpPkg.RespondError(ctx, http.StatusUnauthorized, constants.ErrUnauthorized)
		return
	}

	if err := c.userService.DeleteUserProfile(userID); err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgUserDeleted)
}

// UploadAvatar handles uploading a user avatar
func (c *UserController) UploadAvatar(ctx *gin.Context) {
	userID, ok := c.getUserIDFromContext(ctx)
	if !ok {
		httpPkg.RespondError(ctx, http.StatusUnauthorized, constants.ErrUnauthorized)
		return
	}

	fileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, "Avatar file is required")
		return
	}

	avatarURL, err := c.userService.UploadUserAvatar(userID, fileHeader)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, gin.H{"avatar_url": avatarURL}, constants.MsgAvatarUploaded)
}

// GetAllUsers handles retrieving all users with pagination
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	users, total, err := c.userService.GetUsersWithPagination(page, limit)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Gender:    string(user.Gender),
			AvatarURL: user.Avatar,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	httpPkg.RespondSuccess(ctx, http.StatusOK, gin.H{
		"users":        userResponses,
		"total_users":  total,
		"total_pages":  totalPages,
		"current_page": page,
	}, constants.MsgUsersFetched)
}

// UpdateUserProfile handles updating a specific user profile (admin or self)
func (c *UserController) UpdateUserProfile(ctx *gin.Context) {
	userIDParam := ctx.Param("id")
	targetUserID, err := uuid.Parse(userIDParam)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	var input dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userService.UpdateUserProfile(targetUserID, input); err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgUserUpdated)
}

// DeleteUserAccount handles deleting a specific user account (admin or self)
func (c *UserController) DeleteUserAccount(ctx *gin.Context) {
	userIDParam := ctx.Param("id")
	targetUserID, err := uuid.Parse(userIDParam)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	if err := c.userService.DeleteUserProfile(targetUserID); err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgUserDeleted)
}
