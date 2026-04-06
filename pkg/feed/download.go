package feed

import (
	"errors"
	"fmt"
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

const (
	stateFile      = "feed-state.json"
	defaultTimeout = 30 * time.Second
	maxFileSize    = 500 * 1024 * 1024
)

var httpClient = &http.Client{
	Timeout: defaultTimeout,
}

func DownloadFeed(feedConfig config.Feed) error {
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

	rules, err := compileRules(feedConfig.Rules)
	if err != nil {
		return fmt.Errorf("failed to compile rules: %w", err)
	}
	excludes, err := compileRules(feedConfig.Exclude)
	if err != nil {
		return fmt.Errorf("failed to compile exclude rules: %w", err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedConfig.URL)
	if err != nil {
		return errors.Join(errors.New("failed to parse feed"), err)
	}
	log.Info().Str("feed", feedConfig.Name).Str("url", feedConfig.URL).Int("items", len(feed.Items)).Str("title", feed.Title).Msg("downloading feed")

	for _, item := range feed.Items {
		matchingRules := len(rules) == 0
		matchingExclude := false

		for _, rule := range rules {
			if matched := rule.MatchString(item.Title); matched {
				matchingRules = true
				break
			}
		}

		for _, exclude := range excludes {
			if matched := exclude.MatchString(item.Title); matched {
				matchingExclude = true
				break
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

func compileRules(rules []config.Rule) ([]*regexp.Regexp, error) {
	var compiled []*regexp.Regexp
	for _, rule := range rules {
		if rule.Type == "regex" {
			re, err := regexp.Compile(rule.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid regex pattern %q: %w", rule.Value, err)
			}
			compiled = append(compiled, re)
		}
	}
	return compiled, nil
}

func downloadItem(url string, filename string) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.ContentLength > maxFileSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", resp.ContentLength, maxFileSize)
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	limitedReader := io.LimitReader(resp.Body, maxFileSize)
	_, err = io.Copy(file, limitedReader)
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
