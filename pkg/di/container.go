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
	WaitlistController  *controllers.WaitlistController
	// Add other controllers here as needed
}

// NewContainer initializes and returns a new dependency container
func NewContainer() *Container {
	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	waitlistRepo := repositories.NewWaitlistRepository()

	// Initialize services
	userService := services.NewUserService(userRepo)
	waitlistService := services.NewWaitlistService(waitlistRepo, userService)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	adminUserController := controllers.NewAdminUserController(userService)
	waitlistController := controllers.NewWaitlistController(waitlistService)

	return &Container{
		UserController:      userController,
		AdminUserController: adminUserController,
		WaitlistController:  waitlistController,
		// Add other controllers here as needed
	}
}
