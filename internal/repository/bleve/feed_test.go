package bleve_test

import (
	"os"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2/mapping"

	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/repository/bleve"
	blevex "github.com/unconditionalday/server/internal/x/bleve"
)

func TestSave(t *testing.T) {
	testCases := []struct {
		name     string
		document app.Feed
		wantErr  bool
	}{
		{
			name:    "document is saved",
			wantErr: false,
			document: app.Feed{
				Title:    "test",
				Link:     "link",
				Language: "it",
				Image:    &app.Image{},
				Summary:  "Lorem Ipsum",
				Source:   "Unconditional Day",
				Date:     time.Time{},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := blevex.NewIndex("test.bleve", mapping.NewIndexMapping())
			if err != nil {
				t.Fatalf("expected bleve to be created")
			}

			f := bleve.NewFeedRepository(b)

			defer os.RemoveAll("test.bleve")

			err = f.Save(tc.document)

			if (err != nil) != tc.wantErr {
				t.Errorf("NewBleve() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

func TestFind(t *testing.T) {
	testCases := []struct {
		name     string
		document app.Feed
		wantErr  bool
	}{
		{
			name:    "document is found",
			wantErr: false,
			document: app.Feed{
				Title:    "test",
				Link:     "link",
				Language: "it",
				Image:    &app.Image{},
				Summary:  "Lorem Ipsum",
				Source:   "Unconditional Day",
				Date:     time.Now(),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := blevex.NewIndex("test.bleve", mapping.NewIndexMapping())
			if err != nil {
				t.Fatalf("expected bleve to be created")
			}

			f := bleve.NewFeedRepository(b)

			defer os.RemoveAll("test.bleve")

			err = f.Save(tc.document)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewBleve() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			found, err := f.Find(tc.document.Title)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewBleve() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if found == nil {
				t.Errorf("expected document to be found")
			}

			if found[0].Title != tc.document.Title {
				t.Errorf("expected document to be found")
			}
		})
	}
}
