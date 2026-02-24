package tracksservice

import (
	"time"

	"github.com/zmb3/spotify/v2"
)

type Pagination struct {
	Page            int  `json:"page"`
	Limit           int  `json:"limit"`
	TotalCount      int  `json:"total_count"`
	TotalPage       int  `json:"total_page"`
	HasNextPage     bool `json:"has_next_page"`
	HasPreviousPage bool `json:"has_previous_page"`
}

type Track struct {
	ID        spotify.ID `gorm:"primaryKey"`
	Name      string     `gorm:"not null"`
	ImageURL  string     `gorm:"not null"`
	Artist    string     `gorm:"not null"`
	Genre     []string   `gorm:"not null"`
	Duration  string     `gorm:"not null"`
	AudioURL  string     `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TrackResponse struct {
	ID         spotify.ID `json:"id"`
	Name       string     `json:"name"`
	ImageURL   string     `json:"image_url"`
	Artist     string     `json:"artist"`
	Genre      []string   `json:"genre"`
	Duration   string     `json:"duration"`
	AudioURL   string     `json:"audio_url"`
	Pagination Pagination `json:"pagination"`
}
