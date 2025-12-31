package counting

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CountingService merupakan interface untuk business logic counting operations
type CountingService interface {
	// Queue operations
	GetCountingQueue(machineID *uint64, dateFrom, dateTo *time.Time) (*QueueResponse, error)
	
	// Detail operations
	GetCountingDetail(id uint64) (*CountingDetailResponse, error)
	GetCountingDetailByPOID(poID uint64) (*CountingDetailResponse, error)
	
	// Action operations
	StartCounting(poID uint64, userID uint64) (*StartCountingResponse, error)
	UpdateResult(id uint64, req UpdateResultRequest, targetQuantity int) (*UpdateResultResponse, error)
	FinalizeCounting(id uint64) (*FinalizeCountingResponse, error)
}

// countingServiceImpl merupakan implementasi CountingService
type countingServiceImpl struct {
	db   *gorm.DB
	repo CountingRepository
}

// NewCountingService membuat instance baru CountingService
func NewCountingService(db *gorm.DB, repo CountingRepository) CountingService {
	return &countingServiceImpl{
		db:   db,
		repo: repo,
	}
}

// GetCountingQueue mengambil list PO yang menunggu penghitungan
func (s *countingServiceImpl) GetCountingQueue(machineID *uint64, dateFrom, dateTo *time.Time) (*QueueResponse, error) {
	// Get queue items dari repository
	items, err := s.repo.GetCountingQueue(machineID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	// Calculate metadata
	total := len(items)
	overdueCount := 0
	for _, item := range items {
		if item.IsOverdue {
			overdueCount++
		}
	}

	response := &QueueResponse{
		Data: items,
		Meta: QueueMetadata{
			Total:        total,
			OverdueCount: overdueCount,
		},
	}

	return response, nil
}

// GetCountingDetail mengambil detail counting record dengan relasi
func (s *countingServiceImpl) GetCountingDetail(id uint64) (*CountingDetailResponse, error) {
	return s.repo.GetCountingDetailWithRelations(id)
}

// GetCountingDetailByPOID mengambil detail counting berdasarkan PO ID
func (s *countingServiceImpl) GetCountingDetailByPOID(poID uint64) (*CountingDetailResponse, error) {
	// Get counting record by PO ID
	counting, err := s.repo.GetByPOID(poID)
	if err != nil {
		return nil, err
	}

	// Get full detail dengan relasi
	return s.repo.GetCountingDetailWithRelations(counting.ID)
}

// StartCounting memulai proses penghitungan untuk PO tertentu
// dengan validasi PO status dan create counting record
func (s *countingServiceImpl) StartCounting(poID uint64, userID uint64) (*StartCountingResponse, error) {
	// Start transaction untuk atomicity
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Validate PO exists dan status correct
	var po struct {
		ID            uint64
		CurrentStatus string
	}
	if err := tx.Table("production_orders").
		Select("id, current_status").
		Where("id = ?", poID).
		First(&po).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("PO tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal validasi PO: %w", err)
	}

	// Check PO status
	if po.CurrentStatus != "WAITING_COUNTING" {
		tx.Rollback()
		return nil, ErrPONotReadyForCounting
	}

	// 2. Check tidak ada counting IN_PROGRESS lain untuk PO ini
	exists, err := s.repo.ExistsInProgressByPOID(poID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if exists {
		tx.Rollback()
		return nil, ErrCountingAlreadyExists
	}

	// 3. Create counting record
	now := time.Now()
	counting := &KhazwalCountingResult{
		ProductionOrderID: poID,
		Status:            CountingInProgress,
		StartedAt:         &now,
		CountedBy:         &userID,
	}

	if err := tx.Create(counting).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal membuat counting record: %w", err)
	}

	// 4. Update PO status ke SEDANG_DIHITUNG
	if err := tx.Table("production_orders").
		Where("id = ?", poID).
		Updates(map[string]interface{}{
			"current_status": "SEDANG_DIHITUNG",
			"updated_at":     time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal update PO status: %w", err)
	}

	// 5. Update po_stage_tracking set started_at untuk stage KHAZWAL_COUNTING
	if err := tx.Table("po_stage_tracking").
		Where("production_order_id = ?", poID).
		Where("stage = ?", "KHAZWAL_COUNTING").
		Update("started_at", now).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal update stage tracking: %w", err)
	}

	// 6. Log activity
	activityLog := map[string]interface{}{
		"user_id":             userID,
		"action":              "START_COUNTING",
		"entity_type":         "COUNTING",
		"entity_id":           counting.ID,
		"production_order_id": poID,
		"description":         fmt.Sprintf("Memulai penghitungan untuk PO %d", poID),
		"created_at":          time.Now(),
	}
	if err := tx.Table("activity_logs").Create(activityLog).Error; err != nil {
		// Log error tapi tidak rollback (non-critical)
		fmt.Printf("Warning: gagal log activity: %v\n", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("gagal commit transaction: %w", err)
	}

	// Return response
	response := &StartCountingResponse{
		ID:                counting.ID,
		ProductionOrderID: poID,
		Status:            string(CountingInProgress),
		StartedAt:         now,
		CountedBy:         userID,
	}

	return response, nil
}

// UpdateResult mengupdate hasil penghitungan (dapat dipanggil multiple times sebelum finalize)
func (s *countingServiceImpl) UpdateResult(id uint64, req UpdateResultRequest, targetQuantity int) (*UpdateResultResponse, error) {
	// 1. Get counting record
	counting, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 2. Validate status (must be IN_PROGRESS)
	if err := ValidateCountingStatus(counting); err != nil {
		return nil, err
	}

	// 3. Validate request dengan business rules
	if err := ValidateUpdateResultRequest(req, targetQuantity); err != nil {
		return nil, err
	}

	// 4. Update counting record
	counting.QuantityGood = req.QuantityGood
	counting.QuantityDefect = req.QuantityDefect
	counting.VarianceReason = req.VarianceReason

	// Calculate auto fields
	counting.UpdateVariance(targetQuantity)
	counting.CalculatePercentages()

	// Serialize defect breakdown
	if len(req.DefectBreakdown) > 0 {
		breakdownJSON, err := SerializeDefectBreakdown(req.DefectBreakdown)
		if err != nil {
			return nil, err
		}
		counting.DefectBreakdown = breakdownJSON
	}

	// 5. Save to database
	if err := s.repo.Update(counting); err != nil {
		return nil, err
	}

	// 6. Return response
	response := &UpdateResultResponse{
		ID:                 counting.ID,
		QuantityGood:       counting.QuantityGood,
		QuantityDefect:     counting.QuantityDefect,
		TotalCounted:       counting.TotalCounted,
		VarianceFromTarget: counting.VarianceFromTarget,
		PercentageGood:     counting.PercentageGood,
		PercentageDefect:   counting.PercentageDefect,
		DefectBreakdown:    req.DefectBreakdown,
		VarianceReason:     counting.VarianceReason,
	}

	return response, nil
}

// FinalizeCounting menyelesaikan penghitungan dengan lock data dan advance PO ke next stage
func (s *countingServiceImpl) FinalizeCounting(id uint64) (*FinalizeCountingResponse, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Get counting record dengan lock
	var counting KhazwalCountingResult
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&counting, id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("counting record tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil counting record: %w", err)
	}

	// 2. Get PO untuk target quantity
	var po struct {
		ID                       uint64
		QuantityTargetLembarBesar int
	}
	if err := tx.Table("production_orders").
		Select("id, quantity_target_lembar_besar").
		Where("id = ?", counting.ProductionOrderID).
		First(&po).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal mengambil PO: %w", err)
	}

	// 3. Validate all requirements
	if err := ValidateFinalizeRequirements(&counting, po.QuantityTargetLembarBesar); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 4. Update counting record
	now := time.Now()
	counting.Status = CountingCompleted
	counting.CompletedAt = &now
	counting.UpdateDuration()

	if err := tx.Save(&counting).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal finalize counting: %w", err)
	}

	// 5. Update PO status dan stage
	if err := tx.Table("production_orders").
		Where("id = ?", counting.ProductionOrderID).
		Updates(map[string]interface{}{
			"current_stage":  "KHAZWAL_CUTTING",
			"current_status": "SIAP_POTONG",
			"updated_at":     time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal update PO status: %w", err)
	}

	// 6. Update po_stage_tracking set completed_at untuk stage KHAZWAL_COUNTING
	if err := tx.Table("po_stage_tracking").
		Where("production_order_id = ?", counting.ProductionOrderID).
		Where("stage = ?", "KHAZWAL_COUNTING").
		Update("completed_at", now).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal update stage tracking: %w", err)
	}

	// 7. Log activity (immutable)
	activityLog := map[string]interface{}{
		"user_id":             counting.CountedBy,
		"action":              "FINALIZE_COUNTING",
		"entity_type":         "COUNTING",
		"entity_id":           counting.ID,
		"production_order_id": counting.ProductionOrderID,
		"description":         fmt.Sprintf("Menyelesaikan penghitungan untuk PO %d", counting.ProductionOrderID),
		"metadata":            buildFinalizeMetadata(&counting),
		"created_at":          time.Now(),
	}
	if err := tx.Table("activity_logs").Create(activityLog).Error; err != nil {
		// Log error tapi tidak rollback
		fmt.Printf("Warning: gagal log activity: %v\n", err)
	}

	// 8. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("gagal commit transaction: %w", err)
	}

	// 9. Return response
	response := &FinalizeCountingResponse{
		ID:              counting.ID,
		Status:          string(CountingCompleted),
		CompletedAt:     *counting.CompletedAt,
		DurationMinutes: *counting.DurationMinutes,
	}

	return response, nil
}

// buildFinalizeMetadata membuat metadata untuk activity log finalize
func buildFinalizeMetadata(counting *KhazwalCountingResult) string {
	metadata := map[string]interface{}{
		"quantity_good":       counting.QuantityGood,
		"quantity_defect":     counting.QuantityDefect,
		"total_counted":       counting.TotalCounted,
		"variance_from_target": counting.VarianceFromTarget,
		"percentage_good":     counting.PercentageGood,
		"percentage_defect":   counting.PercentageDefect,
		"duration_minutes":    counting.DurationMinutes,
	}

	// Parse defect breakdown if exists
	if counting.HasDefectBreakdown() {
		var breakdown []DefectBreakdownItem
		if err := json.Unmarshal(counting.DefectBreakdown, &breakdown); err == nil {
			metadata["defect_breakdown"] = breakdown
		}
	}

	if counting.HasVariance() {
		metadata["variance_reason"] = counting.VarianceReason
	}

	jsonData, _ := json.Marshal(metadata)
	return string(jsonData)
}
