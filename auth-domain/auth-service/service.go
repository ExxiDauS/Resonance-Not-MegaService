package authservice

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   *Repository
	jwtKey []byte
}

func NewService(repo *Repository, secret string) *Service {
	return &Service{repo: repo, jwtKey: []byte(secret)}
}

func (s *Service) Register(email, password string) (*UserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &Credential{
		UserID:       uuid.New(),
		Email:        email,
		PasswordHash: string(hashed),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	// Return response DTO without password
	return &UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		LastLogin: user.LastLogin,
	}, nil
}

func (s *Service) Login(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)

	// email not found
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// password mismatch
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	claims := &Claims{
		UserID: user.UserID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Update last login time (non-critical, log if fails)
	if err := s.repo.UpdateLastLogin(user); err != nil {
		log.Printf("Failed to update last login for user %s: %v", user.UserID, err)
	}

	return token.SignedString(s.jwtKey)
}
