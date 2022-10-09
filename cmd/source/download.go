package source

import (
	"github.com/luigibarbato/isolated-think-source/internal/cobrax"
	"github.com/luigibarbato/isolated-think-source/internal/iox"
	"github.com/luigibarbato/isolated-think-source/internal/netx"
	"github.com/luigibarbato/isolated-think-source/internal/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download rss source",
		Long:  `Download rss source`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sp := cobrax.Flag[string](cmd, "path").(string)

			s := service.NewSource(netx.NewFetcher())
			sd, err := s.Download("https://raw.githubusercontent.com/indipendent-thinker/indipendent-thinker-source/main/source.json")
			if err != nil {
				return err
			}

			if err := iox.WriteJSON(sp, sd); err != nil {
				return err
			}

			return nil
		},
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("IT")

	cmd.Flags().StringP("path", "p", "", "Source path")

	return cmd
}
