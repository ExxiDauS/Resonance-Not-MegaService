package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type SoundCloudConfig struct {
	ClientID     string
	ClientSecret string
}

func LoadSoundCloudConfig() (*SoundCloudConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &SoundCloudConfig{
		ClientID:     os.Getenv("SOUNDCLOUD_CLIENT_ID"),
		ClientSecret: os.Getenv("SOUNDCLOUD_CLIENT_SECRET"),
	}
	return config, nil
}
