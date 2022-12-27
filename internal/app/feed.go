package app

import (
	"time"
)

type FeedRepository interface {
	// Search returns the results of a search query.
	Find(query string) ([]Feed, error)
	// Index indexes a document.
	Save(doc Feed) error
	// Update a document in index.
	Update(doc Feed) error
	// Count the number of documents indexed.
	Count() uint64
	// Delete deletes a document.
	Delete(doc Feed) error
}

type FeedRepository interface {
	// Search returns the results of a search query.
	Find(query string) ([]Feed, error)
	// Index indexes a document.
	Save(doc Feed) error
	// Delete deletes a document.
	Delete(doc Feed) error
	// Update updates a document.
	Update(doc Feed) error
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

func (f Feed) IsValid() bool {
	if f.Title == "" || f.Link == "" || f.Source == "" || f.Date.IsZero() {
		return false
	}

	return true
}

type Image struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
