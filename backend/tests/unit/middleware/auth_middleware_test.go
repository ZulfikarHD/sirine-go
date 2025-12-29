package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"sirine-go/backend/config"
	"sirine-go/backend/middleware"
	"sirine-go/backend/models"
	"sirine-go/backend/services"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB membuat in-memory database untuk testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Gagal membuat test database: %v", err)
	}

	// Migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.UserSession{},
	)
	if err != nil {
		t.Fatalf("Gagal migrate database: %v", err)
	}

	return db
}

// getTestConfig membuat test configuration
func getTestConfig() *config.Config {
	return &config.Config{
		JWTSecret:          "test-secret-key-for-testing-only",
		JWTExpiry:          15 * time.Minute,
		RefreshTokenExpiry: 30 * 24 * time.Hour,
		MaxLoginAttempts:   5,
		LockoutDuration:    15 * time.Minute,
		BcryptCost:         4,
	}
}

// createTestUserWithToken membuat test user dan generate token
func createTestUserWithToken(t *testing.T, db *gorm.DB, cfg *config.Config) (*models.User, string) {
	passwordService := services.NewPasswordService()
	hash, _ := passwordService.HashPassword("TestPass123!")

	user := &models.User{
		NIP:          "12345",
		FullName:     "Test User",
		Email:        "test@example.com",
		Phone:        "08123456789",
		Role:         models.RoleAdmin,
		Department:   models.DepartmentKhazwal,
		PasswordHash: hash,
		Status:       models.UserStatusActive,
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Gagal create test user: %v", err)
	}

	// Generate token
	authService := services.NewAuthService(db, cfg)
	token, err := authService.GenerateJWT(user)
	if err != nil {
		t.Fatalf("Gagal generate JWT: %v", err)
	}

	// Create session
	req := services.LoginRequest{
		NIP:      user.NIP,
		Password: "TestPass123!",
	}
	authService.Login(req, "127.0.0.1", "test-agent")

	return user, token
}

// TestAuthRequired_ValidToken memverifikasi request dengan valid token
func TestAuthRequired_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	db := setupTestDB(t)
	cfg := getTestConfig()
	user, token := createTestUserWithToken(t, db, cfg)

	// Setup router dengan middleware
	router := gin.New()
	authService := services.NewAuthService(db, cfg)
	router.Use(middleware.AuthRequired(authService))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request dengan valid token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusOK)
	}

	// Verify user ada di context
	if w.Body.String() == "" {
		t.Error("Response body tidak boleh kosong")
	}

	_ = user // Use user variable
}

// TestAuthRequired_MissingToken memverifikasi request tanpa token
func TestAuthRequired_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	db := setupTestDB(t)
	cfg := getTestConfig()

	// Setup router dengan middleware
	router := gin.New()
	authService := services.NewAuthService(db, cfg)
	router.Use(middleware.AuthRequired(authService))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request tanpa token
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusUnauthorized)
	}
}

// TestAuthRequired_InvalidToken memverifikasi request dengan invalid token
func TestAuthRequired_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	db := setupTestDB(t)
	cfg := getTestConfig()

	// Setup router dengan middleware
	router := gin.New()
	authService := services.NewAuthService(db, cfg)
	router.Use(middleware.AuthRequired(authService))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request dengan invalid token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-12345")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusUnauthorized)
	}
}

// TestAuthRequired_ExpiredToken memverifikasi request dengan expired token
func TestAuthRequired_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	db := setupTestDB(t)
	cfg := getTestConfig()
	cfg.JWTExpiry = -1 * time.Hour // Expired token
	
	user, token := createTestUserWithToken(t, db, cfg)

	// Setup router dengan middleware menggunakan config normal
	normalCfg := getTestConfig()
	router := gin.New()
	authService := services.NewAuthService(db, normalCfg)
	router.Use(middleware.AuthRequired(authService))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request dengan expired token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusUnauthorized)
	}

	_ = user // Use user variable
}

// TestAuthRequired_MalformedAuthHeader memverifikasi various malformed auth headers
func TestAuthRequired_MalformedAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	db := setupTestDB(t)
	cfg := getTestConfig()

	tests := []struct {
		name        string
		authHeader  string
		expectedCode int
	}{
		{"Empty header", "", http.StatusUnauthorized},
		{"Only Bearer", "Bearer", http.StatusUnauthorized},
		{"No Bearer prefix", "token123", http.StatusUnauthorized},
		{"Wrong prefix", "Basic token123", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			authService := services.NewAuthService(db, cfg)
			router.Use(middleware.AuthRequired(authService))
			router.GET("/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest("GET", "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Status code = %d, expected %d", w.Code, tt.expectedCode)
			}
		})
	}
}

// TestAuthRequired_RevokedSession memverifikasi request dengan revoked session
func TestAuthRequired_RevokedSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	db := setupTestDB(t)
	cfg := getTestConfig()
	user, token := createTestUserWithToken(t, db, cfg)

	// Revoke session
	db.Model(&models.UserSession{}).
		Where("user_id = ?", user.ID).
		Update("is_revoked", true)

	// Setup router dengan middleware
	router := gin.New()
	authService := services.NewAuthService(db, cfg)
	router.Use(middleware.AuthRequired(authService))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request dengan revoked token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusUnauthorized)
	}
}

// BenchmarkAuthRequired mengukur performance auth middleware
func BenchmarkAuthRequired(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	db := setupTestDB(&testing.T{})
	cfg := getTestConfig()
	cfg.BcryptCost = 4
	_, token := createTestUserWithToken(&testing.T{}, db, cfg)

	router := gin.New()
	authService := services.NewAuthService(db, cfg)
	router.Use(middleware.AuthRequired(authService))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
