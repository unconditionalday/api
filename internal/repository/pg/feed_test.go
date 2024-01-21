package pg_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq" // Importa il driver PostgreSQL

	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/repository/pg"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// Imposta la connessione al tuo database di test
	dbUser := os.Getenv("UNCONDITIONAL_API_DATABASE_USER")
	dbName := os.Getenv("UNCONDITIONAL_API_DATABASE_NAME")
	dbPassword := os.Getenv("UNCONDITIONAL_API_DATABASE_PASSWORD")

	dbConfig := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		t.Fatal(err)
	}

	// Esegui migrazioni o inizializza lo schema del database se necessario
	// ...

	// Funzione per pulire il database dopo i test
	cleanup := func() {
		// Pulisci il database (elimina dati di test, ecc.)
		// ...
		db.Close()
	}

	return db, cleanup
}

func TestSave(t *testing.T) {
	// Imposta il database di test
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Inizializza il repository
	f := pg.NewFeedRepository(db)

	testCases := []struct {
		name     string
		document app.Feed
		wantErr  bool
	}{
		{
			name: "document is saved",
			document: app.Feed{
				Title:    "test",
				Link:     "link",
				Language: "it",
				Image:    &app.Image{},
				Summary:  "Lorem Ipsum",
				Source:   "Unconditional Day",
				Date:     time.Time{},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := f.Save(tc.document)

			if (err != nil) != tc.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

func TestFind(t *testing.T) {
	// Imposta il database di test
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Inizializza il repository
	f := pg.NewFeedRepository(db)

	testCases := []struct {
		name     string
		document app.Feed
		wantErr  bool
	}{
		{
			name: "document is found",
			document: app.Feed{
				Title:    "test",
				Link:     "link",
				Language: "it",
				Image:    &app.Image{},
				Summary:  "Lorem Ipsum",
				Source:   "Unconditional Day",
				Date:     time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := f.Save(tc.document)
			if (err != nil) != tc.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			found, err := f.Find(tc.document.Title)
			if (err != nil) != tc.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if len(found) == 0 {
				t.Errorf("expected document to be found")
			}

			if found[0].Title != tc.document.Title {
				t.Errorf("expected document to be found")
			}
		})
	}
}
