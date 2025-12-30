package services

import (
	"encoding/json"
	"sirine-go/backend/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// KhazwalService merupakan service untuk Khazwal Material Preparation operations
// yang mencakup queue management, material prep workflow, dan tracking
type KhazwalService struct {
	db *gorm.DB
}

// NewKhazwalService membuat instance baru dari KhazwalService
func NewKhazwalService(db *gorm.DB) *KhazwalService {
	return &KhazwalService{
		db: db,
	}
}

// GetMaterialPrepQueue mengambil list PO yang menunggu material preparation
// dengan sorting berdasarkan priority score dan due date, yaitu:
// filter by status WAITING_MATERIAL_PREP dan MATERIAL_PREP_IN_PROGRESS,
// dengan pagination dan search functionality
func (s *KhazwalService) GetMaterialPrepQueue(filters QueueFilters) (*QueueResponse, error) {
	// Set default pagination values
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PerPage <= 0 || filters.PerPage > 100 {
		filters.PerPage = 20
	}

	// Query builder dengan base filters
	query := s.db.Model(&models.ProductionOrder{}).
		Where("current_status IN ?", []string{
			string(models.StatusWaitingMaterialPrep),
			string(models.StatusMaterialPrepInProgress),
		}).
		Preload("OBCMaster").
		Preload("KhazwalMaterialPrep")

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

	// Default sorting: priority_score DESC, due_date ASC
	sortBy := "priority_score"
	sortDir := "DESC"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	if filters.SortDir != "" {
		sortDir = filters.SortDir
	}

	// Execute query dengan pagination
	var pos []models.ProductionOrder
	offset := (filters.Page - 1) * filters.PerPage
	if err := query.
		Order(sortBy + " " + sortDir + ", due_date ASC").
		Limit(filters.PerPage).
		Offset(offset).
		Find(&pos).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / filters.PerPage
	if int(total)%filters.PerPage > 0 {
		totalPages++
	}

	return &QueueResponse{
		Items:      pos,
		Total:      int(total),
		Page:       filters.Page,
		PerPage:    filters.PerPage,
		TotalPages: totalPages,
	}, nil
}

// GetMaterialPrepDetail mengambil detail PO beserta material prep info
// dengan preload relations untuk full information display
func (s *KhazwalService) GetMaterialPrepDetail(id uint64) (*models.ProductionOrder, error) {
	var po models.ProductionOrder
	
	// Query dengan full relations preload
	if err := s.db.
		Preload("OBCMaster").
		Preload("KhazwalMaterialPrep.PreparedByUser").
		Preload("StageTracking", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&po, id).Error; err != nil {
		return nil, err
	}

	return &po, nil
}

// StartMaterialPrep memulai proses material preparation untuk PO
// dengan business logic validation dan transaction untuk data consistency
func (s *KhazwalService) StartMaterialPrep(poID uint64, userID uint64) (*models.ProductionOrder, error) {
	// Start transaction untuk ensure data consistency
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get PO dengan lock untuk prevent concurrent updates
	var po models.ProductionOrder
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("KhazwalMaterialPrep").
		First(&po, poID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Validasi: PO harus dalam status WAITING_MATERIAL_PREP
	if po.CurrentStatus != models.StatusWaitingMaterialPrep {
		tx.Rollback()
		return nil, gorm.ErrInvalidData
	}

	// Update PO status ke IN_PROGRESS
	now := time.Now()
	if err := tx.Model(&po).Updates(map[string]interface{}{
		"current_status": models.StatusMaterialPrepInProgress,
		"updated_at":     now,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update KhazwalMaterialPrep status dan timestamps
	if po.KhazwalMaterialPrep != nil {
		if err := tx.Model(&models.KhazwalMaterialPreparation{}).
			Where("id = ?", po.KhazwalMaterialPrep.ID).
			Updates(map[string]interface{}{
				"status":      models.MaterialPrepInProgress,
				"started_at":  now,
				"prepared_by": userID,
				"updated_at":  now,
			}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Create POStageTracking record untuk audit trail
	tracking := models.POStageTracking{
		ProductionOrderID: poID,
		Stage:             models.StageKhazwalMaterialPrep,
		Status:            models.StatusMaterialPrepInProgress,
		StartedAt:         &now,
		HandledBy:         &userID,
		Notes:             "Material preparation dimulai",
	}
	if err := tx.Create(&tracking).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload PO dengan updated data dan relations
	if err := s.db.
		Preload("OBCMaster").
		Preload("KhazwalMaterialPrep.PreparedByUser").
		Preload("StageTracking").
		First(&po, poID).Error; err != nil {
		return nil, err
	}

	return &po, nil
}

// ConfirmPlatRetrieval mengkonfirmasi pengambilan plat dengan barcode validation
// untuk memastikan plat yang diambil sesuai dengan yang ditentukan di SAP
func (s *KhazwalService) ConfirmPlatRetrieval(prepID uint64, scannedCode string, userID uint64) error {
	// Start transaction untuk ensure data consistency
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get material prep dengan lock dan preload OBCMaster untuk validation
	var prep models.KhazwalMaterialPreparation
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("ProductionOrder.OBCMaster").
		First(&prep, prepID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Validasi: prep harus dalam status IN_PROGRESS
	if !prep.IsInProgress() {
		tx.Rollback()
		return gorm.ErrInvalidData
	}

	// Validasi: OBCMaster harus ada
	if prep.ProductionOrder == nil || prep.ProductionOrder.OBCMaster == nil {
		tx.Rollback()
		return gorm.ErrInvalidData
	}

	// Compare scanned code dengan expected plat number dari OBCMaster
	expectedPlat := prep.ProductionOrder.OBCMaster.PlatNumber
	isMatch := (scannedCode == expectedPlat)

	// Update plat retrieved timestamp dan match status
	now := time.Now()
	updates := map[string]interface{}{
		"plat_retrieved_at": now,
		"plat_scanned_code": scannedCode,
		"plat_match":        isMatch,
		"updated_at":        now,
	}

	if err := tx.Model(&prep).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create POStageTracking record untuk audit trail
	notes := "Plat di-scan: " + scannedCode
	if isMatch {
		notes += " (match dengan expected)"
	} else {
		notes += " (tidak match! Expected: " + expectedPlat + ")"
	}
	
	tracking := models.POStageTracking{
		ProductionOrderID: prep.ProductionOrderID,
		Stage:             models.StageKhazwalMaterialPrep,
		Status:            models.StatusMaterialPrepInProgress,
		HandledBy:         &userID,
		Notes:             notes,
	}
	if err := tx.Create(&tracking).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// Jika tidak match, return error
	if !isMatch {
		return gorm.ErrInvalidData
	}

	return nil
}

// UpdateKertasBlanko mengupdate jumlah kertas blanko actual dengan variance calculation
// dimana variance > 5% memerlukan alasan untuk accountability
func (s *KhazwalService) UpdateKertasBlanko(prepID uint64, actualQty int, varianceReason string, userID uint64) error {
	// Start transaction untuk ensure data consistency
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get material prep dengan lock
	var prep models.KhazwalMaterialPreparation
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("ProductionOrder").
		First(&prep, prepID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Validasi: prep harus dalam status IN_PROGRESS
	if !prep.IsInProgress() {
		tx.Rollback()
		return gorm.ErrInvalidData
	}

	// Calculate variance dan variance percentage
	variance := actualQty - prep.KertasBlankoQuantity
	variancePercentage := 0.0
	if prep.KertasBlankoQuantity > 0 {
		variancePercentage = float64(variance) / float64(prep.KertasBlankoQuantity) * 100
	}

	// Validasi: jika variance > 5%, reason harus diisi
	absVariancePercentage := variancePercentage
	if absVariancePercentage < 0 {
		absVariancePercentage = -absVariancePercentage
	}
	if absVariancePercentage > 5.0 {
		if varianceReason == "" {
			tx.Rollback()
			return gorm.ErrInvalidData
		}
	}

	// Update kertas blanko actual, variance, dan variance percentage
	now := time.Now()
	updates := map[string]interface{}{
		"kertas_blanko_actual":             actualQty,
		"kertas_blanko_variance":           variance,
		"kertas_blanko_variance_percentage": variancePercentage,
		"updated_at":                       now,
	}
	if varianceReason != "" {
		updates["kertas_blanko_variance_reason"] = varianceReason
	}

	if err := tx.Model(&prep).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create POStageTracking record untuk audit trail
	notes := "Kertas blanko diupdate: " + string(rune(actualQty))
	if variance != 0 {
		notes += " (variance: " + string(rune(variance)) + ")"
	}
	tracking := models.POStageTracking{
		ProductionOrderID: prep.ProductionOrderID,
		Stage:             models.StageKhazwalMaterialPrep,
		Status:            models.StatusMaterialPrepInProgress,
		HandledBy:         &userID,
		Notes:             notes,
	}
	if err := tx.Create(&tracking).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// UpdateTinta mengupdate informasi tinta yang digunakan dengan checklist validation
// untuk memastikan semua warna tinta tersedia sesuai requirements, serta melakukan
// pengecekan low stock warning untuk inventory management
func (s *KhazwalService) UpdateTinta(prepID uint64, tintaActual interface{}, userID uint64) error {
	// Start transaction untuk ensure data consistency
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get material prep dengan lock
	var prep models.KhazwalMaterialPreparation
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("ProductionOrder").
		First(&prep, prepID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Validasi: prep harus dalam status IN_PROGRESS
	if !prep.IsInProgress() {
		tx.Rollback()
		return gorm.ErrInvalidData
	}

	// Convert tintaActual ke JSON
	tintaJSON, err := json.Marshal(tintaActual)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check for low stock warnings (TODO: integrate dengan SAP inventory API)
	lowStockFlags := s.checkLowStockTinta(tintaActual)
	lowStockJSON, err := json.Marshal(lowStockFlags)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update tinta actual dan low stock flags
	now := time.Now()
	if err := tx.Model(&prep).Updates(map[string]interface{}{
		"tinta_actual":          tintaJSON,
		"tinta_low_stock_flags": lowStockJSON,
		"updated_at":            now,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create POStageTracking record untuk audit trail
	notes := "Tinta berhasil diupdate"
	if len(lowStockFlags) > 0 {
		notes += " (dengan low stock warning)"
	}
	
	tracking := models.POStageTracking{
		ProductionOrderID: prep.ProductionOrderID,
		Stage:             models.StageKhazwalMaterialPrep,
		Status:            models.StatusMaterialPrepInProgress,
		HandledBy:         &userID,
		Notes:             notes,
	}
	if err := tx.Create(&tracking).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// checkLowStockTinta memeriksa low stock warnings untuk tinta
// dengan threshold < 10kg untuk early warning inventory management
// TODO: Integrate dengan SAP inventory API untuk real-time stock check
func (s *KhazwalService) checkLowStockTinta(tintaActual interface{}) map[string]interface{} {
	lowStockFlags := make(map[string]interface{})
	
	// Parse tintaActual untuk extract color information
	tintaMap, ok := tintaActual.(map[string]interface{})
	if !ok {
		return lowStockFlags
	}
	
	// Check each color untuk low stock (placeholder logic)
	// Real implementation harus query SAP inventory API
	colors, ok := tintaMap["colors"].([]interface{})
	if !ok {
		return lowStockFlags
	}
	
	lowStockColors := make([]string, 0)
	for _, colorData := range colors {
		colorInfo, ok := colorData.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Placeholder: check stock_after field (jika ada)
		stockAfter, ok := colorInfo["stock_after"].(float64)
		if ok && stockAfter < 10.0 {
			colorName, _ := colorInfo["color"].(string)
			lowStockColors = append(lowStockColors, colorName)
		}
	}
	
	if len(lowStockColors) > 0 {
		lowStockFlags["low_stock_colors"] = lowStockColors
		lowStockFlags["warning"] = "Beberapa warna tinta memiliki stock rendah"
	}
	
	return lowStockFlags
}

// FinalizeResult merupakan struct untuk response finalize material preparation
// yang mencakup data completion summary untuk display di success screen
type FinalizeResult struct {
	PrepID          uint64    `json:"prep_id"`
	PONumber        int64     `json:"po_number"`
	OBCNumber       string    `json:"obc_number"`
	DurationMinutes int       `json:"duration_minutes"`
	CompletedAt     time.Time `json:"completed_at"`
	PreparedByName  string    `json:"prepared_by_name"`
	PhotosCount     int       `json:"photos_count"`
}

// FinalizeMaterialPrep menyelesaikan proses material preparation
// dengan validation untuk semua steps complete dan notification ke Unit Cetak
func (s *KhazwalService) FinalizeMaterialPrep(prepID uint64, photos []string, notes string, userID uint64) (*FinalizeResult, error) {
	// Start transaction untuk ensure data consistency
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Get material prep dengan lock untuk prevent concurrent updates
	var prep models.KhazwalMaterialPreparation
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("ProductionOrder").
		Preload("PreparedByUser").
		First(&prep, prepID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. Validasi: prep harus dalam status IN_PROGRESS
	if !prep.IsInProgress() {
		tx.Rollback()
		return nil, gorm.ErrInvalidData
	}

	// 3. Validasi: semua steps harus sudah complete
	if prep.PlatRetrievedAt == nil {
		tx.Rollback()
		return nil, gorm.ErrInvalidData
	}
	if prep.KertasBlankoActual == nil {
		tx.Rollback()
		return nil, gorm.ErrInvalidData
	}
	if prep.TintaActual == nil {
		tx.Rollback()
		return nil, gorm.ErrInvalidData
	}

	// 4. Calculate duration dari started_at ke now
	now := time.Now()
	durationMinutes := 0
	if prep.StartedAt != nil {
		duration := now.Sub(*prep.StartedAt)
		durationMinutes = int(duration.Minutes())
	}

	// 5. Convert photos ke JSON
	var photosJSON []byte
	var err error
	if len(photos) > 0 {
		photosJSON, err = json.Marshal(photos)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 6. Update material prep ke COMPLETED
	updates := map[string]interface{}{
		"status":           models.MaterialPrepCompleted,
		"completed_at":     now,
		"duration_minutes": durationMinutes,
		"updated_at":       now,
	}
	if notes != "" {
		updates["notes"] = notes
	}
	if len(photosJSON) > 0 {
		updates["material_photos"] = photosJSON
	}

	if err := tx.Model(&prep).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 7. Update Production Order stage & status ke CETAK
	if err := tx.Model(&models.ProductionOrder{}).
		Where("id = ?", prep.ProductionOrderID).
		Updates(map[string]interface{}{
			"current_stage":  models.StageCetak,
			"current_status": models.StatusReadyForCetak,
			"updated_at":     now,
		}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 8. Create POStageTracking record untuk audit trail
	tracking := models.POStageTracking{
		ProductionOrderID: prep.ProductionOrderID,
		Stage:             models.StageKhazwalMaterialPrep,
		Status:            models.StatusReadyForCetak,
		StartedAt:         prep.StartedAt,
		CompletedAt:       &now,
		HandledBy:         &userID,
		Notes:             "Material preparation selesai. Durasi: " + string(rune(durationMinutes)) + " menit",
	}
	if err := tx.Create(&tracking).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 9. Create notifications ke semua OPERATOR_CETAK users
	var cetakOperators []models.User
	if err := tx.Where("role = ? AND status = ?", models.RoleOperatorCetak, models.StatusActive).
		Find(&cetakOperators).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Buat notifikasi untuk setiap operator cetak
	poNumber := int64(0)
	obcNumber := ""
	if prep.ProductionOrder != nil {
		poNumber = prep.ProductionOrder.PONumber
		obcNumber = prep.ProductionOrder.OBCNumber
	}

	for _, operator := range cetakOperators {
		notification := models.Notification{
			UserID:  operator.ID,
			Title:   "Material Siap - PO #" + obcNumber,
			Message: "Material untuk PO #" + obcNumber + " telah siap. Silakan proses di unit cetak.",
			Type:    models.NotificationInfo,
			IsRead:  false,
		}
		if err := tx.Create(&notification).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 10. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Prepare result
	preparedByName := ""
	if prep.PreparedByUser != nil {
		preparedByName = prep.PreparedByUser.FullName
	}

	return &FinalizeResult{
		PrepID:          prep.ID,
		PONumber:        poNumber,
		OBCNumber:       obcNumber,
		DurationMinutes: durationMinutes,
		CompletedAt:     now,
		PreparedByName:  preparedByName,
		PhotosCount:     len(photos),
	}, nil
}

// QueueFilters merupakan struct untuk filter dan pagination queue
type QueueFilters struct {
	Search   string `form:"search"`
	Priority string `form:"priority"`
	Page     int    `form:"page"`
	PerPage  int    `form:"per_page"`
	SortBy   string `form:"sort_by"`
	SortDir  string `form:"sort_dir"`
}

// QueueResponse merupakan struct untuk paginated queue response
// yang mencakup items, total count, dan pagination metadata
type QueueResponse struct {
	Items      []models.ProductionOrder `json:"items"`
	Total      int                      `json:"total"`
	Page       int                      `json:"page"`
	PerPage    int                      `json:"per_page"`
	TotalPages int                      `json:"total_pages"`
}

// HistoryFilters merupakan struct untuk filter history query
type HistoryFilters struct {
	Search    string `form:"search"`
	StaffID   uint64 `form:"staff_id"`
	DateFrom  string `form:"date_from"`
	DateTo    string `form:"date_to"`
	Page      int    `form:"page"`
	PerPage   int    `form:"per_page"`
}

// HistoryItem merupakan struct untuk history item response
type HistoryItem struct {
	PrepID          uint64    `json:"prep_id"`
	POID            uint64    `json:"po_id"`
	PONumber        int64     `json:"po_number"`
	OBCNumber       string    `json:"obc_number"`
	ProductName     string    `json:"product_name"`
	Priority        string    `json:"priority"`
	Quantity        int       `json:"quantity"`
	CompletedAt     time.Time `json:"completed_at"`
	DurationMinutes int       `json:"duration_minutes"`
	PreparedByID    uint64    `json:"prepared_by_id"`
	PreparedByName  string    `json:"prepared_by_name"`
	PhotosCount     int       `json:"photos_count"`
}

// HistoryResponse merupakan struct untuk paginated history response
type HistoryResponse struct {
	Items      []HistoryItem `json:"items"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalPages int           `json:"total_pages"`
}

// GetMaterialPrepHistory mengambil riwayat material preparation yang sudah COMPLETED
// dengan filter by date range dan staff
func (s *KhazwalService) GetMaterialPrepHistory(filters HistoryFilters) (*HistoryResponse, error) {
	// Set default pagination values
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PerPage <= 0 || filters.PerPage > 100 {
		filters.PerPage = 20
	}

	// Query builder untuk material prep dengan status COMPLETED
	query := s.db.Model(&models.KhazwalMaterialPreparation{}).
		Where("status = ?", models.MaterialPrepCompleted).
		Preload("ProductionOrder").
		Preload("PreparedByUser")

	// Apply staff filter jika ada
	if filters.StaffID > 0 {
		query = query.Where("prepared_by = ?", filters.StaffID)
	}

	// Apply date range filter
	if filters.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", filters.DateFrom)
		if err == nil {
			query = query.Where("completed_at >= ?", dateFrom)
		}
	}
	if filters.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", filters.DateTo)
		if err == nil {
			// Add 1 day untuk include seluruh hari tersebut
			dateTo = dateTo.Add(24 * time.Hour)
			query = query.Where("completed_at < ?", dateTo)
		}
	}

	// Apply search filter
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Joins("LEFT JOIN production_orders ON production_orders.id = khazwal_material_preparations.production_order_id").
			Where("production_orders.po_number LIKE ? OR production_orders.obc_number LIKE ? OR production_orders.product_name LIKE ?",
				searchPattern, searchPattern, searchPattern)
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Execute query dengan pagination dan sorting
	var preps []models.KhazwalMaterialPreparation
	offset := (filters.Page - 1) * filters.PerPage
	if err := query.
		Order("completed_at DESC").
		Limit(filters.PerPage).
		Offset(offset).
		Find(&preps).Error; err != nil {
		return nil, err
	}

	// Transform ke HistoryItem
	items := make([]HistoryItem, 0, len(preps))
	for _, prep := range preps {
		item := HistoryItem{
			PrepID: prep.ID,
		}
		if prep.DurationMinutes != nil {
			item.DurationMinutes = *prep.DurationMinutes
		}

		if prep.CompletedAt != nil {
			item.CompletedAt = *prep.CompletedAt
		}
		if prep.PreparedBy != nil {
			item.PreparedByID = *prep.PreparedBy
		}
		if prep.PreparedByUser != nil {
			item.PreparedByName = prep.PreparedByUser.FullName
		}
		if prep.ProductionOrder != nil {
			item.POID = prep.ProductionOrder.ID
			item.PONumber = prep.ProductionOrder.PONumber
			item.OBCNumber = prep.ProductionOrder.OBCNumber
			item.ProductName = prep.ProductionOrder.ProductName
			item.Priority = string(prep.ProductionOrder.Priority)
			item.Quantity = prep.ProductionOrder.QuantityOrdered
		}

		// Count photos
		if prep.MaterialPhotos != nil {
			var photos []string
			if err := json.Unmarshal(prep.MaterialPhotos, &photos); err == nil {
				item.PhotosCount = len(photos)
			}
		}

		items = append(items, item)
	}

	// Calculate total pages
	totalPages := int(total) / filters.PerPage
	if int(total)%filters.PerPage > 0 {
		totalPages++
	}

	return &HistoryResponse{
		Items:      items,
		Total:      int(total),
		Page:       filters.Page,
		PerPage:    filters.PerPage,
		TotalPages: totalPages,
	}, nil
}

// MonitoringStats merupakan struct untuk supervisor monitoring dashboard
type MonitoringStats struct {
	TotalInQueue        int              `json:"total_in_queue"`
	TotalInProgress     int              `json:"total_in_progress"`
	TotalCompletedToday int              `json:"total_completed_today"`
	AverageDurationMins int              `json:"average_duration_mins"`
	StaffActive         []StaffActivity  `json:"staff_active"`
	RecentCompletions   []RecentComplete `json:"recent_completions"`
}

// StaffActivity merupakan struct untuk staff yang sedang aktif
type StaffActivity struct {
	UserID       uint64     `json:"user_id"`
	Name         string     `json:"name"`
	CurrentPOID  *uint64    `json:"current_po_id"`
	CurrentPO    string     `json:"current_po"`
	ProductName  string     `json:"product_name"`
	StartedAt    *time.Time `json:"started_at"`
	Status       string     `json:"status"`
	DurationMins int        `json:"duration_mins"`
}

// RecentComplete merupakan struct untuk recent completions
type RecentComplete struct {
	PrepID          uint64    `json:"prep_id"`
	PONumber        int64     `json:"po_number"`
	OBCNumber       string    `json:"obc_number"`
	CompletedAt     time.Time `json:"completed_at"`
	DurationMinutes int       `json:"duration_minutes"`
	PreparedByName  string    `json:"prepared_by_name"`
}

// GetMonitoringStats mengambil statistik untuk supervisor monitoring dashboard
func (s *KhazwalService) GetMonitoringStats() (*MonitoringStats, error) {
	stats := &MonitoringStats{}

	// 1. Count total in queue (WAITING_MATERIAL_PREP)
	var inQueue int64
	if err := s.db.Model(&models.ProductionOrder{}).
		Where("current_status = ?", models.StatusWaitingMaterialPrep).
		Count(&inQueue).Error; err != nil {
		return nil, err
	}
	stats.TotalInQueue = int(inQueue)

	// 2. Count total in progress (MATERIAL_PREP_IN_PROGRESS)
	var inProgress int64
	if err := s.db.Model(&models.ProductionOrder{}).
		Where("current_status = ?", models.StatusMaterialPrepInProgress).
		Count(&inProgress).Error; err != nil {
		return nil, err
	}
	stats.TotalInProgress = int(inProgress)

	// 3. Count total completed today
	today := time.Now().Truncate(24 * time.Hour)
	var completedToday int64
	if err := s.db.Model(&models.KhazwalMaterialPreparation{}).
		Where("status = ? AND completed_at >= ?", models.MaterialPrepCompleted, today).
		Count(&completedToday).Error; err != nil {
		return nil, err
	}
	stats.TotalCompletedToday = int(completedToday)

	// 4. Calculate average duration (last 30 days)
	var avgDuration struct {
		Avg float64
	}
	thirtyDaysAgo := time.Now().Add(-30 * 24 * time.Hour)
	s.db.Model(&models.KhazwalMaterialPreparation{}).
		Select("AVG(duration_minutes) as avg").
		Where("status = ? AND completed_at >= ?", models.MaterialPrepCompleted, thirtyDaysAgo).
		Scan(&avgDuration)
	stats.AverageDurationMins = int(avgDuration.Avg)

	// 5. Get active staff (those handling material prep in progress)
	var activePreps []models.KhazwalMaterialPreparation
	if err := s.db.
		Where("status = ?", models.MaterialPrepInProgress).
		Preload("ProductionOrder").
		Preload("PreparedByUser").
		Find(&activePreps).Error; err != nil {
		return nil, err
	}

	staffActive := make([]StaffActivity, 0)
	for _, prep := range activePreps {
		activity := StaffActivity{
			Status: "active",
		}

		if prep.PreparedBy != nil {
			activity.UserID = *prep.PreparedBy
		}
		if prep.PreparedByUser != nil {
			activity.Name = prep.PreparedByUser.FullName
		}
		if prep.ProductionOrder != nil {
			activity.CurrentPOID = &prep.ProductionOrder.ID
			activity.CurrentPO = prep.ProductionOrder.OBCNumber
			activity.ProductName = prep.ProductionOrder.ProductName
		}
		if prep.StartedAt != nil {
			activity.StartedAt = prep.StartedAt
			activity.DurationMins = int(time.Since(*prep.StartedAt).Minutes())
		}

		staffActive = append(staffActive, activity)
	}
	stats.StaffActive = staffActive

	// 6. Get recent completions (last 10)
	var recentPreps []models.KhazwalMaterialPreparation
	if err := s.db.
		Where("status = ?", models.MaterialPrepCompleted).
		Order("completed_at DESC").
		Limit(10).
		Preload("ProductionOrder").
		Preload("PreparedByUser").
		Find(&recentPreps).Error; err != nil {
		return nil, err
	}

	recentCompletions := make([]RecentComplete, 0)
	for _, prep := range recentPreps {
		recent := RecentComplete{
			PrepID: prep.ID,
		}
		if prep.DurationMinutes != nil {
			recent.DurationMinutes = *prep.DurationMinutes
		}

		if prep.CompletedAt != nil {
			recent.CompletedAt = *prep.CompletedAt
		}
		if prep.ProductionOrder != nil {
			recent.PONumber = prep.ProductionOrder.PONumber
			recent.OBCNumber = prep.ProductionOrder.OBCNumber
		}
		if prep.PreparedByUser != nil {
			recent.PreparedByName = prep.PreparedByUser.FullName
		}

		recentCompletions = append(recentCompletions, recent)
	}
	stats.RecentCompletions = recentCompletions

	return stats, nil
}
