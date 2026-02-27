package audiosservice

import "context"

type AudioService interface {
	GetPlayableURLByQuery(ctx context.Context, query string) (string, error)
}

type audioService struct {
	repo *soundCloudAudioRepo
}

func NewAudioService(repo *soundCloudAudioRepo) AudioService {
	return &audioService{repo: repo}
}

// This method perfectly satisfies the tracks-service AudioProvider interface
func (s *audioService) GetPlayableURLByQuery(ctx context.Context, query string) (string, error) {
	// 1. Search and get the protected URL
	protectedStreamURL, err := s.repo.SearchTrackStreamURL(ctx, query)
	if err != nil {
		return "", err
	}

	// 2. Convert it to a public, playable CDN URL for the frontend
	return s.repo.ResolvePlayableCDNURL(ctx, protectedStreamURL)
}
