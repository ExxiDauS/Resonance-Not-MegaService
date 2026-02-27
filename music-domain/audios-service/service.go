package audiosservice

import "context"

type AudioService interface {
	GetPlayableURLByQuery(ctx context.Context, query string) (string, int, error)
}

type audioService struct {
	repo *soundCloudAudioRepo
}

func NewAudioService(repo *soundCloudAudioRepo) AudioService {
	return &audioService{repo: repo}
}

// This method perfectly satisfies the tracks-service AudioProvider interface
func (s *audioService) GetPlayableURLByQuery(ctx context.Context, query string) (string, int, error) {
	// 1. Search and get the protected URL and duration
	protectedStreamURL, durationMs, err := s.repo.SearchTrackStreamURL(ctx, query)
	if err != nil {
		return "", 0, err
	}

	// 2. Convert it to a public, playable CDN URL for the frontend
	cdnURL, err := s.repo.ResolvePlayableCDNURL(ctx, protectedStreamURL)
	if err != nil {
		return "", 0, err
	}

	return cdnURL, durationMs, nil
}
