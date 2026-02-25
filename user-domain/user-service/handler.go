package userservice

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

func (h *Handler) CreateUserProfile(c *gin.Context) {
	var profile Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request payload",
		})
		return
	}
	if err := h.service.CreateUserProfile(&profile); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			c.JSON(409, gin.H{"success": false, "message": "Email already exists"})
		} else {
			c.JSON(500, gin.H{"success": false, "message": "Failed to create user profile", "error": err.Error()})
		}
		return
	}
	c.JSON(201, gin.H{
		"success": true,
		"profile": profile,
	})
}

func (h *Handler) GetAllUserProfiles(c *gin.Context) {
	profiles, err := h.service.GetAllUserProfiles()
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to retrieve user profiles",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success":  true,
		"profiles": profiles,
	})
}

func (h *Handler) GetUserProfileByID(c *gin.Context) {
	userID := c.Param("id")
	profile, err := h.service.GetUserProfileByID(userID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"success": false, "message": "User profile not found"})
		} else {
			c.JSON(500, gin.H{"success": false, "message": "Database error", "error": err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"profile": profile,
	})
}

func (h *Handler) UpdateUserProfile(c *gin.Context) {
	userID := c.Param("id")
	authenticatedUserID := c.GetString("user_id") // From JWT

	// Only allow users to update their own profile
	if userID != authenticatedUserID {
		c.JSON(403, gin.H{
			"success": false,
			"message": "You can only update your own profile",
		})
		return
	}

	// Parse the request body into a map to allow partial updates
	var updates *UpdateProfileInput
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request payload",
		})
		return
	}

	// Validate that at least one field is provided
	if updates.Email == nil && updates.DisplayName == nil && updates.Bio == nil && updates.AvatarURL == nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "At least one field must be provided for update",
		})
		return
	}

	updatedProfile, err := h.service.UpdateUserProfile(userID, updates)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"success": false,
				"message": "User profile not found",
			})
		} else {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Failed to update user profile",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "User profile updated successfully",
		"profile": updatedProfile,
	})
}

func (h *Handler) DeleteUserProfile(c *gin.Context) {
	userID := c.Param("id")
	authenticatedUserID := c.GetString("user_id") // From JWT

	// Only allow users to update their own profile
	if userID != authenticatedUserID {
		c.JSON(403, gin.H{
			"success": false,
			"message": "You can only update your own profile",
		})
		return
	}

	if err := h.service.DeleteUserProfile(userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"success": false,
				"message": "User profile not found",
			})
			return
		} else {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Failed to delete user profile",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "User profile deleted successfully",
	})
}
