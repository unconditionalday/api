package typesense

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func (f *FeedRepository) FindByKeyword(query string) ([]app.Feed, error) {
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
			ID:       doc["id"].(string),
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

func (f *FeedRepository) FindByID(id string) (app.Feed, error) {
	d, err := f.client.Collection("feeds").Document(id).Retrieve(f.ctx)
	if err != nil {
		return app.Feed{}, err
	}

	date, err := time.Parse(time.RFC3339, d["date"].(string))
	if err != nil {
		return app.Feed{}, err
	}

	fStruct := app.Feed{
		ID:       d["id"].(string),
		Title:    d["title"].(string),
		Link:     d["link"].(string),
		Source:   d["source"].(string),
		Language: d["language"].(string),
		Summary:  d["summary"].(string),
		Date:     date,
	}

	return fStruct, nil
}

func (f *FeedRepository) FindBySimilarity(feed app.Feed) ([]app.Feed, error) {
	vQ := fmt.Sprintf("title_summary_embedding:([], id: %s, distance_threshold:0.183)", feed.ID)
	eF := "title_summary_embedding"
	searchParameters := &api.SearchCollectionParams{
		Q:             feed.Title,
		QueryBy:       "title,title_summary_embedding",
		VectorQuery:   &vQ,
		ExcludeFields: &eF,
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
		"id":       generateUniqueID(doc.Link),
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
			"id":       generateUniqueID(doc.Link),
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

func generateUniqueID(link string) string {
	hash := sha256.New()
	hash.Write([]byte(link))
	hashBytes := hash.Sum(nil)

	uniqueID := hex.EncodeToString(hashBytes)

	return uniqueID
}
