# 🔌 API Documentation

## Overview

Central documentation hub untuk semua API endpoints dalam Sirine Go App.

**Base URL**: `http://localhost:8080/api`

---

## General Reference

- [API Design Guidelines](../03-development/api-documentation.md)
- [API Configuration](../05-guides/configuration.md)
- [Error Handling](../05-guides/error-handling.md) - Response format & status codes

---

## Feature APIs

### Authentication (Sprint 1 & 3)
- [**Authentication API**](./authentication.md) - Login, logout, refresh token, forgot/reset password

### Profile Management (Sprint 2, 3, 5)
- [**Profile API**](./profile.md) - Get/update profile, change password, photo upload, activity logs

### User Management (Sprint 2 & 5)
- [**User Management API**](./user-management.md) - CRUD users, search, filters, bulk operations

### Notifications (Sprint 4)
- [**Notifications API**](./notifications.md) - List notifications, mark as read, unread count

### Activity Logs (Sprint 4)
- [**Activity Logs API**](./activity-logs.md) - Audit trail, user activity tracking, statistics (Admin only)

### Gamification (Sprint 5)
- [**Achievements API**](./achievements.md) - Achievements, points system, level progression

---

## Quick API Reference

### 🔐 Authentication
```
POST   /api/auth/login              # Login dengan NIP & password
POST   /api/auth/logout             # Logout & revoke session
POST   /api/auth/refresh            # Refresh JWT token
GET    /api/auth/me                 # Get current user info
POST   /api/auth/forgot-password    # Request reset password
POST   /api/auth/reset-password     # Reset password dengan token
```
📖 [Complete Authentication API Reference](./authentication.md)

### 👤 Profile (Self-Service)
```
GET    /api/profile                 # Get own profile
PUT    /api/profile                 # Update own profile
PUT    /api/profile/password        # Change password
POST   /api/profile/photo           # Upload profile photo
DELETE /api/profile/photo           # Delete profile photo
GET    /api/profile/activity        # Get own activity logs
```
📖 [Complete Profile API Reference](./profile.md)

### 👥 User Management (Admin)
```
GET    /api/users                   # List users (filters, pagination)
GET    /api/users/search            # Search users
GET    /api/users/:id               # Get user detail
POST   /api/users                   # Create user (auto password)
PUT    /api/users/:id               # Update user
DELETE /api/users/:id               # Delete user (soft)
POST   /api/users/bulk-delete       # Bulk delete users
POST   /api/users/import            # Import users from CSV
GET    /api/users/export            # Export users to CSV
```
📖 [Complete User Management API Reference](./user-management.md)

### 🔔 Notifications
```
GET    /api/notifications           # List notifications
GET    /api/notifications/unread-count    # Get unread count
GET    /api/notifications/recent    # Get recent notifications
PUT    /api/notifications/:id/read  # Mark as read
PUT    /api/notifications/read-all  # Mark all as read
DELETE /api/notifications/:id       # Delete notification
```
📖 [Complete Notifications API Reference](./notifications.md)

### 📊 Activity Logs (Admin/Manager)
```
GET    /api/admin/activity-logs     # List activity logs
GET    /api/admin/activity-logs/:id # Get log detail
GET    /api/admin/activity-logs/user/:id    # Get user logs
GET    /api/admin/activity-logs/stats       # Get statistics
```
📖 [Complete Activity Logs API Reference](./activity-logs.md)

### 🎮 Achievements
```
GET    /api/achievements            # List all achievements
GET    /api/profile/achievements    # Get user achievements
GET    /api/profile/stats           # Get gamification stats
POST   /api/admin/achievements/award     # Award achievement (Admin)
GET    /api/admin/users/:id/achievements # Get user achievements (Admin)
```
📖 [Complete Achievements API Reference](./achievements.md)

---

## API Conventions

### Request Headers
```
Content-Type: application/json
Authorization: Bearer {token}    # For protected endpoints
```

### Response Format

**Success**:
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

**Error**:
```json
{
  "success": false,
  "error": "Error message in Indonesian",
  "details": "Additional error details (optional)"
}
```

### HTTP Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation error)
- `401` - Unauthorized (missing/invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error

---

## API Versioning

**Current Version**: 1.1.0

Version history:
- `1.0.0` - Sprint 1: Authentication System
- `1.1.0` - Sprint 2: User Management & Profile

---

## Testing APIs

### Development
```bash
# Use cURL
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN"

# Use Postman/Thunder Client
# Import collection from docs/testing/
```

### Production
```bash
# Replace with production URL
curl https://api.sirine.app/api/users \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📊 API Statistics

**Total Endpoints:** 50+  
**Sprints Covered:** Sprint 1-5  
**Documentation Status:** ✅ Complete

### Endpoints by Category
- Authentication: 6 endpoints
- Profile: 6 endpoints
- User Management: 11 endpoints
- Notifications: 6 endpoints
- Activity Logs: 4 endpoints
- Achievements: 5 endpoints
- Health Check: 1 endpoint

---

## 📚 Additional Resources

- [Complete API Documentation](../03-development/api-documentation.md) - Master API reference
- [Error Handling Guide](../05-guides/error-handling.md) - Response formats & error codes
- [Authentication Guide](../05-guides/authentication/implementation.md) - Auth implementation
- [Configuration Guide](../05-guides/configuration.md) - Environment setup
- [Security Guide](../05-guides/security.md) - Security best practices

---

**Last Updated**: 28 Desember 2025  
**API Version**: 1.5.0  
**Status**: ✅ Production Ready
