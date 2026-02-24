package authservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Register(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
		})
		return
	}

	user, err := h.service.Register(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email already in use",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"user":    user,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
		})
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid email or password",
		})
		return
	}

	c.SetCookie("auth_token", token, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
	})
}
