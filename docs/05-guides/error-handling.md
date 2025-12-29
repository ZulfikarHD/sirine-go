# 🚨 Error Handling Guide - Sirine Go

Panduan comprehensive untuk error handling dalam Sirine Go App, mencakup response format, status codes, dan best practices untuk menangani errors dengan konsisten.

**Last Updated:** 28 Desember 2025

---

## 📋 Overview

Error handling dalam Sirine Go mengikuti standar RESTful API dengan response format yang konsisten, yaitu:
- Consistent JSON structure untuk semua errors
- Descriptive error messages dalam Bahasa Indonesia
- Appropriate HTTP status codes
- Additional error details untuk debugging

---

## 🔍 Response Format

### Success Response

**Structure:**
```json
{
  "success": true,
  "data": { ... },
  "message": "Pesan sukses dalam Bahasa Indonesia"
}
```

**Example:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "full_name": "Super Admin",
      "email": "admin@sirine.local"
    }
  },
  "message": "User berhasil dibuat"
}
```

### Error Response

**Structure:**
```json
{
  "success": false,
  "error": "Pesan error utama dalam Bahasa Indonesia",
  "details": "Additional error details (optional)"
}
```

**Example:**
```json
{
  "success": false,
  "error": "NIP atau password salah"
}
```

### Validation Error Response

**Structure:**
```json
{
  "success": false,
  "error": "Validation error",
  "details": {
    "field_name": "Error message untuk field tersebut"
  }
}
```

**Example:**
```json
{
  "success": false,
  "error": "Validasi gagal",
  "details": {
    "email": "Format email tidak valid",
    "phone": "Nomor telepon harus 10-13 digit"
  }
}
```

---

## 📊 HTTP Status Codes

### Success Codes (2xx)

| Status Code | Meaning | Usage |
|-------------|---------|-------|
| **200 OK** | Success | GET, PUT, DELETE requests berhasil |
| **201 Created** | Resource created | POST request berhasil membuat resource baru |
| **204 No Content** | Success without body | DELETE request berhasil tanpa response body |

### Client Error Codes (4xx)

| Status Code | Meaning | Usage | Example |
|-------------|---------|-------|---------|
| **400 Bad Request** | Invalid request | Validation error, invalid parameters | "Format email tidak valid" |
| **401 Unauthorized** | Authentication required | Token missing/invalid/expired | "Token tidak valid atau sudah expired" |
| **403 Forbidden** | Permission denied | Insufficient permissions | "Anda tidak memiliki akses ke resource ini" |
| **404 Not Found** | Resource not found | Resource ID tidak ditemukan | "User dengan ID tersebut tidak ditemukan" |
| **409 Conflict** | Resource conflict | Duplicate data, constraint violation | "NIP sudah terdaftar" |
| **422 Unprocessable Entity** | Validation error | Business logic validation failed | "Password tidak memenuhi kriteria keamanan" |
| **429 Too Many Requests** | Rate limit exceeded | Terlalu banyak requests | "Terlalu banyak percobaan. Coba lagi dalam 15 menit" |

### Server Error Codes (5xx)

| Status Code | Meaning | Usage | Example |
|-------------|---------|-------|---------|
| **500 Internal Server Error** | Server error | Unexpected server error | "Terjadi kesalahan pada server" |
| **503 Service Unavailable** | Service down | Database connection failed, maintenance | "Layanan sedang dalam pemeliharaan" |

---

## 🎯 Common Error Scenarios

### Authentication Errors

#### 1. Invalid Credentials
```json
{
  "success": false,
  "error": "NIP atau password salah"
}
```
**Status Code:** `401 Unauthorized`

#### 2. Token Expired
```json
{
  "success": false,
  "error": "Token sudah expired. Silakan login kembali"
}
```
**Status Code:** `401 Unauthorized`

#### 3. Token Invalid
```json
{
  "success": false,
  "error": "Token tidak valid"
}
```
**Status Code:** `401 Unauthorized`

#### 4. Account Locked
```json
{
  "success": false,
  "error": "Akun terkunci. Coba lagi dalam 15 menit"
}
```
**Status Code:** `403 Forbidden`

### Validation Errors

#### 1. Required Field Missing
```json
{
  "success": false,
  "error": "Validasi gagal",
  "details": {
    "nip": "NIP wajib diisi",
    "password": "Password wajib diisi"
  }
}
```
**Status Code:** `400 Bad Request`

#### 2. Invalid Format
```json
{
  "success": false,
  "error": "Validasi gagal",
  "details": {
    "email": "Format email tidak valid",
    "phone": "Nomor telepon harus 10-13 digit"
  }
}
```
**Status Code:** `400 Bad Request`

#### 3. Password Strength
```json
{
  "success": false,
  "error": "Password tidak memenuhi kriteria keamanan",
  "details": "Password harus minimal 8 karakter, mengandung huruf besar, huruf kecil, angka, dan simbol"
}
```
**Status Code:** `422 Unprocessable Entity`

### Resource Errors

#### 1. Not Found
```json
{
  "success": false,
  "error": "User dengan ID tersebut tidak ditemukan"
}
```
**Status Code:** `404 Not Found`

#### 2. Duplicate Resource
```json
{
  "success": false,
  "error": "NIP sudah terdaftar di sistem"
}
```
**Status Code:** `409 Conflict`

#### 3. Resource Deleted
```json
{
  "success": false,
  "error": "User telah dihapus dan tidak dapat diakses"
}
```
**Status Code:** `404 Not Found`

### Permission Errors

#### 1. Insufficient Permissions
```json
{
  "success": false,
  "error": "Anda tidak memiliki akses ke resource ini"
}
```
**Status Code:** `403 Forbidden`

#### 2. Role Restriction
```json
{
  "success": false,
  "error": "Fitur ini hanya tersedia untuk ADMIN"
}
```
**Status Code:** `403 Forbidden`

### File Upload Errors

#### 1. File Too Large
```json
{
  "success": false,
  "error": "Ukuran file terlalu besar. Maksimal 5MB"
}
```
**Status Code:** `400 Bad Request`

#### 2. Invalid File Type
```json
{
  "success": false,
  "error": "Format file tidak didukung. Gunakan JPG, PNG, atau WebP"
}
```
**Status Code:** `400 Bad Request`

### Rate Limiting Errors

#### 1. Too Many Login Attempts
```json
{
  "success": false,
  "error": "Terlalu banyak percobaan login. Akun terkunci selama 15 menit"
}
```
**Status Code:** `429 Too Many Requests`

#### 2. API Rate Limit
```json
{
  "success": false,
  "error": "Terlalu banyak requests. Silakan coba lagi nanti"
}
```
**Status Code:** `429 Too Many Requests`

---

## 🛠️ Backend Implementation

### Go/Gin Error Handler

**Standard Error Response Function:**
```go
// pkg/utils/response.go
package utils

import (
    "github.com/gin-gonic/gin"
)

// ErrorResponse mengirim error response dengan format standar
func ErrorResponse(c *gin.Context, statusCode int, message string) {
    c.JSON(statusCode, gin.H{
        "success": false,
        "error":   message,
    })
}

// ErrorWithDetails mengirim error response dengan details
func ErrorWithDetails(c *gin.Context, statusCode int, message string, details interface{}) {
    c.JSON(statusCode, gin.H{
        "success": false,
        "error":   message,
        "details": details,
    })
}

// ValidationError mengirim validation error response
func ValidationError(c *gin.Context, errors map[string]string) {
    c.JSON(400, gin.H{
        "success": false,
        "error":   "Validasi gagal",
        "details": errors,
    })
}
```

**Usage Example:**
```go
// handlers/user_handler.go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    
    // Bind and validate request
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, 400, "Format request tidak valid")
        return
    }
    
    // Business logic validation
    if err := h.service.ValidateUser(&req); err != nil {
        utils.ErrorResponse(c, 422, err.Error())
        return
    }
    
    // Create user
    user, err := h.service.CreateUser(&req)
    if err != nil {
        utils.ErrorResponse(c, 500, "Gagal membuat user")
        return
    }
    
    utils.SuccessResponse(c, 201, "User berhasil dibuat", user)
}
```

---

## 🎨 Frontend Implementation

### Axios Error Interceptor

**Setup:**
```javascript
// src/utils/axios.js
import axios from 'axios'
import { useToast } from '@/composables/useToast'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Response interceptor untuk handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const toast = useToast()
    
    if (error.response) {
      const { status, data } = error.response
      
      // Handle specific status codes
      switch (status) {
        case 401:
          // Unauthorized - redirect to login
          toast.error(data.error || 'Token expired. Silakan login kembali')
          router.push('/auth/login')
          break
          
        case 403:
          // Forbidden
          toast.error(data.error || 'Anda tidak memiliki akses')
          break
          
        case 404:
          // Not Found
          toast.error(data.error || 'Resource tidak ditemukan')
          break
          
        case 422:
          // Validation Error
          toast.error(data.error || 'Validasi gagal')
          break
          
        case 429:
          // Rate Limit
          toast.error(data.error || 'Terlalu banyak requests')
          break
          
        case 500:
          // Server Error
          toast.error('Terjadi kesalahan pada server')
          break
          
        default:
          toast.error(data.error || 'Terjadi kesalahan')
      }
    } else if (error.request) {
      // Network error
      toast.error('Tidak dapat terhubung ke server')
    } else {
      toast.error('Terjadi kesalahan')
    }
    
    return Promise.reject(error)
  }
)

export default api
```

### Vue Component Error Handling

**Example:**
```vue
<script setup>
import { ref } from 'vue'
import api from '@/utils/axios'
import { useToast } from '@/composables/useToast'

const toast = useToast()
const loading = ref(false)
const errors = ref({})

const handleSubmit = async () => {
  loading.value = true
  errors.value = {}
  
  try {
    const response = await api.post('/api/users', formData.value)
    
    if (response.data.success) {
      toast.success(response.data.message)
      // Handle success
    }
  } catch (error) {
    if (error.response?.status === 400) {
      // Validation errors
      errors.value = error.response.data.details || {}
    }
    // Error sudah di-handle oleh interceptor
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <input 
      v-model="formData.email"
      :class="{ 'border-red-500': errors.email }"
    />
    <span v-if="errors.email" class="text-red-500">
      {{ errors.email }}
    </span>
  </form>
</template>
```

---

## ✅ Best Practices

### Backend Best Practices

1. **Consistent Error Messages**
   - Gunakan Bahasa Indonesia untuk user-facing errors
   - Gunakan naming convention yang jelas untuk error types
   - Log detailed errors untuk debugging (English OK)

2. **Security Considerations**
   - Jangan expose sensitive information dalam error messages
   - Prevent enumeration attacks (e.g., "User not found" vs "Invalid credentials")
   - Log security-related errors untuk monitoring

3. **Error Logging**
   ```go
   // Log error untuk debugging, tapi jangan expose ke user
   log.Printf("ERROR: Failed to create user: %v", err)
   utils.ErrorResponse(c, 500, "Gagal membuat user")
   ```

4. **Validation**
   - Validate input pada endpoint handler
   - Return specific validation errors untuk setiap field
   - Use descriptive error messages

### Frontend Best Practices

1. **User-Friendly Messages**
   - Display clear, actionable error messages
   - Provide suggestions untuk fix errors
   - Use toast notifications untuk errors yang tidak blocking

2. **Error Recovery**
   - Provide retry mechanisms untuk network errors
   - Auto-redirect ke login pada 401 errors
   - Preserve user input saat validation fails

3. **Loading States**
   - Show loading indicators during API calls
   - Disable form submission buttons saat processing
   - Provide feedback untuk long-running operations

4. **Offline Handling**
   - Detect offline state
   - Queue requests untuk when online kembali
   - Show appropriate offline messages

---

## 🧪 Testing Error Scenarios

### cURL Testing

**Test 401 Unauthorized:**
```bash
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer invalid_token"
```

**Test 400 Validation Error:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "",
    "email": "invalid-email"
  }'
```

**Test 404 Not Found:**
```bash
curl http://localhost:8080/api/users/99999 \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📞 Support

Jika ada pertanyaan terkait error handling:

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733  
**Timezone:** Asia/Jakarta (WIB)

---

**Last Updated:** 28 Desember 2025  
**Status:** ✅ Production Ready
