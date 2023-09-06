package informer

import (
	"errors"

	"github.com/unconditionalday/server/internal/informer/wiki"
)

var (
	ErrEmptyQuery    = errors.New("query string must not be empty")
	ErrEmptyLanguage = errors.New("language string must not be empty")
)

type Wiki struct {
	client *wiki.Client
}

func NewWiki() *Wiki {
	return &Wiki{
		client: wiki.New(),
	}
}

func (w *Wiki) Search(query string, lang string) (WikiInfo, error) {
	if query == "" {
		return WikiInfo{}, ErrEmptyQuery
	}

	if lang == "" {
		return WikiInfo{}, ErrEmptyLanguage
	}

	args := map[string]string{
		"action":   "query",
		"list":     "search",
		"srprop":   "",
		"srlimit":  "1",
		"srsearch": query,
	}

	res, err := w.client.DoRequest(args, lang)
	if err != nil {
		return WikiInfo{}, err
	}

	if len(res.Query.Search) == 0 {
		return WikiInfo{}, nil
	}

	title := res.Query.Search[0].Title

	wikiPage, err := wiki.MakeWikipediaPage(-1, title, "", false, w.client, lang)
	if wikiPage.Disambiguation != nil {
		title = wikiPage.Disambiguation[0]
		wikiPage, err = wiki.MakeWikipediaPage(-1, title, "", false, w.client, lang)
	}

	if err != nil {
		return WikiInfo{}, err
	}

	summary, err := wikiPage.GetSummary(w.client, lang)
	if err != nil {
		return WikiInfo{}, err
	}

	thumbnail, err := wikiPage.GetThumbURL(w.client, lang)
	if err != nil {
		return WikiInfo{}, err
	}

	return WikiInfo{
		Title:     wikiPage.Title,
		Language:  wikiPage.Language,
		Link:      wikiPage.URL,
		Summary:   summary,
		Thumbnail: thumbnail,
	}, nil
}

type WikiInfo struct {
	Title     string
	Link      string
	Summary   string
	Thumbnail string
	Language  string
}

func (r WikiInfo) IsValid() bool {
	return r.Title != "" && r.Link != "" && r.Summary != "" && r.Thumbnail != "" && r.Language != ""
}
