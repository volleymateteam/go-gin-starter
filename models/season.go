package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Season struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	Name       SeasonNameEnum `gorm:"type:varchar(100);not null"`
	Country    CountryEnum    `gorm:"type:varchar(50);not null"`
	Gender     GenderEnum     `gorm:"type:varchar(10);not null"`
	SeasonType SeasonTypeEnum `gorm:"type:varchar(20);not null"`
	SeasonYear string         `gorm:"type:varchar(10);not null"` // Example: "2024-2025"

	StartDate *time.Time `gorm:"type:timestamp"`
	EndDate   *time.Time `gorm:"type:timestamp"`

	Logo string `gorm:"type:varchar(255);default:'defaults/default-season-logo.png'"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
