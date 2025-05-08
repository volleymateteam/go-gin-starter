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

		// Admin permission-based routes
		admin := auth.Group("/admin")
		{
			// Admin User Management
			admin.GET("/users", middleware.RequirePermission("manage_users"), controllers.GetAllUsers)
			admin.PUT("/users/:id", middleware.RequirePermission("manage_users"), controllers.UpdateUserByAdmin)
			admin.DELETE("/users/:id", middleware.RequirePermission("manage_users"), controllers.DeleteUserByAdmin)

			// Admin Waitlist Management
			admin.GET("/waitlist", middleware.RequirePermission("manage_waitlist"), controllers.GetAllWaitlist)
			admin.POST("/waitlist/:id/approve", middleware.RequirePermission("manage_waitlist"), controllers.ApproveWaitlistEntry)
			admin.DELETE("/waitlist/:id/reject", middleware.RequirePermission("manage_waitlist"), controllers.RejectWaitlistEntry)

			// Admin Team Management
			admin.POST("/teams", middleware.RequirePermission("manage_teams"), controllers.CreateTeam)
			admin.GET("/teams", middleware.RequirePermission("manage_teams"), controllers.GetAllTeams)
			admin.GET("/teams/:id", middleware.RequirePermission("manage_teams"), controllers.GetTeamByID)
			admin.PUT("/teams/:id", middleware.RequirePermission("manage_teams"), controllers.UpdateTeam)
			admin.DELETE("/teams/:id", middleware.RequirePermission("manage_teams"), controllers.DeleteTeam)
			admin.PATCH("/teams/:id/upload-logo", middleware.RequirePermission("manage_teams"), controllers.UploadTeamLogo)

			// Admin Season Management
			admin.POST("/seasons", middleware.RequirePermission("manage_season"), controllers.CreateSeason)
			admin.GET("/seasons", middleware.RequirePermission("manage_season"), controllers.GetAllSeasons)
			admin.GET("/seasons/:id", middleware.RequirePermission("manage_season"), controllers.GetSeasonByID)
			admin.PUT("/seasons/:id", middleware.RequirePermission("manage_season"), controllers.UpdateSeason)
			admin.DELETE("/seasons/:id", middleware.RequirePermission("manage_season"), controllers.DeleteSeason)
			admin.PATCH("/seasons/:id/upload-logo", middleware.RequirePermission("manage_season"), controllers.UploadSeasonLogo)

			// Admin Match Management
			admin.POST("/matches", middleware.RequirePermission("manage_matches"), controllers.CreateMatch)
			admin.GET("/matches", middleware.RequirePermission("manage_matches"), controllers.GetAllMatches)
			admin.GET("/matches/:id", middleware.RequirePermission("manage_matches"), controllers.GetMatchByID)
			admin.PUT("/matches/:id", middleware.RequirePermission("manage_matches"), controllers.UpdateMatch)
			admin.DELETE("/matches/:id", middleware.RequirePermission("manage_matches"), controllers.DeleteMatch)
			admin.PATCH("/matches/:id/upload-video", middleware.RequirePermission("upload_video"), controllers.UploadMatchVideo)
			admin.PATCH("/matches/:id/upload-scout", middleware.RequirePermission("upload_scout"), controllers.UploadMatchScout)
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
