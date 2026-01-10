package cutting

import (
	"errors"
	"time"
)

// Custom errors untuk cutting operations
var (
	ErrCuttingNotFound            = errors.New("cutting record not found")
	ErrPONotFound                 = errors.New("production order not found")
	ErrPONotReadyForCutting       = errors.New("PO not in SIAP_POTONG status")
	ErrCuttingAlreadyStarted      = errors.New("cutting already started for this PO")
	ErrCuttingNotInProgress       = errors.New("cutting not in progress")
	ErrMissingOutputData          = errors.New("output sisiran kiri & kanan must be filled")
	ErrMissingWasteDocumentation  = errors.New("waste > 2% requires reason and photo")
	ErrInvalidWasteData           = errors.New("invalid waste data")
	ErrCuttingAlreadyCompleted    = errors.New("cutting already completed")
	ErrCountingNotCompleted       = errors.New("counting result not completed yet")
)

// Request & Response DTOs

// QueueItemResponse merupakan response DTO untuk item di cutting queue
type QueueItemResponse struct {
	POID                 uint64    `json:"po_id"`
	PONumber             int64     `json:"po_number"`
	OBCNumber            string    `json:"obc_number"`
	Priority             string    `json:"priority"`
	InputLembarBesar     int       `json:"input_lembar_besar"`
	EstimatedOutput      int       `json:"estimated_output"`
	CountingCompletedAt  time.Time `json:"counting_completed_at"`
	WaitingMinutes       int       `json:"waiting_minutes"`
	IsOverdue            bool      `json:"is_overdue"`
}

// QueueFilters merupakan filter parameters untuk queue endpoint
type QueueFilters struct {
	Priority  string `form:"priority"`
	DateFrom  string `form:"date_from"`
	DateTo    string `form:"date_to"`
	SortBy    string `form:"sort_by"`    // priority, date
	SortOrder string `form:"sort_order"` // asc, desc
}

// QueueResponse merupakan response DTO untuk queue endpoint
type QueueResponse struct {
	Data []QueueItemResponse `json:"data"`
	Meta QueueMetadata       `json:"meta"`
}

// QueueMetadata merupakan metadata untuk queue response
type QueueMetadata struct {
	Total        int `json:"total"`
	UrgentCount  int `json:"urgent_count"`
	NormalCount  int `json:"normal_count"`
}

// StartCuttingRequest merupakan request DTO untuk start cutting
type StartCuttingRequest struct {
	CuttingMachine string `json:"cutting_machine" binding:"required"`
	// CutBy diambil dari auth user
}

// StartCuttingResponse merupakan response DTO untuk start cutting
type StartCuttingResponse struct {
	ID                uint64    `json:"id"`
	ProductionOrderID uint64    `json:"production_order_id"`
	InputLembarBesar  int       `json:"input_lembar_besar"`
	ExpectedOutput    int       `json:"expected_output"`
	CuttingMachine    string    `json:"cutting_machine"`
	Status            string    `json:"status"`
	StartedAt         time.Time `json:"started_at"`
	CutBy             uint64    `json:"cut_by"`
}

// UpdateResultRequest merupakan request DTO untuk update cutting results
type UpdateResultRequest struct {
	OutputSisiranKiri  int    `json:"output_sisiran_kiri" binding:"required,min=0"`
	OutputSisiranKanan int    `json:"output_sisiran_kanan" binding:"required,min=0"`
	WasteReason        string `json:"waste_reason"`
	WastePhotoURL      string `json:"waste_photo_url"`
}

// UpdateResultResponse merupakan response DTO untuk update result
type UpdateResultResponse struct {
	ID                 uint64   `json:"id"`
	OutputSisiranKiri  int      `json:"output_sisiran_kiri"`
	OutputSisiranKanan int      `json:"output_sisiran_kanan"`
	TotalOutput        int      `json:"total_output"`
	ExpectedOutput     int      `json:"expected_output"`
	WasteQuantity      int      `json:"waste_quantity"`
	WastePercentage    *float64 `json:"waste_percentage"`
	WasteReason        string   `json:"waste_reason"`
	WastePhotoURL      string   `json:"waste_photo_url"`
}

// FinalizeCuttingRequest merupakan request DTO untuk finalize cutting
type FinalizeCuttingRequest struct {
	// Empty body, all data sudah ada di cutting record
}

// FinalizeCuttingResponse merupakan response DTO untuk finalize cutting
type FinalizeCuttingResponse struct {
	ID              uint64    `json:"id"`
	Status          string    `json:"status"`
	CompletedAt     time.Time `json:"completed_at"`
	DurationMinutes int       `json:"duration_minutes"`
	LabelsGenerated int       `json:"labels_generated"`
}

// CuttingDetailResponse merupakan response DTO untuk detail cutting
type CuttingDetailResponse struct {
	ID                  uint64        `json:"id"`
	ProductionOrderID   uint64        `json:"production_order_id"`
	Status              CuttingStatus `json:"status"`
	StartedAt           *time.Time    `json:"started_at"`
	CompletedAt         *time.Time    `json:"completed_at"`
	InputLembarBesar    int           `json:"input_lembar_besar"`
	ExpectedOutput      int           `json:"expected_output"`
	OutputSisiranKiri   *int          `json:"output_sisiran_kiri"`
	OutputSisiranKanan  *int          `json:"output_sisiran_kanan"`
	TotalOutput         int           `json:"total_output"`
	WasteQuantity       int           `json:"waste_quantity"`
	WastePercentage     *float64      `json:"waste_percentage"`
	WasteReason         string        `json:"waste_reason"`
	WastePhotoURL       string        `json:"waste_photo_url"`
	CuttingMachine      string        `json:"cutting_machine"`
	CutBy               *OperatorInfo `json:"cut_by,omitempty"`
	PO                  *POInfo       `json:"po,omitempty"`
}

// POInfo merupakan info Production Order
type POInfo struct {
	PONumber       int64  `json:"po_number"`
	OBCNumber      string `json:"obc_number"`
	Priority       string `json:"priority"`
	TargetQuantity int    `json:"target_quantity"`
}

// OperatorInfo merupakan info operator/staff
type OperatorInfo struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	NIP  string `json:"nip"`
}
