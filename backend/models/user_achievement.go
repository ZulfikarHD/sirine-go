package models

import "time"

// UserAchievement merupakan model untuk tracking achievements yang telah
// di-unlock oleh user, dengan informasi achievement dan waktu unlock
type UserAchievement struct {
	ID            uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64      `gorm:"not null;index:idx_user_id" json:"user_id"`
	AchievementID uint64      `gorm:"not null;index:idx_achievement_id" json:"achievement_id"`
	UnlockedAt    time.Time   `gorm:"autoCreateTime;index:idx_unlocked_at" json:"unlocked_at"`
	
	// Relations
	User        User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Achievement Achievement `gorm:"foreignKey:AchievementID;constraint:OnDelete:CASCADE" json:"achievement,omitempty"`
}

// TableName menentukan nama tabel untuk UserAchievement model
func (UserAchievement) TableName() string {
	return "user_achievements"
}

// UserAchievementResponse merupakan response structure untuk user achievement
// yang mencakup detail achievement dan status unlock
type UserAchievementResponse struct {
	ID            uint64      `json:"id"`
	Code          string      `json:"code"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	Icon          string      `json:"icon"`
	Points        int         `json:"points"`
	Category      string      `json:"category"`
	IsUnlocked    bool        `json:"is_unlocked"`
	UnlockedAt    *time.Time  `json:"unlocked_at,omitempty"`
}
