package app

import "time"

type FeedRepository interface {
	// Search returns the results of a search query.
	Find(query string) ([]Feed, error)
	// Index indexes a document.
	Save(doc Feed) error
	// Delete deletes a document.
	Delete(doc Feed) error
	// Close closes the database.
	Close() error
}

type Feed struct {
	Title    string    `json:"title"`
	Link     string    `json:"link"`
	Language string    `json:"language"`
	Image    *Image    `json:"image"`
	Summary  string    `json:"summary"`
	Source   string    `json:"source"`
	Date     time.Time `json:"date"`
}

type Image struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
