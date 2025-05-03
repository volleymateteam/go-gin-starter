package dto

import (
	"go-gin-starter/models"
	"time"

	"github.com/google/uuid"
)

type CreateMatchInput struct {
	SeasonID   uuid.UUID        `json:"season_id" binding:"required"`
	HomeTeamID uuid.UUID        `json:"home_team_id" binding:"required"`
	AwayTeamID uuid.UUID        `json:"away_team_id" binding:"required"`
	Round      models.RoundEnum `json:"round" binding:"required"`
	Location   string           `json:"location" binding:"omitempty"`
}

type UpdateMatchInput struct {
	HomeTeamID uuid.UUID        `json:"home_team_id" binding:"omitempty"`
	AwayTeamID uuid.UUID        `json:"away_team_id" binding:"omitempty"`
	Round      models.RoundEnum `json:"round" binding:"omitempty"`
	Location   string           `json:"location" binding:"omitempty"`
	VideoURL   string           `json:"video_url" binding:"omitempty"`
	ScoutJSON  string           `json:"scout_json" binding:"omitempty"`
}

type MatchResponse struct {
	ID         uuid.UUID        `json:"id"`
	SeasonID   uuid.UUID        `json:"season_id"`
	HomeTeamID uuid.UUID        `json:"home_team_id"`
	AwayTeamID uuid.UUID        `json:"away_team_id"`
	Round      models.RoundEnum `json:"round"`
	Location   string           `json:"location"`
	VideoURL   string           `json:"video_url"`
	ScoutJSON  string           `json:"scout_json"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}
