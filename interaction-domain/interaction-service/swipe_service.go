package interactionservice

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type SwipeService interface {
	GetRandomTrack(ctx context.Context) (*TrackResponse, error)
	Swipe(ctx context.Context, req *SwipeRequest) error
	GetUserSwipes(ctx context.Context, userID string) ([]Swipe, error)
	GetSwipe(ctx context.Context, swipe_id string) ([]Swipe, error)
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

func (s *swipeService) GetRandomTrack(ctx context.Context) (*TrackResponse, error) {
	return s.trackClient.GetRandomTrack(ctx)
}

func (s *swipeService) GetSwipe(ctx context.Context, swipe_id string) ([]Swipe, error) {
	return s.repo.GetSwipes(ctx, swipe_id)
}

func (s *swipeService) GetUserSwipes(ctx context.Context, userID string) ([]Swipe, error) {
	return s.repo.GetSwipesByUserID(ctx, userID)
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

	body, err := json.Marshal(swipe)
	if err != nil {
		return err
	}

	err = s.rabbitConn.Publish(
		"",
		"swipe_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("RabbitMQ publish error:", err)
	}

	return nil
}
