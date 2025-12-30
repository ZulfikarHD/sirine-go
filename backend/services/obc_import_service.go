package services

import (
	"errors"
	"fmt"
	"io"
	"math"
	"sirine-go/backend/models"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// OBCImportService merupakan service untuk import data OBC Master dari Excel
// yang mencakup parsing, validation, dan upsert logic
type OBCImportService struct {
	db *gorm.DB
}

// NewOBCImportService membuat instance baru dari OBCImportService
func NewOBCImportService(db *gorm.DB) *OBCImportService {
	return &OBCImportService{db: db}
}

// ImportResult merupakan struktur hasil dari import Excel
// yang mencakup summary dan detail errors
type ImportResult struct {
	TotalRows     int         `json:"total_rows"`
	SuccessCount  int         `json:"success_count"`
	FailedCount   int         `json:"failed_count"`
	FailedRows    []FailedRow `json:"failed_rows"`
	POsGenerated  int         `json:"pos_generated"`
	DurationMs    int64       `json:"duration_ms"`
}

// FailedRow merupakan struktur untuk row yang gagal di-import
type FailedRow struct {
	RowNumber int    `json:"row_number"`
	OBCNumber string `json:"obc_number"`
	Error     string `json:"error"`
}

// excelColumnMapping merupakan mapping dari Excel column header ke field name OBCMaster
// dimana key adalah header di Excel dan value adalah field name di struct
var excelColumnMapping = map[string]string{
	"No OBC":                "OBCNumber",
	"Tgl OBC":               "OBCDate",
	"Material":              "Material",
	"SERI":                  "Seri",
	"WARNA":                 "Warna",
	"KODE_PABRIK":           "FactoryCode",
	"QTY PESAN":             "QuantityOrdered",
	"JHT":                   "JHT",
	"RPB":                   "RPB",
	"HJE":                   "HJE",
	"BPB":                   "BPB",
	"RENCET":                "Rencet",
	"Tgl JTempo":            "DueDate",
	"Perso / non Perso":     "Personalization",
	"Perekat":               "AdhesiveType",
	"GR":                    "GR",
	"No Pelat":              "PlatNumber",
	"Type":                  "Type",
	"Created On":            "CreatedOn",
	"Sales Doc.":            "SalesDocument",
	"Item":                  "ItemCode",
	"Material description":  "MaterialDescription",
	"BUn":                   "BaseUnit",
	"Gol. PCA":              "PCACategory",
	"Kadar Alkohol PCA":     "AlcoholPercentage",
	"Kadar HPTL":            "HPTLContent",
	"Kode Wilayah":          "RegionCode",
	"OBC Awal":              "OBCInitial",
	"Peruntukan":            "Allocation",
	"PESANAN":               "TotalOrderOBC",
	"Plnt":                  "PlantCode",
	"SATUAN":                "Unit",
	"Tahun":                 "ProductionYear",
	"Tarif Per Liter":       "ExciseRatePerLiter",
	"Volume PCA":            "PCAVolume",
	"Warna MMEA":            "MMEAColorCode",
}

// ImportFromExcel melakukan import data OBC dari Excel file
// dengan parsing, validation, dan upsert logic dalam transaction
func (s *OBCImportService) ImportFromExcel(fileReader io.Reader, autoGeneratePO bool) (*ImportResult, error) {
	startTime := time.Now()
	
	// Baca Excel file
	xlsxFile, err := excelize.OpenReader(fileReader)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca Excel file: %w", err)
	}
	defer xlsxFile.Close()

	// Ambil sheet pertama
	sheets := xlsxFile.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("Excel file tidak memiliki sheet")
	}
	sheetName := sheets[0]

	// Baca semua rows
	rows, err := xlsxFile.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca rows dari sheet: %w", err)
	}

	if len(rows) < 2 {
		return nil, errors.New("Excel file harus memiliki minimal 2 rows (header + data)")
	}

	// Parse header untuk mendapatkan column index mapping
	headerRow := rows[0]
	columnIndexMap := s.buildColumnIndexMap(headerRow)

	// Validasi required columns
	requiredColumns := []string{"No OBC"}
	for _, col := range requiredColumns {
		if _, exists := columnIndexMap[col]; !exists {
			return nil, fmt.Errorf("kolom required '%s' tidak ditemukan di Excel", col)
		}
	}

	// Process rows dalam transaction
	result := &ImportResult{
		TotalRows:    len(rows) - 1, // exclude header
		FailedRows:   make([]FailedRow, 0),
	}

	posGenerated := 0

	// Start transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		for i := 1; i < len(rows); i++ {
			rowData := rows[i]
			rowNumber := i + 1

			// Parse row ke OBCMaster struct
			obcMaster, err := s.parseRowToOBCMaster(rowData, columnIndexMap)
			if err != nil {
				result.FailedRows = append(result.FailedRows, FailedRow{
					RowNumber: rowNumber,
					OBCNumber: s.getValueFromRow(rowData, columnIndexMap, "No OBC"),
					Error:     err.Error(),
				})
				result.FailedCount++
				continue
			}

			// Upsert logic: update jika OBC sudah ada, create jika baru
			var existingOBC models.OBCMaster
			findErr := tx.Where("obc_number = ?", obcMaster.OBCNumber).First(&existingOBC).Error
			
			if findErr == nil {
				// Update existing
				obcMaster.ID = existingOBC.ID
				obcMaster.CreatedAt = existingOBC.CreatedAt
				if err := tx.Save(&obcMaster).Error; err != nil {
					result.FailedRows = append(result.FailedRows, FailedRow{
						RowNumber: rowNumber,
						OBCNumber: obcMaster.OBCNumber,
						Error:     fmt.Sprintf("gagal update: %v", err),
					})
					result.FailedCount++
					continue
				}
			} else if errors.Is(findErr, gorm.ErrRecordNotFound) {
				// Create new
				if err := tx.Create(&obcMaster).Error; err != nil {
					result.FailedRows = append(result.FailedRows, FailedRow{
						RowNumber: rowNumber,
						OBCNumber: obcMaster.OBCNumber,
						Error:     fmt.Sprintf("gagal create: %v", err),
					})
					result.FailedCount++
					continue
				}
			} else {
				// Database error
				result.FailedRows = append(result.FailedRows, FailedRow{
					RowNumber: rowNumber,
					OBCNumber: obcMaster.OBCNumber,
					Error:     fmt.Sprintf("database error: %v", findErr),
				})
				result.FailedCount++
				continue
			}

			result.SuccessCount++

			// Auto generate PO jika diminta
			if autoGeneratePO && obcMaster.QuantityOrdered > 0 {
				pos, err := s.generatePOsFromOBCInTx(tx, obcMaster.ID)
				if err == nil {
					posGenerated += len(pos)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction error: %w", err)
	}

	result.POsGenerated = posGenerated
	result.DurationMs = time.Since(startTime).Milliseconds()

	return result, nil
}

// buildColumnIndexMap membuat mapping dari column name ke column index
func (s *OBCImportService) buildColumnIndexMap(headerRow []string) map[string]int {
	indexMap := make(map[string]int)
	for i, header := range headerRow {
		// Trim whitespace dan normalize
		normalizedHeader := strings.TrimSpace(header)
		indexMap[normalizedHeader] = i
	}
	return indexMap
}

// getValueFromRow mengambil value dari row berdasarkan column name
func (s *OBCImportService) getValueFromRow(rowData []string, columnIndexMap map[string]int, columnName string) string {
	if idx, exists := columnIndexMap[columnName]; exists && idx < len(rowData) {
		return strings.TrimSpace(rowData[idx])
	}
	return ""
}

// parseRowToOBCMaster melakukan parsing dari row data ke OBCMaster struct
// dengan type conversion dan validation
func (s *OBCImportService) parseRowToOBCMaster(rowData []string, columnIndexMap map[string]int) (*models.OBCMaster, error) {
	obcMaster := &models.OBCMaster{}

	// OBCNumber (required)
	obcNumber := s.getValueFromRow(rowData, columnIndexMap, "No OBC")
	if obcNumber == "" {
		return nil, errors.New("No OBC tidak boleh kosong")
	}
	obcMaster.OBCNumber = obcNumber

	// String fields
	obcMaster.Material = s.getValueFromRow(rowData, columnIndexMap, "Material")
	obcMaster.Seri = s.getValueFromRow(rowData, columnIndexMap, "SERI")
	obcMaster.Warna = s.getValueFromRow(rowData, columnIndexMap, "WARNA")
	obcMaster.FactoryCode = s.getValueFromRow(rowData, columnIndexMap, "KODE_PABRIK")
	obcMaster.JHT = s.getValueFromRow(rowData, columnIndexMap, "JHT")
	obcMaster.Personalization = s.getValueFromRow(rowData, columnIndexMap, "Perso / non Perso")
	obcMaster.AdhesiveType = s.getValueFromRow(rowData, columnIndexMap, "Perekat")
	obcMaster.GR = s.getValueFromRow(rowData, columnIndexMap, "GR")
	obcMaster.PlatNumber = s.getValueFromRow(rowData, columnIndexMap, "No Pelat")
	obcMaster.Type = s.getValueFromRow(rowData, columnIndexMap, "Type")
	obcMaster.SalesDocument = s.getValueFromRow(rowData, columnIndexMap, "Sales Doc.")
	obcMaster.ItemCode = s.getValueFromRow(rowData, columnIndexMap, "Item")
	obcMaster.MaterialDescription = s.getValueFromRow(rowData, columnIndexMap, "Material description")
	obcMaster.BaseUnit = s.getValueFromRow(rowData, columnIndexMap, "BUn")
	obcMaster.PCACategory = s.getValueFromRow(rowData, columnIndexMap, "Gol. PCA")
	obcMaster.RegionCode = s.getValueFromRow(rowData, columnIndexMap, "Kode Wilayah")
	obcMaster.OBCInitial = s.getValueFromRow(rowData, columnIndexMap, "OBC Awal")
	obcMaster.Allocation = s.getValueFromRow(rowData, columnIndexMap, "Peruntukan")
	obcMaster.PlantCode = s.getValueFromRow(rowData, columnIndexMap, "Plnt")
	obcMaster.Unit = s.getValueFromRow(rowData, columnIndexMap, "SATUAN")
	obcMaster.MMEAColorCode = s.getValueFromRow(rowData, columnIndexMap, "Warna MMEA")

	// Integer fields
	if val := s.getValueFromRow(rowData, columnIndexMap, "QTY PESAN"); val != "" {
		if intVal, err := s.parseInteger(val); err == nil {
			obcMaster.QuantityOrdered = intVal
		}
	}
	if val := s.getValueFromRow(rowData, columnIndexMap, "BPB"); val != "" {
		if intVal, err := s.parseInteger(val); err == nil {
			obcMaster.BPB = intVal
		}
	}
	if val := s.getValueFromRow(rowData, columnIndexMap, "RENCET"); val != "" {
		if intVal, err := s.parseInteger(val); err == nil {
			obcMaster.Rencet = intVal
		}
	}
	if val := s.getValueFromRow(rowData, columnIndexMap, "PESANAN"); val != "" {
		if intVal, err := s.parseInteger(val); err == nil {
			obcMaster.TotalOrderOBC = intVal
		}
	}
	if val := s.getValueFromRow(rowData, columnIndexMap, "Tahun"); val != "" {
		if intVal, err := s.parseInteger(val); err == nil {
			obcMaster.ProductionYear = intVal
		}
	}

	// Float/Decimal fields
	if val := s.getValueFromRow(rowData, columnIndexMap, "RPB"); val != "" {
		if floatVal, err := s.parseFloat(val); err == nil {
			obcMaster.RPB = floatVal
		}
	}
	if val := s.getValueFromRow(rowData, columnIndexMap, "HJE"); val != "" {
		if floatVal, err := s.parseFloat(val); err == nil {
			obcMaster.HJE = floatVal
		}
	}
	if val := s.getValueFromRow(rowData, columnIndexMap, "Kadar Alkohol PCA"); val != "" {
		if floatVal, err := s.parseFloat(val); err == nil {
			obcMaster.AlcoholPercentage = floatVal
		}
	}
	if val := s.getValueFromRow(rowData, columnIndexMap, "Kadar HPTL"); val != "" {
		if floatVal, err := s.parseFloat(val); err == nil {
			obcMaster.HPTLContent = floatVal
		}
	}
	if val := s.getValueFromRow(rowData, columnIndexMap, "Tarif Per Liter"); val != "" {
		if floatVal, err := s.parseFloat(val); err == nil {
			obcMaster.ExciseRatePerLiter = floatVal
		}
	}
	if val := s.getValueFromRow(rowData, columnIndexMap, "Volume PCA"); val != "" {
		if floatVal, err := s.parseFloat(val); err == nil {
			obcMaster.PCAVolume = floatVal
		}
	}

	// Date fields
	if val := s.getValueFromRow(rowData, columnIndexMap, "Tgl OBC"); val != "" {
		if dateVal, err := s.parseDate(val); err == nil {
			obcMaster.OBCDate = dateVal
		}
	}
	if val := s.getValueFromRow(rowData, columnIndexMap, "Tgl JTempo"); val != "" {
		if dateVal, err := s.parseDate(val); err == nil {
			obcMaster.DueDate = dateVal
		}
	}
	if val := s.getValueFromRow(rowData, columnIndexMap, "Created On"); val != "" {
		if dateVal, err := s.parseDate(val); err == nil {
			obcMaster.CreatedOn = dateVal
		}
	}

	return obcMaster, nil
}

// parseInteger melakukan parsing string ke integer dengan handling berbagai format
func (s *OBCImportService) parseInteger(val string) (int, error) {
	// Remove thousands separator
	val = strings.ReplaceAll(val, ",", "")
	val = strings.ReplaceAll(val, ".", "")
	val = strings.TrimSpace(val)
	
	if val == "" {
		return 0, nil
	}
	
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %s", val)
	}
	return intVal, nil
}

// parseFloat melakukan parsing string ke float dengan handling berbagai format
func (s *OBCImportService) parseFloat(val string) (float64, error) {
	// Replace comma dengan dot untuk decimal separator
	val = strings.ReplaceAll(val, ",", ".")
	val = strings.TrimSpace(val)
	
	if val == "" {
		return 0, nil
	}
	
	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float: %s", val)
	}
	return floatVal, nil
}

// parseDate melakukan parsing string atau Excel serial date ke time.Time
func (s *OBCImportService) parseDate(val string) (*time.Time, error) {
	val = strings.TrimSpace(val)
	if val == "" {
		return nil, nil
	}

	// Try parsing as Excel serial date (numeric)
	if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
		// Excel serial date: days since 1900-01-01
		baseDate := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		date := baseDate.Add(time.Duration(floatVal * 24 * float64(time.Hour)))
		return &date, nil
	}

	// Try common date formats
	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"02-01-2006",
		"02.01.2006",
		"2006/01/02",
		"01/02/2006",
	}

	for _, format := range formats {
		if date, err := time.Parse(format, val); err == nil {
			return &date, nil
		}
	}

	return nil, fmt.Errorf("invalid date format: %s", val)
}

// GeneratePOsFromOBC melakukan generate Production Orders dari OBC Master
// dengan formula: Total = QTY + (QTY * 6%), PO Count = CEIL(Total / 40000)
func (s *OBCImportService) GeneratePOsFromOBC(obcID uint64) ([]models.ProductionOrder, error) {
	return s.generatePOsFromOBCInTx(s.db, obcID)
}

// generatePOsFromOBCInTx melakukan generate PO dalam transaction context
func (s *OBCImportService) generatePOsFromOBCInTx(tx *gorm.DB, obcID uint64) ([]models.ProductionOrder, error) {
	// Load OBC Master
	var obc models.OBCMaster
	if err := tx.First(&obc, obcID).Error; err != nil {
		return nil, fmt.Errorf("OBC Master tidak ditemukan: %w", err)
	}

	if obc.QuantityOrdered <= 0 {
		return nil, errors.New("QuantityOrdered harus lebih dari 0")
	}

	// Calculate total dengan buffer 6%
	totalWithBuffer := obc.CalculateTotalWithBuffer()
	
	// Calculate jumlah PO yang diperlukan (max 40000 per PO)
	const maxQuantityPerPO = 40000
	poCount := int(math.Ceil(float64(totalWithBuffer) / float64(maxQuantityPerPO)))

	// Generate PO number base (timestamp-based untuk uniqueness)
	poNumberBase := time.Now().Unix()

	pos := make([]models.ProductionOrder, 0, poCount)

	// Create Production Orders
	for i := 0; i < poCount; i++ {
		// Calculate quantity untuk PO ini
		remainingQty := totalWithBuffer - (i * maxQuantityPerPO)
		poQty := maxQuantityPerPO
		if remainingQty < maxQuantityPerPO {
			poQty = remainingQty
		}

		// Generate unique PO Number
		poNumber := poNumberBase + int64(i+1)

		// Set due date dari OBC atau 30 hari dari sekarang
		dueDate := time.Now().AddDate(0, 0, 30)
		if obc.DueDate != nil {
			dueDate = *obc.DueDate
		}

		// Build product specifications dari OBC data
		productSpecs := map[string]interface{}{
			"material":              obc.Material,
			"seri":                  obc.Seri,
			"warna":                 obc.Warna,
			"factory_code":          obc.FactoryCode,
			"personalization":       obc.Personalization,
			"adhesive_type":         obc.AdhesiveType,
			"type":                  obc.Type,
			"material_description":  obc.MaterialDescription,
			"pca_category":          obc.PCACategory,
			"alcohol_percentage":    obc.AlcoholPercentage,
		}

		po := models.ProductionOrder{
			PONumber:                  poNumber,
			OBCMasterID:               obc.ID,
			// Denormalized fields dari OBCMaster
			OBCNumber:                 obc.OBCNumber,
			ProductName:               obc.MaterialDescription, // Gunakan MaterialDescription sebagai ProductName
			SAPCustomerCode:           obc.ItemCode,             // Map ItemCode sebagai customer code untuk sementara
			SAPProductCode:            obc.ItemCode,
			ProductSpecifications:     productSpecs,
			// PO fields
			QuantityOrdered:           poQty,
			QuantityTargetLembarBesar: poQty, // Default sama dengan ordered
			EstimatedRims:             int(math.Ceil(float64(poQty) / 1000.0)), // Estimasi rim
			OrderDate:                 time.Now(),
			DueDate:                   dueDate,
			Priority:                  models.PriorityNormal,
			PriorityScore:             50,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
		}

		// Calculate priority score
		po.UpdatePriorityScore()

		// Save PO
		if err := tx.Create(&po).Error; err != nil {
			return nil, fmt.Errorf("gagal create Production Order: %w", err)
		}

		pos = append(pos, po)
	}

	return pos, nil
}

// GetOBCMasterByID mengambil OBC Master berdasarkan ID dengan preload POs
func (s *OBCImportService) GetOBCMasterByID(id uint64) (*models.OBCMaster, error) {
	var obc models.OBCMaster
	err := s.db.Preload("ProductionOrders").First(&obc, id).Error
	if err != nil {
		return nil, err
	}
	return &obc, nil
}

// ListOBCMasters mengambil list OBC Masters dengan pagination dan filter
func (s *OBCImportService) ListOBCMasters(page, pageSize int, filters map[string]string) ([]models.OBCMaster, int64, error) {
	var obcs []models.OBCMaster
	var total int64

	query := s.db.Model(&models.OBCMaster{})

	// Apply filters
	if material := filters["material"]; material != "" {
		query = query.Where("material LIKE ?", "%"+material+"%")
	}
	if seri := filters["seri"]; seri != "" {
		query = query.Where("seri LIKE ?", "%"+seri+"%")
	}
	if warna := filters["warna"]; warna != "" {
		query = query.Where("warna LIKE ?", "%"+warna+"%")
	}
	if factoryCode := filters["factory_code"]; factoryCode != "" {
		query = query.Where("factory_code LIKE ?", "%"+factoryCode+"%")
	}
	if obcNumber := filters["obc_number"]; obcNumber != "" {
		query = query.Where("obc_number LIKE ?", "%"+obcNumber+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&obcs).Error; err != nil {
		return nil, 0, err
	}

	return obcs, total, nil
}
