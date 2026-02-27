package tracksservice

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/zmb3/spotify/v2"
)

type AudioProvider interface {
	GetPlayableURLByQuery(ctx context.Context, seacrQuery string) (string, int, error)
}

func formatDuration(ms int) string {
	totalSeconds := ms / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

type Service struct {
	repo          Repository
	spotifyClient *spotify.Client
	audioService  AudioProvider
}

func NewService(repo Repository, spotifyClient *spotify.Client, audioService AudioProvider) *Service {
	return &Service{
		repo:          repo,
		spotifyClient: spotifyClient,
		audioService:  audioService,
	}
}

func (s *Service) getRandomTrackFromSpotify(ctx context.Context) ([]Track, error) {

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

	var tracks []Track

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

		track := Track{
			ID:       item.ID,
			Name:     item.Name,
			Artist:   artists,
			ImageURL: item.Album.Images[0].URL,
			Genre:    genres,
			Duration: "",
		}

		err := s.repo.CreateTrack(&track)
		if err != nil {
			return nil, fmt.Errorf("failed to create track in repo: %w", err)
		}

		tracks = append(tracks, track)
	}
	return tracks, nil
}

func (s *Service) GetAllTracks(ctx context.Context, page int, limit int, name *string, genre *string, artists *string) (*AllTracksResponse, error) {
	tracks, pagination, err := s.repo.GetAllTracks(page, limit, name, genre, artists)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracks: %w", err)
	}
	var trackResponses []Track
	for _, track := range tracks {
		trackResponses = append(trackResponses, Track{
			ID:       track.ID,
			Name:     track.Name,
			Artist:   track.Artist,
			ImageURL: track.ImageURL,
			Genre:    track.Genre,
			Duration: track.Duration,
		})
	}
	return &AllTracksResponse{
		Tracks:     trackResponses,
		Pagination: pagination,
	}, nil
}

func (s *Service) GetTrackByID(ctx context.Context, id string) (*TrackResponse, error) {
	track, err := s.repo.GetTrackByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get track by ID: %w", err)
	}
	if track == nil {
		return nil, fmt.Errorf("track not found")
	}
	audioURL, durationMs, err := s.audioService.GetPlayableURLByQuery(ctx, track.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get audio URL: %w", err)
	}
	duration := formatDuration(durationMs)
	return &TrackResponse{
		ID:       track.ID,
		Name:     track.Name,
		Artist:   track.Artist,
		ImageURL: track.ImageURL,
		Genre:    track.Genre,
		Duration: duration,
		AudioURL: audioURL,
	}, nil
}

func (s *Service) DeleteTrackByID(ctx context.Context, id string) error {
	_, err := s.repo.GetTrackByID(id)
	if err != nil {
		return fmt.Errorf("track not found: %w", err)
	}
	if err := s.repo.DeleteTrack(id); err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}
	return nil
}
