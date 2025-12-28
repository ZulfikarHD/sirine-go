# 🔌 API Documentation

## Overview

Central documentation hub untuk semua API endpoints dalam Sirine Go App.

**Base URL**: `http://localhost:8080/api`

---

## General Reference

- [API Design Guidelines](../development/api-documentation.md)
- [API Configuration](../guides/configuration.md)
- [Error Handling](../guides/error-handling.md)

---

## Feature APIs

### Authentication & Authorization
- [Authentication API](../guides/authentication/api-reference.md) - Login, logout, token management

### User Management
- [User Management & Profile API](./user-management.md) - CRUD users, profile management

### Future Features
- **Sprint 3**: Password Management API (Forgot/reset password, change password)
- **Sprint 4**: Notifications & Activity Logs API
- **Sprint 5**: Achievements & Bulk Operations API

---

## Quick API Reference

### Authentication
```
POST   /api/auth/login      # Login dengan NIP & password
POST   /api/auth/logout     # Logout & revoke session
POST   /api/auth/refresh    # Refresh JWT token
GET    /api/auth/me         # Get current user info
```

### User Management (Admin)
```
GET    /api/users           # List users (filters, pagination)
GET    /api/users/search    # Search users
GET    /api/users/:id       # Get user detail
POST   /api/users           # Create user (auto password)
PUT    /api/users/:id       # Update user
DELETE /api/users/:id       # Delete user (soft)
```

### Profile (Self-Service)
```
GET    /api/profile         # Get own profile
PUT    /api/profile         # Update own profile
```

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

**Last Updated**: 28 Desember 2025
