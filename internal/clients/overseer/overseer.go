package overseer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func GetMedia(baseUrl string, apiKey string) ([]MediaItem, error) {
	var allMedia []MediaItem
	take := 100 // Number of items to fetch per request; adjust based on API's limits.
	skip := 0   // Number of items to skip; used for pagination.

	client := &http.Client{}

	for {
		// Construct the API endpoint with query parameters.
		endpoint, err := url.Parse(fmt.Sprintf("%s/media", baseUrl))
		if err != nil {
			return nil, fmt.Errorf("failed to parse base URL: %v", err)
		}

		// Set query parameters for pagination and sorting.
		query := endpoint.Query()
		query.Set("take", strconv.Itoa(take))
		query.Set("skip", strconv.Itoa(skip))
		query.Set("sort", "added")
		endpoint.RawQuery = query.Encode()

		// Create the HTTP GET request.
		req, err := http.NewRequest("GET", endpoint.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Set required headers.
		req.Header.Set("accept", "application/json")
		req.Header.Set("X-Api-Key", apiKey)

		// Execute the HTTP request.
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Handle non-OK HTTP responses.
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			return nil, fmt.Errorf("received non-OK HTTP status: %s, body: %s", resp.Status, bodyString)
		}

		// Read and parse the response body.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		var apiResp apiResponse
		err = json.Unmarshal(body, &apiResp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %v", err)
		}

		// Map each media item from the API response to the MediaItem struct.
		for _, item := range apiResp.Results {
			mediaItem := MediaItem{
				Id:     item.Id,
				TmdbId: item.TmdbId,
				Title:  item.Title,
				Size:   item.Size,
			}
			allMedia = append(allMedia, mediaItem)
		}

		// Check if all media items have been fetched; exit loop if done.
		totalResults := apiResp.PageInfo.Results
		skip += take
		if skip >= totalResults {
			break
		}
	}

	return allMedia, nil
}

// DeleteMedia sends a DELETE request to the API to remove a media item by its ID.
func DeleteMedia(baseUrl string, apiKey string, mediaId int) error {
	client := &http.Client{}

	// Construct the API endpoint for the DELETE request.
	endpoint := fmt.Sprintf("%s/media/%d", baseUrl, mediaId)

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
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-OK HTTP responses.
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	return nil
}
