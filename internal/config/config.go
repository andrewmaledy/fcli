package config

import (
	"log"
	"os/user"

	"github.com/spf13/viper"
)

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
