package services

import (
	"errors"
	"fmt"
	"go-gin-starter/config"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/pkg/constants"
	"go-gin-starter/pkg/upload"
	"go-gin-starter/repositories"
	"mime/multipart"

	"github.com/google/uuid"
)

// SeasonService defines the interface for season-related business logic
type SeasonService interface {
	CreateSeason(input *dto.CreateSeasonInput) (*dto.SeasonResponse, error)
	GetAllSeasons() ([]dto.SeasonResponse, error)
	GetSeasonByID(id uuid.UUID) (*dto.SeasonResponse, error)
	UpdateSeason(id uuid.UUID, input *dto.UpdateSeasonInput) (*dto.SeasonResponse, error)
	DeleteSeason(id uuid.UUID) error
	UpdateSeasonLogo(seasonID uuid.UUID, logoFilename string) error

	// Deprecated: Use FileUploadService directly from controllers instead.
	UploadAndSaveSeasonLogo(seasonID uuid.UUID, file *multipart.FileHeader) (string, error)
}

// SeasonServiceImpl implements SeasonService
type SeasonServiceImpl struct {
	seasonRepo    repositories.SeasonRepository
	uploadService upload.FileUploadService
}

// NewSeasonService creates a new instance of SeasonService
func NewSeasonService(seasonRepo repositories.SeasonRepository, uploadService upload.FileUploadService) SeasonService {
	return &SeasonServiceImpl{
		seasonRepo:    seasonRepo,
		uploadService: uploadService,
	}
}

// CreateSeason creates a new season
func (s *SeasonServiceImpl) CreateSeason(input *dto.CreateSeasonInput) (*dto.SeasonResponse, error) {
	season := models.Season{
		Name:       input.Name,
		Country:    input.Country,
		Gender:     input.Gender,
		SeasonType: input.SeasonType,
		SeasonYear: input.SeasonYear,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		Logo:       "defaults/default-season-logo.png",
	}

	if err := s.seasonRepo.Create(&season); err != nil {
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
		LogoURL:    fmt.Sprintf("https://%s/logos/seasons/%s", config.AssetCloudFrontDomain, season.Logo),
		CreatedAt:  season.CreatedAt,
		UpdatedAt:  season.UpdatedAt,
	}
	return &response, nil
}

// GetAllSeasons returns all seasons
func (s *SeasonServiceImpl) GetAllSeasons() ([]dto.SeasonResponse, error) {
	seasons, err := s.seasonRepo.GetAll()
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
			LogoURL:    fmt.Sprintf("https://%s/logos/seasons/%s", config.AssetCloudFrontDomain, season.Logo),
			CreatedAt:  season.CreatedAt,
			UpdatedAt:  season.UpdatedAt,
		})
	}
	return responses, nil
}

// GetSeasonByID returns a specific season by ID
func (s *SeasonServiceImpl) GetSeasonByID(id uuid.UUID) (*dto.SeasonResponse, error) {
	season, err := s.seasonRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(constants.ErrSeasonNotFound)
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
		LogoURL:    fmt.Sprintf("https://%s/logos/seasons/%s", config.AssetCloudFrontDomain, season.Logo),
		CreatedAt:  season.CreatedAt,
		UpdatedAt:  season.UpdatedAt,
	}
	return &response, nil
}

// UpdateSeason updates an existing season
func (s *SeasonServiceImpl) UpdateSeason(id uuid.UUID, input *dto.UpdateSeasonInput) (*dto.SeasonResponse, error) {
	season, err := s.seasonRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(constants.ErrSeasonNotFound)
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

	if err := s.seasonRepo.Update(season); err != nil {
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
		LogoURL:    fmt.Sprintf("https://%s/logos/seasons/%s", config.AssetCloudFrontDomain, season.Logo),
		CreatedAt:  season.CreatedAt,
		UpdatedAt:  season.UpdatedAt,
	}
	return &response, nil
}

// DeleteSeason removes a season
func (s *SeasonServiceImpl) DeleteSeason(id uuid.UUID) error {
	return s.seasonRepo.Delete(id)
}

// UpdateSeasonLogo updates the logo of a season
func (s *SeasonServiceImpl) UpdateSeasonLogo(seasonID uuid.UUID, logoFilename string) error {
	season, err := s.seasonRepo.GetByID(seasonID)
	if err != nil {
		return errors.New(constants.ErrSeasonNotFound)
	}

	season.Logo = logoFilename
	return s.seasonRepo.Update(season)
}

// UploadAndSaveSeasonLogo handles uploading and saving a season logo
// Deprecated: This method is kept for backward compatibility but should not be used.
// Use the FileUploadService directly from controllers instead.
func (s *SeasonServiceImpl) UploadAndSaveSeasonLogo(seasonID uuid.UUID, file *multipart.FileHeader) (string, error) {
	// Get season before updating
	_, err := s.seasonRepo.GetByID(seasonID)
	if err != nil {
		return "", errors.New(constants.ErrSeasonNotFound)
	}

	// File validation and upload is now handled by the upload service
	// This method is kept for backward compatibility but should be deprecated
	// in favor of directly using the upload service from the controller

	return "", nil
}

// Legacy functions for backward compatibility
// These will be removed once migration is complete

func CreateSeasonService(input *dto.CreateSeasonInput) (*dto.SeasonResponse, error) {
	seasonRepo := repositories.NewSeasonRepository()
	seasonService := NewSeasonService(seasonRepo, upload.NewFileUploadService())
	return seasonService.CreateSeason(input)
}

func GetAllSeasonsService() ([]dto.SeasonResponse, error) {
	seasonRepo := repositories.NewSeasonRepository()
	seasonService := NewSeasonService(seasonRepo, upload.NewFileUploadService())
	return seasonService.GetAllSeasons()
}

func GetSeasonByIDService(id uuid.UUID) (*dto.SeasonResponse, error) {
	seasonRepo := repositories.NewSeasonRepository()
	seasonService := NewSeasonService(seasonRepo, upload.NewFileUploadService())
	return seasonService.GetSeasonByID(id)
}

func UpdateSeasonService(id uuid.UUID, input *dto.UpdateSeasonInput) (*dto.SeasonResponse, error) {
	seasonRepo := repositories.NewSeasonRepository()
	seasonService := NewSeasonService(seasonRepo, upload.NewFileUploadService())
	return seasonService.UpdateSeason(id, input)
}

func DeleteSeasonService(id uuid.UUID) error {
	seasonRepo := repositories.NewSeasonRepository()
	seasonService := NewSeasonService(seasonRepo, upload.NewFileUploadService())
	return seasonService.DeleteSeason(id)
}

func UpdateSeasonLogoService(seasonID uuid.UUID, logoFilename string) error {
	seasonRepo := repositories.NewSeasonRepository()
	seasonService := NewSeasonService(seasonRepo, upload.NewFileUploadService())
	return seasonService.UpdateSeasonLogo(seasonID, logoFilename)
}
