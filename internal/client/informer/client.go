package informer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

func New(APIKey, baseURL string, http *http.Client) *Client {
	return &Client{
		apiKey:  APIKey,
		baseURL: baseURL,
		http:    http,
	}
}

type Options struct {
	WaitForModel bool `json:"wait_for_model"`
	UseCache     bool `json:"use_cache"`
}

type InputData struct {
	Inputs  []string `json:"inputs"`
	Options Options  `json:"options"`
}

func (c *Client) GetEmbeddings(sentence string) ([]float32, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/get_embedding", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("sentence", sentence)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var embeddings []float32
	if err := json.Unmarshal([]byte(body), &embeddings); err != nil {
		return nil, err
	}

	return embeddings, nil
}
