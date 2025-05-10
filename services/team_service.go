package services

import (
	"errors"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/repositories"
	"go-gin-starter/utils"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
)

// CreateTeamService handles creation of a new team
func CreateTeamService(input *dto.CreateTeamInput) (*models.Team, error) {
	team := models.Team{
		Name:     input.Name,
		Country:  input.Country,
		Gender:   input.Gender,
		SeasonID: input.SeasonID,
		Logo:     "defaults/default-team-logo.png",
	}

	err := repositories.CreateTeam(&team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// GetAllTeamsService fetches all teams
func GetTeamByIDService(id uuid.UUID) (*dto.TeamResponse, error) {
	team, err := repositories.GetTeamByID(id)
	if err != nil {
		return nil, errors.New(utils.ErrTeamNotFound)
	}
	response := mapTeamToResponse(team)
	return &response, nil
}

// GetTeamByIDService fetches a team by ID
func GetAllTeamsService() ([]dto.TeamResponse, error) {
	teams, err := repositories.GetAllTeams()
	if err != nil {
		return nil, err
	}

	var responses []dto.TeamResponse
	for _, team := range teams {
		responses = append(responses, mapTeamToResponse(&team))
	}

	return responses, nil
}

// UpdateTeamService updates a team
func UpdateTeamService(id uuid.UUID, input *dto.UpdateTeamInput) (*dto.TeamResponse, error) {
	team, err := repositories.GetTeamByID(id)
	if err != nil {
		return nil, errors.New(utils.ErrTeamNotFound)
	}

	if input.Name != "" {
		team.Name = input.Name
	}
	if input.Country != "" {
		team.Country = input.Country
	}
	if input.Gender != "" {
		team.Gender = input.Gender
	}
	if input.SeasonID != uuid.Nil {
		team.SeasonID = input.SeasonID
	}

	if err := repositories.UpdateTeam(team); err != nil {
		return nil, err
	}

	response := mapTeamToResponse(team)
	return &response, nil
}

// DeleteTeamService deletes a team by ID
func DeleteTeamService(id uuid.UUID) error {
	return repositories.DeleteTeam(id)
}

// UpdateTeamLogoService updates the logo of a team
func UpdateTeamLogoService(id uuid.UUID, logoFilename string) error {
	team, err := repositories.GetTeamByID(id)
	if err != nil {
		return errors.New(utils.ErrTeamNotFound)
	}

	team.Logo = logoFilename

	return repositories.UpdateTeam(team)
}

// UploadAndSaveTeamLogoService handles validation + saving + DB update
func UploadAndSaveTeamLogoService(teamID uuid.UUID, file *multipart.FileHeader) (string, string, error) {
	if err := utils.ValidateImageFile(file); err != nil {
		return "", "", err
	}

	ext := filepath.Ext(file.Filename)
	newFileName := utils.GenerateTeamLogoFileName(ext)
	savePath := utils.BuildTeamLogoPath(newFileName)

	team, err := repositories.GetTeamByID(teamID)
	if err != nil {
		return "", "", errors.New(utils.ErrTeamNotFound)
	}

	team.Logo = newFileName
	if err := repositories.UpdateTeam(team); err != nil {
		return "", "", err
	}

	return newFileName, savePath, nil
}

func mapTeamToResponse(team *models.Team) dto.TeamResponse {
	logo := team.Logo
	if logo == "" {
		logo = "defaults/default-team-logo.png"
	}

	return dto.TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		Country:   team.Country,
		SeasonID:  team.SeasonID,
		LogoURL:   "/uploads/logos/" + logo,
		CreatedAt: team.CreatedAt,
		UpdatedAt: team.UpdatedAt,
	}
}
