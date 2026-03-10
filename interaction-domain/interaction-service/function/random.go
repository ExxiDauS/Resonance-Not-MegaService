package function

import "math/rand"

func RandomTracks(tracks []string, limit int) []string {

	if len(tracks) <= limit {
		return tracks
	}

	rand.Shuffle(len(tracks), func(i, j int) {
		tracks[i], tracks[j] = tracks[j], tracks[i]
	})

	return tracks[:limit]
}
