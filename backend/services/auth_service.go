package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sirine-go/backend/config"
	"sirine-go/backend/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// AuthService merupakan service untuk authentication dan authorization
// yang mencakup login, logout, JWT generation, dan token validation
type AuthService struct {
	db              *gorm.DB
	passwordService *PasswordService
	config          *config.Config
}

// NewAuthService membuat instance baru dari AuthService
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:              db,
		passwordService: NewPasswordService(),
		config:          cfg,
	}
}

// JWTClaims merupakan custom claims untuk JWT token
type JWTClaims struct {
	UserID     uint64           `json:"user_id"`
	NIP        string           `json:"nip"`
	Email      string           `json:"email"`
	Role       models.UserRole  `json:"role"`
	Department models.Department `json:"department"`
	jwt.RegisteredClaims
}

// LoginRequest merupakan struktur untuk login request
// yang mendukung login dengan NIP atau Email
type LoginRequest struct {
	NIP        string `json:"nip" binding:"required"` // Bisa berisi NIP atau Email
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

// LoginResponse merupakan struktur untuk login response
type LoginResponse struct {
	Token               string           `json:"token"`
	RefreshToken        string           `json:"refresh_token"`
	User                models.SafeUser  `json:"user"`
	RequirePasswordChange bool           `json:"require_password_change"`
}

// Login melakukan autentikasi user dengan NIP atau Email dan password
// serta generate JWT token untuk subsequent requests
func (s *AuthService) Login(req LoginRequest, ipAddress, userAgent string) (*LoginResponse, error) {
	// Find user by NIP atau Email
	// Cek apakah input mengandung @ untuk determine email atau NIP
	var user models.User
	var err error
	
	identifier := req.NIP // Field name tetap NIP untuk backward compatibility
	
	// Try to find by NIP or Email
	err = s.db.Where("nip = ? OR email = ?", identifier, identifier).First(&user).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("NIP/Email atau password salah")
		}
		return nil, err
	}

	// Check jika user terkunci
	if user.IsLocked() {
		return nil, fmt.Errorf("akun Anda terkunci hingga %s karena terlalu banyak percobaan login gagal", 
			user.LockedUntil.Format("15:04:05"))
	}

	// Check jika user inactive
	if !user.IsActive() {
		return nil, errors.New("akun Anda tidak aktif, hubungi administrator")
	}

	// Verify password
	if !s.passwordService.VerifyPassword(user.PasswordHash, req.Password) {
		// Increment failed login attempts
		user.FailedLoginAttempts++
		
		// Lock account jika sudah mencapai limit
		if user.FailedLoginAttempts >= s.config.MaxLoginAttempts {
			lockUntil := time.Now().Add(s.config.LockoutDuration)
			user.LockedUntil = &lockUntil
			s.db.Save(&user)
			
			return nil, fmt.Errorf("terlalu banyak percobaan login gagal, akun Anda dikunci selama %d menit", 
				int(s.config.LockoutDuration.Minutes()))
		}
		
		s.db.Save(&user)
		return nil, errors.New("NIP/Email atau password salah")
	}

	// Reset failed attempts dan update last login
	user.FailedLoginAttempts = 0
	user.LockedUntil = nil
	now := time.Now()
	user.LastLoginAt = &now
	s.db.Save(&user)

	// Generate JWT token
	token, err := s.GenerateJWT(&user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.GenerateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	// Session expiry harus selalu mengikuti refresh token expiry (30 hari)
	// untuk memastikan refresh token mechanism berfungsi dengan baik.
	// RememberMe dapat digunakan untuk extend expiry lebih lama jika diperlukan.
	var expiresAt time.Time
	if req.RememberMe {
		// Extend session untuk remember me (misalnya 90 hari)
		expiresAt = time.Now().Add(s.config.RefreshTokenExpiry * 3)
	} else {
		// Default session expiry mengikuti refresh token (30 hari)
		expiresAt = time.Now().Add(s.config.RefreshTokenExpiry)
	}

	// Save session to database
	tokenHash := hashToken(token)
	refreshTokenHash := hashToken(refreshToken)
	
	session := models.UserSession{
		UserID:           user.ID,
		TokenHash:        tokenHash,
		RefreshTokenHash: refreshTokenHash,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		ExpiresAt:        expiresAt,
		IsRevoked:        false,
	}
	s.db.Create(&session)

	// Log activity
	s.logActivity(user.ID, models.ActionLogin, "user", &user.ID, ipAddress, userAgent)

	return &LoginResponse{
		Token:               token,
		RefreshToken:        refreshToken,
		User:                user.ToSafeUser(),
		RequirePasswordChange: user.MustChangePassword,
	}, nil
}

// Logout melakukan invalidasi token dan menandai session sebagai revoked
func (s *AuthService) Logout(userID uint64, token, ipAddress, userAgent string) error {
	tokenHash := hashToken(token)
	
	// Revoke session
	result := s.db.Model(&models.UserSession{}).
		Where("user_id = ? AND token_hash = ?", userID, tokenHash).
		Update("is_revoked", true)
	
	if result.Error != nil {
		return result.Error
	}

	// Log activity
	s.logActivity(userID, models.ActionLogout, "user", &userID, ipAddress, userAgent)

	return nil
}

// GenerateJWT menghasilkan JWT token untuk user
func (s *AuthService) GenerateJWT(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID:     user.ID,
		NIP:        user.NIP,
		Email:      user.Email,
		Role:       user.Role,
		Department: user.Department,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sirine-go",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// GenerateRefreshToken menghasilkan refresh token dengan expiry lebih lama
func (s *AuthService) GenerateRefreshToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID:     user.ID,
		NIP:        user.NIP,
		Email:      user.Email,
		Role:       user.Role,
		Department: user.Department,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sirine-go-refresh",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// ValidateToken memvalidasi JWT token dan mengembalikan user data
func (s *AuthService) ValidateToken(tokenString string) (*models.User, *JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, nil, errors.New("invalid token")
	}

	// Check jika token sudah revoked
	tokenHash := hashToken(tokenString)
	var session models.UserSession
	if err := s.db.Where("token_hash = ?", tokenHash).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, errors.New("session tidak valid")
		}
		return nil, nil, err
	}

	if !session.IsValid() {
		return nil, nil, errors.New("session expired atau sudah di-revoke")
	}

	// Get user data
	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, nil, err
	}

	if !user.IsActive() {
		return nil, nil, errors.New("user tidak aktif")
	}

	return &user, claims, nil
}

// RefreshAuthToken me-refresh JWT token menggunakan refresh token
func (s *AuthService) RefreshAuthToken(refreshToken string) (*LoginResponse, error) {
	// Parse and validate refresh token JWT structure
	token, err := jwt.ParseWithClaims(refreshToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, errors.New("refresh token tidak valid")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("refresh token tidak valid")
	}

	// Verify this is a refresh token by checking issuer
	if claims.Issuer != "sirine-go-refresh" {
		return nil, errors.New("refresh token tidak valid")
	}

	// Check jika refresh token sudah revoked dengan menggunakan refresh_token_hash
	refreshTokenHash := hashToken(refreshToken)
	var session models.UserSession
	if err := s.db.Where("refresh_token_hash = ?", refreshTokenHash).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("refresh token tidak valid")
		}
		return nil, err
	}

	if !session.IsValid() {
		return nil, errors.New("refresh token expired atau sudah di-revoke")
	}

	// Get user data
	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if !user.IsActive() {
		return nil, errors.New("user tidak aktif")
	}

	// Generate new tokens
	newToken, err := s.GenerateJWT(&user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.GenerateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	// Update session dengan token dan refresh token baru
	// Session expiry tetap mengikuti refresh token expiry (30 hari)
	newTokenHash := hashToken(newToken)
	newRefreshTokenHash := hashToken(newRefreshToken)

	s.db.Model(&models.UserSession{}).
		Where("refresh_token_hash = ?", refreshTokenHash).
		Updates(map[string]interface{}{
			"token_hash":         newTokenHash,
			"refresh_token_hash": newRefreshTokenHash,
			"expires_at":         time.Now().Add(s.config.RefreshTokenExpiry),
		})

	return &LoginResponse{
		Token:               newToken,
		RefreshToken:        newRefreshToken,
		User:                user.ToSafeUser(),
		RequirePasswordChange: user.MustChangePassword,
	}, nil
}

// hashToken menghasilkan SHA256 hash dari token untuk storage
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// logActivity mencatat activity ke activity_logs table
func (s *AuthService) logActivity(userID uint64, action models.ActivityAction, entityType string, entityID *uint64, ipAddress, userAgent string) {
	log := models.ActivityLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}
	s.db.Create(&log)
}
