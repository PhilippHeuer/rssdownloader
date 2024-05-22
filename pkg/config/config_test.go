package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFeedConfigValid(t *testing.T) {
	errs := ValidateFeedConfig(Feed{
		Name:      "Test Feed",
		URL:       "https://example.com/feed.xml",
		Enabled:   true,
		OutputDir: "/tmp",
		Template:  "template.tmpl",
		Rules:     []Rule{},
		Exclude:   []Rule{},
	})

	if len(errs) != 0 {
		t.Errorf("expected no errors, got %d", len(errs))
	}
}

func TestValidateFeedConfigInvalid(t *testing.T) {
	errs := ValidateFeedConfig(Feed{
		Enabled: true,
	})

	assert.NotNil(t, errs)
	assert.ErrorIs(t, errs[0], ErrRequiredConfigFieldMissing)
	assert.Contains(t, errs[0].Error(), "feed name is required")
	assert.ErrorIs(t, errs[1], ErrRequiredConfigFieldMissing)
	assert.Contains(t, errs[1].Error(), "feed URL is required")
	assert.ErrorIs(t, errs[2], ErrRequiredConfigFieldMissing)
	assert.Contains(t, errs[2].Error(), "feed template is required")
	assert.ErrorIs(t, errs[3], ErrRequiredConfigFieldMissing)
	assert.Contains(t, errs[3].Error(), "feed output directory is required")
}
