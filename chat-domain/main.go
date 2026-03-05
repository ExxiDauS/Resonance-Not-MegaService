package main

import (
	chatservice "chat-domain/chat-service"
	configs "chat-domain/configs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize Redis
	rdcfg, err := configs.LoadRedisConfig()

	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     rdcfg.Address,
		Username: rdcfg.Username,
		Password: rdcfg.Password,
	})

	redisService := chatservice.NewRedisService(rdb)
	hub := chatservice.NewHub(redisService, rdb)
	handler := chatservice.NewChatHandler(hub, redisService)

	r := gin.Default()
	r.Use(cors.Default())

	// The frontend connects to this endpoint to join a room
	// Example: ws://localhost:8081/ws/chat/room123?user_id=userA
	r.GET("/ws/chat/:room_id", handler.HandleConnections)

	portcfg, err := configs.LoadPortConfig()
	if err != nil {
		panic(err)
	}

	r.Run(":" + portcfg.Port)
}
