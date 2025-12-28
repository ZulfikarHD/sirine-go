package handlers

import (
	"net/http"
	"sirine-go/backend/models"
	"sirine-go/backend/services"

	"github.com/gin-gonic/gin"
)

// ProfileHandler merupakan handler untuk profile management (self-service)
type ProfileHandler struct {
	userService *services.UserService
}

// NewProfileHandler membuat instance baru dari ProfileHandler
func NewProfileHandler(userService *services.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

// GetProfile mengambil profile user yang sedang login
// @route GET /api/profile
// @access Protected (All authenticated users)
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Get current user dari context (set oleh auth middleware)
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User tidak terautentikasi",
		})
		return
	}

	user := currentUser.(*models.User)

	// Get full user data dari database untuk ensure latest data
	userDetail, err := h.userService.GetUserByID(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userDetail,
	})
}

// UpdateProfile mengupdate profile user sendiri
// User hanya bisa update: full_name, email, phone
// Tidak bisa update: nip, role, department, status
// @route PUT /api/profile
// @access Protected (All authenticated users)
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// Get current user dari context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User tidak terautentikasi",
		})
		return
	}

	user := currentUser.(*models.User)

	// Get user before update untuk activity log
	userBefore, _ := h.userService.GetUserByID(user.ID)

	var req services.UpdateProfileRequest

	// Bind dan validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data profile tidak valid",
			"details": err.Error(),
		})
		return
	}

	// Update profile
	updatedUser, err := h.userService.UpdateProfile(user.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Log activity
	c.Set("activity_action", models.ActionUpdate)
	c.Set("activity_entity_type", "users")
	c.Set("activity_entity_id", user.ID)
	c.Set("activity_changes_before", userBefore)
	c.Set("activity_changes_after", updatedUser)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile berhasil diupdate",
		"data":    updatedUser,
	})
}
