package service

import (
	"encoding/json"

	"github.com/unconditionalday/server/internal/app"
	netx "github.com/unconditionalday/server/internal/x/net"
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
