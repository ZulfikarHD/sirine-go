# 👤 Profile API Reference

Complete API reference untuk profile endpoints dalam Sirine Go App.

**Base URL:** `http://localhost:8080/api/profile`

---

## 📋 Overview

Profile endpoints memungkinkan users untuk manage profile mereka sendiri (self-service), termasuk:
- View dan update profile information
- Change password
- Upload profile photo
- View activity logs

**Authentication:** Semua endpoints memerlukan valid JWT token.

---

## 🔑 Endpoints

### 1. Get Own Profile

Mengambil informasi profile user yang sedang login.

```http
GET /api/profile
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
    "profile_photo": "/uploads/profiles/1.jpg"
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid/expired

---

### 2. Update Own Profile

Update profile information (full_name, email, phone).

```http
PUT /api/profile
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "full_name": "Super Admin Updated",
  "email": "admin.new@sirine.local",
  "phone": "081234567899"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "full_name": "Super Admin Updated",
    "email": "admin.new@sirine.local",
    "phone": "081234567899"
  },
  "message": "Profile berhasil diperbarui"
}
```

**Validation Rules:**
- `full_name`: Required, min 3 characters
- `email`: Required, valid email format
- `phone`: Required, 10-13 digits

**HTTP Status Codes:**
- `200 OK` - Update berhasil
- `400 Bad Request` - Validation error
- `401 Unauthorized` - Token invalid
- `409 Conflict` - Email already exists

---

### 3. Change Password

Change own password (requires current password).

```http
PUT /api/profile/password
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "current_password": "OldPassword@123",
  "new_password": "NewPassword@123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password berhasil diubah. Silakan login kembali"
}
```

**Password Requirements:**
- Minimal 8 karakter
- Mengandung huruf besar, huruf kecil, angka, dan simbol
- Berbeda dari password sebelumnya

**Note:** Setelah password change, user akan di-logout otomatis dan harus login kembali dengan password baru.

**HTTP Status Codes:**
- `200 OK` - Password changed
- `400 Bad Request` - Current password salah
- `401 Unauthorized` - Token invalid
- `422 Unprocessable Entity` - Password tidak memenuhi kriteria

---

### 4. Upload Profile Photo

Upload profile photo (max 5MB, format: JPG/PNG/WebP).

```http
POST /api/profile/photo
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**Request Body:**
```
photo: [file]
```

**Response:**
```json
{
  "success": true,
  "data": {
    "photo_url": "/uploads/profiles/1.jpg"
  },
  "message": "Foto profil berhasil diupload"
}
```

**File Requirements:**
- Max size: 5MB
- Formats: JPG, JPEG, PNG, WebP
- Auto-resize: 200x200px
- Quality: 90% JPEG

**HTTP Status Codes:**
- `200 OK` - Upload berhasil
- `400 Bad Request` - Invalid file type/size
- `401 Unauthorized` - Token invalid

---

### 5. Delete Profile Photo

Hapus profile photo dan kembali ke default avatar.

```http
DELETE /api/profile/photo
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "message": "Foto profil berhasil dihapus"
}
```

**HTTP Status Codes:**
- `200 OK` - Delete berhasil
- `401 Unauthorized` - Token invalid

---

### 6. Get Own Activity Logs

View own activity logs dengan pagination.

```http
GET /api/profile/activity?page=1&page_size=20
Authorization: Bearer {token}
```

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Items per page (default: 20, max: 100)

**Response:**
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "id": 1,
        "action": "UPDATE",
        "entity_type": "profile",
        "entity_id": "1",
        "changes": {
          "before": {
            "email": "old@email.com"
          },
          "after": {
            "email": "new@email.com"
          }
        },
        "ip_address": "127.0.0.1",
        "user_agent": "Mozilla/5.0...",
        "created_at": "2025-12-28T10:00:00+07:00"
      }
    ],
    "total": 15,
    "page": 1,
    "page_size": 20
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid

---

## 🧪 Testing Examples

### cURL Examples

**Get Profile:**
```bash
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Update Profile:**
```bash
curl -X PUT http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "081234567890"
  }'
```

**Upload Photo:**
```bash
curl -X POST http://localhost:8080/api/profile/photo \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "photo=@/path/to/photo.jpg"
```

**Change Password:**
```bash
curl -X PUT http://localhost:8080/api/profile/password \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "OldPassword@123",
    "new_password": "NewPassword@123"
  }'
```

### JavaScript/Axios Examples

**Get Profile:**
```javascript
const response = await axios.get('/api/profile')
const profile = response.data.data
```

**Update Profile:**
```javascript
const response = await axios.put('/api/profile', {
  full_name: 'John Doe',
  email: 'john@example.com',
  phone: '081234567890'
})
```

**Upload Photo:**
```javascript
const formData = new FormData()
formData.append('photo', fileInput.files[0])

const response = await axios.post('/api/profile/photo', formData, {
  headers: {
    'Content-Type': 'multipart/form-data'
  }
})
```

---

## 📚 Related Documentation

- [Authentication API](./authentication.md)
- [User Management API](./user-management.md) - Admin endpoints
- [Error Handling Guide](../05-guides/error-handling.md)
- [File Upload Configuration](../05-guides/configuration.md)

---

**Last Updated:** 28 Desember 2025  
**Sprint:** Sprint 2, 3, 5  
**Status:** ✅ Production Ready
