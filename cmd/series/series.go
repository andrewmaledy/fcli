package series

import (
	"flashbacklabsio/fcli/internal/series"

	"github.com/spf13/cobra"
)

// MoviesCmd represents the movies command
var SeriesCommand = &cobra.Command{
	Use:   "series",
	Short: "Manage movies",
	Long:  `Manage movies through various operations like listing, adding, or deleting.`,
	Run: func(cmd *cobra.Command, args []string) {
		series.HandleSeriesCommand()
	},
}
