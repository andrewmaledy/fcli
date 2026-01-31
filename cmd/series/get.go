package series

import (
	"flashbacklabsio/fcli/internal/series"

	"github.com/spf13/cobra"
)

var getCommand = &cobra.Command{
	Use:   "get",
	Short: "Gets series from Sonarr API",
	Long:  `Retrieve and display series from your Sonarr library.`,
	Run: func(cmd *cobra.Command, args []string) {
		series.HandleGet(sonarrAPIKey, overseerAPIKey, limit)
	},
}

func init() {
	SeriesCommand.AddCommand(getCommand)
}
