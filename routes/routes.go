package routes

import (
	"github.com/gin-gonic/gin"
	"go-gin-starter/controllers"
	"go-gin-starter/middleware"
)


func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
	

		auth := api.Group("/")
		auth.Use(middleware.JWTAuth())
		auth.GET("/profile", controllers.GetProfile)
		auth.PUT("/profile", controllers.UpdateProfile)
		auth.DELETE("/profile", controllers.DeleteProfile)
		auth.GET("/users", controllers.GetAllUsers)
	}
}
