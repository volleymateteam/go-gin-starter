package dto

import (
	"go-gin-starter/models"
	"time"

	"github.com/google/uuid"
)

type CreateTeamInput struct {
	Name     string             `json:"name" binding:"required"`
	Country  models.CountryEnum `json:"country" binding:"required"`
	Gender   models.GenderEnum  `json:"gender" binding:"required,oneof=male female"`
	SeasonID uuid.UUID          `json:"season_id" binding:"required"`
}

type UpdateTeamInput struct {
	Name     string             `json:"name" binding:"omitempty"`
	Country  models.CountryEnum `json:"country" binding:"omitempty"`
	Gender   models.GenderEnum  `json:"gender" binding:"omitempty,oneof=male female"`
	SeasonID uuid.UUID          `json:"season_id" binding:"omitempty"`
}

type TeamResponse struct {
	ID        uuid.UUID          `json:"id"`
	Name      string             `json:"name"`
	Country   models.CountryEnum `json:"country"`
	SeasonID  uuid.UUID          `json:"season_id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
