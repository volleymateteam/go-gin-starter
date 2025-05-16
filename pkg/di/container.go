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
	TeamController      *controllers.TeamController
	MatchController     *controllers.MatchController
	SeasonController    *controllers.SeasonController
	// Add other controllers here as needed
}

// NewContainer initializes and returns a new dependency container
func NewContainer() *Container {
	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	waitlistRepo := repositories.NewWaitlistRepository()
	authRepo := repositories.NewAuthRepository()
	teamRepo := repositories.NewTeamRepository()
	matchRepo := repositories.NewMatchRepository()
	seasonRepo := repositories.NewSeasonRepository()

	// Add other repositories here as needed

	// Initialize services
	userService := services.NewUserService(userRepo)
	waitlistService := services.NewWaitlistService(waitlistRepo, userService)
	authService := services.NewAuthService(authRepo, userRepo)
	teamService := services.NewTeamService(teamRepo)
	matchService := services.NewMatchService(matchRepo, teamRepo, seasonRepo)
	seasonService := services.NewSeasonService(seasonRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	adminUserController := controllers.NewAdminUserController(userService)
	waitlistController := controllers.NewWaitlistController(waitlistService)
	authController := controllers.NewAuthController(authService)
	teamController := controllers.NewTeamController(teamService)
	matchController := controllers.NewMatchController(matchService)
	seasonController := controllers.NewSeasonController(seasonService)

	return &Container{
		UserController:      userController,
		AdminUserController: adminUserController,
		WaitlistController:  waitlistController,
		AuthController:      authController,
		TeamController:      teamController,
		MatchController:     matchController,
		SeasonController:    seasonController,
		// Add other controllers here as needed
	}
}
