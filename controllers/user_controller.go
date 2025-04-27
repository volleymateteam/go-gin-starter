package controllers

import (
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"path/filepath"

	"github.com/google/uuid"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UpdateUserInput struct {
	Username string `json:"username" binding:"omitempty,min=3,max=20"`
	Email    string `json:"email" binding:"omitempty,email"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.IsStrongPassword(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least 8 characters, including uppercase, lowercase, numbers, and special characters."})
		return
	}
	user, err := services.CreateUser(input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	c.JSON(http.StatusCreated, gin.H{"user": response})
}

func GetProfile(c *gin.Context) {
	userINRow, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userINRow.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in token"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK, gin.H{"user": response})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": response})
}

func UpdateProfile(c *gin.Context) {
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	if err := services.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK, gin.H{"user": response})
}

// this endpoint is for admin use only! we will use middleware to check if the user is admin
func GetAllUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	users, totalUsers, err := services.GetUsersWithPagination(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []UserResponse
	for _, user := range users {
		response = append(response, UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	totalPages := (int(totalUsers) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"users":        response,
		"total_users":  totalUsers,
		"total_pages":  totalPages,
		"current_page": page,
	})
}

func DeleteProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	err := services.DeleteUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func ChangePassword(c *gin.Context) {
	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check old password
	if !utils.CheckPasswordHash(input.OldPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid old password"})
		return
	}

	// Validate new password strength
	if !utils.IsStrongPassword(input.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password must contain at least 8 characters, including uppercase, lowercase, numbers, and special characters."})
		return
	}

	// Hash and update new password
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	user.Password = hashedPassword

	if err := services.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func UploadAvatar(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Multipart form
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
		return
	}

	// Check file size - max 2 MB (2 * 1024 * 1024 bytes)
	const maxAvatarSize = 2 * 1024 * 1024

	if file.Size > maxAvatarSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is too large. Max size is 2MB"})
		return
	}

	// Optional: Validate file type (e.g., jpg, png only)
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only .jpg, .jpeg and .png formats are allowed"})
		return
	}

	// Generate a unique filename
	newFileName := uuid.New().String() + ext
	savePath := filepath.Join("uploads/avatars", newFileName)

	// Save file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload avatar"})
		return
	}

	// Update user's avatar field
	user.Avatar = newFileName
	if err := services.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user avatar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar uploaded successfully", "avatar_url": "/uploads/avatars/" + newFileName})
}
