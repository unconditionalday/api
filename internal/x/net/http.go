package netx

import (
	"io/ioutil"
	"net/http"
)

type HttpClient struct{}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

func (h *HttpClient) Download(src string) ([]byte, error) {
	resp, err := http.Get(src)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func HelpAddURL(s string) string {
	if s[0:4] == "http" {
		return s
	}
	return "http:" + s
}
