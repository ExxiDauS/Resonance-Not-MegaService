package main

import (
	authservice "auth-domain/auth-service"
	"auth-domain/configs"
	database "auth-domain/infrastructures/databases"
	"auth-domain/infrastructures/messaging"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database client
	dbClient, err := database.NewPostgresDatabaseClient()
	if err != nil {
		panic("Failed to initialize database client: " + err.Error())
	}

	if err := dbClient.AutoMigrate(&authservice.Credential{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Initialize RabbitMQ publisher
	rabbitMQConfig, err := configs.LoadRabbitMQConfig()
	if err != nil {
		panic("Failed to load RabbitMQ configuration: " + err.Error())
	}

	publisher, err := messaging.NewRabbitMQPublisher(rabbitMQConfig.URL)
	if err != nil {
		panic("Failed to initialize RabbitMQ publisher: " + err.Error())
	}
	defer publisher.Close()

	jwtSecret, err := configs.LoadJWTSecret()
	if err != nil {
		panic("Failed to load JWT secret: " + err.Error())
	}

	repo := authservice.NewRepository(dbClient)
	service := authservice.NewService(repo, jwtSecret, publisher)
	handler := authservice.NewHandler(service)

	port, err := configs.LoadPort()
	if err != nil {
		panic("Failed to load port configuration: " + err.Error())
	}

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/logout", handler.Logout)

	log.Printf("Auth service starting on port %s", port)
	r.Run(port)
}
