package cmd

import (
	"github.com/luigibarbato/isolated-think-source/cmd/source"
	"github.com/spf13/cobra"
)

func NewSourceCommand() *cobra.Command {
	sourceCmd := &cobra.Command{
		Use:   "source",
		Short: "Manage source",
	}

	sourceCmd.AddCommand(source.NewDownloadCmd())

	return sourceCmd
}
