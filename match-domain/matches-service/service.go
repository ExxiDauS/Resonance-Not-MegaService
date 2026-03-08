package matchesservice

import (
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateMatch(match *Match) error {
	return s.repo.CreateMatch(match)
}

func (s *Service) GetMatchesByUserID(userID string) ([]Match, error) {
	userID_UUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetMatchesByUserID(userID_UUID)
}

func (s *Service) GetMatchByID(matchID string) (*Match, error) {
	matchID_UUID, err := uuid.Parse(matchID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetMatchByID(matchID_UUID)
}

func (s *Service) GetAllMatches() ([]Match, error) {
	return s.repo.GetAllMatches()
}

func (s *Service) DeleteMatch(matchID string) error {
	matchID_UUID, err := uuid.Parse(matchID)
	if err != nil {
		return err
	}
	return s.repo.DeleteMatch(matchID_UUID)
}