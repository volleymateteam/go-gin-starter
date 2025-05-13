package main

import (
	"go-gin-starter/config"
        "log"
        "os"
	"go-gin-starter/controllers"
	"go-gin-starter/database"
	"go-gin-starter/middleware"
	"go-gin-starter/models"
	"go-gin-starter/pkg/logger"
	"go-gin-starter/routes"

	_ "go-gin-starter/docs" // swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables and connect to the database
	config.LoadEnv()
	config.InitConfig()

	log.Printf("DEBUG: ENV=%s, VIDEO_CLOUDFRONT_DOMAIN='%s'", 
		 os.Getenv("ENV"), 
	         os.Getenv("VIDEO_CLOUDFRONT_DOMAIN"))

      

	// Initialize structured logger
	logger.Init()
	defer logger.Sync()

	// Connect to the database
	database.ConnectDB()
	database.DB.AutoMigrate(
		&models.User{},
		&models.WaitlistEntry{},
		&models.Season{},
		&models.Team{},
		&models.Match{},
		&models.AdminActionLog{},
		// &models.UserActionLog{},
	)

	// Set Gin mode based on environment
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New() // Using New() instead of Default() for custom logging

	// Setup global middlewares
	r.Use(middleware.ErrorRecovery())
	r.Use(middleware.CORS())
	r.Use(logger.Middleware())

	// Apply rate limiting to all routes except health/readiness checks
	r.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path != "/health" && path != "/readiness" && path != "/liveness" {
			middleware.RateLimiter(10, 20)(c) // 10 requests/second with burst of 20
		} else {
			c.Next()
		}
	})

	// Health check endpoints
	r.GET("/health", controllers.HealthCheck)
	r.GET("/health/details", controllers.DetailedHealthCheck)
	r.GET("/readiness", controllers.ReadinessCheck)
	r.GET("/liveness", controllers.LivenessCheck)

	// Setup API versioning - V1 routes
	v1 := r.Group("/api/v1")
	{
		// Setup swagger
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Setup all API routes under versioned path
		routes.SetupRoutes(v1)
	}

	// Keep legacy routes for backward compatibility
	routes.SetupRoutes(r.Group("/api"))

	// Start the server on the specified port
	port := config.GetEnvWithDefault("PORT", "8080")
	logger.Info("Server starting", zap.String("port", port))
	r.Run(":" + port)
}
