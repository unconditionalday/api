package main

import (
	"github.com/sirupsen/logrus"

	cmd "github.com/unconditionalday/server/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		logrus.Fatal(err)
	}
}
