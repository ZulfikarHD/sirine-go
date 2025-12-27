package models

import (
	"time"

	"gorm.io/gorm"
)

// Example model - ganti dengan model sesuai kebutuhan Anda
type Example struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title" binding:"required"`
	Content   string         `gorm:"type:text" json:"content"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
