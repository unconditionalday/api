 
 
 

package cmd

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/unconditionalday/server/internal/cobrax"
)

type rootConfig struct {
	Spinner          *spinner.Spinner
	Debug            bool
	DisableAnalytics bool
	DisableTty       bool
}

type RootCommand struct {
	*cobra.Command
	config *rootConfig
}

func NewRootCommand(versions map[string]string) *RootCommand {
	cfg := &rootConfig{}
	rootCmd := &RootCommand{
		Command: &cobra.Command{
			Use:           "IT-source",
			Short:         "The individual and agnostic search engine",
			Long:          `Isolated-Thinker is a search engine that allows you to search for information on the web without being tracked by any company.`,
			SilenceUsage:  true,
			SilenceErrors: true,
			PersistentPreRun: func(cmd *cobra.Command, _ []string) {
				// Configure the spinner
				w := logrus.StandardLogger().Out
				if cobrax.Flag[bool](cmd, "no-tty").(bool) {
					f := new(logrus.TextFormatter)
					f.DisableColors = false
					logrus.SetFormatter(f)
				}
				cfg.Spinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(w))

				// Set log level
				if cobrax.Flag[bool](cmd, "debug").(bool) {
					logrus.SetLevel(logrus.DebugLevel)
				} else {
					logrus.SetLevel(logrus.InfoLevel)
				}
			},
		},
		config: cfg,
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("unconditional")

	rootCmd.PersistentFlags().BoolVarP(&rootCmd.config.Debug, "debug", "D", false, "Enables Isolated-Thinker debug output")
	rootCmd.PersistentFlags().BoolVarP(&rootCmd.config.DisableTty, "no-tty", "T", false, "Disable TTY")

	rootCmd.AddCommand(NewServeCommand())
	rootCmd.AddCommand(NewSourceCommand())
	rootCmd.AddCommand(NewIndexCmd())

	return rootCmd
}
