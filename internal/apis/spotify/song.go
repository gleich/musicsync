package spotify

import (
	"fmt"
	"net/http"
	"net/url"

	"go.mattglei.ch/musicsync/internal/apis/applemusic"
)

type Song struct {
	ID     string
	ISRC   string
	Name   string
	Artist string
}

type songResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	ExternalIDs struct {
		ISRC string `json:"isrc"`
	} `json:"external_ids"`
}

type searchResponse struct {
	Tracks struct {
		Items []songResponse
	} `json:"tracks"`
}

func FindAppleMusicSongs(
	client *SpotifyClient,
	appleMusicSongs []applemusic.Song,
) ([]Song, error) {
	songs := []Song{}
	for _, song := range appleMusicSongs {
		params := url.Values{
			"q":     {fmt.Sprintf("isrc:%s", song.ISRC)},
			"type":  {"track"},
			"limit": {"1"},
		}
		resp, err := sendSpotifyAPIRequest[searchResponse](
			client,
			spotifyRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/v1/search?%s", params.Encode()),
			},
		)
		if err != nil {
			return []Song{}, fmt.Errorf("%w failed to search for song with isrc of %s", err, song)
		}
		if len(resp.Tracks.Items) == 0 {

			params.Set("q", fmt.Sprintf("track:\"%s\" artist:\"%s\"", song.Name, song.Artist))
			trackSearchResponse, err := sendSpotifyAPIRequest[searchResponse](
				client,
				spotifyRequest{
					Method: http.MethodGet,
					Path:   fmt.Sprintf("/v1/search?%s", params.Encode()),
				},
			)
			if err != nil {
				return []Song{}, fmt.Errorf(
					"%w failed to search for song with name of \"%s\" and artist of \"%s\"",
					err,
					song.Name,
					song.Artist,
				)
			}
			if len(trackSearchResponse.Tracks.Items) == 0 {
				continue
			}
			resp = trackSearchResponse
		}
		foundSong := resp.Tracks.Items[0]
		songs = append(
			songs,
			Song{ID: foundSong.ID, Artist: foundSong.Artists[0].Name, Name: foundSong.Name},
		)
	}
	return songs, nil
}
