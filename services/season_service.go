package services

import (
	"errors"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/repositories"
	"go-gin-starter/utils"

	"github.com/google/uuid"
)

// CreateSeason creates a new season
func CreateSeason(input *dto.CreateSeasonInput) (*models.Season, error) {
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
	return &season, nil
}

// GetAllSeasons returns all seasons
func GetAllSeasons() ([]models.Season, error) {
	return repositories.GetAllSeasons()
}

// GetSeasonByID returns a specific season by ID
func GetSeasonByID(id uuid.UUID) (*models.Season, error) {
	return repositories.GetSeasonByID(id)
}

// UpdateSeason updates an existing season
func UpdateSeason(id uuid.UUID, input *dto.UpdateSeasonInput) (*models.Season, error) {
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
	return season, nil
}

// DeleteSeason removes a season
func DeleteSeason(id uuid.UUID) error {
	return repositories.DeleteSeason(id)
}
