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
	// Set default values to environment variables or fallback to empty strings

	//getCommand.Flags().StringVar(&radarrAPIKey, "radarr-api-key", "", "API key for Radarr")
	//getCommand.Flags().StringVar(&overseerAPIKey, "overseer-api-key", "", "API key for Overseer")
	//getCommand.Flags().IntVar(&limit, "limit", 10, "Limit of movies to show")

	MoviesCmd.AddCommand(getCommand)
}
