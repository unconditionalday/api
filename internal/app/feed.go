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
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Link     string    `json:"link"`
	Language string    `json:"language"`
	Image    *Image    `json:"image"`
	Summary  string    `json:"summary"`
	Source   string    `json:"source"`
	Date     time.Time `json:"date"`
	Related  []string  `json:"related"`
}

func (f Feed) IsValid() bool {
	if f.ID == "" || f.Title == "" || f.Link == "" || f.Source == "" || f.Date.IsZero() {
		return false
	}

	return true
}

type Image struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
