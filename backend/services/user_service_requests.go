package services

// CreateUserRequest merupakan request structure untuk create user
// dengan comprehensive validation rules
type CreateUserRequest struct {
	NIP        string `json:"nip" binding:"required,max=5"`
	FullName   string `json:"full_name" binding:"required,min=3,max=100"`
	Email      string `json:"email" binding:"required,email,max=255"`
	Phone      string `json:"phone" binding:"required,min=10,max=15"`
	Role       string `json:"role" binding:"required,oneof=ADMIN MANAGER STAFF_KHAZWAL OPERATOR_CETAK QC_INSPECTOR VERIFIKATOR STAFF_KHAZKHIR"`
	Department string `json:"department" binding:"required,oneof=KHAZWAL CETAK VERIFIKASI KHAZKHIR"`
	Shift      string `json:"shift" binding:"omitempty,oneof=PAGI SIANG MALAM"`
}

// UpdateUserRequest merupakan request structure untuk update user (Admin)
type UpdateUserRequest struct {
	FullName   string `json:"full_name" binding:"omitempty,min=3,max=100"`
	Email      string `json:"email" binding:"omitempty,email,max=255"`
	Phone      string `json:"phone" binding:"omitempty,min=10,max=15"`
	Role       string `json:"role" binding:"omitempty,oneof=ADMIN MANAGER STAFF_KHAZWAL OPERATOR_CETAK QC_INSPECTOR VERIFIKATOR STAFF_KHAZKHIR"`
	Department string `json:"department" binding:"omitempty,oneof=KHAZWAL CETAK VERIFIKASI KHAZKHIR"`
	Shift      string `json:"shift" binding:"omitempty,oneof=PAGI SIANG MALAM"`
	Status     string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE SUSPENDED"`
}

// UpdateProfileRequest merupakan request structure untuk update own profile
// User hanya bisa update field tertentu (tidak bisa update role/department)
type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Phone    string `json:"phone" binding:"required,min=10,max=15"`
}

// ChangePasswordRequest merupakan request structure untuk change password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// ForgotPasswordRequest merupakan request structure untuk forgot password
type ForgotPasswordRequest struct {
	NIP   string `json:"nip" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest merupakan request structure untuk reset password
type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// UserFilters merupakan structure untuk filter users list
type UserFilters struct {
	Role       string `form:"role" binding:"omitempty,oneof=ADMIN MANAGER STAFF_KHAZWAL OPERATOR_CETAK QC_INSPECTOR VERIFIKATOR STAFF_KHAZKHIR"`
	Department string `form:"department" binding:"omitempty,oneof=KHAZWAL CETAK VERIFIKASI KHAZKHIR"`
	Status     string `form:"status" binding:"omitempty,oneof=ACTIVE INACTIVE SUSPENDED"`
	Search     string `form:"search" binding:"omitempty,max=100"`
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PerPage    int    `form:"per_page" binding:"omitempty,min=1,max=100"`
}

// BulkDeleteRequest merupakan request structure untuk bulk delete users
type BulkDeleteRequest struct {
	UserIDs []uint64 `json:"user_ids" binding:"required,min=1,max=100,dive,required,gt=0"`
}

// BulkUpdateStatusRequest merupakan request structure untuk bulk update status
type BulkUpdateStatusRequest struct {
	UserIDs []uint64 `json:"user_ids" binding:"required,min=1,max=100,dive,required,gt=0"`
	Status  string   `json:"status" binding:"required,oneof=ACTIVE INACTIVE SUSPENDED"`
}
