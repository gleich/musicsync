package diff

import (
	"strings"

	"go.mattglei.ch/musicsync/internal/apis/applemusic"
	"go.mattglei.ch/musicsync/internal/apis/spotify"
)

func PlaylistDiff(
	appleMusicSongs []applemusic.Song,
	spotifySongs []spotify.Song,
) ([]applemusic.Song, []spotify.Song) {
	var (
		toAdd    []applemusic.Song
		toDelete []spotify.Song
	)

	keyApple := func(s applemusic.Song) string {
		if strings.TrimSpace(s.ISRC) != "" {
			return "isrc:" + strings.ToLower(strings.TrimSpace(s.ISRC))
		}
		name := strings.ToLower(strings.TrimSpace(s.Name))
		artist := strings.ToLower(strings.TrimSpace(s.Artist))
		return "na:" + name + "||" + artist
	}

	keySpotify := func(s spotify.Song) string {
		if strings.TrimSpace(s.ISRC) != "" {
			return "isrc:" + strings.ToLower(strings.TrimSpace(s.ISRC))
		}
		name := strings.ToLower(strings.TrimSpace(s.Name))
		artist := strings.ToLower(strings.TrimSpace(s.Artist))
		return "na:" + name + "||" + artist
	}

	appleSet := make(map[string]struct{}, len(appleMusicSongs))
	for _, s := range appleMusicSongs {
		appleSet[keyApple(s)] = struct{}{}
	}

	spotifySet := make(map[string]struct{}, len(spotifySongs))
	seenSpotify := make(map[string]struct{}, len(spotifySongs))

	for _, s := range spotifySongs {
		k := keySpotify(s)
		spotifySet[k] = struct{}{}

		if _, seen := seenSpotify[k]; seen {
			toDelete = append(toDelete, s)
			continue
		}
		seenSpotify[k] = struct{}{}

		if _, ok := appleSet[k]; !ok {
			toDelete = append(toDelete, s)
		}
	}

	for _, s := range appleMusicSongs {
		if _, ok := spotifySet[keyApple(s)]; !ok {
			toAdd = append(toAdd, s)
		}
	}

	return toAdd, toDelete
}
