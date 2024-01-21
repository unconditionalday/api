package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/unconditionalday/server/internal/cmd/serve"
	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/service"
	"github.com/unconditionalday/server/internal/version"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
)

var (
	ErrAddressNotProvided          = errors.New("server address not provided, please provide it using --address flag")
	ErrPortNotProvided             = errors.New("server port not provided, please provide it using --port flag")
	ErrLogEnvNotProvided           = errors.New("server log-env not provided, please provide it using --log-env flag")
	ErrSourceRepositoryNotProvided = errors.New("source repo not provided, please provide it using --source-repo flag")
	ErrSourceClientKeyNotProvided  = errors.New("source client-key not provided, please provide it using --source-client-key flag")
	ErrAllowedOriginsNotProvided   = errors.New("server allowed origins not provided, please provide it using --allowed-origins flag")
	ErrDatabaseNameNotProvided     = errors.New("database name not provided, please provide it using --database-name flag")
	ErrDatabaseUserNotProvided     = errors.New("database user not provided, please provide it using --database-user flag")
	ErrDatabasePasswordNotProvided = errors.New("database password not provided, please provide it using --database-password flag")
)

func NewServeCommand(version version.Build) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		Long:  `Starts the server`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			l := cobrax.Flag[string](cmd, "log-env").(string)
			if l == "" {
				return ErrLogEnvNotProvided
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

			dbName := cobrax.Flag[string](cmd, "database-name").(string)
			if s == "" {
				return ErrDatabaseNameNotProvided
			}

			dbUser := cobrax.Flag[string](cmd, "database-user").(string)
			if s == "" {
				return ErrDatabaseUserNotProvided
			}

			dbPassword := cobrax.Flag[string](cmd, "database-password").(string)
			if s == "" {
				return ErrDatabasePasswordNotProvided
			}

			params := container.NewParameters(a, s, sk, dbName, dbUser, dbPassword, l, p, ao, version)
			c, _ := container.NewContainer(params)

			sourceService := service.NewSource(c.GetSourceClient(), c.GetParser(), c.GetVersioning(), c.GetLogger())

			source, err := sourceService.Fetch()
			if err != nil {
				return err
			}

			c.Parameters.SourceRelease = &source

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
	cmd.Flags().StringP("log-env", "l", "", "Log Env")
	cmd.Flags().StringP("database-name", "", "", "Database Name")
	cmd.Flags().StringP("database-user", "", "", "Database User")
	cmd.Flags().StringP("database-password", "", "", "Database Password")

	envPrefix := "UNCONDITIONAL_API"
	cobrax.BindFlags(cmd, cobrax.InitEnvs(envPrefix), envPrefix)

	return cmd
}
