package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"
	"time"

	"github.com/google/uuid"
)

// AuthRepository defines the interface for authentication data operations
type AuthRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	UpdateRefreshToken(userID uuid.UUID, token string, expiry time.Time) error
	UpdateResetToken(email, token string, expiry time.Time) error
	GetUserByRefreshToken(token string) (*models.User, error)
	GetUserByResetToken(token string) (*models.User, error)
	UpdatePassword(userID uuid.UUID, hashedPassword string) error
}

// GormAuthRepository implements AuthRepository using GORM
type GormAuthRepository struct{}

// NewAuthRepository creates a new instance of AuthRepository
func NewAuthRepository() AuthRepository {
	return &GormAuthRepository{}
}

// GetUserByEmail finds a user by their email address
func (r *GormAuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateRefreshToken updates a user's refresh token and expiry
func (r *GormAuthRepository) UpdateRefreshToken(userID uuid.UUID, token string, expiry time.Time) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"refresh_token":        token,
		"refresh_token_expiry": expiry,
	}).Error
}

// UpdateResetToken updates a user's password reset token and expiry
func (r *GormAuthRepository) UpdateResetToken(email, token string, expiry time.Time) error {
	return database.DB.Model(&models.User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"reset_password_token":   token,
		"reset_password_expires": expiry,
	}).Error
}

// GetUserByRefreshToken finds a user by their refresh token
func (r *GormAuthRepository) GetUserByRefreshToken(token string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("refresh_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByResetToken finds a user by their password reset token
func (r *GormAuthRepository) GetUserByResetToken(token string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("reset_password_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdatePassword updates a user's password
func (r *GormAuthRepository) UpdatePassword(userID uuid.UUID, hashedPassword string) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password":               hashedPassword,
		"reset_password_token":   nil,
		"reset_password_expires": nil,
	}).Error
}

func SaveResetToken(user *models.User) error {
	return database.DB.Save(user).Error
}

func SaveNewPassword(user *models.User) error {
	return database.DB.Save(user).Error
}
