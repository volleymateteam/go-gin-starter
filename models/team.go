// File: models/team.go

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID       uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name     string      `gorm:"type:varchar(100);not null"`
	Country  CountryEnum `gorm:"type:varchar(50);not null"`
	Gender   GenderEnum  `gorm:"type:varchar(10);not null"`
	SeasonID uuid.UUID   `gorm:"type:uuid;not null"` // FK to Season

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
