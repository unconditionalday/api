package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	blevex "github.com/unconditionalday/server/internal/repository/bleve"
	"github.com/unconditionalday/server/internal/webserver"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
	"go.uber.org/zap"
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

			r, err := blevex.NewBleve(i)
			if err != nil {
				return ErrIndexNotFound
			}

			a := cobrax.Flag[string](cmd, "address").(string)
			p := cobrax.Flag[int](cmd, "port").(int)
			ao := cobrax.FlagSlice(cmd, "allowed-origins")

			sc := webserver.Config{
				Port:           p,
				Address:        a,
				AllowedOrigins: ao,
			}

			l, _ := zap.NewProduction()

			return webserver.NewServer(sc, r, l).Start()
		},
	}

	cmd.Flags().StringP("address", "a", "localhost", "Server address")
	cmd.Flags().IntP("port", "p", 8080, "Server port")
	cmd.Flags().StringP("index", "s", "", "Index path")
	cmd.Flags().String("allowed-origins", "", "Allowed Origins")

	envPrefix := "UNCONDITIONAL_API"
	cobrax.BindFlags(cmd, cobrax.InitEnvs(envPrefix), envPrefix)

	return cmd
}
