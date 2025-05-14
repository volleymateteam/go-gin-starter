package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"go-gin-starter/dto"
	authPkg "go-gin-starter/pkg/auth"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if !authPkg.IsStrongPassword(input.Password) {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrStrongPassword)
		return
	}

	user, err := services.CreateUser(input.Username, input.Email, input.Password, input.Gender)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Gender:    string(user.Gender),
		AvatarURL: "/uploads/avatars/" + user.Avatar, // have to be FIXED
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
	httpPkg.RespondSuccess(c, http.StatusCreated, response, constants.MsgUserRegistered)
}

func Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidCredentials)
		return
	}

	user, err := services.GetUserByEmail(input.Email)
	if err != nil || !authPkg.CheckPasswordHash(input.Password, user.Password) {
		httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrInvalidCredentials)
		return
	}

	accessToken, err := authPkg.GenerateJWT(user.ID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrTokenGenerationFailed)
		return
	}

	refreshToken, err := authPkg.GenerateSecureToken(32)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)

	err = services.UpdateRefreshToken(user.ID, refreshToken, refreshExpiry)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, "failed to save refresh token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidToken)
		return
	}

	user, err := services.GetUserByRefreshToken(input.RefreshToken)
	if err != nil || user.RefreshTokenExpiry.Before(time.Now()) {
		httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrInvalidToken)
		return
	}

	// Generate new access token
	accessToken, err := authPkg.GenerateJWT(user.ID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrTokenGenerationFailed)
		return
	}

	// (Optional) rotate refresh token
	newRefreshToken, err := authPkg.GenerateSecureToken(32)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}
	newExpiry := time.Now().Add(7 * 24 * time.Hour)

	if err := services.UpdateRefreshToken(user.ID, newRefreshToken, newExpiry); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, "failed to update refresh token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

func ForgotPassword(c *gin.Context) {
	var input dto.ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrResetTokenFailed)
		return
	}
	resetToken := hex.EncodeToString(b)
	expiry := time.Now().Add(15 * time.Minute)

	err := services.UpdateResetToken(input.Email, resetToken, expiry)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, gin.H{"reset_token": resetToken}, constants.MsgResetTokenCreated)
}

func ResetPassword(c *gin.Context) {
	var input dto.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if !authPkg.IsStrongPassword(input.NewPassword) {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrStrongPassword)
		return
	}

	err := services.ResetUserPassword(input.Token, input.NewPassword)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgPasswordReset)
}
