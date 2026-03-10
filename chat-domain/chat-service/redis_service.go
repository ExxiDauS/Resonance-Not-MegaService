package chatservice

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type ChatMessage struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

type RedisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{client: client}
}

// SaveMessage pushes a new message to the Redis stream and caps it at ~1000 messages
func (s *RedisService) SaveMessage(ctx context.Context, roomID string, msg ChatMessage) error {
	streamKey := fmt.Sprintf("room:%s:stream", roomID)

	err := s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		MaxLen: 1000,
		Approx: true,
		Values: map[string]interface{}{
			"user_id": msg.UserID,
			"content": msg.Content,
		},
	}).Err()

	return err
}

// GetHistory retrieves all previous messages for a user joining the room
func (s *RedisService) GetHistory(ctx context.Context, roomID string) ([]ChatMessage, error) {
	streamKey := fmt.Sprintf("room:%s:stream", roomID)

	// XRANGE fetches from the beginning ("-") to the end ("+")
	entries, err := s.client.XRange(ctx, streamKey, "-", "+").Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	var history []ChatMessage
	for _, entry := range entries {
		history = append(history, ChatMessage{
			UserID:  entry.Values["user_id"].(string),
			Content: entry.Values["content"].(string),
		})
	}
	return history, nil
}
