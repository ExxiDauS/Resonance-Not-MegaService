package waitingservice

import (
	"time"

	matchesservice "match-domain/matches-service"

	"github.com/google/uuid"
)

type Service struct {
	repo         *Repository
	matchService *matchesservice.Service
}

func NewService(repo *Repository, matchService *matchesservice.Service) *Service {
	return &Service{repo: repo, matchService: matchService}
}

// ProcessWaiting handles the matching logic:
//   - If no one is waiting with the same track_id, add the user to the waiting list.
//   - If someone is already waiting with the same track_id, create a match and remove
//     the matched user from the waiting list.
func (s *Service) ProcessWaiting(userID uuid.UUID, trackID string) error {
	existing, err := s.repo.FindByTrackID(trackID)
	if err != nil {
		return err
	}

	if existing == nil {
		entry := &Waiting{
			ID:        uuid.New(),
			UserID:    userID,
			TrackId:   trackID,
			CreatedAt: time.Now(),
		}
		return s.repo.AddToWaiting(entry)
	}

	// Prevent a user from being matched with themselves
	if existing.UserID == userID {
		return nil
	}

	match := &matchesservice.Match{
		MatchID:   uuid.New(),
		UserAID:   existing.UserID,
		UserBID:   userID,
		TrackedAt: trackID,
		CreatedAt: time.Now(),
	}
	if err := s.matchService.CreateMatch(match); err != nil {
		return err
	}

	return s.repo.DeleteByID(existing.ID)
}
