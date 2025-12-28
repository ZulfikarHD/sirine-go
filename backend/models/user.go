package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole merupakan enum untuk roles dalam sistem
type UserRole string

const (
	RoleAdmin          UserRole = "ADMIN"
	RoleManager        UserRole = "MANAGER"
	RoleStaffKhazwal   UserRole = "STAFF_KHAZWAL"
	RoleOperatorCetak  UserRole = "OPERATOR_CETAK"
	RoleQCInspector    UserRole = "QC_INSPECTOR"
	RoleVerifikator    UserRole = "VERIFIKATOR"
	RoleStaffKhazkhir  UserRole = "STAFF_KHAZKHIR"
)

// Department merupakan enum untuk departemen
type Department string

const (
	DeptKhazwal     Department = "KHAZWAL"
	DeptCetak       Department = "CETAK"
	DeptVerifikasi  Department = "VERIFIKASI"
	DeptKhazkhir    Department = "KHAZKHIR"
)

// Shift merupakan enum untuk shift kerja
type Shift string

const (
	ShiftPagi  Shift = "PAGI"
	ShiftSiang Shift = "SIANG"
	ShiftMalam Shift = "MALAM"
)

// UserStatus merupakan enum untuk status user
type UserStatus string

const (
	StatusActive    UserStatus = "ACTIVE"
	StatusInactive  UserStatus = "INACTIVE"
	StatusSuspended UserStatus = "SUSPENDED"
)

// User merupakan model untuk entitas user dalam sistem
// yang mencakup authentication, authorization, dan profile information
type User struct {
	ID                  uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	NIP                 string         `gorm:"column:nip;type:varchar(5);uniqueIndex;not null" json:"nip" binding:"required,max=5,numeric"`
	FullName            string         `gorm:"type:varchar(255);not null" json:"full_name" binding:"required"`
	Email               string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" binding:"required,email"`
	Phone               string         `gorm:"type:varchar(20)" json:"phone"`
	PasswordHash        string         `gorm:"type:varchar(255);not null" json:"-"` // Hidden dari JSON response
	Role                UserRole       `gorm:"type:enum('ADMIN','MANAGER','STAFF_KHAZWAL','OPERATOR_CETAK','QC_INSPECTOR','VERIFIKATOR','STAFF_KHAZKHIR');not null" json:"role" binding:"required"`
	Department          Department     `gorm:"type:enum('KHAZWAL','CETAK','VERIFIKASI','KHAZKHIR');not null" json:"department" binding:"required"`
	Shift               Shift          `gorm:"type:enum('PAGI','SIANG','MALAM');default:'PAGI'" json:"shift"`
	ProfilePhotoURL     string         `gorm:"type:varchar(500)" json:"profile_photo_url"`
	TotalPoints         int            `gorm:"default:0" json:"total_points"`
	Level               string         `gorm:"type:varchar(20);default:'Bronze'" json:"level"`
	Status              UserStatus     `gorm:"type:enum('ACTIVE','INACTIVE','SUSPENDED');default:'ACTIVE'" json:"status"`
	MustChangePassword  bool           `gorm:"default:true" json:"must_change_password"`
	LastLoginAt         *time.Time     `gorm:"type:timestamp null" json:"last_login_at"`
	FailedLoginAttempts int            `gorm:"default:0" json:"-"` // Hidden dari JSON
	LockedUntil         *time.Time     `gorm:"type:timestamp null" json:"locked_until"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// TableName menentukan nama tabel di database
func (User) TableName() string {
	return "users"
}

// IsLocked memeriksa apakah akun user sedang terkunci
// berdasarkan LockedUntil timestamp
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// IsActive memeriksa apakah user dalam status aktif
func (u *User) IsActive() bool {
	return u.Status == StatusActive && !u.IsLocked()
}

// HasRole memeriksa apakah user memiliki role tertentu
func (u *User) HasRole(roles ...UserRole) bool {
	for _, role := range roles {
		if u.Role == role {
			return true
		}
	}
	return false
}

// IsAdmin memeriksa apakah user adalah admin atau manager
func (u *User) IsAdmin() bool {
	return u.HasRole(RoleAdmin, RoleManager)
}

// SafeUser merupakan struct untuk response user tanpa sensitive data
type SafeUser struct {
	ID              uint64     `json:"id"`
	NIP             string     `json:"nip"`
	FullName        string     `json:"full_name"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Role            UserRole   `json:"role"`
	Department      Department `json:"department"`
	Shift           Shift      `json:"shift"`
	ProfilePhotoURL string     `json:"profile_photo_url"`
	TotalPoints     int        `json:"total_points"`
	Level           string     `json:"level"`
	Status          UserStatus `json:"status"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	LockedUntil     *time.Time `json:"locked_until,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ToSafeUser mengkonversi User menjadi SafeUser untuk response API
func (u *User) ToSafeUser() SafeUser {
	return SafeUser{
		ID:              u.ID,
		NIP:             u.NIP,
		FullName:        u.FullName,
		Email:           u.Email,
		Phone:           u.Phone,
		Role:            u.Role,
		Department:      u.Department,
		Shift:           u.Shift,
		ProfilePhotoURL: u.ProfilePhotoURL,
		TotalPoints:     u.TotalPoints,
		Level:           u.Level,
		Status:          u.Status,
		LastLoginAt:     u.LastLoginAt,
		LockedUntil:     u.LockedUntil,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

// GetLevelFromPoints menghitung level berdasarkan total points
// dengan threshold: Bronze (0-99), Silver (100-499), Gold (500-999), Platinum (1000+)
func GetLevelFromPoints(points int) string {
	switch {
	case points >= 1000:
		return "Platinum"
	case points >= 500:
		return "Gold"
	case points >= 100:
		return "Silver"
	default:
		return "Bronze"
	}
}
