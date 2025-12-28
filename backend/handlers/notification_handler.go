package handlers

import (
	"net/http"
	"sirine-go/backend/models"
	"sirine-go/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NotificationHandler merupakan handler untuk notification endpoints
// yang bertujuan untuk serving notification operations ke frontend
type NotificationHandler struct {
	service *services.NotificationService
}

// NewNotificationHandler membuat instance baru NotificationHandler
func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// GetUserNotifications mengambil semua notifikasi user
// GET /api/notifications?unread_only=true
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	// Get user dari context (set oleh auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	// Parse query param untuk filter unread only
	unreadOnly := c.Query("unread_only") == "true"

	notifications, err := h.service.GetUserNotifications(userID.(uint64), unreadOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil notifikasi",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notifikasi berhasil diambil",
		"data":    notifications,
	})
}

// GetUnreadCount mengambil jumlah notifikasi yang belum dibaca
// GET /api/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	count, err := h.service.GetUnreadCount(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil jumlah notifikasi",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"count": count,
		},
	})
}

// MarkAsRead menandai satu notifikasi sebagai sudah dibaca
// PUT /api/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	// Parse notification ID dari URL param
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID notifikasi tidak valid",
		})
		return
	}

	err = h.service.MarkAsRead(notificationID, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menandai notifikasi sebagai dibaca",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notifikasi berhasil ditandai sebagai dibaca",
	})
}

// MarkAllAsRead menandai semua notifikasi sebagai sudah dibaca
// PUT /api/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	err := h.service.MarkAllAsRead(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menandai semua notifikasi sebagai dibaca",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Semua notifikasi berhasil ditandai sebagai dibaca",
	})
}

// GetRecentNotifications mengambil N notifikasi terbaru
// GET /api/notifications/recent?limit=5
func (h *NotificationHandler) GetRecentNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	// Parse limit dari query param, default 5
	limit := 5
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	notifications, err := h.service.GetRecentNotifications(userID.(uint64), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil notifikasi terbaru",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    notifications,
	})
}

// CreateNotification membuat notifikasi baru (untuk testing atau admin)
// POST /api/notifications
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req struct {
		UserID  uint64                     `json:"user_id" binding:"required"`
		Title   string                     `json:"title" binding:"required,max=255"`
		Message string                     `json:"message" binding:"required"`
		Type    models.NotificationType `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	notification, err := h.service.CreateNotification(req.UserID, req.Title, req.Message, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat notifikasi",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Notifikasi berhasil dibuat",
		"data":    notification,
	})
}

// DeleteNotification menghapus notifikasi
// DELETE /api/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID notifikasi tidak valid",
		})
		return
	}

	err = h.service.DeleteNotification(notificationID, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus notifikasi",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notifikasi berhasil dihapus",
	})
}
