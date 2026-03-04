package interactionservice

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"
)

type SwipeService interface {
	GetRandomTrack(ctx context.Context) (*Track, error)
	Swipe(ctx context.Context, req *SwipeRequest) error
}

type swipeService struct {
	repo        SwipeRepository
	rabbitConn  *amqp.Channel
	trackClient TrackClient
}

func NewSwipeService(
	repo SwipeRepository,
	rabbit *amqp.Channel,
	trackClient TrackClient,
) SwipeService {
	return &swipeService{
		repo:        repo,
		rabbitConn:  rabbit,
		trackClient: trackClient,
	}
}

func (s *swipeService) GetRandomTrack(ctx context.Context) (*Track, error) {
	return s.trackClient.GetRandomTrack(ctx)
}

func (s *swipeService) Swipe(ctx context.Context, req *SwipeRequest) error {
	swipe := &Swipe{
		UserID:  req.UserID,
		TrackID: req.TrackID,
		Action:  req.Action,
	}

	if err := s.repo.Save(ctx, swipe); err != nil {
		return err
	}

	body, _ := json.Marshal(swipe)

	return s.rabbitConn.Publish(
		"",
		"swipe_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
