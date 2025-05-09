package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"
)

// CreateAdminActionLog inserts a new admin action log entry
func CreateAdminActionLog(log *models.AdminActionLog) error {
	return database.DB.Create(log).Error
}

// GetAuditLogs fetches all audit logs, newest first (limit 100 for safety)
func GetAuditLogs() ([]models.AdminActionLog, error) {
	var logs []models.AdminActionLog
	if err := database.DB.Order("created_at DESC").Limit(100).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
