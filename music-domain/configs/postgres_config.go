package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DbURL string
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &DatabaseConfig{
		DbURL: os.Getenv("DATABASE_URL"),
	}

	return config, nil
}
