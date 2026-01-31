package movies

import (
	"flashbacklabsio/fcli/internal/clients/overseer"
	"flashbacklabsio/fcli/internal/clients/radarr"
	"flashbacklabsio/fcli/internal/config"
	"flashbacklabsio/fcli/internal/tui"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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

// initializeRadarrClient creates a Radarr client and tests the connection.
func initializeRadarrClient(conf *config.Configuration) *radarr.RadarrClient {
	radarrClient := radarr.NewRadarrClient(conf.RadarrURL, conf.RadarrAPIKey)

	_, err := radarrClient.GetMovies()
	if err != nil {
		fmt.Println(tui.RenderConnectionError("Radarr", conf.RadarrURL, err))
		os.Exit(1)
	}

	return radarrClient
}

// --- Static table columns for `get` command (still uses bubbles table) ---

func movieTableColumns() []table.Column {
	return []table.Column{
		{Title: "Title", Width: 30},
		{Title: "Year", Width: 6},
		{Title: "Runtime", Width: 9},
		{Title: "Genres", Width: 20},
		{Title: "IMDb", Width: 6},
		{Title: "TMDB", Width: 6},
		{Title: "RT", Width: 5},
		{Title: "Size", Width: 10},
	}
}

func movieToRow(movie radarr.Movie) table.Row {
	return table.Row{
		tui.TruncateString(movie.Title, 28),
		fmt.Sprintf("%d", movie.Year),
		fmt.Sprintf("%dm", movie.Runtime),
		tui.FormatGenres(movie.Genres, 18),
		tui.FormatRating(movie.Ratings.IMDb.Value, "imdb"),
		tui.FormatRating(movie.Ratings.TMDB.Value, "tmdb"),
		tui.FormatRating(movie.Ratings.RottenTomatoes.Value, "rt"),
		tui.FormatSizeGB(movie.Statistics.SizeOnDisk),
	}
}

// --- Interactive list columns and items for `searchanddelete` ---

func movieColumns() []tui.ColumnDef {
	return []tui.ColumnDef{
		{Title: "Title", Width: 30},
		{Title: "Year", Width: 6},
		{Title: "Runtime", Width: 9},
		{Title: "Genres", Width: 20},
		{Title: "IMDb", Width: 6},
		{Title: "TMDB", Width: 6},
		{Title: "RT", Width: 5},
		{Title: "Size", Width: 10},
	}
}

func movieToListItem(movie radarr.Movie) tui.ListItem {
	columns := []string{
		tui.TruncateString(movie.Title, 28),
		fmt.Sprintf("%d", movie.Year),
		fmt.Sprintf("%dm", movie.Runtime),
		tui.FormatGenres(movie.Genres, 18),
		tui.FormatRating(movie.Ratings.IMDb.Value, "imdb"),
		tui.FormatRating(movie.Ratings.TMDB.Value, "tmdb"),
		tui.FormatRating(movie.Ratings.RottenTomatoes.Value, "rt"),
		tui.FormatSizeGB(movie.Statistics.SizeOnDisk),
	}

	label := fmt.Sprintf("%s (%s)", movie.Title, tui.FormatSizeGB(movie.Statistics.SizeOnDisk))

	detail := []tui.KeyValue{
		{Key: "Title", Value: movie.Title},
		{Key: "Year", Value: fmt.Sprintf("%d", movie.Year)},
		{Key: "Runtime", Value: fmt.Sprintf("%d min", movie.Runtime)},
		{Key: "Studio", Value: movie.Studio},
		{Key: "Certification", Value: movie.Certification},
		{Key: "Genres", Value: strings.Join(movie.Genres, ", ")},
	}

	if movie.HasFile {
		mf := movie.MovieFile
		detail = append(detail,
			tui.KeyValue{Key: "Quality", Value: mf.Quality.Quality.Name},
			tui.KeyValue{Key: "Video", Value: mf.MediaInfo.VideoCodec},
			tui.KeyValue{Key: "Audio", Value: fmt.Sprintf("%s %.1fch", mf.MediaInfo.AudioCodec, mf.MediaInfo.AudioChannels)},
		)
		if mf.MediaInfo.VideoDynamicRange != "" {
			detail = append(detail, tui.KeyValue{Key: "HDR", Value: mf.MediaInfo.VideoDynamicRange})
		}
		if mf.MediaInfo.Resolution != "" {
			detail = append(detail, tui.KeyValue{Key: "Resolution", Value: mf.MediaInfo.Resolution})
		}
	}

	imdbStr := tui.FormatRating(movie.Ratings.IMDb.Value, "imdb")
	if movie.Ratings.IMDb.Votes > 0 {
		imdbStr += fmt.Sprintf(" (%d votes)", movie.Ratings.IMDb.Votes)
	}
	detail = append(detail,
		tui.KeyValue{Key: "IMDb", Value: imdbStr},
		tui.KeyValue{Key: "TMDB", Value: tui.FormatRating(movie.Ratings.TMDB.Value, "tmdb")},
		tui.KeyValue{Key: "Rotten Tomatoes", Value: tui.FormatRating(movie.Ratings.RottenTomatoes.Value, "rt")},
		tui.KeyValue{Key: "Size", Value: tui.FormatSizeGB(movie.Statistics.SizeOnDisk)},
		tui.KeyValue{Key: "Path", Value: movie.Path},
	)

	return tui.ListItem{
		Columns:  columns,
		Label:    label,
		Detail:   detail,
		Children: nil,
	}
}

func HandleGet(radarrAPIKey string, overseerAPIKey string, limit int, skip int) {
	conf := config.GetConfig("", radarrAPIKey, overseerAPIKey)
	if len(radarrAPIKey) > 0 {
		conf.RadarrAPIKey = radarrAPIKey
	}
	if len(overseerAPIKey) > 0 {
		conf.OverseerAPIKey = overseerAPIKey
	}

	radarrClient := initializeRadarrClient(conf)

	radarrMovies, err := radarrClient.GetMovies()
	if err != nil {
		fmt.Println(tui.RenderConnectionError("Radarr", conf.RadarrURL, err))
		os.Exit(1)
	}

	sort.Slice(radarrMovies, func(i, j int) bool {
		return radarrMovies[i].SizeOnDisk > radarrMovies[j].SizeOnDisk
	})

	var rows []table.Row
	for i := skip; i < len(radarrMovies) && i < skip+limit; i++ {
		rows = append(rows, movieToRow(radarrMovies[i]))
	}

	header := tui.TitleStyle.Render(fmt.Sprintf("Movies (%d total, showing %d)", len(radarrMovies), len(rows)))
	fmt.Printf("\n%s\n\n%s\n", header, tui.RenderStaticTable(movieTableColumns(), rows))
}

func HandleSearchAndDelete(radarrAPIKey string, overseerAPIKey string, limit int, skip int) {
	conf := config.GetConfig(radarrAPIKey, "", overseerAPIKey)
	radarrClient := radarr.NewRadarrClient(conf.RadarrURL, conf.RadarrAPIKey)
	overseerClient := overseer.NewOverseerClient(conf.OverseerURL, conf.OverseerAPIKey)

	var movieData []radarr.Movie

	fetchFn := func() ([]tui.ListItem, error) {
		movies, err := radarrClient.GetMovies()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch movies from Radarr: %w", err)
		}

		sort.Slice(movies, func(i, j int) bool {
			return movies[i].SizeOnDisk > movies[j].SizeOnDisk
		})

		end := skip + limit
		if end > len(movies) {
			end = len(movies)
		}
		movieData = movies[skip:end]

		var items []tui.ListItem
		for _, m := range movieData {
			items = append(items, movieToListItem(m))
		}
		return items, nil
	}

	deleteFn := func(selections []tui.Selection) []string {
		var results []string

		overseerMedia, err := overseerClient.GetMedia()
		if err != nil {
			results = append(results, tui.WarningStyle.Render("Could not fetch Overseer media: "+err.Error()))
		}

		for _, sel := range selections {
			if sel.ItemIndex >= len(movieData) {
				continue
			}
			movie := movieData[sel.ItemIndex]

			// Delete from Overseer
			if overseerMedia != nil {
				if mediaItem, err := FindMediaItemByTmdbID(movie.TMDBID, overseerMedia); err == nil {
					if err := overseerClient.DeleteMedia(mediaItem.Id); err != nil {
						results = append(results, tui.ErrorStyle.Render(fmt.Sprintf("Overseer delete failed for '%s': %s", movie.Title, err)))
					} else {
						results = append(results, tui.SuccessStyle.Render(fmt.Sprintf("Deleted '%s' from Overseer", movie.Title)))
					}
				}
			}

			// Delete from Radarr
			if err := radarrClient.DeleteMovie(movie.ID); err != nil {
				results = append(results, tui.ErrorStyle.Render(fmt.Sprintf("Radarr delete failed for '%s': %s", movie.Title, err)))
			} else {
				results = append(results, tui.SuccessStyle.Render(fmt.Sprintf("Deleted '%s' from Radarr", movie.Title)))
			}
		}
		return results
	}

	model := tui.NewSelectTableModel("Select Movies to Delete", movieColumns(), fetchFn, deleteFn)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println(tui.ErrorStyle.Render("Fatal error: " + err.Error()))
		os.Exit(1)
	}
}
