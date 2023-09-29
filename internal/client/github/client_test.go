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
	fakeClient := github.New("source", "unconditionalday", key, http.DefaultClient)

	version, err := fakeClient.GetLatestVersion()
	if err != nil {
		t.Errorf("Errore durante il recupero della versione più recente: %v", err)
	}

	log.Info(version)
}
