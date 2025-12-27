package services

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService merupakan service untuk password management
// yang mencakup hashing, verification, dan validation
type PasswordService struct {
	cost int // Bcrypt cost factor (default: 12)
}

// NewPasswordService membuat instance baru dari PasswordService
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: 12, // Cost 12 untuk balance antara security dan performance
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
