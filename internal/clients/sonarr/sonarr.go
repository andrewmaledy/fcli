package sonarr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAllSeries fetches all series from the Sonarr API.
func GetAllSeries(baseURL, apiKey string) ([]Series, error) {
	params := "/series?apikey=" + apiKey
	resp, err := http.Get(baseURL + params)
	if err != nil {
		return nil, fmt.Errorf("error fetching series: %v", err)
	}
	defer resp.Body.Close()

	var series []Series
	if err := json.NewDecoder(resp.Body).Decode(&series); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return series, nil
}

// GetEpisodeFiles fetches all episode files for a specific series from the Sonarr API.
func GetEpiosdeFilesForSeries(baseURL, apiKey string, seriesID int, seasonNumber *int) ([]EpisodeFile, error) {
	params := fmt.Sprintf("/episodefile?seriesId=%d", seriesID)
	req, err := http.NewRequest("GET", baseURL+params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch episode files: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch episode files. Status code: %d", resp.StatusCode)
	}

	var episodeFiles []EpisodeFile
	if err := json.NewDecoder(resp.Body).Decode(&episodeFiles); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// If seasonNumber is provided, filter the results.
	if seasonNumber != nil {
		var filteredFiles []EpisodeFile
		for _, file := range episodeFiles {
			if file.SeasonNumber == *seasonNumber {
				filteredFiles = append(filteredFiles, file)
			}
		}
		return filteredFiles, nil
	}

	return episodeFiles, nil
}

// UpdateSeries sends a PUT request to the Sonarr API to update a series by its ID.
func UpdateSeries(baseURL, apiKey string, series Series) error {
	client := &http.Client{}

	// Convert the Series struct to JSON.
	requestBody, err := json.Marshal(series)
	if err != nil {
		return fmt.Errorf("failed to marshal series: %v", err)
	}
	// Construct the API endpoint for the PUT request.
	endpoint := fmt.Sprintf("%s/series/%d", baseURL, series.ID)

	// Create the HTTP PUT request.
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	// Set required headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Api-Key", apiKey)

	// Execute the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update series: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-OK HTTP responses.
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	return nil
}

// DeleteSeriesAllSeasons deletes all seasons for a specific series from the Sonarr API.
func DeleteSeries(baseURL, apiKey string, seriesID int) error {
	client := &http.Client{}

	// Construct the API endpoint for the DELETE request.
	endpoint := fmt.Sprintf("%s/series/%d?deleteFiles=true", baseURL, seriesID)

	// Create the HTTP DELETE request.
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set required headers.
	req.Header.Set("accept", "*/*")
	req.Header.Set("X-Api-Key", apiKey)

	// Execute the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete series: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-OK HTTP responses.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete series with ID %d. Status code: %d", seriesID, resp.StatusCode)
	}

	return nil
}

// DeleteEpisodeFiles deletes the specified episode files by their IDs from the Sonarr API.
func DeleteEpisodeFiles(baseURL, apiKey string, episodeFiles []EpisodeFile) error {
	// Extract episode file IDs from the provided episode files.
	var episodeFileIds []int
	for _, file := range episodeFiles {
		episodeFileIds = append(episodeFileIds, file.ID)
	}

	// Create the request body with the episode file IDs.
	requestBody := map[string]interface{}{
		"episodeFileIds": episodeFileIds,
	}

	// Marshal the request body to JSON.
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Construct the API endpoint.
	endpoint := fmt.Sprintf("%s/episodefile/bulk", baseURL)

	// Create the HTTP DELETE request.
	req, err := http.NewRequest("DELETE", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set required headers.
	req.Header.Set("accept", "*/*")
	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete episode files: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-OK HTTP responses.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete episode files. Status code: %d", resp.StatusCode)
	}

	return nil
}
