package cutting

import (
	"fmt"
	"time"
)

// Service merupakan interface untuk business logic cutting
type Service interface {
	// Queue operations
	GetCuttingQueue(filters QueueFilters) (*QueueResponse, error)
	
	// Workflow operations
	StartCutting(poID uint64, req StartCuttingRequest, userID uint64) (*StartCuttingResponse, error)
	UpdateCuttingResult(id uint64, req UpdateResultRequest) (*UpdateResultResponse, error)
	FinalizeCutting(id uint64) (*FinalizeCuttingResponse, error)
	
	// Detail operations
	GetCuttingDetail(id uint64) (*CuttingDetailResponse, error)
}

// service merupakan implementasi konkret dari Service interface
type service struct {
	repo Repository
}

// NewService membuat instance baru dari service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetCuttingQueue mengambil list PO yang siap untuk dipotong dengan filters
func (s *service) GetCuttingQueue(filters QueueFilters) (*QueueResponse, error) {
	// Get queue data
	data, err := s.repo.GetCuttingQueue(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get cutting queue: %w", err)
	}
	
	// Get metadata
	meta, err := s.repo.GetQueueMetadata(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue metadata: %w", err)
	}
	
	return &QueueResponse{
		Data: data,
		Meta: meta,
	}, nil
}

// StartCutting memulai proses pemotongan untuk PO tertentu
func (s *service) StartCutting(poID uint64, req StartCuttingRequest, userID uint64) (*StartCuttingResponse, error) {
	// 1. Check if PO exists and in correct status
	_, err := s.repo.GetPOInfo(poID)
	if err != nil {
		return nil, err
	}
	
	// 2. Check if cutting already exists for this PO
	existing, err := s.repo.GetByPOID(poID)
	if err == nil && existing != nil {
		return nil, ErrCuttingAlreadyStarted
	}
	if err != nil && err != ErrCuttingNotFound {
		return nil, fmt.Errorf("failed to check existing cutting: %w", err)
	}
	
	// 3. Get counting result to get input quantity
	countingResult, err := s.repo.GetCountingResultByPOID(poID)
	if err != nil {
		return nil, err
	}
	
	if countingResult.Status != "COMPLETED" {
		return nil, ErrCountingNotCompleted
	}
	
	// 4. Create cutting record
	now := time.Now()
	inputLembarBesar := countingResult.QuantityGood
	expectedOutput := inputLembarBesar * 2
	
	cutting := &KhazwalCuttingResult{
		ProductionOrderID: poID,
		InputLembarBesar:  inputLembarBesar,
		ExpectedOutput:    expectedOutput,
		CuttingMachine:    req.CuttingMachine,
		CutBy:             &userID,
		Status:            CuttingInProgress,
		StartedAt:         &now,
	}
	
	err = s.repo.Create(cutting)
	if err != nil {
		return nil, fmt.Errorf("failed to create cutting record: %w", err)
	}
	
	// 5. Update PO status to SEDANG_DIPOTONG
	err = s.repo.UpdatePOStatus(poID, "KHAZWAL_CUTTING", "SEDANG_DIPOTONG")
	if err != nil {
		return nil, fmt.Errorf("failed to update PO status: %w", err)
	}
	
	// 6. Update po_stage_tracking with started_at
	err = s.repo.UpdatePOStageTracking(poID, "started_at", now)
	if err != nil {
		return nil, fmt.Errorf("failed to update stage tracking: %w", err)
	}
	
	// 7. Build response
	return &StartCuttingResponse{
		ID:                cutting.ID,
		ProductionOrderID: cutting.ProductionOrderID,
		InputLembarBesar:  cutting.InputLembarBesar,
		ExpectedOutput:    cutting.ExpectedOutput,
		CuttingMachine:    cutting.CuttingMachine,
		Status:            string(cutting.Status),
		StartedAt:         *cutting.StartedAt,
		CutBy:             userID,
	}, nil
}

// UpdateCuttingResult mengupdate hasil pemotongan (sisiran kiri & kanan)
func (s *service) UpdateCuttingResult(id uint64, req UpdateResultRequest) (*UpdateResultResponse, error) {
	// 1. Get existing cutting record
	cutting, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// 2. Validate status (must be IN_PROGRESS)
	if !cutting.IsInProgress() {
		return nil, ErrCuttingNotInProgress
	}
	
	// 3. Update output values
	cutting.OutputSisiranKiri = &req.OutputSisiranKiri
	cutting.OutputSisiranKanan = &req.OutputSisiranKanan
	
	// 4. Calculate total, waste, and percentage
	cutting.UpdateTotalOutput()
	cutting.UpdateWaste()
	cutting.UpdateWastePercentage()
	
	// 5. Update waste documentation if provided
	if req.WasteReason != "" {
		cutting.WasteReason = req.WasteReason
	}
	if req.WastePhotoURL != "" {
		cutting.WastePhotoURL = req.WastePhotoURL
	}
	
	// 6. Save to database
	err = s.repo.Update(cutting)
	if err != nil {
		return nil, fmt.Errorf("failed to update cutting result: %w", err)
	}
	
	// 7. Build response
	return &UpdateResultResponse{
		ID:                 cutting.ID,
		OutputSisiranKiri:  *cutting.OutputSisiranKiri,
		OutputSisiranKanan: *cutting.OutputSisiranKanan,
		TotalOutput:        cutting.TotalOutput,
		ExpectedOutput:     cutting.ExpectedOutput,
		WasteQuantity:      cutting.WasteQuantity,
		WastePercentage:    cutting.WastePercentage,
		WasteReason:        cutting.WasteReason,
		WastePhotoURL:      cutting.WastePhotoURL,
	}, nil
}

// FinalizeCutting menyelesaikan proses pemotongan dan generate verification labels
func (s *service) FinalizeCutting(id uint64) (*FinalizeCuttingResponse, error) {
	// 1. Get existing cutting record
	cutting, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// 2. Validate status
	if cutting.IsCompleted() {
		return nil, ErrCuttingAlreadyCompleted
	}
	
	if !cutting.IsInProgress() {
		return nil, ErrCuttingNotInProgress
	}
	
	// 3. Validate data completeness
	err = cutting.ValidateForFinalization()
	if err != nil {
		return nil, err
	}
	
	// 4. Update status and timestamps
	now := time.Now()
	cutting.Status = CuttingCompleted
	cutting.CompletedAt = &now
	cutting.UpdateDuration()
	
	// 5. Save cutting record
	err = s.repo.Update(cutting)
	if err != nil {
		return nil, fmt.Errorf("failed to finalize cutting: %w", err)
	}
	
	// 6. Update PO status to SIAP_VERIFIKASI
	err = s.repo.UpdatePOStatus(cutting.ProductionOrderID, "VERIFIKASI", "SIAP_VERIFIKASI")
	if err != nil {
		return nil, fmt.Errorf("failed to update PO status: %w", err)
	}
	
	// 7. Update po_stage_tracking with completed_at
	err = s.repo.UpdatePOStageTracking(cutting.ProductionOrderID, "completed_at", now)
	if err != nil {
		return nil, fmt.Errorf("failed to update stage tracking: %w", err)
	}
	
	// 8. Generate verification labels (Sprint 3 feature)
	// Total labels = CEIL(total_output / 500)
	labelsGenerated := (cutting.TotalOutput + 499) / 500 // Integer ceiling division
	
	// TODO: Insert ke verification_labels table (Sprint 3)
	// For now, just return the count
	
	// 9. Create notification to Tim Verifikasi (Sprint 3 feature)
	// TODO: Implement notification creation
	
	// 10. Build response
	return &FinalizeCuttingResponse{
		ID:              cutting.ID,
		Status:          string(cutting.Status),
		CompletedAt:     *cutting.CompletedAt,
		DurationMinutes: *cutting.DurationMinutes,
		LabelsGenerated: labelsGenerated,
	}, nil
}

// GetCuttingDetail mengambil detail cutting record
func (s *service) GetCuttingDetail(id uint64) (*CuttingDetailResponse, error) {
	// 1. Get cutting record
	cutting, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// 2. Get PO info
	poInfo, err := s.repo.GetPOInfo(cutting.ProductionOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PO info: %w", err)
	}
	
	// 3. Get operator info if available
	var operatorInfo *OperatorInfo
	if cutting.CutBy != nil {
		operatorInfo, _ = s.repo.GetOperatorInfo(*cutting.CutBy)
	}
	
	// 4. Build response
	return &CuttingDetailResponse{
		ID:                  cutting.ID,
		ProductionOrderID:   cutting.ProductionOrderID,
		Status:              cutting.Status,
		StartedAt:           cutting.StartedAt,
		CompletedAt:         cutting.CompletedAt,
		InputLembarBesar:    cutting.InputLembarBesar,
		ExpectedOutput:      cutting.ExpectedOutput,
		OutputSisiranKiri:   cutting.OutputSisiranKiri,
		OutputSisiranKanan:  cutting.OutputSisiranKanan,
		TotalOutput:         cutting.TotalOutput,
		WasteQuantity:       cutting.WasteQuantity,
		WastePercentage:     cutting.WastePercentage,
		WasteReason:         cutting.WasteReason,
		WastePhotoURL:       cutting.WastePhotoURL,
		CuttingMachine:      cutting.CuttingMachine,
		CutBy:               operatorInfo,
		PO:                  poInfo,
	}, nil
}
