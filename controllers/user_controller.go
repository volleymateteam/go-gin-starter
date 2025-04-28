package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/services"
	"go-gin-starter/utils"
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
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.RespondError(c, http.StatusInternalServerError, "Invalid user ID")
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "User not found")
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Gender:    string(user.Gender),
		AvatarURL: "/uploads/avatars/" + user.Avatar,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
	utils.RespondSuccess(c, response, "User profile retrieved successfully")
}

func UpdateProfile(c *gin.Context) {
	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.RespondError(c, http.StatusInternalServerError, "Invalid user ID")
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "User not found")
		return
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	if err := services.UpdateUser(user); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, nil, "Profile updated successfully")
}

func ChangePassword(c *gin.Context) {
	var input dto.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.RespondError(c, http.StatusInternalServerError, "Invalid user ID")
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "User not found")
		return
	}

	if !utils.CheckPasswordHash(input.OldPassword, user.Password) {
		utils.RespondError(c, http.StatusBadRequest, "Invalid old password")
		return
	}

	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = hashedPassword

	if err := services.UpdateUser(user); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, nil, "Password changed successfully")
}

func DeleteProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.RespondError(c, http.StatusInternalServerError, "Invalid user ID")
		return
	}

	if err := services.DeleteUserByID(userID); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, nil, "Profile deleted successfully")
}

func UploadAvatar(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.RespondError(c, http.StatusInternalServerError, "Invalid user ID")
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "User not found")
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Avatar file is required")
		return
	}

	const maxAvatarSize = 2 * 1024 * 1024
	if file.Size > maxAvatarSize {
		utils.RespondError(c, http.StatusBadRequest, "Avatar file is too large. Max size is 2MB")
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		utils.RespondError(c, http.StatusBadRequest, "Invalid file format, only .jpg, .jpeg, .png allowed")
		return
	}

	newFileName := uuid.New().String() + ext
	savePath := filepath.Join("uploads/avatars", newFileName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to upload avatar")
		return
	}

	user.Avatar = newFileName
	if err := services.UpdateUser(user); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"avatar_url": "/uploads/avatars/" + newFileName}, "Avatar uploaded successfully")
}

func GetAllUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	users, total, err := services.GetUsersWithPagination(page, limit)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Gender:    string(user.Gender),
			AvatarURL: "/uploads/avatars/" + user.Avatar,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	utils.RespondSuccess(c, gin.H{
		"users":        userResponses,
		"total_users":  total,
		"total_pages":  totalPages,
		"current_page": page,
	}, "Users fetched successfully")
}

func UpdateUserProfile(c *gin.Context) {
	userIDParam := c.Param("id")
	targetUserID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var input struct {
		Username string `json:"username" binding:"omitempty,min=3,max=20"`
		Email    string `json:"email" binding:"omitempty,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := services.GetUserByID(targetUserID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "User not found")
		return
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	if err := services.UpdateUser(user); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUserAccount(c *gin.Context) {
	userIDParam := c.Param("id")
	targetUserID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = services.DeleteUserByID(targetUserID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
