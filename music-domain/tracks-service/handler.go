package tracksservice

import (
	"net/http"
	"strings"
	"utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetRandomTracks(c *gin.Context) {
	tracks, err := h.service.getRandomTrackFromSpotify(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tracks,
	})
}

func (h *Handler) GetTrackByID(c *gin.Context) {
	id := c.Param("id")
	track, err := h.service.GetTrackByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if track == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Track not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    track,
	})
}

func (h *Handler) GetAllTracks(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	name := c.Query("name")
	genre := c.Query("genre")
	artists := c.Query("artists")
	pageInt := utils.ParseInt(page)
	limitInt := utils.ParseInt(limit)
	allTracks, err := h.service.GetAllTracks(c.Request.Context(), pageInt, limitInt, &name, &genre, &artists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    allTracks,
	})
}

func (h *Handler) DeleteTrackByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "track ID is required"})
		return
	}
	err := h.service.DeleteTrackByID(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "track not found") {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Track not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Track deleted successfully",
	})
}
