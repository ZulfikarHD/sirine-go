package handlers

import (
	"net/http"
	"sirine-go/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// KhazwalHandler merupakan handler untuk Khazwal Material Preparation endpoints
type KhazwalHandler struct {
	khazwalService *services.KhazwalService
}

// NewKhazwalHandler membuat instance baru dari KhazwalHandler
func NewKhazwalHandler(khazwalService *services.KhazwalService) *KhazwalHandler {
	return &KhazwalHandler{
		khazwalService: khazwalService,
	}
}

// GetQueue mengambil queue list PO untuk material preparation
// dengan filter, search, dan pagination support
// @route GET /api/khazwal/material-prep/queue
// @access STAFF_KHAZWAL, ADMIN, MANAGER
func (h *KhazwalHandler) GetQueue(c *gin.Context) {
	// Parse query filters dari request
	var filters services.QueueFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter query tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Call service untuk get queue
	queueResponse, err := h.khazwalService.GetMaterialPrepQueue(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil queue material preparation",
			"error":   err.Error(),
		})
		return
	}

	// Transform response ke DTO format untuk clean API response
	queueItems := make([]QueueItemDTO, 0, len(queueResponse.Items))
	for _, po := range queueResponse.Items {
		// Get prep status dari material prep relation
		prepStatus := "PENDING"
		if po.KhazwalMaterialPrep != nil {
			prepStatus = string(po.KhazwalMaterialPrep.Status)
		}

		queueItems = append(queueItems, QueueItemDTO{
			ID:            po.ID,
			PONumber:      po.PONumber,
			OBCNumber:     po.OBCNumber,
			ProductName:   po.ProductName,
			Priority:      string(po.Priority),
			PriorityScore: po.PriorityScore,
			DueDate:       po.DueDate.Format("2006-01-02"),
			DaysUntilDue:  po.DaysUntilDue(),
			IsPastDue:     po.IsPastDue(),
			CurrentStatus: string(po.CurrentStatus),
			Quantity:      po.QuantityOrdered,
			PrepStatus:    prepStatus,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Queue berhasil diambil",
		"data": gin.H{
			"items":       queueItems,
			"total":       queueResponse.Total,
			"page":        queueResponse.Page,
			"per_page":    queueResponse.PerPage,
			"total_pages": queueResponse.TotalPages,
		},
	})
}

// GetDetail mengambil detail PO dan material prep info
// dengan full relations untuk complete information display
// @route GET /api/khazwal/material-prep/:id
// @access STAFF_KHAZWAL, ADMIN, MANAGER
func (h *KhazwalHandler) GetDetail(c *gin.Context) {
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
	po, err := h.khazwalService.GetMaterialPrepDetail(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Production Order tidak ditemukan",
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

	// Transform ke DTO dengan computed fields
	detailDTO := DetailDTO{
		ID:                        po.ID,
		PONumber:                  po.PONumber,
		OBCNumber:                 po.OBCNumber,
		SAPCustomerCode:           po.SAPCustomerCode,
		SAPProductCode:            po.SAPProductCode,
		ProductName:               po.ProductName,
		ProductSpecifications:     po.ProductSpecifications,
		QuantityOrdered:           po.QuantityOrdered,
		QuantityTargetLembarBesar: po.QuantityTargetLembarBesar,
		EstimatedRims:             po.EstimatedRims,
		OrderDate:                 po.OrderDate.Format("2006-01-02"),
		DueDate:                   po.DueDate.Format("2006-01-02"),
		Priority:                  string(po.Priority),
		PriorityScore:             po.PriorityScore,
		DaysUntilDue:              po.DaysUntilDue(),
		IsPastDue:                 po.IsPastDue(),
		CurrentStage:              string(po.CurrentStage),
		CurrentStatus:             string(po.CurrentStatus),
		Notes:                     po.Notes,
		KhazwalMaterialPrep:       po.KhazwalMaterialPrep,
		StageTracking:             po.StageTracking,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail PO berhasil diambil",
		"data":    detailDTO,
	})
}

// StartPrep memulai proses material preparation
// dengan validation dan transaction untuk ensure data integrity
// @route POST /api/khazwal/material-prep/:id/start
// @access STAFF_KHAZWAL, ADMIN, MANAGER
func (h *KhazwalHandler) StartPrep(c *gin.Context) {
	// Parse ID dari URL param
	idStr := c.Param("id")
	poID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID dari context (dari auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Format user ID tidak valid",
		})
		return
	}

	// Call service untuk start material prep
	po, err := h.khazwalService.StartMaterialPrep(poID, userID)
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
				"message": "PO tidak dalam status yang valid untuk dimulai",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memulai material preparation",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Material preparation berhasil dimulai",
		"data":    po,
	})
}

// ConfirmPlat mengkonfirmasi pengambilan plat dengan barcode validation
// @route POST /api/khazwal/material-prep/:id/confirm-plat
// @access STAFF_KHAZWAL, ADMIN, MANAGER
func (h *KhazwalHandler) ConfirmPlat(c *gin.Context) {
	// Parse ID dari URL param
	idStr := c.Param("id")
	prepID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Parse request body
	var req ConfirmPlatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Request body tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID dari context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Format user ID tidak valid",
		})
		return
	}

	// Call service untuk confirm plat
	err = h.khazwalService.ConfirmPlatRetrieval(prepID, req.PlatCode, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Material preparation tidak ditemukan",
			})
			return
		}

		// Check for PlatNumberMissingError (data issue)
		if platErr, ok := err.(*services.PlatNumberMissingError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": platErr.Error(),
			})
			return
		}

		if err == gorm.ErrInvalidData {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Kode plat tidak sesuai atau status tidak valid",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengkonfirmasi plat",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plat berhasil dikonfirmasi",
	})
}

// UpdateKertas mengupdate kertas blanko actual dengan variance calculation
// @route PATCH /api/khazwal/material-prep/:id/kertas
// @access STAFF_KHAZWAL, ADMIN, MANAGER
func (h *KhazwalHandler) UpdateKertas(c *gin.Context) {
	// Parse ID dari URL param
	idStr := c.Param("id")
	prepID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Parse request body
	var req UpdateKertasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Request body tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID dari context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Format user ID tidak valid",
		})
		return
	}

	// Call service untuk update kertas
	err = h.khazwalService.UpdateKertasBlanko(prepID, req.ActualQty, req.VarianceReason, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Material preparation tidak ditemukan",
			})
			return
		}

		if err == gorm.ErrInvalidData {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Variance > 5% memerlukan alasan atau status tidak valid",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate kertas blanko",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kertas blanko berhasil diupdate",
	})
}

// UpdateTinta mengupdate tinta actual dengan checklist validation
// @route PATCH /api/khazwal/material-prep/:id/tinta
// @access STAFF_KHAZWAL, ADMIN, MANAGER
func (h *KhazwalHandler) UpdateTinta(c *gin.Context) {
	// Parse ID dari URL param
	idStr := c.Param("id")
	prepID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Parse request body
	var req UpdateTintaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Request body tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID dari context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Format user ID tidak valid",
		})
		return
	}

	// Call service untuk update tinta
	err = h.khazwalService.UpdateTinta(prepID, req.TintaActual, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Material preparation tidak ditemukan",
			})
			return
		}

		if err == gorm.ErrInvalidData {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Status tidak valid untuk update tinta",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate tinta",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tinta berhasil diupdate",
	})
}

// Finalize menyelesaikan material preparation dengan upload photos dan notes
// @route POST /api/khazwal/material-prep/:id/finalize
// @access STAFF_KHAZWAL, ADMIN, MANAGER
func (h *KhazwalHandler) Finalize(c *gin.Context) {
	// Parse ID dari URL param
	idStr := c.Param("id")
	prepID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Parse request body
	var req FinalizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Request body tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID dari context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Format user ID tidak valid",
		})
		return
	}

	// Call service untuk finalize material prep
	result, err := h.khazwalService.FinalizeMaterialPrep(prepID, req.Photos, req.Notes, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Material preparation tidak ditemukan",
			})
			return
		}

		if err == gorm.ErrInvalidData {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Material preparation tidak dapat diselesaikan. Pastikan semua langkah (plat, kertas, tinta) sudah selesai dan status masih IN_PROGRESS.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menyelesaikan material preparation",
			"error":   err.Error(),
		})
		return
	}

	// Transform result ke response DTO
	response := FinalizeResponse{
		PrepID:          result.PrepID,
		PONumber:        result.PONumber,
		OBCNumber:       result.OBCNumber,
		DurationMinutes: result.DurationMinutes,
		CompletedAt:     result.CompletedAt.Format("2006-01-02 15:04:05"),
		PreparedByName:  result.PreparedByName,
		PhotosCount:     result.PhotosCount,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Material preparation berhasil diselesaikan",
		"data":    response,
	})
}

// DTOs untuk API responses

// QueueItemDTO merupakan DTO untuk queue list item response
// yang mencakup essential fields untuk queue display
type QueueItemDTO struct {
	ID            uint64 `json:"id"`
	PONumber      int64  `json:"po_number"`
	OBCNumber     string `json:"obc_number"`
	ProductName   string `json:"product_name"`
	Priority      string `json:"priority"`
	PriorityScore int    `json:"priority_score"`
	DueDate       string `json:"due_date"`
	DaysUntilDue  int    `json:"days_until_due"`
	IsPastDue     bool   `json:"is_past_due"`
	CurrentStatus string `json:"current_status"`
	Quantity      int    `json:"quantity_ordered"`
	PrepStatus    string `json:"prep_status"`
}

// DetailDTO merupakan DTO untuk detail PO response
// yang mencakup full PO data dengan computed fields dan relations
type DetailDTO struct {
	ID                        uint64      `json:"id"`
	PONumber                  int64       `json:"po_number"`
	OBCNumber                 string      `json:"obc_number"`
	SAPCustomerCode           string      `json:"sap_customer_code"`
	SAPProductCode            string      `json:"sap_product_code"`
	ProductName               string      `json:"product_name"`
	ProductSpecifications     interface{} `json:"product_specifications"`
	QuantityOrdered           int         `json:"quantity_ordered"`
	QuantityTargetLembarBesar int         `json:"quantity_target_lembar_besar"`
	EstimatedRims             int         `json:"estimated_rims"`
	OrderDate                 string      `json:"order_date"`
	DueDate                   string      `json:"due_date"`
	Priority                  string      `json:"priority"`
	PriorityScore             int         `json:"priority_score"`
	DaysUntilDue              int         `json:"days_until_due"`
	IsPastDue                 bool        `json:"is_past_due"`
	CurrentStage              string      `json:"current_stage"`
	CurrentStatus             string      `json:"current_status"`
	Notes                     string      `json:"notes"`
	KhazwalMaterialPrep       interface{} `json:"khazwal_material_prep"`
	StageTracking             interface{} `json:"stage_tracking"`
}

// ConfirmPlatRequest merupakan DTO untuk confirm plat request
type ConfirmPlatRequest struct {
	PlatCode string `json:"plat_code" binding:"required"`
}

// UpdateKertasRequest merupakan DTO untuk update kertas blanko request
type UpdateKertasRequest struct {
	ActualQty      int    `json:"actual_qty" binding:"required,min=0"`
	VarianceReason string `json:"variance_reason"`
}

// UpdateTintaRequest merupakan DTO untuk update tinta request
type UpdateTintaRequest struct {
	TintaActual []TintaItem `json:"tinta_actual" binding:"required"`
}

// TintaItem merupakan struct untuk individual tinta item
type TintaItem struct {
	Color    string  `json:"color" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required,min=0"`
	Checked  bool    `json:"checked"`
}

// FinalizeRequest merupakan DTO untuk finalize material preparation request
// yang mencakup optional photos (base64) dan notes
type FinalizeRequest struct {
	Photos []string `json:"photos"`
	Notes  string   `json:"notes"`
}

// FinalizeResponse merupakan DTO untuk finalize material preparation response
// yang mencakup completion summary untuk success screen display
type FinalizeResponse struct {
	PrepID          uint64 `json:"prep_id"`
	PONumber        int64  `json:"po_number"`
	OBCNumber       string `json:"obc_number"`
	DurationMinutes int    `json:"duration_minutes"`
	CompletedAt     string `json:"completed_at"`
	PreparedByName  string `json:"prepared_by_name"`
	PhotosCount     int    `json:"photos_count"`
}

// GetHistory mengambil riwayat material preparation yang sudah selesai
// dengan filter by date range dan staff
// @route GET /api/khazwal/material-prep/history
// @access STAFF_KHAZWAL, SUPERVISOR_KHAZWAL, ADMIN, MANAGER
func (h *KhazwalHandler) GetHistory(c *gin.Context) {
	// Parse query filters dari request
	var filters services.HistoryFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter query tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Call service untuk get history
	historyResponse, err := h.khazwalService.GetMaterialPrepHistory(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil riwayat material preparation",
			"error":   err.Error(),
		})
		return
	}

	// Transform response ke DTO format
	items := make([]HistoryItemDTO, 0, len(historyResponse.Items))
	for _, item := range historyResponse.Items {
		dto := HistoryItemDTO{
			PrepID:          item.PrepID,
			POID:            item.POID,
			PONumber:        item.PONumber,
			OBCNumber:       item.OBCNumber,
			ProductName:     item.ProductName,
			Priority:        item.Priority,
			Quantity:        item.Quantity,
			CompletedAt:     item.CompletedAt.Format("2006-01-02 15:04:05"),
			DurationMinutes: item.DurationMinutes,
			PreparedByID:    item.PreparedByID,
			PreparedByName:  item.PreparedByName,
			PhotosCount:     item.PhotosCount,
		}
		items = append(items, dto)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Riwayat material preparation berhasil diambil",
		"data": gin.H{
			"items":       items,
			"total":       historyResponse.Total,
			"page":        historyResponse.Page,
			"per_page":    historyResponse.PerPage,
			"total_pages": historyResponse.TotalPages,
		},
	})
}

// GetMonitoring mengambil statistik monitoring untuk supervisor dashboard
// @route GET /api/khazwal/monitoring
// @access SUPERVISOR_KHAZWAL, ADMIN, MANAGER
func (h *KhazwalHandler) GetMonitoring(c *gin.Context) {
	// Call service untuk get monitoring stats
	stats, err := h.khazwalService.GetMonitoringStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil statistik monitoring",
			"error":   err.Error(),
		})
		return
	}

	// Transform staff active ke DTO
	staffActive := make([]StaffActivityDTO, 0, len(stats.StaffActive))
	for _, staff := range stats.StaffActive {
		dto := StaffActivityDTO{
			UserID:       staff.UserID,
			Name:         staff.Name,
			CurrentPOID:  staff.CurrentPOID,
			CurrentPO:    staff.CurrentPO,
			ProductName:  staff.ProductName,
			Status:       staff.Status,
			DurationMins: staff.DurationMins,
		}
		if staff.StartedAt != nil {
			startedAtStr := staff.StartedAt.Format("2006-01-02 15:04:05")
			dto.StartedAt = &startedAtStr
		}
		staffActive = append(staffActive, dto)
	}

	// Transform recent completions ke DTO
	recentCompletions := make([]RecentCompleteDTO, 0, len(stats.RecentCompletions))
	for _, recent := range stats.RecentCompletions {
		dto := RecentCompleteDTO{
			PrepID:          recent.PrepID,
			PONumber:        recent.PONumber,
			OBCNumber:       recent.OBCNumber,
			CompletedAt:     recent.CompletedAt.Format("2006-01-02 15:04:05"),
			DurationMinutes: recent.DurationMinutes,
			PreparedByName:  recent.PreparedByName,
		}
		recentCompletions = append(recentCompletions, dto)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Statistik monitoring berhasil diambil",
		"data": gin.H{
			"total_in_queue":        stats.TotalInQueue,
			"total_in_progress":     stats.TotalInProgress,
			"total_completed_today": stats.TotalCompletedToday,
			"average_duration_mins": stats.AverageDurationMins,
			"staff_active":          staffActive,
			"recent_completions":    recentCompletions,
		},
	})
}

// HistoryItemDTO merupakan DTO untuk history item response
type HistoryItemDTO struct {
	PrepID          uint64 `json:"prep_id"`
	POID            uint64 `json:"po_id"`
	PONumber        int64  `json:"po_number"`
	OBCNumber       string `json:"obc_number"`
	ProductName     string `json:"product_name"`
	Priority        string `json:"priority"`
	Quantity        int    `json:"quantity"`
	CompletedAt     string `json:"completed_at"`
	DurationMinutes int    `json:"duration_minutes"`
	PreparedByID    uint64 `json:"prepared_by_id"`
	PreparedByName  string `json:"prepared_by_name"`
	PhotosCount     int    `json:"photos_count"`
}

// StaffActivityDTO merupakan DTO untuk staff activity response
type StaffActivityDTO struct {
	UserID       uint64  `json:"user_id"`
	Name         string  `json:"name"`
	CurrentPOID  *uint64 `json:"current_po_id,omitempty"`
	CurrentPO    string  `json:"current_po,omitempty"`
	ProductName  string  `json:"product_name,omitempty"`
	StartedAt    *string `json:"started_at,omitempty"`
	Status       string  `json:"status"`
	DurationMins int     `json:"duration_mins"`
}

// RecentCompleteDTO merupakan DTO untuk recent completions response
type RecentCompleteDTO struct {
	PrepID          uint64 `json:"prep_id"`
	PONumber        int64  `json:"po_number"`
	OBCNumber       string `json:"obc_number"`
	CompletedAt     string `json:"completed_at"`
	DurationMinutes int    `json:"duration_minutes"`
	PreparedByName  string `json:"prepared_by_name"`
}
