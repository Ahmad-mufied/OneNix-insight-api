package handler

import (
	"encoding/json"
	"google-custom-search/model"
	"google-custom-search/repository"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// NewsHandler handles requests related to news.
// It contains a cache client and a DynamoDB client.
type NewsHandler struct {
	Cache *repository.MemcachedClient
	DB    *repository.DynamoDBClient
}

// GetLatestNews handles the API request to get the latest news.
// It first checks the cache for the latest news. If the news is not in the cache,
// it fetches the news from DynamoDB, caches the result, and returns it.
// If there is an error during any of these operations, it returns an appropriate error response.
func (h *NewsHandler) GetLatestNews(c echo.Context) error {
	cacheKey := "latest_news"
	cachedData, err := h.Cache.Get(cacheKey)
	if err == nil {
		// Return cached data if available
		var newsList []model.News
		if err := json.Unmarshal(cachedData, &newsList); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse cached data"})
		}
		return c.JSON(http.StatusOK, newsList)
	}

	// Cache miss, fetch data from DynamoDB
	newsList, err := h.DB.GetAllNews()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch news"})
	}
	log.Println("Serving news from DynamoDB.")

	// Cache the data
	cacheData, err := json.Marshal(newsList)
	if err == nil {
		log.Println("Caching news data.")
		err := h.Cache.Set(cacheKey, cacheData, int32(24*time.Hour.Seconds()))
		if err != nil {
			log.Printf("Error caching data: %v", err)
		}
	}

	return c.JSON(http.StatusOK, newsList)
}
