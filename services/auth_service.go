package services

import (
	"errors"
	"time"

	"go-gin-starter/database"
	"go-gin-starter/models"
	"go-gin-starter/utils"
)

// GetUserByEmail finds a user by their email (used for login and forgot password)
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateResetToken sets a password reset token and expiry for the user
func UpdateResetToken(email, token string, expiry time.Time) error {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	user.ResetPasswordToken = &token
	user.ResetPasswordExpires = &expiry

	return database.DB.Save(&user).Error
}

// ResetUserPassword changes the user's password based on a valid reset token
func ResetUserPassword(token, newPassword string) error {
	var user models.User
	if err := database.DB.Where("reset_password_token = ?", token).First(&user).Error; err != nil {
		return errors.New("invalid or expired token")
	}

	// Check if the reset token is still valid
	if user.ResetPasswordExpires == nil || user.ResetPasswordExpires.Before(time.Now()) {
		return errors.New("reset token has expired")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.ResetPasswordToken = nil
	user.ResetPasswordExpires = nil

	return database.DB.Save(&user).Error
}
