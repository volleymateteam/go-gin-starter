package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindByID(id uuid.UUID) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(user *models.User) error
	GetWithPagination(limit, offset int) ([]models.User, int64, error)
}

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct{}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository() UserRepository {
	return &GormUserRepository{}
}

// FindByID retrieves a user by ID
func (r *GormUserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create adds a new user to the database
func (r *GormUserRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

// Update modifies an existing user
func (r *GormUserRepository) Update(user *models.User) error {
	return database.DB.Save(user).Error
}

// Delete removes a user (soft delete with GORM)
func (r *GormUserRepository) Delete(user *models.User) error {
	return database.DB.Delete(user).Error
}

// GetWithPagination retrieves users with pagination
func (r *GormUserRepository) GetWithPagination(limit, offset int) ([]models.User, int64, error) {
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
