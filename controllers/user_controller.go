package controllers

import (
	"fmt"
	"go-gin-starter/dto"
	authPkg "go-gin-starter/pkg/auth"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/pkg/storage"
	"go-gin-starter/services"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Gender:    string(user.Gender),
		AvatarURL: user.Avatar,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
	httpPkg.RespondSuccess(c, http.StatusOK, response, constants.MsgProfileFetched)
}

func UpdateProfile(c *gin.Context) {
	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	if err := services.UpdateUser(user); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgUserUpdated)
}

func ChangePassword(c *gin.Context) {
	var input dto.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	if !authPkg.CheckPasswordHash(input.OldPassword, user.Password) {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrPasswordMismatch)
		return
	}

	hashedPassword, err := authPkg.HashPassword(input.NewPassword)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = hashedPassword

	if err := services.UpdateUser(user); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgPasswordChanged)
}

func DeleteProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
		return
	}

	if err := services.DeleteUserByID(userID); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgUserDeleted)
}

func UploadAvatar(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, "Avatar file is required")
		return
	}

	const maxAvatarSize = 2 * 1024 * 1024
	if fileHeader.Size > maxAvatarSize {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrAvatarTooLarge)
		return
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidFileType)
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}
	defer src.Close()

	newFileName := uuid.New().String() + ext
	objectKey := fmt.Sprintf("avatars/%s", newFileName)

	contentType := fileHeader.Header.Get("Content-Type")

	avatarURL, err := storage.UploadFileToS3(src, objectKey, contentType)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrUploadFailed)
		return
	}

	user.Avatar = avatarURL
	if err := services.UpdateUser(user); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, gin.H{"avatar_url": avatarURL}, constants.MsgAvatarUploaded)
}

func GetAllUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	users, total, err := services.GetUsersWithPagination(page, limit)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
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

	httpPkg.RespondSuccess(c, http.StatusOK, gin.H{
		"users":        userResponses,
		"total_users":  total,
		"total_pages":  totalPages,
		"current_page": page,
	}, constants.MsgUsersFetched)
}

func UpdateUserProfile(c *gin.Context) {
	userIDParam := c.Param("id")
	targetUserID, err := uuid.Parse(userIDParam)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	var input struct {
		Username string `json:"username" binding:"omitempty,min=3,max=20"`
		Email    string `json:"email" binding:"omitempty,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := services.GetUserByID(targetUserID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusNotFound, constants.ErrUserNotFound)
		return
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	if err := services.UpdateUser(user); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgUserUpdated)
}

func DeleteUserAccount(c *gin.Context) {
	userIDParam := c.Param("id")
	targetUserID, err := uuid.Parse(userIDParam)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	err = services.DeleteUserByID(targetUserID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgUserDeleted)
}
