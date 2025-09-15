package main

import (
	"fmt"
	"net/http"
	"time"

	"go.mattglei.ch/lcp/pkg/lcp"
	"go.mattglei.ch/musicsync/internal/apis/applemusic"
	"go.mattglei.ch/musicsync/internal/apis/spotify"
	"go.mattglei.ch/musicsync/internal/diff"
	"go.mattglei.ch/musicsync/internal/secrets"
	"go.mattglei.ch/timber"
)

func main() {
	newYork := setupLogger()
	timber.Done("booted")

	secrets.Load()

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

	lcpClient := lcp.Client{Token: secrets.ENV.LcpToken}
	playlists, err := lcp.FetchAppleMusicSyncedPlaylists(&lcpClient)
	if err != nil {
		timber.Fatal(err, "failed to fetch playlists to sync")
	}

	for _, playlist := range playlists {
		fmt.Println()
		timber.Debug(spotifyClient.Tokens.AccessToken)
		appleMusicIDs, err := applemusic.PlaylistSongs(&httpClient, playlist.AppleMusicID)
		if err != nil {
			timber.Fatal(err, "failed to get apple music playlist")
		}
		timber.Done("[1/9] Found", len(appleMusicIDs), "from playlist in APPLE MUSIC")

		appleMusicSongs, err := applemusic.PlaylistISRCs(&httpClient, appleMusicIDs)
		if err != nil {
			timber.Fatal(err, "failed to get isrc for", len(appleMusicIDs), "ids from apple music")
		}
		timber.Done(
			"[2/9] Got",
			len(appleMusicSongs),
			"global isrc values for songs in APPLE MUSIC",
		)

		spotifySongs, err := spotify.PlaylistSongs(&spotifyClient, playlist.SpotifyID)
		if err != nil {
			timber.Fatal(err, "failed to get playlist data")
		}
		timber.Done(
			"[3/9] Found",
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
		timber.Done("[4/9] Got playlist version snapshot")

		toAdd, toDelete := diff.PlaylistDiff(appleMusicSongs, spotifySongs)
		timber.Done("[5/9]", len(toAdd), "songs to add.", len(toDelete), "songs to remove.")

		var songsToAdd []spotify.Song
		if len(toAdd) != 0 {
			songsToAdd, err = spotify.FindAppleMusicSongs(&spotifyClient, toAdd)
			if err != nil {
				timber.Fatal(err, "failed to find isrcs in spotify")
			}
			timber.Done("[6/9]", "Found", len(songsToAdd), "in spotify from Apple Music")
		} else {
			timber.Info("[6/9]", "Skipping as there are no longs to lookup in Apple Music")
		}

		if len(toDelete) != 0 {
			timber.Info("Deleting:")
			for _, song := range toDelete {
				timber.Info(fmt.Sprintf("- \"%s\" by \"%s\"", song.Name, song.Artist))
			}
			err = spotify.EditSongs(
				&spotifyClient,
				playlist.SpotifyID,
				toDelete,
				&spotifyPlaylistSnapshotID,
			)
			if err != nil {
				timber.Fatal(err, "failed to remove songs to playlist")
			}
			timber.Done("[7/9]", "Removed", len(toDelete), "songs")
		} else {
			timber.Info("[7/9] Skipped as there are no songs to remove")
		}

		if len(toAdd) != 0 {
			timber.Info("Adding:")
			for _, song := range songsToAdd {
				timber.Info(fmt.Sprintf("+ \"%s\" by \"%s\"", song.Name, song.Artist))
			}
			err = spotify.EditSongs(&spotifyClient, playlist.SpotifyID, songsToAdd, nil)
			if err != nil {
				timber.Fatal(err, "failed to add songs to playlist")
			}
			timber.Done("[8/9]", "Added", len(toAdd), "songs")
		} else {
			timber.Info("[8/9] Skipped as there are no songs to add")
		}

		if len(toAdd) != 0 || len(toDelete) != 0 {
			err = spotify.UpdateDescription(
				&spotifyClient,
				playlist.SpotifyID,
				playlist.AppleMusicID,
				newYork,
			)
			if err != nil {
				timber.Fatal(err, "failed to updated playlist description")
			}
			timber.Info("[9/9] Updated playlist description")
		} else {
			timber.Info("[9/9] Skipped as playlist didn't get updated")
		}
		fmt.Println()
	}
}

func setupLogger() *time.Location {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		timber.Fatal(err, "failed to load new york timezone")
	}
	timber.Timezone(ny)
	timber.TimeFormat("01/02 03:04:05 PM MST")
	return ny
}
