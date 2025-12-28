package services

import (
	"errors"
	"sirine-go/backend/models"
	"time"

	"gorm.io/gorm"
)

// AchievementService merupakan service untuk handling achievement operations
// yang mencakup fetching achievements, awarding achievements, dan checking criteria
type AchievementService struct {
	db                   *gorm.DB
	notificationService  *NotificationService
}

// NewAchievementService membuat instance baru dari AchievementService
// dengan dependency injection untuk database dan notification service
func NewAchievementService(db *gorm.DB, notificationService *NotificationService) *AchievementService {
	return &AchievementService{
		db:                  db,
		notificationService: notificationService,
	}
}

// GetAllAchievements mengambil semua achievements yang aktif
// untuk ditampilkan ke user dengan informasi unlock status
func (s *AchievementService) GetAllAchievements() ([]models.Achievement, error) {
	var achievements []models.Achievement
	
	if err := s.db.Where("is_active = ?", true).
		Order("category, points").
		Find(&achievements).Error; err != nil {
		return nil, err
	}
	
	return achievements, nil
}

// GetUserAchievements mengambil achievements user dengan status unlock
// untuk menampilkan progress dan achievements yang telah dicapai
func (s *AchievementService) GetUserAchievements(userID uint64) ([]models.UserAchievementResponse, error) {
	// Get all active achievements
	var allAchievements []models.Achievement
	if err := s.db.Where("is_active = ?", true).
		Order("category, points").
		Find(&allAchievements).Error; err != nil {
		return nil, err
	}

	// Get user's unlocked achievements
	var unlockedAchievements []models.UserAchievement
	if err := s.db.Where("user_id = ?", userID).
		Find(&unlockedAchievements).Error; err != nil {
		return nil, err
	}

	// Create map untuk quick lookup
	unlockedMap := make(map[uint64]*time.Time)
	for _, ua := range unlockedAchievements {
		unlockedMap[ua.AchievementID] = &ua.UnlockedAt
	}

	// Build response dengan merge data
	var response []models.UserAchievementResponse
	for _, achievement := range allAchievements {
		unlockedAt, isUnlocked := unlockedMap[achievement.ID]
		
		response = append(response, models.UserAchievementResponse{
			ID:          achievement.ID,
			Code:        achievement.Code,
			Name:        achievement.Name,
			Description: achievement.Description,
			Icon:        achievement.Icon,
			Points:      achievement.Points,
			Category:    string(achievement.Category),
			IsUnlocked:  isUnlocked,
			UnlockedAt:  unlockedAt,
		})
	}

	return response, nil
}

// AwardAchievement memberikan achievement ke user dan update total points
// serta mengirim notification achievement unlock
func (s *AchievementService) AwardAchievement(userID uint64, achievementCode string) error {
	// Get achievement by code
	var achievement models.Achievement
	if err := s.db.Where("code = ? AND is_active = ?", achievementCode, true).
		First(&achievement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("achievement tidak ditemukan")
		}
		return err
	}

	// Check jika sudah unlocked
	var existing models.UserAchievement
	err := s.db.Where("user_id = ? AND achievement_id = ?", userID, achievement.ID).
		First(&existing).Error
	
	if err == nil {
		// Sudah unlocked
		return nil
	}
	
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Begin transaction untuk award achievement dan update points
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user achievement record
	userAchievement := models.UserAchievement{
		UserID:        userID,
		AchievementID: achievement.ID,
		UnlockedAt:    time.Now(),
	}
	
	if err := tx.Create(&userAchievement).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update user total points dan level
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	newPoints := user.TotalPoints + achievement.Points
	newLevel := models.GetLevelFromPoints(newPoints)

	if err := tx.Model(&user).Updates(map[string]interface{}{
		"total_points": newPoints,
		"level":        newLevel,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// Send notification (async, tidak perlu rollback jika gagal)
	go s.notificationService.CreateNotification(
		userID,
		"Achievement Unlocked! 🎉",
		achievement.Name+" - "+achievement.Description,
		"ACHIEVEMENT",
	)

	return nil
}

// CheckAndAwardFirstLogin memeriksa dan memberikan achievement first login
// yang dipanggil saat user berhasil login
func (s *AchievementService) CheckAndAwardFirstLogin(userID uint64) error {
	// Check jika ini adalah login pertama dengan melihat last_login_at
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Jika last_login_at nil atau ini login pertama, award achievement
	if user.LastLoginAt == nil {
		return s.AwardAchievement(userID, "FIRST_LOGIN")
	}

	return nil
}

// CheckAndAwardProfileComplete memeriksa dan memberikan achievement profile complete
// yang dipanggil saat user melengkapi profile dengan foto
func (s *AchievementService) CheckAndAwardProfileComplete(userID uint64) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Check jika profile lengkap (ada foto, email, phone)
	if user.ProfilePhotoURL != "" && user.Email != "" && user.Phone != "" {
		return s.AwardAchievement(userID, "PROFILE_COMPLETE")
	}

	return nil
}

// CheckAndAwardLoginStreak memeriksa dan memberikan achievement login streak
// berdasarkan consecutive login days (untuk future implementation dengan tracking)
func (s *AchievementService) CheckAndAwardLoginStreak(userID uint64, streakDays int) error {
	switch streakDays {
	case 7:
		return s.AwardAchievement(userID, "WEEK_STREAK")
	case 30:
		return s.AwardAchievement(userID, "MONTH_STREAK")
	}
	return nil
}

// CheckAndAwardTimeBasedLogin memeriksa dan memberikan achievement berdasarkan jam login
// untuk Early Bird (sebelum 07:00) dan Night Owl (setelah 20:00)
func (s *AchievementService) CheckAndAwardTimeBasedLogin(userID uint64) error {
	now := time.Now()
	hour := now.Hour()

	// Early bird: sebelum 07:00
	if hour < 7 {
		// Count early logins (future: implement counter tracking)
		// Untuk simplicity, award langsung untuk demo
		return s.AwardAchievement(userID, "EARLY_BIRD")
	}

	// Night owl: setelah 20:00
	if hour >= 20 {
		return s.AwardAchievement(userID, "NIGHT_OWL")
	}

	return nil
}

// GetUserStats mengambil statistik user untuk gamification dashboard
// yang mencakup total points, level, achievements unlocked, dan progress
func (s *AchievementService) GetUserStats(userID uint64) (map[string]interface{}, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Count unlocked achievements
	var unlockedCount int64
	s.db.Model(&models.UserAchievement{}).
		Where("user_id = ?", userID).
		Count(&unlockedCount)

	// Count total achievements
	var totalCount int64
	s.db.Model(&models.Achievement{}).
		Where("is_active = ?", true).
		Count(&totalCount)

	// Calculate percentage
	percentage := 0
	if totalCount > 0 {
		percentage = int((float64(unlockedCount) / float64(totalCount)) * 100)
	}

	// Get points to next level
	nextLevel := ""
	pointsToNext := 0
	switch user.Level {
	case "Bronze":
		nextLevel = "Silver"
		pointsToNext = 100 - user.TotalPoints
	case "Silver":
		nextLevel = "Gold"
		pointsToNext = 500 - user.TotalPoints
	case "Gold":
		nextLevel = "Platinum"
		pointsToNext = 1000 - user.TotalPoints
	case "Platinum":
		nextLevel = "Max Level"
		pointsToNext = 0
	}

	return map[string]interface{}{
		"total_points":          user.TotalPoints,
		"level":                 user.Level,
		"achievements_unlocked": unlockedCount,
		"total_achievements":    totalCount,
		"percentage":            percentage,
		"next_level":            nextLevel,
		"points_to_next":        pointsToNext,
	}, nil
}
