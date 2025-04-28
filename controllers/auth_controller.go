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
		utils.RespondError(c, http.StatusBadRequest, utils.ErrStrongPassword)
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
	utils.RespondSuccess(c, http.StatusCreated, response, utils.MsgUserRegistered)
}

func Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidCredentials)
		return
	}

	user, err := services.GetUserByEmail(input.Email)
	if err != nil || !utils.CheckPasswordHash(input.Password, user.Password) {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrInvalidCredentials)
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrTokenGenerationFailed)
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
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrResetTokenFailed)
		return
	}
	resetToken := hex.EncodeToString(b)
	expiry := time.Now().Add(15 * time.Minute)

	err := services.UpdateResetToken(input.Email, resetToken, expiry)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{"reset_token": resetToken}, utils.MsgResetTokenCreated)
}

func ResetPassword(c *gin.Context) {
	var input dto.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if !utils.IsStrongPassword(input.NewPassword) {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrStrongPassword)
		return
	}

	err := services.ResetUserPassword(input.Token, input.NewPassword)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, nil, utils.MsgPasswordReset)
}
