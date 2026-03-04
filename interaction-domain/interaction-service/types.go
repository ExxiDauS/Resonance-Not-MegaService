package interactionservice

import "time"

type Track struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

type SwipeRequest struct {
	UserID  string `json:"user_id"`
	TrackID string `json:"track_id"`
	Action  string `json:"action"` // like / dislike
}

type Swipe struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	UserID    string `gorm:"index"`
	TrackID   string `gorm:"index"`
	Action    string `gorm:"type:varchar(20)"`
	CreatedAt time.Time
}
