package dto

import (
	"go-gin-starter/models"
	"time"

	"github.com/google/uuid"
)

// CreateSeasonInput defines the fields for creating a season
type CreateSeasonInput struct {
	Name       models.SeasonNameEnum `json:"name" binding:"required"`
	Country    models.CountryEnum    `json:"country" binding:"required"`
	Gender     models.GenderEnum     `json:"gender" binding:"required"`
	SeasonType models.SeasonTypeEnum `json:"season_type" binding:"required"`
	SeasonYear string                `json:"season_year" binding:"required,len=9"` // Example: "2024-2025"
	StartDate  *time.Time            `json:"start_date" binding:"omitempty"`
	EndDate    *time.Time            `json:"end_date" binding:"omitempty"`
}

// UpdateSeasonInput defines the fields for updating a season
type UpdateSeasonInput struct {
	Name       models.SeasonNameEnum `json:"name" binding:"omitempty"`
	Country    models.CountryEnum    `json:"country" binding:"omitempty"`
	Gender     models.GenderEnum     `json:"gender" binding:"omitempty"`
	SeasonType models.SeasonTypeEnum `json:"season_type" binding:"omitempty"`
	SeasonYear string                `json:"season_year" binding:"omitempty,len=9"`
	StartDate  *time.Time            `json:"start_date" binding:"omitempty"`
	EndDate    *time.Time            `json:"end_date" binding:"omitempty"`
}

// SeasonResponse defines the structure for returning a season
type SeasonResponse struct {
	ID         uuid.UUID             `json:"id"`
	Name       models.SeasonNameEnum `json:"name"`
	Country    models.CountryEnum    `json:"country"`
	Gender     models.GenderEnum     `json:"gender"`
	SeasonType models.SeasonTypeEnum `json:"season_type"`
	SeasonYear string                `json:"season_year"`
	StartDate  *time.Time            `json:"start_date,omitempty"`
	EndDate    *time.Time            `json:"end_date,omitempty"`
	LogoURL    string                `json:"logo_url"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}
