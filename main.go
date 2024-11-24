package main

import (
	"encoding/json"
	"fmt"
	"google-custom-search/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type SearchResult struct {
	Items []struct {
		Title       string `json:"title"`
		Link        string `json:"link"`
		Snippet     string `json:"snippet"`
		DisplayLink string `json:"displayLink"`
	} `json:"items"`
}

func searchGoogle(query string) {
	// Replace spaces with '+' for the query
	query = url.QueryEscape(query)
	query = strings.ReplaceAll(query, "%20", "+")

	urlGoogle := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&key=%s&cx=%s",
		query,
		config.GoogleCustomSearchEngineAPIKey,
		config.GoogleCustomSearchEngineID,
	)

	resp, err := http.Get(urlGoogle)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP Error: %d\n", resp.StatusCode)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var result SearchResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, item := range result.Items {
		fmt.Printf("Title: %s\nLink: %s\nSnippet: %s\nDisplay Link: %s\n\n",
			item.Title, item.Link, item.Snippet, item.DisplayLink)
	}
}

func main() {
	searchGoogle("study abroad news")
}
