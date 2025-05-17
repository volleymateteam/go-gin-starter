package routes

import (
	"go-gin-starter/middleware"
	"go-gin-starter/pkg/di"

	"github.com/gin-gonic/gin"
)

// SetupRoutes registers all routes on the given router group
func SetupRoutes(router gin.IRouter) {
	// Create the container and get the controllers
	container := di.NewContainer()
	userCtrl := container.UserController
	adminUserCtrl := container.AdminUserController
	adminPermissionsCtrl := container.AdminUserPermissionsController
	adminAuditCtrl := container.AdminAuditController
	waitlistCtrl := container.WaitlistController
	authCtrl := container.AuthController
	teamCtrl := container.TeamController
	matchCtrl := container.MatchController
	seasonCtrl := container.SeasonController
	healthCtrl := container.HealthController

	// Health check routes
	router.GET("/health", healthCtrl.HealthCheck)
	router.GET("/health/detailed", healthCtrl.DetailedHealthCheck)
	router.GET("/health/liveness", healthCtrl.LivenessCheck)
	router.GET("/health/readiness", healthCtrl.ReadinessCheck)

	// Public Routes (No authentication)
	router.POST("/register", authCtrl.Register)
	router.POST("/login", authCtrl.Login)
	router.POST("/refresh-token", authCtrl.RefreshToken)
	router.POST("/password/forgot", authCtrl.ForgotPassword)
	router.POST("/password/reset", authCtrl.ResetPassword)
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

	// Public read-only season routes (available to all authenticated users)
	auth.GET("/seasons", seasonCtrl.GetAllSeasons)
	auth.GET("/seasons/:id", seasonCtrl.GetSeasonByID)

	// Public read-only team routes (available to all authenticated users)
	auth.GET("/teams", teamCtrl.GetAllTeams)
	auth.GET("/teams/:id", teamCtrl.GetTeamByID)

	// Public read-only match routes (available to all authenticated users)
	auth.GET("/matches", matchCtrl.GetAllMatches)
	auth.GET("/matches/:id", matchCtrl.GetMatchByID)

	// Admin permission-based routes
	admin := auth.Group("/admin")
	{
		// Admin User Management
		admin.GET("/users", middleware.RequirePermission("manage_users"), userCtrl.GetAllUsers)
		admin.PUT("/users/:id", middleware.RequirePermission("manage_users"), adminUserCtrl.UpdateUserByAdmin)
		admin.DELETE("/users/:id", middleware.RequirePermission("manage_users"), adminUserCtrl.DeleteUserByAdmin)

		// Admin User Permissions Management
		admin.PATCH("/users/:id/permissions", middleware.RequirePermission("manage_users"), adminPermissionsCtrl.UpdateUserPermissions)
		admin.GET("/users/:id/permissions", middleware.RequirePermission("manage_users"), adminPermissionsCtrl.GetUserPermissions)
		admin.PATCH("/users/:id/permissions/reset", middleware.RequirePermission("manage_users"), adminPermissionsCtrl.ResetUserPermissions)

		// Admin Audit Logging
		admin.GET("/audit-logs", middleware.RequirePermission("view_audit_logs"), adminAuditCtrl.GetAuditLogs)

		// Admin Waitlist Management
		admin.GET("/waitlist", middleware.RequirePermission("manage_waitlist"), waitlistCtrl.GetAllWaitlist)
		admin.POST("/waitlist/:id/approve", middleware.RequirePermission("manage_waitlist"), waitlistCtrl.ApproveWaitlistEntry)
		admin.DELETE("/waitlist/:id/reject", middleware.RequirePermission("manage_waitlist"), waitlistCtrl.RejectWaitlistEntry)

		// Admin Team Management
		admin.POST("/teams", middleware.RequirePermission("manage_teams"), teamCtrl.CreateTeam)
		admin.GET("/teams", middleware.RequirePermission("manage_teams"), teamCtrl.GetAllTeams)
		admin.GET("/teams/:id", middleware.RequirePermission("manage_teams"), teamCtrl.GetTeamByID)
		admin.PUT("/teams/:id", middleware.RequirePermission("manage_teams"), teamCtrl.UpdateTeam)
		admin.DELETE("/teams/:id", middleware.RequirePermission("manage_teams"), teamCtrl.DeleteTeam)
		admin.PATCH("/teams/:id/upload-logo", middleware.RequirePermission("manage_teams"), teamCtrl.UploadTeamLogo)

		// Admin Season Management
		admin.POST("/seasons", middleware.RequirePermission("manage_season"), seasonCtrl.CreateSeason)
		admin.GET("/seasons", middleware.RequirePermission("manage_season"), seasonCtrl.GetAllSeasons)
		admin.GET("/seasons/:id", middleware.RequirePermission("manage_season"), seasonCtrl.GetSeasonByID)
		admin.PUT("/seasons/:id", middleware.RequirePermission("manage_season"), seasonCtrl.UpdateSeason)
		admin.DELETE("/seasons/:id", middleware.RequirePermission("manage_season"), seasonCtrl.DeleteSeason)
		admin.PATCH("/seasons/:id/upload-logo", middleware.RequirePermission("manage_season"), seasonCtrl.UploadSeasonLogo)

		// Admin Match Management
		admin.POST("/matches", middleware.RequirePermission("manage_matches"), matchCtrl.CreateMatch)
		admin.GET("/matches", middleware.RequirePermission("manage_matches"), matchCtrl.GetAllMatches)
		admin.GET("/matches/:id", middleware.RequirePermission("manage_matches"), matchCtrl.GetMatchByID)
		admin.PUT("/matches/:id", middleware.RequirePermission("manage_matches"), matchCtrl.UpdateMatch)
		admin.DELETE("/matches/:id", middleware.RequirePermission("manage_matches"), matchCtrl.DeleteMatch)
		admin.PATCH("/matches/:id/upload-video", middleware.RequirePermission("upload_video"), matchCtrl.UploadMatchVideo)
		admin.GET("/matches/:id/scout/preview", middleware.RequirePermission("upload_scout"), matchCtrl.PreviewScoutMetadata)
		admin.PATCH("/matches/:id/upload-scout", middleware.RequirePermission("upload_scout"), matchCtrl.UploadMatchScout)
	}

	// AdminOrSelf routes
	user := auth.Group("/users")
	user.Use(middleware.AdminOrSelf())
	{
		user.PUT("/:id/update", userCtrl.UpdateUserProfile)
		user.DELETE("/:id/delete", userCtrl.DeleteUserAccount)
	}
}
