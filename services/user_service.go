package services

import (
	"go-gin-starter/database"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/repositories"
	"go-gin-starter/utils"
	"time"

	"github.com/google/uuid"
)

// CreateUser registers a new user
func CreateUser(username, email, password, gender string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Set default avatar based on gender
	var avatar string
	switch gender {
	case "male":
		avatar = "defaults/default-male.png"
	case "female":
		avatar = "defaults/default-female.png"
	default:
		avatar = "defaults/default-male.png"
	}

	var genderEnum models.GenderEnum
	switch gender {
	case "male":
		genderEnum = models.GenderMale
	case "female":
		genderEnum = models.GenderFemale
	default:
		genderEnum = models.GenderOther
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Gender:   genderEnum,
		Avatar:   avatar,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID fetches a user based on their UUID
func GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates the user's information
func UpdateUser(user *models.User) error {
	return database.DB.Save(user).Error
}

// DeleteUserByID soft-deletes a user based on UUID
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

// GetUsersWithPagination retrieves a paginated list of users
func GetUsersWithPagination(page int, limit int) ([]models.User, int64, error) {
	var users []models.User
	var totalUsers int64

	offset := (page - 1) * limit

	if err := database.DB.Model(&models.User{}).Where("deleted_at IS NULL").Count(&totalUsers).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.Where("deleted_at IS NULL").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalUsers, nil
}

// AdminUpdateUser updates user fields selectively
func AdminUpdateUser(id uuid.UUID, input *dto.AdminUpdateUserInput) (*models.User, error) {
	user, err := repositories.FindUserByID(id)
	if err != nil {
		return nil, err
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Gender != "" {
		user.Gender = input.Gender
	}
	if input.Role != "" {
		user.Role = input.Role
	}

	user.UpdatedAt = time.Now()

	if err := repositories.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUserPermissions(userID uuid.UUID, permissions []string) error {
	user, err := repositories.FindUserByID(userID)
	if err != nil {
		return err
	}

	user.ExtraPermissions = permissions // GORM automaps []string to postgres array or jsonb

	if err := repositories.UpdateUser(user); err != nil {
		return err
	}

	return nil
}
