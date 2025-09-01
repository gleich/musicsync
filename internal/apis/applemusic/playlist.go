package applemusic

import (
	"fmt"
	"net/http"
)

type PlaylistResponse struct {
	Data []struct {
		Attributes struct {
			PlayParams struct {
				ReportingID string `json:"reportingId"`
			} `json:"playParams"`
		} `json:"attributes"`
	} `json:"data"`
	Next string `json:"next"`
}

func PlaylistSongs(client *http.Client, id string) ([]string, error) {
	path := fmt.Sprintf("/v1/me/library/playlists/%s/tracks", id)
	ids := []string{}
	for {
		resp, err := SendAppleMusicAPIRequest[PlaylistResponse](client, path)
		if err != nil {
			return []string{}, fmt.Errorf(
				"%w failed to get apply music playlist data for: %s",
				err,
				path,
			)
		}
		for _, track := range resp.Data {
			ids = append(ids, track.Attributes.PlayParams.ReportingID)
		}

		if resp.Next == "" {
			break
		}
		path = resp.Next
	}
	return ids, nil
}
