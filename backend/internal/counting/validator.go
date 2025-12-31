package counting

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ValidationErrors merupakan collection of validation errors
var (
	ErrInvalidQuantity           = errors.New("quantity_good dan quantity_defect harus >= 0")
	ErrDefectBreakdownRequired   = errors.New("defect_breakdown wajib diisi karena persentase rusak > 5%")
	ErrDefectBreakdownSumMismatch = errors.New("total defect_breakdown harus sama dengan quantity_defect")
	ErrVarianceReasonRequired    = errors.New("variance_reason wajib diisi karena ada selisih dari target")
	ErrCountingNotInProgress     = errors.New("counting tidak dalam status IN_PROGRESS")
	ErrCountingAlreadyCompleted  = errors.New("counting sudah selesai dan tidak bisa diubah")
	ErrRequiredFieldsMissing     = errors.New("field quantity_good dan quantity_defect wajib diisi")
	ErrPONotReadyForCounting     = errors.New("PO belum siap untuk penghitungan")
	ErrCountingAlreadyExists     = errors.New("counting untuk PO ini sudah ada")
)

// Constants untuk validation thresholds
const (
	DefectBreakdownThreshold = 5.0  // 5% - jika rusak > 5%, wajib breakdown
	ToleranceThreshold       = 2.0  // 2% - untuk warning threshold
	MaxWaitingMinutesWarning = 120  // 2 jam - untuk overdue warning
)

// ValidateUpdateResultRequest memvalidasi request untuk update counting result
// dengan business rules: defect breakdown required jika > 5%, variance reason required jika != 0
func ValidateUpdateResultRequest(req UpdateResultRequest, targetQuantity int) error {
	// Validate quantities
	if req.QuantityGood < 0 || req.QuantityDefect < 0 {
		return ErrInvalidQuantity
	}

	totalCounted := req.QuantityGood + req.QuantityDefect
	if totalCounted == 0 {
		return nil // Allow saving zero initially
	}

	// Calculate percentages
	percentageDefect := float64(req.QuantityDefect) / float64(totalCounted) * 100
	
	// Validate defect breakdown requirement (> 5%)
	if percentageDefect > DefectBreakdownThreshold {
		if len(req.DefectBreakdown) == 0 {
			return ErrDefectBreakdownRequired
		}
		
		// Validate breakdown sum
		if err := ValidateDefectBreakdownSum(req.DefectBreakdown, req.QuantityDefect); err != nil {
			return err
		}
	}

	// Validate variance reason requirement
	variance := totalCounted - targetQuantity
	if variance != 0 && req.VarianceReason == "" {
		return ErrVarianceReasonRequired
	}

	return nil
}

// ValidateDefectBreakdownSum memvalidasi bahwa sum dari breakdown sama dengan quantity_defect
func ValidateDefectBreakdownSum(breakdown []DefectBreakdownItem, quantityDefect int) error {
	sum := 0
	for _, item := range breakdown {
		sum += item.Quantity
	}
	
	if sum != quantityDefect {
		return fmt.Errorf("%w: expected %d, got %d", ErrDefectBreakdownSumMismatch, quantityDefect, sum)
	}
	
	return nil
}

// ValidateFinalizeRequirements memvalidasi bahwa semua required fields sudah diisi untuk finalize
func ValidateFinalizeRequirements(counting *KhazwalCountingResult, targetQuantity int) error {
	// Check status
	if !counting.IsInProgress() {
		return ErrCountingNotInProgress
	}

	// Check required fields
	if counting.QuantityGood == 0 && counting.QuantityDefect == 0 {
		return ErrRequiredFieldsMissing
	}

	// Check defect breakdown requirement
	if counting.PercentageDefect != nil && *counting.PercentageDefect > DefectBreakdownThreshold {
		if !counting.HasDefectBreakdown() {
			return ErrDefectBreakdownRequired
		}
		
		// Validate breakdown sum
		var breakdown []DefectBreakdownItem
		if err := json.Unmarshal(counting.DefectBreakdown, &breakdown); err == nil {
			if err := ValidateDefectBreakdownSum(breakdown, counting.QuantityDefect); err != nil {
				return err
			}
		}
	}

	// Check variance reason requirement
	variance := counting.TotalCounted - targetQuantity
	if variance != 0 && counting.VarianceReason == "" {
		return ErrVarianceReasonRequired
	}

	return nil
}

// ValidateCountingStatus memvalidasi status counting untuk update operations
func ValidateCountingStatus(counting *KhazwalCountingResult) error {
	if counting.IsCompleted() {
		return ErrCountingAlreadyCompleted
	}
	
	if !counting.IsInProgress() {
		return ErrCountingNotInProgress
	}
	
	return nil
}

// CalculateWaitingMinutes menghitung waktu tunggu dalam menit dari print completion
func CalculateWaitingMinutes(printCompletedAt *time.Time) int {
	if printCompletedAt == nil {
		return 0
	}
	duration := time.Since(*printCompletedAt)
	return int(duration.Minutes())
}

// IsOverdue memeriksa apakah PO sudah overdue (> 2 jam menunggu)
func IsOverdue(waitingMinutes int) bool {
	return waitingMinutes > MaxWaitingMinutesWarning
}

// ParseDefectBreakdown mem-parse JSON defect breakdown ke slice of DefectBreakdownItem
func ParseDefectBreakdown(jsonData []byte) ([]DefectBreakdownItem, error) {
	if jsonData == nil || len(jsonData) == 0 {
		return []DefectBreakdownItem{}, nil
	}

	var breakdown []DefectBreakdownItem
	if err := json.Unmarshal(jsonData, &breakdown); err != nil {
		return nil, fmt.Errorf("gagal parse defect_breakdown: %w", err)
	}

	return breakdown, nil
}

// SerializeDefectBreakdown men-serialize slice of DefectBreakdownItem ke JSON
func SerializeDefectBreakdown(breakdown []DefectBreakdownItem) ([]byte, error) {
	if len(breakdown) == 0 {
		return nil, nil
	}

	data, err := json.Marshal(breakdown)
	if err != nil {
		return nil, fmt.Errorf("gagal serialize defect_breakdown: %w", err)
	}

	return data, nil
}
