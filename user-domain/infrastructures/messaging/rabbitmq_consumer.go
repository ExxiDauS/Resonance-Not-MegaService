package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	userservice "user-domain/user-service"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserCreatedEvent struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	service *userservice.Service
}

func NewRabbitMQConsumer(url string, service *userservice.Service) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	queueName := "user_created"

	// Declare queue (idempotent operation)
	_, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Set QoS to process one message at a time
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	log.Printf("Successfully connected to RabbitMQ and declared queue: %s", queueName)

	return &RabbitMQConsumer{
		conn:    conn,
		channel: channel,
		queue:   queueName,
		service: service,
	}, nil
}

func (c *RabbitMQConsumer) StartConsuming() error {
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack (set to false for manual acknowledgment)
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Waiting for messages on queue: %s", c.queue)

	// Start consuming messages in a goroutine
	go func() {
		for msg := range msgs {
			if err := c.handleMessage(msg); err != nil {
				log.Printf("Error handling message: %v", err)
				// Reject the message and requeue it
				msg.Nack(false, true)
			} else {
				// Acknowledge the message
				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (c *RabbitMQConsumer) handleMessage(msg amqp.Delivery) error {
	var event UserCreatedEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	log.Printf("Received user created event: UserID=%s, Email=%s, DisplayName=%s",
		event.UserID, event.Email, event.DisplayName)

	// Parse UUID
	userID, err := uuid.Parse(event.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Create user profile
	profile := &userservice.Profile{
		UserID:      userID,
		Email:       event.Email,
		DisplayName: event.DisplayName,
		Bio:         "",
		AvatarURL:   "",
	}

	if err := c.service.CreateUserProfile(profile); err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}

	log.Printf("Successfully created profile for user: %s", event.UserID)
	return nil
}

func (c *RabbitMQConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	log.Println("RabbitMQ consumer connection closed")
}
