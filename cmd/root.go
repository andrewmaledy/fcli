package cmd

import (
	"fmt"
	"os"

	"flashbacklabsio/fcli/cmd/configure"
	"flashbacklabsio/fcli/cmd/movies"
	"flashbacklabsio/fcli/cmd/series"
	"flashbacklabsio/fcli/internal/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fcli",
	Short: "fcli is a CLI tool for flashbacklabsio",
	Long:  `A CLI tool for managing different services and commands for flashbacklabsio.`,
	// Load config before any subcommand runs. The configure command overrides
	// this with its own PersistentPreRunE so it can run without a config file.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := config.InitConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return nil
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(movies.MoviesCmd)
	rootCmd.AddCommand(series.SeriesCommand)
	rootCmd.AddCommand(configure.ConfigureCmd)
}
