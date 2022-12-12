package cmd

import (
	"github.com/spf13/cobra"

	"github.com/unconditionalday/server/cmd/index"
)

func NewIndexCmd() *cobra.Command {
	indexCmd := &cobra.Command{
		Use:   "index",
		Short: "Manage index",
	}

	indexCmd.AddCommand(index.NewCreateCommand())

	return indexCmd
}
