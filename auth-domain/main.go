package main

import (
	"auth-domain/auth-service"
	database "auth-domain/infrastructures/databases"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database client
	dbClient, err := database.NewPostgresDatabaseClient()
	if err != nil {
		panic("Failed to initialize database client: " + err.Error())
	}

	repo := authservice.NewRepository(dbClient)
	service := authservice.NewService(repo, "mysecretkey")
	handler := authservice.NewHandler(service)

	r := gin.Default()
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.Run(":8080")
}
