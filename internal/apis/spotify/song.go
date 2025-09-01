package spotify

import (
	"fmt"
	"net/http"
	"net/url"

	"go.mattglei.ch/musicsync/internal/apis/applemusic"
	"go.mattglei.ch/timber"
)

type Song struct {
	ID   string
	ISRC string
}

type searchResponse struct {
	Tracks struct {
		Items []struct {
			ID string `json:"id"`
		}
	} `json:"tracks"`
}

func FindAppleMusicSongs(
	client *http.Client,
	token *AccessToken,
	appleMusicSongs []applemusic.Song,
) ([]Song, error) {
	songs := []Song{}
	for _, song := range appleMusicSongs {
		params := url.Values{
			"q":     {fmt.Sprintf("isrc:%s", song.ISRC)},
			"type":  {"track"},
			"limit": {"1"},
		}
		resp, err := SendSpotifyAPIRequest[searchResponse](
			client,
			token,
			fmt.Sprintf("/v1/search?%s", params.Encode()),
		)
		if err != nil {
			return []Song{}, fmt.Errorf("%w failed to search for song with isrc of %s", err, song)
		}
		if len(resp.Tracks.Items) == 0 {
			timber.Warning("isrc for", song.Name, "not found in spotify. ISRC:", song.ISRC)

			params.Set("q", fmt.Sprintf("track:\"%s\" artist:\"%s\"", song.Name, song.Artist))
			trackSearchResponse, err := SendSpotifyAPIRequest[searchResponse](
				client,
				token,
				fmt.Sprintf("/v1/search?%s", params.Encode()),
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
				timber.Warning("second search using", song.Name, "and", song.Artist, "failed")
				continue
			}
			resp = trackSearchResponse
		}
		foundSong := resp.Tracks.Items[0]
		songs = append(songs, Song{ID: foundSong.ID})
	}
	return songs, nil
}
