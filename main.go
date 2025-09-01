package main

import (
	"net/http"
	"time"

	"go.mattglei.ch/musicsync/internal/apis/applemusic"
	"go.mattglei.ch/musicsync/internal/apis/spotify"
	"go.mattglei.ch/musicsync/internal/config"
	"go.mattglei.ch/musicsync/internal/diff"
	"go.mattglei.ch/musicsync/internal/secrets"
	"go.mattglei.ch/timber"
)

func main() {
	setupLogger()
	timber.Done("booted")

	secrets.Load()
	config.Load()

	var (
		client = http.Client{Timeout: 20 * time.Second}
	)

	accessToken, err := spotify.Authorize(&client)
	if err != nil {
		timber.Fatal(err, "failed to authorize spotify")
	}

	appleMusicIDs, err := applemusic.PlaylistSongs(&client, "p.qQXLxPpFA75zg8e")
	if err != nil {
		timber.Fatal(err, "failed to get apple music playlist")
	}

	appleMusicSongs, err := applemusic.PlaylistISRCs(&client, appleMusicIDs)
	if err != nil {
		timber.Fatal(err, "failed to get isrc for", len(appleMusicIDs), "ids from apple music")
	}

	spotifySongs, err := spotify.PlaylistSongs(&client, &accessToken, "6MLAGkQPdSBMjit5O1hrws")
	if err != nil {
		timber.Fatal(err, "failed to get playlist data")
	}

	toAdd, toDelete := diff.PlaylistDiff(appleMusicSongs, spotifySongs)
	timber.Debug("toAdd:", toAdd)
	timber.Debug("toDelete:", toDelete)
}

func setupLogger() {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		timber.Fatal(err, "failed to load new york timezone")
	}
	timber.Timezone(ny)
	timber.TimeFormat("01/02 03:04:05 PM MST")
}
