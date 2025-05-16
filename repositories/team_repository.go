package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"

	"github.com/google/uuid"
)

// TeamRepository defines the interface for team data operations
type TeamRepository interface {
	Create(team *models.Team) error
	GetAll() ([]models.Team, error)
	GetByID(id uuid.UUID) (*models.Team, error)
	Update(team *models.Team) error
	Delete(id uuid.UUID) error
}

// GormTeamRepository implements TeamRepository using GORM
type GormTeamRepository struct{}

// NewTeamRepository creates a new instance of TeamRepository
func NewTeamRepository() TeamRepository {
	return &GormTeamRepository{}
}

// Create adds a new team to the database
func (r *GormTeamRepository) Create(team *models.Team) error {
	return database.DB.Create(team).Error
}

// GetAll retrieves all teams
func (r *GormTeamRepository) GetAll() ([]models.Team, error) {
	var teams []models.Team
	err := database.DB.Find(&teams).Error
	return teams, err
}

// GetByID retrieves a team by ID
func (r *GormTeamRepository) GetByID(id uuid.UUID) (*models.Team, error) {
	var team models.Team
	err := database.DB.First(&team, id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// Update modifies an existing team
func (r *GormTeamRepository) Update(team *models.Team) error {
	return database.DB.Save(team).Error
}

// Delete removes a team by ID
func (r *GormTeamRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&models.Team{}, "id = ?", id).Error
}

// Legacy functions for backward compatibility
// These will be removed once migration is complete

func CreateTeam(team *models.Team) error {
	repo := NewTeamRepository()
	return repo.Create(team)
}

func GetAllTeams() ([]models.Team, error) {
	repo := NewTeamRepository()
	return repo.GetAll()
}

func GetTeamByID(id uuid.UUID) (*models.Team, error) {
	repo := NewTeamRepository()
	return repo.GetByID(id)
}

func UpdateTeam(team *models.Team) error {
	repo := NewTeamRepository()
	return repo.Update(team)
}

func DeleteTeam(id uuid.UUID) error {
	repo := NewTeamRepository()
	return repo.Delete(id)
}
