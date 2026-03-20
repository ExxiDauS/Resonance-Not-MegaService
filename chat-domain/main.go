package main

import (
	chatservice "chat-domain/chat-service"
	configs "chat-domain/configs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"time"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)
	
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chat_service_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "chat_service_request_duration_seconds",
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
	// Initialize Redis
	rdcfg, err := configs.LoadRedisConfig()

	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     rdcfg.Address,
		Username: rdcfg.Username,
		Password: rdcfg.Password,
	})

	redisService := chatservice.NewRedisService(rdb)
	hub := chatservice.NewHub(redisService, rdb)
	handler := chatservice.NewChatHandler(hub, redisService)

	r := gin.Default()
	r.Use(cors.Default())

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
			"service": "chat-service",
		})
	})

	// The frontend connects to this endpoint to join a room
	// Example: ws://localhost:8081/ws/chat/room123?user_id=userA
	r.GET("/ws/chat/:room_id", handler.HandleConnections)

	portcfg, err := configs.LoadPortConfig()
	if err != nil {
		panic(err)
	}

	r.Run(":" + portcfg.Port)
}
