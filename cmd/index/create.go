package index

import (
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"sync"

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

			feedsItems := make([]app.Feed, 0)
			for _, feed := range feeds {
				for _, item := range feed.Items {
					f := app.Feed{
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

			// Definisci il numero massimo di goroutine che desideri eseguire contemporaneamente
			maxGoroutines := 10

			// Creare un canale per comunicare il completamento delle goroutine
			done := make(chan struct{})

			// Creare un semaforo per limitare il numero di goroutine in esecuzione contemporaneamente
			semaphore := make(chan struct{}, maxGoroutines)

			var wg sync.WaitGroup

			for f := range feedsItems {
				wg.Add(1)

				go func(f int) {
					defer wg.Done()

					// Acquisisci un semaforo per limitare il numero di goroutine in esecuzione contemporaneamente
					semaphore <- struct{}{}

					embeddings, err := getRelation(feedsItems[f])
					if err != nil {
						log.Println("Errore nell'ottenere la relazione:", err)
					} else {
						logrus.Info("Embeddings: ", embeddings)
						feedsItems[f].VectorEmbedding = embeddings
					}

					// Rilascia il semaforo una volta terminato
					<-semaphore
				}(f)
			}

			// Attendere il completamento di tutte le goroutine
			go func() {
				wg.Wait()
				close(done)
			}()

			// Attendere il completamento di tutte le goroutine prima di uscire
			<-done

			for _, f := range feedsItems {
				if err := b.Save(f); err != nil {
					return err
				}
			}

			logrus.Info("Index created: ", in)
			logrus.Info("Feeds: ", len(feeds))

			return nil
		},
	}

	cmd.Flags().StringP("source", "s", "", "Source path")
	cmd.Flags().StringP("name", "n", "", "Index Name")

	return cmd
}

func getRelation(source app.Feed) ([][]float64, error) {
	cmd := exec.Command("python3", "/Users/luigibarbato/Dev/Projects/unconditional/informer/main.py", source.Title)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var sentenceEmbeddings [][]float64

	if err := json.Unmarshal(out, &sentenceEmbeddings); err != nil {
		return nil, err
	}

	return sentenceEmbeddings, nil
}
