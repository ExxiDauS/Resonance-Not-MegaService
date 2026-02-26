package main

import (
	"log"
	"music-domain/configs"
	"music-domain/infrastructures/databases"
	"music-domain/infrastructures/spotify_client"
	tracksservice "music-domain/tracks-service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := configs.LoadPortConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := databases.NewPostgresDatabaseClient()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&tracksservice.Track{})

	spotifyClient, err := spotify_client.NewSpotifyClient()
	if err != nil {
		log.Fatalf("Failed to create Spotify client: %v", err)
	}

	trackRepo := tracksservice.NewRepository(db)
	trackService := tracksservice.NewService(trackRepo, spotifyClient)
	trackHandler := tracksservice.NewHandler(trackService)

	r := gin.Default()
	r.GET("/tracks", trackHandler.GetAllTracks)
	r.GET("/tracks/random", trackHandler.GetRandomTracks)
	r.GET("/tracks/:id", trackHandler.GetTrackByID)

	r.Run(":" + cfg.Port)
}
