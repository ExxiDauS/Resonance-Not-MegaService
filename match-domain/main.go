package main

import (
	"log"

	"match-domain/configs"
	database "match-domain/infrastructures/databases"
	"match-domain/infrastructures/messaging"
	matchesservice "match-domain/matches-service"
	waitingservice "match-domain/waiting-service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database client
	dbClient, err := database.NewPostgresDatabaseClient()

	if err != nil {
		panic("Failed to initialize database client: " + err.Error())
	}

	if err := dbClient.AutoMigrate(&matchesservice.Match{}, &waitingservice.Waiting{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	port, err := configs.LoadPort()
	if err != nil {
		panic("Failed to load port configuration: " + err.Error())
	}

	matches_repo := matchesservice.NewRepository(dbClient)
	matches_service := matchesservice.NewService(matches_repo)

	waiting_repo := waitingservice.NewRepository(dbClient)
	waiting_service := waitingservice.NewService(waiting_repo, matches_service)
	handler := matchesservice.NewHandler(matches_service)

	// Initialize RabbitMQ consumer
	rabbitMQConfig, err := configs.LoadRabbitMQConfig()
	if err != nil {
		panic("Failed to load RabbitMQ configuration: " + err.Error())
	}

	consumer, err := messaging.NewRabbitMQConsumer(rabbitMQConfig.URL, waiting_service)
	if err != nil {
		panic("Failed to initialize RabbitMQ consumer: " + err.Error())
	}
	defer consumer.Close()

	// Start consuming messages
	if err := consumer.StartConsuming(); err != nil {
		panic("Failed to start RabbitMQ consumer: " + err.Error())
	}

	r := gin.Default()

	r.GET("/matches/user/:user_id", handler.GetMatchesByUserID)
	r.GET("/matches/:match_id", handler.GetMatchByID)
	r.GET("/matches", handler.GetAllMatches)
	r.DELETE("/matches/:match_id", handler.DeleteMatch)

	log.Printf("Match service starting on port %s", port)
	r.Run(":" + port)
}
