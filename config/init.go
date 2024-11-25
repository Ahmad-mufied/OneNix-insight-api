package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// init initializes the configuration by reading from a .env file or environment variables.
// It sets the following global variables:
// - GoogleCustomSearchEngineAPIKey: API key for Google Custom Search Engine
// - GoogleCustomSearchEngineID: ID for Google Custom Search Engine
// - MemcachedServer: Address of the Memcached server
// - DynamodbRegion: AWS region for DynamoDB
func init() {
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

	// Set global variables from the configuration
	GoogleCustomSearchEngineAPIKey = v.GetString("GOOGLE_CUSTOM_SEARCH_ENGINE_API_KEY")
	GoogleCustomSearchEngineID = v.GetString("GOOGLE_CUSTOM_SEARCH_ENGINE_ID")
	MemcachedServer = fmt.Sprintf("%s:%s", v.GetString("MEMCACHED_HOST"), v.GetString("MEMCACHED_PORT"))
	DynamodbRegion = v.GetString("DYNAMODB_REGION")
}
