package services

import (
	"errors"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/repositories"
	"go-gin-starter/utils"

	"github.com/google/uuid"
)

// CreateSeasonService creates a new season
func CreateSeasonService(input *dto.CreateSeasonInput) (*dto.SeasonResponse, error) {
	season := models.Season{
		Name:       input.Name,
		Country:    input.Country,
		Gender:     input.Gender,
		SeasonType: input.SeasonType,
		SeasonYear: input.SeasonYear,
		Round:      input.Round,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
	}

	if err := repositories.CreateSeason(&season); err != nil {
		return nil, err
	}

	response := dto.SeasonResponse{
		ID:         season.ID,
		Name:       season.Name,
		Country:    season.Country,
		Gender:     season.Gender,
		SeasonType: season.SeasonType,
		SeasonYear: season.SeasonYear,
		StartDate:  season.StartDate,
		EndDate:    season.EndDate,
		Round:      season.Round,
		CreatedAt:  season.CreatedAt,
		UpdatedAt:  season.UpdatedAt,
	}
	return &response, nil
}

// GetAllSeasonsService returns all seasons
func GetAllSeasonsService() ([]dto.SeasonResponse, error) {
	seasons, err := repositories.GetAllSeasons()
	if err != nil {
		return nil, err
	}

	var responses []dto.SeasonResponse
	for _, season := range seasons {
		responses = append(responses, dto.SeasonResponse{
			ID:         season.ID,
			Name:       season.Name,
			Country:    season.Country,
			Gender:     season.Gender,
			SeasonType: season.SeasonType,
			SeasonYear: season.SeasonYear,
			StartDate:  season.StartDate,
			EndDate:    season.EndDate,
			Round:      season.Round,
			CreatedAt:  season.CreatedAt,
			UpdatedAt:  season.UpdatedAt,
		})
	}
	return responses, nil
}

// GetSeasonByIDService returns a specific season by ID
func GetSeasonByIDService(id uuid.UUID) (*dto.SeasonResponse, error) {
	season, err := repositories.GetSeasonByID(id)
	if err != nil {
		return nil, errors.New(utils.ErrSeasonNotFound)
	}

	response := dto.SeasonResponse{
		ID:         season.ID,
		Name:       season.Name,
		Country:    season.Country,
		Gender:     season.Gender,
		SeasonType: season.SeasonType,
		SeasonYear: season.SeasonYear,
		StartDate:  season.StartDate,
		EndDate:    season.EndDate,
		Round:      season.Round,
		CreatedAt:  season.CreatedAt,
		UpdatedAt:  season.UpdatedAt,
	}
	return &response, nil
}

// UpdateSeasonService updates an existing season
func UpdateSeasonService(id uuid.UUID, input *dto.UpdateSeasonInput) (*dto.SeasonResponse, error) {
	season, err := repositories.GetSeasonByID(id)
	if err != nil {
		return nil, errors.New(utils.ErrSeasonNotFound)
	}

	if input.Name != "" {
		season.Name = input.Name
	}
	if input.Country != "" {
		season.Country = input.Country
	}
	if input.Gender != "" {
		season.Gender = input.Gender
	}
	if input.SeasonType != "" {
		season.SeasonType = input.SeasonType
	}
	if input.SeasonYear != "" {
		season.SeasonYear = input.SeasonYear
	}
	if input.Round != "" {
		season.Round = input.Round
	}
	if input.StartDate != nil {
		season.StartDate = input.StartDate
	}
	if input.EndDate != nil {
		season.EndDate = input.EndDate
	}

	if err := repositories.UpdateSeason(season); err != nil {
		return nil, err
	}

	response := dto.SeasonResponse{
		ID:         season.ID,
		Name:       season.Name,
		Country:    season.Country,
		Gender:     season.Gender,
		SeasonType: season.SeasonType,
		SeasonYear: season.SeasonYear,
		StartDate:  season.StartDate,
		EndDate:    season.EndDate,
		Round:      season.Round,
		CreatedAt:  season.CreatedAt,
		UpdatedAt:  season.UpdatedAt,
	}
	return &response, nil
}

// DeleteSeasonService removes a season
func DeleteSeasonService(id uuid.UUID) error {
	return repositories.DeleteSeason(id)
}
