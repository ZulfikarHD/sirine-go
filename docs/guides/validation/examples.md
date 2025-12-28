# Validation Examples - Practical Usage

## Overview

Validation Examples merupakan dokumentasi praktis yang berisi real-world examples dari penggunaan server-side validation di Sirine Go application yang bertujuan untuk membantu developer memahami implementation patterns, yaitu: contoh request/response dengan berbagai validation scenarios, integration dengan frontend Vue, dan best practices untuk error handling.

Dokumen ini menunjukkan bagaimana validation bekerja dari request hingga response dengan contoh konkret.

---

## Example 1: Login Validation

### Request Structure
```go
type LoginRequest struct {
    NIP        string `json:"nip" binding:"required"`
    Password   string `json:"password" binding:"required"`
    RememberMe bool   `json:"remember_me"`
}
```

### Handler Implementation
```go
func (h *AuthHandler) Login(c *gin.Context) {
    var req services.LoginRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        // Return validation errors dengan pesan dalam Bahasa Indonesia
        c.JSON(http.StatusBadRequest, utils.NewValidationErrorResponse(err))
        return
    }
    
    // Process login...
}
```

### Test Cases

#### ✅ Valid Request
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "99999",
    "password": "Admin@123",
    "remember_me": true
  }'
```

**Response:**
```json
{
    "success": true,
    "message": "Login berhasil",
    "data": {
        "token": "eyJhbGc...",
        "refresh_token": "eyJhbGc...",
        "user": { ... }
    }
}
```

#### ❌ Invalid Request - Missing Fields
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "",
    "password": ""
  }'
```

**Response:**
```json
{
    "success": false,
    "message": "Data yang dikirim tidak valid",
    "errors": {
        "NIP": "NIP harus diisi",
        "Password": "Password harus diisi"
    }
}
```

---

## Example 2: Create User Validation (Sprint 2)

### Request Structure
```go
type CreateUserRequest struct {
    NIP        string `json:"nip" binding:"required,max=5"`
    FullName   string `json:"full_name" binding:"required,min=3,max=100"`
    Email      string `json:"email" binding:"required,email,max=255"`
    Phone      string `json:"phone" binding:"required,min=10,max=15"`
    Role       string `json:"role" binding:"required,oneof=ADMIN MANAGER STAFF_KHAZWAL OPERATOR_CETAK QC_INSPECTOR VERIFIKATOR STAFF_KHAZKHIR"`
    Department string `json:"department" binding:"required,oneof=KHAZWAL CETAK VERIFIKASI KHAZKHIR"`
    Shift      string `json:"shift" binding:"omitempty,oneof=PAGI SIANG MALAM"`
}
```

### Handler Implementation
```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req services.CreateUserRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, utils.NewValidationErrorResponse(err))
        return
    }
    
    // Process create user...
}
```

### Test Cases

#### ✅ Valid Request
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "nip": "12345",
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "081234567890",
    "role": "OPERATOR_CETAK",
    "department": "CETAK",
    "shift": "PAGI"
  }'
```

**Response:**
```json
{
    "success": true,
    "message": "User berhasil dibuat",
    "data": {
        "user": { ... },
        "temporary_password": "TempPass123!"
    }
}
```

#### ❌ Invalid Request - Invalid Email
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "nip": "12345",
    "full_name": "John Doe",
    "email": "invalid-email",
    "phone": "081234567890",
    "role": "OPERATOR_CETAK",
    "department": "CETAK"
  }'
```

**Response:**
```json
{
    "success": false,
    "message": "Data yang dikirim tidak valid",
    "errors": {
        "Email": "Email harus berupa alamat email yang valid"
    }
}
```

#### ❌ Invalid Request - Invalid Role
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "nip": "12345",
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "081234567890",
    "role": "INVALID_ROLE",
    "department": "CETAK"
  }'
```

**Response:**
```json
{
    "success": false,
    "message": "Data yang dikirim tidak valid",
    "errors": {
        "Role": "Role tidak valid"
    }
}
```

#### ❌ Invalid Request - Multiple Errors
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "nip": "123456",
    "full_name": "Jo",
    "email": "invalid",
    "phone": "08123",
    "role": "INVALID",
    "department": "INVALID"
  }'
```

**Response:**
```json
{
    "success": false,
    "message": "Data yang dikirim tidak valid",
    "errors": {
        "NIP": "NIP maksimal 5 karakter",
        "FullName": "Nama Lengkap minimal 3 karakter",
        "Email": "Email harus berupa alamat email yang valid",
        "Phone": "Nomor Telepon minimal 10 karakter",
        "Role": "Role tidak valid",
        "Department": "Departemen tidak valid"
    }
}
```

---

## Example 3: Change Password Validation (Sprint 3)

### Request Structure
```go
type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password" binding:"required"`
    NewPassword     string `json:"new_password" binding:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
```

### Handler Implementation
```go
func (h *PasswordHandler) ChangePassword(c *gin.Context) {
    var req services.ChangePasswordRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, utils.NewValidationErrorResponse(err))
        return
    }
    
    // Additional business logic validation
    if err := h.passwordService.ValidatePasswordPolicy(req.NewPassword); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }
    
    // Process change password...
}
```

### Test Cases

#### ✅ Valid Request
```bash
curl -X PUT http://localhost:8080/api/profile/password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "current_password": "OldPass123!",
    "new_password": "NewPass123!",
    "confirm_password": "NewPass123!"
  }'
```

**Response:**
```json
{
    "success": true,
    "message": "Password berhasil diubah"
}
```

#### ❌ Invalid Request - Password Mismatch
```bash
curl -X PUT http://localhost:8080/api/profile/password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "current_password": "OldPass123!",
    "new_password": "NewPass123!",
    "confirm_password": "DifferentPass123!"
  }'
```

**Response:**
```json
{
    "success": false,
    "message": "Data yang dikirim tidak valid",
    "errors": {
        "ConfirmPassword": "ConfirmPassword harus sama dengan NewPassword"
    }
}
```

#### ❌ Invalid Request - Password Too Short
```bash
curl -X PUT http://localhost:8080/api/profile/password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "current_password": "OldPass123!",
    "new_password": "short",
    "confirm_password": "short"
  }'
```

**Response:**
```json
{
    "success": false,
    "message": "Data yang dikirim tidak valid",
    "errors": {
        "NewPassword": "NewPassword minimal 8 karakter"
    }
}
```

---

## Example 4: Bulk Operations Validation

### Request Structure
```go
type BulkDeleteRequest struct {
    UserIDs []uint64 `json:"user_ids" binding:"required,min=1,max=100,dive,required,gt=0"`
}
```

### Handler Implementation
```go
func (h *UserHandler) BulkDelete(c *gin.Context) {
    var req services.BulkDeleteRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, utils.NewValidationErrorResponse(err))
        return
    }
    
    // Process bulk delete...
}
```

### Test Cases

#### ✅ Valid Request
```bash
curl -X POST http://localhost:8080/api/users/bulk-delete \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "user_ids": [1, 2, 3, 4, 5]
  }'
```

**Response:**
```json
{
    "success": true,
    "message": "5 users berhasil dihapus"
}
```

#### ❌ Invalid Request - Empty Array
```bash
curl -X POST http://localhost:8080/api/users/bulk-delete \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "user_ids": []
  }'
```

**Response:**
```json
{
    "success": false,
    "message": "Data yang dikirim tidak valid",
    "errors": {
        "UserIDs": "UserIDs minimal 1 karakter"
    }
}
```

#### ❌ Invalid Request - Invalid ID (0)
```bash
curl -X POST http://localhost:8080/api/users/bulk-delete \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "user_ids": [1, 2, 0, 4]
  }'
```

**Response:**
```json
{
    "success": false,
    "message": "Data yang dikirim tidak valid",
    "errors": {
        "UserIDs[2]": "UserIDs[2] harus lebih besar dari 0"
    }
}
```

---

## Example 5: Query Parameter Validation

### Request Structure
```go
type UserFilters struct {
    Role       string `form:"role" binding:"omitempty,oneof=ADMIN MANAGER STAFF_KHAZWAL OPERATOR_CETAK QC_INSPECTOR VERIFIKATOR STAFF_KHAZKHIR"`
    Department string `form:"department" binding:"omitempty,oneof=KHAZWAL CETAK VERIFIKASI KHAZKHIR"`
    Status     string `form:"status" binding:"omitempty,oneof=ACTIVE INACTIVE SUSPENDED"`
    Search     string `form:"search" binding:"omitempty,max=100"`
    Page       int    `form:"page" binding:"omitempty,min=1"`
    PerPage    int    `form:"per_page" binding:"omitempty,min=1,max=100"`
}
```

### Handler Implementation
```go
func (h *UserHandler) GetUsers(c *gin.Context) {
    var filters services.UserFilters
    
    if err := c.ShouldBindQuery(&filters); err != nil {
        c.JSON(http.StatusBadRequest, utils.NewValidationErrorResponse(err))
        return
    }
    
    // Process get users...
}
```

### Test Cases

#### ✅ Valid Request
```bash
curl -X GET "http://localhost:8080/api/users?role=ADMIN&department=KHAZWAL&page=1&per_page=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**
```json
{
    "success": true,
    "message": "Data users berhasil diambil",
    "data": {
        "users": [ ... ],
        "pagination": {
            "page": 1,
            "per_page": 20,
            "total": 100,
            "total_pages": 5
        }
    }
}
```

#### ❌ Invalid Request - Invalid Role
```bash
curl -X GET "http://localhost:8080/api/users?role=INVALID_ROLE" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**
```json
{
    "success": false,
    "message": "Data yang dikirim tidak valid",
    "errors": {
        "Role": "Role tidak valid"
    }
}
```

#### ❌ Invalid Request - Invalid Pagination
```bash
curl -X GET "http://localhost:8080/api/users?page=0&per_page=200" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**
```json
{
    "success": false,
    "message": "Data yang dikirim tidak valid",
    "errors": {
        "Page": "Page harus lebih besar atau sama dengan 1",
        "PerPage": "PerPage maksimal 100 karakter"
    }
}
```

---

## Frontend Integration

### Handling Validation Errors di Vue

```javascript
// composables/useApi.js
const handleError = (error) => {
  if (error.response?.status === 400 && error.response?.data?.errors) {
    // Validation errors
    const validationErrors = error.response.data.errors
    
    // Display field-specific errors
    Object.keys(validationErrors).forEach(field => {
      // Show error untuk each field
      showFieldError(field, validationErrors[field])
    })
  } else {
    // Generic error
    showToast(error.response?.data?.message || 'Terjadi kesalahan')
  }
}

// Usage in component
const handleSubmit = async () => {
  try {
    await api.post('/api/users', formData)
    showToast('User berhasil dibuat', 'success')
  } catch (error) {
    handleError(error)
  }
}
```

### Display Validation Errors
```vue
<template>
  <form @submit.prevent="handleSubmit">
    <div>
      <label>Email</label>
      <input v-model="form.email" type="email" />
      <p v-if="errors.Email" class="text-red-600">{{ errors.Email }}</p>
    </div>
    
    <div>
      <label>Password</label>
      <input v-model="form.password" type="password" />
      <p v-if="errors.Password" class="text-red-600">{{ errors.Password }}</p>
    </div>
    
    <button type="submit">Submit</button>
  </form>
</template>

<script setup>
import { ref } from 'vue'
import { useApi } from '@/composables/useApi'

const api = useApi()
const form = ref({ email: '', password: '' })
const errors = ref({})

const handleSubmit = async () => {
  errors.value = {}
  
  try {
    await api.post('/api/auth/login', form.value)
  } catch (error) {
    if (error.response?.data?.errors) {
      errors.value = error.response.data.errors
    }
  }
}
</script>
```

---

## Summary

### Key Points:

1. **Validation di Struct Level**: Menggunakan binding tags untuk declarative validation
2. **Automatic Validation**: Gin automatically validates saat `ShouldBindJSON()` atau `ShouldBindQuery()`
3. **User-Friendly Errors**: Translation helper untuk convert ke Bahasa Indonesia
4. **Field-Specific Errors**: Return map dengan field name sebagai key
5. **Frontend Integration**: Easy integration dengan Vue untuk display errors per field

### Benefits:

✅ **Type-Safe**: Compile-time checking  
✅ **Declarative**: Validation rules di struct tags  
✅ **Reusable**: Validation helper dapat digunakan di semua handlers  
✅ **User-Friendly**: Error messages dalam Bahasa Indonesia  
✅ **Consistent**: Consistent error response format  
✅ **Frontend-Ready**: Easy integration dengan frontend frameworks  

### Next Steps untuk Sprint 2:

1. Implement user handlers dengan validation
2. Test semua validation scenarios
3. Update frontend untuk handle validation errors
4. Add custom validators jika diperlukan (phone format, NIP format, etc)
