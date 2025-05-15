package di

import (
	"go-gin-starter/controllers"
	"go-gin-starter/repositories"
	"go-gin-starter/services"
)

// Container holds all the dependency instances
type Container struct {
	UserController      *controllers.UserController
	AdminUserController *controllers.AdminUserController
	// Add other controllers here as needed
}

// NewContainer initializes and returns a new dependency container
func NewContainer() *Container {
	// Initialize repositories
	userRepo := repositories.NewUserRepository()

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	adminUserController := controllers.NewAdminUserController(userService)

	return &Container{
		UserController:      userController,
		AdminUserController: adminUserController,
		// Add other controllers here as needed
	}
}
