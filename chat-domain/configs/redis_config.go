package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type RedisConfig struct {
	Address  string
	Username string
	Password string
}

func LoadRedisConfig() (*RedisConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &RedisConfig{
		Address:  os.Getenv("REDIS_ADDRESS"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	return config, nil
}
