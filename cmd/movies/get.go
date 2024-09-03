package movies

import (
	"flashbacklabsio/fcli/internal/movies"

	"github.com/spf13/cobra"
)

// searchAndDeleteCmd represents the searchanddelete subcommand
var getCommand = &cobra.Command{
	Use:   "get",
	Short: "gets movies from radarr API.",
	Long:  `Search for movies based on criteria.`,
	Run: func(cmd *cobra.Command, args []string) {
		movies.HandleGet(radarrAPIKey, overseerAPIKey, limit, skip)
	},
}

func init() {
	MoviesCmd.AddCommand(getCommand)
}
