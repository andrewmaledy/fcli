package radarr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RadarrClient holds the base URL and API key for the Radarr API.
type RadarrClient struct {
	BaseURL string
	APIKey  string
}

// NewRadarrClient creates a new instance of RadarrClient with the given base URL and API key.
func NewRadarrClient(baseURL, apiKey string) *RadarrClient {
	return &RadarrClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

// FetchMovies retrieves the list of movies from Radarr.
func (client *RadarrClient) FetchMovies() ([]Movie, error) {
	params := "/movie?excludeLocalCovers=false&apikey=" + client.APIKey
	resp, err := http.Get(client.BaseURL + params)
	if err != nil {
		return nil, fmt.Errorf("error fetching movies: %w", err)
	}
	defer resp.Body.Close()

	var movies []Movie
	if err := json.NewDecoder(resp.Body).Decode(&movies); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return filterMovies(movies), nil
}

func filterMovies(movies []Movie) []Movie {
	var filtered []Movie
	for _, movie := range movies {
		if movie.SizeOnDisk > 0 {
			filtered = append(filtered, movie)
		}
	}
	return filtered
}

// DeleteMovie sends a DELETE request to the Radarr API to remove a movie by its ID.
func (client *RadarrClient) DeleteMovie(movieID int) error {
	// Construct the API endpoint for the DELETE request.
	endpoint := fmt.Sprintf("%s/movie/%d?deleteFiles=true&apikey=%s", client.BaseURL, movieID, client.APIKey)

	// Create the HTTP DELETE request.
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers.
	req.Header.Set("accept", "*/*")

	// Execute the HTTP request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-OK HTTP responses.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete movie with ID %d. Status code: %d", movieID, resp.StatusCode)
	}

	return nil
}
