package middleware

import (
	"net/http"
	"sirine-go/backend/config"
	"sirine-go/backend/services"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthMiddleware merupakan middleware untuk validasi JWT token
// yang memastikan user terautentikasi sebelum mengakses protected routes
func AuthMiddleware(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	authService := services.NewAuthService(db, cfg)
	
	return func(c *gin.Context) {
		// Get token dari Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header tidak ditemukan",
			})
			c.Abort()
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Format authorization header tidak valid, gunakan: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token tidak boleh kosong",
			})
			c.Abort()
			return
		}

		// Validate token
		user, claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token tidak valid atau sudah expired",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Set user dan claims ke context untuk digunakan di handlers
		c.Set("user", user)
		c.Set("claims", claims)
		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)
		
		c.Next()
	}
}

// OptionalAuthMiddleware merupakan middleware untuk optional authentication
// yang tidak memblokir request jika token tidak ada, tapi akan set user jika ada
func OptionalAuthMiddleware(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	authService := services.NewAuthService(db, cfg)
	
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			user, claims, err := authService.ValidateToken(token)
			if err == nil {
				c.Set("user", user)
				c.Set("claims", claims)
				c.Set("user_id", user.ID)
				c.Set("user_role", user.Role)
			}
		}

		c.Next()
	}
}
