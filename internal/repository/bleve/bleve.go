package bleve

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"

	"github.com/unconditionalday/server/internal/app"
)

type Bleve struct {
	client bleve.Index
}

func NewBleve(path string) (*Bleve, error) {
	b, err := bleve.Open(path)
	if err != nil {
		return nil, err
	}

	return &Bleve{client: b}, nil
}

func NewBleveIndex(path string, mapping mapping.IndexMapping) (*Bleve, error) {
	b, err := bleve.New(path, mapping)
	if err != nil {
		return nil, err
	}

	return &Bleve{client: b}, nil
}

func (b *Bleve) Find(query string) ([]app.Feed, error) {
	q := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(q)
	// We need to say to bleve to return all fields of the document
	searchRequest.Fields = []string{"*"}
	searchResult, err := b.client.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var feeds []app.Feed
	for _, hit := range searchResult.Hits {
		date, err := time.Parse(time.RFC3339, hit.Fields["date"].(string))
		if err != nil {
			return nil, err
		}

		f := app.Feed{
			Title:    hit.Fields["title"].(string),
			Link:     hit.Fields["link"].(string),
			Source:   hit.Fields["source"].(string),
			Language: hit.Fields["language"].(string),
			Summary:  hit.Fields["summary"].(string),
			Date:     date,
		}

		if hit.Fields["image.url"] != "" {
			f.Image = &app.Image{
				Title: hit.Fields["image.title"].(string),
				URL:   hit.Fields["image.url"].(string),
			}
		}

		feeds = append(feeds, f)
	}

	return feeds, nil
}

func (b *Bleve) Save(doc app.Feed) error {
	if err := b.client.Index(doc.Link, doc); err != nil {
		return err
	}

	return nil
}

func (b *Bleve) Delete(doc app.Feed) error {
	if err := b.client.Delete(doc.Link); err != nil {
		return err
	}

	return nil
}

func (b *Bleve) Close() error {
	if err := b.client.Close(); err != nil {
		return err
	}

	return nil
}
