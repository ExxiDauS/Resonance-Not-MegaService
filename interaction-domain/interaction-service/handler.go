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

func (h *SwipeHandler) GetSwipe(c *gin.Context) {
	swipeID := c.Param("swipeId")
	swipes, err := h.service.GetSwipe(c.Request.Context(), swipeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, swipes)
}

func (h *SwipeHandler) GetSwipeByUser(c *gin.Context) {
	userID := c.Param("userId")
	swipes, err := h.service.GetUserSwipes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, swipes)
}

// PlaylistHandler is a placeholder for future playlist-related endpoints
type PlaylistHandler struct {
	service PlaylistService
}

func NewPlaylistHandler(s PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{s}
}

// CREATE
func (h *PlaylistHandler) CreatePlaylist(c *gin.Context) {
	var req struct {
		UserID string `json:"userId" binding:"required"`
		Name   string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playlist, err := h.service.CreatePlaylist(req.UserID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, playlist)
}

// READ
func (h *PlaylistHandler) GetPlaylist(c *gin.Context) {
	id := c.Param("id")

	playlist, err := h.service.GetPlaylist(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, playlist)
}

// READ
func (h *PlaylistHandler) GetUserPlaylists(c *gin.Context) {
	userID := c.Param("userId")

	playlists, err := h.service.GetUserPlaylists(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, playlists)
}

// UPDATE
func (h *PlaylistHandler) UpdatePlaylistName(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdatePlaylistName(id, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// DELETE
func (h *PlaylistHandler) DeletePlaylist(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeletePlaylist(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ADD TRACK
func (h *PlaylistHandler) AddTrack(c *gin.Context) {
	playlistID := c.Param("id")
	trackID := c.Param("trackId")

	if err := h.service.AddTrack(playlistID, trackID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "track added"})
}

// REMOVE TRACK
func (h *PlaylistHandler) RemoveTrack(c *gin.Context) {

	playlistID := c.Param("id")
	trackID := c.Param("trackId")

	if err := h.service.RemoveTrack(playlistID, trackID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "track removed"})
}

func (h *PlaylistHandler) GetRecommended(c *gin.Context) {

	userID := c.Param("userId")

	playlist, err := h.service.GetRecommended(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, playlist)
}
