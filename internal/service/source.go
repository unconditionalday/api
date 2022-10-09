package service

import (
	"encoding/json"

	"github.com/luigibarbato/isolated-think-source/internal/app"
	"github.com/luigibarbato/isolated-think-source/internal/netx"
)

type Source struct {
	client netx.Client
}

func NewSource(client netx.Client) *Source {
	return &Source{client: client}
}

func (s *Source) Download(path string) (app.Source, error) {
	data, err := s.client.Download(path)
	if err != nil {
		return app.Source{}, err
	}

	var source app.Source
	if err := json.Unmarshal(data, &source); err != nil {
		return app.Source{}, err
	}

	return source, nil
}
