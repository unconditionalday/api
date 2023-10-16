package index

import (
	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/informer"
	"github.com/unconditionalday/server/internal/service"
	"go.uber.org/zap"
)

func PopulateIndex(c *container.Container, source app.Source, sourceService *service.Source, informer *informer.Informer) error {
	feeds, err := sourceService.FetchFeeds(source)
	if err != nil {
		c.GetLogger().Error("Can't fetch new feeds", zap.Error(err))
		return err
	}

	pool := app.NewWorkerPool(10)

	pool.Start()

	feedsToStore := len(feeds)
	for i, f := range feeds {
		taskID := i

		result, err := c.GetFeedRepository().Find(f.Link)
		if err != nil {
			c.GetLogger().Error("Unexpected error", zap.Error(err))
		}

		if len(result) > 0 {
			continue
		}

		pool.SubmitTask(func() {
			// s, err := informer.GetSimilarity(f.Summary)
			// if err != nil {
			// 	c.GetLogger().Error("Error during fetch similarity", zap.Error(err))
			// }

			// f.Similarity = s

			err = c.GetFeedRepository().Save(f)
			if err != nil {
				c.GetLogger().Error("Can't save feed", zap.String("Feed", f.Link), zap.Error(err))
			}

			feedsToStore--

			c.GetLogger().Debug("Task completed", zap.Int("ID", taskID))
			c.GetLogger().Debug("Start new task", zap.Int("feeds remained", feedsToStore))
		})
	}

	pool.Stop()

	return nil
}
