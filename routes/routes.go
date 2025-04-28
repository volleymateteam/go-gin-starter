package routes

import (
	"go-gin-starter/controllers"
	"go-gin-starter/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Public routes (no auth)
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.POST("/password/forgot", controllers.ForgotPassword)
		api.POST("/password/reset", controllers.ResetPassword)

		// Public waitlist route (Anyone can submit waitlist)
		api.POST("/waitlist/submit", controllers.SubmitWaitlist)

		// Protected routes (need JWT)
		auth := api.Group("/")
		auth.Use(middleware.JWTAuth())

		// Normal authenticated user routes
		auth.GET("/profile", controllers.GetProfile)
		auth.POST("/profile/upload-avatar", controllers.UploadAvatar)
		auth.PUT("/profile", controllers.UpdateProfile)
		auth.DELETE("/profile", controllers.DeleteProfile)
		auth.PUT("/profile/change-password", controllers.ChangePassword)

		// Admin-only routes
		admin := auth.Group("/admin")
		admin.Use(middleware.AdminOnly())
		{
			admin.GET("/users", controllers.GetAllUsers)
			admin.PUT("/users/:id", controllers.UpdateUserByAdmin)
			admin.DELETE("/users/:id", controllers.DeleteUserByAdmin)

			admin.GET("/waitlist", controllers.GetAllWaitlist)
		}

		// AdminOrSelf routes
		user := auth.Group("/users")
		user.Use(middleware.AdminOrSelf())
		{
			user.PUT("/:id/update", controllers.UpdateUserProfile)
			user.DELETE("/:id/delete", controllers.DeleteUserAccount)
		}
	}
}
