package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type SpotifyConfig struct {
	SpotifyClientID     string
	SpotifyClientSecret string
}

func LoadSpotifyConfig() (*SpotifyConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &SpotifyConfig{
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}

	return config, nil
}
