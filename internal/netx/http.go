package netx

import (
	"io/ioutil"
	"net/http"
)

type Fetcher struct{}

func NewFetcher() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) Download(src string) ([]byte, error) {
	resp, err := http.Get(src)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
