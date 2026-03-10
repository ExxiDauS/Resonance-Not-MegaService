package interactionservice

import (
	"time"
)

type TrackResponse struct {
	Data    []MusicTrack `json:"data"`
	Success bool         `json:"success"`
}

type MusicTrack struct {
	ID     string `json:"ID"`
	Name   string `json:"Name"`
	Artist string `json:"Artist"`
}

type PlaylistTrack struct {
	ID         string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PlaylistID string `gorm:"index;not null"`
	TrackID    string `gorm:"index;not null"`
	AddedAt    time.Time
}

type PersonalPlaylist struct {
	ID     string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID string `gorm:"type:uuid;index"`
	Name   string

	Track []string `gorm:"-"`
}

type SwipeRequest struct {
	UserID  string `json:"userId" binding:"required,uuid"`
	TrackID string `json:"trackId" binding:"required"`
	Action  string `json:"action" binding:"required,oneof=like dislike"`
}

type Swipe struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;index" json:"user_id"`
	TrackID   string    `gorm:"index" json:"track_id"`
	Action    string    `gorm:"type:varchar(20)" json:"action"`
	CreatedAt time.Time `json:"created_at"`
}
