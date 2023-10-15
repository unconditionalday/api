package index

import (
	"errors"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/service"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
	iox "github.com/unconditionalday/server/internal/x/io"
)

var (
	ErrIndexNotProvided           = errors.New("index not provided, please provide it using --index flag")
	ErrSourceNotProvided          = errors.New("source not provided, please provide it using --source flag")
	ErrSourceClientKeyNotProvided = errors.New("source client-key not provided, please provide it using --source-client-key flag")
	ErrLogEnvNotProvided          = errors.New("log-env not provided, please provide it using --log-env flag")
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

			s := cobrax.Flag[string](cmd, "source").(string)
			if s == "" {
				return ErrSourceNotProvided
			}

			l := cobrax.Flag[string](cmd, "log-env").(string)
			if l == "" {
				return ErrLogEnvNotProvided
			}

			params := container.NewDefaultParameters()
			params.FeedIndex = i
			params.LogEnv = l

			c, _ := container.NewContainer(params)

			sourceService := service.NewSource(c.GetSourceClient(), c.GetParser(), c.GetVersioning(), c.GetLogger())

			var source app.Source
			source, err := iox.ReadJSON(s, source)
			if err != nil {
				return err
			}

			feeds, err := sourceService.FetchFeeds(source)
			if err != nil {
				c.GetLogger().Error("Can't fetch feeds", zap.Error(err))
			}

			actualFeeds := c.GetFeedRepository().Count()
			for _, f := range feeds {
				err := c.GetFeedRepository().Save(f)
				if err != nil {
					c.GetLogger().Error("Can't save feed", zap.String("Feed", f.Link), zap.Error(err))
					continue
				}
			}

			c.GetLogger().Info("New Documents indexed", zap.Uint64("Count", c.GetFeedRepository().Count()-actualFeeds))

			return nil
		},
	}

	cmd.Flags().StringP("source", "s", "", "Source Path")
	cmd.Flags().StringP("name", "n", "", "Index Name")
	cmd.Flags().StringP("log-env", "l", "", "Log Env")

	envPrefix := "UNCONDITIONAL_API"
	cobrax.BindFlags(cmd, cobrax.InitEnvs(envPrefix), envPrefix)

	return cmd
}
