package services

import (
	"errors"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/repositories"
	"go-gin-starter/utils"

	"github.com/google/uuid"
)

// CreateMatchService handles creation of a new match
func CreateMatchService(input *dto.CreateMatchInput) (*dto.MatchResponse, error) {
	match := models.Match{
		SeasonID:   input.SeasonID,
		HomeTeamID: input.HomeTeamID,
		AwayTeamID: input.AwayTeamID,
		Round:      input.Round,
		Location:   input.Location,
	}

	if err := repositories.CreateMatch(&match); err != nil {
		return nil, err
	}

	response := dto.MatchResponse{
		ID:         match.ID,
		SeasonID:   match.SeasonID,
		HomeTeamID: match.HomeTeamID,
		AwayTeamID: match.AwayTeamID,
		Round:      match.Round,
		Location:   match.Location,
		VideoURL:   match.VideoURL,
		ScoutJSON:  match.ScoutJSON,
		CreatedAt:  match.CreatedAt,
		UpdatedAt:  match.UpdatedAt,
	}

	return &response, nil
}

// GetAllMatchesService returns all matches
func GetAllMatchesService() ([]dto.MatchResponse, error) {
	matches, err := repositories.GetAllMatches()
	if err != nil {
		return nil, err
	}

	var responses []dto.MatchResponse
	for _, m := range matches {
		responses = append(responses, dto.MatchResponse{
			ID:         m.ID,
			SeasonID:   m.SeasonID,
			HomeTeamID: m.HomeTeamID,
			AwayTeamID: m.AwayTeamID,
			Round:      m.Round,
			Location:   m.Location,
			VideoURL:   m.VideoURL,
			ScoutJSON:  m.ScoutJSON,
			CreatedAt:  m.CreatedAt,
			UpdatedAt:  m.UpdatedAt,
		})
	}

	return responses, nil
}

// GetMatchByIDService returns a single match by ID
func GetMatchByIDService(id uuid.UUID) (*dto.MatchResponse, error) {
	match, err := repositories.GetMatchByID(id)
	if err != nil {
		return nil, errors.New(utils.ErrMatchNotFound)
	}

	return &dto.MatchResponse{
		ID:         match.ID,
		SeasonID:   match.SeasonID,
		HomeTeamID: match.HomeTeamID,
		AwayTeamID: match.AwayTeamID,
		Round:      match.Round,
		Location:   match.Location,
		VideoURL:   match.VideoURL,
		ScoutJSON:  match.ScoutJSON,
		CreatedAt:  match.CreatedAt,
		UpdatedAt:  match.UpdatedAt,
	}, nil
}

// UpdateMatchService updates an existing match
func UpdateMatchService(id uuid.UUID, input *dto.UpdateMatchInput) (*dto.MatchResponse, error) {
	match, err := repositories.GetMatchByID(id)
	if err != nil {
		return nil, errors.New(utils.ErrMatchNotFound)
	}

	if input.HomeTeamID != uuid.Nil {
		match.HomeTeamID = input.HomeTeamID
	}
	if input.AwayTeamID != uuid.Nil {
		match.AwayTeamID = input.AwayTeamID
	}
	if input.Round != "" {
		match.Round = input.Round
	}
	if input.Location != "" {
		match.Location = input.Location
	}
	if input.VideoURL != "" {
		match.VideoURL = input.VideoURL
	}
	if input.ScoutJSON != "" {
		match.ScoutJSON = input.ScoutJSON
	}

	if err := repositories.UpdateMatch(match); err != nil {
		return nil, err
	}

	return &dto.MatchResponse{
		ID:         match.ID,
		SeasonID:   match.SeasonID,
		HomeTeamID: match.HomeTeamID,
		AwayTeamID: match.AwayTeamID,
		Round:      match.Round,
		Location:   match.Location,
		VideoURL:   match.VideoURL,
		ScoutJSON:  match.ScoutJSON,
		CreatedAt:  match.CreatedAt,
		UpdatedAt:  match.UpdatedAt,
	}, nil
}

// DeleteMatchService deletes a match
func DeleteMatchService(id uuid.UUID) error {
	return repositories.DeleteMatch(id)
}
