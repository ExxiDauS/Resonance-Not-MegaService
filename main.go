package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func newReverseProxy(target string) (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Customize the director to preserve the Host header
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[Gateway] Proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"success":false,"message":"service unavailable"}`))
	}

	return proxy, nil
}

func proxyHandler(proxy *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	// Service URLs (defaults for Docker Compose networking)
	authURL := getEnv("AUTH_SERVICE_URL", "http://auth-domain:8080")
	userURL := getEnv("USER_SERVICE_URL", "http://user-domain:8081")
	chatURL := getEnv("CHAT_SERVICE_URL", "http://chat-domain:8082")
	musicURL := getEnv("MUSIC_SERVICE_URL", "http://music-domain:8083")

	// Create reverse proxies
	authProxy, err := newReverseProxy(authURL)
	if err != nil {
		log.Fatalf("Failed to create auth proxy: %v", err)
	}

	userProxy, err := newReverseProxy(userURL)
	if err != nil {
		log.Fatalf("Failed to create user proxy: %v", err)
	}

	chatProxy, err := newReverseProxy(chatURL)
	if err != nil {
		log.Fatalf("Failed to create chat proxy: %v", err)
	}

	musicProxy, err := newReverseProxy(musicURL)
	if err != nil {
		log.Fatalf("Failed to create music proxy: %v", err)
	}

	r := gin.Default()

	// Global CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	// --- Auth Domain ---
	// Routes: POST /register, POST /login, POST /logout
	auth := r.Group("/api/auth")
	{
		auth.Any("/*path", func(c *gin.Context) {
			// Strip /api/auth prefix before forwarding
			c.Request.URL.Path = c.Param("path")
			authProxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	// --- User Domain ---
	// Routes: GET /me, POST /profiles, GET /profiles, GET /profiles/:id, PATCH /profiles/:id, DELETE /profiles/:id
	users := r.Group("/api/users")
	{
		users.Any("/*path", func(c *gin.Context) {
			// Strip /api/users prefix before forwarding
			c.Request.URL.Path = c.Param("path")
			userProxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	// --- Music Domain ---
	// Routes: GET /tracks, GET /tracks/random, GET /tracks/:id, DELETE /tracks/:id
	music := r.Group("/api/music")
	{
		music.Any("/*path", func(c *gin.Context) {
			// Strip /api/music prefix before forwarding
			c.Request.URL.Path = c.Param("path")
			musicProxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	// --- Chat Domain (WebSocket) ---
	// Routes: GET /ws/chat/:room_id
	r.Any("/ws/*path", func(c *gin.Context) {
		// Forward WebSocket requests directly (no prefix stripping needed)
		c.Request.URL.Path = "/ws" + c.Param("path")
		chatProxy.ServeHTTP(c.Writer, c.Request)
	})

	port := getEnv("GATEWAY_PORT", "8000")
	log.Printf("🚀 API Gateway starting on port %s", port)
	log.Printf("   Auth   → %s", authURL)
	log.Printf("   User   → %s", userURL)
	log.Printf("   Chat   → %s", chatURL)
	log.Printf("   Music  → %s", musicURL)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start gateway: %v", err)
	}
}
