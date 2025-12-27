package handlers

import (
	"net/http"
	"sirine-go/backend/models"
	"sirine-go/backend/services"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthHandler merupakan handler untuk authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler membuat instance baru dari AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles POST /api/auth/login
// @Summary User login dengan NIP dan password
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "Login credentials"
// @Success 200 {object} services.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data yang dikirim tidak valid",
			"error":   err.Error(),
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Refresh token harus disertakan",
			"error":   err.Error(),
		})
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
