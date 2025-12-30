package services

import (
	"encoding/json"
	"sirine-go/backend/models"
	"time"

	"gorm.io/gorm"
)

// CetakService merupakan service untuk Unit Cetak operations
// yang mencakup queue management dan detail retrieval untuk PO yang siap cetak
type CetakService struct {
	db *gorm.DB
}

// NewCetakService membuat instance baru dari CetakService
func NewCetakService(db *gorm.DB) *CetakService {
	return &CetakService{
		db: db,
	}
}

// CetakQueueFilters merupakan struct untuk filter dan pagination queue cetak
type CetakQueueFilters struct {
	Status   string `form:"status"`
	Priority string `form:"priority"`
	Search   string `form:"search"`
	Page     int    `form:"page"`
	PerPage  int    `form:"per_page"`
}

// OBCMasterInfo merupakan struct untuk OBC Master information di response
type OBCMasterInfo struct {
	ID                  uint64 `json:"id"`
	OBCNumber           string `json:"obc_number"`
	Material            string `json:"material"`
	MaterialDescription string `json:"material_description"`
	Seri                string `json:"seri"`
	Warna               string `json:"warna"`
	FactoryCode         string `json:"factory_code"`
	PlatNumber          string `json:"plat_number"`
	Personalization     string `json:"personalization"`
}

// CetakQueueItem merupakan struct untuk response item queue cetak
// yang mencakup informasi PO dan material readiness
type CetakQueueItem struct {
	POID            uint64          `json:"po_id"`
	PONumber        int64           `json:"po_number"`
	OBCNumber       string          `json:"obc_number"`
	ProductName     string          `json:"product_name"`
	Priority        string          `json:"priority"`
	PriorityScore   int             `json:"priority_score"`
	Quantity        int             `json:"quantity"`
	DueDate         time.Time       `json:"due_date"`
	DaysUntilDue    int             `json:"days_until_due"`
	IsPastDue       bool            `json:"is_past_due"`
	MaterialReadyAt time.Time       `json:"material_ready_at"`
	PreparedByID    uint64          `json:"prepared_by_id"`
	PreparedByName  string          `json:"prepared_by_name"`
	MaterialPhotos  []string        `json:"material_photos"`
	Notes           string          `json:"notes"`
	PrepID          uint64          `json:"prep_id"`
	OBCMaster       *OBCMasterInfo  `json:"obc_master"`
}

// CetakQueueResponse merupakan struct untuk paginated queue response
type CetakQueueResponse struct {
	Items      []CetakQueueItem `json:"items"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
}

// CetakDetail merupakan struct untuk detail PO cetak
// yang mencakup full information termasuk material prep data
type CetakDetail struct {
	POID                      uint64          `json:"po_id"`
	PONumber                  int64           `json:"po_number"`
	OBCNumber                 string          `json:"obc_number"`
	SAPCustomerCode           string          `json:"sap_customer_code"`
	SAPProductCode            string          `json:"sap_product_code"`
	ProductName               string          `json:"product_name"`
	ProductSpecifications     interface{}     `json:"product_specifications"`
	QuantityOrdered           int             `json:"quantity_ordered"`
	QuantityTargetLembarBesar int             `json:"quantity_target_lembar_besar"`
	EstimatedRims             int             `json:"estimated_rims"`
	OrderDate                 time.Time       `json:"order_date"`
	DueDate                   time.Time       `json:"due_date"`
	Priority                  string          `json:"priority"`
	PriorityScore             int             `json:"priority_score"`
	DaysUntilDue              int             `json:"days_until_due"`
	IsPastDue                 bool            `json:"is_past_due"`
	CurrentStatus             string          `json:"current_status"`
	Notes                     string          `json:"notes"`
	OBCMaster                 *OBCMasterInfo  `json:"obc_master"`
	MaterialPrep              *MaterialPrep   `json:"material_prep"`
}

// MaterialPrep merupakan struct untuk material preparation detail
type MaterialPrep struct {
	PrepID               uint64    `json:"prep_id"`
	Status               string    `json:"status"`
	StartedAt            *time.Time `json:"started_at"`
	CompletedAt          *time.Time `json:"completed_at"`
	DurationMinutes      int       `json:"duration_minutes"`
	PreparedByID         uint64    `json:"prepared_by_id"`
	PreparedByName       string    `json:"prepared_by_name"`
	SAPPlatCode          string    `json:"sap_plat_code"`
	PlatRetrievedAt      *time.Time `json:"plat_retrieved_at"`
	KertasBlankoQuantity int       `json:"kertas_blanko_quantity"`
	KertasBlankoActual   *int      `json:"kertas_blanko_actual"`
	KertasBlankoVariance int       `json:"kertas_blanko_variance"`
	TintaRequirements    interface{} `json:"tinta_requirements"`
	TintaActual          interface{} `json:"tinta_actual"`
	MaterialPhotos       []string  `json:"material_photos"`
	Notes                string    `json:"notes"`
}

// GetCetakQueue mengambil list PO dengan status READY_FOR_CETAK
// yang sudah selesai material preparation dan siap untuk dicetak
func (s *CetakService) GetCetakQueue(filters CetakQueueFilters) (*CetakQueueResponse, error) {
	// Set default pagination values
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PerPage <= 0 || filters.PerPage > 100 {
		filters.PerPage = 20
	}

	// Query builder dengan filter status READY_FOR_CETAK
	query := s.db.Model(&models.ProductionOrder{}).
		Where("current_status = ?", models.StatusReadyForCetak).
		Preload("OBCMaster").
		Preload("KhazwalMaterialPrep.PreparedByUser")

	// Apply priority filter jika ada
	if filters.Priority != "" {
		query = query.Where("priority = ?", filters.Priority)
	}

	// Apply search filter untuk PO number, OBC number, atau product name
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where(
			"po_number LIKE ? OR obc_number LIKE ? OR product_name LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Count total records untuk pagination
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Execute query dengan pagination dan sorting
	var pos []models.ProductionOrder
	offset := (filters.Page - 1) * filters.PerPage
	if err := query.
		Order("priority_score DESC, due_date ASC").
		Limit(filters.PerPage).
		Offset(offset).
		Find(&pos).Error; err != nil {
		return nil, err
	}

	// Transform ke CetakQueueItem
	items := make([]CetakQueueItem, 0, len(pos))
	for _, po := range pos {
		item := CetakQueueItem{
			POID:          po.ID,
			PONumber:      po.PONumber,
			OBCNumber:     po.OBCNumber,
			ProductName:   po.ProductName,
			Priority:      string(po.Priority),
			PriorityScore: po.PriorityScore,
			Quantity:      po.QuantityOrdered,
			DueDate:       po.DueDate,
			DaysUntilDue:  po.DaysUntilDue(),
			IsPastDue:     po.IsPastDue(),
		}

		// Add OBC Master info jika ada
		if po.OBCMaster != nil {
			item.OBCMaster = &OBCMasterInfo{
				ID:                  po.OBCMaster.ID,
				OBCNumber:           po.OBCMaster.OBCNumber,
				Material:            po.OBCMaster.Material,
				MaterialDescription: po.OBCMaster.MaterialDescription,
				Seri:                po.OBCMaster.Seri,
				Warna:               po.OBCMaster.Warna,
				FactoryCode:         po.OBCMaster.FactoryCode,
				PlatNumber:          po.OBCMaster.PlatNumber,
				Personalization:     po.OBCMaster.Personalization,
			}
		}

		// Add material prep info jika ada
		if po.KhazwalMaterialPrep != nil {
			prep := po.KhazwalMaterialPrep
			if prep.CompletedAt != nil {
				item.MaterialReadyAt = *prep.CompletedAt
			}
			if prep.PreparedBy != nil {
				item.PreparedByID = *prep.PreparedBy
			}
			if prep.PreparedByUser != nil {
				item.PreparedByName = prep.PreparedByUser.FullName
			}
			item.Notes = prep.Notes
			item.PrepID = prep.ID

			// Parse material photos dari JSON
			if prep.MaterialPhotos != nil {
				var photos []string
				if err := json.Unmarshal(prep.MaterialPhotos, &photos); err == nil {
					item.MaterialPhotos = photos
				}
			}
		}

		items = append(items, item)
	}

	// Calculate total pages
	totalPages := int(total) / filters.PerPage
	if int(total)%filters.PerPage > 0 {
		totalPages++
	}

	return &CetakQueueResponse{
		Items:      items,
		Total:      int(total),
		Page:       filters.Page,
		PerPage:    filters.PerPage,
		TotalPages: totalPages,
	}, nil
}

// GetCetakDetail mengambil detail PO untuk cetak termasuk material photos dan prep info
func (s *CetakService) GetCetakDetail(poID uint64) (*CetakDetail, error) {
	var po models.ProductionOrder

	// Query dengan full relations preload
	if err := s.db.
		Preload("OBCMaster").
		Preload("KhazwalMaterialPrep.PreparedByUser").
		First(&po, poID).Error; err != nil {
		return nil, err
	}

	// Validasi: PO harus dalam status READY_FOR_CETAK
	if po.CurrentStatus != models.StatusReadyForCetak {
		return nil, gorm.ErrInvalidData
	}

	detail := &CetakDetail{
		POID:                      po.ID,
		PONumber:                  po.PONumber,
		OBCNumber:                 po.OBCNumber,
		SAPCustomerCode:           po.SAPCustomerCode,
		SAPProductCode:            po.SAPProductCode,
		ProductName:               po.ProductName,
		ProductSpecifications:     po.ProductSpecifications,
		QuantityOrdered:           po.QuantityOrdered,
		QuantityTargetLembarBesar: po.QuantityTargetLembarBesar,
		EstimatedRims:             po.EstimatedRims,
		OrderDate:                 po.OrderDate,
		DueDate:                   po.DueDate,
		Priority:                  string(po.Priority),
		PriorityScore:             po.PriorityScore,
		DaysUntilDue:              po.DaysUntilDue(),
		IsPastDue:                 po.IsPastDue(),
		CurrentStatus:             string(po.CurrentStatus),
		Notes:                     po.Notes,
	}

	// Add OBC Master info jika ada
	if po.OBCMaster != nil {
		detail.OBCMaster = &OBCMasterInfo{
			ID:                  po.OBCMaster.ID,
			OBCNumber:           po.OBCMaster.OBCNumber,
			Material:            po.OBCMaster.Material,
			MaterialDescription: po.OBCMaster.MaterialDescription,
			Seri:                po.OBCMaster.Seri,
			Warna:               po.OBCMaster.Warna,
			FactoryCode:         po.OBCMaster.FactoryCode,
			PlatNumber:          po.OBCMaster.PlatNumber,
			Personalization:     po.OBCMaster.Personalization,
		}
	}

	// Add material prep detail jika ada
	if po.KhazwalMaterialPrep != nil {
		prep := po.KhazwalMaterialPrep
		materialPrep := &MaterialPrep{
			PrepID:               prep.ID,
			Status:               string(prep.Status),
			StartedAt:            prep.StartedAt,
			CompletedAt:          prep.CompletedAt,
			SAPPlatCode:          prep.SAPPlatCode,
			PlatRetrievedAt:      prep.PlatRetrievedAt,
			KertasBlankoQuantity: prep.KertasBlankoQuantity,
			KertasBlankoActual:   prep.KertasBlankoActual,
			Notes:                prep.Notes,
		}

		// Handle pointer fields
		if prep.DurationMinutes != nil {
			materialPrep.DurationMinutes = *prep.DurationMinutes
		}
		if prep.KertasBlankoVariance != nil {
			materialPrep.KertasBlankoVariance = *prep.KertasBlankoVariance
		}

		if prep.PreparedBy != nil {
			materialPrep.PreparedByID = *prep.PreparedBy
		}
		if prep.PreparedByUser != nil {
			materialPrep.PreparedByName = prep.PreparedByUser.FullName
		}

		// Parse tinta requirements
		if prep.TintaRequirements != nil {
			var tintaReq interface{}
			if err := json.Unmarshal(prep.TintaRequirements, &tintaReq); err == nil {
				materialPrep.TintaRequirements = tintaReq
			}
		}

		// Parse tinta actual
		if prep.TintaActual != nil {
			var tintaAct interface{}
			if err := json.Unmarshal(prep.TintaActual, &tintaAct); err == nil {
				materialPrep.TintaActual = tintaAct
			}
		}

		// Parse material photos
		if prep.MaterialPhotos != nil {
			var photos []string
			if err := json.Unmarshal(prep.MaterialPhotos, &photos); err == nil {
				materialPrep.MaterialPhotos = photos
			}
		}

		detail.MaterialPrep = materialPrep
	}

	return detail, nil
}
