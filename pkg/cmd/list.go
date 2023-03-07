package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/PhilippHeuer/rssdownloader/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		cfg, cfgErr := config.LoadConfig("feeds.yaml")
		if cfgErr != nil {
			log.Fatal().Err(cfgErr).Msg("failed to load configuration")
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		_, _ = fmt.Fprintln(w, "NAME\tURI")
		for _, feed := range cfg.Feeds {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", feed.Name, feed.URL)
		}
		_ = w.Flush()
	},
}
