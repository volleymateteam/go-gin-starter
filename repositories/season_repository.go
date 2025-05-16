package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"

	"github.com/google/uuid"
)

// SeasonRepository defines the interface for season data operations
type SeasonRepository interface {
	Create(season *models.Season) error
	GetAll() ([]models.Season, error)
	GetByID(id uuid.UUID) (*models.Season, error)
	Update(season *models.Season) error
	Delete(id uuid.UUID) error
}

// GormSeasonRepository implements SeasonRepository using GORM
type GormSeasonRepository struct{}

// NewSeasonRepository creates a new instance of SeasonRepository
func NewSeasonRepository() SeasonRepository {
	return &GormSeasonRepository{}
}

// Create inserts a new season
func (r *GormSeasonRepository) Create(season *models.Season) error {
	return database.DB.Create(season).Error
}

// GetAll fetches all seasons
func (r *GormSeasonRepository) GetAll() ([]models.Season, error) {
	var seasons []models.Season
	if err := database.DB.Find(&seasons).Error; err != nil {
		return nil, err
	}
	return seasons, nil
}

// GetByID fetches a single season by ID
func (r *GormSeasonRepository) GetByID(id uuid.UUID) (*models.Season, error) {
	var season models.Season
	if err := database.DB.First(&season, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &season, nil
}

// Update updates an existing season
func (r *GormSeasonRepository) Update(season *models.Season) error {
	return database.DB.Save(season).Error
}

// Delete soft deletes a season
func (r *GormSeasonRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&models.Season{}, "id = ?", id).Error
}

// Legacy functions for backward compatibility
// These will be removed once migration is complete

func CreateSeason(season *models.Season) error {
	repo := NewSeasonRepository()
	return repo.Create(season)
}

func GetAllSeasons() ([]models.Season, error) {
	repo := NewSeasonRepository()
	return repo.GetAll()
}

func GetSeasonByID(id uuid.UUID) (*models.Season, error) {
	repo := NewSeasonRepository()
	return repo.GetByID(id)
}

func UpdateSeason(season *models.Season) error {
	repo := NewSeasonRepository()
	return repo.Update(season)
}

func DeleteSeason(id uuid.UUID) error {
	repo := NewSeasonRepository()
	return repo.Delete(id)
}
