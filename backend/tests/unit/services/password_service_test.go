package services_test

import (
	"sirine-go/backend/services"
	"testing"
)

// TestHashPassword memverifikasi password hashing dengan bcrypt
func TestHashPassword(t *testing.T) {
	service := services.NewPasswordService()
	password := "TestPassword123!"

	hash, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword gagal: %v", err)
	}

	if hash == "" {
		t.Error("Hash tidak boleh kosong")
	}

	if hash == password {
		t.Error("Hash tidak boleh sama dengan plaintext password")
	}

	// Verify hash dapat di-verify
	if !service.VerifyPassword(hash, password) {
		t.Error("VerifyPassword gagal untuk password yang valid")
	}
}

// TestHashPasswordEmpty memverifikasi error handling untuk empty password
func TestHashPasswordEmpty(t *testing.T) {
	service := services.NewPasswordService()

	_, err := service.HashPassword("")
	if err == nil {
		t.Error("HashPassword harus return error untuk empty password")
	}
}

// TestVerifyPassword memverifikasi password verification
func TestVerifyPassword(t *testing.T) {
	service := services.NewPasswordService()
	password := "CorrectPassword123!"

	hash, _ := service.HashPassword(password)

	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Correct password", "CorrectPassword123!", true},
		{"Wrong password", "WrongPassword123!", false},
		{"Empty password", "", false},
		{"Case sensitive", "correctpassword123!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.VerifyPassword(hash, tt.password)
			if result != tt.expected {
				t.Errorf("VerifyPassword(%s) = %v, expected %v", tt.password, result, tt.expected)
			}
		})
	}
}

// TestValidatePasswordPolicy memverifikasi password policy enforcement
func TestValidatePasswordPolicy(t *testing.T) {
	service := services.NewPasswordService()

	tests := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Valid password",
			password: "ValidPass123!",
			wantErr:  false,
		},
		{
			name:     "Too short",
			password: "Short1!",
			wantErr:  true,
			errMsg:   "password minimal 8 karakter",
		},
		{
			name:     "No uppercase",
			password: "nouppercas123!",
			wantErr:  true,
			errMsg:   "password harus mengandung minimal 1 huruf besar",
		},
		{
			name:     "No number",
			password: "NoNumber!",
			wantErr:  true,
			errMsg:   "password harus mengandung minimal 1 angka",
		},
		{
			name:     "No special char",
			password: "NoSpecial123",
			wantErr:  true,
			errMsg:   "password harus mengandung minimal 1 karakter spesial (!@#$%^&*, dll)",
		},
		{
			name:     "All requirements met",
			password: "Perfect123!@#",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidatePasswordPolicy(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("Error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestGetPasswordStrength memverifikasi password strength calculation
func TestGetPasswordStrength(t *testing.T) {
	service := services.NewPasswordService()

	tests := []struct {
		name     string
		password string
		minLevel int // Minimum expected strength level
	}{
		{"Very weak", "abc", 0},
		{"Weak", "abcd1234", 1},
		{"Medium", "Abcd1234", 2},
		{"Strong", "Abcd1234!", 3},
		{"Very strong", "Abcd1234!@#$", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strength := service.GetPasswordStrength(tt.password)
			if strength < tt.minLevel {
				t.Errorf("GetPasswordStrength(%s) = %d, want >= %d", tt.password, strength, tt.minLevel)
			}
		})
	}
}

// TestValidateAndHash memverifikasi combined validation dan hashing
func TestValidateAndHash(t *testing.T) {
	service := services.NewPasswordService()

	// Valid password
	validPassword := "ValidPass123!"
	hash, err := service.ValidateAndHash(validPassword)
	if err != nil {
		t.Errorf("ValidateAndHash gagal untuk valid password: %v", err)
	}
	if hash == "" {
		t.Error("Hash tidak boleh kosong untuk valid password")
	}

	// Invalid password
	invalidPassword := "weak"
	_, err = service.ValidateAndHash(invalidPassword)
	if err == nil {
		t.Error("ValidateAndHash harus return error untuk invalid password")
	}
}

// BenchmarkHashPassword mengukur performance password hashing
func BenchmarkHashPassword(b *testing.B) {
	service := services.NewPasswordService()
	password := "BenchmarkPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.HashPassword(password)
	}
}

// BenchmarkVerifyPassword mengukur performance password verification
func BenchmarkVerifyPassword(b *testing.B) {
	service := services.NewPasswordService()
	password := "BenchmarkPassword123!"
	hash, _ := service.HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.VerifyPassword(hash, password)
	}
}
