package waitingservice

import (
	"time"

	"github.com/google/uuid"
)

type Waiting struct {
	ID 	  uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TrackId string `gorm:"type:text" json:"track_id"`
	CreatedAt time.Time `json:"created_at"`
}