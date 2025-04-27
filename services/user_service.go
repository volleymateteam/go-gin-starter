package services

import (
	"errors"
	"time"

	"go-gin-starter/database"
	"go-gin-starter/models"
	"go-gin-starter/utils"

	"github.com/google/uuid"
)

func CreateUser(username, email, password, gender string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Default avatar based on gender
	var avatar string
	switch gender {
	case "male":
		avatar = "defaults/default-male.png"
	case "female":
		avatar = "defaults/default-female.png"
	default:
		avatar = "defaults/default-male.png"
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Gender:   gender,
		Avatar:   avatar,
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

// UpdateResetToken sets a reset token and expiry time
func UpdateResetToken(email, token string, expiry time.Time) error {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	user.ResetPasswordToken = &token
	user.ResetPasswordExpires = &expiry

	return database.DB.Save(&user).Error
}

// ResetUserPassword resets password using token
func ResetUserPassword(token, newPassword string) error {
	var user models.User
	if err := database.DB.Where("reset_password_token = ?", token).First(&user).Error; err != nil {
		return errors.New("invalid or expired token")
	}

	// Check expiry
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
