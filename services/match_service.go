package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"go-gin-starter/config"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	scoutPkg "go-gin-starter/pkg/scout"
	storagePkg "go-gin-starter/pkg/storage"
	"go-gin-starter/repositories"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

// CreateMatchService handles creation of a new match
func CreateMatchService(input *dto.CreateMatchInput) (*dto.MatchResponse, error) {
	match := models.Match{
		SeasonID:   input.SeasonID,
		HomeTeamID: input.HomeTeamID,
		AwayTeamID: input.AwayTeamID,
		Round:      input.Round,
		Location:   input.Location,
	}

	if err := repositories.CreateMatch(&match); err != nil {
		return nil, err
	}

	// Fetch related names
	season, _ := repositories.GetSeasonByID(match.SeasonID)
	homeTeam, _ := repositories.GetTeamByID(match.HomeTeamID)
	awayTeam, _ := repositories.GetTeamByID(match.AwayTeamID)

	response := dto.MatchResponse{
		ID:           match.ID,
		SeasonID:     match.SeasonID,
		SeasonName:   getSeasonName(season),
		HomeTeamID:   match.HomeTeamID,
		HomeTeamName: getTeamName(homeTeam),
		AwayTeamID:   match.AwayTeamID,
		AwayTeamName: getTeamName(awayTeam),
		Round:        match.Round,
		Location:     match.Location,
		VideoURL:     match.VideoURL,
		ScoutJSON:    match.ScoutJSON,
		CreatedAt:    match.CreatedAt,
		UpdatedAt:    match.UpdatedAt,
	}

	return &response, nil
}

// GetAllMatchesService returns all matches
func GetAllMatchesService() ([]dto.MatchListResponse, error) {
	matches, err := repositories.GetAllMatches()
	if err != nil {
		return nil, err
	}

	var responses []dto.MatchListResponse
	for _, m := range matches {
		season, _ := repositories.GetSeasonByID(m.SeasonID)
		homeTeam, _ := repositories.GetTeamByID(m.HomeTeamID)
		awayTeam, _ := repositories.GetTeamByID(m.AwayTeamID)

		status := "missing"
		if m.ScoutJSON != "" {
			status = "available"
		}

		responses = append(responses, dto.MatchListResponse{
			ID:           m.ID,
			SeasonID:     m.SeasonID,
			SeasonName:   getSeasonName(season),
			HomeTeamID:   m.HomeTeamID,
			HomeTeamName: getTeamName(homeTeam),
			AwayTeamID:   m.AwayTeamID,
			AwayTeamName: getTeamName(awayTeam),
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

// GetMatchByIDService returns a single match by ID
func GetMatchByIDService(id uuid.UUID) (*dto.MatchResponse, error) {
	match, err := repositories.GetMatchByID(id)
	if err != nil {
		return nil, errors.New(constants.ErrMatchNotFound)
	}

	season, _ := repositories.GetSeasonByID(match.SeasonID)
	homeTeam, _ := repositories.GetTeamByID(match.HomeTeamID)
	awayTeam, _ := repositories.GetTeamByID(match.AwayTeamID)

	// Fetch and parse JSON from S3
	var jsonData map[string]interface{}
	if match.ScoutJSON != "" {
		jsonData, err = httpPkg.FetchJSONFromS3(match.ScoutJSON)
		if err != nil {
			log.Printf("Failed to fetch JSON from S3: %v", err)
		}
	}

	return &dto.MatchResponse{
		ID:           match.ID,
		SeasonID:     match.SeasonID,
		SeasonName:   getSeasonName(season),
		HomeTeamID:   match.HomeTeamID,
		HomeTeamName: getTeamName(homeTeam),
		AwayTeamID:   match.AwayTeamID,
		AwayTeamName: getTeamName(awayTeam),
		Round:        match.Round,
		Location:     match.Location,
		VideoURL:     match.VideoURL,
		ScoutJSON:    match.ScoutJSON,
		JsonData:     jsonData,
		CreatedAt:    match.CreatedAt,
		UpdatedAt:    match.UpdatedAt,
	}, nil
}

// UpdateMatchService updates an existing match
func UpdateMatchService(id uuid.UUID, input *dto.UpdateMatchInput) (*dto.MatchResponse, error) {
	match, err := repositories.GetMatchByID(id)
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

	if err := repositories.UpdateMatch(match); err != nil {
		return nil, err
	}

	season, _ := repositories.GetSeasonByID(match.SeasonID)
	homeTeam, _ := repositories.GetTeamByID(match.HomeTeamID)
	awayTeam, _ := repositories.GetTeamByID(match.AwayTeamID)

	return &dto.MatchResponse{
		ID:           match.ID,
		SeasonID:     match.SeasonID,
		SeasonName:   getSeasonName(season),
		HomeTeamID:   match.HomeTeamID,
		HomeTeamName: getTeamName(homeTeam),
		AwayTeamID:   match.AwayTeamID,
		AwayTeamName: getTeamName(awayTeam),
		Round:        match.Round,
		Location:     match.Location,
		VideoURL:     match.VideoURL,
		ScoutJSON:    match.ScoutJSON,
		CreatedAt:    match.CreatedAt,
		UpdatedAt:    match.UpdatedAt,
	}, nil
}

// DeleteMatchService deletes a match
func DeleteMatchService(id uuid.UUID) error {
	return repositories.DeleteMatch(id)
}

// helper to format season name
func getSeasonName(season *models.Season) string {
	if season == nil {
		return ""
	}
	return string(season.Name) + " " + season.SeasonYear
}

// helper to get team name
func getTeamName(team *models.Team) string {
	if team == nil {
		return ""
	}
	return team.Name
}

// UploadMatchVideoService handles uploading a match video to S3
// Path: services/match_service.go
func UploadMatchVideoService(matchID uuid.UUID, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	match, err := repositories.GetMatchByID(matchID)
	if err != nil {
		return "", errors.New(constants.ErrMatchNotFound)
	}

	season, err := repositories.GetSeasonByID(match.SeasonID)
	if err != nil {
		return "", errors.New(constants.ErrSeasonNotFound)
	}

	// Create AWS session and uploader
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AWSRegion),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(sess)

	// Upload and get S3 key
	s3Key, err := storagePkg.UploadMatchVideoToS3(
		uploader,
		file,
		fileHeader,
		match.ID.String(),
		season.SeasonYear,
		string(season.Country),
		string(season.Name),
		string(season.Gender),
	)
	if err != nil {
		return "", err
	}

	// Generate public S3 URL manually
	videoURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.AWSBucketName, config.AWSRegion, s3Key)

	// Save to DB
	match.VideoURL = videoURL
	if err := repositories.UpdateMatch(match); err != nil {
		return "", err
	}

	return videoURL, nil
}

// UploadMatchScoutService handles uploading and parsing a scout (.dvw) file
func UploadMatchScoutService(matchID uuid.UUID, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	match, err := repositories.GetMatchByID(matchID)
	if err != nil {
		return "", errors.New(constants.ErrMatchNotFound)
	}

	if strings.ToLower(filepath.Ext(fileHeader.Filename)) != ".dvw" {
		return "", errors.New("invalid file type: only .dvw supported")
	}

	// Upload original .dvw file to S3
	s3InputKey := fmt.Sprintf("scouts/%s.dvw", matchID.String())
	contentType := fileHeader.Header.Get("Content-Type")

	_, err = storagePkg.UploadFileToS3(file, s3InputKey, contentType)
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

	// Upload parsed JSON to S3
	s3OutputKey := fmt.Sprintf("scout-files/%s.json", matchID.String())
	jsonURL, err := storagePkg.UploadBytesToS3(jsonBytes, s3OutputKey, "application/json")
	if err != nil {
		return "", fmt.Errorf("failed to upload scout json: %w", err)
	}

	// Save JSON URL to DB
	match.ScoutJSON = jsonURL
	if err := repositories.UpdateMatch(match); err != nil {
		return "", err
	}

	return jsonURL, nil
}
