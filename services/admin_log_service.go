package services

import (
	"go-gin-starter/models"
	"go-gin-starter/repositories"

	"github.com/google/uuid"
)

// LogAdminAction logs an admin action to database
func LogAdminAction(adminID uuid.UUID, actionType string, targetUserID, targetTeamID, targetSeasonID, targetMatchID *uuid.UUID, metadata map[string]interface{}) error {
	log := &models.AdminActionLog{
		AdminID:        adminID,
		ActionType:     actionType,
		TargetUserID:   targetUserID,
		TargetTeamID:   targetTeamID,
		TargetSeasonID: targetSeasonID,
		TargetMatchID:  targetMatchID,
		Metadata:       metadata,
	}

	return repositories.CreateAdminActionLog(log)
}

// GetAuditLogs returns audit logs from repository
func GetAuditLogs() ([]models.AdminActionLog, error) {
	return repositories.GetAuditLogs()
}
