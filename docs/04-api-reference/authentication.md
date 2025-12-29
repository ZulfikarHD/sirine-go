# 🔐 Authentication API Reference

Complete API reference untuk authentication endpoints dalam Sirine Go App.

**Base URL:** `http://localhost:8080/api/auth`

---

## 📋 Overview

Authentication menggunakan JWT (JSON Web Token) dengan access token dan refresh token untuk secure session management.

### Token Types
- **Access Token:** Valid 15 menit, digunakan untuk API requests
- **Refresh Token:** Valid 30 hari, digunakan untuk refresh access token

### Rate Limiting
- Max 5 login attempts per user dalam 15 menit
- Account akan terkunci selama 15 menit setelah 5 failed attempts

---

## 🔑 Endpoints

### 1. Login

Login menggunakan NIP dan password.

```http
POST /api/auth/login
Content-Type: application/json
```

**Request Body:**
```json
{
  "nip": "99999",
  "password": "Admin@123"
}
```

**Response (Success):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "nip": "99999",
      "full_name": "Super Admin",
      "email": "admin@sirine.local",
      "phone": "081234567890",
      "role": "ADMIN",
      "department": "ADMIN",
      "shift": "PAGI",
      "status": "ACTIVE",
      "total_points": 100,
      "level": "Silver",
      "profile_photo": "/uploads/profiles/1.jpg"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "require_password_change": false
  },
  "message": "Login berhasil"
}
```

**Response (Failed - Invalid Credentials):**
```json
{
  "success": false,
  "error": "NIP atau password salah"
}
```

**Response (Failed - Account Locked):**
```json
{
  "success": false,
  "error": "Akun terkunci. Coba lagi dalam 15 menit"
}
```

**HTTP Status Codes:**
- `200 OK` - Login berhasil
- `401 Unauthorized` - Invalid credentials
- `403 Forbidden` - Account locked
- `429 Too Many Requests` - Too many attempts

---

### 2. Logout

Logout dan revoke session.

```http
POST /api/auth/logout
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "message": "Logout berhasil"
}
```

**HTTP Status Codes:**
- `200 OK` - Logout berhasil
- `401 Unauthorized` - Token invalid

---

### 3. Get Current User

Mengambil informasi user yang sedang login.

```http
GET /api/auth/me
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "nip": "99999",
    "full_name": "Super Admin",
    "email": "admin@sirine.local",
    "phone": "081234567890",
    "role": "ADMIN",
    "department": "ADMIN",
    "shift": "PAGI",
    "status": "ACTIVE",
    "total_points": 100,
    "level": "Silver",
    "profile_photo": "/uploads/profiles/1.jpg",
    "created_at": "2025-12-20T10:00:00+07:00"
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid/expired

---

### 4. Refresh Token

Refresh JWT token yang expired.

```http
POST /api/auth/refresh
Content-Type: application/json
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Token refreshed
- `401 Unauthorized` - Invalid refresh token

---

### 5. Forgot Password

Request reset password link via email.

```http
POST /api/auth/forgot-password
Content-Type: application/json
```

**Request Body:**
```json
{
  "nip_or_email": "99999"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Link reset password telah dikirim ke email Anda"
}
```

**Note:** Response selalu success untuk prevent email enumeration attack.

**HTTP Status Codes:**
- `200 OK` - Request processed (always returns success)

---

### 6. Reset Password

Reset password menggunakan token dari email.

```http
POST /api/auth/reset-password
Content-Type: application/json
```

**Request Body:**
```json
{
  "token": "abc123def456...",
  "new_password": "NewPassword@123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password berhasil direset"
}
```

**HTTP Status Codes:**
- `200 OK` - Password reset berhasil
- `400 Bad Request` - Invalid/expired token
- `422 Unprocessable Entity` - Password tidak memenuhi kriteria

---

## 🔒 Security Notes

### Password Requirements
- Minimal 8 karakter
- Mengandung huruf besar (A-Z)
- Mengandung huruf kecil (a-z)
- Mengandung angka (0-9)
- Mengandung simbol (!@#$%^&*, dll)

### Token Security
- Access token expires dalam 15 menit
- Refresh token expires dalam 30 hari
- Password reset token expires dalam 1 jam dan single-use
- Tokens disimpan securely di backend (tidak di-share)

### Rate Limiting
- Login endpoint: Max 5 attempts dalam 15 menit per user
- Failed attempts akan lock account untuk 15 menit
- Forgot password endpoint: Max 3 requests dalam 1 jam

---

## 🧪 Testing Examples

### cURL Examples

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"nip":"99999","password":"Admin@123"}'
```

**Get Current User:**
```bash
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Refresh Token:**
```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'
```

### JavaScript/Axios Examples

**Login:**
```javascript
const response = await axios.post('/api/auth/login', {
  nip: '99999',
  password: 'Admin@123'
})

const { token, refresh_token, user } = response.data.data
```

**Set Authorization Header:**
```javascript
axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
```

---

## 📚 Related Documentation

- [Error Handling Guide](../05-guides/error-handling.md)
- [Security Best Practices](../05-guides/security.md)
- [Authentication Implementation Guide](../05-guides/authentication/implementation.md)
- [User Management API](./user-management.md)

---

**Last Updated:** 28 Desember 2025  
**Sprint:** Sprint 1 & 3  
**Status:** ✅ Production Ready
