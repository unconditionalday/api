package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	blevex "github.com/unconditionalday/server/internal/repository/bleve"
	"github.com/unconditionalday/server/internal/webserver"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
)

var (
	ErrIndexNotFound    = errors.New("index not found, please create it first using source command")
	ErrIndexNotProvided = errors.New("index not provided, please provide it using --index flag")
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		Long:  `Starts the server`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			i := cobrax.Flag[string](cmd, "index").(string)
			if i == "" {
				return ErrIndexNotProvided
			}

			r, err := blevex.New(i)
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
	viper.SetEnvPrefix("unconditional")

	cmd.Flags().StringP("address", "a", "localhost", "Server address")
	cmd.Flags().IntP("port", "p", 8080, "Server port")
	cmd.Flags().StringP("index", "s", "", "Index path")

	return cmd
}
