package services

import (
	"go-gin-starter/models"

	"github.com/google/uuid"
)

// Global userService instance
var globalUserService UserService

// InitGlobalServices initializes global service references
func InitGlobalServices(userService UserService) {
	globalUserService = userService
}

// GetUserByID wrapper for backward compatibility
func GetUserByID(id uuid.UUID) (*models.User, error) {
	return globalUserService.GetUserByID(id)
}

// CreateUser wrapper for backward compatibility
func CreateUser(username, email, password, gender string) (*models.User, error) {
	return globalUserService.CreateUser(username, email, password, gender)
}

// UpdateUser wrapper for backward compatibility
func UpdateUser(user *models.User) error {
	return globalUserService.UpdateUser(user)
}

// DeleteUserByID wrapper for backward compatibility
func DeleteUserByID(id uuid.UUID) error {
	return globalUserService.DeleteUserByID(id)
}
