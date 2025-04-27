package routes

import (
	"go-gin-starter/controllers"
	"go-gin-starter/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.POST("/password/forgot", controllers.ForgotPassword)
		api.POST("/password/reset", controllers.ResetPassword)

		auth := api.Group("/")
		auth.Use(middleware.JWTAuth())
		auth.GET("/profile", controllers.GetProfile)
		auth.POST("/profile/upload-avatar", controllers.UploadAvatar)
		auth.PUT("/profile", controllers.UpdateProfile)
		auth.DELETE("/profile", controllers.DeleteProfile)
		auth.GET("/users", controllers.GetAllUsers)
		auth.PUT("/profile/change-password", controllers.ChangePassword)
	}
}
