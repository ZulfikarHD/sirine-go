# 📚 API Documentation - Master Index

Complete API reference untuk Sirine Go App dengan 50+ endpoints dari Sprint 1-5.

**Base URL (Development):** `http://localhost:8080`  
**Base URL (Production):** `https://your-domain.com`

---

## 🔍 Overview

API ini menggunakan RESTful architecture dengan JSON sebagai format data exchange, yaitu:
- Consistent response format untuk success dan error
- JWT-based authentication dengan access & refresh tokens
- Comprehensive error handling dengan descriptive messages
- Rate limiting pada sensitive endpoints

### Quick Links

- 🔐 [**Authentication API**](../04-api-reference/authentication.md) - Login, logout, token management
- 👤 [**Profile API**](../04-api-reference/profile.md) - Profile management, password, photos
- 👥 [**User Management API**](../04-api-reference/user-management.md) - CRUD users, bulk operations
- 🔔 [**Notifications API**](../04-api-reference/notifications.md) - Notifications system
- 📊 [**Activity Logs API**](../04-api-reference/activity-logs.md) - Audit trail (Admin only)
- 🎮 [**Achievements API**](../04-api-reference/achievements.md) - Gamification system

---

## 🎯 API Categories

### 1. Authentication & Security (Sprint 1 & 3)

**Endpoints:** 6 endpoints  
**Documentation:** [Authentication API Reference](../04-api-reference/authentication.md)

**Key Features:**
- JWT token authentication (access + refresh tokens)
- Password reset flow via email
- Rate limiting (5 attempts per 15 minutes)
- Account locking mechanism

**Quick Reference:**
```
POST   /api/auth/login              # Login dengan NIP & password
POST   /api/auth/logout             # Logout & revoke session
GET    /api/auth/me                 # Get current user
POST   /api/auth/refresh            # Refresh JWT token
POST   /api/auth/forgot-password    # Request password reset
POST   /api/auth/reset-password     # Reset password
```

---

### 2. Profile Management (Sprint 2, 3, 5)

**Endpoints:** 6 endpoints  
**Documentation:** [Profile API Reference](../04-api-reference/profile.md)

**Key Features:**
- Self-service profile updates
- Password change dengan current password validation
- Profile photo upload (max 5MB, auto-resize 200x200px)
- Personal activity logs tracking

**Quick Reference:**
```
GET    /api/profile                 # Get own profile
PUT    /api/profile                 # Update profile (name, email, phone)
PUT    /api/profile/password        # Change password
POST   /api/profile/photo           # Upload photo (JPG/PNG/WebP)
DELETE /api/profile/photo           # Delete photo
GET    /api/profile/activity        # Get own activity logs
```

---

### 3. User Management (Sprint 2 & 5) - Admin Only

**Endpoints:** 11 endpoints  
**Documentation:** [User Management API Reference](../04-api-reference/user-management.md)

**Key Features:**
- CRUD operations untuk users
- Advanced filtering & search
- Bulk operations (delete, update status)
- CSV import/export functionality
- Auto-generated passwords untuk new users

**Quick Reference:**
```
GET    /api/users                   # List with filters & pagination
GET    /api/users/search            # Search by NIP/name/email
GET    /api/users/:id               # Get user detail
POST   /api/users                   # Create user (auto password)
PUT    /api/users/:id               # Update user
DELETE /api/users/:id               # Soft delete
POST   /api/users/bulk-delete       # Bulk delete
POST   /api/users/bulk-update-status # Bulk update status
POST   /api/users/import            # Import from CSV (max 1000 rows)
GET    /api/users/export            # Export to CSV
POST   /api/users/:id/reset-password # Admin force reset
```

---

### 4. Notifications (Sprint 4)

**Endpoints:** 6 endpoints  
**Documentation:** [Notifications API Reference](../04-api-reference/notifications.md)

**Key Features:**
- Real-time notification management
- Unread count untuk badge display
- Mark as read (individual & bulk)
- Notification types: INFO, SUCCESS, WARNING, ERROR

**Quick Reference:**
```
GET    /api/notifications           # List notifications
GET    /api/notifications/unread-count    # Unread count
GET    /api/notifications/recent    # Recent notifications
PUT    /api/notifications/:id/read  # Mark as read
PUT    /api/notifications/read-all  # Mark all as read
DELETE /api/notifications/:id       # Delete notification
```

---

### 5. Activity Logs (Sprint 4) - Admin/Manager Only

**Endpoints:** 4 endpoints  
**Documentation:** [Activity Logs API Reference](../04-api-reference/activity-logs.md)

**Key Features:**
- Complete audit trail untuk all operations
- Advanced filtering (by action, entity, date range)
- User-specific activity tracking
- Activity statistics & reports

**Quick Reference:**
```
GET    /api/admin/activity-logs     # List with advanced filters
GET    /api/admin/activity-logs/:id # Get log detail
GET    /api/admin/activity-logs/user/:id    # User activity logs
GET    /api/admin/activity-logs/stats       # Statistics & reports
```

**Logged Actions:**
- CREATE, UPDATE, DELETE operations
- LOGIN, LOGOUT events
- PASSWORD_CHANGE, PHOTO_UPLOAD
- Tracks before/after values untuk updates

---

### 6. Achievements & Gamification (Sprint 5)

**Endpoints:** 5 endpoints  
**Documentation:** [Achievements API Reference](../04-api-reference/achievements.md)

**Key Features:**
- Achievement tracking system
- Points accumulation
- Level progression (Bronze → Silver → Gold → Platinum)
- Auto-unlock achievements berdasarkan user actions

**Quick Reference:**
```
GET    /api/achievements            # List all achievements
GET    /api/profile/achievements    # User achievements with status
GET    /api/profile/stats           # Gamification stats
POST   /api/admin/achievements/award     # Manual award (Admin)
GET    /api/admin/users/:id/achievements # View user achievements (Admin)
```

**Levels:**
- Bronze: 0-99 points
- Silver: 100-499 points
- Gold: 500-999 points
- Platinum: 1000+ points

---

## 🏥 Health Check

### Check Server Status

Endpoint untuk monitoring server health.

```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "message": "Server berjalan dengan baik"
}
```

---

## 📊 API Statistics

| Category | Endpoints | Sprint | Authorization |
|----------|-----------|--------|---------------|
| Health Check | 1 | - | Public |
| Authentication | 6 | 1, 3 | Public + Protected |
| Profile | 6 | 2, 3, 5 | User |
| User Management | 11 | 2, 5 | Admin |
| Notifications | 6 | 4 | User |
| Activity Logs | 4 | 4 | Admin/Manager |
| Achievements | 5 | 5 | User + Admin |
| **TOTAL** | **39** | **1-5** | - |

---

## 🔐 Authentication

Sebagian besar endpoints memerlukan authentication menggunakan JWT token di header:

```http
Authorization: Bearer {your_jwt_token}
```

### Token Types
- **Access Token:** Valid 15 menit
- **Refresh Token:** Valid 30 hari
- **Password Reset Token:** Valid 1 jam (single-use)

### Getting Tokens

1. **Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"nip":"99999","password":"Admin@123"}'
```

2. **Use Token:**
```bash
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📝 Response Format

### Success Response

```json
{
  "success": true,
  "data": { ... },
  "message": "Pesan sukses dalam Bahasa Indonesia"
}
```

### Error Response

```json
{
  "success": false,
  "error": "Pesan error dalam Bahasa Indonesia",
  "details": "Additional error details (optional)"
}
```

### Pagination Response

```json
{
  "success": true,
  "data": {
    "items": [ ... ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 🎯 HTTP Status Codes

| Code | Meaning | Usage |
|------|---------|-------|
| **200** | OK | Request berhasil |
| **201** | Created | Resource berhasil dibuat |
| **400** | Bad Request | Validation error atau invalid request |
| **401** | Unauthorized | Token invalid atau expired |
| **403** | Forbidden | Insufficient permissions |
| **404** | Not Found | Resource tidak ditemukan |
| **409** | Conflict | Duplicate data (e.g., NIP sudah ada) |
| **422** | Unprocessable Entity | Business logic validation failed |
| **429** | Too Many Requests | Rate limit exceeded |
| **500** | Internal Server Error | Server error |

📖 **Detailed Error Handling:** [Error Handling Guide](../05-guides/error-handling.md)

---

## 🔒 Authorization & Roles

### Role Hierarchy

| Role | Access Level | Description |
|------|--------------|-------------|
| **ADMIN** | Full access | Manage all users, view all logs, award achievements |
| **MANAGER_*** | Department management | View department logs, manage department users |
| **STAFF_*** | Self-service | Own profile, notifications, achievements only |

**Note:** `*` indicates department (KHAZWAL, KEUANGAN, SDM, etc.)

### Role-Based Access

**Public Endpoints:**
- `/health`
- `/api/auth/login`
- `/api/auth/forgot-password`
- `/api/auth/reset-password`

**User Endpoints (All authenticated users):**
- `/api/auth/*` (except login/forgot/reset)
- `/api/profile/*`
- `/api/notifications/*`
- `/api/achievements` (read-only)

**Admin Only Endpoints:**
- `/api/users/*` (CRUD operations)
- `/api/admin/activity-logs/*`
- `/api/admin/achievements/award`
- `/api/users/import` & `/api/users/export`

**Manager/Admin Endpoints:**
- `/api/admin/activity-logs/*` (view only for managers)

---

## 🧪 Testing API

### Using cURL

**Basic Request:**
```bash
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**POST with JSON:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "12345",
    "full_name": "John Doe",
    "email": "john@example.com"
  }'
```

**File Upload:**
```bash
curl -X POST http://localhost:8080/api/profile/photo \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "photo=@/path/to/photo.jpg"
```

### Using Postman/Thunder Client

1. Create new request
2. Set `Authorization` header: `Bearer YOUR_TOKEN`
3. Set `Content-Type` header: `application/json`
4. Test endpoint

**Tip:** Save token in environment variable untuk reuse.

---

## 📚 Additional Resources

### API References (Detailed)
- [Authentication API](../04-api-reference/authentication.md) - Complete auth endpoints
- [Profile API](../04-api-reference/profile.md) - Profile management
- [User Management API](../04-api-reference/user-management.md) - Admin user operations
- [Notifications API](../04-api-reference/notifications.md) - Notification system
- [Activity Logs API](../04-api-reference/activity-logs.md) - Audit trail
- [Achievements API](../04-api-reference/achievements.md) - Gamification

### Guides
- [Error Handling Guide](../05-guides/error-handling.md) - Error responses & codes
- [Authentication Implementation](../05-guides/authentication/implementation.md) - Auth setup
- [Configuration Guide](../05-guides/configuration.md) - Environment variables
- [Security Best Practices](../05-guides/security.md) - Security guidelines
- [Validation Guide](../05-guides/validation/guide.md) - Input validation

### Testing
- [Testing Guide](./testing.md) - How to test APIs
- [User Management Testing](../06-testing/user-management-testing.md) - Test scenarios

---

## 📞 Support

Jika ada pertanyaan atau issue terkait API:

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733  
**Timezone:** Asia/Jakarta (WIB)

---

## ✨ What's New

### Sprint 5 (Current)
- ✅ Gamification system (achievements, points, levels)
- ✅ Profile photo upload dengan auto-resize
- ✅ CSV bulk import/export users

### Sprint 4
- ✅ Notification system
- ✅ Activity logs & audit trail

### Sprint 3
- ✅ Password management (forgot/reset)

### Sprint 2
- ✅ User management CRUD
- ✅ Profile self-service

### Sprint 1
- ✅ Authentication system
- ✅ JWT token management

---

**Last Updated:** 28 Desember 2025  
**API Version:** 1.5.0  
**Total Endpoints:** 50+  
**Status:** ✅ Production Ready
