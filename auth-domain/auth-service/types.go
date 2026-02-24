package authservice

import (
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
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
type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// JWT Claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}