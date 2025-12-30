package handlers

import (
	"net/http"
	"sirine-go/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OBCHandler merupakan handler untuk OBC Master endpoints
// yang mencakup import Excel, list, detail, dan generate PO
type OBCHandler struct {
	obcService *services.OBCImportService
}

// NewOBCHandler membuat instance baru dari OBCHandler
func NewOBCHandler(obcService *services.OBCImportService) *OBCHandler {
	return &OBCHandler{
		obcService: obcService,
	}
}

// Import melakukan upload dan import Excel file OBC Master
// dengan parsing, validation, dan optional auto PO generation
// @route POST /api/obc/import
// @access ADMIN, PPIC
func (h *OBCHandler) Import(c *gin.Context) {
	// Parse multipart form untuk file upload
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File tidak ditemukan",
			"error":   "Parameter 'file' diperlukan untuk upload Excel",
		})
		return
	}
	defer file.Close()

	// Validate file extension
	if header.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format file tidak valid",
			"error":   "Hanya file Excel (.xlsx) yang diperbolehkan",
		})
		return
	}

	// Parse auto_generate_po flag (default: false)
	autoGeneratePO := c.DefaultQuery("auto_generate_po", "false") == "true"

	// Call service untuk import
	result, err := h.obcService.ImportFromExcel(file, autoGeneratePO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal import Excel file",
			"error":   err.Error(),
		})
		return
	}

	// Return result dengan detail success/failed rows
	statusCode := http.StatusOK
	message := "Import berhasil"
	if result.FailedCount > 0 {
		statusCode = http.StatusPartialContent
		message = "Import selesai dengan beberapa error"
	}

	c.JSON(statusCode, gin.H{
		"success": result.FailedCount == 0,
		"message": message,
		"data": gin.H{
			"total_rows":     result.TotalRows,
			"success_count":  result.SuccessCount,
			"failed_count":   result.FailedCount,
			"failed_rows":    result.FailedRows,
			"pos_generated":  result.POsGenerated,
			"duration_ms":    result.DurationMs,
			"file_name":      header.Filename,
			"file_size":      header.Size,
		},
	})
}

// List mengambil list OBC Masters dengan pagination dan filter
// untuk display table atau dropdown selection
// @route GET /api/obc
// @access ADMIN, PPIC, MANAGER
func (h *OBCHandler) List(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Parse filter parameters
	filters := map[string]string{
		"material":     c.Query("material"),
		"seri":         c.Query("seri"),
		"warna":        c.Query("warna"),
		"factory_code": c.Query("factory_code"),
		"obc_number":   c.Query("obc_number"),
	}

	// Call service untuk get list
	obcs, total, err := h.obcService.ListOBCMasters(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil list OBC Master",
			"error":   err.Error(),
		})
		return
	}

	// Transform ke DTO untuk clean response
	obcItems := make([]OBCMasterListDTO, 0, len(obcs))
	for _, obc := range obcs {
		var obcDate, dueDate, createdOn string
		if obc.OBCDate != nil {
			obcDate = obc.OBCDate.Format("2006-01-02")
		}
		if obc.DueDate != nil {
			dueDate = obc.DueDate.Format("2006-01-02")
		}
		if obc.CreatedOn != nil {
			createdOn = obc.CreatedOn.Format("2006-01-02")
		}

		obcItems = append(obcItems, OBCMasterListDTO{
			ID:                  obc.ID,
			OBCNumber:           obc.OBCNumber,
			OBCDate:             obcDate,
			Material:            obc.Material,
			Seri:                obc.Seri,
			Warna:               obc.Warna,
			FactoryCode:         obc.FactoryCode,
			QuantityOrdered:     obc.QuantityOrdered,
			MaterialDescription: obc.MaterialDescription,
			DueDate:             dueDate,
			Personalization:     obc.Personalization,
			CreatedOn:           createdOn,
			CreatedAt:           obc.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Calculate pagination metadata
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "List OBC Master berhasil diambil",
		"data": gin.H{
			"items":       obcItems,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}

// Detail mengambil detail OBC Master beserta Production Orders terkait
// untuk display complete information
// @route GET /api/obc/:id
// @access ADMIN, PPIC, MANAGER
func (h *OBCHandler) Detail(c *gin.Context) {
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

	// Call service untuk get detail dengan preload POs
	obc, err := h.obcService.GetOBCMasterByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "OBC Master tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	// Transform ke DTO dengan full details
	var obcDate, dueDate, createdOn string
	if obc.OBCDate != nil {
		obcDate = obc.OBCDate.Format("2006-01-02")
	}
	if obc.DueDate != nil {
		dueDate = obc.DueDate.Format("2006-01-02")
	}
	if obc.CreatedOn != nil {
		createdOn = obc.CreatedOn.Format("2006-01-02")
	}

	// Transform Production Orders
	poItems := make([]POSummaryDTO, 0, len(obc.ProductionOrders))
	for _, po := range obc.ProductionOrders {
		poItems = append(poItems, POSummaryDTO{
			ID:              po.ID,
			PONumber:        po.PONumber,
			QuantityOrdered: po.QuantityOrdered,
			Priority:        string(po.Priority),
			CurrentStage:    string(po.CurrentStage),
			CurrentStatus:   string(po.CurrentStatus),
			DueDate:         po.DueDate.Format("2006-01-02"),
			CreatedAt:       po.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	obcDetail := OBCMasterDetailDTO{
		ID:                   obc.ID,
		OBCNumber:            obc.OBCNumber,
		OBCDate:              obcDate,
		Material:             obc.Material,
		Seri:                 obc.Seri,
		Warna:                obc.Warna,
		FactoryCode:          obc.FactoryCode,
		QuantityOrdered:      obc.QuantityOrdered,
		JHT:                  obc.JHT,
		RPB:                  obc.RPB,
		HJE:                  obc.HJE,
		BPB:                  obc.BPB,
		Rencet:               obc.Rencet,
		DueDate:              dueDate,
		Personalization:      obc.Personalization,
		AdhesiveType:         obc.AdhesiveType,
		GR:                   obc.GR,
		PlatNumber:           obc.PlatNumber,
		Type:                 obc.Type,
		CreatedOn:            createdOn,
		SalesDocument:        obc.SalesDocument,
		ItemCode:             obc.ItemCode,
		MaterialDescription:  obc.MaterialDescription,
		BaseUnit:             obc.BaseUnit,
		PCACategory:          obc.PCACategory,
		AlcoholPercentage:    obc.AlcoholPercentage,
		HPTLContent:          obc.HPTLContent,
		RegionCode:           obc.RegionCode,
		OBCInitial:           obc.OBCInitial,
		Allocation:           obc.Allocation,
		TotalOrderOBC:        obc.TotalOrderOBC,
		PlantCode:            obc.PlantCode,
		Unit:                 obc.Unit,
		ProductionYear:       obc.ProductionYear,
		ExciseRatePerLiter:   obc.ExciseRatePerLiter,
		PCAVolume:            obc.PCAVolume,
		MMEAColorCode:        obc.MMEAColorCode,
		CreatedAt:            obc.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:            obc.UpdatedAt.Format("2006-01-02 15:04:05"),
		ProductionOrders:     poItems,
		TotalPOs:             len(poItems),
		TotalPOQuantity:      obc.GetTotalPOQuantity(),
		TotalWithBuffer:      obc.CalculateTotalWithBuffer(),
		IsPersonalized:       obc.IsPersonalized(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail OBC Master berhasil diambil",
		"data":    obcDetail,
	})
}

// GeneratePO melakukan generate Production Orders dari OBC Master
// dengan formula otomatis berdasarkan quantity
// @route POST /api/obc/:id/generate-po
// @access ADMIN, PPIC
func (h *OBCHandler) GeneratePO(c *gin.Context) {
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

	// Call service untuk generate PO
	pos, err := h.obcService.GeneratePOsFromOBC(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal generate Production Order",
			"error":   err.Error(),
		})
		return
	}

	// Transform PO result ke summary
	poSummaries := make([]POSummaryDTO, 0, len(pos))
	for _, po := range pos {
		poSummaries = append(poSummaries, POSummaryDTO{
			ID:              po.ID,
			PONumber:        po.PONumber,
			QuantityOrdered: po.QuantityOrdered,
			Priority:        string(po.Priority),
			CurrentStage:    string(po.CurrentStage),
			CurrentStatus:   string(po.CurrentStatus),
			DueDate:         po.DueDate.Format("2006-01-02"),
			CreatedAt:       po.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Production Orders berhasil di-generate",
		"data": gin.H{
			"pos_generated":    len(pos),
			"production_orders": poSummaries,
		},
	})
}

// DTOs untuk clean API response

// OBCMasterListDTO merupakan DTO untuk list OBC Masters
type OBCMasterListDTO struct {
	ID                  uint64 `json:"id"`
	OBCNumber           string `json:"obc_number"`
	OBCDate             string `json:"obc_date"`
	Material            string `json:"material"`
	Seri                string `json:"seri"`
	Warna               string `json:"warna"`
	FactoryCode         string `json:"factory_code"`
	QuantityOrdered     int    `json:"quantity_ordered"`
	MaterialDescription string `json:"material_description"`
	DueDate             string `json:"due_date"`
	Personalization     string `json:"personalization"`
	CreatedOn           string `json:"created_on"`
	CreatedAt           string `json:"created_at"`
}

// OBCMasterDetailDTO merupakan DTO untuk detail OBC Master dengan semua fields
type OBCMasterDetailDTO struct {
	ID                  uint64          `json:"id"`
	OBCNumber           string          `json:"obc_number"`
	OBCDate             string          `json:"obc_date"`
	Material            string          `json:"material"`
	Seri                string          `json:"seri"`
	Warna               string          `json:"warna"`
	FactoryCode         string          `json:"factory_code"`
	QuantityOrdered     int             `json:"quantity_ordered"`
	JHT                 string          `json:"jht"`
	RPB                 float64         `json:"rpb"`
	HJE                 float64         `json:"hje"`
	BPB                 int             `json:"bpb"`
	Rencet              int             `json:"rencet"`
	DueDate             string          `json:"due_date"`
	Personalization     string          `json:"personalization"`
	AdhesiveType        string          `json:"adhesive_type"`
	GR                  string          `json:"gr"`
	PlatNumber          string          `json:"plat_number"`
	Type                string          `json:"type"`
	CreatedOn           string          `json:"created_on"`
	SalesDocument       string          `json:"sales_document"`
	ItemCode            string          `json:"item_code"`
	MaterialDescription string          `json:"material_description"`
	BaseUnit            string          `json:"base_unit"`
	PCACategory         string          `json:"pca_category"`
	AlcoholPercentage   float64         `json:"alcohol_percentage"`
	HPTLContent         float64         `json:"hptl_content"`
	RegionCode          string          `json:"region_code"`
	OBCInitial          string          `json:"obc_initial"`
	Allocation          string          `json:"allocation"`
	TotalOrderOBC       int             `json:"total_order_obc"`
	PlantCode           string          `json:"plant_code"`
	Unit                string          `json:"unit"`
	ProductionYear      int             `json:"production_year"`
	ExciseRatePerLiter  float64         `json:"excise_rate_per_liter"`
	PCAVolume           float64         `json:"pca_volume"`
	MMEAColorCode       string          `json:"mmea_color_code"`
	CreatedAt           string          `json:"created_at"`
	UpdatedAt           string          `json:"updated_at"`
	ProductionOrders    []POSummaryDTO  `json:"production_orders"`
	TotalPOs            int             `json:"total_pos"`
	TotalPOQuantity     int             `json:"total_po_quantity"`
	TotalWithBuffer     int             `json:"total_with_buffer"`
	IsPersonalized      bool            `json:"is_personalized"`
}

// POSummaryDTO merupakan DTO untuk summary Production Order
type POSummaryDTO struct {
	ID              uint64 `json:"id"`
	PONumber        int64  `json:"po_number"`
	QuantityOrdered int    `json:"quantity_ordered"`
	Priority        string `json:"priority"`
	CurrentStage    string `json:"current_stage"`
	CurrentStatus   string `json:"current_status"`
	DueDate         string `json:"due_date"`
	CreatedAt       string `json:"created_at"`
}
