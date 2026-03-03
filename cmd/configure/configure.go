package configure

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"flashbacklabsio/fcli/internal/config"
	"flashbacklabsio/fcli/internal/tui"

	"github.com/spf13/cobra"
)

// ConfigureCmd runs an interactive setup wizard for fcli credentials.
var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure fcli service URLs and API keys",
	Long:  `Interactively set up API keys and URLs for Radarr, Sonarr, and Overseer/Jellyseerr.`,
	// Override parent PersistentPreRunE — no config file required to run configure.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	Run:               runConfigure,
}

// redact masks all but the last 6 characters of a token.
func redact(s string) string {
	if len(s) <= 6 {
		return strings.Repeat("*", len(s))
	}
	return "..." + s[len(s)-6:]
}

func prompt(scanner *bufio.Scanner, label, current string) string {
	if current != "" {
		fmt.Printf("  %s [%s]: ", tui.DetailKeyStyle.Render(label), tui.DimStyle.Render(current))
	} else {
		fmt.Printf("  %s: ", tui.DetailKeyStyle.Render(label))
	}
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return current
	}
	return input
}

func promptSecret(scanner *bufio.Scanner, label, current string) string {
	if current != "" {
		fmt.Printf("  %s [%s]: ", tui.DetailKeyStyle.Render(label), tui.DimStyle.Render(redact(current)))
	} else {
		fmt.Printf("  %s: ", tui.DetailKeyStyle.Render(label))
	}
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return current
	}
	return input
}

func runConfigure(cmd *cobra.Command, args []string) {
	fmt.Println()
	fmt.Println(tui.TitleStyle.Render("fcli Configuration Setup"))
	fmt.Println(tui.SubtitleStyle.Render("  URLs accept any format: sonarr.example.com, http://..., or https://..."))
	fmt.Println()

	// Load existing config if available (ignore errors — first-time setup may have none).
	_ = config.InitConfig()
	existing := config.GetConfig("", "", "")

	scanner := bufio.NewScanner(os.Stdin)

	// --- Radarr ---
	fmt.Println(tui.SelectionHeaderStyle.Render("Radarr"))
	radarrURL := config.NormalizeURL(prompt(scanner, "URL", existing.RadarrURL))
	radarrAPIKey := promptSecret(scanner, "API Key", existing.RadarrAPIKey)
	fmt.Println()

	// --- Sonarr ---
	fmt.Println(tui.SelectionHeaderStyle.Render("Sonarr"))
	sonarrURL := config.NormalizeURL(prompt(scanner, "URL", existing.SonarrURL))
	sonarrAPIKey := promptSecret(scanner, "API Key", existing.SonarrAPIKey)
	fmt.Println()

	// --- Overseer / Jellyseerr ---
	fmt.Println(tui.SelectionHeaderStyle.Render("Overseer / Jellyseerr"))
	overseerURL := config.NormalizeURL(prompt(scanner, "URL", existing.OverseerURL))
	overseerAPIKey := promptSecret(scanner, "API Key", existing.OverseerAPIKey)
	fmt.Println()

	conf := &config.Configuration{
		RadarrURL:      radarrURL,
		RadarrAPIKey:   radarrAPIKey,
		SonarrURL:      sonarrURL,
		SonarrAPIKey:   sonarrAPIKey,
		OverseerURL:    overseerURL,
		OverseerAPIKey: overseerAPIKey,
	}

	if err := config.WriteConfig(conf); err != nil {
		fmt.Println(tui.ErrorStyle.Render("Failed to save configuration: " + err.Error()))
		os.Exit(1)
	}

	fmt.Println(tui.SuccessStyle.Render("Configuration saved to ~/.fcli-config.yaml"))
	fmt.Println()
}
