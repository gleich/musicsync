package main

import (
	"net/http"
	"time"

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

	// playlist, err := applemusic.SendAppleMusicAPIRequest[applemusic.PlaylistResponse](
	// 	&client,
	// 	"/v1/me/library/playlists/p.AWXoZoxHLrvpJlY/tracks",
	// )
	// if err != nil {
	// 	timber.Fatal(err, "failed to make apple music api request")
	// }

	// _, err = applemusic.PlaylistISRCs(&client, playlist)
	// if err != nil {
	// 	timber.Fatal(err, "failed to load playlist isrcs")
	// }

	playlist, err := spotify.SendSpotifyAPIRequest[spotify.PlaylistResponse](
		&client,
		&accessToken,
		"/v1/playlists/6MLAGkQPdSBMjit5O1hrws",
	)
	if err != nil {
		timber.Fatal(err, "failed to load playlist")
	}

	timber.Debug("track length:", len(playlist.Tracks.Items))
	for _, track := range playlist.Tracks.Items {
		timber.Debug(track.Track.ExternalIDs.ISRC)
	}
}

func setupLogger() {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		timber.Fatal(err, "failed to load new york timezone")
	}
	timber.Timezone(ny)
	timber.TimeFormat("01/02 03:04:05 PM MST")
}
