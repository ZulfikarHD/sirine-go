package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// AchievementCategory merupakan kategori untuk achievement
// yang digunakan untuk grouping achievements berdasarkan tipenya
type AchievementCategory string

const (
	AchievementCategoryLogin        AchievementCategory = "LOGIN"
	AchievementCategoryProductivity AchievementCategory = "PRODUCTIVITY"
	AchievementCategoryQuality      AchievementCategory = "QUALITY"
	AchievementCategoryMilestone    AchievementCategory = "MILESTONE"
)

// AchievementCriteria merupakan struktur untuk criteria unlock achievement
// yang berisi kondisi yang harus dipenuhi untuk unlock achievement
type AchievementCriteria struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value,omitempty"`
}

// Scan mengimplementasikan sql.Scanner interface untuk database scanning
func (ac *AchievementCriteria) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, ac)
}

// Value mengimplementasikan driver.Valuer interface untuk database storage
func (ac AchievementCriteria) Value() (driver.Value, error) {
	return json.Marshal(ac)
}

// Achievement merupakan model untuk achievement dalam gamification system
// yang mencakup informasi achievement, points, criteria, dan status
type Achievement struct {
	ID          uint64              `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string              `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name        string              `gorm:"type:varchar(255);not null" json:"name"`
	Description string              `gorm:"type:text;not null" json:"description"`
	Icon        string              `gorm:"type:varchar(100)" json:"icon"`
	Points      int                 `gorm:"default:0;not null" json:"points"`
	Category    AchievementCategory `gorm:"type:enum('LOGIN','PRODUCTIVITY','QUALITY','MILESTONE');default:'MILESTONE'" json:"category"`
	Criteria    AchievementCriteria `gorm:"type:json" json:"criteria"`
	IsActive    bool                `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName menentukan nama tabel untuk Achievement model
func (Achievement) TableName() string {
	return "achievements"
}
