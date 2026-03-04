package interactionservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SwipeHandler struct {
	service SwipeService
}

func NewSwipeHandler(s SwipeService) *SwipeHandler {
	return &SwipeHandler{service: s}
}

func (h *SwipeHandler) GetRandomTrack(c *gin.Context) {
	track, err := h.service.GetRandomTrack(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, track)
}

func (h *SwipeHandler) Swipe(c *gin.Context) {
	var req SwipeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Swipe(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "swiped successfully"})
}
