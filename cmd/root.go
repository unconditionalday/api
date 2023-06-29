// Copyright (c) 2017-present SIGHUP s.r.l All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/unconditionalday/server/internal/x/cobra"
)

type rootConfig struct {
	Debug bool
}

type RootCommand struct {
	*cobra.Command
	config *rootConfig
}

func NewRootCommand(versions map[string]string) *RootCommand {
	cfg := &rootConfig{}
	rootCmd := &RootCommand{
		Command: &cobra.Command{
			Use:           "unconditional-server [command]",
			Short:         "Unconditional's server engine that indexes and serves the content",
			SilenceUsage:  true,
			SilenceErrors: true,
			PersistentPreRun: func(cmd *cobra.Command, _ []string) {
				// Set log level
				if cobrax.Flag[bool](cmd, "debug").(bool) {
					logrus.SetLevel(logrus.DebugLevel)
				} else {
					logrus.SetLevel(logrus.InfoLevel)
				}
			},
		},
		config: cfg,
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("unconditional")

	rootCmd.PersistentFlags().BoolVarP(&rootCmd.config.Debug, "debug", "D", false, "Enables debug output")

	rootCmd.AddCommand(NewServeCommand())
	rootCmd.AddCommand(NewSourceCommand())
	rootCmd.AddCommand(NewIndexCmd())

	return rootCmd
}
