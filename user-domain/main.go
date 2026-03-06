package main

import (
	database "user-domain/infrastructures/databases"
	"user-domain/infrastructures/messaging"
	userservice "user-domain/user-service"

	"log"
	"time"
	"user-domain/configs"
	"user-domain/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database client
	dbClient, err := database.NewPostgresDatabaseClient()

	if err != nil {
		panic("Failed to initialize database client: " + err.Error())
	}

	if err := dbClient.AutoMigrate(&userservice.Profile{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	repo := userservice.NewRepository(dbClient)
	service := userservice.NewService(repo)
	handler := userservice.NewHandler(service)

	// Initialize RabbitMQ consumer
	rabbitMQConfig, err := configs.LoadRabbitMQConfig()
	if err != nil {
		panic("Failed to load RabbitMQ configuration: " + err.Error())
	}

	consumer, err := messaging.NewRabbitMQConsumer(rabbitMQConfig.URL, service)
	if err != nil {
		panic("Failed to initialize RabbitMQ consumer: " + err.Error())
	}
	defer consumer.Close()

	// Start consuming messages
	if err := consumer.StartConsuming(); err != nil {
		panic("Failed to start RabbitMQ consumer: " + err.Error())
	}

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

	jwtSecret, err := configs.LoadJWTSecret()
	if err != nil {
		panic("Failed to load JWT secret: " + err.Error())
	}

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))

	protected.GET("/me", handler.GetMe)
	protected.POST("/profiles", handler.CreateUserProfile)
	protected.GET("/profiles", handler.GetAllUserProfiles)
	protected.GET("/profiles/:id", handler.GetUserProfileByID)
	protected.PATCH("/profiles/:id", handler.UpdateUserProfile)
	protected.DELETE("/profiles/:id", handler.DeleteUserProfile)

	log.Printf("User service starting on port %s", port)
	log.Printf("RabbitMQ consumer is running and waiting for messages...")
	r.Run(port)
}
