package matchesservice

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateMatch(c *gin.Context) {
	var match Match
	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request payload",
		})
		return
	}

	if err := h.service.CreateMatch(&match); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			c.JSON(409, gin.H{
				"success": false,
				"message": "Match already exists",
			})
		} else {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Failed to create match",
			})
		}
	} else {
		c.JSON(201, gin.H{
			"success": true,
			"message": "Match created successfully",
			"match":   match,
		})
	}
}

func (h *Handler) GetMatchesByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	matches, err := h.service.GetMatchesByUserID(userID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"success": false,
				"message": "No matches found for the user",
			})
		} else {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Failed to retrieve matches",
				"error":   err.Error(),
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"success": true,
		"matches": matches,
	})
}

func (h *Handler) GetMatchByID(c *gin.Context) {
	matchID := c.Param("match_id")
	match, err := h.service.GetMatchByID(matchID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"success": false,
				"message": "Match not found",
			})
		} else {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Failed to retrieve match",
				"error":   err.Error(),
			})
		}
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"match":   match,
	})
}

func (h *Handler) GetAllMatches(c *gin.Context) {
	matches, err := h.service.GetAllMatches()
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to retrieve matches",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"matches": matches,
	})
}

func (h *Handler) DeleteMatch(c *gin.Context) {
	matchID := c.Param("match_id")
	if err := h.service.DeleteMatch(matchID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"success": false,
				"message": "Match not found",
			})
		} else {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Failed to delete match",
				"error":   err.Error(),
			})
		}
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "Match deleted successfully",
	})
}
