package interactionservice

import (
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type PlaylistRepository interface {
	Create(playlist *PersonalPlaylist) error
	FindByID(id string) (*PersonalPlaylist, error)
	FindByUserID(userID string) ([]PersonalPlaylist, error)
	UpdateName(id string, name string) error
	Delete(id string) error
	GetTracks(playlistID string) ([]PlaylistTrack, error)
	AddTrack(playlistID string, trackID string) error
	RemoveTrack(playlistID string, trackID string) error

	GetOrCreateRecommendedPlaylist(userID string) (*PersonalPlaylist, error)
	GetUserSwipedTracks(userID string) ([]string, error)
	ClearPlaylistTracks(playlistID string) error
}

type playlistRepository struct {
	db *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) PlaylistRepository {
	return &playlistRepository{db}
}

func (r *playlistRepository) Create(p *PersonalPlaylist) error {
	return r.db.Create(p).Error
}

func (r *playlistRepository) FindByID(id string) (*PersonalPlaylist, error) {

	var playlist PersonalPlaylist
	err := r.db.
		Where("id = ?", id).
		First(&playlist).Error
	if err != nil {
		return nil, err
	}

	var tracks []PlaylistTrack
	err = r.db.
		Where("playlist_id = ?", id).
		Find(&tracks).Error
	if err != nil {
		return nil, err
	}

	for _, t := range tracks {
		playlist.Track = append(playlist.Track, t.TrackID)
	}

	return &playlist, nil
}

func (r *playlistRepository) FindByUserID(userID string) ([]PersonalPlaylist, error) {
	var playlists []PersonalPlaylist
	err := r.db.
		Where("user_id = ?", userID).
		Find(&playlists).Error
	if err != nil {
		return nil, err
	}
	if len(playlists) == 0 {
		return playlists, nil
	}
	var playlistIDs []string
	for _, p := range playlists {
		playlistIDs = append(playlistIDs, p.ID)
	}
	var tracks []PlaylistTrack
	err = r.db.
		Where("playlist_id IN ?", playlistIDs).
		Find(&tracks).Error

	if err != nil {
		return nil, err
	}
	trackMap := map[string][]string{}
	for _, t := range tracks {
		trackMap[t.PlaylistID] = append(trackMap[t.PlaylistID], t.TrackID)
	}
	for i := range playlists {
		playlists[i].Track = trackMap[playlists[i].ID]
	}
	return playlists, nil
}

func (r *playlistRepository) UpdateName(id string, name string) error {
	return r.db.Model(&PersonalPlaylist{}).
		Where("id = ?", id).
		Update("name", name).Error
}

func (r *playlistRepository) Delete(id string) error {
	return r.db.Delete(&PersonalPlaylist{}, "id = ?", id).Error
}

func (r *playlistRepository) GetTracks(playlistID string) ([]PlaylistTrack, error) {
	var tracks []PlaylistTrack
	err := r.db.
		Where("playlist_id = ?", playlistID).
		Find(&tracks).Error
	return tracks, err
}

func (r *playlistRepository) AddTrack(playlistID, trackID string) error {

	var existing PlaylistTrack

	err := r.db.
		Where("playlist_id = ? AND track_id = ?", playlistID, trackID).
		First(&existing).Error

	if err == nil {
		return nil
	}

	return r.db.Create(&PlaylistTrack{
		PlaylistID: playlistID,
		TrackID:    trackID,
	}).Error
}

func (r *playlistRepository) RemoveTrack(playlistID string, trackID string) error {
	return r.db.
		Where("playlist_id = ? AND track_id = ?", playlistID, trackID).
		Delete(&PlaylistTrack{}).Error
}

func (r *playlistRepository) GetUserSwipedTracks(userID string) ([]string, error) {

	var swipes []Swipe

	err := r.db.
		Where("user_id = ? AND action = ?", userID, "like").
		Find(&swipes).Error

	if err != nil {
		return nil, err
	}

	var trackIDs []string
	for _, s := range swipes {
		trackIDs = append(trackIDs, s.TrackID)
	}

	return trackIDs, nil
}

func (r *playlistRepository) GetOrCreateRecommendedPlaylist(userID string) (*PersonalPlaylist, error) {

	var playlist PersonalPlaylist

	err := r.db.
		Where("user_id = ? AND name = ?", userID, "Recommended").
		First(&playlist).Error

	if err == nil {
		return &playlist, nil
	}

	if err == gorm.ErrRecordNotFound {

		playlist = PersonalPlaylist{
			ID:     uuid.NewString(),
			UserID: userID,
			Name:   "Recommended",
		}

		if err := r.db.Create(&playlist).Error; err != nil {
			return nil, err
		}

		return &playlist, nil
	}

	return nil, err
}

func (r *playlistRepository) ClearPlaylistTracks(playlistID string) error {
	return r.db.Where("playlist_id = ?", playlistID).
		Delete(&PlaylistTrack{}).Error
}
