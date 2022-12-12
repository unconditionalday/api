package cmd

import (
	"github.com/spf13/cobra"

	"github.com/unconditionalday/server/cmd/source"
)

func NewSourceCommand() *cobra.Command {
	sourceCmd := &cobra.Command{
		Use:   "source",
		Short: "Manage source",
	}

	sourceCmd.AddCommand(source.NewDownloadCmd())

	return sourceCmd
}
