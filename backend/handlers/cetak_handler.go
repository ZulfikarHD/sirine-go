package handlers

import (
	"net/http"
	"sirine-go/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CetakHandler merupakan handler untuk Unit Cetak endpoints
// yang mencakup queue retrieval dan detail view untuk PO siap cetak
type CetakHandler struct {
	cetakService *services.CetakService
}

// NewCetakHandler membuat instance baru dari CetakHandler
func NewCetakHandler(cetakService *services.CetakService) *CetakHandler {
	return &CetakHandler{
		cetakService: cetakService,
	}
}

// GetQueue mengambil queue list PO yang siap untuk dicetak
// dengan filter, search, dan pagination support
// @route GET /api/cetak/queue
// @access OPERATOR_CETAK, SUPERVISOR_CETAK, ADMIN, MANAGER
func (h *CetakHandler) GetQueue(c *gin.Context) {
	// Parse query filters dari request
	var filters services.CetakQueueFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter query tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Call service untuk get queue
	queueResponse, err := h.cetakService.GetCetakQueue(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil queue cetak",
			"error":   err.Error(),
		})
		return
	}

	// Transform response untuk API dengan format yang konsisten
	items := make([]CetakQueueItemDTO, 0, len(queueResponse.Items))
	for _, item := range queueResponse.Items {
		dto := CetakQueueItemDTO{
			POID:           item.POID,
			PONumber:       item.PONumber,
			OBCNumber:      item.OBCNumber,
			ProductName:    item.ProductName,
			Priority:       item.Priority,
			PriorityScore:  item.PriorityScore,
			Quantity:       item.Quantity,
			DueDate:        item.DueDate.Format("2006-01-02"),
			DaysUntilDue:   item.DaysUntilDue,
			IsPastDue:      item.IsPastDue,
			PreparedByID:   item.PreparedByID,
			PreparedByName: item.PreparedByName,
			MaterialPhotos: item.MaterialPhotos,
			Notes:          item.Notes,
			PrepID:         item.PrepID,
		}

		// Format material ready at jika ada
		if !item.MaterialReadyAt.IsZero() {
			dto.MaterialReadyAt = item.MaterialReadyAt.Format("2006-01-02 15:04:05")
		}

		items = append(items, dto)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Queue cetak berhasil diambil",
		"data": gin.H{
			"items":       items,
			"total":       queueResponse.Total,
			"page":        queueResponse.Page,
			"per_page":    queueResponse.PerPage,
			"total_pages": queueResponse.TotalPages,
		},
	})
}

// GetDetail mengambil detail PO untuk cetak termasuk material photos
// dengan full relations untuk complete information display
// @route GET /api/cetak/queue/:id
// @access OPERATOR_CETAK, SUPERVISOR_CETAK, ADMIN, MANAGER
func (h *CetakHandler) GetDetail(c *gin.Context) {
	// Parse ID dari URL param
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Call service untuk get detail
	detail, err := h.cetakService.GetCetakDetail(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Production Order tidak ditemukan",
			})
			return
		}

		if err == gorm.ErrInvalidData {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "PO tidak dalam status siap cetak",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil detail PO",
			"error":   err.Error(),
		})
		return
	}

	// Transform ke DTO
	dto := CetakDetailDTO{
		POID:                      detail.POID,
		PONumber:                  detail.PONumber,
		OBCNumber:                 detail.OBCNumber,
		SAPCustomerCode:           detail.SAPCustomerCode,
		SAPProductCode:            detail.SAPProductCode,
		ProductName:               detail.ProductName,
		ProductSpecifications:     detail.ProductSpecifications,
		QuantityOrdered:           detail.QuantityOrdered,
		QuantityTargetLembarBesar: detail.QuantityTargetLembarBesar,
		EstimatedRims:             detail.EstimatedRims,
		OrderDate:                 detail.OrderDate.Format("2006-01-02"),
		DueDate:                   detail.DueDate.Format("2006-01-02"),
		Priority:                  detail.Priority,
		PriorityScore:             detail.PriorityScore,
		DaysUntilDue:              detail.DaysUntilDue,
		IsPastDue:                 detail.IsPastDue,
		CurrentStatus:             detail.CurrentStatus,
		Notes:                     detail.Notes,
	}

	// Transform material prep jika ada
	if detail.MaterialPrep != nil {
		prep := detail.MaterialPrep
		materialPrepDTO := &MaterialPrepDTO{
			PrepID:               prep.PrepID,
			Status:               prep.Status,
			DurationMinutes:      prep.DurationMinutes,
			PreparedByID:         prep.PreparedByID,
			PreparedByName:       prep.PreparedByName,
			SAPPlatCode:          prep.SAPPlatCode,
			KertasBlankoQuantity: prep.KertasBlankoQuantity,
			KertasBlankoActual:   prep.KertasBlankoActual,
			KertasBlankoVariance: prep.KertasBlankoVariance,
			TintaRequirements:    prep.TintaRequirements,
			TintaActual:          prep.TintaActual,
			MaterialPhotos:       prep.MaterialPhotos,
			Notes:                prep.Notes,
		}

		if prep.StartedAt != nil {
			materialPrepDTO.StartedAt = prep.StartedAt.Format("2006-01-02 15:04:05")
		}
		if prep.CompletedAt != nil {
			materialPrepDTO.CompletedAt = prep.CompletedAt.Format("2006-01-02 15:04:05")
		}
		if prep.PlatRetrievedAt != nil {
			materialPrepDTO.PlatRetrievedAt = prep.PlatRetrievedAt.Format("2006-01-02 15:04:05")
		}

		dto.MaterialPrep = materialPrepDTO
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail PO berhasil diambil",
		"data":    dto,
	})
}

// DTOs untuk API responses

// CetakQueueItemDTO merupakan DTO untuk queue list item response
type CetakQueueItemDTO struct {
	POID            uint64   `json:"po_id"`
	PONumber        int64    `json:"po_number"`
	OBCNumber       string   `json:"obc_number"`
	ProductName     string   `json:"product_name"`
	Priority        string   `json:"priority"`
	PriorityScore   int      `json:"priority_score"`
	Quantity        int      `json:"quantity"`
	DueDate         string   `json:"due_date"`
	DaysUntilDue    int      `json:"days_until_due"`
	IsPastDue       bool     `json:"is_past_due"`
	MaterialReadyAt string   `json:"material_ready_at,omitempty"`
	PreparedByID    uint64   `json:"prepared_by_id,omitempty"`
	PreparedByName  string   `json:"prepared_by_name,omitempty"`
	MaterialPhotos  []string `json:"material_photos,omitempty"`
	Notes           string   `json:"notes,omitempty"`
	PrepID          uint64   `json:"prep_id,omitempty"`
}

// CetakDetailDTO merupakan DTO untuk detail PO response
type CetakDetailDTO struct {
	POID                      uint64           `json:"po_id"`
	PONumber                  int64            `json:"po_number"`
	OBCNumber                 string           `json:"obc_number"`
	SAPCustomerCode           string           `json:"sap_customer_code"`
	SAPProductCode            string           `json:"sap_product_code"`
	ProductName               string           `json:"product_name"`
	ProductSpecifications     interface{}      `json:"product_specifications"`
	QuantityOrdered           int              `json:"quantity_ordered"`
	QuantityTargetLembarBesar int              `json:"quantity_target_lembar_besar"`
	EstimatedRims             int              `json:"estimated_rims"`
	OrderDate                 string           `json:"order_date"`
	DueDate                   string           `json:"due_date"`
	Priority                  string           `json:"priority"`
	PriorityScore             int              `json:"priority_score"`
	DaysUntilDue              int              `json:"days_until_due"`
	IsPastDue                 bool             `json:"is_past_due"`
	CurrentStatus             string           `json:"current_status"`
	Notes                     string           `json:"notes"`
	MaterialPrep              *MaterialPrepDTO `json:"material_prep,omitempty"`
}

// MaterialPrepDTO merupakan DTO untuk material preparation detail
type MaterialPrepDTO struct {
	PrepID               uint64      `json:"prep_id"`
	Status               string      `json:"status"`
	StartedAt            string      `json:"started_at,omitempty"`
	CompletedAt          string      `json:"completed_at,omitempty"`
	DurationMinutes      int         `json:"duration_minutes"`
	PreparedByID         uint64      `json:"prepared_by_id"`
	PreparedByName       string      `json:"prepared_by_name"`
	SAPPlatCode          string      `json:"sap_plat_code"`
	PlatRetrievedAt      string      `json:"plat_retrieved_at,omitempty"`
	KertasBlankoQuantity int         `json:"kertas_blanko_quantity"`
	KertasBlankoActual   *int        `json:"kertas_blanko_actual,omitempty"`
	KertasBlankoVariance int         `json:"kertas_blanko_variance"`
	TintaRequirements    interface{} `json:"tinta_requirements,omitempty"`
	TintaActual          interface{} `json:"tinta_actual,omitempty"`
	MaterialPhotos       []string    `json:"material_photos,omitempty"`
	Notes                string      `json:"notes,omitempty"`
}
