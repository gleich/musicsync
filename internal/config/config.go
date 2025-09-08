package config

import (
	"github.com/BurntSushi/toml"
	"go.mattglei.ch/timber"
)

var Configuration ConfigurationData

type ConfigurationData struct {
	Playlists []Playlist `toml:"playlists"`
}

type Playlist struct {
	Name         string `toml:"name"`
	AppleMusicID string `toml:"apple_music"`
	SpotifyID    string `toml:"spotify"`
}

func Load() {
	_, err := toml.DecodeFile("config.toml", &Configuration)
	if err != nil {
		timber.Fatal(err, "failed to load configuration")
	}
	timber.Done("loaded", len(Configuration.Playlists), "playlists")
}
