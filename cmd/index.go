package cmd

import (
	"github.com/unconditionalday/server/cmd/index"

	"github.com/spf13/cobra"
)

func NewIndexCmd() *cobra.Command {
	indexCmd := &cobra.Command{
		Use:   "index",
		Short: "Manage index",
	}

	indexCmd.AddCommand(index.NewCreateCommand())

	return indexCmd
}
