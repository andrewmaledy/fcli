package radarr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RadarrClient holds the base URL and API key for the Radarr API.
type RadarrClient struct {
	BaseURL    string
	APIKey     string
	httpClient *http.Client
}

// NewRadarrClient creates a new instance of RadarrClient with the given base URL and API key.
func NewRadarrClient(baseURL, apiKey string) *RadarrClient {
	return &RadarrClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *RadarrClient) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Api-Key", c.APIKey)
}

// GetMovies retrieves the list of movies from Radarr.
func (c *RadarrClient) GetMovies() ([]Movie, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/movie?excludeLocalCovers=false", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching movies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching movies", resp.StatusCode)
	}

	var movies []Movie
	if err := json.NewDecoder(resp.Body).Decode(&movies); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return movies, nil
}

// DeleteMovie sends a DELETE request to the Radarr API to remove a movie by its ID.
func (c *RadarrClient) DeleteMovie(movieID int) error {
	endpoint := fmt.Sprintf("%s/movie/%d?deleteFiles=true", c.BaseURL, movieID)

	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete movie with ID %d. Status code: %d", movieID, resp.StatusCode)
	}
	return nil
}
