package informer

import (
	"errors"

	"github.com/unconditionalday/server/internal/informer/wiki"
)

type Wiki struct {
	client wiki.Client
}

func NewWiki() *Wiki {
	return &Wiki{
		client: wiki.New(),
	}
}

func (w *Wiki) Search(query string, lang string) (Result, error) {
	if query == "" {
		return Result{}, errors.New("query string must not be empty")
	}

	if lang == "" {
		return Result{}, errors.New("language string must not be empty")
	}

	args := map[string]string{
		"action":   "query",
		"list":     "search",
		"srprop":   "",
		"srlimit":  "1",
		"srsearch": query,
	}

	res, err := w.client.RequestWikiApi(args, lang)
	if err != nil {
		return Result{}, err
	}

	title := res.Query.Search[0].Title

	wikiPage, err := wiki.MakeWikipediaPage(-1, title, "", false, w.client, lang)
	if err != nil {
		if err.Error() == "disambiguation" {
			return Result{}, errors.New("ambiguous result")
		}
	}

	summary, err := wikiPage.GetSummary(w.client, lang)
	if err != nil {
		return Result{}, err
	}

	thumbnail, err := wikiPage.GetThumbURL(w.client, lang)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Title:     wikiPage.Title,
		Language:  wikiPage.Language,
		Link:      wikiPage.URL,
		Summary:   summary,
		Thumbnail: thumbnail,
	}, nil
}

type Result struct {
	Title     string
	Link      string
	Summary   string
	Thumbnail string
	Language  string
}

func (r Result) IsValid() bool {
	return r.Title != "" && r.Link != "" && r.Summary != "" && r.Thumbnail != "" && r.Language != ""
}
