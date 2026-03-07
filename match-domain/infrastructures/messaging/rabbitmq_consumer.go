package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	waitingservice "match-domain/waiting-service"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SwipeEvent struct {
	UserID  string `json:"user_id"`
	TrackId string `json:"track_id"`
	Action  string `json:"action"` // "like" or "dislike"
}

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	service *waitingservice.Service
}

func NewRabbitMQConsumer(url string, service *waitingservice.Service) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	queueName := "swipe_queue"

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
	var event SwipeEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	log.Printf("Received swipe event: UserID=%s, TrackId=%s, Action=%s",
		event.UserID, event.TrackId, event.Action)

	// Parse UUID
	userID, err := uuid.Parse(event.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}
	

	if err := c.service.ProcessWaiting(userID, event.TrackId); err != nil {
		return fmt.Errorf("failed to process waiting: %w", err)
	}

	log.Printf("Successfully processed swipe event for UserID=%s, TrackId=%s", event.UserID, event.TrackId)
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
