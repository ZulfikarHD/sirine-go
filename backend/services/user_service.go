package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"sirine-go/backend/models"
	"strings"

	"gorm.io/gorm"
)

// UserService merupakan service untuk user management operations
// yang mencakup CRUD, search, dan filter functionality
type UserService struct {
	db              *gorm.DB
	passwordService *PasswordService
}

// NewUserService membuat instance baru dari UserService
func NewUserService(db *gorm.DB, passwordService *PasswordService) *UserService {
	return &UserService{
		db:              db,
		passwordService: passwordService,
	}
}

// UserListResponse merupakan response structure untuk list users dengan pagination
type UserListResponse struct {
	Users      []models.SafeUser `json:"users"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	TotalPages int               `json:"total_pages"`
}

// GetAllUsers mengambil list users dengan filters dan pagination
func (s *UserService) GetAllUsers(filters UserFilters) (*UserListResponse, error) {
	var users []models.User
	var total int64

	// Set default pagination
	page := filters.Page
	if page < 1 {
		page = 1
	}
	perPage := filters.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	// Build query dengan filters
	query := s.db.Model(&models.User{})

	// Filter by role
	if filters.Role != "" {
		query = query.Where("role = ?", filters.Role)
	}

	// Filter by department
	if filters.Department != "" {
		query = query.Where("department = ?", filters.Department)
	}

	// Filter by status
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Search by NIP atau Full Name
	if filters.Search != "" {
		searchTerm := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where("LOWER(nip) LIKE ? OR LOWER(full_name) LIKE ?", searchTerm, searchTerm)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("gagal menghitung total users: %w", err)
	}

	// Apply pagination
	offset := (page - 1) * perPage
	query = query.Offset(offset).Limit(perPage)

	// Order by created_at desc (newest first)
	query = query.Order("created_at DESC")

	// Execute query
	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data users: %w", err)
	}

	// Convert to SafeUser
	safeUsers := make([]models.SafeUser, len(users))
	for i, user := range users {
		safeUsers[i] = user.ToSafeUser()
	}

	// Calculate total pages
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &UserListResponse{
		Users:      safeUsers,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetUserByID mengambil user berdasarkan ID
func (s *UserService) GetUserByID(id uint64) (*models.SafeUser, error) {
	var user models.User

	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user dengan ID %d tidak ditemukan", id)
		}
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	safeUser := user.ToSafeUser()
	return &safeUser, nil
}

// GetUserByNIP mengambil user berdasarkan NIP
func (s *UserService) GetUserByNIP(nip string) (*models.User, error) {
	var user models.User

	if err := s.db.Where("nip = ?", nip).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user dengan NIP %s tidak ditemukan", nip)
		}
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	return &user, nil
}

// CreateUserResponse merupakan response setelah create user
type CreateUserResponse struct {
	User            models.SafeUser `json:"user"`
	GeneratedPassword string          `json:"generated_password"`
}

// generateRandomPassword membuat random password yang memenuhi policy
func (s *UserService) generateRandomPassword() (string, error) {
	// Generate 12 character random password dengan uppercase, lowercase, number, dan special char
	const (
		lowerChars   = "abcdefghijklmnopqrstuvwxyz"
		upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digitChars   = "0123456789"
		specialChars = "@#$%&*"
		allChars     = lowerChars + upperChars + digitChars + specialChars
		length       = 12
	)

	// Ensure at least 1 dari setiap category
	password := make([]byte, length)

	// 1 uppercase
	password[0] = upperChars[s.randomInt(len(upperChars))]
	// 1 lowercase
	password[1] = lowerChars[s.randomInt(len(lowerChars))]
	// 1 digit
	password[2] = digitChars[s.randomInt(len(digitChars))]
	// 1 special
	password[3] = specialChars[s.randomInt(len(specialChars))]

	// Fill rest randomly
	for i := 4; i < length; i++ {
		password[i] = allChars[s.randomInt(len(allChars))]
	}

	// Shuffle untuk avoid predictable pattern
	for i := length - 1; i > 0; i-- {
		j := s.randomInt(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// randomInt menghasilkan random integer dari 0 hingga max-1
func (s *UserService) randomInt(max int) int {
	b := make([]byte, 1)
	rand.Read(b)
	return int(b[0]) % max
}

// CreateUser membuat user baru dengan auto-generated password
func (s *UserService) CreateUser(req CreateUserRequest) (*CreateUserResponse, error) {
	// Validate NIP uniqueness
	var existingUser models.User
	if err := s.db.Where("nip = ?", req.NIP).First(&existingUser).Error; err == nil {
		return nil, errors.New("NIP sudah terdaftar dalam sistem")
	}

	// Validate Email uniqueness
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email sudah terdaftar dalam sistem")
	}

	// Generate random password
	randomPassword, err := s.generateRandomPassword()
	if err != nil {
		return nil, fmt.Errorf("gagal generate password: %w", err)
	}

	// Hash password
	hashedPassword, err := s.passwordService.HashPassword(randomPassword)
	if err != nil {
		return nil, fmt.Errorf("gagal hash password: %w", err)
	}

	// Set default shift jika tidak ada
	shift := req.Shift
	if shift == "" {
		shift = string(models.ShiftPagi)
	}

	// Create user
	user := models.User{
		NIP:                req.NIP,
		FullName:           req.FullName,
		Email:              req.Email,
		Phone:              req.Phone,
		PasswordHash:       hashedPassword,
		Role:               models.UserRole(req.Role),
		Department:         models.Department(req.Department),
		Shift:              models.Shift(shift),
		Status:             models.StatusActive,
		MustChangePassword: true, // User harus change password saat first login
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("gagal membuat user baru: %w", err)
	}

	return &CreateUserResponse{
		User:              user.ToSafeUser(),
		GeneratedPassword: randomPassword,
	}, nil
}

// UpdateUser mengupdate user data (Admin only)
func (s *UserService) UpdateUser(id uint64, req UpdateUserRequest) (*models.SafeUser, error) {
	var user models.User

	// Find user
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user dengan ID %d tidak ditemukan", id)
		}
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	// Update fields yang ada
	updates := make(map[string]interface{})

	if req.FullName != "" {
		updates["full_name"] = req.FullName
	}
	if req.Email != "" {
		// Validate email uniqueness (exclude current user)
		var existingUser models.User
		if err := s.db.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("email sudah digunakan oleh user lain")
		}
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Department != "" {
		updates["department"] = req.Department
	}
	if req.Shift != "" {
		updates["shift"] = req.Shift
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	// Perform update
	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("gagal update user: %w", err)
	}

	// Reload user untuk get updated data
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data user setelah update: %w", err)
	}

	safeUser := user.ToSafeUser()
	return &safeUser, nil
}

// DeleteUser melakukan soft delete pada user
func (s *UserService) DeleteUser(id uint64) error {
	var user models.User

	// Find user
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user dengan ID %d tidak ditemukan", id)
		}
		return fmt.Errorf("gagal mengambil data user: %w", err)
	}

	// Prevent deleting self (akan di-check di handler level juga)
	// Soft delete
	if err := s.db.Delete(&user).Error; err != nil {
		return fmt.Errorf("gagal menghapus user: %w", err)
	}

	return nil
}

// SearchUsers melakukan search users berdasarkan query
func (s *UserService) SearchUsers(query string) ([]models.SafeUser, error) {
	var users []models.User

	searchTerm := "%" + strings.ToLower(query) + "%"
	
	if err := s.db.Where("LOWER(nip) LIKE ? OR LOWER(full_name) LIKE ? OR LOWER(email) LIKE ?", 
		searchTerm, searchTerm, searchTerm).
		Limit(10). // Limit untuk search results
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("gagal search users: %w", err)
	}

	// Convert to SafeUser
	safeUsers := make([]models.SafeUser, len(users))
	for i, user := range users {
		safeUsers[i] = user.ToSafeUser()
	}

	return safeUsers, nil
}

// UpdateProfile mengupdate profile user sendiri (self-service)
// Hanya bisa update: full_name, email, phone
func (s *UserService) UpdateProfile(userID uint64, req UpdateProfileRequest) (*models.SafeUser, error) {
	var user models.User

	// Find user
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	// Validate email uniqueness (exclude current user)
	if req.Email != user.Email {
		var existingUser models.User
		if err := s.db.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			return nil, errors.New("email sudah digunakan oleh user lain")
		}
	}

	// Update allowed fields only
	user.FullName = req.FullName
	user.Email = req.Email
	user.Phone = req.Phone

	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("gagal update profile: %w", err)
	}

	safeUser := user.ToSafeUser()
	return &safeUser, nil
}

// GetUserByEmail mengambil user berdasarkan email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user dengan email %s tidak ditemukan", email)
		}
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	return &user, nil
}

// BulkDeleteUsers melakukan soft delete pada multiple users
func (s *UserService) BulkDeleteUsers(userIDs []uint64, currentUserID uint64) (int, error) {
	// Filter out current user dari list
	var filteredIDs []uint64
	for _, id := range userIDs {
		if id != currentUserID {
			filteredIDs = append(filteredIDs, id)
		}
	}

	if len(filteredIDs) == 0 {
		return 0, errors.New("tidak ada user yang dapat dihapus")
	}

	// Perform bulk soft delete
	result := s.db.Where("id IN ?", filteredIDs).Delete(&models.User{})
	if result.Error != nil {
		return 0, fmt.Errorf("gagal melakukan bulk delete: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}

// BulkUpdateStatus mengupdate status multiple users
func (s *UserService) BulkUpdateStatus(userIDs []uint64, status string, currentUserID uint64) (int, error) {
	// Filter out current user dari list
	var filteredIDs []uint64
	for _, id := range userIDs {
		if id != currentUserID {
			filteredIDs = append(filteredIDs, id)
		}
	}

	if len(filteredIDs) == 0 {
		return 0, errors.New("tidak ada user yang dapat diupdate")
	}

	// Perform bulk update
	result := s.db.Model(&models.User{}).
		Where("id IN ?", filteredIDs).
		Update("status", status)

	if result.Error != nil {
		return 0, fmt.Errorf("gagal melakukan bulk update status: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}
