package spotify

import (
	"fmt"
	"net/http"
	"strings"
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

type Song struct {
	ID   string
	ISRC string
}

func PlaylistSongs(client *http.Client, token *AccessToken, id string) ([]Song, error) {
	path := fmt.Sprintf("/v1/playlists/%s/tracks", id)
	songs := []Song{}
	for {
		resp, err := SendSpotifyAPIRequest[PlaylistResponse](client, token, path)
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
		path = strings.TrimPrefix(resp.Next, "https://api.spotify.com/")
	}

	return songs, nil
}
