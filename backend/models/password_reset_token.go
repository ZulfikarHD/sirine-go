package models

import (
	"time"
)

// PasswordResetToken merupakan model untuk token reset password
// dengan expiration tracking untuk security purposes
type PasswordResetToken struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64     `gorm:"not null;index" json:"user_id"`
	TokenHash string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"-"` // SHA256 hash dari token
	ExpiresAt time.Time  `gorm:"type:timestamp;not null;index" json:"expires_at"`
	UsedAt    *time.Time `gorm:"type:timestamp null" json:"used_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName menentukan nama tabel di database
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// IsValid memeriksa apakah token masih valid
// berdasarkan expiration dan usage status
func (t *PasswordResetToken) IsValid() bool {
	if t.UsedAt != nil {
		return false // Token sudah digunakan
	}
	return time.Now().Before(t.ExpiresAt)
}

// IsExpired memeriksa apakah token sudah expired
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed memeriksa apakah token sudah digunakan
func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}
