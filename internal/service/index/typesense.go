package index

import (
	"context"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/unconditionalday/server/internal/app"
)

type FeedService struct {
	client *typesense.Client
	ctx    context.Context
}

func (f *FeedService) Create(doc app.Feed) error {
	schema := &api.CollectionSchema{
		Fields: []api.Field{
			{
				Name: "title",
				Type: "string",
			},
			{
				Name: "link",
				Type: "string",
			},
			{
				Name: "source",
				Type: "string",
			},
			{
				Name: "language",
				Type: "string",
			},
			{
				Name: "summary",
				Type: "string",
			},
			{
				Name: "title_summary_embedding",
				Type: "float[]",
				Embed: &struct {
					From        []string "json:\"from\""
					ModelConfig struct {
						AccessToken  *string "json:\"access_token,omitempty\""
						ApiKey       *string "json:\"api_key,omitempty\""
						ClientId     *string "json:\"client_id,omitempty\""
						ClientSecret *string "json:\"client_secret,omitempty\""
						ModelName    string  "json:\"model_name\""
						ProjectId    *string "json:\"project_id,omitempty\""
					} "json:\"model_config\""
				}{
					From: []string{"title", "summary"},
					ModelConfig: struct {
						AccessToken  *string "json:\"access_token,omitempty\""
						ApiKey       *string "json:\"api_key,omitempty\""
						ClientId     *string "json:\"client_id,omitempty\""
						ClientSecret *string "json:\"client_secret,omitempty\""
						ModelName    string  "json:\"model_name\""
						ProjectId    *string "json:\"project_id,omitempty\""
					}{
						ModelName: "ts/all-MiniLM-L12-v2",
					},
				},
			},
		},
		Name: "feeds",
	}

	if _, err := f.client.Collections().Create(context.Background(), schema); err != nil {
		return err
	}

	return nil
}

func (f *FeedService) Migrate(doc app.Feed) error {
	c := f.client.Collection("feeds")
	s := &api.CollectionUpdateSchema{
		Fields: []api.Field{
			{
				Name: "title",
				Type: "string",
			},
			{
				Name: "link",
				Type: "string",
			},
			{
				Name: "source",
				Type: "string",
			},
			{
				Name: "language",
				Type: "string",
			},
			{
				Name: "summary",
				Type: "string",
			},
			{
				Name: "title_summary_embedding",
				Type: "float[]",
				Embed: &struct {
					From        []string "json:\"from\""
					ModelConfig struct {
						AccessToken  *string "json:\"access_token,omitempty\""
						ApiKey       *string "json:\"api_key,omitempty\""
						ClientId     *string "json:\"client_id,omitempty\""
						ClientSecret *string "json:\"client_secret,omitempty\""
						ModelName    string  "json:\"model_name\""
						ProjectId    *string "json:\"project_id,omitempty\""
					} "json:\"model_config\""
				}{
					From: []string{"title", "summary"},
					ModelConfig: struct {
						AccessToken  *string "json:\"access_token,omitempty\""
						ApiKey       *string "json:\"api_key,omitempty\""
						ClientId     *string "json:\"client_id,omitempty\""
						ClientSecret *string "json:\"client_secret,omitempty\""
						ModelName    string  "json:\"model_name\""
						ProjectId    *string "json:\"project_id,omitempty\""
					}{
						ModelName: "ts/all-MiniLM-L12-v2",
					},
				},
			},
		},
	}

	_, err := c.Update(f.ctx, s)

	if err != nil {
		return err
	}

	return nil
}
