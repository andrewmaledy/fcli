package series

import (
	"flashbacklabsio/fcli/internal/clients/overseer"
	"flashbacklabsio/fcli/internal/clients/sonarr"
	"flashbacklabsio/fcli/internal/config"
	"flashbacklabsio/fcli/internal/tui"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleSeriesCommand is the entry point for the series command.
func HandleSeriesCommand() {
	fmt.Println("Series management subcommands can be found here. Supply --help to see available series commands.")
}

// FindMediaItemByTvdbId searches for a media item by its TVDB ID.
func FindMediaItemByTvdbId(tvdbId int, mediaItems []overseer.Media) (*overseer.Media, error) {
	for _, item := range mediaItems {
		if item.TvdbId == tvdbId {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("no matching MediaItem found for TvdbId %d", tvdbId)
}

// filterSeasons returns only seasons that have files on disk.
func filterSeasons(seasons []sonarr.Season) []sonarr.Season {
	var filtered []sonarr.Season
	for _, season := range seasons {
		if season.Statistics.SizeOnDisk != 0 {
			filtered = append(filtered, season)
		}
	}
	return filtered
}

// --- Static table for `get` command (still uses bubbles table) ---

func seriesTableColumns() []table.Column {
	return []table.Column{
		{Title: "Title", Width: 30},
		{Title: "Year", Width: 6},
		{Title: "Seasons", Width: 9},
		{Title: "Episodes", Width: 10},
		{Title: "Genres", Width: 20},
		{Title: "Rating", Width: 8},
		{Title: "Size", Width: 10},
	}
}

func seriesToRow(s sonarr.Series) table.Row {
	seasonCount := len(filterSeasons(s.Seasons))
	return table.Row{
		tui.TruncateString(s.Title, 28),
		fmt.Sprintf("%d", s.Year),
		fmt.Sprintf("%d", seasonCount),
		tui.FormatEpisodes(s.Statistics.EpisodeFileCount, s.Statistics.TotalEpisodeCount),
		tui.FormatGenres(s.Genres, 18),
		tui.FormatRating(s.Ratings.Value, ""),
		tui.FormatSizeGB(s.Statistics.SizeOnDisk),
	}
}

// --- Interactive list columns and items for `searchanddelete` ---

func seriesColumns() []tui.ColumnDef {
	return []tui.ColumnDef{
		{Title: "Title", Width: 30},
		{Title: "Year", Width: 6},
		{Title: "Seasons", Width: 9},
		{Title: "Episodes", Width: 10},
		{Title: "Genres", Width: 20},
		{Title: "Rating", Width: 8},
		{Title: "Size", Width: 10},
	}
}

func seriesToListItem(s sonarr.Series) tui.ListItem {
	seasons := filterSeasons(s.Seasons)

	columns := []string{
		tui.TruncateString(s.Title, 28),
		fmt.Sprintf("%d", s.Year),
		fmt.Sprintf("%d", len(seasons)),
		tui.FormatEpisodes(s.Statistics.EpisodeFileCount, s.Statistics.TotalEpisodeCount),
		tui.FormatGenres(s.Genres, 18),
		tui.FormatRating(s.Ratings.Value, ""),
		tui.FormatSizeGB(s.Statistics.SizeOnDisk),
	}

	label := fmt.Sprintf("%s (%s)", s.Title, tui.FormatSizeGB(s.Statistics.SizeOnDisk))

	detail := []tui.KeyValue{
		{Key: "Title", Value: s.Title},
		{Key: "Year", Value: fmt.Sprintf("%d", s.Year)},
		{Key: "Network", Value: s.Network},
		{Key: "Status", Value: s.Status},
		{Key: "Runtime", Value: fmt.Sprintf("%d min", s.Runtime)},
		{Key: "Genres", Value: strings.Join(s.Genres, ", ")},
	}

	ratingStr := tui.FormatRating(s.Ratings.Value, "")
	if s.Ratings.Votes > 0 {
		ratingStr += fmt.Sprintf(" (%d votes)", s.Ratings.Votes)
	}
	detail = append(detail,
		tui.KeyValue{Key: "Rating", Value: ratingStr},
		tui.KeyValue{Key: "Seasons", Value: fmt.Sprintf("%d", len(seasons))},
		tui.KeyValue{Key: "Episodes", Value: tui.FormatEpisodes(s.Statistics.EpisodeFileCount, s.Statistics.TotalEpisodeCount)},
		tui.KeyValue{Key: "Size", Value: tui.FormatSizeGB(s.Statistics.SizeOnDisk)},
		tui.KeyValue{Key: "Path", Value: s.Path},
	)

	// Build children from seasons
	var children []tui.ListItem
	for _, season := range seasons {
		childColumns := []string{
			fmt.Sprintf("Season %d", season.SeasonNumber),
			"",
			"",
			tui.FormatEpisodes(season.Statistics.EpisodeFileCount, season.Statistics.TotalEpisodeCount),
			"",
			"",
			tui.FormatSizeGB(season.Statistics.SizeOnDisk),
		}

		childLabel := fmt.Sprintf("%s - S%d (%s)",
			s.Title, season.SeasonNumber, tui.FormatSizeGB(season.Statistics.SizeOnDisk))

		childDetail := []tui.KeyValue{
			{Key: "Series", Value: s.Title},
			{Key: "Season", Value: fmt.Sprintf("%d", season.SeasonNumber)},
			{Key: "Episodes", Value: tui.FormatEpisodes(season.Statistics.EpisodeFileCount, season.Statistics.TotalEpisodeCount)},
			{Key: "Size", Value: tui.FormatSizeGB(season.Statistics.SizeOnDisk)},
			{Key: "Monitored", Value: fmt.Sprintf("%v", season.Monitored)},
		}

		children = append(children, tui.ListItem{
			Columns:  childColumns,
			Label:    childLabel,
			Detail:   childDetail,
			Children: nil,
		})
	}

	return tui.ListItem{
		Columns:  columns,
		Label:    label,
		Detail:   detail,
		Children: children,
	}
}

func initializeSonarrClient(conf *config.Configuration) *sonarr.SonarrClient {
	sonarrClient := sonarr.NewSonarrClient(conf.SonarrURL, conf.SonarrAPIKey)

	_, err := sonarrClient.GetAllSeries()
	if err != nil {
		fmt.Println(tui.RenderConnectionError("Sonarr", conf.SonarrURL, err))
		os.Exit(1)
	}

	return sonarrClient
}

// HandleGet displays all series in a styled static table.
func HandleGet(sonarrAPIKey string, overseerAPIKey string, limit int) {
	conf := config.GetConfig("", sonarrAPIKey, overseerAPIKey)
	sonarrClient := initializeSonarrClient(conf)

	allSeries, err := sonarrClient.GetAllSeries()
	if err != nil {
		fmt.Println(tui.RenderConnectionError("Sonarr", conf.SonarrURL, err))
		os.Exit(1)
	}

	sort.Slice(allSeries, func(i, j int) bool {
		return allSeries[i].Statistics.SizeOnDisk > allSeries[j].Statistics.SizeOnDisk
	})

	var rows []table.Row
	for i := 0; i < len(allSeries) && i < limit; i++ {
		rows = append(rows, seriesToRow(allSeries[i]))
	}

	header := tui.TitleStyle.Render(fmt.Sprintf("Series (%d total, showing %d)", len(allSeries), len(rows)))
	fmt.Printf("\n%s\n\n%s\n", header, tui.RenderStaticTable(seriesTableColumns(), rows))
}

// HandleSearchAndDeleteSeries manages the interactive search and delete.
// Series with seasons auto-expand when highlighted. Users can select
// entire series or individual seasons.
func HandleSearchAndDeleteSeries(sonarrAPIKey string, overseerAPIKey string, limit int) {
	conf := config.GetConfig("", sonarrAPIKey, overseerAPIKey)
	sonarrClient := sonarr.NewSonarrClient(conf.SonarrURL, conf.SonarrAPIKey)
	overseerClient := overseer.NewOverseerClient(conf.OverseerURL, conf.OverseerAPIKey)

	var seriesData []sonarr.Series

	fetchFn := func() ([]tui.ListItem, error) {
		allSeries, err := sonarrClient.GetAllSeries()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch series from Sonarr: %w", err)
		}

		sort.Slice(allSeries, func(i, j int) bool {
			return allSeries[i].Statistics.SizeOnDisk > allSeries[j].Statistics.SizeOnDisk
		})

		end := limit
		if end > len(allSeries) {
			end = len(allSeries)
		}
		seriesData = allSeries[:end]

		var items []tui.ListItem
		for _, s := range seriesData {
			items = append(items, seriesToListItem(s))
		}
		return items, nil
	}

	deleteFn := func(selections []tui.Selection) []string {
		var results []string

		// Group selections by series
		type seriesAction struct {
			deleteWhole bool
			seasonIdxs  []int
		}
		actions := make(map[int]*seriesAction)

		for _, sel := range selections {
			if sel.ItemIndex >= len(seriesData) {
				continue
			}
			a, ok := actions[sel.ItemIndex]
			if !ok {
				a = &seriesAction{}
				actions[sel.ItemIndex] = a
			}
			if sel.ChildIndex == -1 {
				a.deleteWhole = true
			} else {
				a.seasonIdxs = append(a.seasonIdxs, sel.ChildIndex)
			}
		}

		for itemIdx, action := range actions {
			s := seriesData[itemIdx]
			seasons := filterSeasons(s.Seasons)

			if action.deleteWhole || len(action.seasonIdxs) == len(seasons) {
				results = append(results, deleteEntireSeries(s, sonarrClient, overseerClient)...)
				continue
			}

			// Delete individual seasons
			for _, seasonIdx := range action.seasonIdxs {
				if seasonIdx >= len(seasons) {
					continue
				}
				season := seasons[seasonIdx]
				episodeFiles, err := sonarrClient.GetEpisodeFilesForSeries(s.ID, &season.SeasonNumber)
				if err != nil {
					results = append(results, tui.ErrorStyle.Render(
						fmt.Sprintf("Failed to get episode files for %s S%d: %s",
							s.Title, season.SeasonNumber, err)))
					continue
				}
				if err := sonarrClient.DeleteEpisodeFiles(episodeFiles); err != nil {
					results = append(results, tui.ErrorStyle.Render(
						fmt.Sprintf("Failed to delete %s S%d episodes: %s",
							s.Title, season.SeasonNumber, err)))
				} else {
					results = append(results, tui.SuccessStyle.Render(
						fmt.Sprintf("Deleted Season %d of '%s' from Sonarr",
							season.SeasonNumber, s.Title)))
				}
			}

			// Unmonitor deleted seasons
			for i := range s.Seasons {
				for _, seasonIdx := range action.seasonIdxs {
					if seasonIdx < len(seasons) && s.Seasons[i].SeasonNumber == seasons[seasonIdx].SeasonNumber {
						s.Seasons[i].Monitored = false
					}
				}
			}
			if err := sonarrClient.UpdateSeries(s); err != nil {
				results = append(results, tui.WarningStyle.Render(
					fmt.Sprintf("Failed to unmonitor deleted seasons: %s", err)))
			} else {
				results = append(results, tui.SuccessStyle.Render("Unmonitored deleted seasons in Sonarr"))
			}
		}

		return results
	}

	model := tui.NewSelectTableModel("Select Series to Delete", seriesColumns(), fetchFn, deleteFn)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println(tui.ErrorStyle.Render("Fatal error: " + err.Error()))
		os.Exit(1)
	}
}

// deleteEntireSeries deletes a series from Sonarr and its request from Overseer.
func deleteEntireSeries(s sonarr.Series, sonarrClient *sonarr.SonarrClient, overseerClient *overseer.OverseerClient) []string {
	var results []string

	err := sonarrClient.DeleteSeries(s.ID)
	if err != nil {
		results = append(results, tui.ErrorStyle.Render(fmt.Sprintf("Sonarr delete failed for '%s': %s", s.Title, err)))
	} else {
		results = append(results, tui.SuccessStyle.Render(fmt.Sprintf("Deleted '%s' from Sonarr", s.Title)))
	}

	mediaItems, err := overseerClient.GetMedia()
	if err != nil {
		results = append(results, tui.WarningStyle.Render("Could not fetch Overseer media: "+err.Error()))
		return results
	}

	media, err := FindMediaItemByTvdbId(s.TvdbID, mediaItems)
	if err != nil {
		return results
	}

	err = overseerClient.DeleteMedia(media.Id)
	if err != nil {
		results = append(results, tui.ErrorStyle.Render(fmt.Sprintf("Overseer delete failed for '%s': %s", s.Title, err)))
	} else {
		results = append(results, tui.SuccessStyle.Render(fmt.Sprintf("Deleted '%s' from Overseer", s.Title)))
	}

	return results
}
