package chatservice

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Configure properly for your frontend
	},
}

type ChatHandler struct {
	hub   *Hub
	redis *RedisService
}

func NewChatHandler(hub *Hub, redisService *RedisService) *ChatHandler {
	return &ChatHandler{hub: hub, redis: redisService}
}

func (h *ChatHandler) HandleConnections(c *gin.Context) {
	roomID := c.Param("room_id")
	userID := c.Query("user_id") // In production, extract this from the JWT token

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		Conn:   conn,
		UserID: userID,
		RoomID: roomID,
	}

	// 1. Fetch and send chat history immediately upon joining
	ctx := c.Request.Context()
	history, err := h.redis.GetHistory(ctx, roomID)
	if err == nil {
		for _, msg := range history {
			conn.WriteJSON(msg)
		}
	}

	// 2. Register client to receive real-time updates via XREAD
	h.hub.Register(client)
	defer h.hub.Unregister(client)

	// 3. Listen for incoming messages from this user and push to Redis Stream
	for {
		var msg ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg.UserID = userID // Enforce identity

		// Push to Redis (XADD). The background Hub listener will pick this up
		// and broadcast it back down to all users in the room.
		h.redis.SaveMessage(ctx, roomID, msg)
	}
}
