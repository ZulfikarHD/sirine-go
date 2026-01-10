package cutting

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Repository merupakan interface untuk database operations cutting
type Repository interface {
	// Queue operations
	GetCuttingQueue(filters QueueFilters) ([]QueueItemResponse, error)
	GetQueueMetadata(filters QueueFilters) (QueueMetadata, error)
	
	// CRUD operations
	Create(cutting *KhazwalCuttingResult) error
	GetByID(id uint64) (*KhazwalCuttingResult, error)
	GetByPOID(poID uint64) (*KhazwalCuttingResult, error)
	Update(cutting *KhazwalCuttingResult) error
	
	// Related data operations
	GetCountingResultByPOID(poID uint64) (*CountingResult, error)
	UpdatePOStatus(poID uint64, stage, status string) error
	UpdatePOStageTracking(poID uint64, field string, value time.Time) error
	GetPOInfo(poID uint64) (*POInfo, error)
	GetOperatorInfo(userID uint64) (*OperatorInfo, error)
}

// repository merupakan implementasi konkret dari Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository membuat instance baru dari repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CountingResult merupakan struct simplified untuk query counting result
type CountingResult struct {
	ID                uint64    `gorm:"column:id"`
	ProductionOrderID uint64    `gorm:"column:production_order_id"`
	QuantityGood      int       `gorm:"column:quantity_good"`
	Status            string    `gorm:"column:status"`
	CompletedAt       time.Time `gorm:"column:completed_at"`
}

// GetCuttingQueue mengambil list PO yang siap untuk dipotong
func (r *repository) GetCuttingQueue(filters QueueFilters) ([]QueueItemResponse, error) {
	var results []QueueItemResponse
	
	query := r.db.Table("production_orders as po").
		Select(`
			po.id as po_id,
			po.po_number,
			po.obc_number,
			po.priority,
			kcr.quantity_good as input_lembar_besar,
			kcr.quantity_good * 2 as estimated_output,
			kcr.completed_at as counting_completed_at,
			TIMESTAMPDIFF(MINUTE, kcr.completed_at, NOW()) as waiting_minutes,
			CASE WHEN TIMESTAMPDIFF(MINUTE, kcr.completed_at, NOW()) > 60 THEN 1 ELSE 0 END as is_overdue
		`).
		Joins("INNER JOIN khazwal_counting_results kcr ON kcr.production_order_id = po.id").
		Where("po.current_stage = ?", "KHAZWAL_CUTTING").
		Where("po.current_status = ?", "SIAP_POTONG").
		Where("po.deleted_at IS NULL").
		Where("kcr.status = ?", "COMPLETED")
	
	// Apply filters
	if filters.Priority != "" {
		query = query.Where("po.priority = ?", filters.Priority)
	}
	
	if filters.DateFrom != "" {
		query = query.Where("kcr.completed_at >= ?", filters.DateFrom)
	}
	
	if filters.DateTo != "" {
		query = query.Where("kcr.completed_at <= ?", filters.DateTo)
	}
	
	// Default sorting: priority (URGENT first) + FIFO (oldest completed_at first)
	sortBy := filters.SortBy
	sortOrder := filters.SortOrder
	
	if sortOrder == "" {
		sortOrder = "ASC"
	}
	
	if sortBy == "date" {
		query = query.Order(fmt.Sprintf("kcr.completed_at %s", sortOrder))
	} else {
		// Default: priority + FIFO
		query = query.Order(`
			CASE po.priority 
				WHEN 'URGENT' THEN 1 
				WHEN 'HIGH' THEN 2
				WHEN 'NORMAL' THEN 3 
				WHEN 'LOW' THEN 4
				ELSE 5 
			END ASC,
			kcr.completed_at ASC
		`)
	}
	
	err := query.Find(&results).Error
	return results, err
}

// GetQueueMetadata mengambil metadata queue (total, counts by priority)
func (r *repository) GetQueueMetadata(filters QueueFilters) (QueueMetadata, error) {
	var meta QueueMetadata
	
	query := r.db.Table("production_orders as po").
		Joins("INNER JOIN khazwal_counting_results kcr ON kcr.production_order_id = po.id").
		Where("po.current_stage = ?", "KHAZWAL_CUTTING").
		Where("po.current_status = ?", "SIAP_POTONG").
		Where("po.deleted_at IS NULL").
		Where("kcr.status = ?", "COMPLETED")
	
	// Apply same filters as queue
	if filters.Priority != "" {
		query = query.Where("po.priority = ?", filters.Priority)
	}
	
	if filters.DateFrom != "" {
		query = query.Where("kcr.completed_at >= ?", filters.DateFrom)
	}
	
	if filters.DateTo != "" {
		query = query.Where("kcr.completed_at <= ?", filters.DateTo)
	}
	
	// Get total count
	var total int64
	query.Count(&total)
	meta.Total = int(total)
	
	// Get count by priority (without filters)
	var urgentCount int64
	r.db.Table("production_orders as po").
		Joins("INNER JOIN khazwal_counting_results kcr ON kcr.production_order_id = po.id").
		Where("po.current_stage = ?", "KHAZWAL_CUTTING").
		Where("po.current_status = ?", "SIAP_POTONG").
		Where("po.deleted_at IS NULL").
		Where("kcr.status = ?", "COMPLETED").
		Where("po.priority = ?", "URGENT").
		Count(&urgentCount)
	meta.UrgentCount = int(urgentCount)
	
	var normalCount int64
	r.db.Table("production_orders as po").
		Joins("INNER JOIN khazwal_counting_results kcr ON kcr.production_order_id = po.id").
		Where("po.current_stage = ?", "KHAZWAL_CUTTING").
		Where("po.current_status = ?", "SIAP_POTONG").
		Where("po.deleted_at IS NULL").
		Where("kcr.status = ?", "COMPLETED").
		Where("po.priority = ?", "NORMAL").
		Count(&normalCount)
	meta.NormalCount = int(normalCount)
	
	return meta, nil
}

// Create membuat cutting result baru
func (r *repository) Create(cutting *KhazwalCuttingResult) error {
	return r.db.Create(cutting).Error
}

// GetByID mengambil cutting result berdasarkan ID
func (r *repository) GetByID(id uint64) (*KhazwalCuttingResult, error) {
	var cutting KhazwalCuttingResult
	err := r.db.Where("id = ?", id).First(&cutting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCuttingNotFound
		}
		return nil, err
	}
	return &cutting, nil
}

// GetByPOID mengambil cutting result berdasarkan production order ID
func (r *repository) GetByPOID(poID uint64) (*KhazwalCuttingResult, error) {
	var cutting KhazwalCuttingResult
	err := r.db.Where("production_order_id = ?", poID).First(&cutting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCuttingNotFound
		}
		return nil, err
	}
	return &cutting, nil
}

// Update mengupdate cutting result
func (r *repository) Update(cutting *KhazwalCuttingResult) error {
	return r.db.Save(cutting).Error
}

// GetCountingResultByPOID mengambil counting result berdasarkan PO ID
func (r *repository) GetCountingResultByPOID(poID uint64) (*CountingResult, error) {
	var result CountingResult
	err := r.db.Table("khazwal_counting_results").
		Where("production_order_id = ?", poID).
		First(&result).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCountingNotCompleted
		}
		return nil, err
	}
	
	return &result, nil
}

// UpdatePOStatus mengupdate status production order
func (r *repository) UpdatePOStatus(poID uint64, stage, status string) error {
	return r.db.Table("production_orders").
		Where("id = ?", poID).
		Updates(map[string]interface{}{
			"current_stage":  stage,
			"current_status": status,
			"updated_at":     time.Now(),
		}).Error
}

// UpdatePOStageTracking mengupdate timestamp di po_stage_tracking
func (r *repository) UpdatePOStageTracking(poID uint64, field string, value time.Time) error {
	// Check if record exists
	var count int64
	r.db.Table("po_stage_tracking").
		Where("production_order_id = ?", poID).
		Where("stage = ?", "KHAZWAL_CUTTING").
		Count(&count)
	
	if count == 0 {
		// Create new record
		return r.db.Table("po_stage_tracking").Create(map[string]interface{}{
			"production_order_id": poID,
			"stage":               "KHAZWAL_CUTTING",
			field:                 value,
			"created_at":          time.Now(),
			"updated_at":          time.Now(),
		}).Error
	}
	
	// Update existing record
	return r.db.Table("po_stage_tracking").
		Where("production_order_id = ?", poID).
		Where("stage = ?", "KHAZWAL_CUTTING").
		Updates(map[string]interface{}{
			field:        value,
			"updated_at": time.Now(),
		}).Error
}

// GetPOInfo mengambil info production order
func (r *repository) GetPOInfo(poID uint64) (*POInfo, error) {
	var info POInfo
	err := r.db.Table("production_orders").
		Select("po_number, obc_number, priority, target_quantity").
		Where("id = ?", poID).
		First(&info).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPONotFound
		}
		return nil, err
	}
	
	return &info, nil
}

// GetOperatorInfo mengambil info operator/staff
func (r *repository) GetOperatorInfo(userID uint64) (*OperatorInfo, error) {
	var info OperatorInfo
	err := r.db.Table("users").
		Select("id, full_name as name, nip").
		Where("id = ?", userID).
		First(&info).Error
	
	if err != nil {
		return nil, err
	}
	
	return &info, nil
}
