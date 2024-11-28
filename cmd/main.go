package main

import (
	"context"
	"errors"
	"github.com/go-co-op/gocron"
	"github.com/labstack/echo/v4"
	"google-custom-search/config"
	"google-custom-search/handler"
	"google-custom-search/repository"
	"google-custom-search/router"
	"google-custom-search/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	memcachedRepo := repository.NewMemcachedClient(config.MemcachedServer)
	mongoRepo := repository.NewMongoRepository(config.DB.Database("news").Collection("news"))
	newsHandler := handler.NewsHandler{Cache: memcachedRepo, DB: mongoRepo}

	if config.AutoFetchSwitch {
		initializeAndStartCrawlerService(mongoRepo)
	}

	e := echo.New()
	startAndGracefullyStopServer(e, &newsHandler)
}

func startAndGracefullyStopServer(e *echo.Echo, newsHandler *handler.NewsHandler) {
	// Register routes
	router.RegisterRoutes(e, newsHandler)

	port := config.Viper.GetString("WEB_SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...", port)

	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: e,
	}

	go func() {
		if err := e.StartServer(server); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

func initializeAndStartCrawlerService(mongoRepo *repository.MongoRepository) {
	gSearchAPI := utils.GoogleSearchAPI{
		DB: mongoRepo,
	}

	// Set up cron job for daily crawler execution
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Day().Do(func() {
		log.Println("Executing daily crawler task...")
		gSearchAPI.RunTask()
	})
	scheduler.StartAsync()
}
