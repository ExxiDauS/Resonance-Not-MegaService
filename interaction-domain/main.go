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
	db, err := database.NewPostgresDatabaseClient()
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&interactionservice.Swipe{},
		&interactionservice.PlaylistTrack{},
		&interactionservice.PersonalPlaylist{},
		&interactionservice.RecommendedPlaylist{},
	); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}

	trackClient := interactionservice.NewHTTPTrackClient(
		"http://localhost:8080",
	)

	repo := interactionservice.NewSwipeRepository(db)
	svc := interactionservice.NewSwipeService(repo, ch, trackClient)
	h := interactionservice.NewSwipeHandler(svc)

	playlistRepo := interactionservice.NewPlaylistRepository(db)
	playlistSvc := interactionservice.NewPlaylistService(playlistRepo)
	playlistHandler := interactionservice.NewPlaylistHandler(playlistSvc)

	port, err := configs.LoadPort()
	if err != nil {
		panic("Failed to load port configuration: " + err.Error())
	}

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

	r.Run(port)
}
