package handlers

import (
	"net/http"
	"sirine-go/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AchievementHandler merupakan handler untuk achievement endpoints
// yang mengelola operations terkait gamification achievements
type AchievementHandler struct {
	achievementService *services.AchievementService
}

// NewAchievementHandler membuat instance baru dari AchievementHandler
func NewAchievementHandler(achievementService *services.AchievementService) *AchievementHandler {
	return &AchievementHandler{
		achievementService: achievementService,
	}
}

// GetAllAchievements mengambil semua available achievements
// GET /api/achievements
func (h *AchievementHandler) GetAllAchievements(c *gin.Context) {
	achievements, err := h.achievementService.GetAllAchievements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data achievements",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data achievements",
		"data":    achievements,
	})
}

// GetUserAchievements mengambil achievements user dengan unlock status
// GET /api/profile/achievements
func (h *AchievementHandler) GetUserAchievements(c *gin.Context) {
	// Get user ID dari context (dari auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	achievements, err := h.achievementService.GetUserAchievements(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil achievements user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil achievements user",
		"data":    achievements,
	})
}

// GetUserStats mengambil statistik gamification user
// GET /api/profile/stats
func (h *AchievementHandler) GetUserStats(c *gin.Context) {
	// Get user ID dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	stats, err := h.achievementService.GetUserStats(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil statistik user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil statistik user",
		"data":    stats,
	})
}

// AwardAchievementRequest merupakan request body untuk award achievement
type AwardAchievementRequest struct {
	UserID          uint64 `json:"user_id" binding:"required"`
	AchievementCode string `json:"achievement_code" binding:"required"`
}

// AwardAchievement memberikan achievement ke user (Admin only)
// POST /api/admin/achievements/award
func (h *AchievementHandler) AwardAchievement(c *gin.Context) {
	var req AwardAchievementRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if err := h.achievementService.AwardAchievement(req.UserID, req.AchievementCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memberikan achievement",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Achievement berhasil diberikan",
	})
}

// GetAchievementsByUserID mengambil achievements user berdasarkan user ID (Admin only)
// GET /api/admin/users/:id/achievements
func (h *AchievementHandler) GetAchievementsByUserID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User ID tidak valid",
		})
		return
	}

	achievements, err := h.achievementService.GetUserAchievements(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil achievements user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil achievements user",
		"data":    achievements,
	})
}
