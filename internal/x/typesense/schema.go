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

func CreateSchema(client *typesense.Client, schema *api.CollectionSchema) error {
	if _, err := client.Collections().Create(context.Background(), schema); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return ErrCollectionAlreadyExists
		}
	}

	return nil
}
