package main

import (
	"log"
	audiosservice "music-domain/audios-service"
	"music-domain/configs"
	"music-domain/infrastructures/databases"
	"music-domain/infrastructures/soundcloud_client"
	"music-domain/infrastructures/spotify_client"
	tracksservice "music-domain/tracks-service"

	"github.com/gin-gonic/gin"

	"time"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "music_service_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "music_service_request_duration_seconds",
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
	cfg, err := configs.LoadPortConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := databases.NewPostgresDatabaseClient()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&tracksservice.Track{})

	spotifyClient, err := spotify_client.NewSpotifyClient()
	if err != nil {
		log.Fatalf("Failed to create Spotify client: %v", err)
	}

	souncloudClient := soundcloud_client.NewSoundCloudClient()
	if souncloudClient == nil {
		log.Fatalf("Failed to create SoundCloud client")
	}

	audioRepo := audiosservice.NewSoundCloudAudioRepo(souncloudClient)
	audioService := audiosservice.NewAudioService(audioRepo)

	trackRepo := tracksservice.NewRepository(db)
	trackService := tracksservice.NewService(trackRepo, spotifyClient, audioService)
	trackHandler := tracksservice.NewHandler(trackService)

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
			"service": "music-service",
		})
	})

	r.GET("/tracks", trackHandler.GetAllTracks)
	r.GET("/tracks/random", trackHandler.GetRandomTracks)
	r.GET("/tracks/:id", trackHandler.GetTrackByID)
	r.DELETE("/tracks/:id", trackHandler.DeleteTrackByID)

	r.Run(":" + cfg.Port)
}
