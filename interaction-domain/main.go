package main

import (
	"log"

	"interaction-domain/configs"
	database "interaction-domain/infrastructures/databases"
	interactionservice "interaction-domain/interaction-service"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func main() {

	// ---------------- DATABASE ----------------
	db, err := database.NewPostgresDatabaseClient()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&interactionservice.Swipe{},
		&interactionservice.PlaylistTrack{},
		&interactionservice.PersonalPlaylist{},
		&interactionservice.RecommendedPlaylist{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// ---------------- RABBITMQ CONFIG ----------------
	rabbitMQConfig, err := configs.LoadRabbitMQConfig()
	if err != nil {
		panic("Failed to load RabbitMQ configuration: " + err.Error())
	}

	// ---------------- RABBITMQ CONNECTION ----------------
	conn, err := amqp.Dial(rabbitMQConfig.URL)
	if err != nil {
		panic("Failed to connect to RabbitMQ: " + err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic("Failed to open RabbitMQ channel: " + err.Error())
	}
	defer ch.Close()

	// ---------------- DECLARE QUEUE ----------------
	_, err = ch.QueueDeclare(
		"swipe_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("Failed to declare queue: " + err.Error())
	}

	// ---------------- TRACK CLIENT ----------------
	trackClient := interactionservice.NewHTTPTrackClient(
		"http://localhost:8080",
	)

	// ---------------- SWIPE SERVICE ----------------
	repo := interactionservice.NewSwipeRepository(db)
	svc := interactionservice.NewSwipeService(repo, ch, trackClient)
	h := interactionservice.NewSwipeHandler(svc)

	// ---------------- PLAYLIST SERVICE ----------------
	playlistRepo := interactionservice.NewPlaylistRepository(db)
	playlistSvc := interactionservice.NewPlaylistService(playlistRepo)
	playlistHandler := interactionservice.NewPlaylistHandler(playlistSvc)

	// ---------------- PORT ----------------
	port, err := configs.LoadPort()
	if err != nil {
		panic("Failed to load port configuration: " + err.Error())
	}

	// ---------------- ROUTER ----------------
	r := gin.Default()

	r.GET("/tracks/random", h.GetRandomTrack)
	r.POST("/swipe", h.Swipe)

	r.POST("/playlists", playlistHandler.CreatePlaylist)
	r.GET("/playlists/:id", playlistHandler.GetPlaylist)
	r.GET("/users/:userId/playlists", playlistHandler.GetUserPlaylists)
	r.PUT("/playlists/:id", playlistHandler.UpdatePlaylistName)
	r.DELETE("/playlists/:id", playlistHandler.DeletePlaylist)

	r.POST("/playlists/:id/tracks/:trackId", playlistHandler.AddTrack)
	r.DELETE("/playlists/:id/tracks/:trackId", playlistHandler.RemoveTrack)

	r.GET("/users/:userId/recommended", playlistHandler.GetRecommended)

	log.Println("Server running on", port)
	r.Run(port)
}
