package movies

import (
	"flashbacklabsio/fcli/internal/movies"

	"github.com/spf13/cobra"
)

// searchAndDeleteCmd represents the searchanddelete subcommand
var searchAndDeleteCmd = &cobra.Command{
	Use:   "searchanddelete",
	Short: "Search and delete movies",
	Long:  `Search for movies based on criteria and delete them from the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		movies.HandleSearchAndDelete(radarrAPIKey, overseerAPIKey, limit, skip)
	},
}

func init() {
	// Set default values to environment variables or fallback to empty strings
	MoviesCmd.AddCommand(searchAndDeleteCmd)
}
