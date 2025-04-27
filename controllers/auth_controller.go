package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"go-gin-starter/dto"
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if !utils.IsStrongPassword(input.Password) {
		utils.RespondError(c, http.StatusBadRequest, "Password must contain at least 8 characters, including uppercase, lowercase, numbers, and special characters.")
		return
	}

	user, err := services.CreateUser(input.Username, input.Email, input.Password, input.Gender)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
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
	utils.RespondSuccess(c, response, "User registered successfully")
}

func Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := services.GetUserByEmail(input.Email)
	if err != nil || !utils.CheckPasswordHash(input.Password, user.Password) {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Token generation failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func ForgotPassword(c *gin.Context) {
	var input dto.ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to generate reset token")
		return
	}
	resetToken := hex.EncodeToString(b)
	expiry := time.Now().Add(15 * time.Minute)

	err := services.UpdateResetToken(input.Email, resetToken, expiry)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"reset_token": resetToken}, "Reset token generated successfully")
}

func ResetPassword(c *gin.Context) {
	var input dto.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if !utils.IsStrongPassword(input.NewPassword) {
		utils.RespondError(c, http.StatusBadRequest, "Password must contain at least 8 characters, including uppercase, lowercase, numbers, and special characters.")
		return
	}

	err := services.ResetUserPassword(input.Token, input.NewPassword)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondSuccess(c, nil, "Password reset successfully")
}
