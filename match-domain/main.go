package main

import (
	"log"

	"match-domain/configs"
	database "match-domain/infrastructures/databases"
	"match-domain/infrastructures/messaging"
	matchesservice "match-domain/matches-service"
	waitingservice "match-domain/waiting-service"

	"github.com/gin-gonic/gin"

	"time"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "match_service_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "match_service_request_duration_seconds",
			Help:    "Response time duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(requestDuration)
}

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

	// Middleware for Prometheus metrics
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(c.FullPath(), status).Inc()
		requestDuration.WithLabelValues(c.FullPath()).Observe(duration)
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/user/:user_id", handler.GetMatchesByUserID)
	r.GET("/:match_id", handler.GetMatchByID)
	r.GET("/", handler.GetAllMatches)
	r.DELETE("/:match_id", handler.DeleteMatch)

	log.Printf("Match service starting on port %s", port)
	r.Run(":" + port)
}
