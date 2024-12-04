package handler

import (
	"github.com/labstack/echo/v4"
	"google-custom-search/config"
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

	// check if country is US then set the country to United States
	if filters["country"] == "US" {
		filters["country"] = "United States"
	}

	var newsList []model.News

	// Check cache not enabled
	if !config.CacheSwitch {
		log.Println("Cache not enabled")
		newsList, err := h.DB.List(ctx, filters)
		if err != nil {
			return utils.HandleError(c, utils.NewAPIError(http.StatusInternalServerError, "Failed to fetch news", nil))
		}

		// Check if the list is empty and return an error
		if len(newsList) == 0 {
			return utils.HandleError(c, utils.NewAPIError(http.StatusNotFound, "No news found", nil))
		}

		return c.JSON(http.StatusOK, model.JSONResponse{
			Status:  http.StatusOK,
			Message: "Success fetching news",
			Count:   len(newsList),
			Data:    newsList,
		})
	}

	// Check cache enabled
	newsList, err := h.Cache.GetCachedList(filters)
	if err != nil {
		log.Println("Cache miss")
		newsList, err = h.DB.List(ctx, filters)
		if err != nil {
			return utils.HandleError(c, utils.NewAPIError(http.StatusInternalServerError, "Failed to fetch news", nil))
		}

		// Check if the list is empty and return an error
		if len(newsList) == 0 {
			return utils.HandleError(c, utils.NewAPIError(http.StatusNotFound, "No news found", nil))
		}

		// Set the list to cache
		err = h.Cache.SetCachedList(filters, newsList)
		if err != nil {
			log.Println("Failed to set cache")
		}
	}

	return c.JSON(http.StatusOK, model.JSONResponse{
		Status:  http.StatusOK,
		Message: "Success fetching news",
		Count:   len(newsList),
		Data:    newsList,
	})
}
