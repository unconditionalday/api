package app

import (
	"time"
)

type FeedRepository interface {
	// Search returns the results of a search query.
	FindByKeyword(query string) ([]Feed, error)
	// Search returns the results of a search query by similarity.
	FindBySimilarity(doc Feed) ([]Feed, error)
	FindByID(id string) (Feed, error)
	// Index indexes a document.
	Save(doc Feed) error
	// Update a document in index.
	Update(docs ...Feed) error
	// Count the number of documents indexed.
	Count() uint64
	// Delete deletes a document.
	Delete(doc Feed) error
}

type Feed struct {
	ID       string    `json:"id"`
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
