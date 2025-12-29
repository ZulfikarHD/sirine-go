package models

import (
	"time"

	"gorm.io/gorm"
)

// POStageTracking merupakan model untuk tracking history setiap stage PO
// yang mencakup timestamp, duration, dan user yang handle
type POStageTracking struct {
	ID                uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductionOrderID uint64         `gorm:"index;not null" json:"production_order_id" binding:"required"`
	Stage             POStage        `gorm:"type:enum('KHAZWAL_MATERIAL_PREP','CETAK','VERIFIKASI','KHAZKHIR','COMPLETED');not null" json:"stage" binding:"required"`
	Status            POStatus       `gorm:"type:varchar(50);not null" json:"status" binding:"required"`
	StartedAt         *time.Time     `gorm:"type:timestamp null" json:"started_at"`
	CompletedAt       *time.Time     `gorm:"type:timestamp null" json:"completed_at"`
	DurationMinutes   *int           `gorm:"type:int null" json:"duration_minutes"`
	HandledBy         *uint64        `gorm:"type:bigint unsigned null" json:"handled_by"`
	Notes             string         `gorm:"type:text" json:"notes"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	ProductionOrder *ProductionOrder `gorm:"foreignKey:ProductionOrderID" json:"production_order,omitempty"`
	Handler         *User            `gorm:"foreignKey:HandledBy" json:"handler,omitempty"`
}

// TableName menentukan nama tabel di database
func (POStageTracking) TableName() string {
	return "po_stage_trackings"
}

// IsCompleted memeriksa apakah stage ini sudah selesai
func (pst *POStageTracking) IsCompleted() bool {
	return pst.CompletedAt != nil
}

// IsInProgress memeriksa apakah stage ini sedang dalam progress
func (pst *POStageTracking) IsInProgress() bool {
	return pst.StartedAt != nil && pst.CompletedAt == nil
}

// CalculateDuration menghitung durasi dalam menit antara started dan completed
func (pst *POStageTracking) CalculateDuration() int {
	if pst.StartedAt == nil || pst.CompletedAt == nil {
		return 0
	}
	duration := pst.CompletedAt.Sub(*pst.StartedAt)
	return int(duration.Minutes())
}

// UpdateDuration mengupdate DurationMinutes field
func (pst *POStageTracking) UpdateDuration() {
	duration := pst.CalculateDuration()
	pst.DurationMinutes = &duration
}
