package movies

import (
	"bufio"
	"flashbacklabsio/fcli/internal/clients/overseer"
	"flashbacklabsio/fcli/internal/clients/radarr"
	"flashbacklabsio/fcli/internal/config"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

// HandleMoviesCommand is the entry point for the movies command.
func HandleMoviesCommand() {
	fmt.Println("Movie management subcommands can be found here. Supply --help to see available movie commands.")
}

// FindMediaItemByTmdbID searches for a media item by its TMDB ID.
func FindMediaItemByTmdbID(tmdbID int, media []overseer.Media) (*overseer.Media, error) {
	for _, item := range media {
		if item.TmdbId == tmdbID {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("no matching MovieItem found for TMDBID %d", tmdbID)
}

// DisplayMovies displays a list of movies sorted by size.
func DisplayMovies(movies []radarr.Movie, limit int) {
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].SizeOnDisk > movies[j].SizeOnDisk
	})

	fmt.Println("Movies:")
	for i, movie := range movies {
		if i >= limit {
			break
		}
		fmt.Printf("%d: %s (%.2f GB)\n", i+1, movie.Title, float64(movie.SizeOnDisk)/(1024*1024*1024))
	}
}

// GetUserSelections prompts the user to select movies to delete.
func GetUserSelections(limit int) ([]int, error) {
	fmt.Print(Green + "Select movie numbers to delete (comma-separated): " + Reset)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	// Parse selections
	var selections []int
	for _, selection := range strings.Split(strings.TrimSpace(input), ",") {
		movieIndex, err := strconv.Atoi(strings.TrimSpace(selection))
		if err != nil || movieIndex < 1 || movieIndex > limit {
			fmt.Printf("Invalid selection: %s\n", selection)
			continue
		}
		selections = append(selections, movieIndex)
	}
	return selections, nil
}

// ConfirmDeletion prompts the user to confirm the deletion of a movie.
func ConfirmDeletion(movieTitle string, movieSize int64) bool {
	fmt.Printf(Yellow+"Are you sure you want to delete '%s' (%.2f GB)? (y/N): "+Reset, movieTitle, float64(movieSize)/(1024*1024*1024))
	reader := bufio.NewReader(os.Stdin)
	confirmInput, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(confirmInput)) == "y"
}

// HandleSearchAndDelete manages the search and delete process.
func HandleSearchAndDelete(radarrAPIKey, overseerAPIKey string, limit int) {
	// Initialize and get configuration
	config.InitConfig()
	conf := config.GetConfig()
	radarrClient := radarr.NewRadarrClient(conf.RadarrURL, conf.RadarrAPIKey)
	overseerClient := overseer.NewOverseerClient(conf.OverseerURL, conf.OverseerAPIKey)
	fmt.Printf("Radarr API Endpoint %v\n", conf.RadarrURL)

	// Fetch and display movies from Radarr
	radarrMovies, _ := radarrClient.FetchMovies()
	DisplayMovies(radarrMovies, limit)

	// Get user selections
	selections, err := GetUserSelections(len(radarrMovies))
	if err != nil {
		fmt.Println(Red + err.Error() + Reset)
		return
	}

	// Fetch media items from Overseer
	overseerMedia, err := overseerClient.GetMedia()
	if err != nil {
		fmt.Println(Red + err.Error() + Reset)
		return
	}

	// Process selections
	for _, movieIndex := range selections {
		selectedMovie := radarrMovies[movieIndex-1]

		if ConfirmDeletion(selectedMovie.Title, selectedMovie.SizeOnDisk) {
			// Delete media item from Overseer
			if movieItem, err := FindMediaItemByTmdbID(selectedMovie.TmdbId, overseerMedia); err != nil {
				fmt.Println(Red + err.Error() + Reset)
			} else {
				if err := overseerClient.DeleteMedia(movieItem.Id); err != nil {
					fmt.Println(Red + err.Error() + Reset)
				} else {
					fmt.Printf(Green+"Request '%v' was successfully deleted from Overseer.\n"+Reset, selectedMovie.Title)
				}
			}

			// Delete movie from Radarr
			if err := radarrClient.DeleteMovie(selectedMovie.Id); err != nil {
				fmt.Println(Red + err.Error() + Reset)
			} else {
				fmt.Printf(Green+"Movie '%v' was successfully deleted from Radarr.\n"+Reset, selectedMovie.Title)
			}
		} else {
			fmt.Printf("Skipped deletion of '%s'.\n", selectedMovie.Title)
		}
	}
}
