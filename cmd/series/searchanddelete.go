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

// searchAndDeleteCmd represents the searchanddelete subcommand
var searchAndDeleteCmd = &cobra.Command{
	Use:   "searchanddelete",
	Short: "Search and delete shows/series",
	Long:  `Search for movies based on criteria and delete them from the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		series.HandleSearchAndDeleteSeries(sonarrAPIKey, overseerAPIKey, limit)
	},
}

func init() {
	// Set default values to environment variables or fallback to empty strings

	searchAndDeleteCmd.Flags().StringVar(&sonarrAPIKey, "radarr-api-key", "", "API key for Radarr")
	searchAndDeleteCmd.Flags().StringVar(&overseerAPIKey, "overseer-api-key", "", "API key for Overseer")
	searchAndDeleteCmd.Flags().IntVar(&limit, "limit", 10, "Limit of movies to show")

	SeriesCommand.AddCommand(searchAndDeleteCmd)
}
