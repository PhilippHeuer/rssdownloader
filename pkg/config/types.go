package config

import (
	"time"
)

type Config struct {
	Feeds []Feed `yaml:"feeds"`
}

type Feed struct {
	Name      string `yaml:"name"`
	URL       string `yaml:"url"`
	Rules     []Rule `yaml:"rules"`
	Exclude   []Rule `yaml:"exclude"`
	OutputDir string `yaml:"output"`
	Enabled   bool   `yaml:"enabled"`
}

type Rule struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

type State struct {
	FeedState map[string]FeedState `json:"feeds"`
}

type FeedState struct {
	FetchedAt time.Time `json:"timestamp"`
}
