package models

import (
	"time"

	"gorm.io/gorm"
)

// POPriority merupakan enum untuk priority level Production Order
type POPriority string

const (
	PriorityUrgent POPriority = "URGENT"
	PriorityNormal POPriority = "NORMAL"
	PriorityLow    POPriority = "LOW"
)

// POStage merupakan enum untuk stage Production Order
type POStage string

const (
	StageKhazwalMaterialPrep POStage = "KHAZWAL_MATERIAL_PREP"
	StageCetak               POStage = "CETAK"
	StageVerifikasi          POStage = "VERIFIKASI"
	StageKhazkhir            POStage = "KHAZKHIR"
	StageCompleted           POStage = "COMPLETED"
)

// POStatus merupakan enum untuk status Production Order
type POStatus string

const (
	StatusWaitingMaterialPrep POStatus = "WAITING_MATERIAL_PREP"
	StatusMaterialPrepInProgress POStatus = "MATERIAL_PREP_IN_PROGRESS"
	StatusReadyForCetak       POStatus = "READY_FOR_CETAK"
	StatusCetakInProgress     POStatus = "CETAK_IN_PROGRESS"
	StatusReadyForVerifikasi  POStatus = "READY_FOR_VERIFIKASI"
	StatusVerifikasiInProgress POStatus = "VERIFIKASI_IN_PROGRESS"
	StatusReadyForKhazkhir    POStatus = "READY_FOR_KHAZKHIR"
	StatusKhazkhirInProgress  POStatus = "KHAZKHIR_IN_PROGRESS"
	StatusPOCompleted         POStatus = "COMPLETED"
)

// ProductionOrder merupakan model untuk entitas Production Order
// yang mencakup informasi PO dari SAP dan tracking status.
// Data spesifikasi OBC sekarang di-reference dari OBCMaster, dengan denormalisasi beberapa field untuk performance
type ProductionOrder struct {
	ID                        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PONumber                  int64          `gorm:"uniqueIndex;not null" json:"po_number" binding:"required"`
	OBCMasterID               uint64         `gorm:"index;not null" json:"obc_master_id" binding:"required"`
	
	// Denormalized fields dari OBCMaster untuk performance
	OBCNumber                 string         `gorm:"type:varchar(20);index" json:"obc_number"`
	ProductName               string         `gorm:"type:varchar(255)" json:"product_name"`
	SAPCustomerCode           string         `gorm:"type:varchar(50)" json:"sap_customer_code"`
	SAPProductCode            string         `gorm:"type:varchar(50)" json:"sap_product_code"`
	ProductSpecifications     interface{}    `gorm:"type:json" json:"product_specifications"`
	
	QuantityOrdered           int            `gorm:"not null" json:"quantity_ordered" binding:"required,min=1"`
	QuantityTargetLembarBesar int            `gorm:"not null" json:"quantity_target_lembar_besar" binding:"required,min=1"`
	EstimatedRims             int            `gorm:"not null" json:"estimated_rims" binding:"required,min=1"`
	OrderDate                 time.Time      `gorm:"type:date;not null" json:"order_date" binding:"required"`
	DueDate                   time.Time      `gorm:"type:date;not null" json:"due_date" binding:"required"`
	Priority                  POPriority     `gorm:"type:enum('URGENT','NORMAL','LOW');default:'NORMAL'" json:"priority"`
	PriorityScore             int            `gorm:"default:50" json:"priority_score"`
	CurrentStage              POStage        `gorm:"type:enum('KHAZWAL_MATERIAL_PREP','CETAK','VERIFIKASI','KHAZKHIR','COMPLETED');default:'KHAZWAL_MATERIAL_PREP'" json:"current_stage"`
	CurrentStatus             POStatus       `gorm:"type:enum('WAITING_MATERIAL_PREP','MATERIAL_PREP_IN_PROGRESS','READY_FOR_CETAK','CETAK_IN_PROGRESS','READY_FOR_VERIFIKASI','VERIFIKASI_IN_PROGRESS','READY_FOR_KHAZKHIR','KHAZKHIR_IN_PROGRESS','COMPLETED');default:'WAITING_MATERIAL_PREP'" json:"current_status"`
	Notes                     string         `gorm:"type:text" json:"notes"`
	CreatedAt                 time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt                 time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt                 gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	OBCMaster           *OBCMaster                  `gorm:"foreignKey:OBCMasterID" json:"obc_master,omitempty"`
	KhazwalMaterialPrep *KhazwalMaterialPreparation `gorm:"foreignKey:ProductionOrderID" json:"khazwal_material_prep,omitempty"`
	StageTracking       []POStageTracking           `gorm:"foreignKey:ProductionOrderID" json:"stage_tracking,omitempty"`
}

// TableName menentukan nama tabel di database
func (ProductionOrder) TableName() string {
	return "production_orders"
}

// IsUrgent memeriksa apakah PO berstatus urgent
func (po *ProductionOrder) IsUrgent() bool {
	return po.Priority == PriorityUrgent
}

// IsPastDue memeriksa apakah PO sudah melewati due date
func (po *ProductionOrder) IsPastDue() bool {
	return time.Now().After(po.DueDate)
}

// DaysUntilDue menghitung sisa hari hingga due date
func (po *ProductionOrder) DaysUntilDue() int {
	duration := time.Until(po.DueDate)
	return int(duration.Hours() / 24)
}

// CalculatePriorityScore menghitung priority score berdasarkan priority dan due date
// dimana skor lebih tinggi = lebih urgent
func (po *ProductionOrder) CalculatePriorityScore() int {
	score := 50 // Base score

	// Priority multiplier
	switch po.Priority {
	case PriorityUrgent:
		score += 50
	case PriorityNormal:
		score += 20
	case PriorityLow:
		score += 0
	}

	// Due date factor
	daysUntilDue := po.DaysUntilDue()
	if daysUntilDue < 0 {
		// Past due - add penalty
		score += (daysUntilDue * -10)
	} else if daysUntilDue <= 3 {
		// Very close to due date
		score += 30
	} else if daysUntilDue <= 7 {
		// Close to due date
		score += 15
	}

	return score
}

// UpdatePriorityScore mengupdate priority score
func (po *ProductionOrder) UpdatePriorityScore() {
	po.PriorityScore = po.CalculatePriorityScore()
}
