package tracksservice

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/zmb3/spotify/v2"
)

type Service struct {
	repo          Repository
	spotifyClient *spotify.Client
}

func NewService(repo Repository, spotifyClient *spotify.Client) *Service {
	return &Service{
		repo:          repo,
		spotifyClient: spotifyClient,
	}
}

func (s *Service) getRandomTrackFromSpotify(ctx context.Context) ([]TrackResponse, error) {

	queries := []string{
		// --- BY YEAR RANGE (Popular Eras) ---
		"year:2020-2024", "year:2010-2019", "year:2000-2009", "year:1990-1999",
		"year:1980-1989", "year:1970-1979", "year:1960-1969", "year:1950-1959",

		// --- BY POPULAR GENRES ---
		"genre:pop", "genre:rock", "genre:hip-hop", "genre:electronic",
		"genre:r-n-b", "genre:jazz", "genre:classical", "genre:indie",
		"genre:latin", "genre:k-pop", "genre:country", "genre:metal",
		"genre:lo-fi", "genre:house", "genre:techno", "genre:reggae",

		// --- BY MAJOR LABELS (Great for finding specific aesthetics) ---
		"label:\"Universal Music Group\"", "label:\"Sony Music\"", "label:\"Warner Records\"",
		"label:\"Domino Recording Co\"", "label:\"XL Recordings\"", "label:\"Def Jam\"",
		"label:\"Sub Pop\"", "label:\"Brainfeeder\"", "label:\"Ninja Tune\"",

		// --- DISCOVERY TAGS (Spotify Specific) ---
		"tag:new",     // Albums released in the last 2 weeks
		"tag:hipster", // Albums with the lowest 10% popularity (underground)

		// --- COMBINATION QUERIES (Highly Popular use cases) ---
		"genre:pop year:2024",
		"genre:rock year:1970-1975",
		"genre:hip-hop year:1990-1995",
		"genre:electronic tag:new",
		"genre:indie tag:hipster",
		"genre:jazz label:Blue Note",
		"genre:classical year:2020-2024",

		// --- BY MOOD/STYLE (Using keywords with filters) ---
		"mood:happy", "mood:sad", "mood:chill", "mood:workout",
		"Christmas year:2023", "Halloween genre:rock",
		"acoustic genre:pop", "remix genre:electronic",
	}

	query := queries[rand.Intn(len(queries))]

	result, err := s.spotifyClient.Search(ctx, query, spotify.SearchTypeTrack)
	if err != nil {
		return nil, fmt.Errorf("spotify search failed: %w", err)
	}

	var tracks []TrackResponse

	for _, item := range result.Tracks.Tracks {

		// Concatenate artist names into a single string
		var artists string
		var genres string
		for i, artist := range item.Artists {
			if i > 0 {
				artists += ", "
			}
			artists += artist.Name
			fullArtist, err := s.spotifyClient.GetArtist(ctx, artist.ID)
			if err != nil {
				continue
			}
			if len(fullArtist.Genres) > 0 {
				for j, genre := range fullArtist.Genres {
					if j > 0 {
						genres += ", "
					}
					genres += genre
				}
			}
		}

		tracks = append(tracks, TrackResponse{
			ID:       item.ID,
			Name:     item.Name,
			Artist:   artists,
			ImageURL: item.Album.Images[0].URL,
			Genre:    genres,
			Duration: "0:00",
			AudioURL: "placeholder",
		})
	}
	return tracks, nil
}
