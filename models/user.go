package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StringArray is a custom type to handle string arrays in PostgreSQL
type StringArray []string

// Value implements the driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

// Scan implements the sql.Scanner interface
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

type User struct {
	ID                   uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username             string      `gorm:"unique;not null"`
	Email                string      `gorm:"unique;not null"`
	Password             string      `gorm:"column:hashed_password;not null"`
	Avatar               string      `gorm:"default:'default_avatar.png'"`
	Gender               GenderEnum  `gorm:"type:varchar(10)"`
	Role                 RoleEnum    `gorm:"type:varchar(20);default:'player'"`
	ExtraPermissions     StringArray `gorm:"type:jsonb;default:null" json:"extra_permissions"`
	ResetPasswordToken   *string     `gorm:"type:text"`
	ResetPasswordExpires *time.Time  `gorm:"type:timestamp"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}
