package services

import (
	"errors"
	"fmt"
	"go-gin-starter/config"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/pkg/constants"
	storagePkg "go-gin-starter/pkg/storage"
	"go-gin-starter/pkg/upload"
	"go-gin-starter/repositories"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
)

// TeamService defines the interface for team-related business logic
type TeamService interface {
	CreateTeam(input *dto.CreateTeamInput) (*models.Team, error)
	GetAllTeams() ([]dto.TeamResponse, error)
	GetTeamByID(id uuid.UUID) (*dto.TeamResponse, error)
	UpdateTeam(id uuid.UUID, input *dto.UpdateTeamInput) (*dto.TeamResponse, error)
	DeleteTeam(id uuid.UUID) error
	UpdateTeamLogo(id uuid.UUID, logoFilename string) error

	// Deprecated: Use FileUploadService directly from controllers instead.
	UploadAndSaveTeamLogo(teamID uuid.UUID, file *multipart.FileHeader) (string, string, error)
}

// TeamServiceImpl implements TeamService
type TeamServiceImpl struct {
	teamRepo      repositories.TeamRepository
	uploadService upload.FileUploadService
}

// NewTeamService creates a new instance of TeamService
func NewTeamService(teamRepo repositories.TeamRepository, uploadService upload.FileUploadService) TeamService {
	return &TeamServiceImpl{
		teamRepo:      teamRepo,
		uploadService: uploadService,
	}
}

// CreateTeam handles creation of a new team
func (s *TeamServiceImpl) CreateTeam(input *dto.CreateTeamInput) (*models.Team, error) {
	team := models.Team{
		Name:     input.Name,
		Country:  input.Country,
		Gender:   input.Gender,
		SeasonID: input.SeasonID,
		Logo:     "defaults/default-team-logo.png",
	}

	err := s.teamRepo.Create(&team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// GetAllTeams fetches all teams
func (s *TeamServiceImpl) GetAllTeams() ([]dto.TeamResponse, error) {
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.TeamResponse
	for _, team := range teams {
		responses = append(responses, s.mapTeamToResponse(&team))
	}

	return responses, nil
}

// GetTeamByID fetches a team by ID
func (s *TeamServiceImpl) GetTeamByID(id uuid.UUID) (*dto.TeamResponse, error) {
	team, err := s.teamRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(constants.ErrTeamNotFound)
	}
	response := s.mapTeamToResponse(team)
	return &response, nil
}

// UpdateTeam updates a team
func (s *TeamServiceImpl) UpdateTeam(id uuid.UUID, input *dto.UpdateTeamInput) (*dto.TeamResponse, error) {
	team, err := s.teamRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(constants.ErrTeamNotFound)
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

	if err := s.teamRepo.Update(team); err != nil {
		return nil, err
	}

	response := s.mapTeamToResponse(team)
	return &response, nil
}

// DeleteTeam deletes a team by ID
func (s *TeamServiceImpl) DeleteTeam(id uuid.UUID) error {
	return s.teamRepo.Delete(id)
}

// UpdateTeamLogo updates the logo of a team
func (s *TeamServiceImpl) UpdateTeamLogo(id uuid.UUID, logoFilename string) error {
	team, err := s.teamRepo.GetByID(id)
	if err != nil {
		return errors.New(constants.ErrTeamNotFound)
	}

	team.Logo = logoFilename

	return s.teamRepo.Update(team)
}

// UploadAndSaveTeamLogo handles validation + saving + DB update
// Deprecated: This method is kept for backward compatibility but should not be used.
// Use the FileUploadService directly from controllers instead.
func (s *TeamServiceImpl) UploadAndSaveTeamLogo(teamID uuid.UUID, file *multipart.FileHeader) (string, string, error) {
	// Get team before updating
	team, err := s.teamRepo.GetByID(teamID)
	if err != nil {
		return "", "", errors.New(constants.ErrTeamNotFound)
	}

	// File validation and upload is now handled by the upload service
	// This method is kept for backward compatibility but should be deprecated
	// in favor of directly using the upload service from the controller

	ext := filepath.Ext(file.Filename)
	newFileName := storagePkg.GenerateTeamLogoFileName(ext)
	savePath := storagePkg.BuildTeamLogoPath(newFileName)

	team.Logo = newFileName
	if err := s.teamRepo.Update(team); err != nil {
		return "", "", err
	}

	return newFileName, savePath, nil
}

// Helper function to map team model to response DTO
func (s *TeamServiceImpl) mapTeamToResponse(team *models.Team) dto.TeamResponse {
	logo := team.Logo
	if logo == "" {
		logo = "defaults/default-team-logo.png"
	}

	return dto.TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		Country:   team.Country,
		SeasonID:  team.SeasonID,
		LogoURL:   fmt.Sprintf("https://%s/logos/teams/%s", config.AssetCloudFrontDomain, logo),
		CreatedAt: team.CreatedAt,
		UpdatedAt: team.UpdatedAt,
	}
}
