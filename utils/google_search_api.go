package utils

import (
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

type GoogleSearchAPI struct {
	DB *repository.MongoRepository
}

func FetchNews(query string) ([]*model.News, error) {
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

	var newsList []*model.News
	for _, item := range result.Items {
		thumbnail := ""
		if len(item.Pagemap.CseThumbnail) > 0 {
			thumbnail = item.Pagemap.CseThumbnail[0].Src
		}
		newsList = append(newsList, &model.News{
			Title:     item.Title,
			Date:      time.Now().Format(time.RFC3339),
			Thumbnail: thumbnail,
			Snippet:   item.Snippet,
			Link:      item.Link,
		})
	}

	return newsList, nil
}

func (g *GoogleSearchAPI) FetchAndSaveNews(query, country, degree, major string) {
	results, err := FetchNews(query)
	if err != nil {
		log.Printf("Error fetching news for query %s: %v", query, err)
		return
	}

	err = g.DB.SaveNews(query, country, degree, major, results)
	if err != nil {
		log.Printf("Error saving news for query %s: %v", query, err)
	}
}

func (g *GoogleSearchAPI) RunTask() {

	// Drop the collection before fetching fresh data
	log.Println("Dropping collection")
	err := g.DB.DropCollection()
	if err != nil {
		log.Printf("Error dropping collection: %v", err)
	}

	//countries := []string{"Germany", "United States", "Malaysia", "Australia"}
	//degrees := []string{"Diploma", "Bachelor", "Master", "Doctoral"}
	//majors := []string{"Art", "Science", "Social"}
	countries := []string{"Australia"}
	degrees := []string{"Diploma"}
	majors := []string{"Art"}

	for _, country := range countries {
		for _, degree := range degrees {
			for _, major := range majors {
				query := fmt.Sprintf("study abroad news %s %s %s", country, degree, major)
				g.FetchAndSaveNews(query, country, degree, major)
			}
		}
	}
}
