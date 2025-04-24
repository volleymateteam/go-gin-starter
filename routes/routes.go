package routes

import (
	"github.com/gin-gonic/gin"
	"go-gin-starter/controllers"
	"go-gin-starter/middleware"
)


func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/login", controllers.Login)

		auth := api.Group("/")
		auth.User(middleware.JWTAuth())
		auth.GET("/profile", controllers.Profile)
	}
}
