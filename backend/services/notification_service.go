package services

import (
	"sirine-go/backend/models"
	"time"

	"gorm.io/gorm"
)

// NotificationService merupakan service untuk management notifikasi
// yang bertujuan untuk handling in-app notifications kepada users
type NotificationService struct {
	db *gorm.DB
}

// NewNotificationService membuat instance baru NotificationService
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// GetUserNotifications mengambil semua notifikasi user
// dengan optional filter untuk unread only
func (s *NotificationService) GetUserNotifications(userID uint64, unreadOnly bool) ([]models.Notification, error) {
	var notifications []models.Notification
	query := s.db.Where("user_id = ?", userID)

	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	err := query.Order("created_at DESC").Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

// GetUnreadCount mengambil jumlah notifikasi yang belum dibaca
// untuk badge display di notification bell
func (s *NotificationService) GetUnreadCount(userID uint64) (int64, error) {
	var count int64
	err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

// MarkAsRead menandai satu notifikasi sebagai sudah dibaca
// dengan update timestamp read_at
func (s *NotificationService) MarkAsRead(notificationID uint64, userID uint64) error {
	now := time.Now()
	result := s.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// MarkAllAsRead menandai semua notifikasi user sebagai sudah dibaca
// untuk bulk action dari notification center
func (s *NotificationService) MarkAllAsRead(userID uint64) error {
	now := time.Now()
	return s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error
}

// CreateNotification membuat notifikasi baru untuk user
// dengan type yang ditentukan (INFO, SUCCESS, WARNING, ERROR, ACHIEVEMENT)
func (s *NotificationService) CreateNotification(userID uint64, title, message string, notifType models.NotificationType) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    notifType,
		IsRead:  false,
	}

	err := s.db.Create(notification).Error
	if err != nil {
		return nil, err
	}

	return notification, nil
}

// GetRecentNotifications mengambil N notifikasi terbaru
// untuk quick preview di notification bell dropdown
func (s *NotificationService) GetRecentNotifications(userID uint64, limit int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&notifications).Error

	if err != nil {
		return nil, err
	}

	return notifications, nil
}

// DeleteNotification menghapus notifikasi (soft delete)
// untuk user yang ingin cleanup notification list
func (s *NotificationService) DeleteNotification(notificationID uint64, userID uint64) error {
	result := s.db.Where("id = ? AND user_id = ?", notificationID, userID).
		Delete(&models.Notification{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
