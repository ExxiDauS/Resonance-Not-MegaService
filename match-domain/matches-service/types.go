package matchesservice

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrDuplicateMatch = errors.New("match already exists for this track")

type Match struct {
	MatchID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"match_id"`
	UserAID   uuid.UUID `gorm:"type:uuid;not null" json:"user_a_id"`
	UserBID   uuid.UUID `gorm:"type:uuid;not null" json:"user_b_id"`
	TrackID   string    `gorm:"type:text" json:"track_id"`
	CreatedAt time.Time `json:"created_at"`
}
