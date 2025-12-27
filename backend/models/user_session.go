package models

import (
	"time"
)

// UserSession merupakan model untuk tracking active sessions dengan JWT tokens
// yang mencakup device information, IP address, dan expiration tracking
type UserSession struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint64    `gorm:"not null;index" json:"user_id"`
	TokenHash        string    `gorm:"type:varchar(255);not null;index" json:"-"` // SHA256 hash dari JWT token
	RefreshTokenHash string    `gorm:"type:varchar(255);index" json:"-"`          // Hash dari refresh token
	DeviceInfo       string    `gorm:"type:varchar(500)" json:"device_info"`
	IPAddress        string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent        string    `gorm:"type:text" json:"user_agent"`
	ExpiresAt        time.Time `gorm:"type:timestamp;not null;index" json:"expires_at"`
	IsRevoked        bool      `gorm:"default:false;index" json:"is_revoked"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName menentukan nama tabel di database
func (UserSession) TableName() string {
	return "user_sessions"
}

// IsValid memeriksa apakah session masih valid
// berdasarkan expiration dan revocation status
func (s *UserSession) IsValid() bool {
	if s.IsRevoked {
		return false
	}
	return time.Now().Before(s.ExpiresAt)
}

// IsExpired memeriksa apakah session sudah expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
