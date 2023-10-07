package main

import (
	"github.com/sirupsen/logrus"

	cmd "github.com/unconditionalday/server/cmd"
)

var (
	releaseVersion   string
	gitCommit string
)

func main() {
	v := map[string]string{
		"releaseVersion":   "unknown",
		"gitCommit": "unknown",
	}

	if releaseVersion != "" {
		v["releaseVersion"] = releaseVersion
	}

	if gitCommit != "" {
		v["gitCommit"] = gitCommit
	}

	if err := cmd.NewRootCommand(v).Execute(); err != nil {
		logrus.Fatal(err)
	}
}
