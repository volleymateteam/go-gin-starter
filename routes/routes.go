package routes

import (
	"go-gin-starter/controllers"
	"go-gin-starter/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Public Routes (No authentication)
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.POST("/password/forgot", controllers.ForgotPassword)
		api.POST("/password/reset", controllers.ResetPassword)
		api.POST("/waitlist/submit", controllers.SubmitWaitlist)

		// Authenticated Routes (JWT required)
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
			// Admin User Management
			admin.GET("/users", controllers.GetAllUsers)
			admin.PUT("/users/:id", controllers.UpdateUserByAdmin)
			admin.DELETE("/users/:id", controllers.DeleteUserByAdmin)

			// Admin Waitlist Management
			admin.GET("/waitlist", controllers.GetAllWaitlist)
			admin.POST("/waitlist/:id/approve", controllers.ApproveWaitlistEntry)
			admin.DELETE("/waitlist/:id/reject", controllers.RejectWaitlistEntry)

			// Admin Team Management
			admin.POST("/teams", controllers.CreateTeam)
			admin.GET("/teams", controllers.GetAllTeams)
			admin.GET("/teams/:id", controllers.GetTeamByID)
			admin.PUT("/teams/:id", controllers.UpdateTeam)
			admin.DELETE("/teams/:id", controllers.DeleteTeam)
			// admin.POST("/teams/:id/assign", controllers.AssignTeamToUser)
			// admin.POST("/teams/:id/remove", controllers.RemoveTeamFromUser)
			// admin.GET("/teams/:id/members", controllers.GetTeamMembers)

			// Admin Season Management
			admin.POST("/seasons", controllers.CreateSeason)
			admin.GET("/seasons", controllers.GetAllSeasons)
			admin.GET("/seasons/:id", controllers.GetSeasonByID)
			admin.PUT("/seasons/:id", controllers.UpdateSeason)
			admin.DELETE("/seasons/:id", controllers.DeleteSeason)
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
