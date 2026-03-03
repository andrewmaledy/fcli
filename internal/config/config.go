package config

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/spf13/viper"
)

// Configuration holds the necessary API configuration.
type Configuration struct {
	RadarrURL      string
	RadarrAPIKey   string
	OverseerURL    string
	OverseerAPIKey string
	SonarrAPIKey   string
	SonarrURL      string
}

// NormalizeURL ensures a URL has the https:// scheme and no trailing slash.
// Accepts inputs like "sonarr.example.com", "http://sonarr.example.com", or "https://sonarr.example.com/api/v3".
func NormalizeURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	rawURL = strings.TrimRight(rawURL, "/")
	if rawURL == "" {
		return rawURL
	}
	if strings.HasPrefix(rawURL, "http://") {
		return "https://" + rawURL[7:]
	}
	if strings.HasPrefix(rawURL, "https://") {
		return rawURL
	}
	return "https://" + rawURL
}

// InitConfig initializes the configuration using viper.
// Returns an error if the config file is not found or cannot be read.
func InitConfig() error {
	viper.SetConfigName(".fcli-config")
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	viper.AddConfigPath(usr.HomeDir)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config file not found — run 'fcli configure' to set up your credentials\n(%w)", err)
	}
	return nil
}

// WriteConfig saves the given configuration to ~/.fcli-config.yaml.
func WriteConfig(conf *Configuration) error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	viper.SetConfigType("yaml")
	viper.Set("radarr.url", conf.RadarrURL)
	viper.Set("radarr.apiKey", conf.RadarrAPIKey)
	viper.Set("sonarr.url", conf.SonarrURL)
	viper.Set("sonarr.apiKey", conf.SonarrAPIKey)
	viper.Set("overseer.url", conf.OverseerURL)
	viper.Set("overseer.apiKey", conf.OverseerAPIKey)

	return viper.WriteConfigAs(usr.HomeDir + "/.fcli-config.yaml")
}

func assignAPIKey(key string, fallback string) string {
	if len(key) > 0 {
		return key
	}
	return viper.GetString(fallback)
}

// GetConfig returns a Configuration struct populated with values from viper.
func GetConfig(sonarrAPIKey string, radarrAPIKey string, overseerAPIKey string) *Configuration {
	return &Configuration{
		RadarrURL:      viper.GetString("radarr.url"),
		RadarrAPIKey:   assignAPIKey(radarrAPIKey, "radarr.apiKey"),
		SonarrURL:      viper.GetString("sonarr.url"),
		SonarrAPIKey:   assignAPIKey(sonarrAPIKey, "sonarr.apiKey"),
		OverseerURL:    viper.GetString("overseer.url"),
		OverseerAPIKey: assignAPIKey(overseerAPIKey, "overseer.apiKey"),
	}
}
