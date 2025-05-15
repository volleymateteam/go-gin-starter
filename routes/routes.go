package routes

import (
	"go-gin-starter/controllers"
	"go-gin-starter/middleware"
	"go-gin-starter/pkg/di"

	"github.com/gin-gonic/gin"
)

// SetupRoutes registers all routes on the given router group
func SetupRoutes(router gin.IRouter) {
	// Create the container and get the user controller
	container := di.NewContainer()
	userCtrl := container.UserController
	adminUserCtrl := container.AdminUserController
	waitlistCtrl := container.WaitlistController

	// Public Routes (No authentication)
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST("/refresh-token", controllers.RefreshToken)
	router.POST("/password/forgot", controllers.ForgotPassword)
	router.POST("/password/reset", controllers.ResetPassword)
	router.POST("/waitlist/submit", waitlistCtrl.SubmitWaitlist)

	// Authenticated Routes (JWT required)
	auth := router.Group("/")
	auth.Use(middleware.JWTAuth())

	// Normal authenticated user routes
	auth.GET("/profile", userCtrl.GetProfile)
	auth.POST("/profile/upload-avatar", userCtrl.UploadAvatar)
	auth.PUT("/profile", userCtrl.UpdateProfile)
	auth.DELETE("/profile", userCtrl.DeleteProfile)
	auth.PUT("/profile/change-password", userCtrl.ChangePassword)

	// Admin permission-based routes
	admin := auth.Group("/admin")
	{
		// Admin User Management
		admin.GET("/users", middleware.RequirePermission("manage_users"), userCtrl.GetAllUsers)
		admin.PUT("/users/:id", middleware.RequirePermission("manage_users"), adminUserCtrl.UpdateUserByAdmin)
		admin.DELETE("/users/:id", middleware.RequirePermission("manage_users"), adminUserCtrl.DeleteUserByAdmin)
		admin.PATCH("/users/:id/permissions", middleware.RequirePermission("manage_users"), adminUserCtrl.UpdateUserPermissions)
		admin.GET("/users/:id/permissions", middleware.RequirePermission("manage_users"), adminUserCtrl.GetUserPermissions)
		admin.PATCH("/users/:id/permissions/reset", middleware.RequirePermission("manage_users"), adminUserCtrl.ResetUserPermissions)

		// Admin Audit Logging
		admin.GET("/audit-logs", middleware.RequirePermission("view_audit_logs"), adminUserCtrl.GetAuditLogs)

		// Admin Waitlist Management
		admin.GET("/waitlist", middleware.RequirePermission("manage_waitlist"), waitlistCtrl.GetAllWaitlist)
		admin.POST("/waitlist/:id/approve", middleware.RequirePermission("manage_waitlist"), waitlistCtrl.ApproveWaitlistEntry)
		admin.DELETE("/waitlist/:id/reject", middleware.RequirePermission("manage_waitlist"), waitlistCtrl.RejectWaitlistEntry)

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
		user.PUT("/:id/update", userCtrl.UpdateUserProfile)
		user.DELETE("/:id/delete", userCtrl.DeleteUserAccount)
	}
}
