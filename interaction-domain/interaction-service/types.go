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
	CreatedAt  time.Time
}

type PersonalPlaylist struct {
	ID     string `gorm:"primaryKey"`
	UserID string `gorm:"index"`
	Name   string

	Track []string `gorm:"-"`
}

type RecommendedPlaylist struct {
	UserID string `gorm:"primaryKey"`

	Tracks []string `gorm:"serializer:json"`
}

type SwipeRequest struct {
	UserID  string `json:"userId" binding:"required"`
	TrackID string `json:"trackId" binding:"required"`
	Action  string `json:"action" binding:"required,oneof=like dislike"`
}

type Swipe struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	UserID    string `gorm:"index"`
	TrackID   string `gorm:"index"`
	Action    string `gorm:"type:varchar(20)"`
	CreatedAt time.Time
}
