package main

import (
	"github.com/go-co-op/gocron"
	"github.com/labstack/echo/v4"
	"google-custom-search/config"
	"google-custom-search/handler"
	"google-custom-search/repository"
	"google-custom-search/router"
	"google-custom-search/utils"
	"log"
	"time"
)

func main() {
	// Load environment variables
	memcachedHost := config.MemcachedServer // ElasticCache endpoint
	dynamoRegion := config.DynamodbRegion   // AWS region

	// Initialize Memcached client with ElasticCache endpoint
	memcachedClient := repository.NewMemcachedClient(memcachedHost)

	// Initialize DynamoDB client
	dynamoDBClient, err := repository.NewDynamoDBClient(dynamoRegion)
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB client: %v", err)
	}

	// Initialize crawler service
	crawlerService := utils.Crawler{
		Cache: memcachedClient,
		DB:    dynamoDBClient,
	}

	// Set up Echo
	e := echo.New()

	// Initialize API handlers
	newsHandler := handler.NewsHandler{
		Cache: memcachedClient,
		DB:    dynamoDBClient,
	}

	// Register routes
	router.RegisterRoutes(e, &newsHandler)

	// Set up cron job for daily crawler execution
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Day().Do(func() {
		log.Println("Executing daily crawler task...")
		crawlerService.FetchAndSaveNews()
	})
	scheduler.StartAsync()

	// Start Echo server

	port := config.Viper.GetString("WEB_SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
