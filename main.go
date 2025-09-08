package main

import (
	"fmt"
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

	for _, playlist := range config.Configuration.Playlists {
		fmt.Println()
		timber.Info("Running sync for playlist:", playlist.Name)
		appleMusicIDs, err := applemusic.PlaylistSongs(&httpClient, playlist.AppleMusicID)
		if err != nil {
			timber.Fatal(err, "failed to get apple music playlist")
		}
		timber.Done("[1/8] Found", len(appleMusicIDs), "from playlist in APPLE MUSIC")

		appleMusicSongs, err := applemusic.PlaylistISRCs(&httpClient, appleMusicIDs)
		if err != nil {
			timber.Fatal(err, "failed to get isrc for", len(appleMusicIDs), "ids from apple music")
		}
		timber.Done(
			"[2/8] Got",
			len(appleMusicSongs),
			"global isrc values for songs in APPLE MUSIC",
		)

		spotifySongs, err := spotify.PlaylistSongs(&spotifyClient, playlist.SpotifyID)
		if err != nil {
			timber.Fatal(err, "failed to get playlist data")
		}
		timber.Done(
			"[3/8] Found",
			len(spotifySongs),
			"songs in the current SPOTIFY playlist",
		)

		spotifyPlaylistSnapshotID, err := spotify.PlaylistSnapshot(
			&spotifyClient,
			playlist.SpotifyID,
		)
		if err != nil {
			timber.Fatal(err, "failed to get snapshot id for playlist")
		}
		timber.Done("[4/8] Got playlist version snapshot")

		toAdd, toDelete := diff.PlaylistDiff(appleMusicSongs, spotifySongs)
		timber.Done("[5/8]", len(toAdd), "songs to add.", len(toDelete), "songs to remove.")

		songsToAdd, err := spotify.FindAppleMusicSongs(&spotifyClient, toAdd)
		if err != nil {
			timber.Fatal(err, "failed to find isrcs in spotify")
		}
		timber.Done("[6/8]", len(toAdd), "songs to add.", len(toDelete), "songs to remove.")

		if len(toDelete) != 0 {
			err = spotify.EditSongs(
				&spotifyClient,
				playlist.SpotifyID,
				toDelete,
				&spotifyPlaylistSnapshotID,
			)
			if err != nil {
				timber.Fatal(err, "failed to remove songs to playlist")
			}
			timber.Done("[7/8]", "Removed", len(toDelete), "songs")
		} else {
			timber.Info("[7/8] Skipped as there are no songs to remove")
		}

		err = spotify.EditSongs(&spotifyClient, playlist.SpotifyID, songsToAdd, nil)
		if err != nil {
			timber.Fatal(err, "failed to add songs to playlist")
		}
		timber.Done("[8/8]", "Added", len(toAdd), "songs")
		fmt.Println()
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
