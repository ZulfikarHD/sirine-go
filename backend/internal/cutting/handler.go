package cutting

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler merupakan HTTP handler untuk cutting endpoints
type Handler struct {
	service Service
}

// NewHandler membuat instance baru dari Handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetCuttingQueue menangani GET /api/khazwal/cutting/queue
// @Summary Get cutting queue
// @Description Mengambil list PO yang siap untuk dipotong dengan filters
// @Tags Cutting
// @Accept json
// @Produce json
// @Param priority query string false "Filter by priority (URGENT, HIGH, NORMAL, LOW)"
// @Param date_from query string false "Filter by date from (YYYY-MM-DD)"
// @Param date_to query string false "Filter by date to (YYYY-MM-DD)"
// @Param sort_by query string false "Sort by field (priority, date)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {object} QueueResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/khazwal/cutting/queue [get]
func (h *Handler) GetCuttingQueue(c *gin.Context) {
	// Parse filters from query params
	var filters QueueFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}
	
	// Get queue dari service
	response, err := h.service.GetCuttingQueue(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get cutting queue",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// GetCuttingDetail menangani GET /api/khazwal/cutting/:id
// @Summary Get cutting detail
// @Description Mengambil detail cutting record berdasarkan ID
// @Tags Cutting
// @Accept json
// @Produce json
// @Param id path int true "Cutting ID"
// @Success 200 {object} CuttingDetailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/khazwal/cutting/{id} [get]
func (h *Handler) GetCuttingDetail(c *gin.Context) {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid cutting ID",
			"details": "ID must be a positive integer",
		})
		return
	}
	
	// Get detail dari service
	response, err := h.service.GetCuttingDetail(id)
	if err != nil {
		if err == ErrCuttingNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cutting record not found",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get cutting detail",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// StartCutting menangani POST /api/khazwal/cutting/po/:po_id/start
// @Summary Start cutting process
// @Description Memulai proses pemotongan untuk PO tertentu
// @Tags Cutting
// @Accept json
// @Produce json
// @Param po_id path int true "Production Order ID"
// @Param request body StartCuttingRequest true "Start cutting request"
// @Success 200 {object} StartCuttingResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/khazwal/cutting/po/{po_id}/start [post]
func (h *Handler) StartCutting(c *gin.Context) {
	// Parse PO ID from URL
	poIDStr := c.Param("po_id")
	poID, err := strconv.ParseUint(poIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid production order ID",
			"details": "ID must be a positive integer",
		})
		return
	}
	
	// Parse request body
	var req StartCuttingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	// Get user ID dari context (dari auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	
	userIDUint64, ok := userID.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}
	
	// Start cutting via service
	response, err := h.service.StartCutting(poID, req, userIDUint64)
	if err != nil {
		// Handle specific errors
		switch err {
		case ErrPONotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Production order not found",
			})
		case ErrPONotReadyForCutting:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "PO is not ready for cutting",
			})
		case ErrCuttingAlreadyStarted:
			c.JSON(http.StatusConflict, gin.H{
				"error": "Cutting already started for this PO",
			})
		case ErrCountingNotCompleted:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Counting result not completed yet",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to start cutting",
				"details": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// UpdateCuttingResult menangani PATCH /api/khazwal/cutting/:id/result
// @Summary Update cutting result
// @Description Mengupdate hasil pemotongan (sisiran kiri & kanan)
// @Tags Cutting
// @Accept json
// @Produce json
// @Param id path int true "Cutting ID"
// @Param request body UpdateResultRequest true "Update result request"
// @Success 200 {object} UpdateResultResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/khazwal/cutting/{id}/result [patch]
func (h *Handler) UpdateCuttingResult(c *gin.Context) {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid cutting ID",
			"details": "ID must be a positive integer",
		})
		return
	}
	
	// Parse request body
	var req UpdateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	// Update result via service
	response, err := h.service.UpdateCuttingResult(id, req)
	if err != nil {
		switch err {
		case ErrCuttingNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cutting record not found",
			})
		case ErrCuttingNotInProgress:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cutting is not in progress",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update cutting result",
				"details": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// FinalizeCutting menangani POST /api/khazwal/cutting/:id/finalize
// @Summary Finalize cutting
// @Description Menyelesaikan proses pemotongan dan generate verification labels
// @Tags Cutting
// @Accept json
// @Produce json
// @Param id path int true "Cutting ID"
// @Success 200 {object} FinalizeCuttingResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/khazwal/cutting/{id}/finalize [post]
func (h *Handler) FinalizeCutting(c *gin.Context) {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid cutting ID",
			"details": "ID must be a positive integer",
		})
		return
	}
	
	// Finalize via service
	response, err := h.service.FinalizeCutting(id)
	if err != nil {
		switch err {
		case ErrCuttingNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cutting record not found",
			})
		case ErrCuttingNotInProgress:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cutting is not in progress",
			})
		case ErrCuttingAlreadyCompleted:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cutting already completed",
			})
		case ErrMissingOutputData:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Output data (sisiran kiri & kanan) must be filled before finalization",
			})
		case ErrMissingWasteDocumentation:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Waste exceeds 2%, reason and photo are required",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to finalize cutting",
				"details": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, response)
}
