package scout

import (
	"encoding/json"
	"fmt"
	"go-gin-starter/dto"
	httpPkg "go-gin-starter/pkg/http"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ScoutParseResponse struct {
	Success  bool                   `json:"success"`
	MatchID  string                 `json:"match_id"`
	JsonData map[string]interface{} `json:"json_data"`
	Message  string                 `json:"message"`
}

// ExtractScoutMetadata extracts basic metadata from a scout JSON file
func ExtractScoutMetadata(jsonURL string) (*dto.ScoutMetadataResponse, error) {
	// Fetch JSON from S3
	jsonData, err := httpPkg.FetchJSONFromS3(jsonURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JSON: %w", err)
	}

	// Extract metadata from the JSON structure
	matchInfo, ok := jsonData["match_info"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid match_info structure")
	}

	// Extract home and away team info
	home, ok := matchInfo["home"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid home team structure")
	}
	away, ok := matchInfo["away"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid away team structure")
	}

	// Parse match date from metadata
	metadata, ok := jsonData["metadata"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid metadata structure")
	}

	// Try to parse the date from metadata
	var matchDate time.Time
	if generator, ok := metadata["generator"].(map[string]interface{}); ok {
		if day, ok := generator["day"].(string); ok {
			matchDate, _ = time.Parse("01/02/2006 15.04.05", day)
		}
	}

	// Build response
	response := &dto.ScoutMetadataResponse{
		Competition: "Bundesliga Men", // This should be extracted from the JSON if available
		Season:      "2024-2025",      // This should be extracted from the JSON if available
		HomeTeam:    home["name"].(string),
		AwayTeam:    away["name"].(string),
		HomeScore:   int(home["sets_won"].(float64)),
		AwayScore:   int(away["sets_won"].(float64)),
		MatchDate:   matchDate,
		Location:    matchInfo["location"].(string),
	}

	return response, nil
}

// CallPythonParser sends the .dvw file path to the Python microservice via query param
func CallPythonParser(s3Key string) (*ScoutParseResponse, error) {
	url := os.Getenv("PYTHON_SCOUT_PARSER_URL") + "?input_file=" + url.QueryEscape(s3Key)

	req, err := http.NewRequest("POST", url, nil) // No body needed
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("parser error: status %d", resp.StatusCode)
	}

	var result ScoutParseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("parser returned success=false")
	}

	return &result, nil
}
