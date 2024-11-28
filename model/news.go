package model

type News struct {
	Title     string `json:"title" bson:"title"`
	Date      string `json:"date" bson:"date"`
	Thumbnail string `json:"thumbnail" bson:"thumbnail"`
	Snippet   string `json:"snippet" bson:"snippet"`
	Link      string `json:"link" bson:"link"`
}

type SearchResult struct {
	Country string  `bson:"country"`
	Degree  string  `bson:"degree"`
	Major   string  `bson:"major"`
	Results []*News `bson:"news"`
}
