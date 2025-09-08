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
		httpClient    = http.Client{Timeout: 20 * time.Second}
		spotifyClient = spotify.SpotifyClient{
			HttpClient: &httpClient,
			Tokens:     &spotify.Tokens{RefreshToken: secrets.ENV.SpotifyRefreshToken},
		}
	)

	err := spotifyClient.Authorize()
	if err != nil {
		timber.Fatal(err, "failed to authorize spotify")
	}

	appleMusicIDs, err := applemusic.PlaylistSongs(&httpClient, "p.AWXoZoxHLrvpJlY")
	if err != nil {
		timber.Fatal(err, "failed to get apple music playlist")
	}

	appleMusicSongs, err := applemusic.PlaylistISRCs(&httpClient, appleMusicIDs)
	if err != nil {
		timber.Fatal(err, "failed to get isrc for", len(appleMusicIDs), "ids from apple music")
	}

	spotifyPlaylistID := "6MLAGkQPdSBMjit5O1hrws"
	spotifySongs, err := spotify.PlaylistSongs(&spotifyClient, spotifyPlaylistID)
	if err != nil {
		timber.Fatal(err, "failed to get playlist data")
	}

	toAdd, toDelete := diff.PlaylistDiff(appleMusicSongs, spotifySongs)
	timber.Debug("toAdd:", toAdd)
	timber.Debug("toDelete:", toDelete)

	songsToAdd, err := spotify.FindAppleMusicSongs(&spotifyClient, toAdd)
	if err != nil {
		timber.Fatal(err, "failed to find isrcs in spotify")
	}

	err = spotify.AddSongs(&spotifyClient, spotifyPlaylistID, songsToAdd)
	if err != nil {
		timber.Fatal(err, "failed to add songs to playlist")
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
