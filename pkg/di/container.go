package di

import (
	"go-gin-starter/controllers"
	"go-gin-starter/pkg/upload"
	"go-gin-starter/repositories"
	"go-gin-starter/services"
)

// Container holds all the dependency instances
type Container struct {
	UserController                 *controllers.UserController
	AdminUserController            *controllers.AdminUserController
	AdminUserPermissionsController *controllers.AdminUserPermissionsController
	AdminAuditController           *controllers.AdminAuditController
	WaitlistController             *controllers.WaitlistController
	AuthController                 *controllers.AuthController
	TeamController                 *controllers.TeamController
	MatchController                *controllers.MatchController
	SeasonController               *controllers.SeasonController
	HealthController               *controllers.HealthController
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

	// Initialize utility services
	uploadService := upload.NewFileUploadService()

	// Initialize services
	userService := services.NewUserService(userRepo)
	waitlistService := services.NewWaitlistService(waitlistRepo, userService)
	authService := services.NewAuthService(authRepo, userRepo)
	teamService := services.NewTeamService(teamRepo, uploadService)
	matchService := services.NewMatchService(matchRepo, teamRepo, seasonRepo)
	seasonService := services.NewSeasonService(seasonRepo, uploadService)

	// Initialize controllers
	userController := controllers.NewUserController(userService, uploadService)
	adminUserController := controllers.NewAdminUserController(userService)
	adminUserPermissionsController := controllers.NewAdminUserPermissionsController(userService)
	adminAuditController := controllers.NewAdminAuditController()
	waitlistController := controllers.NewWaitlistController(waitlistService)
	authController := controllers.NewAuthController(authService)
	teamController := controllers.NewTeamController(teamService, uploadService)
	matchController := controllers.NewMatchController(matchService)
	seasonController := controllers.NewSeasonController(seasonService, uploadService)
	healthController := controllers.NewHealthController()

	return &Container{
		UserController:                 userController,
		AdminUserController:            adminUserController,
		AdminUserPermissionsController: adminUserPermissionsController,
		AdminAuditController:           adminAuditController,
		WaitlistController:             waitlistController,
		AuthController:                 authController,
		TeamController:                 teamController,
		MatchController:                matchController,
		SeasonController:               seasonController,
		HealthController:               healthController,
		// Add other controllers here as needed
	}
}
