package typesense

import (
	"context"
	"errors"
	"strings"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

var (
	ErrCollectionAlreadyExists = errors.New("the collection already exists")
)

func CreateOrUpdateCollection(client *typesense.Client, schema *api.CollectionSchema) error {
	if _, err := client.Collections().Create(context.Background(), schema); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return updateCollection(client, schema)
		}
	}

	return nil
}

func updateCollection(client *typesense.Client, schema *api.CollectionSchema) error {
	u := &api.CollectionUpdateSchema{
		Fields: schema.Fields,
	}

	if _, err := client.Collection(schema.Name).Update(context.Background(), u); err != nil{
		return err
	}

	return nil
}
