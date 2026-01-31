package series

import (
	"flashbacklabsio/fcli/internal/series"

	"github.com/spf13/cobra"
)

var searchAndDeleteCmd = &cobra.Command{
	Use:   "searchanddelete",
	Short: "Search and delete shows/series",
	Long:  `Search for series and delete them or individual seasons.`,
	Run: func(cmd *cobra.Command, args []string) {
		series.HandleSearchAndDeleteSeries(sonarrAPIKey, overseerAPIKey, limit)
	},
}

func init() {
	SeriesCommand.AddCommand(searchAndDeleteCmd)
}
