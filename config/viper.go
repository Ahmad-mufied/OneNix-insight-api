package config

import (
	"github.com/spf13/viper"
	"log"
)

var Viper *viper.Viper

func InitViper() {
	// Initialize Viper
	v := viper.New()

	// Set the configuration file to .env
	v.SetConfigFile(".env")
	v.SetConfigType("dotenv")
	v.AddConfigPath("./")

	if err := v.ReadInConfig(); err != nil {
		log.Println(err)
		log.Println("Loaded From Environment Variables")
		v.AutomaticEnv()
	} else {
		log.Println("Loaded .env file")
	}

	Viper = v
}
