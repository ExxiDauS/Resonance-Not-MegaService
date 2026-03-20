package main

import (
	"log"

	"interaction-domain/configs"
	database "interaction-domain/infrastructures/databases"
	interactionservice "interaction-domain/interaction-service"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"

	"time"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "interaction_service_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "interaction_service_request_duration_seconds",
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

	// ---------------- DATABASE ----------------
	db, err := database.NewPostgresDatabaseClient()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&interactionservice.Swipe{},
		&interactionservice.PlaylistTrack{},
		&interactionservice.PersonalPlaylist{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// ---------------- RABBITMQ CONFIG ----------------
	rabbitMQConfig, err := configs.LoadRabbitMQConfig()
	if err != nil {
		panic("Failed to load RabbitMQ configuration: " + err.Error())
	}

	// ---------------- RABBITMQ CONNECTION ----------------
	conn, err := amqp.Dial(rabbitMQConfig.URL)
	if err != nil {
		panic("Failed to connect to RabbitMQ: " + err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic("Failed to open RabbitMQ channel: " + err.Error())
	}
	defer ch.Close()

	// ---------------- DECLARE QUEUE ----------------
	_, err = ch.QueueDeclare(
		"swipe_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("Failed to declare queue: " + err.Error())
	}

	// ---------------- SWIPE SERVICE ----------------
	repo := interactionservice.NewSwipeRepository(db)
	svc := interactionservice.NewSwipeService(repo, ch)
	h := interactionservice.NewSwipeHandler(svc)

	// ---------------- PLAYLIST SERVICE ----------------
	playlistRepo := interactionservice.NewPlaylistRepository(db)
	playlistSvc := interactionservice.NewPlaylistService(playlistRepo)
	playlistHandler := interactionservice.NewPlaylistHandler(playlistSvc)

	// ---------------- PORT ----------------
	port, err := configs.LoadPort()
	if err != nil {
		panic("Failed to load port configuration: " + err.Error())
	}

	// ---------------- ROUTER ----------------
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
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "interaction-service",
		})
	})

	r.POST("/swipe", h.Swipe)
	r.GET("/swipes/:swipeId", h.GetSwipe)
	r.GET("/users/:userId/swipes", h.GetSwipeByUser)

	r.POST("/playlists", playlistHandler.CreatePlaylist)
	r.GET("/playlists/:id", playlistHandler.GetPlaylist)
	r.GET("/users/:userId/playlists", playlistHandler.GetUserPlaylists)
	r.PATCH("/playlists/:id", playlistHandler.UpdatePlaylistName)
	r.DELETE("/playlists/:id", playlistHandler.DeletePlaylist)

	r.POST("/playlists/:id/tracks/:trackId", playlistHandler.AddTrack)
	r.DELETE("/playlists/:id/tracks/:trackId", playlistHandler.RemoveTrack)

	r.GET("/users/:userId/recommended", playlistHandler.GetRecommended)

	log.Println("Server running on", port)
	r.Run(":" + port)
}
