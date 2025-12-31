package counting

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CountingStatus merupakan enum untuk status Khazwal Counting
type CountingStatus string

const (
	CountingPending    CountingStatus = "PENDING"
	CountingInProgress CountingStatus = "IN_PROGRESS"
	CountingCompleted  CountingStatus = "COMPLETED"
)

// DefectBreakdownItem merupakan struct untuk item breakdown kerusakan
type DefectBreakdownItem struct {
	Type     string `json:"type" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

// KhazwalCountingResult merupakan model untuk entitas Hasil Penghitungan
// yang mencakup jumlah baik, rusak, breakdown kerusakan, dan variance tracking
type KhazwalCountingResult struct {
	ID                   uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductionOrderID    uint64         `gorm:"uniqueIndex;not null" json:"production_order_id" binding:"required"`
	
	// Counting Results (Lembar Besar)
	QuantityGood         int            `gorm:"not null;default:0" json:"quantity_good" binding:"min=0"`
	QuantityDefect       int            `gorm:"not null;default:0" json:"quantity_defect" binding:"min=0"`
	TotalCounted         int            `gorm:"->;type:integer GENERATED ALWAYS AS (quantity_good + quantity_defect) STORED" json:"total_counted"`
	VarianceFromTarget   *int           `gorm:"type:int null" json:"variance_from_target"`
	
	// Percentages
	PercentageGood       *float64       `gorm:"type:decimal(5,2)" json:"percentage_good"`
	PercentageDefect     *float64       `gorm:"type:decimal(5,2)" json:"percentage_defect"`
	
	// Defect Breakdown
	DefectBreakdown      datatypes.JSON `gorm:"type:jsonb" json:"defect_breakdown"`
	
	// Status & Timing
	Status               CountingStatus `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"`
	StartedAt            *time.Time     `gorm:"type:timestamp null" json:"started_at"`
	CompletedAt          *time.Time     `gorm:"type:timestamp null" json:"completed_at"`
	DurationMinutes      *int           `gorm:"type:int null" json:"duration_minutes"`
	
	// Staff
	CountedBy            *uint64        `gorm:"type:bigint unsigned null" json:"counted_by"`
	
	// Notes
	VarianceReason       string         `gorm:"type:text" json:"variance_reason"`
	
	// Timestamps
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName menentukan nama tabel di database
func (KhazwalCountingResult) TableName() string {
	return "khazwal_counting_results"
}

// IsPending memeriksa apakah counting masih pending
func (kcr *KhazwalCountingResult) IsPending() bool {
	return kcr.Status == CountingPending
}

// IsInProgress memeriksa apakah counting sedang dalam progress
func (kcr *KhazwalCountingResult) IsInProgress() bool {
	return kcr.Status == CountingInProgress
}

// IsCompleted memeriksa apakah counting sudah selesai
func (kcr *KhazwalCountingResult) IsCompleted() bool {
	return kcr.Status == CountingCompleted
}

// CalculateDuration menghitung durasi dalam menit antara started dan completed
func (kcr *KhazwalCountingResult) CalculateDuration() int {
	if kcr.StartedAt == nil || kcr.CompletedAt == nil {
		return 0
	}
	duration := kcr.CompletedAt.Sub(*kcr.StartedAt)
	return int(duration.Minutes())
}

// UpdateDuration mengupdate DurationMinutes field
func (kcr *KhazwalCountingResult) UpdateDuration() {
	duration := kcr.CalculateDuration()
	kcr.DurationMinutes = &duration
}

// CalculateVariance menghitung variance dari target quantity
func (kcr *KhazwalCountingResult) CalculateVariance(targetQuantity int) int {
	return kcr.TotalCounted - targetQuantity
}

// UpdateVariance mengupdate VarianceFromTarget field
func (kcr *KhazwalCountingResult) UpdateVariance(targetQuantity int) {
	variance := kcr.CalculateVariance(targetQuantity)
	kcr.VarianceFromTarget = &variance
}

// CalculatePercentages menghitung percentage good dan defect
func (kcr *KhazwalCountingResult) CalculatePercentages() {
	if kcr.TotalCounted == 0 {
		return
	}
	
	percentageGood := float64(kcr.QuantityGood) / float64(kcr.TotalCounted) * 100
	percentageDefect := float64(kcr.QuantityDefect) / float64(kcr.TotalCounted) * 100
	
	kcr.PercentageGood = &percentageGood
	kcr.PercentageDefect = &percentageDefect
}

// HasDefectBreakdown memeriksa apakah defect breakdown sudah diisi
func (kcr *KhazwalCountingResult) HasDefectBreakdown() bool {
	return kcr.DefectBreakdown != nil && len(kcr.DefectBreakdown) > 0
}

// HasVariance memeriksa apakah ada variance dari target
func (kcr *KhazwalCountingResult) HasVariance() bool {
	if kcr.VarianceFromTarget == nil {
		return false
	}
	return *kcr.VarianceFromTarget != 0
}

// DefectPercentageExceedsThreshold memeriksa apakah persentase rusak melebihi threshold
func (kcr *KhazwalCountingResult) DefectPercentageExceedsThreshold(threshold float64) bool {
	if kcr.PercentageDefect == nil {
		return false
	}
	return *kcr.PercentageDefect > threshold
}

// Request & Response DTOs

// QueueItemResponse merupakan response DTO untuk item di counting queue
type QueueItemResponse struct {
	POID              uint64    `json:"po_id"`
	PONumber          int64     `json:"po_number"`
	OBCNumber         string    `json:"obc_number"`
	TargetQuantity    int       `json:"target_quantity"`
	PrintCompletedAt  time.Time `json:"print_completed_at"`
	WaitingMinutes    int       `json:"waiting_minutes"`
	IsOverdue         bool      `json:"is_overdue"`
	Machine           *MachineInfo `json:"machine,omitempty"`
	Operator          *OperatorInfo `json:"operator,omitempty"`
}

// MachineInfo merupakan info mesin cetak
type MachineInfo struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// OperatorInfo merupakan info operator cetak
type OperatorInfo struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	NIP  string `json:"nip"`
}

// CountingDetailResponse merupakan response DTO untuk detail counting
type CountingDetailResponse struct {
	ID                   uint64                 `json:"id"`
	ProductionOrderID    uint64                 `json:"production_order_id"`
	Status               CountingStatus         `json:"status"`
	StartedAt            *time.Time             `json:"started_at"`
	QuantityGood         int                    `json:"quantity_good"`
	QuantityDefect       int                    `json:"quantity_defect"`
	TotalCounted         int                    `json:"total_counted"`
	VarianceFromTarget   *int                   `json:"variance_from_target"`
	PercentageGood       *float64               `json:"percentage_good"`
	PercentageDefect     *float64               `json:"percentage_defect"`
	DefectBreakdown      []DefectBreakdownItem  `json:"defect_breakdown"`
	VarianceReason       string                 `json:"variance_reason"`
	CountedBy            *OperatorInfo          `json:"counted_by,omitempty"`
	PO                   *POInfo                `json:"po,omitempty"`
	PrintInfo            *PrintInfo             `json:"print_info,omitempty"`
}

// POInfo merupakan info Production Order
type POInfo struct {
	PONumber       int64  `json:"po_number"`
	OBCNumber      string `json:"obc_number"`
	TargetQuantity int    `json:"target_quantity"`
}

// PrintInfo merupakan info print job summary
type PrintInfo struct {
	MachineName  string    `json:"machine_name"`
	OperatorName string    `json:"operator_name"`
	FinalizedAt  time.Time `json:"finalized_at"`
}

// StartCountingRequest merupakan request DTO untuk start counting
type StartCountingRequest struct {
	// Empty body, user dari auth token
}

// UpdateResultRequest merupakan request DTO untuk update counting results
type UpdateResultRequest struct {
	QuantityGood    int                   `json:"quantity_good" binding:"required,min=0"`
	QuantityDefect  int                   `json:"quantity_defect" binding:"required,min=0"`
	DefectBreakdown []DefectBreakdownItem `json:"defect_breakdown"`
	VarianceReason  string                `json:"variance_reason"`
}

// FinalizeCountingRequest merupakan request DTO untuk finalize counting
type FinalizeCountingRequest struct {
	// Empty body, all data sudah ada di counting record
}

// QueueResponse merupakan response DTO untuk queue endpoint
type QueueResponse struct {
	Data []QueueItemResponse `json:"data"`
	Meta QueueMetadata       `json:"meta"`
}

// QueueMetadata merupakan metadata untuk queue response
type QueueMetadata struct {
	Total         int `json:"total"`
	OverdueCount  int `json:"overdue_count"`
}

// StartCountingResponse merupakan response DTO untuk start counting
type StartCountingResponse struct {
	ID                uint64    `json:"id"`
	ProductionOrderID uint64    `json:"production_order_id"`
	Status            string    `json:"status"`
	StartedAt         time.Time `json:"started_at"`
	CountedBy         uint64    `json:"counted_by"`
}

// UpdateResultResponse merupakan response DTO untuk update result
type UpdateResultResponse struct {
	ID                 uint64   `json:"id"`
	QuantityGood       int      `json:"quantity_good"`
	QuantityDefect     int      `json:"quantity_defect"`
	TotalCounted       int      `json:"total_counted"`
	VarianceFromTarget *int     `json:"variance_from_target"`
	PercentageGood     *float64 `json:"percentage_good"`
	PercentageDefect   *float64 `json:"percentage_defect"`
	DefectBreakdown    []DefectBreakdownItem `json:"defect_breakdown"`
	VarianceReason     string   `json:"variance_reason"`
}

// FinalizeCountingResponse merupakan response DTO untuk finalize counting
type FinalizeCountingResponse struct {
	ID              uint64    `json:"id"`
	Status          string    `json:"status"`
	CompletedAt     time.Time `json:"completed_at"`
	DurationMinutes int       `json:"duration_minutes"`
}
