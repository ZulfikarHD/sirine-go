package cutting

import (
	"time"

	"gorm.io/gorm"
)

// CuttingStatus merupakan enum untuk status Khazwal Cutting
type CuttingStatus string

const (
	CuttingPending    CuttingStatus = "PENDING"
	CuttingInProgress CuttingStatus = "IN_PROGRESS"
	CuttingCompleted  CuttingStatus = "COMPLETED"
)

// KhazwalCuttingResult merupakan model untuk entitas Hasil Pemotongan
// yang mencakup input lembar besar, output sisiran kiri & kanan, dan waste tracking
type KhazwalCuttingResult struct {
	ID                  uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductionOrderID   uint64         `gorm:"uniqueIndex;not null" json:"production_order_id" binding:"required"`
	
	// Input dari Counting Result
	InputLembarBesar    int            `gorm:"not null;default:0" json:"input_lembar_besar" binding:"min=0"`
	ExpectedOutput      int            `gorm:"not null;default:0" json:"expected_output" binding:"min=0"` // input × 2
	
	// Output Results (Sisiran Kiri & Kanan)
	OutputSisiranKiri   *int           `gorm:"type:int null" json:"output_sisiran_kiri" binding:"omitempty,min=0"`
	OutputSisiranKanan  *int           `gorm:"type:int null" json:"output_sisiran_kanan" binding:"omitempty,min=0"`
	TotalOutput         int            `gorm:"not null;default:0" json:"total_output"`
	
	// Waste Tracking
	WasteQuantity       int            `gorm:"not null;default:0" json:"waste_quantity"`
	WastePercentage     *float64       `gorm:"type:decimal(5,2)" json:"waste_percentage"`
	WasteReason         string         `gorm:"type:text" json:"waste_reason"`
	WastePhotoURL       string         `gorm:"type:varchar(500)" json:"waste_photo_url"`
	
	// Machine & Staff
	CuttingMachine      string         `gorm:"type:varchar(100)" json:"cutting_machine" binding:"required"`
	CutBy               *uint64        `gorm:"type:bigint unsigned null" json:"cut_by"`
	
	// Status & Timing
	Status              CuttingStatus  `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"`
	StartedAt           *time.Time     `gorm:"type:timestamp null" json:"started_at"`
	CompletedAt         *time.Time     `gorm:"type:timestamp null" json:"completed_at"`
	DurationMinutes     *int           `gorm:"type:int null" json:"duration_minutes"`
	
	// Timestamps
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName menentukan nama tabel di database
func (KhazwalCuttingResult) TableName() string {
	return "khazwal_cutting_results"
}

// IsPending memeriksa apakah cutting masih pending
func (kcr *KhazwalCuttingResult) IsPending() bool {
	return kcr.Status == CuttingPending
}

// IsInProgress memeriksa apakah cutting sedang dalam progress
func (kcr *KhazwalCuttingResult) IsInProgress() bool {
	return kcr.Status == CuttingInProgress
}

// IsCompleted memeriksa apakah cutting sudah selesai
func (kcr *KhazwalCuttingResult) IsCompleted() bool {
	return kcr.Status == CuttingCompleted
}

// CalculateDuration menghitung durasi dalam menit antara started dan completed
func (kcr *KhazwalCuttingResult) CalculateDuration() int {
	if kcr.StartedAt == nil || kcr.CompletedAt == nil {
		return 0
	}
	duration := kcr.CompletedAt.Sub(*kcr.StartedAt)
	return int(duration.Minutes())
}

// UpdateDuration mengupdate DurationMinutes field
func (kcr *KhazwalCuttingResult) UpdateDuration() {
	duration := kcr.CalculateDuration()
	kcr.DurationMinutes = &duration
}

// CalculateTotalOutput menghitung total output dari sisiran kiri & kanan
func (kcr *KhazwalCuttingResult) CalculateTotalOutput() int {
	kiri := 0
	kanan := 0
	
	if kcr.OutputSisiranKiri != nil {
		kiri = *kcr.OutputSisiranKiri
	}
	if kcr.OutputSisiranKanan != nil {
		kanan = *kcr.OutputSisiranKanan
	}
	
	return kiri + kanan
}

// UpdateTotalOutput mengupdate TotalOutput field
func (kcr *KhazwalCuttingResult) UpdateTotalOutput() {
	kcr.TotalOutput = kcr.CalculateTotalOutput()
}

// CalculateWaste menghitung waste quantity
func (kcr *KhazwalCuttingResult) CalculateWaste() int {
	return kcr.ExpectedOutput - kcr.TotalOutput
}

// UpdateWaste mengupdate WasteQuantity field
func (kcr *KhazwalCuttingResult) UpdateWaste() {
	kcr.WasteQuantity = kcr.CalculateWaste()
}

// CalculateWastePercentage menghitung waste percentage
func (kcr *KhazwalCuttingResult) CalculateWastePercentage() float64 {
	if kcr.ExpectedOutput == 0 {
		return 0
	}
	return (float64(kcr.WasteQuantity) / float64(kcr.ExpectedOutput)) * 100
}

// UpdateWastePercentage mengupdate WastePercentage field
func (kcr *KhazwalCuttingResult) UpdateWastePercentage() {
	percentage := kcr.CalculateWastePercentage()
	kcr.WastePercentage = &percentage
}

// WasteExceedsThreshold memeriksa apakah waste melebihi threshold (default 2%)
func (kcr *KhazwalCuttingResult) WasteExceedsThreshold(threshold float64) bool {
	if kcr.WastePercentage == nil {
		return false
	}
	return *kcr.WastePercentage > threshold
}

// RequiresWasteDocumentation memeriksa apakah waste memerlukan dokumentasi (reason & photo)
func (kcr *KhazwalCuttingResult) RequiresWasteDocumentation() bool {
	return kcr.WasteExceedsThreshold(2.0)
}

// HasWasteDocumentation memeriksa apakah waste documentation sudah lengkap
func (kcr *KhazwalCuttingResult) HasWasteDocumentation() bool {
	return kcr.WasteReason != "" && kcr.WastePhotoURL != ""
}

// ValidateForFinalization memeriksa apakah data ready untuk finalisasi
func (kcr *KhazwalCuttingResult) ValidateForFinalization() error {
	if kcr.OutputSisiranKiri == nil || kcr.OutputSisiranKanan == nil {
		return ErrMissingOutputData
	}
	
	if kcr.RequiresWasteDocumentation() && !kcr.HasWasteDocumentation() {
		return ErrMissingWasteDocumentation
	}
	
	return nil
}
