package feed

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/PhilippHeuer/rssdownloader/pkg/config"
	"github.com/rs/zerolog/log"
)

const stateFile = "feed-state.json"

func DownloadFeed(feedConfig config.Feed, timeout *int64, userAgent *string) error {
	if err := os.MkdirAll(feedConfig.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	state, stateErr := config.LoadState(filepath.Join(feedConfig.OutputDir, stateFile))
	if stateErr != nil {
		log.Fatal().Err(stateErr).Msg("failed to load state")
	}

	if _, ok := state.FeedState[feedConfig.Name]; !ok {
		state.FeedState[feedConfig.Name] = config.FeedState{FetchedAt: time.Now().UTC()}
	}
	stateItem := state.FeedState[feedConfig.Name]
	newestItem := &stateItem.FetchedAt

	filter, err := NewFilter(feedConfig.Rules, feedConfig.Exclude)
	if err != nil {
		return err
	}

	feed, err := ParseFeed(feedConfig.URL, timeout, userAgent)
	if err != nil {
		return fmt.Errorf("failed to parse feed: %w", err)
	}
	log.Info().Str("feed", feedConfig.Name).Str("url", feedConfig.URL).Int("items", len(feed.Items)).Str("title", feed.Title).Msg("downloading feed")

	for _, item := range feed.Items {
		log.Trace().
			Bool("matches", filter.Matches(item.Title)).
			Str("item", item.Title).
			Time("published_at", *item.PublishedParsed).
			Time("state_last_download_at", state.FeedState[feedConfig.Name].FetchedAt).
			Msg("processing item")

		if filter.Matches(item.Title) && item.PublishedParsed.After(state.FeedState[feedConfig.Name].FetchedAt) {
			filename := RenderTemplate(feedConfig.Template, item)
			log.Info().Str("item", item.Title).Str("filename", filename).Msg("downloading item")
			if dlErr := DownloadFile(item.Link, filepath.Join(feedConfig.OutputDir, filename), timeout); dlErr != nil {
				log.Error().Err(dlErr).Str("item", item.Title).Msg("failed to download item")
			}

			if item.PublishedParsed.After(*newestItem) {
				newestItem = item.PublishedParsed
			}
		}
	}

	stateItem.FetchedAt = *newestItem
	state.FeedState[feedConfig.Name] = stateItem
	err = config.SaveState(filepath.Join(feedConfig.OutputDir, stateFile), state)
	return err
}
