package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/unconditionalday/server/internal/cobrax"
	blevex "github.com/unconditionalday/server/internal/repository/bleve"
	"github.com/unconditionalday/server/internal/webserver"
)

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "run all services",
		Long:  `run all services`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			i := cobrax.Flag[string](cmd, "index").(string)

			r, err := blevex.NewBleve(i)
			if err != nil {
				return ErrIndexNotFound
			}

			a := cobrax.Flag[string](cmd, "address").(string)
			p := cobrax.Flag[int](cmd, "port").(int)

			sc := webserver.ServerConfig{
				Address: a,
				Port:    p,
			}

			if err := webserver.NewServer(sc, r).Start(); err != nil {
				return err
			}

			for {
				select {
				case <-cmd.Context().Done():
					return nil
				}
			}
		},
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("unconditional")

	cmd.Flags().StringP("address", "a", "", "Server address")
	cmd.Flags().IntP("port", "p", 8080, "Server port")
	cmd.Flags().StringP("index", "s", "", "Index path")

	return cmd
}
