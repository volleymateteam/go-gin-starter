package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"

	"github.com/google/uuid"
)

func CreateTeam(team *models.Team) error {
	return database.DB.Create(team).Error
}

func GetAllTeams() ([]models.Team, error) {
	var teams []models.Team
	err := database.DB.Find(&teams).Error
	return teams, err
}

func GetTeamByID(id uuid.UUID) (*models.Team, error) {
	var team models.Team
	err := database.DB.First(&team, id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func UpdateTeam(team *models.Team) error {
	return database.DB.Save(team).Error
}

func DeleteTeam(id uuid.UUID) error {
	return database.DB.Delete(&models.Team{}, "id = ?", id).Error
}
