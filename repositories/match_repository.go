package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"

	"github.com/google/uuid"
)

// CreateMatch inserts a new match record
func CreateMatch(match *models.Match) error {
	return database.DB.Create(match).Error
}

// GetAllMatches fetches all matches
func GetAllMatches() ([]models.Match, error) {
	var matches []models.Match
	if err := database.DB.Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

// GetMatchByID fetches a match by ID
func GetMatchByID(id uuid.UUID) (*models.Match, error) {
	var match models.Match
	if err := database.DB.First(&match, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

// UpdateMatch updates an existing match
func UpdateMatch(match *models.Match) error {
	return database.DB.Save(match).Error
}

// DeleteMatch soft deletes a match by ID
func DeleteMatch(id uuid.UUID) error {
	return database.DB.Delete(&models.Match{}, "id = ?", id).Error
}
