package handlers

import (
	"net/http"
	"sirine-go/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ActivityLogHandler merupakan handler untuk activity log endpoints
// yang bertujuan untuk serving audit trail data ke admin
type ActivityLogHandler struct {
	service *services.ActivityLogService
}

// NewActivityLogHandler membuat instance baru ActivityLogHandler
func NewActivityLogHandler(service *services.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{service: service}
}

// GetActivityLogs mengambil activity logs dengan filters dan pagination
// GET /api/admin/activity-logs
func (h *ActivityLogHandler) GetActivityLogs(c *gin.Context) {
	var filters services.ActivityLogFilters
	var pagination services.Pagination

	// Bind query params ke filters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter filter tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Bind query params ke pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter pagination tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Default pagination values
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 20
	}

	response, err := h.service.GetActivityLogs(filters, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil activity logs",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Activity logs berhasil diambil",
		"data":    response.Data,
		"meta": gin.H{
			"total":       response.Total,
			"page":        response.Page,
			"page_size":   response.PageSize,
			"total_pages": response.TotalPages,
		},
	})
}

// GetUserActivity mengambil activity logs untuk specific user
// GET /api/admin/activity-logs/user/:id
func (h *ActivityLogHandler) GetUserActivity(c *gin.Context) {
	// Parse user ID dari URL param
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID user tidak valid",
		})
		return
	}

	// Parse limit dari query param, default 50
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := h.service.GetUserActivity(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil activity logs user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Activity logs user berhasil diambil",
		"data":    logs,
	})
}

// GetActivityLogByID mengambil detail single activity log
// GET /api/admin/activity-logs/:id
func (h *ActivityLogHandler) GetActivityLogByID(c *gin.Context) {
	// Parse log ID dari URL param
	logID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID activity log tidak valid",
		})
		return
	}

	log, err := h.service.GetActivityLogByID(logID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil detail activity log",
			"error":   err.Error(),
		})
		return
	}

	// Parse changes untuk display
	changes, err := log.GetChanges()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal parsing changes data",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail activity log berhasil diambil",
		"data": gin.H{
			"id":          log.ID,
			"user_id":     log.UserID,
			"user":        log.User,
			"action":      log.Action,
			"entity_type": log.EntityType,
			"entity_id":   log.EntityID,
			"changes":     changes,
			"ip_address":  log.IPAddress,
			"user_agent":  log.UserAgent,
			"created_at":  log.CreatedAt,
		},
	})
}

// GetActivityStats mengambil statistik activity
// GET /api/admin/activity-logs/stats
func (h *ActivityLogHandler) GetActivityStats(c *gin.Context) {
	// Parse date range dari query params
	var req struct {
		StartDate string `form:"start_date"`
		EndDate   string `form:"end_date"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Default: last 7 days
	// Implementasi parsing date bisa ditambahkan sesuai kebutuhan
	// Untuk sekarang, return empty stats
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Activity stats berhasil diambil",
		"data": gin.H{
			"stats": map[string]int64{},
		},
	})
}

// GetMyActivity mengambil activity logs untuk current user
// GET /api/profile/activity
func (h *ActivityLogHandler) GetMyActivity(c *gin.Context) {
	// Get user dari context (set oleh auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	// Parse limit dari query param, default 20
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := h.service.GetUserActivity(userID.(uint64), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil riwayat aktivitas",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Riwayat aktivitas berhasil diambil",
		"data":    logs,
	})
}
