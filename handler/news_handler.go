package handler

import (
	"github.com/labstack/echo/v4"
	"google-custom-search/model"
	"google-custom-search/repository"
	"google-custom-search/utils"
	"log"
	"net/http"
)

// NewsHandler handles requests related to news.
// It contains a cache client and a DynamoDB client.
type NewsHandler struct {
	Cache *repository.MemcachedClient
	DB    *repository.MongoRepository
}

func (h *NewsHandler) GetLatestNews(c echo.Context) error {
	ctx := c.Request().Context()

	// Extract filters from query parameters
	filters := map[string]string{
		"country": c.QueryParam("country"),
		"degree":  c.QueryParam("degree"),
		"major":   c.QueryParam("major"),
	}

	// Check cache first
	newsList, err := h.Cache.GetCachedList(filters)
	if err == nil {
		log.Println("Cache hit")
		return c.JSON(http.StatusOK, model.JSONResponse{
			Status:  http.StatusOK,
			Message: "Success fetching news",
			Count:   len(newsList),
			Data:    newsList,
		})
	}

	// If cache miss, query MongoDB
	log.Println("Cache miss")
	newsList, err = h.DB.List(ctx, filters)
	if err != nil {
		return utils.HandleError(c, utils.NewAPIError(http.StatusInternalServerError, "Failed to fetch news", nil))
	}

	// Check if the list is empty and return an error
	if len(newsList) == 0 {
		return utils.HandleError(c, utils.NewAPIError(http.StatusNotFound, "No news found", nil))
	}

	// Store result in cache
	log.Println("Storing result in cache")
	err = h.Cache.SetCachedList(filters, newsList)
	if err != nil {
		log.Printf("Error storing result in cache: %v\n", err)
		return utils.HandleError(c, utils.NewAPIError(http.StatusInternalServerError, "Failed to fetch news", nil))
	}

	return c.JSON(http.StatusOK, model.JSONResponse{
		Status:  http.StatusOK,
		Message: "Success fetching news",
		Count:   len(newsList),
		Data:    newsList,
	})
}
