package cmd

import (
	"time"

	"github.com/SlyMarbo/rss"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/luigibarbato/isolated-think-source/internal/app"
	"github.com/luigibarbato/isolated-think-source/internal/cobrax"
	"github.com/luigibarbato/isolated-think-source/internal/iox"
	blevex "github.com/luigibarbato/isolated-think-source/internal/repository/bleve"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "check if there are new rss source and update index",
		Long:  `check if there are new rss source and update index`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			i := cobrax.Flag[string](cmd, "name").(string)
			s := cobrax.Flag[string](cmd, "source").(string)

			var source app.Source
			source, err := iox.ReadJSON(s, source)
			if err != nil {
				return err
			}

			for {
				r, err := blevex.NewBleve(i)
				if err != nil {
					return err
				}
				feeds := make([]*rss.Feed, 0)
				for _, s := range source {
					feed, err := rss.Fetch(s.URL)
					if err != nil {
						logrus.Warnf("error fetching feed %s: %s", s.URL, err)
						continue
					}

					feeds = append(feeds, feed)
				}

				for _, feed := range feeds {
					for _, item := range feed.Items {
						f := app.Feed{
							Title:    item.Title,
							Link:     item.Link,
							Source:   feed.Title,
							Language: feed.Language,
							Image: &app.Image{
								Title: feed.Image.Title,
								URL:   feed.Image.URL,
							},
							Summary: strip.StripTags(item.Summary),
							Date:    item.Date,
						}

						logrus.Info("delete: ", f.Title)

						r.Delete(f)
						r.Index(f.Title, f)
						logrus.Info("feed indexed: ", f.Title)
						time.Sleep(1 * time.Minute)
					}
				}

			}
		},
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("IT")

	cmd.Flags().StringP("source", "s", "", "Source path")
	cmd.Flags().StringP("name", "n", "", "Index Name")

	return cmd
}
