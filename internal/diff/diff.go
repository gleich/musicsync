package diff

import "slices"

func PlaylistDiff(appleMusicSongs []string, spotifySongs []string) ([]string, []string) {
	var (
		toAdd    []string
		toDelete []string
	)

	for _, appleMusicSong := range appleMusicSongs {
		if !slices.Contains(spotifySongs, appleMusicSong) {
			toAdd = append(toAdd, appleMusicSong)
		}
	}

	for _, spotifySong := range spotifySongs {
		if !slices.Contains(appleMusicSongs, spotifySong) {
			toDelete = append(toDelete, spotifySong)
		}
	}

	return toAdd, toDelete
}
