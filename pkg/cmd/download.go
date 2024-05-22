package cmd

import (
	"errors"

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

		// config
		cfg, cfgErr := config.LoadConfig(configFile)
		if cfgErr != nil {
			log.Fatal().Err(cfgErr).Msg("failed to load configuration")
		}

		for i, f := range cfg.Feeds {
			// skip disabled feeds
			if f.Enabled != true {
				log.Debug().Str("feed", f.Name).Str("url", f.URL).Msg("skipping disabled feed")
				continue
			}

			err := processFeed(f)
			if err != nil {
				log.Error().Err(err).Int("index", i).Str("feed", f.Name).Str("url", f.URL).Msg("failed to query feed")
			}
		}
	},
}

func processFeed(f config.Feed) error {
	// skip invalid feeds
	errs := config.ValidateFeedConfig(f)
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	err := feed.DownloadFeed(f)
	if err != nil {
		return err
	}

	return nil
}
