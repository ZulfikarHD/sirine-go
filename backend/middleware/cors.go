package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS merupakan middleware untuk handle Cross-Origin Resource Sharing
// yang mengizinkan frontend dari berbagai origin untuk mengakses API
func CORS() gin.HandlerFunc {
	// Ambil allowed origins dari environment variable
	// Format: CORS_ORIGINS="http://localhost:5173,http://localhost:8080"
	allowedOrigins := []string{"http://localhost:5173", "http://localhost:8080"}
	
	if corsOrigins := os.Getenv("CORS_ORIGINS"); corsOrigins != "" {
		allowedOrigins = strings.Split(corsOrigins, ",")
		// Trim whitespace dari setiap origin
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}
	
	// Untuk development, allow all origins jika CORS_ALLOW_ALL=true
	if os.Getenv("CORS_ALLOW_ALL") == "true" {
		return cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: false, // Tidak bisa true jika AllowAllOrigins true
			MaxAge:           12 * time.Hour,
		})
	}
	
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
