package github_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/labstack/gommon/log"

	"github.com/unconditionalday/server/internal/client/github"
)

func TestGetLatestVersion(t *testing.T) {
	key := os.Getenv("UNCONDITIONAL_API_SOURCE_CLIENT_KEY")
	c := github.New("source", "unconditionalday", key, http.DefaultClient)

	version, err := c.GetLatestVersion()
	if err != nil {
		t.Errorf("Fail: %v", err)
	}

	log.Info(version)
}
