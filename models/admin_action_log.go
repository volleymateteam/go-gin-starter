// models/admin_action_log.go

package models

import (
	"time"

	"github.com/google/uuid"
)

// AdminActionLog represents an admin action log entry
type AdminActionLog struct {
	ID             uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AdminID        uuid.UUID              `gorm:"type:uuid;not null"`
	ActionType     string                 `gorm:"type:varchar(50);not null"`
	TargetUserID   *uuid.UUID             `gorm:"type:uuid"`
	TargetTeamID   *uuid.UUID             `gorm:"type:uuid"`
	TargetSeasonID *uuid.UUID             `gorm:"type:uuid"`
	TargetMatchID  *uuid.UUID             `gorm:"type:uuid"`
	Metadata       map[string]interface{} `gorm:"type:jsonb"`
	CreatedAt      time.Time              `gorm:"autoCreateTime"`
}
