package services_test

import (
	"sirine-go/backend/config"
	"sirine-go/backend/models"
	"sirine-go/backend/services"
	"testing"
	"time"

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
		&models.PasswordResetToken{},
		&models.ActivityLog{},
	)
	if err != nil {
		t.Fatalf("Gagal migrate database: %v", err)
	}

	return db
}

// createTestUser membuat test user di database
func createTestUser(t *testing.T, db *gorm.DB) *models.User {
	passwordService := services.NewPasswordService()
	hash, _ := passwordService.HashPassword("TestPass123!")

	user := &models.User{
		NIP:          "12345",
		FullName:     "Test User",
		Email:        "test@example.com",
		Phone:        "08123456789",
		Role:         models.RoleAdmin,
		Department:   models.DeptKhazwal,
		PasswordHash: hash,
		Status:       models.StatusActive,
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Gagal create test user: %v", err)
	}

	return user
}

// getTestConfig membuat test configuration
func getTestConfig() *config.Config {
	return &config.Config{
		JWTSecret:          "test-secret-key-for-testing-only",
		JWTExpiry:          15 * time.Minute,
		RefreshTokenExpiry: 30 * 24 * time.Hour,
		MaxLoginAttempts:   5,
		LockoutDuration:    15 * time.Minute,
		BcryptCost:         4, // Lower cost untuk testing
	}
}

// TestLogin_Success memverifikasi login berhasil dengan credentials valid
func TestLogin_Success(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(t, db)

	req := services.LoginRequest{
		NIP:        user.NIP,
		Password:   "TestPass123!",
		RememberMe: false,
	}

	response, err := authService.Login(req, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Login gagal: %v", err)
	}

	if response.Token == "" {
		t.Error("Token tidak boleh kosong")
	}

	if response.RefreshToken == "" {
		t.Error("Refresh token tidak boleh kosong")
	}

	if response.User.NIP != user.NIP {
		t.Errorf("User NIP = %v, expected %v", response.User.NIP, user.NIP)
	}
}

// TestLogin_WithEmail memverifikasi login dengan email
func TestLogin_WithEmail(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(t, db)

	req := services.LoginRequest{
		NIP:        user.Email, // Login dengan email
		Password:   "TestPass123!",
		RememberMe: false,
	}

	response, err := authService.Login(req, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Login dengan email gagal: %v", err)
	}

	if response.User.Email != user.Email {
		t.Errorf("User email = %v, expected %v", response.User.Email, user.Email)
	}
}

// TestLogin_InvalidPassword memverifikasi login gagal dengan password salah
func TestLogin_InvalidPassword(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(t, db)

	req := services.LoginRequest{
		NIP:      user.NIP,
		Password: "WrongPassword123!",
	}

	_, err := authService.Login(req, "127.0.0.1", "test-agent")
	if err == nil {
		t.Error("Login harus gagal dengan password salah")
	}

	// Verify failed attempts incremented
	var updatedUser models.User
	db.First(&updatedUser, user.ID)
	if updatedUser.FailedLoginAttempts != 1 {
		t.Errorf("FailedLoginAttempts = %d, expected 1", updatedUser.FailedLoginAttempts)
	}
}

// TestLogin_UserNotFound memverifikasi login gagal dengan NIP tidak terdaftar
func TestLogin_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	req := services.LoginRequest{
		NIP:      "99999",
		Password: "TestPass123!",
	}

	_, err := authService.Login(req, "127.0.0.1", "test-agent")
	if err == nil {
		t.Error("Login harus gagal dengan NIP tidak terdaftar")
	}
}

// TestLogin_AccountLocked memverifikasi login gagal untuk locked account
func TestLogin_AccountLocked(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(t, db)

	// Lock account
	lockUntil := time.Now().Add(15 * time.Minute)
	user.LockedUntil = &lockUntil
	db.Save(user)

	req := services.LoginRequest{
		NIP:      user.NIP,
		Password: "TestPass123!",
	}

	_, err := authService.Login(req, "127.0.0.1", "test-agent")
	if err == nil {
		t.Error("Login harus gagal untuk locked account")
	}
}

// TestLogin_InactiveUser memverifikasi login gagal untuk inactive user
func TestLogin_InactiveUser(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(t, db)

	// Set user inactive
	user.Status = models.StatusInactive
	db.Save(user)

	req := services.LoginRequest{
		NIP:      user.NIP,
		Password: "TestPass123!",
	}

	_, err := authService.Login(req, "127.0.0.1", "test-agent")
	if err == nil {
		t.Error("Login harus gagal untuk inactive user")
	}
}

// TestLogin_AccountLockAfterMaxAttempts memverifikasi account lock setelah max attempts
func TestLogin_AccountLockAfterMaxAttempts(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	cfg.MaxLoginAttempts = 3
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(t, db)

	req := services.LoginRequest{
		NIP:      user.NIP,
		Password: "WrongPassword123!",
	}

	// Attempt login 3 times dengan wrong password
	for i := 0; i < 3; i++ {
		authService.Login(req, "127.0.0.1", "test-agent")
	}

	// Verify account is locked
	var updatedUser models.User
	db.First(&updatedUser, user.ID)
	if updatedUser.LockedUntil == nil {
		t.Error("Account harus terkunci setelah max attempts")
	}
}

// TestGenerateJWT memverifikasi JWT token generation
func TestGenerateJWT(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(t, db)

	token, err := authService.GenerateJWT(user)
	if err != nil {
		t.Fatalf("GenerateJWT gagal: %v", err)
	}

	if token == "" {
		t.Error("Token tidak boleh kosong")
	}

	// Validate token
	validatedUser, claims, err := authService.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken gagal: %v", err)
	}

	if validatedUser.ID != user.ID {
		t.Errorf("User ID = %v, expected %v", validatedUser.ID, user.ID)
	}

	if claims.NIP != user.NIP {
		t.Errorf("Claims NIP = %v, expected %v", claims.NIP, user.NIP)
	}
}

// TestGenerateRefreshToken memverifikasi refresh token generation
func TestGenerateRefreshToken(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(t, db)

	refreshToken, err := authService.GenerateRefreshToken(user)
	if err != nil {
		t.Fatalf("GenerateRefreshToken gagal: %v", err)
	}

	if refreshToken == "" {
		t.Error("Refresh token tidak boleh kosong")
	}
}

// TestValidateToken_InvalidToken memverifikasi validation untuk invalid token
func TestValidateToken_InvalidToken(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	_, _, err := authService.ValidateToken("invalid-token")
	if err == nil {
		t.Error("ValidateToken harus return error untuk invalid token")
	}
}

// TestLogout memverifikasi logout functionality
func TestLogout(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(t, db)

	// Login first
	req := services.LoginRequest{
		NIP:      user.NIP,
		Password: "TestPass123!",
	}
	response, _ := authService.Login(req, "127.0.0.1", "test-agent")

	// Logout
	err := authService.Logout(user.ID, response.Token, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Logout gagal: %v", err)
	}

	// Verify session is revoked
	var session models.UserSession
	db.Where("user_id = ?", user.ID).First(&session)
	if !session.IsRevoked {
		t.Error("Session harus di-revoke setelah logout")
	}
}

// TestRefreshAuthToken memverifikasi token refresh functionality
func TestRefreshAuthToken(t *testing.T) {
	db := setupTestDB(t)
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(t, db)

	// Login first
	req := services.LoginRequest{
		NIP:      user.NIP,
		Password: "TestPass123!",
	}
	loginResponse, _ := authService.Login(req, "127.0.0.1", "test-agent")

	// Create session for refresh token
	refreshTokenHash := hashTokenForTest(loginResponse.RefreshToken)
	session := models.UserSession{
		UserID:           user.ID,
		TokenHash:        refreshTokenHash,
		RefreshTokenHash: refreshTokenHash,
		IPAddress:        "127.0.0.1",
		UserAgent:        "test-agent",
		ExpiresAt:        time.Now().Add(24 * time.Hour),
		IsRevoked:        false,
	}
	db.Create(&session)

	// Refresh token
	response, err := authService.RefreshAuthToken(loginResponse.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshAuthToken gagal: %v", err)
	}

	if response.Token == "" {
		t.Error("New token tidak boleh kosong")
	}

	if response.Token == loginResponse.Token {
		t.Error("New token harus berbeda dari old token")
	}
}

// Helper function untuk hash token dalam test
func hashTokenForTest(token string) string {
	// Implement hash function atau copy dari auth_service
	return token
}

// BenchmarkLogin mengukur performance login operation
func BenchmarkLogin(b *testing.B) {
	db := setupTestDB(&testing.T{})
	cfg := getTestConfig()
	cfg.BcryptCost = 4 // Lower cost untuk benchmark
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(&testing.T{}, db)

	req := services.LoginRequest{
		NIP:      user.NIP,
		Password: "TestPass123!",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		authService.Login(req, "127.0.0.1", "test-agent")
	}
}

// BenchmarkGenerateJWT mengukur performance JWT generation
func BenchmarkGenerateJWT(b *testing.B) {
	db := setupTestDB(&testing.T{})
	cfg := getTestConfig()
	authService := services.NewAuthService(db, cfg)

	user := createTestUser(&testing.T{}, db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		authService.GenerateJWT(user)
	}
}
