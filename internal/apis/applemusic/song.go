package applemusic

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type CatalogSongsResponse struct {
	Data []struct {
		Attributes struct {
			ISRC string `json:"isrc"`
		}
	}
}

func PlaylistISRCs(client *http.Client, playlist PlaylistResponse) ([]string, error) {
	// break down songs into chunks of 300 songs (limit for searching using this endpoint)
	var (
		groups     = [][]string{{}}
		added      = 0
		groupIndex = 0
	)
	for _, song := range playlist.Data {
		if len(groups[groupIndex]) > 300 {
			groups = append(groups, []string{})
			groupIndex++
		}
		groups[groupIndex] = append(groups[groupIndex], song.Attributes.PlayParams.ReportingID)
		added++
	}

	isrcs := []string{}
	for _, group := range groups {
		ids := strings.Join(group, ",")
		params := url.Values{"ids": {ids}}
		songs, err := SendAppleMusicAPIRequest[CatalogSongsResponse](
			client,
			fmt.Sprintf("/v1/catalog/us/songs?%s", params.Encode()),
		)
		if err != nil {
			return []string{}, fmt.Errorf(
				"%w failed to get catalog data for following ids: %s",
				err,
				ids,
			)
		}
		for _, song := range songs.Data {
			isrcs = append(isrcs, song.Attributes.ISRC)
		}
	}

	return isrcs, nil
}
