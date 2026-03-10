package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type PortConfig struct {
	Port string
}

func LoadPortConfig() (*PortConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &PortConfig{
		Port: os.Getenv("PORT"),
	}

	return config, nil
}
