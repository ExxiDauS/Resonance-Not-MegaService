package matchesservice

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	MatchID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"match_id"`
	UserAID   uuid.UUID `gorm:"type:uuid;not null" json:"user_a_id"`
	UserBID   uuid.UUID `gorm:"type:uuid;not null" json:"user_b_id"`
	TrackedAt string `gorm:"type:text" json:"tracked_at"`
	CreatedAt time.Time `json:"created_at"`
}