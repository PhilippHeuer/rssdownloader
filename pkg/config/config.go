package config

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(filename string) (Config, error) {
	var config Config

	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func LoadState(filename string) (State, error) {
	var config State

	file, err := os.Open(filename)
	if err != nil {
		return State{FeedState: map[string]FeedState{}}, nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return State{}, err
	}

	return config, nil
}

func SaveState(filename string, state State) {
	file, _ := json.MarshalIndent(state, "", " ")
	_ = os.WriteFile(filename, file, 0644)
}
