package model

type News struct {
	ID        string `json:"id"`        // Unique ID (e.g., hash of URL)
	Title     string `json:"title"`     // Title of the news
	Date      string `json:"date"`      // Date of publication
	Thumbnail string `json:"thumbnail"` // Image thumbnail URL
	Snippet   string `json:"snippet"`   // First paragraph or description
	Link      string `json:"link"`      // News article URL
}
