package services

import (
	"go-gin-starter/database"
	"go-gin-starter/models"
	"go-gin-starter/utils"

	"github.com/google/uuid"
)


func CreateUser(username, email, password string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *models.User) error {
	return database.DB.Save(user).Error
}


func GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}


func DeleteUserByID(id uuid.UUID) error {
	var user models.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		return err
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}


func GetUsersWithPagination(page int, limit int) ([]models.User, int64, error) {
	var users []models.User
	var totalUsers int64

	offset := (page - 1) * limit

	// Count total users (without deleted)
	if err := database.DB.Model(&models.User{}).Where("deleted_at IS NULL").Count(&totalUsers).Error; err != nil {
		return nil, 0, err
	}

	// Fetch users with limit and offset
	if err := database.DB.Where("deleted_at IS NULL").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalUsers, nil
}
