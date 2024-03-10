package index

import (
	"errors"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/repository/typesense"
	"github.com/unconditionalday/server/internal/service"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
	typesensex "github.com/unconditionalday/server/internal/x/typesense"
)

var (
	ErrIndexNotProvided            = errors.New("index not provided, please provide it using --index flag")
	ErrFeedRepoHostNotProvided     = errors.New("feed repository host not provided, please provide it using --feed-repo-host flag")
	ErrFeedRepoKeyNotProvided      = errors.New("feed repository key not provided, please provide it using --feed-repo-key flag")
	ErrSourceNotProvided           = errors.New("source not provided, please provide it using --source flag")
	ErrSourceClientKeyNotProvided  = errors.New("source client-key not provided, please provide it using --source-client-key flag")
	ErrSourceRepositoryNotProvided = errors.New("source repo not provided, please provide it using --source-repo flag")
	ErrLogEnvNotProvided           = errors.New("log-env not provided, please provide it using --log-env flag")
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

			frh := cobrax.Flag[string](cmd, "feed-repo-host").(string)
			if i == "" {
				return ErrFeedRepoHostNotProvided
			}

			frk := cobrax.Flag[string](cmd, "feed-repo-key").(string)
			if frk == "" {
				return ErrFeedRepoKeyNotProvided
			}

			sr := cobrax.Flag[string](cmd, "source-repo").(string)
			if sr == "" {
				return ErrSourceRepositoryNotProvided
			}

			sk := cobrax.Flag[string](cmd, "source-client-key").(string)
			if sk == "" {
				return ErrSourceClientKeyNotProvided
			}

			l := cobrax.Flag[string](cmd, "log-env").(string)
			if l == "" {
				return ErrLogEnvNotProvided
			}

			params := container.NewDefaultParameters()
			params.FeedRepositoryIndex = i
			params.FeedRepositoryHost = frh
			params.FeedRepositoryKey = frk
			params.SourceClientKey = sk
			params.SourceRepository = sr
			params.LogEnv = l

			c, _ := container.NewContainer(params)

			t := c.GetTypesenseClient()

			feedSchema := typesense.GetFeedSchema(t)
			if err := typesensex.CreateOrUpdateCollection(t, feedSchema); err != nil {
				return err
			}

			sourceService := service.NewSource(c.GetSourceClient(), c.GetParser(), c.GetVersioning(), c.GetLogger())

			s, err := sourceService.Fetch()
			if err != nil {
				return err
			}

			feeds, err := sourceService.FetchFeeds(s.Data)
			if err != nil {
				c.GetLogger().Error("Can't fetch feeds", zap.Error(err))
			}

			if err := c.GetFeedRepository().Update(feeds...); err != nil {
				c.GetLogger().Error("Can't save feed", zap.Error(err))
			}

			c.GetLogger().Info("Index created", zap.String("Name", i))
			c.GetLogger().Info("Documents indexed", zap.Uint64("Count", c.GetFeedRepository().Count()))

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Index Name")
	cmd.Flags().StringP("log-env", "l", "", "Log Env")
	cmd.Flags().StringP("feed-repo-host", "", "", "Feed's repository host")
	cmd.Flags().StringP("feed-repo-key", "", "", "Feed's repository API's key")
	cmd.Flags().String("source-repo", "", "Source Repository")
	cmd.Flags().String("source-client-key", "", "Source Client Key")

	envPrefix := "UNCONDITIONAL_API"
	cobrax.BindFlags(cmd, cobrax.InitEnvs(envPrefix), envPrefix)

	return cmd
}
