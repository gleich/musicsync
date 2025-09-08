package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.mattglei.ch/musicsync/internal/utils"
)

type PlaylistResponse struct {
	SnapshotID string `json:"snapshot_id"`
}

type PlaylistTracksResponse struct {
	Items []struct {
		Track struct {
			ID          string `json:"id"`
			ExternalIDs struct {
				ISRC string `json:"isrc"`
			} `json:"external_ids"`
		} `json:"track"`
	} `json:"items"`
	Next string `json:"next"`
}

type addSongsPayload struct {
	URIs []string `json:"uris"`
}

type removeSongsPayload struct {
	Tracks     []track `json:"tracks"`
	SnapshotID string  `json:"snapshot_id"`
}

type track struct {
	URI string `json:"uri"`
}

func PlaylistSnapshot(client *SpotifyClient, id string) (string, error) {
	req := spotifyRequest{Method: http.MethodGet, Path: fmt.Sprintf("/v1/playlists/%s", id)}
	resp, err := sendSpotifyAPIRequest[PlaylistResponse](client, req)
	if err != nil {
		return "", fmt.Errorf("%w failed to make request for playlist data", err)
	}
	return resp.SnapshotID, nil
}

func PlaylistSongs(client *SpotifyClient, id string) ([]Song, error) {
	req := spotifyRequest{Method: http.MethodGet, Path: fmt.Sprintf("/v1/playlists/%s/tracks", id)}
	songs := []Song{}
	for {
		resp, err := sendSpotifyAPIRequest[PlaylistTracksResponse](client, req)
		if err != nil {
			return []Song{}, fmt.Errorf(
				"%w failed to get spotify playlist data for: %s",
				err,
				id,
			)
		}
		for _, track := range resp.Items {
			songs = append(songs, Song{
				ID:   track.Track.ID,
				ISRC: track.Track.ExternalIDs.ISRC,
			})
		}

		if resp.Next == "" {
			break
		}
		req.Path = strings.TrimPrefix(resp.Next, "https://api.spotify.com/")
	}

	return songs, nil
}

func EditSongs(client *SpotifyClient, id string, songs []Song, snapshotID *string) error {
	batches := utils.Batch(songs, 100)
	var method string
	if snapshotID == nil {
		method = http.MethodGet
	} else {
		method = http.MethodDelete
	}

	for _, batch := range batches {
		tracks := []string{}
		for _, song := range batch {
			tracks = append(tracks, fmt.Sprintf("spotify:track:%s", song.ID))
		}

		var payload any
		if snapshotID == nil {
			payload = addSongsPayload{URIs: tracks}
		} else {
			removeTracks := []track{}
			for _, t := range tracks {
				removeTracks = append(removeTracks, track{URI: t})
			}
			payload = removeSongsPayload{Tracks: removeTracks, SnapshotID: *snapshotID}
		}
		binary, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("%w failed to json marshal payload", err)
		}

		fmt.Println(string(binary))

		_, err = sendSpotifyAPIRequest[any](
			client,
			spotifyRequest{
				Method: method,
				Path:   fmt.Sprintf("/v1/playlists/%s/tracks", id),
				Body:   bytes.NewReader(binary),
			},
		)
		if err != nil {
			return fmt.Errorf("%w failed to send spotify api request", err)
		}
	}

	return nil
}
