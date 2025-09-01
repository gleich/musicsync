package spotify

import (
	"fmt"
	"net/http"
	"strings"
)

type PlaylistResponse struct {
	Items []struct {
		Track struct {
			ExternalIDs struct {
				ISRC string `json:"isrc"`
			} `json:"external_ids"`
		} `json:"track"`
	} `json:"items"`
	Next string `json:"next"`
}

func PlaylistISRCs(client *http.Client, token *AccessToken, id string) ([]string, error) {
	path := fmt.Sprintf("/v1/playlists/%s/tracks", id)
	isrcs := []string{}
	for {
		resp, err := SendSpotifyAPIRequest[PlaylistResponse](client, token, path)
		if err != nil {
			return []string{}, fmt.Errorf("%w failed to get spotify playlist data for: %s", err, id)
		}
		for _, track := range resp.Items {
			isrcs = append(isrcs, track.Track.ExternalIDs.ISRC)
		}

		if resp.Next == "" {
			break
		}
		path = strings.TrimPrefix(resp.Next, "https://api.spotify.com/")
	}

	return isrcs, nil
}
