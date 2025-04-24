package services

import (
	"errors"
	"go-gin-starter/database"
	"go-gin-starter/models"
	"go-gin-starter/utils"
)


func CreateUser(username, email, password string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
