package services

import (
	"sirine-go/backend/models"
	"time"

	"gorm.io/gorm"
)

// ActivityLogService merupakan service untuk audit logging
// yang bertujuan untuk tracking critical actions dalam sistem
type ActivityLogService struct {
	db *gorm.DB
}

// NewActivityLogService membuat instance baru ActivityLogService
func NewActivityLogService(db *gorm.DB) *ActivityLogService {
	return &ActivityLogService{db: db}
}

// ActivityLogFilters merupakan struct untuk filtering activity logs
// yang mencakup user, action, entity type, dan date range
type ActivityLogFilters struct {
	UserID     *uint64               `form:"user_id"`
	Action     *models.ActivityAction `form:"action"`
	EntityType *string               `form:"entity_type"`
	StartDate  *time.Time            `form:"start_date"`
	EndDate    *time.Time            `form:"end_date"`
	Search     string                `form:"search"` // Search by entity_id or user name
}

// Pagination merupakan struct untuk pagination params
type Pagination struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}

// ActivityLogResponse merupakan response dengan pagination info
type ActivityLogResponse struct {
	Data       []models.ActivityLog `json:"data"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// GetActivityLogs mengambil activity logs dengan filters dan pagination
// untuk admin audit logs viewer
func (s *ActivityLogService) GetActivityLogs(filters ActivityLogFilters, pagination Pagination) (*ActivityLogResponse, error) {
	var logs []models.ActivityLog
	var total int64

	// Default pagination
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 20
	}

	query := s.db.Model(&models.ActivityLog{}).Preload("User")

	// Apply filters
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Action != nil {
		query = query.Where("action = ?", *filters.Action)
	}

	if filters.EntityType != nil {
		query = query.Where("entity_type = ?", *filters.EntityType)
	}

	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", *filters.StartDate)
	}

	if filters.EndDate != nil {
		// Add 1 day untuk inclusive end date
		endDate := filters.EndDate.Add(24 * time.Hour)
		query = query.Where("created_at < ?", endDate)
	}

	if filters.Search != "" {
		// Search by entity_id atau join dengan user name
		query = query.Where("CAST(entity_id AS CHAR) LIKE ? OR entity_type LIKE ?",
			"%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// Get paginated data
	offset := (pagination.Page - 1) * pagination.PageSize
	err = query.Order("created_at DESC").
		Limit(pagination.PageSize).
		Offset(offset).
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pagination.PageSize
	if int(total)%pagination.PageSize > 0 {
		totalPages++
	}

	return &ActivityLogResponse{
		Data:       logs,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetUserActivity mengambil activity logs untuk specific user
// untuk display di user detail page
func (s *ActivityLogService) GetUserActivity(userID uint64, limit int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog

	if limit < 1 {
		limit = 50
	}

	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}

// CreateActivityLog mencatat activity log baru
// dipanggil oleh activity_logger middleware
func (s *ActivityLogService) CreateActivityLog(log *models.ActivityLog) error {
	return s.db.Create(log).Error
}

// GetActivityLogByID mengambil single activity log dengan detail
// untuk view changes detail
func (s *ActivityLogService) GetActivityLogByID(id uint64) (*models.ActivityLog, error) {
	var log models.ActivityLog
	err := s.db.Preload("User").First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetActivityStats mengambil statistik activity untuk dashboard
// dengan breakdown by action type
func (s *ActivityLogService) GetActivityStats(startDate, endDate time.Time) (map[string]int64, error) {
	type ActionCount struct {
		Action models.ActivityAction `json:"action"`
		Count  int64                `json:"count"`
	}

	var results []ActionCount
	err := s.db.Model(&models.ActivityLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("action").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[string(result.Action)] = result.Count
	}

	return stats, nil
}

// DeleteOldLogs menghapus activity logs lebih dari X hari
// untuk maintenance dan cleanup old audit data
func (s *ActivityLogService) DeleteOldLogs(olderThanDays int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -olderThanDays)
	result := s.db.Where("created_at < ?", cutoffDate).Delete(&models.ActivityLog{})
	return result.RowsAffected, result.Error
}
