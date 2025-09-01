package main

import (
	"net/http"
	"time"

	"go.mattglei.ch/musicsync/internal/apis/applemusic"
	"go.mattglei.ch/musicsync/internal/apis/spotify"
	"go.mattglei.ch/musicsync/internal/config"
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

	isrcs, err := spotify.PlaylistISRCs(&client, &accessToken, "5SnoWhWIJRmJNkvdxCpMAe")
	if err != nil {
		timber.Fatal(err, "failed to get playlist data")
	}

	timber.Debug(len(isrcs), "isrcs loaded")

	ids, err := applemusic.Playlists(&client, "p.AWXoZoxHLrvpJlY")
	if err != nil {
		timber.Fatal(err, "failed to get apple music playlist")
	}
	timber.Debug(len(ids), "songs from apple music playlist loaded")
}

func setupLogger() {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		timber.Fatal(err, "failed to load new york timezone")
	}
	timber.Timezone(ny)
	timber.TimeFormat("01/02 03:04:05 PM MST")
}
