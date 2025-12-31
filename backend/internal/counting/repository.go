package counting

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CountingRepository merupakan interface untuk database operations counting
type CountingRepository interface {
	// Queue operations
	GetCountingQueue(machineID *uint64, dateFrom, dateTo *time.Time) ([]QueueItemResponse, error)
	
	// CRUD operations
	GetByID(id uint64) (*KhazwalCountingResult, error)
	GetByPOID(poID uint64) (*KhazwalCountingResult, error)
	GetCountingDetailWithRelations(id uint64) (*CountingDetailResponse, error)
	Create(counting *KhazwalCountingResult) error
	Update(counting *KhazwalCountingResult) error
	
	// Check operations
	ExistsInProgressByPOID(poID uint64) (bool, error)
}

// countingRepositoryImpl merupakan implementasi CountingRepository
type countingRepositoryImpl struct {
	db *gorm.DB
}

// NewCountingRepository membuat instance baru CountingRepository
func NewCountingRepository(db *gorm.DB) CountingRepository {
	return &countingRepositoryImpl{db: db}
}

// GetCountingQueue mengambil list PO yang menunggu penghitungan (FIFO)
// dengan join ke print_job_summaries untuk mendapatkan info mesin dan operator
func (r *countingRepositoryImpl) GetCountingQueue(machineID *uint64, dateFrom, dateTo *time.Time) ([]QueueItemResponse, error) {
	var results []QueueItemResponse

	query := r.db.Table("production_orders po").
		Select(`
			po.id as po_id,
			po.po_number,
			po.obc_number,
			po.quantity_target_lembar_besar as target_quantity,
			pjs.finalized_at as print_completed_at,
			EXTRACT(EPOCH FROM (NOW() - pjs.finalized_at)) / 60 as waiting_minutes,
			CASE WHEN EXTRACT(EPOCH FROM (NOW() - pjs.finalized_at)) / 60 > 120 THEN true ELSE false END as is_overdue,
			m.id as machine_id,
			m.name as machine_name,
			m.code as machine_code,
			u.id as operator_id,
			u.name as operator_name,
			u.nip as operator_nip
		`).
		Joins("INNER JOIN print_job_summaries pjs ON pjs.production_order_id = po.id").
		Joins("LEFT JOIN machines m ON m.id = pjs.machine_id").
		Joins("LEFT JOIN users u ON u.id = pjs.operator_id").
		Where("po.current_status = ?", "WAITING_COUNTING").
		Where("pjs.finalized_at IS NOT NULL").
		Order("pjs.finalized_at ASC") // FIFO

	// Apply optional filters
	if machineID != nil {
		query = query.Where("pjs.machine_id = ?", *machineID)
	}
	if dateFrom != nil {
		query = query.Where("DATE(pjs.finalized_at) >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("DATE(pjs.finalized_at) <= ?", *dateTo)
	}

	// Execute query dengan raw scan ke struct
	rows, err := query.Rows()
	if err != nil {
		return nil, fmt.Errorf("gagal query counting queue: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item QueueItemResponse
		var machineID, operatorID *uint64
		var machineName, machineCode, operatorName, operatorNIP *string
		var waitingMinutes float64

		err := rows.Scan(
			&item.POID,
			&item.PONumber,
			&item.OBCNumber,
			&item.TargetQuantity,
			&item.PrintCompletedAt,
			&waitingMinutes,
			&item.IsOverdue,
			&machineID,
			&machineName,
			&machineCode,
			&operatorID,
			&operatorName,
			&operatorNIP,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan queue item: %w", err)
		}

		item.WaitingMinutes = int(waitingMinutes)

		// Populate machine info
		if machineID != nil && machineName != nil && machineCode != nil {
			item.Machine = &MachineInfo{
				ID:   *machineID,
				Name: *machineName,
				Code: *machineCode,
			}
		}

		// Populate operator info
		if operatorID != nil && operatorName != nil && operatorNIP != nil {
			item.Operator = &OperatorInfo{
				ID:   *operatorID,
				Name: *operatorName,
				NIP:  *operatorNIP,
			}
		}

		results = append(results, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating queue results: %w", err)
	}

	return results, nil
}

// GetByID mengambil counting record berdasarkan ID
func (r *countingRepositoryImpl) GetByID(id uint64) (*KhazwalCountingResult, error) {
	var counting KhazwalCountingResult
	
	if err := r.db.First(&counting, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("counting record tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil counting record: %w", err)
	}
	
	return &counting, nil
}

// GetByPOID mengambil counting record berdasarkan Production Order ID
func (r *countingRepositoryImpl) GetByPOID(poID uint64) (*KhazwalCountingResult, error) {
	var counting KhazwalCountingResult
	
	if err := r.db.Where("production_order_id = ?", poID).First(&counting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("counting record untuk PO tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil counting record: %w", err)
	}
	
	return &counting, nil
}

// Create membuat counting record baru
func (r *countingRepositoryImpl) Create(counting *KhazwalCountingResult) error {
	if err := r.db.Create(counting).Error; err != nil {
		return fmt.Errorf("gagal membuat counting record: %w", err)
	}
	return nil
}

// Update mengupdate counting record
func (r *countingRepositoryImpl) Update(counting *KhazwalCountingResult) error {
	if err := r.db.Save(counting).Error; err != nil {
		return fmt.Errorf("gagal update counting record: %w", err)
	}
	return nil
}

// ExistsInProgressByPOID memeriksa apakah ada counting IN_PROGRESS untuk PO tertentu
func (r *countingRepositoryImpl) ExistsInProgressByPOID(poID uint64) (bool, error) {
	var count int64
	
	if err := r.db.Model(&KhazwalCountingResult{}).
		Where("production_order_id = ?", poID).
		Where("status = ?", CountingInProgress).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("gagal check existing counting: %w", err)
	}
	
	return count > 0, nil
}

// GetCountingDetailWithRelations mengambil counting detail dengan relasi PO dan print info
// untuk digunakan di GET /counting/:id endpoint
func (r *countingRepositoryImpl) GetCountingDetailWithRelations(id uint64) (*CountingDetailResponse, error) {
	var response CountingDetailResponse

	// Query dengan joins
	query := r.db.Table("khazwal_counting_results kcr").
		Select(`
			kcr.id,
			kcr.production_order_id,
			kcr.status,
			kcr.started_at,
			kcr.quantity_good,
			kcr.quantity_defect,
			kcr.total_counted,
			kcr.variance_from_target,
			kcr.percentage_good,
			kcr.percentage_defect,
			kcr.defect_breakdown,
			kcr.variance_reason,
			po.po_number,
			po.obc_number,
			po.quantity_target_lembar_besar as target_quantity,
			u.id as counted_by_id,
			u.name as counted_by_name,
			u.nip as counted_by_nip,
			pjs.machine_id,
			m.name as machine_name,
			pjs.operator_id,
			op.name as operator_name,
			pjs.finalized_at
		`).
		Joins("INNER JOIN production_orders po ON po.id = kcr.production_order_id").
		Joins("LEFT JOIN users u ON u.id = kcr.counted_by").
		Joins("LEFT JOIN print_job_summaries pjs ON pjs.production_order_id = kcr.production_order_id").
		Joins("LEFT JOIN machines m ON m.id = pjs.machine_id").
		Joins("LEFT JOIN users op ON op.id = pjs.operator_id").
		Where("kcr.id = ?", id)

	// Scan hasil ke temporary struct
	var result struct {
		ID                 uint64
		ProductionOrderID  uint64
		Status             string
		StartedAt          *time.Time
		QuantityGood       int
		QuantityDefect     int
		TotalCounted       int
		VarianceFromTarget *int
		PercentageGood     *float64
		PercentageDefect   *float64
		DefectBreakdown    []byte
		VarianceReason     string
		PONumber           int64
		OBCNumber          string
		TargetQuantity     int
		CountedByID        *uint64
		CountedByName      *string
		CountedByNIP       *string
		MachineID          *uint64
		MachineName        *string
		OperatorID         *uint64
		OperatorName       *string
		FinalizedAt        *time.Time
	}

	if err := query.Scan(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("counting detail tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil counting detail: %w", err)
	}

	// Populate response
	response.ID = result.ID
	response.ProductionOrderID = result.ProductionOrderID
	response.Status = CountingStatus(result.Status)
	response.StartedAt = result.StartedAt
	response.QuantityGood = result.QuantityGood
	response.QuantityDefect = result.QuantityDefect
	response.TotalCounted = result.TotalCounted
	response.VarianceFromTarget = result.VarianceFromTarget
	response.PercentageGood = result.PercentageGood
	response.PercentageDefect = result.PercentageDefect
	response.VarianceReason = result.VarianceReason

	// Parse defect breakdown
	if result.DefectBreakdown != nil {
		breakdown, err := ParseDefectBreakdown(result.DefectBreakdown)
		if err == nil {
			response.DefectBreakdown = breakdown
		}
	}

	// Populate PO info
	response.PO = &POInfo{
		PONumber:       result.PONumber,
		OBCNumber:      result.OBCNumber,
		TargetQuantity: result.TargetQuantity,
	}

	// Populate counted by info
	if result.CountedByID != nil && result.CountedByName != nil && result.CountedByNIP != nil {
		response.CountedBy = &OperatorInfo{
			ID:   *result.CountedByID,
			Name: *result.CountedByName,
			NIP:  *result.CountedByNIP,
		}
	}

	// Populate print info
	if result.MachineName != nil && result.OperatorName != nil && result.FinalizedAt != nil {
		response.PrintInfo = &PrintInfo{
			MachineName:  *result.MachineName,
			OperatorName: *result.OperatorName,
			FinalizedAt:  *result.FinalizedAt,
		}
	}

	return &response, nil
}
