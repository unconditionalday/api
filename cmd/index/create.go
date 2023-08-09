package index

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/SlyMarbo/rss"
	"github.com/blevesearch/bleve/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/parser"
	blevex "github.com/unconditionalday/server/internal/repository/bleve"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
	iox "github.com/unconditionalday/server/internal/x/io"
)

var ErrSourceNotFound = errors.New("source not found, please download it first using source command")

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates the index",
		Long:  `Creates the index`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sp := cobrax.Flag[string](cmd, "source").(string)

			source, err := iox.ReadJSON(sp, app.Source{})
			if err != nil {
				return err
			}

			in := cobrax.Flag[string](cmd, "name").(string)
			b, err := blevex.NewIndex(in, bleve.NewIndexMapping())
			if err != nil {
				return err
			}

			feeds := make([]*rss.Feed, 0)
			for _, s := range source {
				feed, err := rss.Fetch(s.URL)
				if err != nil {
					// logrus.Warnf("error fetching feed %s: %s", s.URL, err)
					continue
				}

				feeds = append(feeds, feed)
			}

			p := parser.NewParser()

			var feedsItems []app.Feed
			for _, feed := range feeds {
				for _, item := range feed.Items {
					f := app.Feed{
						ID:       item.ID,
						Title:    p.Parse(item.Title),
						Link:     item.Link,
						Source:   feed.Title,
						Language: feed.Language,
						Image: &app.Image{
							Title: feed.Image.Title,
							URL:   feed.Image.URL,
						},
						Summary: p.Parse(item.Summary),
						Date:    item.Date,
					}

					if !f.IsValid() {
						logrus.Warn("invalid feed: ", f)
						continue
					}

					feedsItems = append(feedsItems, f)
				}
			}

			for i, f := range feedsItems {
				for j := i + 1; j < len(feedsItems); j++ {
					if err := getRelation(f, feedsItems[j]); err != nil {
						logrus.Warn("error getting relation: ", err)
					}
				}

				if err := b.Save(f); err != nil {
					logrus.Warn("error saving feed: ", err)
				}
			}

			logrus.Info("Index created: ", in)
			logrus.Info("Documents indexed: ", len(feeds))

			return nil
		},
	}

	cmd.Flags().StringP("source", "s", "", "Source path")
	cmd.Flags().StringP("name", "n", "", "Index Name")

	return cmd
}

func getRelation(source, target app.Feed) error {
	for _, r := range source.Related {
		if r == target.ID {
			return nil
		}
	}

	for _, r := range target.Related {
		if r == source.ID {
			return nil
		}
	}

	cmd := exec.Command("python3", "/Users/luigibarbato/Dev/Projects/unconditional/informer/relation.py", source.Title, target.Title)
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	outConv, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		return err
	}

	fmt.Println(source.Title, target.Title, outConv)

	if outConv > 0.8 {
		source.Related = append(source.Related, target.ID)
		target.Related = append(target.Related, source.ID)
	}

	return nil
}
