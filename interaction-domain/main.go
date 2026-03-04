package main

import (
	"log"

	"interaction-domain/configs"
	database "interaction-domain/infrastructures/databases"
	interactionservice "interaction-domain/interaction-service"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func main() {
	db, err := database.NewPostgresDatabaseClient()
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate
	if db.AutoMigrate(&interactionservice.Swipe{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// RabbitMQ
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	ch.QueueDeclare("swipe_queue", true, false, false, false, nil)

	trackClient := interactionservice.NewHTTPTrackClient(
		"",
	)

	repo := interactionservice.NewSwipeRepository(db)
	svc := interactionservice.NewSwipeService(repo, ch, trackClient)
	h := interactionservice.NewSwipeHandler(svc)

	port, err := configs.LoadPort()
	if err != nil {
		panic("Failed to load port configuration: " + err.Error())
	}

	r := gin.Default()
	r.GET("/tracks/random", h.GetRandomTrack)
	r.POST("/swipe", h.Swipe)

	r.Run(port)
}
