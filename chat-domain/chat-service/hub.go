package chatservice

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Conn   *websocket.Conn
	UserID string
	RoomID string
}

type Hub struct {
	// Maps roomID -> list of active clients
	rooms map[string]map[*Client]bool
	mu    sync.RWMutex
	redis *RedisService
	rdb   *redis.Client
}

func NewHub(redisService *RedisService, rdb *redis.Client) *Hub {
	return &Hub{
		rooms: make(map[string]map[*Client]bool),
		redis: redisService,
		rdb:   rdb,
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[client.RoomID] == nil {
		h.rooms[client.RoomID] = make(map[*Client]bool)
		// Start a Redis Stream listener for this new room
		go h.listenToRoomStream(client.RoomID)
	}
	h.rooms[client.RoomID][client] = true
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[client.RoomID][client]; ok {
		delete(h.rooms[client.RoomID], client)
		client.Conn.Close()
	}
}

// listenToRoomStream blocks and waits for new messages in the Redis Stream
func (h *Hub) listenToRoomStream(roomID string) {
	ctx := context.Background()
	streamKey := fmt.Sprintf("room:%s:stream", roomID)
	lastID := "$" // Start listening for ONLY new messages

	for {
		// Stop listening if no one is in the room anymore to save routines
		h.mu.RLock()
		if len(h.rooms[roomID]) == 0 {
			h.mu.RUnlock()
			return
		}
		h.mu.RUnlock()

		streams, err := h.rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{streamKey, lastID},
			Count:   10,
			Block:   0, // Block indefinitely until a message arrives
		}).Result()

		if err != nil {
			log.Printf("Error reading stream for room %s: %v", roomID, err)
			time.Sleep(1 * time.Second) // Prevent tight loop on error
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				lastID = message.ID

				chatMsg := ChatMessage{
					UserID:  message.Values["user_id"].(string),
					Content: message.Values["content"].(string),
				}

				// Broadcast to all connected WebSockets in this room
				h.broadcastToRoom(roomID, chatMsg)
			}
		}
	}
}

func (h *Hub) broadcastToRoom(roomID string, msg ChatMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.rooms[roomID] {
		err := client.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error writing to client %s: %v", client.UserID, err)
			client.Conn.Close()
			delete(h.rooms[roomID], client)
		}
	}
}
