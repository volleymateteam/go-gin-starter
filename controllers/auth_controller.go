package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	authService services.AuthService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register handles user registration
func (c *AuthController) Register(ctx *gin.Context) {
	var input dto.RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.authService.Register(input.Username, input.Email, input.Password, input.Gender)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
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
	httpPkg.RespondSuccess(ctx, http.StatusCreated, response, constants.MsgUserRegistered)
}

// Login handles user authentication
func (c *AuthController) Login(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidCredentials)
		return
	}

	accessToken, refreshToken, err := c.authService.Login(input.Email, input.Password)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken handles token refresh
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidToken)
		return
	}

	accessToken, refreshToken, err := c.authService.RefreshToken(input.RefreshToken)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// ForgotPassword handles password reset requests
func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var input dto.ForgotPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resetToken, err := c.authService.ForgotPassword(input.Email)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, gin.H{"reset_token": resetToken}, constants.MsgResetTokenCreated)
}

// ResetPassword handles password reset using a token
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var input dto.ResetPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := c.authService.ResetPassword(input.Token, input.NewPassword)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgPasswordReset)
}
