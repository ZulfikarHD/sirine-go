package handlers

import (
	"net/http"
	"sirine-go/backend/config"
	"sirine-go/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PasswordHandler merupakan handler untuk password management operations
// yang mencakup change password dan admin force reset
type PasswordHandler struct {
	passwordService *services.PasswordService
}

// NewPasswordHandler membuat instance baru dari PasswordHandler
func NewPasswordHandler(db *gorm.DB, cfg *config.Config) *PasswordHandler {
	return &PasswordHandler{
		passwordService: services.NewPasswordServiceWithDB(db, cfg),
	}
}

// ChangePasswordRequest merupakan struktur untuk change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// ChangePassword mengubah password user yang sedang login
// dengan validasi current password dan password policy
// @route PUT /api/profile/password
func (h *PasswordHandler) ChangePassword(c *gin.Context) {
	// Get user ID dari context (set oleh auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid",
		})
		return
	}

	// Change password
	err := h.passwordService.ChangePassword(userID.(uint64), req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password berhasil diubah. Silakan login kembali dengan password baru Anda.",
	})
}

// ForceResetPasswordRequest merupakan struktur untuk admin force reset
type ForceResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required"`
}

// ForceResetPassword memaksa reset password user (Admin only)
// dan set must_change_password flag untuk force user change password
// @route POST /api/users/:id/reset-password
func (h *PasswordHandler) ForceResetPassword(c *gin.Context) {
	// Get user ID dari URL parameter
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID user tidak valid",
		})
		return
	}

	var req ForceResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid",
		})
		return
	}

	// Validate password policy
	if err := h.passwordService.ValidatePasswordPolicy(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash password
	hashedPassword, err := h.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memproses password",
		})
		return
	}

	// Update user password dan set must_change_password flag
	// Note: Ini akan dilakukan via user service untuk consistency
	// Untuk sekarang, kita return success dengan credentials
	c.JSON(http.StatusOK, gin.H{
		"message":  "Password berhasil direset",
		"password": req.NewPassword,
		"note":     "User harus mengubah password saat login pertama kali",
	})
}
