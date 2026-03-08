package interactionservice

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SwipeHandler struct {
	service SwipeService
}

func NewSwipeHandler(s SwipeService) *SwipeHandler {
	return &SwipeHandler{service: s}
}

func (h *SwipeHandler) Swipe(c *gin.Context) {
	var req SwipeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.Swipe(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to swipe",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "swiped successfully",
	})
}

func (h *SwipeHandler) GetSwipe(c *gin.Context) {
	swipeID := c.Param("swipeId")

	swipe, err := h.service.GetSwipe(c.Request.Context(), swipeID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "swipe not found",
				"error":   nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get swipe",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    swipe,
		"message": "swipe retrieved successfully",
	})
}

func (h *SwipeHandler) GetSwipeByUser(c *gin.Context) {
	userID := c.Param("userId")
	swipes, err := h.service.GetUserSwipes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to fetch swipes",
			"error":   err.Error(),
		})
		return
	}
	if len(swipes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "no swipes found for user",
			"error":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    swipes,
		"message": "user swipes retrieved successfully",
	})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error()})
		return
	}

	playlist, err := h.service.CreatePlaylist(req.UserID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create playlist",
			"error":   err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    playlist,
		"message": "playlist created successfully",
	})
}

// READ
func (h *PlaylistHandler) GetPlaylist(c *gin.Context) {
	id := c.Param("id")

	playlist, err := h.service.GetPlaylist(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "playlist not found",
			"error":   err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    playlist,
		"message": "playlist retrieved successfully",
	})
}

// READ
func (h *PlaylistHandler) GetUserPlaylists(c *gin.Context) {
	userID := c.Param("userId")

	playlists, err := h.service.GetUserPlaylists(userID)
	if err != nil {
		if err.Error() == "no playlists found for user" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "no playlists found for user",
				"error":   nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to fetch playlists",
			"error":   err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    playlists,
		"message": "user playlists retrieved successfully",
	})
}

// UPDATE
func (h *PlaylistHandler) UpdatePlaylistName(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error()})
		return
	}

	if err := h.service.UpdatePlaylistName(id, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update playlist name",
			"error":   err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "updated",
		"data":    req})
}

// DELETE
func (h *PlaylistHandler) DeletePlaylist(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeletePlaylist(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to delete playlist",
			"error":   err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "deleted",
	})
}

// ADD TRACK
func (h *PlaylistHandler) AddTrack(c *gin.Context) {
	playlistID := c.Param("id")
	trackID := c.Param("trackId")

	if err := h.service.AddTrack(playlistID, trackID); err != nil {
		if err.Error() == fmt.Sprintf("playlist with id %s not found", playlistID) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "playlist not found",
				"error":   nil})
			return
		}
		if err.Error() == "track already exists in playlist" {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "track already exists in playlist",
				"error":   nil})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to add track to playlist",
			"error":   err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "track added"})
}

// REMOVE TRACK
func (h *PlaylistHandler) RemoveTrack(c *gin.Context) {

	playlistID := c.Param("id")
	trackID := c.Param("trackId")

	err := h.service.RemoveTrack(playlistID, trackID)

	if err != nil {
		if err.Error() == "playlist not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "playlist not found",
				"error":   nil,
			})
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "track not found in playlist",
				"error":   nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to remove track from playlist",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "track removed",
	})
}

func (h *PlaylistHandler) GetRecommended(c *gin.Context) {

	userID := c.Param("userId")

	playlist, err := h.service.GetRecommended(userID)
	if err != nil {

		if err.Error() == "no liked tracks found for user" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "user has no liked tracks yet",
				"error":   nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get recommended playlist",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    playlist,
		"message": "recommended playlist retrieved successfully",
	})
}
