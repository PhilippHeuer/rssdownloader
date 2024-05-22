package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrConfigFileInvalid          = fmt.Errorf("failed to load configuration from file")
	ErrStateFileInvalid           = fmt.Errorf("failed to load state from file")
	ErrFailedToWriteState         = fmt.Errorf("failed to write state to file")
	ErrRequiredConfigFieldMissing = fmt.Errorf("required field missing in config")
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
		return Config{}, errors.Join(ErrConfigFileInvalid, err)
	}

	return config, nil
}

func ValidateFeedConfig(feed Feed) (errs []error) {
	if feed.Name == "" {
		errs = append(errs, errors.Join(ErrRequiredConfigFieldMissing, fmt.Errorf("feed name is required")))
	}
	if feed.URL == "" {
		errs = append(errs, errors.Join(ErrRequiredConfigFieldMissing, fmt.Errorf("feed URL is required")))
	}
	if feed.Template == "" {
		errs = append(errs, errors.Join(ErrRequiredConfigFieldMissing, fmt.Errorf("feed template is required")))
	}
	if feed.OutputDir == "" {
		errs = append(errs, errors.Join(ErrRequiredConfigFieldMissing, fmt.Errorf("feed output directory is required")))
	}

	return errs
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
		return State{}, errors.Join(ErrStateFileInvalid, err)
	}

	return config, nil
}

func SaveState(filename string, state State) error {
	file, _ := json.MarshalIndent(state, "", " ")
	err := os.WriteFile(filename, file, 0644)
	if err != nil {
		return errors.Join(ErrFailedToWriteState, err)
	}

	return nil
}
