package wiki

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
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

// Create a new WikiClient
func New() *Client {
	return &Client{
		userAgent: "unconditional.day",
		URL:       "https://%v.wikipedia.org/w/api.php",
		lastCall:  time.Now(),
		cache:     MakeWikiCache(cacheExpiration, maxCacheMemory),
	}
}

/*
Make a request to the Wikipedia API using the given search parameters.

Returns a RequestResult (You can see the model in the models.go file)
*/
func (c *Client) DoRequest(args map[string]string, wikiLang string) (RequestResult, error) {
	const ReqPerSec = 199
	const ApiGap = time.Second / ReqPerSec

	url := fmt.Sprintf(c.URL, wikiLang)
	// Make new request object
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RequestResult{}, err
	}
	// Add header
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

	// Make GET request
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
	// Read body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return RequestResult{}, err
	}
	// Parse
	var result RequestResult
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return RequestResult{}, err
	}
	c.cache.Add(full_url, result)
	return result, nil
}

/*
Update the last time we call the API (API should)
*/
func (c *Client) updateLastCall(now time.Time) {
	c.lastCall = now
}
