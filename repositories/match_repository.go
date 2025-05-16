package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"

	"github.com/google/uuid"
)

// MatchRepository defines the interface for match data operations
type MatchRepository interface {
	Create(match *models.Match) error
	GetAll() ([]models.Match, error)
	GetByID(id uuid.UUID) (*models.Match, error)
	Update(match *models.Match) error
	Delete(id uuid.UUID) error
}

// GormMatchRepository implements MatchRepository using GORM
type GormMatchRepository struct{}

// NewMatchRepository creates a new instance of MatchRepository
func NewMatchRepository() MatchRepository {
	return &GormMatchRepository{}
}

// Create inserts a new match record
func (r *GormMatchRepository) Create(match *models.Match) error {
	return database.DB.Create(match).Error
}

// GetAll fetches all matches
func (r *GormMatchRepository) GetAll() ([]models.Match, error) {
	var matches []models.Match
	if err := database.DB.Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

// GetByID fetches a match by ID
func (r *GormMatchRepository) GetByID(id uuid.UUID) (*models.Match, error) {
	var match models.Match
	if err := database.DB.First(&match, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

// Update updates an existing match
func (r *GormMatchRepository) Update(match *models.Match) error {
	return database.DB.Save(match).Error
}

// Delete soft deletes a match by ID
func (r *GormMatchRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&models.Match{}, "id = ?", id).Error
}

// Legacy functions for backward compatibility
// These will be removed once migration is complete

func CreateMatch(match *models.Match) error {
	repo := NewMatchRepository()
	return repo.Create(match)
}

func GetAllMatches() ([]models.Match, error) {
	repo := NewMatchRepository()
	return repo.GetAll()
}

func GetMatchByID(id uuid.UUID) (*models.Match, error) {
	repo := NewMatchRepository()
	return repo.GetByID(id)
}

func UpdateMatch(match *models.Match) error {
	repo := NewMatchRepository()
	return repo.Update(match)
}

func DeleteMatch(id uuid.UUID) error {
	repo := NewMatchRepository()
	return repo.Delete(id)
}
