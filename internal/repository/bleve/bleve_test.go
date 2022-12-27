package bleve_test

import (
	"os"
	"testing"
	"time"

	"github.com/blevesearch/bleve/mapping"
	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/repository/bleve"
)

func TestBleveIndex(t *testing.T) {
	testCases := []struct {
		name    string
		b       *bleve.Bleve
		wantErr bool
	}{
		{
			name:    "bleve is created",
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			tc.b, err = bleve.NewBleveIndex("test.bleve", mapping.NewIndexMapping())

			defer os.RemoveAll("test.bleve")

			if tc.b == nil {
				t.Fatalf("expected bleve to be created")
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("NewBleve() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

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
			b, err := bleve.NewBleveIndex("test.bleve", mapping.NewIndexMapping())
			if b == nil {
				t.Fatalf("expected bleve to be created")
			}

			defer os.RemoveAll("test.bleve")

			err = b.Save(tc.document)

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
			b, err := bleve.NewBleveIndex("test.bleve", mapping.NewIndexMapping())
			if b == nil {
				t.Fatalf("expected bleve to be created")
			}

			defer os.RemoveAll("test.bleve")

			err = b.Save(tc.document)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewBleve() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			found, err := b.Find(tc.document.Title)
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

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name     string
		document app.Feed
		wantErr  bool
	}{
		{
			name:    "document is updated",
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
			b, err := bleve.NewBleveIndex("test.bleve", mapping.NewIndexMapping())
			if b == nil {
				t.Fatalf("expected bleve to be created")
			}

			defer os.RemoveAll("test.bleve")

			err = b.Save(tc.document)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewBleve() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			tc.document.Summary = "New Summary"

			err = b.Update(tc.document)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewBleve() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			found, err := b.Find(tc.document.Title)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewBleve() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if found == nil {
				t.Errorf("expected document to be found")
			}

			if found[0].Summary == tc.document.Summary {
				t.Errorf("expected document to be found")
			}

		})
	}
}
