package authservice

import (
	"auth-domain/infrastructures/messaging"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *Repository
	jwtKey    []byte
	publisher *messaging.RabbitMQPublisher
}

func NewService(repo *Repository, secret string, publisher *messaging.RabbitMQPublisher) *Service {
	return &Service{
		repo:      repo,
		jwtKey:    []byte(secret),
		publisher: publisher,
	}
}

func (s *Service) Register(email, password, displayName string) (*UserResponse, error) {
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

	// Publish message to RabbitMQ for user profile creation
	if err := s.publisher.PublishUserCreated(user.UserID.String(), email, displayName); err != nil {
		log.Printf("Failed to publish user created event for user %s: %v", user.UserID, err)
		// Note: We don't fail the registration if message publishing fails
		// The user account is already created in the database
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
