package repository

import "github.com/luigibarbato/isolated-think-source/internal/app"

type Repository interface {
	// Search returns the results of a search query.
	Find(query string) ([]app.Feed, error)
	// Index indexes a document.
	Save(doc app.Feed) error
	// Delete deletes a document.
	Delete(doc app.Feed) error
	// Update updates a document.
	Update(doc app.Feed) error
	// Index indexes a document.
	Index(id string, doc app.Feed) error
	// Close closes the database.
	Close() error
}
