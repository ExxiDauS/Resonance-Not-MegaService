package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type RabbitMQConfig struct {
	URL string
}

func LoadRabbitMQConfig() (*RabbitMQConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	user := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")

	if host == "" || port == "" || user == "" || password == "" {
		return nil, fmt.Errorf("RabbitMQ configuration missing in environment variables")
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	return &RabbitMQConfig{URL: url}, nil
}
