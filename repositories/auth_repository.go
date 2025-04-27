package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"
)

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByResetToken(token string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("reset_password_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func SaveResetToken(user *models.User) error {
	return database.DB.Save(user).Error
}

func SaveNewPassword(user *models.User) error {
	return database.DB.Save(user).Error
}
