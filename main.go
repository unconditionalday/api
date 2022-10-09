package main

import (
	"github.com/luigibarbato/isolated-think-source/cmd"
	"github.com/sirupsen/logrus"
)

var (
	version   = "unknown"
	gitCommit = "unknown"
	buildTime = "unknown"
	goVersion = "unknown"
	osArch    = "unknown"
)

func main() {
	versions := map[string]string{
		"version":   version,
		"gitCommit": gitCommit,
		"buildTime": buildTime,
		"goVersion": goVersion,
		"osArch":    osArch,
	}

	if err := cmd.NewRootCommand(versions).Execute(); err != nil {
		logrus.Fatal(err)
	}
}
