package models

import (
	"time"

	"github.com/google/uuid"
)

type WaitlistEntry struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email     string    `gorm:"unique;not null"`
	Source    string    `gorm:"type:text"` // (optional) "landing_page", "mobile_app", etc.
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
