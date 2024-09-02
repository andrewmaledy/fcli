package series

import (
	"bufio"
	"flashbacklabsio/fcli/internal/clients/overseer"
	"flashbacklabsio/fcli/internal/clients/sonarr"
	"flashbacklabsio/fcli/internal/config"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
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

// Function to filter seasons with SizeOnDisk not equal to 0
func filterSeasons(seasons []sonarr.Season) []sonarr.Season {
	var filteredSeasons []sonarr.Season
	for _, season := range seasons {
		if season.Statistics.SizeOnDisk != 0 {
			filteredSeasons = append(filteredSeasons, season)
		}
	}
	return filteredSeasons
}

// HandleSeriesCommand is the entry point for the series command
func HandleSeriesCommand() {
	fmt.Println("Series management sub commands can be found here. Supply --help to see available series commands.")
	// Add logic here
}
func HandleSearchAndDeleteSeries(sonarrAPIKey string, overseerAPIKey string, limit int) {

	// Initialize and get configuration
	config.InitConfig()
	conf := config.GetConfig()
	if len(sonarrAPIKey) > 0 {
		conf.SonarrAPIKey = sonarrAPIKey
	}
	if len(overseerAPIKey) > 0 {
		conf.OverseerAPIKey = overseerAPIKey
	}
	overseerClient := overseer.NewOverseerClient(conf.OverseerURL, conf.OverseerAPIKey)
	sonarrClient := sonarr.NewSonarrClient(conf.SonarrURL, conf.SonarrAPIKey)
	fmt.Printf("Sonarr API Endpoint: %v\n", conf.SonarrURL)

	// Fetch and display series from Sonarr
	sonarrSeries, err := sonarrClient.GetAllSeries()
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
	selectedSeries.Seasons = filterSeasons(selectedSeries.Seasons) //filter out series that aren't actually present.

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
		if strings.ToLower(confirmInput) != "y" {
			fmt.Printf("Skipped deletion of series '%s'.\n", selectedSeries.Title)
		} else {
			err = sonarrClient.DeleteSeries(selectedSeries.ID)
			if err != nil {
				fmt.Printf("Error deleting series: %v\n", err)
			} else {
				fmt.Printf(Green+"Series '%s' successfully deleted from Sonarr.\n"+Reset, selectedSeries.Title)
			}

			// Delete corresponding request from Overseer
			mediaItems, err := overseerClient.GetMedia()
			if err != nil {
				fmt.Printf("Error fetching media: %v\n", err)
			} else {
				media, err := FindMediaItemByTvdbId(selectedSeries.TvdbId, mediaItems)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					err = overseerClient.DeleteMedia(media.Id)
					if err != nil {
						fmt.Printf("Error deleting request from Overseer: %v\n", err)
					} else {
						fmt.Printf(Green+"Request '%s' successfully deleted from Overseer.\n"+Reset, selectedSeries.Title)
					}
				}
			}
		}
	} else {
		// Delete selected episodefiles
		selectedSeason := &selectedSeries.Seasons[seasonIndex-1]
		selectedSeason.Monitored = false
		fmt.Printf(Yellow+"Are you sure you want to delete Season %d of '%s'? (y/N): "+Reset, selectedSeason.SeasonNumber, selectedSeries.Title)
		confirmInput, _ := reader.ReadString('\n')
		confirmInput = strings.TrimSpace(confirmInput)
		if strings.ToLower(confirmInput) != "y" {
			fmt.Printf("Skipped deletion of Season %d of series '%s'.\n", selectedSeason.SeasonNumber, selectedSeries.Title)
		} else {
			episodeFiles, err := sonarrClient.GetEpiosdeFilesForSeries(selectedSeries.ID, &selectedSeason.SeasonNumber)
			if err != nil {
				fmt.Printf("Error getting season episode files: %v\n", err)
			}
			err = sonarrClient.DeleteEpisodeFiles(episodeFiles)
			if err != nil {
				fmt.Printf("Error deleting episodes: %v\n", err)
			} else {
				fmt.Printf(Green+"Season %d of series '%s' successfully deleted from Sonarr.\n"+Reset, selectedSeason.SeasonNumber, selectedSeries.Title)
			}

		}
		// Update the series to unmonitor the deleted season.
		err := sonarrClient.UpdateSeries(selectedSeries)
		if err != nil {
			fmt.Printf("Error removing season %d monitoring. This means the series will be downloaded automatically again. ERROR: %v\n", selectedSeason.SeasonNumber, err)
		} else {
			fmt.Printf(Green+"Season %d of series '%s' successfully unmonitored in Sonarr.\n"+Reset, selectedSeason.SeasonNumber, selectedSeries.Title)
		}

	}
}
