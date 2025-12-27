package models_test

import (
	"sirine-go/backend/models"
	"testing"
	"time"
)

// TestUserIsLocked memverifikasi user lock status checking
func TestUserIsLocked(t *testing.T) {
	tests := []struct {
		name        string
		lockedUntil *time.Time
		expected    bool
	}{
		{
			name:        "Not locked - nil",
			lockedUntil: nil,
			expected:    false,
		},
		{
			name: "Locked - future time",
			lockedUntil: func() *time.Time {
				t := time.Now().Add(10 * time.Minute)
				return &t
			}(),
			expected: true,
		},
		{
			name: "Not locked - past time",
			lockedUntil: func() *time.Time {
				t := time.Now().Add(-10 * time.Minute)
				return &t
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{
				LockedUntil: tt.lockedUntil,
			}
			result := user.IsLocked()
			if result != tt.expected {
				t.Errorf("IsLocked() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestUserIsActive memverifikasi user active status checking
func TestUserIsActive(t *testing.T) {
	tests := []struct {
		name        string
		status      models.UserStatus
		lockedUntil *time.Time
		expected    bool
	}{
		{
			name:        "Active and not locked",
			status:      models.StatusActive,
			lockedUntil: nil,
			expected:    true,
		},
		{
			name:     "Active but locked",
			status:   models.StatusActive,
			lockedUntil: func() *time.Time {
				t := time.Now().Add(10 * time.Minute)
				return &t
			}(),
			expected: false,
		},
		{
			name:        "Inactive",
			status:      models.StatusInactive,
			lockedUntil: nil,
			expected:    false,
		},
		{
			name:        "Suspended",
			status:      models.StatusSuspended,
			lockedUntil: nil,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{
				Status:      tt.status,
				LockedUntil: tt.lockedUntil,
			}
			result := user.IsActive()
			if result != tt.expected {
				t.Errorf("IsActive() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestUserHasRole memverifikasi role checking
func TestUserHasRole(t *testing.T) {
	user := &models.User{
		Role: models.RoleAdmin,
	}

	tests := []struct {
		name     string
		roles    []models.UserRole
		expected bool
	}{
		{
			name:     "Has role - single",
			roles:    []models.UserRole{models.RoleAdmin},
			expected: true,
		},
		{
			name:     "Has role - multiple",
			roles:    []models.UserRole{models.RoleManager, models.RoleAdmin, models.RoleStaffKhazwal},
			expected: true,
		},
		{
			name:     "Does not have role",
			roles:    []models.UserRole{models.RoleOperatorCetak},
			expected: false,
		},
		{
			name:     "Empty roles",
			roles:    []models.UserRole{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.HasRole(tt.roles...)
			if result != tt.expected {
				t.Errorf("HasRole() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestUserIsAdmin memverifikasi admin checking
func TestUserIsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     models.UserRole
		expected bool
	}{
		{
			name:     "Admin role",
			role:     models.RoleAdmin,
			expected: true,
		},
		{
			name:     "Manager role",
			role:     models.RoleManager,
			expected: true,
		},
		{
			name:     "Staff role",
			role:     models.RoleStaffKhazwal,
			expected: false,
		},
		{
			name:     "Operator role",
			role:     models.RoleOperatorCetak,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{
				Role: tt.role,
			}
			result := user.IsAdmin()
			if result != tt.expected {
				t.Errorf("IsAdmin() = %v, expected %v for role %s", result, tt.expected, tt.role)
			}
		})
	}
}

// TestUserToSafeUser memverifikasi safe user conversion
func TestUserToSafeUser(t *testing.T) {
	now := time.Now()
	user := &models.User{
		ID:              1,
		NIP:             "99999",
		FullName:        "Test User",
		Email:           "test@example.com",
		Phone:           "081234567890",
		PasswordHash:    "secret_hash",
		Role:            models.RoleAdmin,
		Department:      models.DeptKhazwal,
		Shift:           models.ShiftPagi,
		ProfilePhotoURL: "https://example.com/photo.jpg",
		Status:          models.StatusActive,
		LastLoginAt:     &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	safeUser := user.ToSafeUser()

	// Verify all fields copied correctly
	if safeUser.ID != user.ID {
		t.Errorf("SafeUser.ID = %v, expected %v", safeUser.ID, user.ID)
	}
	if safeUser.NIP != user.NIP {
		t.Errorf("SafeUser.NIP = %v, expected %v", safeUser.NIP, user.NIP)
	}
	if safeUser.FullName != user.FullName {
		t.Errorf("SafeUser.FullName = %v, expected %v", safeUser.FullName, user.FullName)
	}
	if safeUser.Email != user.Email {
		t.Errorf("SafeUser.Email = %v, expected %v", safeUser.Email, user.Email)
	}
	if safeUser.Role != user.Role {
		t.Errorf("SafeUser.Role = %v, expected %v", safeUser.Role, user.Role)
	}

	// Verify password hash is NOT included (implicit - SafeUser doesn't have PasswordHash field)
	// This is a compile-time check, not runtime
}

// TestUserSessionIsValid memverifikasi session validation
func TestUserSessionIsValid(t *testing.T) {
	tests := []struct {
		name      string
		isRevoked bool
		expiresAt time.Time
		expected  bool
	}{
		{
			name:      "Valid session",
			isRevoked: false,
			expiresAt: time.Now().Add(10 * time.Minute),
			expected:  true,
		},
		{
			name:      "Revoked session",
			isRevoked: true,
			expiresAt: time.Now().Add(10 * time.Minute),
			expected:  false,
		},
		{
			name:      "Expired session",
			isRevoked: false,
			expiresAt: time.Now().Add(-10 * time.Minute),
			expected:  false,
		},
		{
			name:      "Revoked and expired",
			isRevoked: true,
			expiresAt: time.Now().Add(-10 * time.Minute),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &models.UserSession{
				IsRevoked: tt.isRevoked,
				ExpiresAt: tt.expiresAt,
			}
			result := session.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestPasswordResetTokenIsValid memverifikasi reset token validation
func TestPasswordResetTokenIsValid(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		usedAt    *time.Time
		expiresAt time.Time
		expected  bool
	}{
		{
			name:      "Valid token",
			usedAt:    nil,
			expiresAt: now.Add(1 * time.Hour),
			expected:  true,
		},
		{
			name:      "Used token",
			usedAt:    &now,
			expiresAt: now.Add(1 * time.Hour),
			expected:  false,
		},
		{
			name:      "Expired token",
			usedAt:    nil,
			expiresAt: now.Add(-1 * time.Hour),
			expected:  false,
		},
		{
			name:      "Used and expired",
			usedAt:    &now,
			expiresAt: now.Add(-1 * time.Hour),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := &models.PasswordResetToken{
				UsedAt:    tt.usedAt,
				ExpiresAt: tt.expiresAt,
			}
			result := token.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
