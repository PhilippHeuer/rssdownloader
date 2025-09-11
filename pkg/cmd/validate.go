package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/PhilippHeuer/rssdownloader/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("config", "c", "feeds.yaml", "config file")
	validateCmd.Flags().BoolP("quiet", "q", false, "suppress output, only exit code matters")
}

var validateCmd = &cobra.Command{
	Use:     "validate",
	Aliases: []string{"v"},
	Short:   "Validate the configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, _ := cmd.Flags().GetString("config")
		quiet, _ := cmd.Flags().GetBool("quiet")

		cfg, cfgErr := config.LoadConfig(configFile)
		if cfgErr != nil {
			return fmt.Errorf("failed to load configuration: %w", cfgErr)
		}

		hasErrors := false
		if !quiet {
			tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
			_, _ = fmt.Fprintln(tw, "FEED\tSTATUS\tDETAILS")

			for _, feed := range cfg.Feeds {
				errs := config.ValidateFeedConfig(feed)
				if len(errs) > 0 {
					hasErrors = true
					for _, e := range errs {
						_, _ = fmt.Fprintf(tw, "%s\tINVALID\t%s\n", feed.Name, e.Error())
					}
				} else {
					_, _ = fmt.Fprintf(tw, "%s\tOK\t-\n", feed.Name)
				}
			}

			tw.Flush()
		} else {
			// Quiet mode: still validate feeds, but don't print
			for _, feed := range cfg.Feeds {
				errs := config.ValidateFeedConfig(feed)
				if len(errs) > 0 {
					hasErrors = true
				}
			}
		}

		if hasErrors {
			return fmt.Errorf("validation failed")
		}
		return nil
	},
}
