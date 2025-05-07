package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchJSONFromS3(jsonURL string) (map[string]interface{}, error) {
	resp, err := http.Get(jsonURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JSON: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("S3 returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	return jsonData, nil
}
