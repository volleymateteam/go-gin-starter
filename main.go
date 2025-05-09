package main

import (
	"go-gin-starter/config"
	"go-gin-starter/database"
	"go-gin-starter/middleware"
	"go-gin-starter/models"
	"go-gin-starter/routes"

	_ "go-gin-starter/docs" // swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load environment variables and connect to the database
	config.LoadEnv()
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

	r := gin.Default()

	// Setup global middlewares
	r.Use(middleware.ErrorRecovery())
	r.Use(middleware.RequestLogger())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Setup swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup all API routes
	routes.SetupRoutes(r)

	// Start the server
	r.Run(":8080")
}
