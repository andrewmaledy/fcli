package cmd

import (
	"os"

	"flashbacklabsio/fcli/cmd/movies"
	"flashbacklabsio/fcli/cmd/series"
	"flashbacklabsio/fcli/internal/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fcli",
	Short: "fcli is a CLI tool for flashbacklabsio",
	Long:  `A CLI tool for managing different services and commands for flashbacklabsio.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add subcommands here
	config.InitConfig()
	rootCmd.AddCommand(movies.MoviesCmd)
	rootCmd.AddCommand(series.SeriesCommand)
}
