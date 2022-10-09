package cmd

import (
	"github.com/luigibarbato/isolated-think-source/cmd/index"
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
