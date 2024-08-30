package movies

import (
	"bufio"
	"flashbacklabsio/fcli/internal/clients/overseer"
	"flashbacklabsio/fcli/internal/clients/radarr"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

// HandleMoviesCommand is the entry point for the movies command
func HandleMoviesCommand() {
	fmt.Println("Movie management sub commands can be found here. Supply --help to see available movie commands.")
	// Add logic here
}

func FindMediaItemByTmdbID(tmdbID int, media []overseer.Media) (*overseer.Media, error) {
	for _, item := range media {
		if item.TmdbId == tmdbID {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("no matching MovieItem found for TMDBID %d", tmdbID)
}

func HandleSearchAndDelete(radarrAPIKey string, overseerAPIKey string, limit int) {

	if radarrAPIKey == "" {
		radarrAPIKey = viper.GetString("radarr.apiKey")
	}
	if overseerAPIKey == "" {
		overseerAPIKey = viper.GetString("overseer.apiKey")

	}
	radarrURL := viper.GetString("radarr.url")
	overseerURL := viper.GetString("overseer.url")
	fmt.Printf("Radarr API Endpoint %v\n", radarrURL)
	// Fetch and display movies from Radarr
	radarrMovies := radarr.FetchMovies(radarrURL, radarrAPIKey)
	sort.Slice(radarrMovies, func(i, j int) bool {
		return radarrMovies[i].SizeOnDisk > radarrMovies[j].SizeOnDisk
	})
	// Print movies with a number for selection
	fmt.Println("Movies:")

	for i, movie := range radarrMovies {
		if i > limit-1 {
			break
		}
		fmt.Printf("%d: %s (%.2f GB)\n", i+1, movie.Title, (float64(movie.SizeOnDisk) / (1024 * 1024 * 1024)))
	}

	// Ask user to select movies to delete
	fmt.Print(Green + "Select movie numbers to delete (comma-separated): " + Reset)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	selections := strings.Split(input[:len(input)-1], ",")

	result, err := overseer.GetMedia(overseerURL, overseerAPIKey)

	if err != nil {
		fmt.Print(err.Error())
	}

	for _, selection := range selections {
		selection = strings.TrimSpace(selection)
		movieIndex, err := strconv.Atoi(selection)
		if err != nil || movieIndex < 1 || movieIndex > len(radarrMovies) {
			fmt.Printf("Invalid selection: %s\n", selection)
			continue
		}

		selectedMovie := radarrMovies[movieIndex-1]
		fmt.Printf(Yellow+"Are you sure you want to delete '%s' (%d GB)? (y/N): "+Reset, selectedMovie.Title, selectedMovie.SizeOnDisk/(1024*1024*1024))
		confirmInput, _ := reader.ReadString('\n')
		confirmInput = strings.TrimSpace(confirmInput)
		// Proceed with Overseer delete
		if strings.ToLower(confirmInput) == "y" {

			movieItem, err := FindMediaItemByTmdbID(selectedMovie.TmdbId, result)
			if err != nil {
				fmt.Println(err.Error())
			} else {

				fmt.Printf("Found matching media item for %v. Overseer ID: %d \n", selectedMovie.Title, movieItem.Id)
				err = overseer.DeleteMedia(overseerURL, overseerAPIKey, movieItem.Id)
				if err != nil {
					fmt.Print(err.Error())
				} else {
					fmt.Printf(Green+"Request '%v' was successfully deleted from Jellyseer.\n"+Reset, selectedMovie.Title)
				}
			}
			err = radarr.DeleteMovie(radarrURL, radarrAPIKey, selectedMovie.Id)
			if err != nil {
				fmt.Print(err.Error())
			} else {
				fmt.Printf(Green+"Movie '%v' was successfully deleted from Radarr.\n"+Reset, selectedMovie.Title)
			}

		} else { // cancelling deletion
			fmt.Printf("Skipped deletion of '%s'.\n", selectedMovie.Title)
		}
	}

}
