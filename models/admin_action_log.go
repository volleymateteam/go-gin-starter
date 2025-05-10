// models/admin_action_log.go
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type JSONBMap map[string]interface{}

func (j JSONBMap) Value() (driver.Value, error) {
	bytes, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (j *JSONBMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONBMap)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to type assert JSONBMap to []byte")
	}

	return json.Unmarshal(bytes, j)
}

// AdminActionLog represents an admin action log entry
type AdminActionLog struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AdminID        uuid.UUID  `gorm:"type:uuid;not null"`
	ActionType     string     `gorm:"type:varchar(50);not null"`
	TargetUserID   *uuid.UUID `gorm:"type:uuid"`
	TargetTeamID   *uuid.UUID `gorm:"type:uuid"`
	TargetSeasonID *uuid.UUID `gorm:"type:uuid"`
	TargetMatchID  *uuid.UUID `gorm:"type:uuid"`
	Metadata       JSONBMap   `gorm:"type:jsonb"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
}
