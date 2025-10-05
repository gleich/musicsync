package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.mattglei.ch/musicsync/internal/utils"
)

type PlaylistResponse struct {
	SnapshotID string `json:"snapshot_id"`
}

type PlaylistTracksResponse struct {
	Items []struct {
		Track songResponse `json:"track"`
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
				ID:     track.Track.ID,
				ISRC:   track.Track.ExternalIDs.ISRC,
				Name:   track.Track.Name,
				Artist: track.Track.Artists[0].Name,
			})
		}

		if resp.Next == "" {
			break
		}

		req.Path = strings.TrimPrefix(resp.Next, "https://api.spotify.com")
	}

	return songs, nil
}

func EditSongs(client *SpotifyClient, id string, songs []Song, snapshotID *string) error {
	batches := utils.Batch(songs, 100)
	var method string
	if snapshotID == nil {
		method = http.MethodPost
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

func UpdateDescription(
	client *SpotifyClient,
	spotifyID string,
	appleMusicID string,
	location *time.Location,
) error {
	description := fmt.Sprintf(
		"https://mattglei.ch/music/playlists/%s. Auto updated %s.",
		appleMusicID,
		time.Now().In(location).Format("January 2 2006 at 3:04pm MST"),
	)

	binary, err := json.Marshal(struct {
		Description string `json:"description"`
	}{Description: description})
	if err != nil {
		return fmt.Errorf("%w failed to marshal JSON", err)
	}

	_, err = sendSpotifyAPIRequest[any](client, spotifyRequest{
		Method:           http.MethodPut,
		Path:             fmt.Sprintf("/v1/playlists/%s", spotifyID),
		Body:             bytes.NewReader(binary),
		NotExpectingJSON: true,
	})
	if err != nil {
		return fmt.Errorf("%w failed to send spotify api request", err)
	}

	return nil
}
