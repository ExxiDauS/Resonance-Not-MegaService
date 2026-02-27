package utils

import (
	"strconv"
	"time"
)

func ParseSpotifyDate(dateStr string) time.Time {
	// Spotify release dates vary in precision
	layouts := []string{"2006-01-02", "2006-01", "2006"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t
		}
	}
	return time.Now() // Fallback
}

func ParseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
