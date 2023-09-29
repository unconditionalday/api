package serve

import (
	"time"

	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/service"
	"go.uber.org/zap"
)

func UpdateResources(source *app.SourceRelease, s *service.Source, c *container.Container) {
	srcReleasesChan := make(chan *app.SourceRelease)
	ticker := time.NewTicker(4 * time.Second)

	go updateSource(srcReleasesChan, source, s, c.GetLogger())

	for {
		select {
		case newSource := <-srcReleasesChan:
			feeds, err := s.FetchFeeds(newSource.Source)
			if err != nil {
				c.GetLogger().Error("Can't fetch new feeds", zap.Error(err))
			}

			for _, f := range feeds {
				if err := c.GetFeedRepository().Update(f); err != nil {
					c.GetLogger().Error("Can't index new feed", zap.Error(err))
				}
			}
			c.GetLogger().Debug("Feeds updated")
		case <-ticker.C:
			feeds, err := s.FetchFeeds(source.Source)
			if err != nil {
				c.GetLogger().Error("Can't fetch new feeds", zap.Error(err))
			}

			for _, f := range feeds {
				if err := c.GetFeedRepository().Update(f); err != nil {
					c.GetLogger().Error("Can't index new feed", zap.Error(err))
				}
			}

			c.GetLogger().Debug("Feeds updated")
		}
	}
}

func updateSource(sourceChan chan *app.SourceRelease, s *app.SourceRelease, sourceService *service.Source, l *zap.Logger) {
	for {

		time.Sleep(2 * time.Minute)

		currentVersion := s.Version
		res, err := sourceService.Update(s)
		if err != nil {
			l.Error("Can't comprare sources version", zap.Error(err))
			continue
		}

		if !res {
			l.Debug("No new version found, source not updated.", zap.String("Current version", currentVersion))
			continue
		}

		sourceChan <- s

		l.Debug("New version found, source updated.", zap.String("Old version", currentVersion), zap.String("New version", s.Version))
	}
}
