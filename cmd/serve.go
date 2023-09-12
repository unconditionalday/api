package cmd

import (
	"errors"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/parser"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
	"go.uber.org/zap"
)

var (
	ErrIndexNotFound             = errors.New("index not found, please create it first using source command")
	ErrIndexNotProvided          = errors.New("index not provided, please provide it using --index flag")
	ErrAddressNotProvided        = errors.New("server address not provided, please provide it using --address flag")
	ErrPortNotProvided           = errors.New("server port not provided, please provide it using --port flag")
	ErrAllowedOriginsNotProvided = errors.New("server allowed origins not provided, please provide it using --allowed-origins flag")
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		Long:  `Starts the server`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			i := cobrax.Flag[string](cmd, "index").(string)
			if i == "" {
				return ErrIndexNotProvided
			}

			a := cobrax.Flag[string](cmd, "address").(string)
			if a == "" {
				return ErrAddressNotProvided
			}

			p := cobrax.Flag[int](cmd, "port").(int)
			if p == 0 {
				return ErrPortNotProvided
			}

			ao := cobrax.FlagSlice(cmd, "allowed-origins")
			if ao == nil {
				return ErrAllowedOriginsNotProvided
			}

			params := container.NewParameters(a, i, p, ao)
			c, _ := container.NewContainer(params)

			sourceChan := make(chan app.Source)

			go func(sourceChan chan app.Source, s app.SourceService, l *zap.Logger) {
				for {
					updateSource(sourceChan, c.GetSourceService(), l)
					time.Sleep(2 * time.Minute)
				}
			}(sourceChan, c.GetSourceService(), c.GetLogger())

			go func(sourceChan chan app.Source, index app.FeedRepository, l *zap.Logger, pa *parser.Parser) {
				for {
					select {
					case s := <-sourceChan:
						feeds := fetchNewFeeds(s, l, pa)
						updateIndex(index, feeds, l)
					}
				}

			}(sourceChan, c.GetFeedRepository(), c.GetLogger(), c.GetParser())

			return c.GetAPIServer().Start()
		},
	}

	cmd.Flags().StringP("address", "a", "localhost", "Server address")
	cmd.Flags().IntP("port", "p", 8080, "Server port")
	cmd.Flags().StringP("index", "s", "", "Index path")
	cmd.Flags().String("allowed-origins", "", "Allowed Origins")

	envPrefix := "UNCONDITIONAL_API"
	cobrax.BindFlags(cmd, cobrax.InitEnvs(envPrefix), envPrefix)

	return cmd
}

func updateSource(sourceChan chan app.Source, service app.SourceService, l *zap.Logger) {
	s, err := service.Download("https://raw.githubusercontent.com/unconditionalday/source/main/source.json")
	if err != nil {
		return
	}

	sourceChan <- s

	l.Info("Update Source")
}

func updateIndex(index app.FeedRepository, feeds []app.Feed, l *zap.Logger) {
	for _, f := range feeds {
		index.Update(f)
	}

	l.Info("Update Index")
}

func fetchNewFeeds(source app.Source, l *zap.Logger, parser *parser.Parser) []app.Feed {
	feeds := make([]*rss.Feed, 0)
	for _, s := range source {
		feed, err := rss.Fetch(s.URL)
		if err != nil {
			logrus.Warnf("error fetching feed %s: %s", s.URL, err)
			continue
		}

		feeds = append(feeds, feed)
	}

	items := make([]app.Feed, 0)
	for _, feed := range feeds {
		for _, item := range feed.Items {
			f := app.Feed{
				Title:    parser.Parse(item.Title),
				Link:     item.Link,
				Source:   feed.Title,
				Language: feed.Language,
				Image: &app.Image{
					Title: feed.Image.Title,
					URL:   feed.Image.URL,
				},
				Summary: parser.Parse(item.Summary),
				Date:    item.Date,
			}

			items = append(items, f)
		}
	}

	return items
}
