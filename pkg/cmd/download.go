package cmd

import (
	"github.com/PhilippHeuer/rssdownloader/pkg/config"
	"github.com/PhilippHeuer/rssdownloader/pkg/feed"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("config", "c", "feeds.yaml", "config file")
}

var downloadCmd = &cobra.Command{
	Use:     "download",
	Aliases: []string{"dl"},
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")
		cfg, cfgErr := config.LoadConfig(configFile)
		if cfgErr != nil {
			log.Fatal().Err(cfgErr).Msg("failed to load configuration")
		}

		for _, f := range cfg.Feeds {
			if f.Enabled != true {
				log.Debug().Str("feed", f.Name).Str("url", f.URL).Msg("skipping disabled feed")
				continue
			}

			_ = feed.DownloadFeed(f)
		}
	},
}
