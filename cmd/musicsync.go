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
		spotifyClient = spotify.Client{
			HttpClient: &httpClient,
			Tokens:     &spotify.Tokens{RefreshToken: secrets.ENV.SpotifyRefreshToken},
		}
	)

	err := spotifyClient.Authorize()
	if err != nil {
		timber.Fatal(err, "failed to authorize spotify")
	}

	lcpClient := lcp.Client{Token: secrets.ENV.LcpToken}

	for {
		err = updateCycle(&httpClient, &spotifyClient, &lcpClient, newYork)
		if err != nil {
			timber.Warning("encountered error while trying to update", err.Error())
		}
	}
}

func updateCycle(
	httpClient *http.Client,
	spotifyClient *spotify.Client,
	lcpClient *lcp.Client,
	newYork *time.Location,
) error {
	playlists, err := lcp.FetchAppleMusicSyncedPlaylists(lcpClient)
	if err != nil {
		return fmt.Errorf("%w failed to fetch playlists to sync from lcp", err)
	}

	for _, playlist := range playlists {
		fmt.Println()
		timber.Info("Processing", playlist.Name)
		appleMusicIDs, err := applemusic.PlaylistSongs(httpClient, playlist.AppleMusicID)
		if err != nil {
			return fmt.Errorf("%w failed to get apple music playlist", err)
		}
		timber.Done("[1/9] Found", len(appleMusicIDs), "from playlist in APPLE MUSIC")

		appleMusicSongs, err := applemusic.PlaylistISRCs(httpClient, appleMusicIDs)
		if err != nil {
			return fmt.Errorf(
				"%w failed to get isrc for %d ids from apple music",
				err,
				len(appleMusicIDs),
			)
		}
		timber.Done(
			"[2/9] Got",
			len(appleMusicSongs),
			"global isrc values for songs in APPLE MUSIC",
		)

		spotifySongs, err := spotify.PlaylistSongs(spotifyClient, playlist.SpotifyID)
		if err != nil {
			return fmt.Errorf("%w failed to get playlist data", err)
		}
		timber.Done(
			"[3/9] Found",
			len(spotifySongs),
			"songs in the current SPOTIFY playlist",
		)

		spotifyPlaylistSnapshotID, err := spotify.PlaylistSnapshot(
			spotifyClient,
			playlist.SpotifyID,
		)
		if err != nil {
			return fmt.Errorf("%w failed to get snapshot id for playlist", err)
		}
		timber.Done("[4/9] Got playlist version snapshot")

		toAdd, toDelete := diff.PlaylistDiff(appleMusicSongs, spotifySongs)
		timber.Done("[5/9]", "Found playlist diff")

		var songsToAdd []spotify.Song
		if len(toAdd) != 0 {
			songsToAdd, err = spotify.FindAppleMusicSongs(spotifyClient, toAdd)
			if err != nil {
				return fmt.Errorf("%w failed to find isrcs in spotify", err)
			}
			timber.Done("[6/9]", "Found", len(songsToAdd), "in spotify from Apple Music")
		} else {
			timber.Info("[6/9]", "Skipping as there are no songs add")
		}
		songsToAdd, toDelete = diff.FilterPlaylists(songsToAdd, toDelete)

		if len(toDelete) != 0 {
			timber.Info("Deleting", len(toDelete), "songs")
			for _, song := range toDelete {
				timber.Info(fmt.Sprintf("- \"%s\" by \"%s\"", song.Name, song.Artist))
			}
			err = spotify.EditSongs(
				spotifyClient,
				playlist.SpotifyID,
				toDelete,
				&spotifyPlaylistSnapshotID,
			)
			if err != nil {
				return fmt.Errorf("%w failed to remove songs from playlist", err)
			}
			timber.Done("[7/9]", "Removed", len(toDelete), "songs")
		} else {
			timber.Info("[7/9] Skipped as there are no songs to remove")
		}

		if len(songsToAdd) != 0 {
			timber.Info("Adding", len(songsToAdd), "songs")
			for _, song := range songsToAdd {
				timber.Info(fmt.Sprintf("+ \"%s\" by \"%s\"", song.Name, song.Artist))
			}
			err = spotify.EditSongs(spotifyClient, playlist.SpotifyID, songsToAdd, nil)
			if err != nil {
				return fmt.Errorf("%w failed to add songs to playlist", err)
			}
			timber.Done("[8/9]", "Added", len(toAdd), "songs")
		} else {
			timber.Info("[8/9] Skipped as there are no songs to add")
		}

		if len(toAdd) != 0 || len(toDelete) != 0 {
			err = spotify.UpdateDescription(
				spotifyClient,
				playlist.SpotifyID,
				playlist.AppleMusicID,
				newYork,
			)
			if err != nil {
				return fmt.Errorf("%w failed to update playlist description", err)
			}
			timber.Info("[9/9] Updated playlist description")
		} else {
			timber.Info("[9/9] Skipped as playlist didn't get updated")
		}

		timber.Info("Waiting 5 minutes before syncing next playlist")
		time.Sleep(5 * time.Minute)
	}

	return nil
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
