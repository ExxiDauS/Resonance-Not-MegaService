package userservice

import (
	"time"

	"github.com/google/uuid"
)

// Database Model
type Profile struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id" binding:"required,uuid"`
	Email       string    `gorm:"unique;not null" json:"email" binding:"required,email"`
	DisplayName string    `gorm:"not null" json:"display_name" binding:"required"`
	Bio         string    `gorm:"type:text" json:"bio"`
	AvatarURL   string    `gorm:"type:text" json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateProfileInput struct {
	DisplayName *string `json:"display_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}
