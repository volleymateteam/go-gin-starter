package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                   uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username             string     `gorm:"unique;not null"`
	Email                string     `gorm:"unique;not null"`
	Password             string     `gorm:"not null"`
	Avatar               string     `gorm:"default:'default_avatar.png'"`
	Gender               GenderEnum `gorm:"type:varchar(10)"`
	Role                 RoleEnum   `gorm:"type:varchar(20);default:'player'"`
	ExtraPermissions     []string   `gorm:"type:jsonb;default:'[]'" json:"extra_permissions"`
	ResetPasswordToken   *string    `gorm:"type:text"`
	ResetPasswordExpires *time.Time `gorm:"type:timestamp"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}
