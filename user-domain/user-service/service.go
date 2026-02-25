package userservice

import (
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUserProfile(profile *Profile) error {
	return s.repo.CreateUserProfile(profile)
}

func (s *Service) GetAllUserProfiles() ([]Profile, error) {
	return s.repo.GetAllUserProfiles()
}

func (s *Service) GetUserProfileByID(userID string) (*Profile, error) {
	userID_UUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserProfileByID(userID_UUID)
}

func (s *Service) UpdateUserProfile(userID string, updates *UpdateProfileInput) error {
	userID_UUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return s.repo.UpdateUserProfile(userID_UUID, updates)
}

func (s *Service) DeleteUserProfile(userID string) error {
	userID_UUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return s.repo.DeleteUserProfile(userID_UUID)
}
