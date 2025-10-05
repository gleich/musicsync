package diff

import (
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

	for _, appleMusicSong := range appleMusicSongs {
		contains := false
		for _, spotifySong := range spotifySongs {
			if spotifySong.ISRC == appleMusicSong.ISRC ||
				(spotifySong.Name == appleMusicSong.Name && spotifySong.Artist == appleMusicSong.Artist) {
				contains = true
				break
			}
		}
		if !contains {
			toAdd = append(toAdd, appleMusicSong)
		}
	}

	for _, spotifySong := range spotifySongs {
		contains := false
		for _, appleMusicSong := range appleMusicSongs {
			if spotifySong.ISRC == appleMusicSong.ISRC ||
				(spotifySong.Name == appleMusicSong.Name && spotifySong.Artist == appleMusicSong.Artist) {
				contains = true
			}
		}
		if !contains {
			toDelete = append(toDelete, spotifySong)
		}
	}

	return toAdd, toDelete
}
