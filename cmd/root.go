package cmd

import (
	"github.com/spf13/cobra"
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
			Use:           "unconditional-api [command]",
			Short:         "Unconditional api engine that indexes and serves the content",
			SilenceUsage:  true,
			SilenceErrors: true,
		},
		config: cfg,
	}

	rootCmd.PersistentFlags().BoolVarP(&rootCmd.config.Debug, "debug", "D", false, "Enables debug output")

	rootCmd.AddCommand(NewServeCommand())
	rootCmd.AddCommand(NewIndexCmd())
	rootCmd.AddCommand(NewSourceCommand())

	return rootCmd
}
