package sonarr

import (
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

// GetSeasonsForSeries fetches all seasons for a specific series from the Sonarr API.
func GetSeasonsForSeries(baseURL, apiKey string, seriesID int) ([]Season, error) {
	params := fmt.Sprintf("/series/%d?apikey=%s", seriesID, apiKey)
	resp, err := http.Get(baseURL + params)
	if err != nil {
		return nil, fmt.Errorf("error fetching seasons for series %d: %v", seriesID, err)
	}
	defer resp.Body.Close()

	var series Series
	if err := json.NewDecoder(resp.Body).Decode(&series); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return series.Seasons, nil
}

// DeleteSeasonForSeries deletes a specific season for a series from the Sonarr API.
func DeleteSeasonForSeries(baseURL, apiKey string, seriesID int, seasonNumber int) error {
	return nil
	client := &http.Client{}

	// Construct the API endpoint for the DELETE request.
	endpoint := fmt.Sprintf("%s/episodefile?seriesId=%d&seasonNumber=%d", baseURL, seriesID, seasonNumber)

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
		return fmt.Errorf("failed to delete season for series: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-OK HTTP responses.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete season for series with ID %d. Status code: %d", seriesID, resp.StatusCode)
	}

	return nil
}

// DeleteSeriesAllSeasons deletes all seasons for a specific series from the Sonarr API.
func DeleteSeriesAllSeasons(baseURL, apiKey string, seriesID int) error {
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
