package movies

import (
	"flashbacklabsio/fcli/internal/movies"

	"github.com/spf13/cobra"
)

var (
	radarrAPIKey   string
	overseerAPIKey string
	limit          int
	skip           int
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

func init() {
	// Set default values to environment variables or fallback to empty strings

	MoviesCmd.PersistentFlags().IntVar(&limit, "limit", 10, "Limit of movies to show")
	MoviesCmd.PersistentFlags().StringVar(&radarrAPIKey, "radarr-api-key", "", "API key for Radarr")
	MoviesCmd.PersistentFlags().StringVar(&overseerAPIKey, "overseer-api-key", "", "API key for Overseer")
	MoviesCmd.PersistentFlags().IntVar(&skip, "skip", 0, "Pagination skip. Start printing after the skip.")
}
