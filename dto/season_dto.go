package dto

import (
	"go-gin-starter/models"
	"time"
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
	Round      models.RoundEnum      `json:"round" binding:"required"`
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
	Round      models.RoundEnum      `json:"round" binding:"omitempty"`
}

// SeasonResponse defines the structure for returning a season
type SeasonResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Country    string `json:"country"`
	Gender     string `json:"gender"`
	SeasonType string `json:"season_type"`
	SeasonYear string `json:"season_year"`
	Round      string `json:"round"`
	StartDate  string `json:"start_date,omitempty"`
	EndDate    string `json:"end_date,omitempty"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
