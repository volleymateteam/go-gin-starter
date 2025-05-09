package repositories

import (
	"fmt"
	"go-gin-starter/database"
	"go-gin-starter/models"
)

// CreateAdminActionLog inserts a new admin action log entry
func CreateAdminActionLog(log *models.AdminActionLog) error {
	// return database.DB.Create(log).Error
	result := database.DB.Create(log)
	fmt.Printf("CreateAdminActionLog SQL: %+v\n", result.Statement.SQL.String())
	fmt.Printf("CreateAdminActionLog Error: %v\n", result.Error)
	return result.Error
}

// GetAuditLogs fetches audit logs with optional actionType filter and pagination
func GetAuditLogs(actionType string, offset, limit int) ([]models.AdminActionLog, error) {
	var logs []models.AdminActionLog
	query := database.DB.Order("created_at DESC")

	if actionType != "" {
		query = query.Where("action_type = ?", actionType)
	}

	if err := query.Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}
