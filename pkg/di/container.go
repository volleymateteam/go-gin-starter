package di

import (
	"go-gin-starter/controllers"
	"go-gin-starter/pkg/upload"
	"go-gin-starter/pkg/video"
	"go-gin-starter/repositories"
	"go-gin-starter/services"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
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
	TestController                 *controllers.TestController
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

	// Initialize AWS/video queue for match service
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		panic("Failed to create AWS session: " + err.Error())
	}
	s3Client := s3.New(sess)
	sqsClient := sqs.New(sess)
	videoProcessor := video.NewVideoProcessor(sess, s3Client, os.Getenv("AWS_BUCKET_NAME"))
	videoQueue := video.NewQueueManager(
		sqsClient,
		os.Getenv("VIDEO_PROCESSING_QUEUE_URL"),
		videoProcessor,
	)

	// Initialize services
	userService := services.NewUserService(userRepo)
	waitlistService := services.NewWaitlistService(waitlistRepo, userService)
	authService := services.NewAuthService(authRepo, userRepo)
	teamService := services.NewTeamService(teamRepo, uploadService)
	matchService := services.NewMatchService(matchRepo, teamRepo, seasonRepo, videoQueue)
	seasonService := services.NewSeasonService(seasonRepo, uploadService)

	// Initialize global service references for backward compatibility
	services.InitGlobalServices(userService)

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
	testController := controllers.NewTestController(videoQueue)

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
		TestController:                 testController,
		// Add other controllers here as needed
	}
}
