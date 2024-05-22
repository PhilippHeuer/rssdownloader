package cmd

import (
	"testing"

	"github.com/PhilippHeuer/rssdownloader/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestProcessFeedConfigErr(t *testing.T) {
	err := processFeed(config.Feed{
		Enabled: true,
	})

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, config.ErrRequiredConfigFieldMissing)
}

func TestProcessFeedConfigValid(t *testing.T) {
	// TODO: need to mock http response
}
