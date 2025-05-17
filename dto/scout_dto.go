package dto

import "time"

type ScoutMetadataResponse struct {
	Competition string    `json:"competition"`     // e.g., "Bundesliga Men"
	Season      string    `json:"season"`          // e.g., "2024-2025"
	HomeTeam    string    `json:"home_team"`       // e.g., "Bitterfeld"
	AwayTeam    string    `json:"away_team"`       // e.g., "BR Volley"
	HomeScore   int       `json:"home_score"`      // e.g., 2
	AwayScore   int       `json:"away_score"`      // e.g., 3
	MatchDate   time.Time `json:"match_date"`      // e.g., "2024-02-07T14:00:00Z"
	Location    string    `json:"location"`        // Optional
	Round       string    `json:"round,omitempty"` // Optional
}
