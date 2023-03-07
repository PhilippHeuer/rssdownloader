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
	Enabled   bool   `yaml:"enabled"`
	OutputDir string `yaml:"output"`
	Template  string `yaml:"template"`
	Rules     []Rule `yaml:"rules"`
	Exclude   []Rule `yaml:"exclude"`
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
