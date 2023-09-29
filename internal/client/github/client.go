package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/unconditionalday/server/internal/app"
)

var (
	clientAPIVersion        = "X-GitHub-Api-Version"
	clientDefaultAPIVersion = "2022-11-28"
	clientAcceptType        = "application/vnd.github.v3+json"

	ErrResourceNotFound = errors.New("resource not found")
)

type Client struct {
	Repository string
	Owner      string
	Secret     string
	http       *http.Client
}

func New(repo, owner, secret string, http *http.Client) *Client {
	return &Client{
		Repository: repo,
		Owner:      owner,
		Secret:     secret,
		http:       http,
	}
}

func (c *Client) GetLatestVersion() (string, error) {
	q := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", c.Owner, c.Repository)

	req, err := c.prepareRequest(q)
	if err != nil {
		return "", err
	}

	var release Release
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&release); err != nil {
		return "", err
	}

	if release.TagName == "" {
		return "", ErrResourceNotFound
	}

	return release.TagName, nil
}

func (c *Client) Download(version string) (app.SourceRelease, error) {
	q := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", c.Owner, c.Repository, version)

	req, err := c.prepareRequest(q)
	if err != nil {
		return app.SourceRelease{}, err
	}

	var release Release
	resp, err := c.http.Do(req)
	if err != nil {
		return app.SourceRelease{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&release); err != nil {
		return app.SourceRelease{}, err
	}

	url := ""
	for _, a := range release.Assets {
		if a.Name == "source.json" {
			url = a.BrowserDownloadURL
		}
	}

	sourceJson, err := c.download(url)
	if err != nil {
		return app.SourceRelease{}, nil
	}

	return app.SourceRelease{
		Source:  sourceJson,
		Version: release.TagName,
	}, nil
}

func (c *Client) download(releaseUrl string) (app.Source, error) {
	var buf io.ReadWriter
	req, err := http.NewRequest("GET", releaseUrl, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Secret))

	resp, err := c.http.Do(req)
	if err != nil {
		return app.Source{}, err
	}

	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return app.Source{}, nil
	}

	var s app.Source
	if err := json.Unmarshal(res, &s); err != nil {
		return app.Source{}, err
	}

	return s, nil
}

func (c *Client) prepareRequest(url string) (*http.Request, error) {
	var buf io.ReadWriter
	req, err := http.NewRequest("GET", url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set(clientAPIVersion, clientDefaultAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", clientAcceptType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Secret))

	return req, nil
}

type Release struct {
	TagName    string         `json:"tag_name,omitempty"`
	Name       string         `json:"name,omitempty"`
	Draft      bool           `json:"draft,omitempty"`
	Prerelease bool           `json:"prerelease,omitempty"`
	Assets     []ReleaseAsset `json:"assets,omitempty"`
}

type ReleaseAsset struct {
	Name               string `json:"name,omitempty"`
	Label              string `json:"label,omitempty"`
	State              string `json:"state,omitempty"`
	ContentType        string `json:"content_type,omitempty"`
	Size               int    `json:"size,omitempty"`
	BrowserDownloadURL string `json:"browser_download_url,omitempty"`
}
