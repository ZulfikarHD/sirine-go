package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/smtp"
	"regexp"
	"sirine-go/backend/config"
	"sirine-go/backend/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// PasswordService merupakan service untuk password management
// yang mencakup hashing, verification, validation, change, dan reset password
type PasswordService struct {
	cost   int         // Bcrypt cost factor (default: 12)
	db     *gorm.DB    // Database connection untuk reset token management
	config *config.Config // Configuration untuk email service
}

// NewPasswordService membuat instance baru dari PasswordService
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: 12, // Cost 12 untuk balance antara security dan performance
	}
}

// NewPasswordServiceWithDB membuat instance PasswordService dengan database dan config
// untuk fitur password reset yang memerlukan database dan email service
func NewPasswordServiceWithDB(db *gorm.DB, cfg *config.Config) *PasswordService {
	return &PasswordService{
		cost:   cfg.BcryptCost,
		db:     db,
		config: cfg,
	}
}

// HashPassword menghasilkan bcrypt hash dari plaintext password
// dengan cost factor 12 untuk security yang optimal
func (s *PasswordService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password tidak boleh kosong")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword memverifikasi apakah plaintext password match dengan hash
func (s *PasswordService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ValidatePasswordPolicy memvalidasi password berdasarkan policy yang ditetapkan
// yaitu: minimum 8 karakter, minimal 1 uppercase, 1 number, dan 1 special character
func (s *PasswordService) ValidatePasswordPolicy(password string) error {
	if len(password) < 8 {
		return errors.New("password minimal 8 karakter")
	}

	// Check minimal 1 uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return errors.New("password harus mengandung minimal 1 huruf besar")
	}

	// Check minimal 1 number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return errors.New("password harus mengandung minimal 1 angka")
	}

	// Check minimal 1 special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		return errors.New("password harus mengandung minimal 1 karakter spesial (!@#$%^&*, dll)")
	}

	return nil
}

// ValidateAndHash memvalidasi password policy kemudian menghasilkan hash
// untuk simplify user creation process
func (s *PasswordService) ValidateAndHash(password string) (string, error) {
	if err := s.ValidatePasswordPolicy(password); err != nil {
		return "", err
	}
	return s.HashPassword(password)
}

// GetPasswordStrength menghitung strength password (0-4)
// dimana: 0 = Very Weak, 1 = Weak, 2 = Medium, 3 = Strong, 4 = Very Strong
func (s *PasswordService) GetPasswordStrength(password string) int {
	strength := 0

	// Length check
	if len(password) >= 8 {
		strength++
	}
	if len(password) >= 12 {
		strength++
	}

	// Character diversity
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		strength++
	}
	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		strength++
	}
	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		strength++
	}
	if regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password) {
		strength++
	}

	// Normalize to 0-4 scale
	if strength >= 5 {
		return 4 // Very Strong
	}
	return strength - 1
}

// ChangePassword mengubah password user dengan validasi current password
// dan memastikan new password tidak sama dengan current password
func (s *PasswordService) ChangePassword(userID uint64, currentPassword, newPassword string) error {
	if s.db == nil {
		return errors.New("database connection tidak tersedia")
	}

	// Get user dari database
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user tidak ditemukan")
		}
		return err
	}

	// Verify current password
	if !s.VerifyPassword(user.PasswordHash, currentPassword) {
		return errors.New("password saat ini tidak valid")
	}

	// Check new password tidak sama dengan current password
	if currentPassword == newPassword {
		return errors.New("password baru tidak boleh sama dengan password saat ini")
	}

	// Validate new password policy
	if err := s.ValidatePasswordPolicy(newPassword); err != nil {
		return err
	}

	// Hash new password
	newHash, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password dan reset must_change_password flag
	updates := map[string]interface{}{
		"password_hash":       newHash,
		"must_change_password": false,
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return err
	}

	// Revoke all existing sessions untuk force re-login
	s.db.Model(&models.UserSession{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true)

	return nil
}

// GenerateResetToken menghasilkan token untuk password reset
// dengan expiry 1 jam dan menyimpannya ke database
func (s *PasswordService) GenerateResetToken(userID uint64) (string, error) {
	if s.db == nil {
		return "", errors.New("database connection tidak tersedia")
	}

	// Generate random token (32 bytes = 64 hex characters)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// Hash token untuk storage (SHA256)
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// Invalidate semua existing reset tokens untuk user ini
	s.db.Model(&models.PasswordResetToken{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Update("used_at", time.Now())

	// Create new reset token dengan expiry 1 hour
	resetToken := models.PasswordResetToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.db.Create(&resetToken).Error; err != nil {
		return "", err
	}

	return token, nil
}

// ResetPassword mereset password user menggunakan reset token
// dan memvalidasi token validity sebelum reset
func (s *PasswordService) ResetPassword(token, newPassword string) error {
	if s.db == nil {
		return errors.New("database connection tidak tersedia")
	}

	// Hash token untuk lookup
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// Find reset token
	var resetToken models.PasswordResetToken
	if err := s.db.Where("token_hash = ?", tokenHash).First(&resetToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("token reset tidak valid atau sudah digunakan")
		}
		return err
	}

	// Check token validity
	if !resetToken.IsValid() {
		if resetToken.IsExpired() {
			return errors.New("token reset sudah kadaluarsa, silakan request ulang")
		}
		if resetToken.IsUsed() {
			return errors.New("token reset sudah digunakan")
		}
		return errors.New("token reset tidak valid")
	}

	// Validate new password policy
	if err := s.ValidatePasswordPolicy(newPassword); err != nil {
		return err
	}

	// Hash new password
	newHash, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user password
	var user models.User
	if err := s.db.First(&user, resetToken.UserID).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"password_hash":       newHash,
		"must_change_password": false,
		"failed_login_attempts": 0,
		"locked_until":        nil,
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return err
	}

	// Mark token as used
	now := time.Now()
	resetToken.UsedAt = &now
	s.db.Save(&resetToken)

	// Revoke all existing sessions
	s.db.Model(&models.UserSession{}).
		Where("user_id = ?", user.ID).
		Update("is_revoked", true)

	return nil
}

// SendResetEmail mengirim email dengan link reset password ke user
// menggunakan SMTP configuration dari config
func (s *PasswordService) SendResetEmail(email, token string) error {
	if s.config == nil {
		return errors.New("config tidak tersedia")
	}

	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return success untuk prevent email enumeration
			return nil
		}
		return err
	}

	// Build reset URL
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.config.FrontendURL, token)

	// Email subject dan body
	subject := "Reset Password - Sistem Produksi Pita Cukai"
	body := fmt.Sprintf(`
Halo %s,

Anda menerima email ini karena ada permintaan untuk mereset password akun Anda.

Klik link berikut untuk mereset password Anda:
%s

Link ini akan kadaluarsa dalam 1 jam.

Jika Anda tidak melakukan permintaan ini, abaikan email ini.

Terima kasih,
Tim Sistem Produksi Pita Cukai
`, user.FullName, resetURL)

	// Send email
	return s.sendEmail(email, subject, body)
}

// sendEmail mengirim email menggunakan SMTP configuration
func (s *PasswordService) sendEmail(to, subject, body string) error {
	// Build email message
	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.config.EmailFromAddress, to, subject, body)

	// SMTP authentication
	auth := smtp.PlainAuth(
		"",
		s.config.EmailUsername,
		s.config.EmailPassword,
		s.config.EmailSMTPHost,
	)

	// Send email via SMTP
	addr := fmt.Sprintf("%s:%d", s.config.EmailSMTPHost, s.config.EmailSMTPPort)
	err := smtp.SendMail(
		addr,
		auth,
		s.config.EmailFromAddress,
		[]string{to},
		[]byte(message),
	)

	return err
}

// RequestPasswordReset membuat reset token dan mengirim email sekaligus
// untuk simplify forgot password flow
func (s *PasswordService) RequestPasswordReset(nipOrEmail string) error {
	if s.db == nil {
		return errors.New("database connection tidak tersedia")
	}

	// Find user by NIP or Email
	var user models.User
	if err := s.db.Where("nip = ? OR email = ?", nipOrEmail, nipOrEmail).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return success untuk prevent enumeration
			return nil
		}
		return err
	}

	// Check jika user memiliki email
	if user.Email == "" {
		return errors.New("user tidak memiliki email terdaftar, hubungi administrator")
	}

	// Generate reset token
	token, err := s.GenerateResetToken(user.ID)
	if err != nil {
		return err
	}

	// Send email
	return s.SendResetEmail(user.Email, token)
}
