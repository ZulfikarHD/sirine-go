package middleware

import (
	"net/http"
	"sirine-go/backend/models"

	"github.com/gin-gonic/gin"
)

// RequireRole merupakan middleware untuk role-based access control
// yang memastikan user memiliki salah satu role yang diizinkan
func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user dari context (harus sudah diset oleh AuthMiddleware)
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User tidak terautentikasi",
			})
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error mengambil data user",
			})
			c.Abort()
			return
		}

		// Check apakah user memiliki salah satu role yang diizinkan
		hasPermission := false
		for _, allowedRole := range allowedRoles {
			if user.Role == allowedRole {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Anda tidak memiliki akses untuk resource ini",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin merupakan middleware shorthand untuk admin dan manager only
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin, models.RoleManager)
}

// RequireDepartment merupakan middleware untuk department-based access control
func RequireDepartment(allowedDepartments ...models.Department) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User tidak terautentikasi",
			})
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error mengambil data user",
			})
			c.Abort()
			return
		}

		// Admin dan Manager dapat akses semua department
		if user.IsAdmin() {
			c.Next()
			return
		}

		// Check department permission
		hasPermission := false
		for _, allowedDept := range allowedDepartments {
			if user.Department == allowedDept {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Anda tidak memiliki akses untuk departemen ini",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
