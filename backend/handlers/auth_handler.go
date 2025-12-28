package handlers

import (
	"net/http"
	"sirine-go/backend/config"
	"sirine-go/backend/models"
	"sirine-go/backend/services"
	"sirine-go/backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler merupakan handler untuk authentication endpoints
// yang mencakup login, logout, refresh token, dan password reset
type AuthHandler struct {
	authService     *services.AuthService
	passwordService *services.PasswordService
}

// NewAuthHandler membuat instance baru dari AuthHandler
func NewAuthHandler(authService *services.AuthService, db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService:     authService,
		passwordService: services.NewPasswordServiceWithDB(db, cfg),
	}
}

// Login handles POST /api/auth/login
// @Summary User login dengan NIP atau Email dan password
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "Login credentials (NIP atau Email)"
// @Success 200 {object} services.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		// Return validation errors dengan pesan dalam Bahasa Indonesia
		c.JSON(http.StatusBadRequest, utils.NewValidationErrorResponse(err))
		return
	}

	// Get client info
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Perform login
	response, err := h.authService.Login(req, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data":    response,
	})
}

// Logout handles POST /api/auth/logout
// @Summary User logout dan invalidate token
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user dari context (diset oleh AuthMiddleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	user := userInterface.(*models.User)
	
	// Get token dari header
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Get client info
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Perform logout
	err := h.authService.Logout(user.ID, token, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Logout gagal",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout berhasil",
	})
}

// GetCurrentUser handles GET /api/auth/me
// @Summary Get current authenticated user info
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.SafeUser
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get user dari context (diset oleh AuthMiddleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	user := userInterface.(*models.User)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data user berhasil diambil",
		"data":    user.ToSafeUser(),
	})
}

// RefreshToken handles POST /api/auth/refresh
// @Summary Refresh JWT token menggunakan refresh token
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} services.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// Return validation errors dengan pesan dalam Bahasa Indonesia
		c.JSON(http.StatusBadRequest, utils.NewValidationErrorResponse(err))
		return
	}

	// Refresh the token
	response, err := h.authService.RefreshAuthToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token berhasil di-refresh",
		"data":    response,
	})
}

// ForgotPasswordRequest merupakan struktur untuk forgot password request
type ForgotPasswordRequest struct {
	NIPOrEmail string `json:"nip_or_email" binding:"required"`
}

// ForgotPassword handles POST /api/auth/forgot-password
// @Summary Request password reset link via email
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "NIP atau Email"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewValidationErrorResponse(err))
		return
	}

	// Request password reset (akan send email jika user ditemukan)
	err := h.passwordService.RequestPasswordReset(req.NIPOrEmail)
	if err != nil {
		// Log error tapi tetap return success untuk prevent enumeration
		// (jangan kasih tau user apakah email exist atau tidak)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Jika NIP/Email terdaftar, link reset password telah dikirim ke email Anda. Silakan cek inbox atau spam folder.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jika NIP/Email terdaftar, link reset password telah dikirim ke email Anda. Silakan cek inbox atau spam folder.",
	})
}

// ResetPasswordRequest merupakan struktur untuk reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ResetPassword handles POST /api/auth/reset-password
// @Summary Reset password menggunakan token dari email
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Token dan password baru"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewValidationErrorResponse(err))
		return
	}

	// Reset password dengan token
	err := h.passwordService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password berhasil direset. Silakan login dengan password baru Anda.",
	})
}

// ErrorResponse untuk swagger documentation
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SuccessResponse untuk swagger documentation
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
