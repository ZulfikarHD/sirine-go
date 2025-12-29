package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"sirine-go/backend/middleware"
	"sirine-go/backend/models"
	"testing"

	"github.com/gin-gonic/gin"
)

// setupRouterWithUser membuat router dengan user context untuk testing
func setupRouterWithUser(user *models.User, requiredRoles ...models.UserRole) *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	
	// Middleware untuk inject user ke context (simulate AuthRequired)
	router.Use(func(c *gin.Context) {
		c.Set("user", user)
		c.Next()
	})
	
	// Apply role middleware
	router.Use(middleware.RequireRole(requiredRoles...))
	
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	return router
}

// TestRequireRole_AdminAccess memverifikasi admin dapat akses admin-only endpoint
func TestRequireRole_AdminAccess(t *testing.T) {
	user := &models.User{
		ID:         1,
		NIP:        "12345",
		FullName:   "Admin User",
		Role:       models.RoleAdmin,
		Department: models.DepartmentKhazwal,
		Status:     models.UserStatusActive,
	}

	router := setupRouterWithUser(user, models.RoleAdmin)

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusOK)
	}
}

// TestRequireRole_ManagerAccess memverifikasi manager dapat akses admin/manager endpoint
func TestRequireRole_ManagerAccess(t *testing.T) {
	user := &models.User{
		ID:         2,
		NIP:        "12346",
		FullName:   "Manager User",
		Role:       models.RoleManager,
		Department: models.DepartmentKhazwal,
		Status:     models.UserStatusActive,
	}

	router := setupRouterWithUser(user, models.RoleAdmin, models.RoleManager)

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusOK)
	}
}

// TestRequireRole_UnauthorizedAccess memverifikasi staff tidak dapat akses admin endpoint
func TestRequireRole_UnauthorizedAccess(t *testing.T) {
	user := &models.User{
		ID:         3,
		NIP:        "12347",
		FullName:   "Staff User",
		Role:       models.RoleStaffKhazwal,
		Department: models.DepartmentKhazwal,
		Status:     models.UserStatusActive,
	}

	router := setupRouterWithUser(user, models.RoleAdmin)

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusForbidden)
	}
}

// TestRequireRole_MultipleAllowedRoles memverifikasi multiple allowed roles
func TestRequireRole_MultipleAllowedRoles(t *testing.T) {
	tests := []struct {
		name         string
		userRole     models.UserRole
		allowedRoles []models.UserRole
		expectedCode int
	}{
		{
			name:         "Admin access admin/manager endpoint",
			userRole:     models.RoleAdmin,
			allowedRoles: []models.UserRole{models.RoleAdmin, models.RoleManager},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Manager access admin/manager endpoint",
			userRole:     models.RoleManager,
			allowedRoles: []models.UserRole{models.RoleAdmin, models.RoleManager},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Staff denied admin/manager endpoint",
			userRole:     models.RoleStaffKhazwal,
			allowedRoles: []models.UserRole{models.RoleAdmin, models.RoleManager},
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "Operator access operator/staff endpoint",
			userRole:     models.RoleOperatorCetak,
			allowedRoles: []models.UserRole{models.RoleOperatorCetak, models.RoleStaffKhazwal},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{
				ID:         1,
				NIP:        "12345",
				FullName:   "Test User",
				Role:       tt.userRole,
				Department: models.DepartmentKhazwal,
				Status:     models.UserStatusActive,
			}

			router := setupRouterWithUser(user, tt.allowedRoles...)

			req := httptest.NewRequest("GET", "/protected", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Status code = %d, expected %d", w.Code, tt.expectedCode)
			}
		})
	}
}

// TestRequireRole_NoUserInContext memverifikasi error jika user tidak ada di context
func TestRequireRole_NoUserInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	
	// No user injection middleware
	router.Use(middleware.RequireRole(models.RoleAdmin))
	
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusUnauthorized)
	}
}

// TestRequireRole_AllRoles memverifikasi semua role dapat akses jika tidak ada restriction
func TestRequireRole_AllRoles(t *testing.T) {
	allRoles := []models.UserRole{
		models.RoleAdmin,
		models.RoleManager,
		models.RoleStaffKhazwal,
		models.RoleOperatorCetak,
		models.RoleQCInspector,
		models.RoleVerifikator,
		models.RoleStaffKhazkhir,
	}

	for _, role := range allRoles {
		t.Run(string(role), func(t *testing.T) {
			user := &models.User{
				ID:         1,
				NIP:        "12345",
				FullName:   "Test User",
				Role:       role,
				Department: models.DepartmentKhazwal,
				Status:     models.UserStatusActive,
			}

			// Allow all roles
			router := setupRouterWithUser(user, allRoles...)

			req := httptest.NewRequest("GET", "/protected", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Role %s: Status code = %d, expected %d", role, w.Code, http.StatusOK)
			}
		})
	}
}

// TestRequireRole_EmptyAllowedRoles memverifikasi behavior dengan empty allowed roles
func TestRequireRole_EmptyAllowedRoles(t *testing.T) {
	user := &models.User{
		ID:         1,
		NIP:        "12345",
		FullName:   "Test User",
		Role:       models.RoleAdmin,
		Department: models.DepartmentKhazwal,
		Status:     models.UserStatusActive,
	}

	// No roles specified - should deny all
	router := setupRouterWithUser(user)

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Status code = %d, expected %d for empty allowed roles", w.Code, http.StatusForbidden)
	}
}

// TestRequireRole_RoleHierarchy memverifikasi role hierarchy
func TestRequireRole_RoleHierarchy(t *testing.T) {
	tests := []struct {
		name          string
		userRole      models.UserRole
		requiredRole  models.UserRole
		shouldSucceed bool
	}{
		{
			name:          "Admin accessing Manager endpoint",
			userRole:      models.RoleAdmin,
			requiredRole:  models.RoleManager,
			shouldSucceed: true, // Admin biasanya punya akses ke Manager endpoints
		},
		{
			name:          "Manager accessing Admin endpoint",
			userRole:      models.RoleManager,
			requiredRole:  models.RoleAdmin,
			shouldSucceed: false, // Manager tidak bisa akses Admin-only
		},
		{
			name:          "Staff accessing Operator endpoint",
			userRole:      models.RoleStaffKhazwal,
			requiredRole:  models.RoleOperatorCetak,
			shouldSucceed: false, // Different roles, no access
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{
				ID:         1,
				NIP:        "12345",
				FullName:   "Test User",
				Role:       tt.userRole,
				Department: models.DepartmentKhazwal,
				Status:     models.UserStatusActive,
			}

			// Setup router dengan required roles termasuk admin
			allowedRoles := []models.UserRole{tt.requiredRole}
			if tt.userRole == models.RoleAdmin {
				// Admin biasanya selalu included
				allowedRoles = append(allowedRoles, models.RoleAdmin)
			}

			router := setupRouterWithUser(user, allowedRoles...)

			req := httptest.NewRequest("GET", "/protected", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			expectedCode := http.StatusForbidden
			if tt.shouldSucceed {
				expectedCode = http.StatusOK
			}

			if w.Code != expectedCode {
				t.Errorf("Status code = %d, expected %d", w.Code, expectedCode)
			}
		})
	}
}

// BenchmarkRequireRole mengukur performance role middleware
func BenchmarkRequireRole(b *testing.B) {
	user := &models.User{
		ID:         1,
		NIP:        "12345",
		FullName:   "Test User",
		Role:       models.RoleAdmin,
		Department: models.DepartmentKhazwal,
		Status:     models.UserStatusActive,
	}

	router := setupRouterWithUser(user, models.RoleAdmin, models.RoleManager)

	req := httptest.NewRequest("GET", "/protected", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
