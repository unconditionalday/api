package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/unconditionalday/server/internal/cmd/serve"
	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/service"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
)

var (
	ErrIndexNotProvided            = errors.New("index not provided, please provide it using --index flag")
	ErrAddressNotProvided          = errors.New("server address not provided, please provide it using --address flag")
	ErrPortNotProvided             = errors.New("server port not provided, please provide it using --port flag")
	ErrSourceRepositoryNotProvided = errors.New("source repo not provided, please provide it using --source-repo flag")
	ErrAllowedOriginsNotProvided   = errors.New("server allowed origins not provided, please provide it using --allowed-origins flag")
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

			s := cobrax.Flag[string](cmd, "source-repo").(string)
			if s == "" {
				return ErrSourceRepositoryNotProvided
			}

			sk := cobrax.Flag[string](cmd, "source-client-key").(string)
			if s == "" {
				return ErrSourceRepositoryNotProvided
			}

			a := cobrax.Flag[string](cmd, "address").(string)
			if a == "" {
				return ErrAddressNotProvided
			}

			p := cobrax.Flag[int](cmd, "port").(int)
			if p == 0 {
				return ErrPortNotProvided
			}

			ao := cobrax.FlagSlice(cmd, "allowed-origins")
			if ao == nil {
				return ErrAllowedOriginsNotProvided
			}

			params := container.NewParameters(a, i, s, sk, p, ao,)
			c, _ := container.NewContainer(params)

			sourceService := service.NewSource(c.GetSourceClient(), c.GetParser(), c.GetVersioning(), c.GetLogger())

			source, err := sourceService.Download()
			if err != nil {
				return err
			}

			go serve.UpdateResources(&source, sourceService, c)

			return c.GetAPIServer().Start()
		},
	}

	cmd.Flags().StringP("address", "a", "localhost", "Server address")
	cmd.Flags().IntP("port", "p", 8080, "Server port")
	cmd.Flags().StringP("index", "s", "", "Index path")
	cmd.Flags().String("allowed-origins", "", "Allowed Origins")
	cmd.Flags().String("source-repo", "", "Source Repository")
	cmd.Flags().String("source-client-key", "", "Source Client Key")
	
	envPrefix := "UNCONDITIONAL_API"
	cobrax.BindFlags(cmd, cobrax.InitEnvs(envPrefix), envPrefix)

	return cmd
}
