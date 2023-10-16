package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/unconditionalday/server/internal/cmd/index"
	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/informer"
	"github.com/unconditionalday/server/internal/service"
	"github.com/unconditionalday/server/internal/version"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
)

var (
	ErrIndexNotProvided               = errors.New("index not provided, please provide it using --index flag")
	ErrAddressNotProvided             = errors.New("server address not provided, please provide it using --address flag")
	ErrPortNotProvided                = errors.New("server port not provided, please provide it using --port flag")
	ErrLogEnvNotProvided              = errors.New("server log-env not provided, please provide it using --log-env flag")
	ErrInformerScriptsPathNotProvided = errors.New("informer-scripts-path not provided, please provide it using --informer-scripts-path flag")
	ErrSourceRepositoryNotProvided    = errors.New("source repo not provided, please provide it using --source-repo flag")
	ErrSourceClientKeyNotProvided     = errors.New("source client-key not provided, please provide it using --source-client-key flag")
	ErrAllowedOriginsNotProvided      = errors.New("server allowed origins not provided, please provide it using --allowed-origins flag")
)

func NewServeCommand(version version.Build) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		Long:  `Starts the server`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			i := cobrax.Flag[string](cmd, "index").(string)
			if i == "" {
				return ErrIndexNotProvided
			}

			l := cobrax.Flag[string](cmd, "log-env").(string)
			if l == "" {
				return ErrLogEnvNotProvided
			}

			ip := cobrax.Flag[string](cmd, "informer-scripts-path").(string)
			if ip == "" {
				return ErrInformerScriptsPathNotProvided
			}

			s := cobrax.Flag[string](cmd, "source-repo").(string)
			if s == "" {
				return ErrSourceRepositoryNotProvided
			}

			sk := cobrax.Flag[string](cmd, "source-client-key").(string)
			if s == "" {
				return ErrSourceClientKeyNotProvided
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

			params := container.NewParameters(a, i, s, sk, ip, l, p, ao, version)
			c, _ := container.NewContainer(params)

			informer := informer.NewInformer(c.GetRunner(), ip, c.GetLogger())

			sourceService := service.NewSource(c.GetSourceClient(), c.GetParser(), c.GetVersioning(), c.GetLogger())

			source, err := sourceService.Fetch()
			if err != nil {
				return err
			}

			c.Parameters.SourceRelease = &source

			go index.UpdateResources(&source, sourceService, informer, c)

			return c.GetAPIServer().Start()
		},
	}

	cmd.Flags().StringP("address", "a", "localhost", "Server address")
	cmd.Flags().IntP("port", "p", 8080, "Server port")
	cmd.Flags().StringP("index", "s", "", "Index path")
	cmd.Flags().String("allowed-origins", "", "Allowed Origins")
	cmd.Flags().String("source-repo", "", "Source Repository")
	cmd.Flags().String("source-client-key", "", "Source Client Key")
	cmd.Flags().StringP("log-env", "l", "", "Log Env")
	cmd.Flags().StringP("informer-scripts-path", "i", "", "Informer scripts path")

	envPrefix := "UNCONDITIONAL_API"
	cobrax.BindFlags(cmd, cobrax.InitEnvs(envPrefix), envPrefix)

	return cmd
}
