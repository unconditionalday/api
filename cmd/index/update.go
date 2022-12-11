package index

import (
	"errors"

	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/cobrax"
	"github.com/unconditionalday/server/internal/iox"
	blevex "github.com/unconditionalday/server/internal/repository/bleve"

	"github.com/SlyMarbo/rss"
	"github.com/blevesearch/bleve"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	ErrSourceNotExists = errors.New("source not found, please download it first")
)

func NewUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates the index",
		Long:  `Updates the index`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sp := cobrax.Flag[string](cmd, "source").(string)

			var source app.Source
			source, err := iox.ReadJSON(sp, source)
			if err != nil {
				return err
			}

			in := cobrax.Flag[string](cmd, "name").(string)
			b, err := blevex.NewBleveIndex(in, bleve.NewIndexMapping())
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

					b.Index(f.Title, f)
				}
			}

			logrus.Infof("Index created: ", in)
			logrus.Infof("Documents indexed: ", len(feeds))

			return nil
		},
	}

	cmd.Flags().StringP("source", "s", "", "Source path")
	cmd.Flags().StringP("name", "n", "", "Index Name")

	return cmd
}
