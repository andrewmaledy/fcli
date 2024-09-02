package config

import (
	"log"
	"os/user"

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

// InitConfig initializes the configuration using viper.
func InitConfig() {
	viper.SetConfigName(".fcli-config")
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(usr.HomeDir)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

// GetConfig returns a Configuration struct populated with values from viper.
func GetConfig() *Configuration {
	return &Configuration{
		RadarrURL:      viper.GetString("radarr.url"),
		RadarrAPIKey:   viper.GetString("radarr.apiKey"),
		OverseerURL:    viper.GetString("overseer.url"),
		OverseerAPIKey: viper.GetString("overseer.apiKey"),
		SonarrAPIKey:   viper.GetString("sonarr.apiKey"),
		SonarrURL:      viper.GetString("sonarr.url"),
	}
}
