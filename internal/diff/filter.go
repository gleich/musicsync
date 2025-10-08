package diff

import "go.mattglei.ch/musicsync/internal/apis/spotify"

func FilterPlaylists(
	toAdd []spotify.Song,
	toDelete []spotify.Song,
) ([]spotify.Song, []spotify.Song) {
	var (
		filteredToAdd    []spotify.Song
		filteredToDelete []spotify.Song
	)

	for _, songToAdd := range toAdd {
		contains := false
		for _, songToRemove := range toDelete {
			if songToAdd.ID == songToRemove.ID ||
				(songToAdd.Artist == songToRemove.Artist && songToAdd.Name == songToRemove.Name) {
				contains = true
				break
			}
		}
		if !contains {
			filteredToAdd = append(filteredToAdd, songToAdd)
		}
	}

	for _, songToRemove := range toDelete {
		contains := false
		for _, songToAdd := range toAdd {
			if songToRemove.ID == songToAdd.ID ||
				(songToRemove.Artist == songToAdd.Artist && songToRemove.Name == songToAdd.Name) {
				contains = true
				break
			}
		}
		if !contains {
			filteredToDelete = append(filteredToDelete, songToRemove)
		}
	}

	return filteredToAdd, filteredToDelete
}
