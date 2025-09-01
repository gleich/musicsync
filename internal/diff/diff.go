package diff

import (
	"go.mattglei.ch/musicsync/internal/apis/spotify"
)

func PlaylistDiff(
	appleMusicSongs []string,
	spotifySongs []spotify.Song,
) ([]string, []spotify.Song) {
	var (
		toAdd    []string
		toDelete []spotify.Song
	)

	for _, appleMusicSong := range appleMusicSongs {
		var contains = false
		for _, spotifySong := range spotifySongs {
			if spotifySong.ISRC == appleMusicSong {
				contains = true
				break
			}
		}
		if !contains {
			toAdd = append(toAdd, appleMusicSong)
		}
	}

	for _, spotifySong := range spotifySongs {
		var contains = false
		for _, appleMusicSong := range appleMusicSongs {
			if spotifySong.ISRC == appleMusicSong {
				contains = true
			}
		}
		if !contains {
			toDelete = append(toDelete, spotifySong)
		}
	}

	return toAdd, toDelete
}
