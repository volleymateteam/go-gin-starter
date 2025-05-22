package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/pkg/logger"
	scoutPkg "go-gin-starter/pkg/scout"
	storagePkg "go-gin-starter/pkg/storage"
	"go-gin-starter/pkg/video"
	"go-gin-starter/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// MatchService defines the interface for match-related business logic
type MatchService interface {
	CreateMatch(input *dto.CreateMatchInput) (*dto.MatchResponse, error)
	GetAllMatches() ([]dto.MatchListResponse, error)
	GetMatchByID(id uuid.UUID) (*dto.MatchResponse, error)
	UpdateMatch(id uuid.UUID, input *dto.UpdateMatchInput) (*dto.MatchResponse, error)
	DeleteMatch(id uuid.UUID) error
	UploadMatchVideo(matchID uuid.UUID, file io.Reader, fileHeader *multipart.FileHeader) (string, error)
	UploadMatchScout(matchID uuid.UUID, file io.Reader, fileHeader *multipart.FileHeader) (string, error)
}

// MatchServiceImpl implements MatchService
type MatchServiceImpl struct {
	matchRepo  repositories.MatchRepository
	teamRepo   repositories.TeamRepository
	seasonRepo repositories.SeasonRepository
	videoQueue *video.QueueManager
}

// NewMatchService creates a new instance of MatchService
func NewMatchService(
	matchRepo repositories.MatchRepository,
	teamRepo repositories.TeamRepository,
	seasonRepo repositories.SeasonRepository,
	videoQueue *video.QueueManager,
) MatchService {
	return &MatchServiceImpl{
		matchRepo:  matchRepo,
		teamRepo:   teamRepo,
		seasonRepo: seasonRepo,
		videoQueue: videoQueue,
	}
}

// CreateMatch handles creation of a new match
func (s *MatchServiceImpl) CreateMatch(input *dto.CreateMatchInput) (*dto.MatchResponse, error) {
	match := models.Match{
		SeasonID:   input.SeasonID,
		HomeTeamID: input.HomeTeamID,
		AwayTeamID: input.AwayTeamID,
		Round:      input.Round,
		Location:   input.Location,
	}

	season, err := s.seasonRepo.GetByID(match.SeasonID)
	if err != nil {
		return nil, errors.New("season not found")
	}

	match.Gender = season.Gender
	match.Competition = string(season.Name)

	if err := s.matchRepo.Create(&match); err != nil {
		return nil, err
	}

	// Fetch related names
	homeTeam, _ := s.teamRepo.GetByID(match.HomeTeamID)
	awayTeam, _ := s.teamRepo.GetByID(match.AwayTeamID)

	response := dto.MatchResponse{
		ID:           match.ID,
		SeasonID:     match.SeasonID,
		SeasonName:   s.getSeasonName(season),
		HomeTeamID:   match.HomeTeamID,
		HomeTeamName: s.getTeamName(homeTeam),
		AwayTeamID:   match.AwayTeamID,
		AwayTeamName: s.getTeamName(awayTeam),
		Round:        match.Round,
		Location:     match.Location,
		VideoURL:     match.VideoURL,
		ScoutJSON:    match.ScoutJSON,
		CreatedAt:    match.CreatedAt,
		UpdatedAt:    match.UpdatedAt,
	}

	return &response, nil
}

// GetAllMatches returns all matches
func (s *MatchServiceImpl) GetAllMatches() ([]dto.MatchListResponse, error) {
	matches, err := s.matchRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.MatchListResponse
	for _, m := range matches {
		season, _ := s.seasonRepo.GetByID(m.SeasonID)
		homeTeam, _ := s.teamRepo.GetByID(m.HomeTeamID)
		awayTeam, _ := s.teamRepo.GetByID(m.AwayTeamID)

		status := "missing"
		if m.ScoutJSON != "" {
			status = "available"
		}

		responses = append(responses, dto.MatchListResponse{
			ID:           m.ID,
			SeasonID:     m.SeasonID,
			SeasonName:   s.getSeasonName(season),
			HomeTeamID:   m.HomeTeamID,
			HomeTeamName: s.getTeamName(homeTeam),
			AwayTeamID:   m.AwayTeamID,
			AwayTeamName: s.getTeamName(awayTeam),
			Round:        m.Round,
			Location:     m.Location,
			VideoURL:     m.VideoURL,
			ScoutJSONURL: m.ScoutJSON,
			JsonStatus:   status,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
		})
	}

	return responses, nil
}

// GetMatchByID returns a single match by ID
func (s *MatchServiceImpl) GetMatchByID(id uuid.UUID) (*dto.MatchResponse, error) {
	match, err := s.matchRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(constants.ErrMatchNotFound)
	}

	season, _ := s.seasonRepo.GetByID(match.SeasonID)
	homeTeam, _ := s.teamRepo.GetByID(match.HomeTeamID)
	awayTeam, _ := s.teamRepo.GetByID(match.AwayTeamID)

	// Fetch and parse JSON from S3
	var jsonData map[string]interface{}
	if match.ScoutJSON != "" {
		jsonData, _ = fetchJSONFromS3(match.ScoutJSON)
	}

	return &dto.MatchResponse{
		ID:           match.ID,
		SeasonID:     match.SeasonID,
		SeasonName:   s.getSeasonName(season),
		HomeTeamID:   match.HomeTeamID,
		HomeTeamName: s.getTeamName(homeTeam),
		AwayTeamID:   match.AwayTeamID,
		AwayTeamName: s.getTeamName(awayTeam),
		Round:        match.Round,
		Location:     match.Location,
		VideoURL:     match.VideoURL,
		ThumbnailURL: match.ThumbnailURL,
		ScoutJSON:    match.ScoutJSON,
		JsonData:     jsonData,
		CreatedAt:    match.CreatedAt,
		UpdatedAt:    match.UpdatedAt,
	}, nil
}

// UpdateMatch updates an existing match
func (s *MatchServiceImpl) UpdateMatch(id uuid.UUID, input *dto.UpdateMatchInput) (*dto.MatchResponse, error) {
	match, err := s.matchRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(constants.ErrMatchNotFound)
	}

	if input.HomeTeamID != uuid.Nil {
		match.HomeTeamID = input.HomeTeamID
	}
	if input.AwayTeamID != uuid.Nil {
		match.AwayTeamID = input.AwayTeamID
	}
	if input.Round != "" {
		match.Round = input.Round
	}
	if input.Location != "" {
		match.Location = input.Location
	}
	if input.VideoURL != "" {
		match.VideoURL = input.VideoURL
	}
	if input.ScoutJSON != "" {
		match.ScoutJSON = input.ScoutJSON
	}

	if err := s.matchRepo.Update(match); err != nil {
		return nil, err
	}

	season, _ := s.seasonRepo.GetByID(match.SeasonID)
	homeTeam, _ := s.teamRepo.GetByID(match.HomeTeamID)
	awayTeam, _ := s.teamRepo.GetByID(match.AwayTeamID)

	return &dto.MatchResponse{
		ID:           match.ID,
		SeasonID:     match.SeasonID,
		SeasonName:   s.getSeasonName(season),
		HomeTeamID:   match.HomeTeamID,
		HomeTeamName: s.getTeamName(homeTeam),
		AwayTeamID:   match.AwayTeamID,
		AwayTeamName: s.getTeamName(awayTeam),
		Round:        match.Round,
		Location:     match.Location,
		VideoURL:     match.VideoURL,
		ScoutJSON:    match.ScoutJSON,
		CreatedAt:    match.CreatedAt,
		UpdatedAt:    match.UpdatedAt,
	}, nil
}

// DeleteMatch deletes a match
func (s *MatchServiceImpl) DeleteMatch(id uuid.UUID) error {
	return s.matchRepo.Delete(id)
}

// UploadMatchVideo handles uploading a match video to S3
func (s *MatchServiceImpl) UploadMatchVideo(
	matchID uuid.UUID,
	file io.Reader,
	fileHeader *multipart.FileHeader,
) (string, error) {
	match, err := s.matchRepo.GetByID(matchID)
	if err != nil {
		return "", errors.New(constants.ErrMatchNotFound)
	}

	season, err := s.seasonRepo.GetByID(match.SeasonID)
	if err != nil {
		return "", errors.New(constants.ErrSeasonNotFound)
	}

	// Create the folder structure
	safeSeasonName := strings.ReplaceAll(strings.ToLower(string(season.Name)), " ", "_")
	safeSeasonYear := strings.ReplaceAll(season.SeasonYear, "/", "_")
	safeGender := strings.ToLower(string(season.Gender))
	safeCountry := strings.ToLower(string(season.Country))

	// Generate paths
	basePath := fmt.Sprintf("videos/%s_%s/%s_%s/%s",
		safeSeasonYear,
		safeCountry,
		safeSeasonName,
		safeGender,
		matchID.String())

	logger.Info("Upload path",
		zap.String("rawKey", basePath),
		zap.String("compressedKey", basePath))

	rawKey := fmt.Sprintf("%s/%s/%s%s",
		basePath,
		video.RawVideoFolder,
		uuid.New().String(),
		filepath.Ext(fileHeader.Filename))

	compressedKey := fmt.Sprintf("%s/%s/%s.mp4",
		basePath,
		video.CompressedFolder,
		uuid.New().String())

	// Upload raw video to S3
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	_, err = storagePkg.UploadBytesToS3(buf.Bytes(), rawKey, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}

	// Create and enqueue processing job
	job := &video.VideoProcessingJob{
		MatchID:   matchID.String(),
		InputKey:  rawKey,
		OutputKey: compressedKey,
	}

	if err := s.videoQueue.EnqueueVideo(job); err != nil {
		logger.Error("Failed to enqueue video job",
			zap.String("match_id", matchID.String()),
			zap.Error(err))
	}

	// build CloudFront compressed video URL
	compressedURL := fmt.Sprintf("https://%s/%s", os.Getenv("VIDEO_CLOUDFRONT_DOMAIN"), compressedKey)

	// Save the compressed URL instead of raw URL
	match.VideoURL = compressedURL
	if err := s.matchRepo.Update(match); err != nil {
		return "", err
	}

	return compressedURL, nil
}

// UploadMatchScout handles uploading and processing a match scout file
func (s *MatchServiceImpl) UploadMatchScout(matchID uuid.UUID,
	file io.Reader,
	fileHeader *multipart.FileHeader,
) (string, error) {
	match, err := s.matchRepo.GetByID(matchID)
	if err != nil {
		return "", errors.New(constants.ErrMatchNotFound)
	}

	if filepath.Ext(fileHeader.Filename) != ".dvw" {
		return "", errors.New("invalid file type: only .dvw supported")
	}

	// Read file into memory
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	// Upload original .dvw file to S3
	s3InputKey := fmt.Sprintf("scouts/%s.dvw", matchID.String())
	contentType := fileHeader.Header.Get("Content-Type")

	_, err = storagePkg.UploadBytesToS3(buf.Bytes(), s3InputKey, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload .dvw file: %w", err)
	}

	// Call Python microservice to parse the uploaded file
	parsedResult, err := scoutPkg.CallPythonParser(s3InputKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse scout file: %w", err)
	}

	// Convert parsed JSON to []byte
	jsonBytes, err := json.Marshal(parsedResult.JsonData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal parsed data: %w", err)
	}

	// Upload parsed JSON to S3 using the scout CloudFront domain
	s3OutputKey := fmt.Sprintf("scout-files/%s.json", matchID.String())
	jsonURL, err := storagePkg.UploadBytesToS3(jsonBytes, s3OutputKey, "application/json")
	if err != nil {
		return "", fmt.Errorf("failed to upload scout json: %w", err)
	}

	// Save JSON URL to DB
	match.ScoutJSON = jsonURL
	if err := s.matchRepo.Update(match); err != nil {
		return "", err
	}

	return jsonURL, nil
}

// Helper function to fetch JSON from S3
func fetchJSONFromS3(url string) (map[string]interface{}, error) {
	if url == "" {
		return nil, nil
	}

	// Use the existing HTTP utility to fetch and parse JSON
	jsonData, err := httpPkg.FetchJSONFromS3(url)
	if err != nil {
		logger.Error("Failed to fetch JSON from S3", zap.Error(err), zap.String("url", url))
		return nil, err
	}

	return jsonData, nil
}

// Helper to format season name
func (s *MatchServiceImpl) getSeasonName(season *models.Season) string {
	if season == nil {
		return ""
	}
	return string(season.Name) + " " + season.SeasonYear
}

// Helper to get team name
func (s *MatchServiceImpl) getTeamName(team *models.Team) string {
	if team == nil {
		return ""
	}
	return team.Name
}
