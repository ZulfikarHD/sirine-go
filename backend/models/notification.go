package models

import (
	"time"
)

// NotificationType merupakan enum untuk tipe notifikasi
type NotificationType string

const (
	NotificationInfo        NotificationType = "INFO"
	NotificationSuccess     NotificationType = "SUCCESS"
	NotificationWarning     NotificationType = "WARNING"
	NotificationError       NotificationType = "ERROR"
	NotificationAchievement NotificationType = "ACHIEVEMENT"
)

// Notification merupakan model untuk in-app notifications
// yang digunakan untuk menampilkan notifikasi kepada user
type Notification struct {
	ID        uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64           `gorm:"not null;index:idx_user_read,priority:1" json:"user_id"`
	Title     string           `gorm:"type:varchar(255);not null" json:"title"`
	Message   string           `gorm:"type:text;not null" json:"message"`
	Type      NotificationType `gorm:"type:enum('INFO','SUCCESS','WARNING','ERROR','ACHIEVEMENT');default:'INFO'" json:"type"`
	IsRead    bool             `gorm:"default:false;index:idx_user_read,priority:2" json:"is_read"`
	ReadAt    *time.Time       `gorm:"type:timestamp null" json:"read_at"`
	CreatedAt time.Time        `gorm:"autoCreateTime" json:"created_at"`

	// Relationship
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName menentukan nama tabel di database
func (Notification) TableName() string {
	return "notifications"
}

// MarkAsRead menandai notifikasi sebagai sudah dibaca
func (n *Notification) MarkAsRead() {
	n.IsRead = true
	now := time.Now()
	n.ReadAt = &now
}
