package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unconditionalday/server/internal/version"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
)

var (
	releaseVersion = "unknown"
	gitCommit      = "unknown"
)

func NewRootCommand() *cobra.Command {
	v := version.Build{
		Version: releaseVersion,
		Commit:  gitCommit,
	}

	rootCmd := &cobra.Command{
		Use:           "unconditional-api [command]",
		Short:         "Unconditional api engine that indexes and serves the content",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(rootCmd *cobra.Command, _ []string) error {
			commitVersion := cobrax.Flag[string](rootCmd, "build-commit-version").(string)
			if commitVersion != "" {
				v.Commit = commitVersion
			}

			releaseVersion := cobrax.Flag[string](rootCmd, "build-release-version").(string)
			if releaseVersion != "" {
				v.Version = releaseVersion
			}

			return nil
		},
	}

	rootCmd.Flags().String("build-commit-version", "", "The Build GH commit version")
	rootCmd.Flags().String("build-release-version", "", "The Build release version")

	envPrefix := "UNCONDITIONAL_API"
	cobrax.BindFlags(rootCmd, cobrax.InitEnvs(envPrefix), envPrefix)

	rootCmd.AddCommand(NewServeCommand(v))
	rootCmd.AddCommand(NewIndexCmd())
	rootCmd.AddCommand(NewSourceCommand())

	return rootCmd
}
