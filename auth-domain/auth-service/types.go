package authservice

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Database Model
type Credential struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	LastLogin    *time.Time
}

// Request DTOs
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name" binding:"required,min=1"`
}

type LoginRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
}

// Response DTOs
type UserResponse struct {
	UserID    uuid.UUID  `json:"user_id"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}

// JWT Claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
