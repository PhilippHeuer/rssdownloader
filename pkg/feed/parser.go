package feed

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	defaultTimeout   = 30 * time.Second
	defaultUserAgent = "RSSDownloader/1.0"
)

func newClient(timeoutSeconds *int64) *http.Client {
	timeout := defaultTimeout
	if timeoutSeconds != nil {
		timeout = time.Duration(*timeoutSeconds) * time.Second
	}

	return &http.Client{Timeout: timeout}
}

func newUserAgent(userAgentOverride *string) string {
	if userAgentOverride != nil && *userAgentOverride != "" {
		return *userAgentOverride
	}
	return defaultUserAgent
}

func ParseFeed(url string, timeoutSeconds *int64, userAgentOverride *string) (*gofeed.Feed, error) {
	client := newClient(timeoutSeconds)
	ua := newUserAgent(userAgentOverride)

	fp := gofeed.NewParser()
	fp.UserAgent = ua
	fp.Client = client

	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed at %s: %w", url, err)
	}
	return feed, nil
}

func GetClient(timeoutSeconds *int64) *http.Client {
	return newClient(timeoutSeconds)
}
