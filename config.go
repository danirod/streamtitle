// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	toml "github.com/pelletier/go-toml/v2"
)

const CONFIG_FILE = "StreamTitle/config.toml"
const CREDENTIALS_FILE = "StreamTitle/credentials.toml"

func configFile() string {
	file, err := xdg.ConfigFile(CONFIG_FILE)
	if err != nil {
		logger.Fatal(err)
	}
	return file
}

func credentialsFile() string {
	file, err := xdg.StateFile(CREDENTIALS_FILE)
	if err != nil {
		logger.Fatal(err)
	}
	return file
}

type appProfile struct {
	Game     int      `toml:"game" default:"0"`    // The game ID to use
	Language string   `toml:"language" default:""` // The language ISO code to use
	Title    string   `toml:"title" default:""`    // The title to set
	Tags     []string `toml:"tags" default:""`     // The tags to use
}

type appConfig struct {
	ClientId     string                `toml:"client_id"`     // The client ID used in OAuth
	ClientSecret string                `toml:"client_secret"` // the client secret used in OAuth
	Profiles     map[string]appProfile `toml:"profiles"`      // A set of profiles already created
}

func (cfg *appConfig) profiles() []string {
	profiles := []string{}
	for p := range cfg.Profiles {
		profiles = append(profiles, p)
	}
	return profiles
}

func (cfg *appConfig) profile(name string) (*appProfile, bool) {
	profile, found := cfg.Profiles[name]
	if found {
		return &profile, true
	}
	return nil, false
}

func (cfg *appConfig) read() error {
	path := configFile()
	logger.Print("Reading config file from ", path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	decoder.DisallowUnknownFields()
	return decoder.Decode(cfg)
}

type appState struct {
	AccessToken  string `toml:"access_token" default:""`  // The last used access token, in case is still valid
	RefreshToken string `toml:"refresh_token" default:""` // The last used refresh token, in case is still valid
}

func (state *appState) read() error {
	path := credentialsFile()
	logger.Print("Reading credentials from ", path)

	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		decoder := toml.NewDecoder(file)
		decoder.DisallowUnknownFields()
		return decoder.Decode(state)
	} else if errors.Is(err, os.ErrNotExist) {
		return nil
	} else {
		return err
	}
}

func (state *appState) write() error {
	path := credentialsFile()
	logger.Print("Updating credentials in ", path)

	// Make sure that the directory exists.
	dirname := filepath.Dir(path)
	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
		return err
	}

	// Create the config file.
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	err = encoder.Encode(state)
	return err
}
