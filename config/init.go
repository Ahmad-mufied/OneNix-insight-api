package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func init() {

	// Initialize Viper
	LoadConfig()

	// Set global variables from the configuration
	GoogleCustomSearchEngineAPIKey = Viper.GetString("GOOGLE_CUSTOM_SEARCH_ENGINE_API_KEY")
	GoogleCustomSearchEngineID = Viper.GetString("GOOGLE_CUSTOM_SEARCH_ENGINE_ID")
	MemcachedServer = fmt.Sprintf("%s:%s", Viper.GetString("MEMCACHED_HOST"), Viper.GetString("MEMCACHED_PORT"))
	// Init Mongo
	InitMongo()

	AutoFetchSwitch = Viper.GetBool("AUTO_FETCH_SWITCH")
}

func LoadConfig() {
	v := viper.New()
	Viper = v
	// Set the configuration file to .env
	v.SetConfigFile(".env")
	v.SetConfigType("dotenv")
	v.AddConfigPath("./")

	// Read the configuration file
	if err := v.ReadInConfig(); err != nil {
		log.Println(err)
		log.Println("Loaded From Environment Variables")
		v.AutomaticEnv()
	} else {
		log.Println("Loaded .env file")
	}
}
