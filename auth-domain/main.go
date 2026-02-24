package main

import (
	authservice "auth-domain/auth-service"
	"auth-domain/configs"
	database "auth-domain/infrastructures/databases"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database client
	dbClient, err := database.NewPostgresDatabaseClient()
	if err != nil {
		panic("Failed to initialize database client: " + err.Error())
	}

	dbClient.AutoMigrate(&authservice.Credential{})

	jwtSecret, err := configs.LoadJWTSecret()
	if err != nil {
		panic("Failed to load JWT secret: " + err.Error())
	}

	repo := authservice.NewRepository(dbClient)
	service := authservice.NewService(repo, jwtSecret)
	handler := authservice.NewHandler(service)

	r := gin.Default()
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/logout", handler.Logout)
	r.Run(":8080")
}
