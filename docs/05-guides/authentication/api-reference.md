# Authentication API Documentation

**Feature**: Authentication System  
**Version**: 1.0.0  
**Last Updated**: 27 Desember 2025

---

## Overview

Authentication API merupakan RESTful endpoints untuk user authentication yang bertujuan untuk manage login, logout, token refresh, dan user session, yaitu: JWT-based authentication dengan refresh token mechanism, role-based access control, dan comprehensive session tracking.

**Base URL**: `http://localhost:8080/api/auth`

---

## Endpoints Summary

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| POST | `/login` | ❌ No | Login dengan NIP dan password |
| POST | `/refresh` | ❌ No | Refresh JWT token |
| POST | `/logout` | ✅ Yes | Logout dan revoke session |
| GET | `/me` | ✅ Yes | Get current user info |

---

## POST /api/auth/login

**Description**: Login dengan NIP dan password untuk mendapatkan JWT access token dan refresh token.

**Authentication**: None (Public endpoint)

### Request

**Headers**:
```
Content-Type: application/json
```

**Body**:
```json
{
  "nip": "99999",
  "password": "Admin@123",
  "remember_me": false
}
```

**Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `nip` | string | ✅ Yes | User NIP (5 digit angka) |
| `password` | string | ✅ Yes | User password (min 8 char) |
| `remember_me` | boolean | ❌ No | Extend refresh token expiry (default: false) |

### Response

#### Success (200 OK)

```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJuaXAiOiI5OTk5OSIsInJvbGUiOiJBRE1JTiIsImV4cCI6MTczNTI5NjAwMH0.signature",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ0eXBlIjoicmVmcmVzaCIsImV4cCI6MTczNzg4ODAwMH0.signature",
    "user": {
      "id": 1,
      "nip": "99999",
      "full_name": "Administrator",
      "email": "admin@sirine.local",
      "phone": "081234567890",
      "role": "ADMIN",
      "department": "KHAZWAL",
      "shift": "PAGI",
      "status": "ACTIVE",
      "must_change_password": false,
      "last_login_at": "2025-12-27T14:30:00+07:00",
      "created_at": "2025-12-20T10:00:00+07:00",
      "updated_at": "2025-12-27T14:30:00+07:00"
    },
    "require_password_change": false
  }
}
```

**Response Fields**:

| Field | Type | Description |
|-------|------|-------------|
| `token` | string | JWT access token (expires in 15 min) |
| `refresh_token` | string | Refresh token (expires in 30 days) |
| `user` | object | User data (password excluded) |
| `require_password_change` | boolean | Flag untuk force password change |

#### Error Responses

**Invalid Credentials (401 Unauthorized)**:
```json
{
  "success": false,
  "message": "NIP atau password salah"
}
```

**Account Locked (401 Unauthorized)**:
```json
{
  "success": false,
  "message": "Akun Anda terkunci hingga 14:45:00 karena terlalu banyak percobaan login gagal"
}
```

**Account Inactive (403 Forbidden)**:
```json
{
  "success": false,
  "message": "Akun Anda tidak aktif. Silakan hubungi administrator"
}
```

**Validation Error (400 Bad Request)**:
```json
{
  "success": false,
  "message": "NIP dan password harus diisi"
}
```

### Examples

**cURL**:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "99999",
    "password": "Admin@123",
    "remember_me": true
  }'
```

**JavaScript (Axios)**:
```javascript
const response = await axios.post('/api/auth/login', {
  nip: '99999',
  password: 'Admin@123',
  remember_me: true
});

const { token, refresh_token, user } = response.data.data;
```

**Go**:
```go
loginReq := LoginRequest{
    NIP:        "99999",
    Password:   "Admin@123",
    RememberMe: true,
}

resp, err := authService.Login(loginReq, ipAddress, userAgent)
```

### Business Logic

1. **Validation**:
   - NIP dan password tidak boleh kosong
   - NIP harus 5 digit angka
   - Password min 8 karakter

2. **Authentication Flow**:
   - Find user by NIP
   - Check if account locked (`locked_until > now`)
   - Check if account active (`status = ACTIVE`)
   - Verify password dengan bcrypt
   - Reset `failed_login_attempts` to 0 on success
   - Increment `failed_login_attempts` on failure
   - Lock account after 5 failed attempts (15 min lockout)

3. **Token Generation**:
   - Generate JWT access token (15 min expiry)
   - Generate refresh token (30 days expiry)
   - Hash tokens dengan SHA256 untuk storage
   - Create `user_session` record

4. **Activity Logging**:
   - Log LOGIN action ke `activity_logs`
   - Record IP address dan user agent
   - Update `last_login_at` timestamp

---

## GET /api/auth/me

**Description**: Get current authenticated user information.

**Authentication**: Required (Bearer token)

### Request

**Headers**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Parameters**: None

### Response

#### Success (200 OK)

```json
{
  "success": true,
  "message": "Data user berhasil diambil",
  "data": {
    "id": 1,
    "nip": "99999",
    "full_name": "Administrator",
    "email": "admin@sirine.local",
    "phone": "081234567890",
    "role": "ADMIN",
    "department": "KHAZWAL",
    "shift": "PAGI",
    "status": "ACTIVE",
    "must_change_password": false,
    "last_login_at": "2025-12-27T14:30:00+07:00",
    "created_at": "2025-12-20T10:00:00+07:00",
    "updated_at": "2025-12-27T14:30:00+07:00"
  }
}
```

#### Error Responses

**Missing Token (401 Unauthorized)**:
```json
{
  "success": false,
  "message": "Authorization header tidak ditemukan"
}
```

**Invalid Token (401 Unauthorized)**:
```json
{
  "success": false,
  "message": "Token tidak valid atau sudah expired"
}
```

**User Not Found (401 Unauthorized)**:
```json
{
  "success": false,
  "message": "User tidak ditemukan"
}
```

### Examples

**cURL**:
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

**JavaScript (Axios)**:
```javascript
// Token auto-injected by axios interceptor
const response = await axios.get('/api/auth/me');
const user = response.data.data;
```

**Go**:
```go
req, _ := http.NewRequest("GET", "/api/auth/me", nil)
req.Header.Set("Authorization", "Bearer "+token)

resp, err := client.Do(req)
```

---

## POST /api/auth/logout

**Description**: Logout dan revoke current session. Token tidak dapat digunakan lagi setelah logout.

**Authentication**: Required (Bearer token)

### Request

**Headers**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Body**: None

### Response

#### Success (200 OK)

```json
{
  "success": true,
  "message": "Logout berhasil"
}
```

#### Error Responses

**Missing Token (401 Unauthorized)**:
```json
{
  "success": false,
  "message": "Authorization header tidak ditemukan"
}
```

**Invalid Token (401 Unauthorized)**:
```json
{
  "success": false,
  "message": "Token tidak valid atau sudah expired"
}
```

### Examples

**cURL**:
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

**JavaScript (Axios)**:
```javascript
await axios.post('/api/auth/logout');
// Clear local storage
localStorage.removeItem('auth_token');
localStorage.removeItem('refresh_token');
```

**Go**:
```go
err := authService.Logout(userID, token, ipAddress, userAgent)
```

### Business Logic

1. **Token Validation**:
   - Extract user dari JWT claims
   - Verify token signature dan expiry

2. **Session Revocation**:
   - Hash token dengan SHA256
   - Find session di `user_sessions` by token_hash
   - Set `is_revoked = true`

3. **Activity Logging**:
   - Log LOGOUT action ke `activity_logs`
   - Record IP address dan user agent
   - Timestamp recorded automatically

---

## POST /api/auth/refresh

**Description**: Refresh JWT access token menggunakan refresh token. Digunakan untuk extend session tanpa re-login.

**Authentication**: None (Uses refresh token)

### Request

**Headers**:
```
Content-Type: application/json
```

**Body**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `refresh_token` | string | ✅ Yes | Refresh token dari login response |

### Response

#### Success (200 OK)

```json
{
  "success": true,
  "message": "Token berhasil di-refresh",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "nip": "99999",
      "full_name": "Administrator",
      "email": "admin@sirine.local",
      "role": "ADMIN",
      "department": "KHAZWAL",
      "shift": "PAGI",
      "status": "ACTIVE"
    }
  }
}
```

#### Error Responses

**Invalid Refresh Token (401 Unauthorized)**:
```json
{
  "success": false,
  "message": "Refresh token tidak valid atau sudah expired"
}
```

**Revoked Session (401 Unauthorized)**:
```json
{
  "success": false,
  "message": "Session telah di-revoke. Silakan login kembali"
}
```

### Examples

**cURL**:
```bash
REFRESH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}"
```

**JavaScript (Axios)**:
```javascript
// Auto-handled by axios response interceptor on 401
const response = await axios.post('/api/auth/refresh', {
  refresh_token: localStorage.getItem('refresh_token')
});

const { token, refresh_token } = response.data.data;
localStorage.setItem('auth_token', token);
localStorage.setItem('refresh_token', refresh_token);
```

**Go**:
```go
resp, err := authService.RefreshAuthToken(refreshToken)
```

### Business Logic

1. **Refresh Token Validation**:
   - Parse JWT refresh token
   - Verify signature dan expiry
   - Check token type = "refresh"

2. **Session Validation**:
   - Find session by refresh_token_hash
   - Check `is_revoked = false`
   - Check `expires_at > now`

3. **New Token Generation**:
   - Generate new JWT access token (15 min)
   - Generate new refresh token (30 days)
   - Update session record dengan new token hashes
   - Return new tokens + user data

---

## Authentication Flow

### Initial Login Flow

```
Client                    Backend                   Database
  │                         │                          │
  ├─ POST /api/auth/login ──▶│                          │
  │  (NIP, password)         │                          │
  │                          ├─ Find user by NIP ──────▶│
  │                          │◀─ User data ─────────────┤
  │                          │                          │
  │                          ├─ Verify password         │
  │                          ├─ Generate JWT (15 min)   │
  │                          ├─ Generate refresh (30d)  │
  │                          │                          │
  │                          ├─ Create session ─────────▶│
  │                          ├─ Log activity ───────────▶│
  │                          │                          │
  │◀─ 200 OK ────────────────┤                          │
  │  (token, refresh_token)  │                          │
  │                          │                          │
  ├─ Store tokens            │                          │
  │  (localStorage)          │                          │
```

### Token Refresh Flow

```
Client                    Backend                   Database
  │                         │                          │
  ├─ GET /api/auth/me ──────▶│                          │
  │  (expired token)         │                          │
  │◀─ 401 Unauthorized ──────┤                          │
  │                          │                          │
  ├─ POST /api/auth/refresh ▶│                          │
  │  (refresh_token)         │                          │
  │                          ├─ Validate refresh token  │
  │                          ├─ Find session ───────────▶│
  │                          │◀─ Session data ──────────┤
  │                          │                          │
  │                          ├─ Generate new JWT        │
  │                          ├─ Generate new refresh    │
  │                          ├─ Update session ─────────▶│
  │                          │                          │
  │◀─ 200 OK ────────────────┤                          │
  │  (new tokens)            │                          │
  │                          │                          │
  ├─ Retry original request ▶│                          │
  │  (new token)             │                          │
  │◀─ 200 OK ────────────────┤                          │
```

### Logout Flow

```
Client                    Backend                   Database
  │                         │                          │
  ├─ POST /api/auth/logout ─▶│                          │
  │  (token)                 │                          │
  │                          ├─ Validate token          │
  │                          ├─ Find session ───────────▶│
  │                          ├─ Revoke session ─────────▶│
  │                          │  (is_revoked = true)     │
  │                          ├─ Log activity ───────────▶│
  │                          │                          │
  │◀─ 200 OK ────────────────┤                          │
  │                          │                          │
  ├─ Clear localStorage      │                          │
  ├─ Redirect to /login      │                          │
```

---

## Error Codes

| HTTP Code | Error Type | Description |
|-----------|------------|-------------|
| 200 | Success | Request berhasil |
| 400 | Bad Request | Validation error, missing required fields |
| 401 | Unauthorized | Invalid credentials, expired token, missing token |
| 403 | Forbidden | Account inactive/suspended, insufficient permissions |
| 500 | Internal Server Error | Database error, unexpected error |

---

## Rate Limiting

**Current Implementation**: Per-user rate limiting dengan account lockout

- Max failed attempts: **5**
- Lockout duration: **15 minutes**
- Counter reset: On successful login

**Future Enhancement**: IP-based rate limiting untuk prevent distributed attacks

---

## Security Notes

1. **Token Storage**:
   - Access token: Short-lived (15 min) untuk minimize exposure
   - Refresh token: Long-lived (30 days) untuk better UX
   - Tokens hashed dengan SHA256 di database (not plaintext)

2. **Password Security**:
   - Bcrypt hashing dengan cost 12 (~200ms)
   - Password never returned di API responses
   - Password policy enforced (min 8 char, uppercase, number, special)

3. **Session Security**:
   - IP address dan user agent tracked
   - Session revocation on logout
   - Activity logging untuk audit trail

4. **CORS**:
   - Allowed origins configured di backend
   - Credentials allowed untuk cookie support (future)

---

## Related Documentation

- **Feature Documentation**: [Authentication System](../features/AUTHENTICATION.md)
- **Test Plan**: [AUTH Test Plan](../06-testing/AUTH-test-plan.md)
- **User Journeys**: [Authentication User Journeys](../05-guides/authentication-user-journeys.md)
- **Testing Guide**: [TESTING_GUIDE.md](../../TESTING_GUIDE.md)

---

**Last Updated**: 27 Desember 2025  
**Version**: 1.0.0 - Sprint 1  
**Status**: ✅ API Documentation Complete
