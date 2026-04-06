package cmd

import (
	"testing"

	"github.com/PhilippHeuer/rssdownloader/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestProcessFeedConfigErr(t *testing.T) {
	err := processFeed(config.Feed{
		Enabled: true,
	}, nil, nil)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, config.ErrRequiredConfigFieldMissing)
}

func TestProcessFeedConfigValid(t *testing.T) {
	feedConfig := config.Feed{
		Name:      "test-feed",
		URL:       "https://example.com/feed.xml",
		Enabled:   true,
		OutputDir: "/tmp",
		Template:  "{title}",
		Rules:     []config.Rule{},
		Exclude:   []config.Rule{},
	}

	validErrs := config.ValidateFeedConfig(feedConfig)
	assert.Empty(t, validErrs, "valid config should have no validation errors")
}
