package bleve

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

type Bleve struct{}

func New(path string) (bleve.Index, error) {
	b, err := bleve.Open(path)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func NewIndex(path string, mapping mapping.IndexMapping) (bleve.Index, error) {
	b, err := bleve.New(path, mapping)
	if err != nil {
		return nil, err
	}

	return b, nil
}
