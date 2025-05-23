package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Match struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SeasonID   uuid.UUID `gorm:"type:uuid;not null"`
	HomeTeamID uuid.UUID `gorm:"type:uuid;not null"`
	AwayTeamID uuid.UUID `gorm:"type:uuid;not null"`
	Round      RoundEnum `gorm:"type:varchar(30);not null"`

	Competition string     `gorm:"type:varchar(100);not null"`
	Gender      GenderEnum `gorm:"type:varchar(10);not null"`

	Location     string `gorm:"type:varchar(100)"`
	VideoURL     string `gorm:"type:text"` // optional
	ThumbnailURL string `gorm:"type:text"` // optional
	ScoutJSON    string `gorm:"type:text"` // optional

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
