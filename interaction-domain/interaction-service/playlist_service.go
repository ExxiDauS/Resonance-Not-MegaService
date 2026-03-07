package interactionservice

import (
	"errors"
	"interaction-domain/interaction-service/function"

	"github.com/google/uuid"
)

type PlaylistService interface {
	CreatePlaylist(userID, name string) (*PersonalPlaylist, error)
	GetPlaylist(id string) (*PersonalPlaylist, error)
	GetUserPlaylists(userID string) ([]PersonalPlaylist, error)
	UpdatePlaylistName(id, name string) error
	DeletePlaylist(id string) error
	GetTracks(playlistID string) ([]PlaylistTrack, error)
	AddTrack(playlistID, trackID string) error
	RemoveTrack(playlistID, trackID string) error
	UpdateRecommended(userID string) error
	GetRecommended(userID string) (*PersonalPlaylist, error)
}

type playlistService struct {
	repo PlaylistRepository
}

func NewPlaylistService(r PlaylistRepository) PlaylistService {
	return &playlistService{r}
}

func (s *playlistService) CreatePlaylist(userID, name string) (*PersonalPlaylist, error) {
	playlist := &PersonalPlaylist{
		ID:     uuid.NewString(),
		UserID: userID,
		Name:   name,
	}
	return playlist, s.repo.Create(playlist)
}

func (s *playlistService) GetPlaylist(id string) (*PersonalPlaylist, error) {
	return s.repo.FindByID(id)
}

func (s *playlistService) GetUserPlaylists(userID string) ([]PersonalPlaylist, error) {
	return s.repo.FindByUserID(userID)
}

func (s *playlistService) UpdatePlaylistName(id, name string) error {
	return s.repo.UpdateName(id, name)
}

func (s *playlistService) DeletePlaylist(id string) error {
	return s.repo.Delete(id)
}

func (s *playlistService) GetTracks(playlistID string) ([]PlaylistTrack, error) {
	return s.repo.GetTracks(playlistID)
}

func (s *playlistService) AddTrack(playlistID, trackID string) error {
	return s.repo.AddTrack(playlistID, trackID)
}

func (s *playlistService) RemoveTrack(playlistID, trackID string) error {
	return s.repo.RemoveTrack(playlistID, trackID)
}
func (s *playlistService) UpdateRecommended(userID string) error {

	tracks, err := s.repo.GetUserSwipedTracks(userID)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return errors.New("no liked tracks found for user")
	}

	randomTracks := function.RandomTracks(tracks, 10)

	playlist, err := s.repo.GetOrCreateRecommendedPlaylist(userID)
	if err != nil {
		return err
	}

	// clear old tracks
	err = s.repo.ClearPlaylistTracks(playlist.ID)
	if err != nil {
		return err
	}

	for _, track := range randomTracks {
		err = s.repo.AddTrack(playlist.ID, track)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *playlistService) GetRecommended(userID string) (*PersonalPlaylist, error) {
	err := s.UpdateRecommended(userID)
	if err != nil {
		return nil, err
	}

	playlist, err := s.repo.GetOrCreateRecommendedPlaylist(userID)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(playlist.ID)
}
