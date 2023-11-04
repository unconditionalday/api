package bleve

import (
	"time"

	"github.com/blevesearch/bleve/v2"

	"github.com/unconditionalday/server/internal/app"
)

type FeedRepository struct {
	client bleve.Index
}

func NewFeedRepository(client bleve.Index) *FeedRepository {
	return &FeedRepository{
		client: client,
	}
}

func (f *FeedRepository) Find(query string) ([]app.Feed, error) {
	q := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(q)
	// We need to say to bleve to return all fields of the document
	searchRequest.Fields = []string{"*"}
	searchResult, err := f.client.Search(searchRequest)
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

func (f *FeedRepository) Exists(id string) (bool, error) {
	feed, err := f.client.Document(id)
	if err != nil {
		return false, err
	}

	return feed == nil, nil
}

func (f *FeedRepository) Save(doc app.Feed) error {
	if err := f.client.Index(doc.Link, doc); err != nil {
		return err
	}

	return nil
}

func (b *FeedRepository) Update(doc app.Feed) error {
	return b.Save(doc)
}

func (b *FeedRepository) Count() uint64 {
	c, _ := b.client.DocCount()

	return c
}

func (f *FeedRepository) Delete(doc app.Feed) error {
	if err := f.client.Delete(doc.Link); err != nil {
		return err
	}

	return nil
}
