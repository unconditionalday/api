package index

import (
	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/service"
	"go.uber.org/zap"
)

func PopulateIndex(c *container.Container, source app.Source, sourceService *service.Source) error {
	feeds, err := sourceService.FetchFeeds(source)
	if err != nil {
		c.GetLogger().Error("Can't fetch new feeds", zap.Error(err))
		return err
	}

	feedsToStore := getFeedsToStore(feeds, c)

	for _, f := range feedsToStore {
		embeddings, err := c.GetInformerClient().GetEmbeddings(f.Summary)
		if err != nil {
			c.GetLogger().Error("Can't store new feeds", zap.String("Feed", f.Link), zap.Error(err))
			return err
		}

		f.Similarity = embeddings

		err = c.GetFeedRepository().Save(f)
		if err != nil {
			c.GetLogger().Error("Can't store new feeds", zap.String("Feed", f.Link), zap.Error(err))
		}
	}

	return nil
}

func getFeedsToStore(feeds []app.Feed, c *container.Container) []app.Feed {
	feedsToStore := make([]app.Feed, 0)
	for _, f := range feeds {
		exists, err := c.GetFeedRepository().Exists(f.Link)
		if err != nil {
			c.GetLogger().Error("Unexpected error", zap.Error(err))
		}

		if !exists {
			continue
		}

		feedsToStore = append(feedsToStore, f)
	}

	return feedsToStore
}
