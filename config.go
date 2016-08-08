package main

import (
	"errors"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Banlist  string
	Logfiles []string
	MinScore float64
}

func NewConfig(filename string) (*Config, error) {
	// Load configuration file
	config := &Config{}
	meta, err := toml.DecodeFile(filename, config)
	if err != nil {
		return nil, err
	}

	// Ensure Banlist exists
	if !meta.IsDefined("Banlist") {
		return nil, errors.New("Config: must define Banlist")
	}

	// Ensure Logfiles exists
	if !meta.IsDefined("Logfiles") {
		return nil, errors.New("Config: must define Logfiles")
	}

	// Ensure MinScore exists
	if !meta.IsDefined("MinScore") {
		return nil, errors.New("Config: must define MinScore")
	}

	return config, nil
}
