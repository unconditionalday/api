package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unconditionalday/server/internal/version"
)

func NewRootCommand(versions map[string]string) *cobra.Command {
	v := version.Build{
		Version: versions["releaseVersion"],
		Commit:  versions["gitCommit"],
	}

	rootCmd := &cobra.Command{
		Use:           "unconditional-api [command]",
		Short:         "Unconditional api engine that indexes and serves the content",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(NewServeCommand(v))
	rootCmd.AddCommand(NewIndexCmd())
	rootCmd.AddCommand(NewSourceCommand())

	return rootCmd
}
