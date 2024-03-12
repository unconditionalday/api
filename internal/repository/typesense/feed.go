package typesense

import (
	"context"
	"time"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/unconditionalday/server/internal/app"
)

type FeedRepository struct {
	client *typesense.Client
	ctx    context.Context
}

func NewFeedRepository(client *typesense.Client) *FeedRepository {
	return &FeedRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (f *FeedRepository) Find(query string) ([]app.Feed, error) {
	searchParameters := &api.SearchCollectionParams{
		Q:       query,
		QueryBy: "title, summary",
	}
	searchResult, err := f.client.Collection("feeds").Documents().Search(f.ctx, searchParameters)
	if err != nil {
		return nil, err
	}

	feeds := make([]app.Feed, len(*searchResult.Hits))
	for i, x := range *searchResult.Hits {
		doc := *x.Document

		date, err := time.Parse(time.RFC3339, doc["date"].(string))
		if err != nil {
			return nil, err
		}

		f := app.Feed{
			Title:    doc["title"].(string),
			Link:     doc["link"].(string),
			Source:   doc["source"].(string),
			Language: doc["language"].(string),
			Summary:  doc["summary"].(string),
			Date:     date,
		}

		feeds[i] = f
	}

	return feeds, nil
}


func (f *FeedRepository) FindBySimilarity(query string) ([]app.Feed, error){
	searchParameters := &api.SearchCollectionParams{
		Q:       query,
		QueryBy: "title_summary_embedding",
	}
	searchResult, err := f.client.Collection("feeds").Documents().Search(f.ctx, searchParameters)
	if err != nil {
		return nil, err
	}

	feeds := make([]app.Feed, len(*searchResult.Hits))
	for i, x := range *searchResult.Hits {
		doc := *x.Document

		date, err := time.Parse(time.RFC3339, doc["date"].(string))
		if err != nil {
			return nil, err
		}

		f := app.Feed{
			Title:    doc["title"].(string),
			Link:     doc["link"].(string),
			Source:   doc["source"].(string),
			Language: doc["language"].(string),
			Summary:  doc["summary"].(string),
			Date:     date,
		}

		feeds[i] = f
	}

	return feeds, nil
}


func (f *FeedRepository) Save(doc app.Feed) error {
	docMap := map[string]interface{}{
		"id":       doc.Link,
		"title":    doc.Title,
		"link":     doc.Link,
		"source":   doc.Source,
		"language": doc.Language,
		"summary":  doc.Summary,
		"date":     doc.Date.Format(time.RFC3339),
	}

	// Perform the save/indexing operation
	_, err := f.client.Collection("feeds").Documents().Create(f.ctx, docMap)
	if err != nil {
		return err
	}

	return nil
}

func (f *FeedRepository) Update(docs ...app.Feed) error {
	docsMap := make([]interface{}, len(docs))

	for i, doc := range docs {
		// Convert app.Feed to map[string]interface{} for updating
		docMap := map[string]interface{}{
			"id":       doc.Link,
			"title":    doc.Title,
			"link":     doc.Link,
			"source":   doc.Source,
			"language": doc.Language,
			"summary":  doc.Summary,
			"date":     doc.Date.Format(time.RFC3339),
		}

		docsMap[i] = docMap
	}

	upsertAction := "upsert"
	params := &api.ImportDocumentsParams{
		Action: &upsertAction,
	}

	// Perform the update operation
	_, err := f.client.Collection("feeds").Documents().Import(f.ctx, docsMap, params)
	if err != nil {
		return err
	}

	return nil
}

func (f *FeedRepository) Count() uint64 {
	// Perform the operation to get document count
	coll, err := f.client.Collection("feeds").Retrieve(f.ctx)
	if err != nil || coll.NumDocuments == nil {
		return 0
	}

	return uint64(*coll.NumDocuments)
}

func (f *FeedRepository) Delete(doc app.Feed) error {
	if _, err := f.client.Collection("feeds").Document(doc.Link).Delete(f.ctx); err != nil {
		return err
	}

	return nil
}
