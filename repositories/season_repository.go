package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"

	"github.com/google/uuid"
)

// CreateSeason inserts a new season
func CreateSeason(season *models.Season) error {
	return database.DB.Create(season).Error
}

// GetAllSeasons fetches all seasons
func GetAllSeasons() ([]models.Season, error) {
	var seasons []models.Season
	if err := database.DB.Find(&seasons).Error; err != nil {
		return nil, err
	}
	return seasons, nil
}

// GetSeasonByID fetches a single season by ID
func GetSeasonByID(id uuid.UUID) (*models.Season, error) {
	var season models.Season
	if err := database.DB.First(&season, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &season, nil
}

// UpdateSeason updates an existing season
func UpdateSeason(season *models.Season) error {
	return database.DB.Save(season).Error
}

// DeleteSeason soft deletes a season
func DeleteSeason(id uuid.UUID) error {
	return database.DB.Delete(&models.Season{}, "id = ?", id).Error
}
