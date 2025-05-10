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
		AvatarURL: "/uploads/avatars/" + user.Avatar,
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

	token, err := authPkg.GenerateJWT(user.ID)
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrTokenGenerationFailed)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
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
