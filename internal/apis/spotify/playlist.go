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

func PlaylistSongs(client *SpotifyClient, id string) ([]Song, error) {
	req := spotifyRequest{Method: http.MethodGet, Path: fmt.Sprintf("/v1/playlists/%s/tracks", id)}
	songs := []Song{}
	for {
		resp, err := sendSpotifyAPIRequest[PlaylistResponse](client, req)
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

func AddSongs(client *SpotifyClient, id string, songs []Song) error {
	batches := utils.Batch(songs, 100)

	for _, batch := range batches {
		payload := addSongsPayload{URIs: []string{}}
		for _, song := range batch {
			payload.URIs = append(payload.URIs, fmt.Sprintf("spotify:track:%s", song.ID))
		}

		binary, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("%w failed to wrap payload", err)
		}

		_, err = sendSpotifyAPIRequest[any](
			client,
			spotifyRequest{
				Method: http.MethodPost,
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
