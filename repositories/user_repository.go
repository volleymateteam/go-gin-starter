package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"

	"github.com/google/uuid"
)

func FindUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func UpdateUser(user *models.User) error {
	return database.DB.Save(user).Error
}

func DeleteUser(user *models.User) error {
	return database.DB.Delete(user).Error
}

func GetUsersWithPagination(limit int, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := database.DB.Model(&models.User{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.Where("deleted_at IS NULL").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
