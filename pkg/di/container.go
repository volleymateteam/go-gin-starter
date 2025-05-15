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
	AuthController      *controllers.AuthController
	// Add other controllers here as needed
}

// NewContainer initializes and returns a new dependency container
func NewContainer() *Container {
	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	waitlistRepo := repositories.NewWaitlistRepository()
	authRepo := repositories.NewAuthRepository()

	// Initialize services
	userService := services.NewUserService(userRepo)
	waitlistService := services.NewWaitlistService(waitlistRepo, userService)
	authService := services.NewAuthService(authRepo, userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	adminUserController := controllers.NewAdminUserController(userService)
	waitlistController := controllers.NewWaitlistController(waitlistService)
	authController := controllers.NewAuthController(authService)

	return &Container{
		UserController:      userController,
		AdminUserController: adminUserController,
		WaitlistController:  waitlistController,
		AuthController:      authController,
		// Add other controllers here as needed
	}
}
