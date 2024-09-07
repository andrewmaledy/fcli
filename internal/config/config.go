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
func assignAPIKey(key string, fallback string) string {
	if len(key) > 0 {
		return key
	}
	return viper.GetString(fallback)
}

// GetConfig returns a Configuration struct populated with values from viper.
func GetConfig(sonarrAPIKey string, radarrAPIKey string, overseerAPIKey string) *Configuration {
	InitConfig()

	conf := &Configuration{
		RadarrURL:    viper.GetString("radarr.url"),
		OverseerURL:  viper.GetString("overseer.url"),
		SonarrAPIKey: viper.GetString("sonarr.apiKey"),
		SonarrURL:    viper.GetString("sonarr.url"),
	}
	conf.SonarrAPIKey = assignAPIKey(sonarrAPIKey, "sonarr.apiKey")
	conf.RadarrAPIKey = assignAPIKey(radarrAPIKey, "radarr.apiKey")
	conf.OverseerAPIKey = assignAPIKey(overseerAPIKey, "overseer.apiKey")
	return conf
}
