package applemusic

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"go.mattglei.ch/musicsync/internal/utils"
)

type Song struct {
	Name   string `json:"name"`
	ISRC   string `json:"isrc"`
	Artist string `json:"artistName"`
}

type CatalogSongsResponse struct {
	Data []struct {
		Attributes Song
	}
}

func PlaylistISRCs(client *http.Client, ids []string) ([]Song, error) {
	groups := utils.Batch(ids, 300)
	songs := []Song{}
	for _, group := range groups {
		ids := strings.Join(group, ",")
		params := url.Values{"ids": {ids}}
		searchedSongs, err := SendAppleMusicAPIRequest[CatalogSongsResponse](
			client,
			fmt.Sprintf("/v1/catalog/us/songs?%s", params.Encode()),
		)
		if err != nil {
			return []Song{}, fmt.Errorf(
				"%w failed to get catalog data for following ids: %s",
				err,
				ids,
			)
		}
		for _, song := range searchedSongs.Data {
			songs = append(songs, song.Attributes)
		}
	}

	return songs, nil
}
