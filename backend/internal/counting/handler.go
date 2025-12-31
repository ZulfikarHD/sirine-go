package counting

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CountingHandler merupakan HTTP handler untuk counting endpoints
type CountingHandler struct {
	service CountingService
	db      *gorm.DB
}

// NewCountingHandler membuat instance baru CountingHandler
func NewCountingHandler(db *gorm.DB) *CountingHandler {
	repo := NewCountingRepository(db)
	service := NewCountingService(db, repo)
	
	return &CountingHandler{
		service: service,
		db:      db,
	}
}

// GetCountingQueue menghandle GET /api/khazwal/counting/queue
// untuk mengambil list PO yang menunggu penghitungan (FIFO sorted)
func (h *CountingHandler) GetCountingQueue(c *gin.Context) {
	// Parse optional query parameters
	var machineID *uint64
	if machineIDStr := c.Query("machine_id"); machineIDStr != "" {
		id, err := strconv.ParseUint(machineIDStr, 10, 64)
		if err == nil {
			machineID = &id
		}
	}

	var dateFrom, dateTo *time.Time
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = &parsed
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = &parsed
		}
	}

	// Get queue dari service
	response, err := h.service.GetCountingQueue(machineID, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil counting queue",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response.Data,
		"meta":    response.Meta,
	})
}

// GetCountingDetail menghandle GET /api/khazwal/counting/:id
// untuk mengambil detail counting record dengan relasi
func (h *CountingHandler) GetCountingDetail(c *gin.Context) {
	// Parse ID dari URL params
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	// Get detail dari service
	detail, err := h.service.GetCountingDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Counting record tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    detail,
	})
}

// StartCounting menghandle POST /api/khazwal/counting/:po_id/start
// untuk memulai proses penghitungan dengan create counting record dan update PO status
func (h *CountingHandler) StartCounting(c *gin.Context) {
	// Parse PO ID dari URL params
	poID, err := strconv.ParseUint(c.Param("po_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "PO ID tidak valid",
		})
		return
	}

	// Get user ID dari auth context (middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak terautentikasi",
		})
		return
	}

	// Start counting via service
	response, err := h.service.StartCounting(poID, userID.(uint64))
	if err != nil {
		// Determine appropriate status code
		statusCode := http.StatusInternalServerError
		if err == ErrPONotReadyForCounting {
			statusCode = http.StatusBadRequest
		} else if err == ErrCountingAlreadyExists {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Penghitungan berhasil dimulai",
		"data":    response,
	})
}

// UpdateCountingResult menghandle PATCH /api/khazwal/counting/:id/result
// untuk update hasil penghitungan (dapat dipanggil multiple times sebelum finalize)
func (h *CountingHandler) UpdateCountingResult(c *gin.Context) {
	// Parse counting ID dari URL params
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	// Parse request body
	var req UpdateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Request body tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Get counting untuk ambil PO ID
	counting, err := h.service.GetCountingDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Counting record tidak ditemukan",
		})
		return
	}

	// Query target quantity dari PO
	var po struct {
		QuantityTargetLembarBesar int
	}
	if err := h.db.Table("production_orders").
		Select("quantity_target_lembar_besar").
		Where("id = ?", counting.ProductionOrderID).
		Scan(&po).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil target quantity",
		})
		return
	}

	// Update result via service
	response, err := h.service.UpdateResult(id, req, po.QuantityTargetLembarBesar)
	if err != nil {
		// Determine status code berdasarkan error type
		statusCode := http.StatusInternalServerError
		if err == ErrInvalidQuantity || 
		   err == ErrDefectBreakdownRequired || 
		   err == ErrDefectBreakdownSumMismatch || 
		   err == ErrVarianceReasonRequired {
			statusCode = http.StatusUnprocessableEntity
		} else if err == ErrCountingNotInProgress || err == ErrCountingAlreadyCompleted {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hasil penghitungan berhasil disimpan",
		"data":    response,
	})
}

// FinalizeCounting menghandle POST /api/khazwal/counting/:id/finalize
// untuk finalize counting dengan lock data dan advance PO ke next stage
func (h *CountingHandler) FinalizeCounting(c *gin.Context) {
	// Parse counting ID dari URL params
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	// Finalize via service
	response, err := h.service.FinalizeCounting(id)
	if err != nil {
		// Determine status code
		statusCode := http.StatusInternalServerError
		if err == ErrRequiredFieldsMissing || 
		   err == ErrDefectBreakdownRequired || 
		   err == ErrVarianceReasonRequired ||
		   err == ErrCountingNotInProgress {
			statusCode = http.StatusBadRequest
		} else if err == ErrDefectBreakdownSumMismatch {
			statusCode = http.StatusUnprocessableEntity
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Penghitungan berhasil diselesaikan. PO siap untuk pemotongan.",
		"data":    response,
	})
}
