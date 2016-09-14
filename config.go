/*
 *  zanproxy: a proxy detector for Zandronum
 *  Copyright (C) 2016  Alex Mayfield <alexmax2742@gmail.com>
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"errors"

	"github.com/BurntSushi/toml"
)

// Config contains the configuration for the program.
type Config struct {
	Banlist    string
	Email      string
	Logfiles   []string
	MinScore   float64
	BanMessage string
}

// NewConfig creates a new instance of Config from a configuration file.
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

	// Ensure Email exists
	if !meta.IsDefined("Email") {
		return nil, errors.New("Config: must define Email")
	}

	// Ensure Logfiles exists
	if !meta.IsDefined("Logfiles") {
		return nil, errors.New("Config: must define Logfiles")
	}

	// Ensure MinScore exists
	if !meta.IsDefined("MinScore") {
		return nil, errors.New("Config: must define MinScore")
	}

	// Ensure BanMessage exists, else use a sane default
	if !meta.IsDefined("BanMessage") {
		config.BanMessage = "You have been banned on suspicion of proxy use.  If you believe this is in error, please contact the administrators."
	}

	return config, nil
}
