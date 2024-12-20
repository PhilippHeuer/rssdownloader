package feed

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/PhilippHeuer/rssdownloader/pkg/config"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
)

const stateFile = "feed-state.json"

func DownloadFeed(feedConfig config.Feed) error {
	// state
	state, stateErr := config.LoadState(filepath.Join(feedConfig.OutputDir, stateFile))
	if stateErr != nil {
		log.Fatal().Err(stateErr).Msg("failed to load state")
	}

	// init state
	if _, ok := state.FeedState[feedConfig.Name]; !ok {
		state.FeedState[feedConfig.Name] = config.FeedState{FetchedAt: time.Now().UTC()}
	}
	stateItem := state.FeedState[feedConfig.Name]
	newestItem := &stateItem.FetchedAt

	// parse feed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedConfig.URL)
	if err != nil {
		return errors.Join(errors.New("failed to parse feed"), err)
	}
	log.Info().Str("feed", feedConfig.Name).Str("url", feedConfig.URL).Int("items", len(feed.Items)).Str("title", feed.Title).Msg("downloading feed")

	for _, item := range feed.Items {
		matchingRules := false
		matchingExclude := false

		for _, rule := range feedConfig.Rules {
			if rule.Type == "regex" {
				if matched, _ := regexp.MatchString(rule.Value, item.Title); matched {
					matchingRules = true
				}
			}
		}

		for _, exclude := range feedConfig.Exclude {
			if exclude.Type == "regex" {
				if matched, _ := regexp.MatchString(exclude.Value, item.Title); matched {
					matchingExclude = true
				}
			}
		}

		log.Trace().Bool("rules_match", matchingRules).Bool("exclude_match", matchingExclude).Str("item", item.Title).Time("published_at", *item.PublishedParsed).Time("state_last_download_at", state.FeedState[feedConfig.Name].FetchedAt).Msg("processing item")
		if matchingRules && !matchingExclude && item.PublishedParsed.After(state.FeedState[feedConfig.Name].FetchedAt) {
			fileName := renderFileNameTemplate(feedConfig.Template, item)
			log.Info().Str("item", item.Title).Str("filename", fileName).Msg("downloading item")
			if dlErr := downloadItem(item.Link, filepath.Join(feedConfig.OutputDir, fileName)); dlErr != nil {
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
	if err != nil {
		return err
	}

	return nil
}

func downloadItem(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func renderFileNameTemplate(template string, item *gofeed.Item) string {
	if template == "" {
		template = "{title}"
	}

	template = strings.ReplaceAll(template, "{title}", item.Title)
	if runtime.GOOS == "windows" {
		template = sanitizeFileNameForWindows(template)
	}
	return template
}

func sanitizeFileNameForWindows(filename string) string {
	invalidChars := []string{":", "\\", "/", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, char, "-")
	}
	return filename
}
