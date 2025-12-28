# User Management & Profile API Documentation

**Feature**: User Management & Profile System  
**Version**: 1.1.0  
**Last Updated**: 28 Desember 2025

---

## Overview

User Management & Profile API merupakan RESTful endpoints untuk manage users (Admin) dan profile management (self-service), yaitu: CRUD operations untuk users, search & filter functionality, bulk operations, dan self-service profile updates.

**Base URL**: `http://localhost:8080/api`

---

## Endpoints Summary

### User Management (Admin/Manager)

| Method | Endpoint | Auth Required | Roles | Description |
|--------|----------|---------------|-------|-------------|
| GET | `/users` | ✅ Yes | Admin, Manager | List users dengan filters & pagination |
| GET | `/users/search` | ✅ Yes | Admin, Manager | Search users by NIP/name/email |
| GET | `/users/:id` | ✅ Yes | Admin, Manager | Get user detail by ID |
| POST | `/users` | ✅ Yes | Admin | Create new user |
| PUT | `/users/:id` | ✅ Yes | Admin | Update user data |
| DELETE | `/users/:id` | ✅ Yes | Admin | Soft delete user |
| POST | `/users/bulk-delete` | ✅ Yes | Admin | Bulk soft delete users |
| POST | `/users/bulk-update-status` | ✅ Yes | Admin | Bulk update user status |

### Profile Management (Self-Service)

| Method | Endpoint | Auth Required | Roles | Description |
|--------|----------|---------------|-------|-------------|
| GET | `/profile` | ✅ Yes | All | Get own profile |
| PUT | `/profile` | ✅ Yes | All | Update own profile (limited fields) |

---

## User Management APIs

### GET /api/users

**Description**: List users dengan filters dan pagination untuk admin panel.

**Authentication**: Required (Bearer token)  
**Authorization**: Admin, Manager

#### Request

**Headers**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters**:

| Field | Type | Required | Description | Default |
|-------|------|----------|-------------|---------|
| `page` | integer | ❌ No | Page number | 1 |
| `per_page` | integer | ❌ No | Users per page (max 100) | 20 |
| `role` | string | ❌ No | Filter by role (ADMIN, STAFF_KHAZWAL, etc) | - |
| `department` | string | ❌ No | Filter by department (KHAZWAL, CETAK, etc) | - |
| `status` | string | ❌ No | Filter by status (ACTIVE, INACTIVE, SUSPENDED) | - |
| `search` | string | ❌ No | Search by NIP or name (case-insensitive) | - |

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": 1,
        "nip": "99999",
        "full_name": "Administrator",
        "email": "admin@sirine.local",
        "phone": "081234567890",
        "role": "ADMIN",
        "department": "KHAZWAL",
        "shift": "PAGI",
        "profile_photo_url": "",
        "status": "ACTIVE",
        "last_login_at": "2025-12-28T10:30:00+07:00",
        "created_at": "2025-12-20T10:00:00+07:00",
        "updated_at": "2025-12-28T10:30:00+07:00"
      }
    ],
    "total": 25,
    "page": 1,
    "per_page": 20,
    "total_pages": 2
  }
}
```

#### Examples

**cURL**:
```bash
# List all users
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN"

# With filters
curl "http://localhost:8080/api/users?role=ADMIN&page=1&per_page=10" \
  -H "Authorization: Bearer $TOKEN"

# Search
curl "http://localhost:8080/api/users?search=Admin" \
  -H "Authorization: Bearer $TOKEN"
```

**JavaScript**:
```javascript
// Using apiClient (auto auth headers)
const response = await apiClient.get('/users', {
  params: {
    role: 'ADMIN',
    page: 1,
    per_page: 20
  }
});

const { users, total, page, total_pages } = response.data.data;
```

---

### POST /api/users

**Description**: Create new user dengan auto-generated password.

**Authentication**: Required  
**Authorization**: Admin only

#### Request

**Headers**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Body**:
```json
{
  "nip": "12345",
  "full_name": "Test User",
  "email": "test@example.com",
  "phone": "081234567890",
  "role": "STAFF_KHAZWAL",
  "department": "KHAZWAL",
  "shift": "PAGI"
}
```

**Fields**:

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `nip` | string | ✅ Yes | Max 5 char, unique |
| `full_name` | string | ✅ Yes | Min 3, max 100 char |
| `email` | string | ✅ Yes | Valid email, unique |
| `phone` | string | ✅ Yes | Min 10, max 15 char, format: 08xxx |
| `role` | string | ✅ Yes | Enum: ADMIN, MANAGER, STAFF_KHAZWAL, etc |
| `department` | string | ✅ Yes | Enum: KHAZWAL, CETAK, VERIFIKASI, KHAZKHIR |
| `shift` | string | ❌ No | Enum: PAGI, SIANG, MALAM (default: PAGI) |

#### Response

**Success (201 Created)**:
```json
{
  "success": true,
  "message": "User berhasil dibuat",
  "data": {
    "user": {
      "id": 2,
      "nip": "12345",
      "full_name": "Test User",
      "email": "test@example.com",
      "phone": "081234567890",
      "role": "STAFF_KHAZWAL",
      "department": "KHAZWAL",
      "shift": "PAGI",
      "status": "ACTIVE",
      "created_at": "2025-12-28T14:30:00+07:00"
    },
    "generated_password": "A7b@xY3k9Mz2"
  }
}
```

**⚠️ Important**: The `generated_password` is only shown once. User must change password on first login (`must_change_password = true`).

#### Examples

**JavaScript**:
```javascript
const response = await apiClient.post('/users', {
  nip: '12345',
  full_name: 'Test User',
  email: 'test@example.com',
  phone: '081234567890',
  role: 'STAFF_KHAZWAL',
  department: 'KHAZWAL',
  shift: 'PAGI'
});

const { user, generated_password } = response.data.data;
// Display password to admin
alert(`Password: ${generated_password}`);
```

---

### PUT /api/users/:id

**Description**: Update user data (Admin can update all fields).

**Authentication**: Required  
**Authorization**: Admin only

#### Request

**Headers**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Body** (all fields optional):
```json
{
  "full_name": "Updated Name",
  "email": "updated@example.com",
  "phone": "089876543210",
  "role": "OPERATOR_CETAK",
  "department": "CETAK",
  "shift": "SIANG",
  "status": "ACTIVE"
}
```

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "message": "User berhasil diupdate",
  "data": {
    "id": 2,
    "nip": "12345",
    "full_name": "Updated Name",
    "email": "updated@example.com",
    "role": "OPERATOR_CETAK",
    "department": "CETAK",
    "status": "ACTIVE",
    "updated_at": "2025-12-28T15:00:00+07:00"
  }
}
```

---

### DELETE /api/users/:id

**Description**: Soft delete user (data retained untuk audit).

**Authentication**: Required  
**Authorization**: Admin only

**⚠️ Important**: Cannot delete self (current logged-in user).

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "message": "User berhasil dihapus"
}
```

**Error (400 Bad Request)** - Trying to delete self:
```json
{
  "error": "Anda tidak dapat menghapus akun sendiri"
}
```

---

### GET /api/users/search

**Description**: Search users by NIP, name, or email (for autocomplete, etc).

**Authentication**: Required  
**Authorization**: Admin, Manager

#### Request

**Query Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `q` | string | ✅ Yes | Search query (min 1 char) |

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "nip": "99999",
      "full_name": "Administrator",
      "email": "admin@sirine.local",
      "role": "ADMIN",
      "department": "KHAZWAL"
    }
  ]
}
```

**Note**: Max 10 results returned untuk performance.

---

## Profile Management APIs

### GET /api/profile

**Description**: Get own profile information (self-service).

**Authentication**: Required  
**Authorization**: All authenticated users

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "nip": "99999",
    "full_name": "Administrator",
    "email": "admin@sirine.local",
    "phone": "081234567890",
    "role": "ADMIN",
    "department": "KHAZWAL",
    "shift": "PAGI",
    "profile_photo_url": "",
    "status": "ACTIVE",
    "last_login_at": "2025-12-28T10:30:00+07:00",
    "created_at": "2025-12-20T10:00:00+07:00",
    "updated_at": "2025-12-28T10:30:00+07:00"
  }
}
```

---

### PUT /api/profile

**Description**: Update own profile (limited fields only).

**Authentication**: Required  
**Authorization**: All authenticated users

**⚠️ Restrictions**: User can ONLY update:
- ✅ `full_name`
- ✅ `email` 
- ✅ `phone`

**Cannot update**: NIP, role, department, shift, status (admin-only fields).

#### Request

**Body**:
```json
{
  "full_name": "Updated Name",
  "email": "newemail@example.com",
  "phone": "089999999999"
}
```

**Validation**:
- `full_name`: Required, min 3, max 100 char
- `email`: Required, valid email format, unique
- `phone`: Required, min 10, max 15 char, format: 08xxx

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "message": "Profile berhasil diupdate",
  "data": {
    "id": 1,
    "nip": "99999",
    "full_name": "Updated Name",
    "email": "newemail@example.com",
    "phone": "089999999999",
    "role": "ADMIN",
    "department": "KHAZWAL",
    "updated_at": "2025-12-28T15:30:00+07:00"
  }
}
```

**Error (400 Bad Request)** - Email already used:
```json
{
  "error": "email sudah digunakan oleh user lain"
}
```

#### Examples

**JavaScript**:
```javascript
const response = await apiClient.put('/profile', {
  full_name: 'Updated Name',
  email: 'newemail@example.com',
  phone: '089999999999'
});

// Refresh auth store dengan updated data
await authStore.fetchCurrentUser();
```

---

## Activity Logging

All user management operations automatically logged ke `activity_logs` table:

**Logged Actions**:
- ✅ CREATE user
- ✅ UPDATE user
- ✅ DELETE user  
- ✅ UPDATE profile (self-service)

**Captured Data**:
- User ID (who performed action)
- Action type (CREATE/UPDATE/DELETE)
- Entity type (users)
- Entity ID (affected user ID)
- Before/After changes (JSON)
- IP address
- User agent
- Timestamp

**Example Activity Log**:
```json
{
  "id": 1,
  "user_id": 1,
  "action": "UPDATE",
  "entity_type": "users",
  "entity_id": 2,
  "changes": {
    "before": {
      "full_name": "Old Name",
      "role": "STAFF_KHAZWAL"
    },
    "after": {
      "full_name": "New Name",
      "role": "OPERATOR_CETAK"
    }
  },
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "created_at": "2025-12-28T15:00:00+07:00"
}
```

---

## Role-Based Access Control

### User Management Access

| Role | View Users | Create | Update | Delete |
|------|------------|--------|--------|--------|
| **Admin** | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |
| **Manager** | ✅ Yes | ❌ No | ❌ No | ❌ No |
| **Staff** | ❌ No | ❌ No | ❌ No | ❌ No |

### Profile Management Access

| Role | View Own | Update Own |
|------|----------|------------|
| **All Users** | ✅ Yes | ✅ Yes (limited) |

---

## Error Codes

| HTTP Code | Error Type | Description |
|-----------|------------|-------------|
| 200 | Success | Request successful |
| 201 | Created | User created successfully |
| 400 | Bad Request | Validation error, duplicate NIP/email |
| 401 | Unauthorized | Missing/invalid auth token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | User not found |
| 500 | Internal Server Error | Database error |

---

## Common Error Responses

**Duplicate NIP (400)**:
```json
{
  "error": "NIP sudah terdaftar dalam sistem"
}
```

**Duplicate Email (400)**:
```json
{
  "error": "email sudah terdaftar dalam sistem"
}
```

**User Not Found (404)**:
```json
{
  "error": "user dengan ID 999 tidak ditemukan"
}
```

**Insufficient Permissions (403)**:
```json
{
  "error": "Anda tidak memiliki akses ke resource ini"
}
```

---

## Related Documentation

- **Authentication API**: [auth-api-reference.md](./authentication/api-reference.md)
- **Sprint 2 Summary**: [SPRINT2_SUMMARY.md](../../SPRINT2_SUMMARY.md)
- **Testing Guide**: [SPRINT2_TESTING_GUIDE.md](../../SPRINT2_TESTING_GUIDE.md)
- **Activity Logs**: See Sprint 4 documentation (future)

---

**Last Updated**: 28 Desember 2025  
**Version**: 1.1.0 - Sprint 2  
**Status**: ✅ Complete & Production Ready
