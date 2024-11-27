package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"google-custom-search/config"
	"google-custom-search/model"
	"google-custom-search/repository"
	"log"
	"net/http"
	"strings"
	"time"
)

type Crawler struct {
	Cache *repository.MemcachedClient
	DB    *repository.DynamoDBClient
}

func (c *Crawler) FetchAndSaveNews() {

	// Fetch fresh data
	query := "study abroad news"
	results, err := searchGoogle(query)
	if err != nil {
		log.Println("Error fetching news:", err)
		return
	}

	cacheKey := "latest_news"
	cachedData, err := c.Cache.Get(cacheKey)
	if err == nil {
		log.Println("Serving news from cache.")
		// Cached data exists, return it
		var newsList []model.News
		if err := json.Unmarshal(cachedData, &newsList); err != nil {
			log.Printf("Error unmarshalling cached data: %v", err)
		} else {
			return // Skip saving if cache is valid
		}
	}

	var savedNews []model.News
	for _, news := range results {
		// Check for duplicates before saving to the database
		exists, err := c.DB.CheckNewsExists(news.ID)
		if err != nil {
			log.Printf("Failed to check news: %s, error: %v\n", news.Title, err)
			continue
		}

		if !exists {
			err = c.DB.SaveNews(news)
			if err != nil {
				log.Printf("Failed to save news: %s, error: %v\n", news.Title, err)
			} else {
				savedNews = append(savedNews, news)
				log.Printf("News saved: %s\n", news.Title)
			}
		}

	}

	log.Println("Serving news from Google.")

	// Cache the newly fetched data
	if len(savedNews) > 0 {
		cacheData, err := json.Marshal(savedNews)
		if err != nil {
			log.Printf("Error marshalling data for cache: %v", err)
		} else {
			c.Cache.Set(cacheKey, cacheData, int32(24*time.Hour.Seconds())) // Cache for 24 hours
		}
	}
}

func searchGoogle(query string) ([]model.News, error) {
	query = strings.ReplaceAll(query, " ", "+")
	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/customsearch/v1?q=%s&key=%s&cx=%s",
		query, config.GoogleCustomSearchEngineAPIKey, config.GoogleCustomSearchEngineID,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Error: %d", resp.StatusCode)
	}

	log.Println("Fetched news from Google")
	var result struct {
		Items []struct {
			Title   string `json:"title"`
			Link    string `json:"link"`
			Snippet string `json:"snippet"`
			Pagemap struct {
				CseThumbnail []struct {
					Src string `json:"src"`
				} `json:"cse_thumbnail"`
			} `json:"pagemap"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var newsList []model.News
	for _, item := range result.Items {
		thumbnail := ""
		if len(item.Pagemap.CseThumbnail) > 0 {
			thumbnail = item.Pagemap.CseThumbnail[0].Src
		}
		newsList = append(newsList, model.News{
			ID:        GenerateID(item.Link),
			Title:     item.Title,
			Date:      time.Now().Format(time.RFC3339),
			Thumbnail: thumbnail,
			Snippet:   item.Snippet,
			Link:      item.Link,
		})
	}

	return newsList, nil
}

func GenerateID(url string) string {
	hash := sha256.Sum256([]byte(url))
	return hex.EncodeToString(hash[:])
}
