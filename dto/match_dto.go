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
	ScoutJSON  string           `json:"scout_json_url" binding:"omitempty"`
}

type MatchResponse struct {
	ID             uuid.UUID         `json:"id"`
	SeasonID       uuid.UUID         `json:"season_id"`
	SeasonName     string            `json:"season_name"`
	HomeTeamID     uuid.UUID         `json:"home_team_id"`
	HomeTeamName   string            `json:"home_team_name"`
	AwayTeamID     uuid.UUID         `json:"away_team_id"`
	AwayTeamName   string            `json:"away_team_name"`
	Round          models.RoundEnum  `json:"round"`
	Location       string            `json:"location"`
	VideoURL       string            `json:"video_url"`
	VideoQualities map[string]string `json:"video_urls"`
	ThumbnailURL   string            `json:"thumbnail_url"`
	ScoutJSON      string            `json:"scout_json_url"`
	JsonData       interface{}       `json:"json_data"`
	// JsonData     map[string]interface{} `json:"json_data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MatchListResponse struct {
	ID           uuid.UUID        `json:"id"`
	SeasonID     uuid.UUID        `json:"season_id"`
	SeasonName   string           `json:"season_name"`
	HomeTeamID   uuid.UUID        `json:"home_team_id"`
	HomeTeamName string           `json:"home_team_name"`
	AwayTeamID   uuid.UUID        `json:"away_team_id"`
	AwayTeamName string           `json:"away_team_name"`
	Round        models.RoundEnum `json:"round"`
	Location     string           `json:"location"`
	VideoURL     string           `json:"video_url"`
	ScoutJSONURL string           `json:"scout_json_url"`
	JsonStatus   string           `json:"json_status"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}
