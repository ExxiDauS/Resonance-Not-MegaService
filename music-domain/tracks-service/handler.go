package tracksservice

import (
	"net/http"
	"strconv"

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
	c.JSON(http.StatusOK, tracks)
}

func (h *Handler) GetTrackByID(c *gin.Context) {
	id := c.Param("id")
	track, err := h.service.GetTrackByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if track == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}
	c.JSON(http.StatusOK, track)
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func (h *Handler) GetAllTracks(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	name := c.Query("name")
	genre := c.Query("genre")
	artists := c.Query("artists")
	pageInt := parseInt(page)
	limitInt := parseInt(limit)
	allTracks, err := h.service.GetAllTracks(c.Request.Context(), pageInt, limitInt, &name, &genre, &artists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, allTracks)
}
