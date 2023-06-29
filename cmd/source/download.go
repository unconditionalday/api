package source

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/unconditionalday/server/internal/service"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
	iox "github.com/unconditionalday/server/internal/x/io"
	netx "github.com/unconditionalday/server/internal/x/net"
)

func NewDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download rss source",
		Long:  `Download rss source`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sp := cobrax.Flag[string](cmd, "path").(string)

			s := service.NewSource(netx.NewHttpClient())
			sd, err := s.Download("https://raw.githubusercontent.com/unconditionalday/source/main/source.json")
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
	viper.SetEnvPrefix("unconditional")

	cmd.Flags().StringP("path", "p", "", "Source path")

	return cmd
}
