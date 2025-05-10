package controllers

import (
	"go-gin-starter/database"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck provides basic health status
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// DetailedHealthCheck provides comprehensive system health information
func DetailedHealthCheck(c *gin.Context) {
	// Check database connection
	dbStatus := "ok"
	dbLatency := time.Duration(0)
	startTime := time.Now()

	sqlDB, err := database.DB.DB()
	if err != nil {
		dbStatus = "error: failed to get database connection"
	} else {
		// Ping database with timeout
		pingErr := sqlDB.Ping()
		dbLatency = time.Since(startTime)

		if pingErr != nil {
			dbStatus = "error: " + pingErr.Error()
		}
	}

	// Memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// System info
	hostname, _ := os.Hostname()

	// Build response
	healthData := gin.H{
		"status":      "ok",
		"time":        time.Now().Format(time.RFC3339),
		"environment": os.Getenv("ENV"),
		"version":     os.Getenv("APP_VERSION"),
		"hostname":    hostname,
		"system": gin.H{
			"go_version": runtime.Version(),
			"goroutines": runtime.NumGoroutine(),
			"cpu_cores":  runtime.NumCPU(),
			"memory_usage": gin.H{
				"alloc_mb":       m.Alloc / 1024 / 1024,
				"total_alloc_mb": m.TotalAlloc / 1024 / 1024,
				"sys_mb":         m.Sys / 1024 / 1024,
				"gc_cycles":      m.NumGC,
			},
		},
		"dependencies": gin.H{
			"database": gin.H{
				"status":  dbStatus,
				"latency": dbLatency.String(),
			},
		},
	}

	// If there's a database error, update overall status
	if dbStatus != "ok" {
		healthData["status"] = "degraded"
	}

	c.JSON(http.StatusOK, healthData)
}

// LivenessCheck for kubernetes liveness probe
func LivenessCheck(c *gin.Context) {
	// Only check if server is running
	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgHealthy)
}

// ReadinessCheck for kubernetes readiness probe
func ReadinessCheck(c *gin.Context) {
	// Check if application can handle traffic by testing DB connection
	sqlDB, err := database.DB.DB()
	if err != nil {
		httpPkg.RespondError(c, http.StatusServiceUnavailable, constants.ErrDatabaseConnection)
		return
	}

	// Ping database with timeout
	err = sqlDB.Ping()
	if err != nil {
		httpPkg.RespondError(c, http.StatusServiceUnavailable, constants.ErrDatabaseConnection)
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgHealthy)
}
