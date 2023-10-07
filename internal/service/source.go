package service

import (
	"slices"

	"github.com/SlyMarbo/rss"
	"go.uber.org/zap"

	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/parser"
	"github.com/unconditionalday/server/internal/version"
)

type Source struct {
	client  app.SourceClient
	parser  *parser.Parser
	logger  *zap.Logger
	version version.Versioning
	current *app.SourceRelease
}

func NewSource(client app.SourceClient, parser *parser.Parser, versioning version.Versioning, logger *zap.Logger) *Source {
	return &Source{
		client:  client,
		parser:  parser,
		logger:  logger,
		version: versioning,
	}
}

func (s *Source) Fetch() (app.SourceRelease, error) {
	version, err := s.client.GetLatestVersion()
	if err != nil {
		return app.SourceRelease{}, err
	}

	if s.current != nil {
		res, err := s.version.Lower(s.current.Version, version)
		if !res && err == nil {
			return *s.current, nil
		}
	}

	source, err := s.client.Download(version)
	if err != nil {
		return app.SourceRelease{}, err
	}

	s.current = &app.SourceRelease{
		Data:         slices.Clone(source.Data),
		Version:      source.Version,
		LastUpdateAt: source.LastUpdateAt,
	}

	return *s.current, nil
}

func (s *Source) Update(src *app.SourceRelease) (bool, error) {
	release, err := s.Fetch()
	if err != nil {
		return false, err
	}

	res, err := s.version.Lower(src.Version, release.Version)
	if err != nil {
		return false, err
	}

	if !res {
		return false, nil
	}

	src.Version = release.Version
	src.Data = slices.Clone(release.Data)

	return true, nil
}

func (s *Source) FetchFeeds(src app.Source) ([]app.Feed, error) {
	rssFeeds := make([]*rss.Feed, 0)
	for _, source := range src {
		feed, err := rss.Fetch(source.URL)
		if err != nil {
			s.logger.Warn("error fetching feed", zap.String("Url", source.URL), zap.Error(err))
			continue
		}

		rssFeeds = append(rssFeeds, feed)
	}

	feeds := make([]app.Feed, 0)
	for _, rssFeed := range rssFeeds {
		for _, item := range rssFeed.Items {
			f := app.Feed{
				Title:    s.parser.Parse(item.Title),
				Link:     item.Link,
				Source:   rssFeed.Title,
				Language: rssFeed.Language,
				Image: &app.Image{
					Title: rssFeed.Image.Title,
					URL:   rssFeed.Image.URL,
				},
				Summary: s.parser.Parse(item.Summary),
				Date:    item.Date,
			}

			if !f.IsValid() {
				s.logger.Warn("invalid feed", zap.String("Feed title", f.Title))
				continue
			}

			feeds = append(feeds, f)
		}
	}

	return feeds, nil
}
