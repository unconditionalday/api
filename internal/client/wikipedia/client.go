package wikipedia

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/unconditionalday/server/internal/search"
)

type Client struct {
	userAgent string
	URL       string
	lastCall  time.Time
	cache     *Cache
}

const (
	cacheExpiration = 12 * time.Hour
	maxCacheMemory  = 500
)

var (
	ErrEmptyQuery           = errors.New("query string must not be empty")
	ErrEmptyLanguage        = errors.New("language string must not be empty")
	ErrDisambiguationResult = errors.New("disambiguation result")
)

func NewClient() *Client {
	return &Client{
		userAgent: "unconditional.day",
		URL:       "https://%v.wikipedia.org/w/api.php",
		lastCall:  time.Now(),
		cache:     MakeWikiCache(cacheExpiration, maxCacheMemory),
	}
}

func (c *Client) doRequest(args map[string]string, wikiLang string) (RequestResult, error) {
	const ReqPerSec = 199
	const ApiGap = time.Second / ReqPerSec

	url := fmt.Sprintf(c.URL, wikiLang)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RequestResult{}, err
	}
	request.Header.Set("User-Agent", c.userAgent)
	q := request.URL.Query()
	// Add parameters
	if args["format"] == "" {
		args["format"] = "json"
	}
	if args["action"] == "" {
		args["action"] = "query"
	}
	for k, v := range args {
		q.Add(k, v)
	}
	request.URL.RawQuery = q.Encode()
	now := time.Now()
	if now.Sub(c.lastCall) < ApiGap {
		wait := c.lastCall.Add(ApiGap).Sub(now)
		time.Sleep(wait)
		now = time.Now()
	}
	// Check in cache
	full_url := request.URL.String()
	r, err := c.cache.Get(full_url)
	if err == nil {
		return r, nil
	}

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(request)
	defer c.updateLastCall(now)
	if err != nil {
		return RequestResult{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return RequestResult{}, errors.New("unable to fetch the results")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return RequestResult{}, err
	}
	var result RequestResult
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return RequestResult{}, err
	}
	c.cache.Add(full_url, result)
	return result, nil
}

func (c *Client) updateLastCall(now time.Time) {
	c.lastCall = now
}

func (w *Client) FetchContextDetails(query string, lang string) (search.ContextDetails, error) {
	if query == "" {
		return search.ContextDetails{}, ErrEmptyQuery
	}

	if lang == "" {
		return search.ContextDetails{}, ErrEmptyLanguage
	}

	args := map[string]string{
		"action":   "query",
		"list":     "search",
		"srprop":   "",
		"srlimit":  "1",
		"srsearch": query,
	}

	res, err := w.doRequest(args, lang)
	if err != nil {
		return search.ContextDetails{}, err
	}

	if len(res.Query.Search) == 0 {
		return search.ContextDetails{}, nil
	}

	title := res.Query.Search[0].Title

	wikiPage, err := MakeWikipediaPage(-1, title, "", false, w, lang)
	if len(wikiPage.Disambiguation) != 0 {
		return search.ContextDetails{}, nil
	}

	if err != nil {
		return search.ContextDetails{}, err
	}

	summary, err := wikiPage.GetSummary(w, lang)
	if err != nil {
		return search.ContextDetails{}, err
	}

	thumbnail, err := wikiPage.GetThumbURL(w, lang)
	if err != nil {
		return search.ContextDetails{}, err
	}

	s := search.ContextDetails{
		Title:     wikiPage.Title,
		Language:  wikiPage.Language,
		Link:      wikiPage.URL,
		Summary:   summary,
		Thumbnail: thumbnail,
	}

	if !s.IsValid() {
		return search.ContextDetails{}, nil
	}

	return s, nil
}
