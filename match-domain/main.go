package main

import (
	"log"

	"match-domain/configs"
	database "match-domain/infrastructures/databases"
	matchesservice "match-domain/matches-service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database client
	dbClient, err := database.NewPostgresDatabaseClient()

	if err != nil {
		panic("Failed to initialize database client: " + err.Error())
	}

	if err := dbClient.AutoMigrate(&matchesservice.Match{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	port, err := configs.LoadPort()
	if err != nil {
		panic("Failed to load port configuration: " + err.Error())
	}

	repo := matchesservice.NewRepository(dbClient)
	service := matchesservice.NewService(repo)
	handler := matchesservice.NewHandler(service)

	r := gin.Default()

	r.POST("/matches", handler.CreateMatch)
	r.GET("/matches/user/:user_id", handler.GetMatchesByUserID)
	r.GET("/matches/:match_id", handler.GetMatchByID)
	r.GET("/matches", handler.GetAllMatches)
	r.DELETE("/matches/:match_id", handler.DeleteMatch)

	log.Printf("Match service starting on port %s", port)
	r.Run(port)
}
