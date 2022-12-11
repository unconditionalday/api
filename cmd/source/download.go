package source

import (
<<<<<<< HEAD
	"errors"
=======
	"github.com/unconditionalday/server/internal/cobrax"
	"github.com/unconditionalday/server/internal/iox"
	"github.com/unconditionalday/server/internal/netx"
	"github.com/unconditionalday/server/internal/service"
>>>>>>> ed637fb (chore: refactor)

	"github.com/spf13/cobra"

	"github.com/unconditionalday/server/internal/container"
	"github.com/unconditionalday/server/internal/service"
	cobrax "github.com/unconditionalday/server/internal/x/cobra"
	iox "github.com/unconditionalday/server/internal/x/io"
)

var (
	ErrSourceRepositoryNotProvided = errors.New("source repo not provided, please provide it using --source-repo flag")
	ErrSourceClientKeyNotProvided  = errors.New("source client-key not provided, please provide it using --source-client-key flag")
	ErrSourcePathNotProvided       = errors.New("source path not provided, please provide it using --path flag")
)

func NewDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download rss source",
		Long:  `Download rss source`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sp := cobrax.Flag[string](cmd, "path").(string)
			if sp == "" {
				return ErrSourcePathNotProvided
			}

			s := cobrax.Flag[string](cmd, "source-repo").(string)
			if s == "" {
				return ErrSourceRepositoryNotProvided
			}

			sk := cobrax.Flag[string](cmd, "source-client-key").(string)
			if s == "" {
				return ErrSourceClientKeyNotProvided
			}

			params := container.NewDefaultParameters()
			params.SourceRepository = s
			params.SourceClientKey = sk

			c, _ := container.NewContainer(params)

			sourceService := service.NewSource(
				c.GetSourceClient(),
				c.GetParser(),
				c.GetVersioning(),
				c.GetLogger())

			sourceRelease, err := sourceService.Fetch()
			if err != nil {
				return err
			}

			if err := iox.WriteJSON(sp, sourceRelease.Data); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP("source-repo", "s", "", "Source Repository URL")
	cmd.Flags().String("source-client-key", "", "Source Client Key")
	cmd.Flags().StringP("path", "p", "", "Source path")

	envPrefix := "UNCONDITIONAL_API"
	cobrax.BindFlags(cmd, cobrax.InitEnvs(envPrefix), envPrefix)

	return cmd
}
