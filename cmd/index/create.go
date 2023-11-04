package index

import (
	"errors"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/unconditionalday/server/internal/cmd/index"
	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/service"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
)

var (
	ErrIndexNotProvided            = errors.New("index not provided, please provide it using --index flag")
	ErrSourceNotProvided           = errors.New("source not provided, please provide it using --source flag")
	ErrLogEnvNotProvided           = errors.New("log-env not provided, please provide it using --log-env flag")
	ErrFeedInformerKeyNotProvided  = errors.New("feed-informer-key not provided, please provide it using --feed-informer-key flag")
	ErrSourceRepositoryNotProvided = errors.New("source repo not provided, please provide it using --source-repo flag")
	ErrSourceClientKeyNotProvided  = errors.New("source client-key not provided, please provide it using --source-client-key flag")
)

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates the index",
		Long:  `Creates the index`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			i := cobrax.Flag[string](cmd, "name").(string)
			if i == "" {
				return ErrIndexNotProvided
			}

			fk := cobrax.Flag[string](cmd, "feed-informer-key").(string)
			if fk == "" {
				return ErrFeedInformerKeyNotProvided
			}

			s := cobrax.Flag[string](cmd, "source").(string)
			if s == "" {
				return ErrSourceNotProvided
			}

			sr := cobrax.Flag[string](cmd, "source-repo").(string)
			if s == "" {
				return ErrSourceRepositoryNotProvided
			}

			sk := cobrax.Flag[string](cmd, "source-client-key").(string)
			if s == "" {
				return ErrSourceClientKeyNotProvided
			}

			l := cobrax.Flag[string](cmd, "log-env").(string)
			if l == "" {
				return ErrLogEnvNotProvided
			}

			params := container.NewDefaultParameters()
			params.FeedIndex = i
			params.LogEnv = l
			params.SourceClientKey = sk
			params.SourceRepository = sr
			params.InformerClientKey = fk

			c, _ := container.NewContainer(params)

			sourceService := service.NewSource(c.GetSourceClient(), c.GetParser(), c.GetVersioning(), c.GetLogger())

			source, err := sourceService.Fetch()
			if err != nil {
				return err
			}

			index.PopulateIndex(c, source.Data, sourceService)

			c.GetLogger().Info("Index created", zap.String("Name", i))
			c.GetLogger().Info("Documents indexed", zap.Uint64("Count", c.GetFeedRepository().Count()))

			return nil
		},
	}

	cmd.Flags().StringP("source", "s", "", "Source Path")
	cmd.Flags().String("source-repo", "", "Source Repository")
	cmd.Flags().String("source-client-key", "", "Source Client Key")
	cmd.Flags().StringP("name", "n", "", "Index Name")
	cmd.Flags().StringP("log-env", "l", "", "Log Env")
	cmd.Flags().String("feed-informer-key", "", "Feed Informer Key")

	envPrefix := "UNCONDITIONAL_API"
	cobrax.BindFlags(cmd, cobrax.InitEnvs(envPrefix), envPrefix)

	return cmd
}
