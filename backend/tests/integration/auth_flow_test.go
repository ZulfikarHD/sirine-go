package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sirine-go/backend/config"
	"sirine-go/backend/handlers"
	"sirine-go/backend/middleware"
	"sirine-go/backend/models"
	"sirine-go/backend/routes"
	"sirine-go/backend/services"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestAuthFlow merupakan comprehensive test untuk complete authentication flow
// yang mencakup login, access protected endpoints, refresh token, dan logout
type TestAuthFlow struct {
	db     *gorm.DB
	router *gin.Engine
	cfg    *config.Config
}

// setupTestEnv membuat test environment dengan database dan router
func setupTestEnv(t *testing.T) *TestAuthFlow {
	// Create in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Gagal membuat test database: %v", err)
	}

	// Migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.UserSession{},
		&models.PasswordResetToken{},
		&models.ActivityLog{},
		&models.Notification{},
	)
	if err != nil {
		t.Fatalf("Gagal migrate database: %v", err)
	}

	// Create test config
	cfg := &config.Config{
		JWTSecret:          "test-secret-key-for-integration-testing",
		JWTExpiry:          15 * time.Minute,
		RefreshTokenExpiry: 30 * 24 * time.Hour,
		MaxLoginAttempts:   5,
		LockoutDuration:    15 * time.Minute,
		BcryptCost:         4, // Lower cost untuk testing
		FrontendURL:        "http://localhost:5173",
	}

	// Setup router dengan routes
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Setup services
	authService := services.NewAuthService(db, cfg)
	userService := services.NewUserService(db)
	passwordService := services.NewPasswordServiceWithDB(db, cfg)
	
	// Setup handlers
	authHandler := handlers.NewAuthHandler(authService, passwordService)
	userHandler := handlers.NewUserHandler(userService, authService)
	
	// Setup routes (simplified untuk testing)
	routes.SetupAuthRoutes(router, authHandler, authService)
	routes.SetupUserRoutes(router, userHandler, authService)

	return &TestAuthFlow{
		db:     db,
		router: router,
		cfg:    cfg,
	}
}

// createTestUser membuat test user di database
func (tf *TestAuthFlow) createTestUser(t *testing.T) *models.User {
	passwordService := services.NewPasswordService()
	hash, _ := passwordService.HashPassword("TestPass123!")

	user := &models.User{
		NIP:          "12345",
		FullName:     "Integration Test User",
		Email:        "integration@example.com",
		Phone:        "08123456789",
		Role:         models.RoleAdmin,
		Department:   models.DepartmentKhazwal,
		PasswordHash: hash,
		Status:       models.UserStatusActive,
	}

	if err := tf.db.Create(user).Error; err != nil {
		t.Fatalf("Gagal create test user: %v", err)
	}

	return user
}

// TestCompleteLoginFlow memverifikasi complete login flow
func TestCompleteLoginFlow(t *testing.T) {
	tf := setupTestEnv(t)
	user := tf.createTestUser(t)

	// Step 1: Login dengan valid credentials
	loginReq := map[string]interface{}{
		"nip":         user.NIP,
		"password":    "TestPass123!",
		"remember_me": false,
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	tf.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login gagal: status code = %d, body = %s", w.Code, w.Body.String())
	}

	var loginResp struct {
		Success bool `json:"success"`
		Data    struct {
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("Gagal parse login response: %v", err)
	}

	if !loginResp.Success {
		t.Fatal("Login response success = false")
	}

	if loginResp.Data.Token == "" {
		t.Fatal("Token kosong")
	}

	token := loginResp.Data.Token

	// Step 2: Access protected endpoint dengan token
	req2 := httptest.NewRequest("GET", "/api/auth/me", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()

	tf.router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Access protected endpoint gagal: status code = %d", w2.Code)
	}

	// Step 3: Logout
	req3 := httptest.NewRequest("POST", "/api/auth/logout", nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()

	tf.router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Logout gagal: status code = %d", w3.Code)
	}

	// Step 4: Verify token tidak bisa digunakan setelah logout
	req4 := httptest.NewRequest("GET", "/api/auth/me", nil)
	req4.Header.Set("Authorization", "Bearer "+token)
	w4 := httptest.NewRecorder()

	tf.router.ServeHTTP(w4, req4)

	if w4.Code != http.StatusUnauthorized {
		t.Errorf("Token masih valid setelah logout: status code = %d", w4.Code)
	}
}

// TestLoginWithInvalidCredentials memverifikasi failed login attempts
func TestLoginWithInvalidCredentials(t *testing.T) {
	tf := setupTestEnv(t)
	user := tf.createTestUser(t)

	// Attempt 1: Wrong password
	loginReq := map[string]interface{}{
		"nip":      user.NIP,
		"password": "WrongPassword123!",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	tf.router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("Login seharusnya gagal dengan wrong password")
	}

	// Verify failed attempts incremented
	var updatedUser models.User
	tf.db.First(&updatedUser, user.ID)
	if updatedUser.FailedLoginAttempts != 1 {
		t.Errorf("FailedLoginAttempts = %d, expected 1", updatedUser.FailedLoginAttempts)
	}
}

// TestAccountLockoutAfterMaxAttempts memverifikasi account lockout mechanism
func TestAccountLockoutAfterMaxAttempts(t *testing.T) {
	tf := setupTestEnv(t)
	tf.cfg.MaxLoginAttempts = 3
	user := tf.createTestUser(t)

	// Attempt login 3 times dengan wrong password
	loginReq := map[string]interface{}{
		"nip":      user.NIP,
		"password": "WrongPassword123!",
	}

	for i := 0; i < 3; i++ {
		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		tf.router.ServeHTTP(w, req)
	}

	// Verify account is locked
	var updatedUser models.User
	tf.db.First(&updatedUser, user.ID)
	if updatedUser.LockedUntil == nil {
		t.Error("Account seharusnya locked setelah max attempts")
	}

	// Verify login dengan correct password juga gagal
	correctReq := map[string]interface{}{
		"nip":      user.NIP,
		"password": "TestPass123!",
	}

	body, _ := json.Marshal(correctReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	tf.router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("Login seharusnya gagal karena account locked")
	}
}

// TestRefreshTokenFlow memverifikasi refresh token functionality
func TestRefreshTokenFlow(t *testing.T) {
	tf := setupTestEnv(t)
	user := tf.createTestUser(t)

	// Step 1: Login
	loginReq := map[string]interface{}{
		"nip":      user.NIP,
		"password": "TestPass123!",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	tf.router.ServeHTTP(w, req)

	var loginResp struct {
		Success bool `json:"success"`
		Data    struct {
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &loginResp)

	oldToken := loginResp.Data.Token
	refreshToken := loginResp.Data.RefreshToken

	// Step 2: Refresh token
	refreshReq := map[string]interface{}{
		"refresh_token": refreshToken,
	}

	body2, _ := json.Marshal(refreshReq)
	req2 := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	tf.router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Refresh token gagal: status code = %d, body = %s", w2.Code, w2.Body.String())
	}

	var refreshResp struct {
		Success bool `json:"success"`
		Data    struct {
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &refreshResp)

	newToken := refreshResp.Data.Token

	if newToken == "" {
		t.Error("New token kosong")
	}

	if newToken == oldToken {
		t.Error("New token harus berbeda dari old token")
	}

	// Step 3: Verify new token bisa digunakan
	req3 := httptest.NewRequest("GET", "/api/auth/me", nil)
	req3.Header.Set("Authorization", "Bearer "+newToken)
	w3 := httptest.NewRecorder()

	tf.router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("New token tidak valid: status code = %d", w3.Code)
	}
}

// TestLoginWithRememberMe memverifikasi remember me functionality
func TestLoginWithRememberMe(t *testing.T) {
	tf := setupTestEnv(t)
	user := tf.createTestUser(t)

	// Login dengan remember_me = true
	loginReq := map[string]interface{}{
		"nip":         user.NIP,
		"password":    "TestPass123!",
		"remember_me": true,
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	tf.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login dengan remember_me gagal: status code = %d", w.Code)
	}

	// Verify session expiry lebih lama
	var session models.UserSession
	tf.db.Where("user_id = ?", user.ID).Order("created_at DESC").First(&session)

	expectedExpiry := time.Now().Add(tf.cfg.RefreshTokenExpiry)
	if session.ExpiresAt.Before(expectedExpiry.Add(-1 * time.Hour)) {
		t.Error("Session expiry terlalu pendek untuk remember_me")
	}
}

// TestLoginWithEmail memverifikasi login dengan email instead of NIP
func TestLoginWithEmail(t *testing.T) {
	tf := setupTestEnv(t)
	user := tf.createTestUser(t)

	// Login dengan email
	loginReq := map[string]interface{}{
		"nip":      user.Email, // Field name tetap nip, tapi value adalah email
		"password": "TestPass123!",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	tf.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Login dengan email gagal: status code = %d, body = %s", w.Code, w.Body.String())
	}
}

// TestMultipleSimultaneousSessions memverifikasi multiple sessions untuk same user
func TestMultipleSimultaneousSessions(t *testing.T) {
	tf := setupTestEnv(t)
	user := tf.createTestUser(t)

	// Login dari 2 devices berbeda
	loginReq := map[string]interface{}{
		"nip":      user.NIP,
		"password": "TestPass123!",
	}

	// Device 1
	body1, _ := json.Marshal(loginReq)
	req1 := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("User-Agent", "Device1")
	w1 := httptest.NewRecorder()
	tf.router.ServeHTTP(w1, req1)

	// Device 2
	body2, _ := json.Marshal(loginReq)
	req2 := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("User-Agent", "Device2")
	w2 := httptest.NewRecorder()
	tf.router.ServeHTTP(w2, req2)

	// Verify both sessions exist
	var sessionCount int64
	tf.db.Model(&models.UserSession{}).Where("user_id = ? AND is_revoked = ?", user.ID, false).Count(&sessionCount)

	if sessionCount != 2 {
		t.Errorf("Session count = %d, expected 2", sessionCount)
	}
}

// TestActivityLogCreation memverifikasi activity logs dibuat untuk auth events
func TestActivityLogCreation(t *testing.T) {
	tf := setupTestEnv(t)
	user := tf.createTestUser(t)

	// Login
	loginReq := map[string]interface{}{
		"nip":      user.NIP,
		"password": "TestPass123!",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	tf.router.ServeHTTP(w, req)

	// Verify activity log created
	var logCount int64
	tf.db.Model(&models.ActivityLog{}).
		Where("user_id = ? AND action = ?", user.ID, models.ActionLogin).
		Count(&logCount)

	if logCount == 0 {
		t.Error("Activity log tidak dibuat untuk login event")
	}
}
