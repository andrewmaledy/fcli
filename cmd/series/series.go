package series

import (
	"flashbacklabsio/fcli/internal/series"

	"github.com/spf13/cobra"
)

var (
	sonarrAPIKey   string
	overseerAPIKey string
	limit          int
)

// SeriesCommand represents the series command
var SeriesCommand = &cobra.Command{
	Use:   "series",
	Short: "Manage TV series",
	Long:  `Manage TV series through various operations like listing, adding, or deleting.`,
	Run: func(cmd *cobra.Command, args []string) {
		series.HandleSeriesCommand()
	},
}

func init() {
	SeriesCommand.PersistentFlags().IntVar(&limit, "limit", 10, "Limit of series to show")
	SeriesCommand.PersistentFlags().StringVar(&sonarrAPIKey, "sonarr-api-key", "", "API key for Sonarr")
	SeriesCommand.PersistentFlags().StringVar(&overseerAPIKey, "overseer-api-key", "", "API key for Overseer")
}
