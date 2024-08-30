package series

import (
	"bufio"
	"flashbacklabsio/fcli/internal/clients/overseer"
	"flashbacklabsio/fcli/internal/clients/sonarr"
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

func FindMediaItemByTvdbId(tvdbId int, mediaItems []overseer.Media) (*overseer.Media, error) {
	for _, item := range mediaItems {
		if item.TvdbId == tvdbId {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("no matching MediaItem found for TvdbId %d", tvdbId)
}

// HandleSeriesCommand is the entry point for the series command
func HandleSeriesCommand() {
	fmt.Println("Series management sub commands can be found here. Supply --help to see available series commands.")
	// Add logic here
}
func HandleSearchAndDeleteSeries(sonarrAPIKey string, overseerAPIKey string, limit int) {

	if sonarrAPIKey == "" {
		sonarrAPIKey = viper.GetString("sonarr.apiKey")
	}
	if overseerAPIKey == "" {
		overseerAPIKey = viper.GetString("overseer.apiKey")
	}
	sonarrURL := viper.GetString("sonarr.url")
	overseerURL := viper.GetString("overseer.url")
	fmt.Printf("Sonarr API Endpoint: %v\n", sonarrURL)

	// Fetch and display series from Sonarr
	sonarrSeries, err := sonarr.GetAllSeries(sonarrURL, sonarrAPIKey)
	if err != nil {
		fmt.Printf("Error fetching series: %v\n", err)
		return
	}

	sort.Slice(sonarrSeries, func(i, j int) bool {
		return sonarrSeries[i].Statistics.SizeOnDisk > sonarrSeries[j].Statistics.SizeOnDisk
	})

	for i, series := range sonarrSeries {
		if i > limit-1 {
			break
		}

		fmt.Printf("%d: %s (%.2f GB)\n", i+1, series.Title, float64(series.Statistics.SizeOnDisk)/(1024*1024*1024))
	}

	// Ask user to select a series
	fmt.Print(Green + "Select series number to view seasons or delete (0 to delete entire series): " + Reset)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	seriesIndex, err := strconv.Atoi(strings.TrimSpace(input))

	if err != nil || seriesIndex < 0 || seriesIndex > len(sonarrSeries) {
		fmt.Printf("Invalid selection: %s\n", input)
		return
	}

	if seriesIndex == 0 {
		fmt.Println("No series selected. Exiting.")
		return
	}

	selectedSeries := sonarrSeries[seriesIndex-1]
	fmt.Printf("Selected series: %s\n", selectedSeries.Title)

	fmt.Println("Seasons:")
	for i, season := range selectedSeries.Seasons {
		fmt.Printf("%d: Season %d (%.2f GB)\n", i+1, season.SeasonNumber, float64(season.Statistics.SizeOnDisk)/(1024*1024*1024))
	}

	// Ask user to select a season or delete the entire series
	fmt.Print(Green + "Select season number to delete or enter 0 to delete the entire series: " + Reset)
	seasonInput, _ := reader.ReadString('\n')
	seasonIndex, err := strconv.Atoi(strings.TrimSpace(seasonInput))
	if err != nil || seasonIndex < 0 || seasonIndex > len(selectedSeries.Seasons) {
		fmt.Printf("Invalid selection: %s\n", seasonInput)
		return
	}

	if seasonIndex == 0 {
		// Delete entire series
		fmt.Printf(Yellow+"Are you sure you want to delete the entire series '%s'? (y/N): "+Reset, selectedSeries.Title)
		confirmInput, _ := reader.ReadString('\n')
		confirmInput = strings.TrimSpace(confirmInput)
		if strings.ToLower(confirmInput) == "y" {
			err = sonarr.DeleteSeriesAllSeasons(sonarrURL, sonarrAPIKey, selectedSeries.ID)
			if err != nil {
				fmt.Printf("Error deleting series: %v\n", err)
			} else {
				fmt.Printf(Green+"Series '%s' successfully deleted from Sonarr.\n"+Reset, selectedSeries.Title)
			}

			// Delete corresponding request from Overseer
			mediaItems, err := overseer.GetMedia(overseerURL, overseerAPIKey)
			if err != nil {
				fmt.Printf("Error fetching media: %v\n", err)
			} else {
				media, err := FindMediaItemByTvdbId(selectedSeries.TvdbId, mediaItems)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					err = overseer.DeleteMedia(overseerURL, overseerAPIKey, media.Id)
					if err != nil {
						fmt.Printf("Error deleting request from Overseer: %v\n", err)
					} else {
						fmt.Printf(Green+"Request '%s' successfully deleted from Overseer.\n"+Reset, selectedSeries.Title)
					}
				}
			}
		} else {
			fmt.Printf("Skipped deletion of series '%s'.\n", selectedSeries.Title)
		}
	} else {
		// Delete selected season
		selectedSeason := selectedSeries.Seasons[seasonIndex-1]
		fmt.Printf(Yellow+"Are you sure you want to delete Season %d of '%s'? (y/N): "+Reset, selectedSeason.SeasonNumber, selectedSeries.Title)
		confirmInput, _ := reader.ReadString('\n')
		confirmInput = strings.TrimSpace(confirmInput)
		if strings.ToLower(confirmInput) == "y" {
			err = sonarr.DeleteSeasonForSeries(sonarrURL, sonarrAPIKey, selectedSeries.ID, selectedSeason.SeasonNumber)
			if err != nil {
				fmt.Printf("Error deleting season: %v\n", err)
			} else {
				fmt.Printf(Green+"Season %d of series '%s' successfully deleted from Sonarr.\n"+Reset, selectedSeason.SeasonNumber, selectedSeries.Title)
			}

			// Update Overseer request if applicable
			mediaItems, err := overseer.GetMedia(overseerURL, overseerAPIKey)
			if err != nil {
				fmt.Printf("Error fetching media: %v\n", err)
			} else {
				mediaItem, err := FindMediaItemByTvdbId(selectedSeries.TvdbId, mediaItems)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					// Assuming Overseer allows partial deletion or disabling of seasons
					// If not supported, delete the whole request
					// err = overseer.DisableSeason(overseerURL, overseerAPIKey, mediaItem.Id, selectedSeason.SeasonNumber)
					err = overseer.DeleteMedia(overseerURL, overseerAPIKey, mediaItem.Id)
					if err != nil {
						fmt.Printf("Error updating request in Overseer: %v\n", err)
					} else {
						fmt.Printf(Green+"Request for Season %d of series '%s' successfully deleted from Overseer.\n"+Reset, selectedSeason.SeasonNumber, selectedSeries.Title)
					}
				}
			}
		} else {
			fmt.Printf("Skipped deletion of Season %d of series '%s'.\n", selectedSeason.SeasonNumber, selectedSeries.Title)
		}
	}
}
