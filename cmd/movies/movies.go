package movies

import (
	"flashbacklabsio/fcli/internal/movies"

	"github.com/spf13/cobra"
)

// MoviesCmd represents the movies command
var MoviesCmd = &cobra.Command{
	Use:   "movies",
	Short: "Manage movies",
	Long:  `Manage movies through various operations like listing, adding, or deleting.`,
	Run: func(cmd *cobra.Command, args []string) {
		movies.HandleMoviesCommand()
	},
}
