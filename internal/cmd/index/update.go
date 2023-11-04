package index

import (
	"time"

	"go.uber.org/zap"

	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/service"
)

var (
	SourceUpdateInterval = 4 * time.Hour
	FeedsUpdateInterval  = 1 * time.Hour
)

func UpdateResources(source *app.SourceRelease, s *service.Source, c *container.Container) {
	srcReleasesChan := make(chan *app.SourceRelease)
	feedsTicker := time.NewTicker(FeedsUpdateInterval)

	go updateSource(srcReleasesChan, source, s, c.GetLogger())

	for {
		select {
		case newSource := <-srcReleasesChan:
			PopulateIndex(c, newSource.Data, s)
			c.GetLogger().Debug("Feeds updated")
		case <-feedsTicker.C:
			PopulateIndex(c, source.Data, s)
			c.GetLogger().Debug("Feeds updated")
		}
	}
}

func updateSource(sourceChan chan *app.SourceRelease, s *app.SourceRelease, sourceService *service.Source, l *zap.Logger) {
	for {
		time.Sleep(SourceUpdateInterval)

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
