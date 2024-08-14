package radarr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchMovies(baseUrl, apiKey string) []Movie {
	params := "/movie?excludeLocalCovers=false&apikey=" + apiKey
	resp, err := http.Get(baseUrl + params)
	if err != nil {
		fmt.Println("Error fetching movies:", err)
		return nil
	}
	defer resp.Body.Close()

	var movies []Movie
	if err := json.NewDecoder(resp.Body).Decode(&movies); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil
	}

	return filterMovies(movies)
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

// DeleteMediaFile sends a DELETE request to the Radarr API to remove a MediaFile by its ID.
func DeleteMovie(baseURL, apiKey string, movieID int) error {
	client := &http.Client{}

	// Construct the API endpoint for the DELETE request.
	endpoint := fmt.Sprintf("%s/movie/%d", baseURL, movieID)

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
		return fmt.Errorf("failed to delete media file: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-OK HTTP responses.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to Movie with ID %d. Status code: %d", movieID, resp.StatusCode)
	}

	return nil
}
