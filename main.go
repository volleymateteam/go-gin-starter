package main

import(
	"go-gin-starter/config",
	"go-gin-starter/routes",
	"github.com/gin-gonic/gin"
	_ "go-gin-starter/docs" // swagger docs
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func main() {
	config.LoadEnv()
	database.ConnectDB()
	database.DB.AutoMigrate(&models.User{})

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.SetupRoutes(r)

	r.Run(":8000")
}
