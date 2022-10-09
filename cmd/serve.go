package cmd

import (
	"errors"

	"github.com/luigibarbato/isolated-think-source/internal/cobrax"
	blevex "github.com/luigibarbato/isolated-think-source/internal/repository/bleve"
	"github.com/luigibarbato/isolated-think-source/internal/webserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrIndexNotFound = errors.New("index not found, please create it first using source command")
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		Long:  `Starts the server`,
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

			return webserver.NewServer(sc, r).Start()
		},
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("IT")

	cmd.Flags().StringP("address", "a", "", "Server address")
	cmd.Flags().IntP("port", "p", 8080, "Server port")
	cmd.Flags().StringP("index", "s", "", "Index path")

	return cmd
}
