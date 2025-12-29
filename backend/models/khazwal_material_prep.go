package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// MaterialPrepStatus merupakan enum untuk status Khazwal Material Preparation
type MaterialPrepStatus string

const (
	MaterialPrepPending    MaterialPrepStatus = "PENDING"
	MaterialPrepInProgress MaterialPrepStatus = "IN_PROGRESS"
	MaterialPrepCompleted  MaterialPrepStatus = "COMPLETED"
)

// KhazwalMaterialPreparation merupakan model untuk entitas Persiapan Material
// yang mencakup pengambilan plat, kertas blanko, dan tinta
type KhazwalMaterialPreparation struct {
	ID                         uint64             `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductionOrderID          uint64             `gorm:"uniqueIndex;not null" json:"production_order_id" binding:"required"`
	SAPPlatCode                string             `gorm:"type:varchar(50);not null" json:"sap_plat_code" binding:"required"`
	KertasBlankoQuantity       int                `gorm:"not null" json:"kertas_blanko_quantity" binding:"required,min=1"`
	TintaRequirements          datatypes.JSON     `gorm:"type:json;not null" json:"tinta_requirements" binding:"required"`
	PlatRetrievedAt            *time.Time         `gorm:"type:timestamp null" json:"plat_retrieved_at"`
	KertasBlankoActual         *int               `gorm:"type:int null" json:"kertas_blanko_actual"`
	KertasBlankoVariance       *int               `gorm:"type:int null" json:"kertas_blanko_variance"`
	KertasBlankoVarianceReason string             `gorm:"type:varchar(500)" json:"kertas_blanko_variance_reason"`
	TintaActual                datatypes.JSON     `gorm:"type:json" json:"tinta_actual"`
	MaterialPhotos             datatypes.JSON     `gorm:"type:json" json:"material_photos"`
	Status                     MaterialPrepStatus `gorm:"type:enum('PENDING','IN_PROGRESS','COMPLETED');default:'PENDING'" json:"status"`
	StartedAt                  *time.Time         `gorm:"type:timestamp null" json:"started_at"`
	CompletedAt                *time.Time         `gorm:"type:timestamp null" json:"completed_at"`
	DurationMinutes            *int               `gorm:"type:int null" json:"duration_minutes"`
	PreparedBy                 *uint64            `gorm:"type:bigint unsigned null" json:"prepared_by"`
	Notes                      string             `gorm:"type:text" json:"notes"`
	CreatedAt                  time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt                  time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt                  gorm.DeletedAt     `gorm:"index" json:"-"`

	// Relationships
	ProductionOrder *ProductionOrder `gorm:"foreignKey:ProductionOrderID" json:"production_order,omitempty"`
	PreparedByUser  *User            `gorm:"foreignKey:PreparedBy" json:"prepared_by_user,omitempty"`
}

// TableName menentukan nama tabel di database
func (KhazwalMaterialPreparation) TableName() string {
	return "khazwal_material_preparations"
}

// IsPending memeriksa apakah material prep masih pending
func (kmp *KhazwalMaterialPreparation) IsPending() bool {
	return kmp.Status == MaterialPrepPending
}

// IsInProgress memeriksa apakah material prep sedang dalam progress
func (kmp *KhazwalMaterialPreparation) IsInProgress() bool {
	return kmp.Status == MaterialPrepInProgress
}

// IsCompleted memeriksa apakah material prep sudah selesai
func (kmp *KhazwalMaterialPreparation) IsCompleted() bool {
	return kmp.Status == MaterialPrepCompleted
}

// HasVariance memeriksa apakah ada variance dalam kertas blanko
func (kmp *KhazwalMaterialPreparation) HasVariance() bool {
	if kmp.KertasBlankoVariance == nil {
		return false
	}
	return *kmp.KertasBlankoVariance != 0
}

// CalculateDuration menghitung durasi dalam menit antara started dan completed
func (kmp *KhazwalMaterialPreparation) CalculateDuration() int {
	if kmp.StartedAt == nil || kmp.CompletedAt == nil {
		return 0
	}
	duration := kmp.CompletedAt.Sub(*kmp.StartedAt)
	return int(duration.Minutes())
}

// UpdateDuration mengupdate DurationMinutes field
func (kmp *KhazwalMaterialPreparation) UpdateDuration() {
	duration := kmp.CalculateDuration()
	kmp.DurationMinutes = &duration
}

// CalculateVariance menghitung variance kertas blanko
func (kmp *KhazwalMaterialPreparation) CalculateVariance() int {
	if kmp.KertasBlankoActual == nil {
		return 0
	}
	return *kmp.KertasBlankoActual - kmp.KertasBlankoQuantity
}

// UpdateVariance mengupdate KertasBlankoVariance field
func (kmp *KhazwalMaterialPreparation) UpdateVariance() {
	variance := kmp.CalculateVariance()
	kmp.KertasBlankoVariance = &variance
}
