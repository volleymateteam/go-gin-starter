package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type ScoutParseResponse struct {
	Success  bool                   `json:"success"`
	MatchID  string                 `json:"match_id"`
	JsonData map[string]interface{} `json:"json_data"`
	Message  string                 `json:"message"`
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
