package typesense

import (
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

func GetFeedSchema(client *typesense.Client) *api.CollectionSchema {
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
						ModelName: "ts/multilingual-e5-large",
					},
				},
			},
		},
		Name: "feeds",
	}

	return schema
}
