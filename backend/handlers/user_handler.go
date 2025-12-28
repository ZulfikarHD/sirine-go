package handlers

import (
	"net/http"
	"sirine-go/backend/models"
	"sirine-go/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler merupakan handler untuk user management operations (Admin)
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler membuat instance baru dari UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetAllUsers mengambil list users dengan filters dan pagination
// @route GET /api/users
// @access Admin, Manager
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var filters services.UserFilters

	// Bind query parameters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parameter filter tidak valid",
			"details": err.Error(),
		})
		return
	}

	// Get users dari service
	response, err := h.userService.GetAllUsers(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data users",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetUserByID mengambil detail user berdasarkan ID
// @route GET /api/users/:id
// @access Admin, Manager
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Parse ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID user tidak valid",
		})
		return
	}

	// Get user dari service
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// CreateUser membuat user baru dengan auto-generated password
// @route POST /api/users
// @access Admin only
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req services.CreateUserRequest

	// Bind dan validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data user tidak valid",
			"details": err.Error(),
		})
		return
	}

	// Create user
	response, err := h.userService.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Log activity (akan ditambahkan oleh activity_logger middleware)
	c.Set("activity_action", models.ActionCreate)
	c.Set("activity_entity_type", "users")
	c.Set("activity_entity_id", response.User.ID)
	c.Set("activity_changes_after", response.User)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User berhasil dibuat",
		"data":    response,
	})
}

// UpdateUser mengupdate data user
// @route PUT /api/users/:id
// @access Admin only
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Parse ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID user tidak valid",
		})
		return
	}

	// Get user before update untuk activity log
	userBefore, _ := h.userService.GetUserByID(id)

	var req services.UpdateUserRequest

	// Bind dan validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data update tidak valid",
			"details": err.Error(),
		})
		return
	}

	// Update user
	user, err := h.userService.UpdateUser(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Log activity
	c.Set("activity_action", models.ActionUpdate)
	c.Set("activity_entity_type", "users")
	c.Set("activity_entity_id", id)
	c.Set("activity_changes_before", userBefore)
	c.Set("activity_changes_after", user)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User berhasil diupdate",
		"data":    user,
	})
}

// DeleteUser melakukan soft delete pada user
// @route DELETE /api/users/:id
// @access Admin only
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Parse ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID user tidak valid",
		})
		return
	}

	// Get current user untuk prevent self-delete
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User tidak terautentikasi",
		})
		return
	}

	user := currentUser.(*models.User)
	if user.ID == id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Anda tidak dapat menghapus akun sendiri",
		})
		return
	}

	// Get user before delete untuk activity log
	userBefore, _ := h.userService.GetUserByID(id)

	// Delete user
	if err := h.userService.DeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Log activity
	c.Set("activity_action", models.ActionDelete)
	c.Set("activity_entity_type", "users")
	c.Set("activity_entity_id", id)
	c.Set("activity_changes_before", userBefore)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User berhasil dihapus",
	})
}

// SearchUsers melakukan search users
// @route GET /api/users/search
// @access Admin, Manager
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query parameter 'q' diperlukan",
		})
		return
	}

	users, err := h.userService.SearchUsers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal melakukan search",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

// BulkDeleteUsers melakukan soft delete pada multiple users
// @route POST /api/users/bulk-delete
// @access Admin only
func (h *UserHandler) BulkDeleteUsers(c *gin.Context) {
	var req services.BulkDeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request tidak valid",
			"details": err.Error(),
		})
		return
	}

	// Get current user ID
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User tidak terautentikasi",
		})
		return
	}

	user := currentUser.(*models.User)

	// Perform bulk delete
	affected, err := h.userService.BulkDeleteUsers(req.UserIDs, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Bulk delete berhasil",
		"data": gin.H{
			"affected_count": affected,
		},
	})
}

// BulkUpdateStatus mengupdate status multiple users
// @route POST /api/users/bulk-update-status
// @access Admin only
func (h *UserHandler) BulkUpdateStatus(c *gin.Context) {
	var req services.BulkUpdateStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request tidak valid",
			"details": err.Error(),
		})
		return
	}

	// Get current user ID
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User tidak terautentikasi",
		})
		return
	}

	user := currentUser.(*models.User)

	// Perform bulk update
	affected, err := h.userService.BulkUpdateStatus(req.UserIDs, req.Status, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Bulk update status berhasil",
		"data": gin.H{
			"affected_count": affected,
		},
	})
}
