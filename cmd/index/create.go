package index

import (
	"github.com/spf13/cobra"
	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/service"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
	"go.uber.org/zap"
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

			s := cobrax.Flag[int](cmd, "source-repo").(string)
			if s == "" {
				return ErrSourceRepositoryNotProvided
			}

			params := container.NewDefaultParameters()
			params.FeedIndex = i
			params.SourceRepository = s

			c, _ := container.NewContainer(params)

			sourceService := service.NewSource(c.GetSourceClient(), c.GetParser(), c.GetVersioning(), c.GetLogger())

			source, err := sourceService.Download()
			if err != nil {
				return err
			}

			feeds, err := sourceService.FetchFeeds(source.Source)
			if err != nil {
				c.GetLogger().Error("Can't fetch feeds", zap.Error(err))
			}

			for _, f := range feeds {
				err := c.GetFeedRepository().Save(f)
				if err != nil {
					c.GetLogger().Error("Can't save feed", zap.String("Feed", f.Link), zap.Error(err))
					continue
				}
			}

			c.GetLogger().Info("Index created: ", zap.String("Name", i))
			c.GetLogger().Info("Documents indexed: ", zap.Uint64("Count", c.GetFeedRepository().Count()))

			return nil
		},
	}

	cmd.Flags().StringP("source-repo", "s", "", "Source Repository URL")
	cmd.Flags().StringP("name", "n", "", "Index Name")

	return cmd
}
