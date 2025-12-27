# Authentication System - Sprint 1

**Status**: ✅ Complete  
**Priority**: P0 (Critical MVP)  
**Sprint**: 1 (Week 1)  
**Last Updated**: 27 Desember 2025

---

## Overview

Authentication System merupakan foundational security layer yang bertujuan untuk mengamankan akses aplikasi dengan JWT-based authentication, yaitu: login/logout flow, session management, role-based access control, dan rate limiting untuk mencegah brute force attacks. Sistem ini mencakup complete backend API dengan Go + Gin framework serta modern iOS-inspired frontend design menggunakan Vue 3 dengan Indigo & Fuchsia gradient theme.

---

## User Stories

| ID | User Story | Acceptance Criteria | Status |
|----|------------|---------------------|--------|
| AUTH-001 | Sebagai user, saya ingin login dengan NIP dan password agar dapat mengakses sistem | - Form login dengan NIP (5 digit) dan password<br>- Validasi credentials<br>- Redirect ke dashboard sesuai role<br>- Token tersimpan di localStorage | ✅ Complete |
| AUTH-002 | Sebagai user, saya ingin tetap logged in setelah refresh page | - Token persistence di localStorage<br>- Auto-restore auth state<br>- Tidak perlu login ulang | ✅ Complete |
| AUTH-003 | Sebagai user, saya ingin logout dari sistem | - Button logout di navbar dropdown<br>- Token dihapus<br>- Redirect ke login page<br>- Session di-revoke di backend | ✅ Complete |
| AUTH-004 | Sebagai admin, saya ingin akses dilindungi berdasarkan role | - Protected routes dengan auth guard<br>- Role-based route access<br>- 403 untuk unauthorized access | ✅ Complete |
| AUTH-005 | Sebagai sistem, saya ingin mencegah brute force login | - Rate limiting: max 5 failed attempts<br>- Account lockout selama 15 menit<br>- Error message dengan countdown | ✅ Complete |
| AUTH-006 | Sebagai user, saya ingin token auto-refresh saat expired | - JWT expiry 15 menit<br>- Refresh token 30 hari<br>- Auto-refresh on 401<br>- Seamless user experience | ✅ Complete |

---

## Business Rules

| Rule ID | Rule | Validation | Impact |
|---------|------|------------|--------|
| BR-AUTH-001 | NIP harus 5 digit angka | Frontend: max 5 char, numeric only<br>Backend: varchar(5) unique | Login gagal jika format salah |
| BR-AUTH-002 | Password policy: min 8 char, 1 uppercase, 1 number, 1 special char | Backend: regex validation<br>Frontend: password strength indicator | Password weak ditolak |
| BR-AUTH-003 | JWT token expiry 15 menit | Backend: JWT_EXPIRY config<br>Auto-refresh mechanism | Session security |
| BR-AUTH-004 | Refresh token expiry 30 hari | Backend: REFRESH_TOKEN_EXPIRY config | Remember me functionality |
| BR-AUTH-005 | Max 5 failed login attempts → lockout 15 menit | Backend: failed_login_attempts counter<br>locked_until timestamp | Prevent brute force |
| BR-AUTH-006 | Session tracking dengan token hash | Backend: SHA256 hash di user_sessions<br>Revoke on logout | Audit trail |
| BR-AUTH-007 | Role-based dashboard routing | Frontend: router guards<br>Backend: role middleware | ADMIN → /dashboard/admin<br>STAFF → /dashboard/staff |

---

## Technical Implementation

### Backend Architecture

#### Database Schema

**Tables**:
- `users` - User accounts dengan roles, departments, password hash, lockout mechanism
- `user_sessions` - Active sessions dengan token hash, IP tracking, revocation status
- `password_reset_tokens` - Password reset tokens (Sprint 3)
- `activity_logs` - Audit trail untuk LOGIN, LOGOUT, dan actions lainnya
- `notifications` - In-app notifications (Sprint 4)

**Seeded Data**: Admin user dengan NIP `99999`, Password `Admin@123`

> 📄 **Detail Schema**: Lihat `backend/database/setup.sql` untuk complete DDL

#### Models

**User Model** (`backend/models/user.go`):
- Fields: NIP, FullName, Email, PasswordHash, Role, Department, Shift, Status
- Enums: UserRole (7 roles), Department (4 dept), Shift (3 shifts), UserStatus (3 status)
- Methods: `IsLocked()`, `IsActive()`, `HasRole()`, `ToSafeUser()`

**Supporting Models**:
- `UserSession` - Session tracking
- `PasswordResetToken` - Password reset (Sprint 3)
- `ActivityLog` - Audit logging

#### Services

**PasswordService** (`backend/services/password_service.go`):
- `HashPassword()` - Bcrypt dengan cost 12
- `VerifyPassword()` - Constant-time comparison
- `ValidatePasswordPolicy()` - Enforce password rules
- `GetPasswordStrength()` - Strength scoring (0-4)

**AuthService** (`backend/services/auth_service.go`):
- `Login()` - Credentials validation, token generation, session creation
- `Logout()` - Session revocation, activity logging
- `GenerateJWT()` - JWT dengan 15 min expiry
- `GenerateRefreshToken()` - Refresh token dengan 30 days expiry
- `ValidateToken()` - JWT parsing dan validation
- `RefreshAuthToken()` - Token refresh flow

**Security Features**:
- Bcrypt cost 12 untuk password hashing (~200ms)
- JWT HMAC-SHA256 signing
- Token SHA256 hashing untuk database storage
- Rate limiting dengan account lockout
- Activity logging untuk audit trail

#### Middleware

**AuthMiddleware** (`backend/middleware/auth_middleware.go`):
- Extract Bearer token dari Authorization header
- Validate JWT signature dan expiry
- Load user dari database
- Set user ke Gin context
- Return 401 untuk invalid tokens

**RoleMiddleware** (`backend/middleware/role_middleware.go`):
- `RequireRole(roles...)` - Check user has any of specified roles
- `RequireAdmin()` - Shortcut untuk admin-only routes
- `RequireDepartment(depts...)` - Department-based access control

#### API Routes

**Public Routes**:
```
POST   /api/auth/login      # Login dengan NIP & password
POST   /api/auth/refresh    # Refresh JWT token
```

**Protected Routes** (Require JWT):
```
POST   /api/auth/logout     # Logout dan revoke session
GET    /api/auth/me         # Get current user info
```

> 📋 **API Details**: Lihat [Authentication API](../api/authentication.md) untuk complete documentation

---

### Frontend Architecture

#### State Management

**Auth Store** (`frontend/src/stores/auth.js`):
- State: `user`, `token`, `refreshToken`
- Computed: `isAuthenticated`, `userRole`, `userDepartment`, `isAdmin`
- Actions: `setAuth()`, `clearAuth()`, `restoreAuth()`, `hasRole()`
- Persistence: localStorage untuk token & user data

#### Composables

**useAuth** (`frontend/src/composables/useAuth.js`):
- `login(nip, password, rememberMe)` - Login flow dengan error handling
- `logout()` - Logout flow dengan cleanup
- `checkAuth()` - Verify current auth status
- `getDashboardRoute()` - Role-based routing
- `triggerHapticFeedback(type)` - iOS-style vibration

**useApi** (`frontend/src/composables/useApi.js`):
- Axios instance dengan base URL configuration
- Request interceptor: Auto-inject Bearer token
- Response interceptor: Auto-refresh on 401, error handling
- Retry mechanism untuk failed requests

#### Router Guards

**Navigation Guards** (`frontend/src/router/index.js`):
- Guest-only routes (login) → redirect jika sudah authenticated
- Protected routes → redirect ke login jika belum authenticated
- Role-based access control dengan meta.roles
- Redirect back to intended page after login

#### Pages

**Login Page** (`frontend/src/views/auth/Login.vue`):
- Glass effect card dengan backdrop blur
- NIP input (numeric only, max 5 digits)
- Password input dengan show/hide toggle
- Remember me checkbox
- Spring entrance animation, shake animation on error
- Haptic feedback (success/error)

**Dashboards**:
- `AdminDashboard.vue` - Stats, quick actions, activity feed
- `StaffDashboard.vue` - Tasks, performance metrics, notifications

**Profile Page** (`frontend/src/views/profile/Profile.vue`):
- User avatar dengan initial, profile info, badges
- Account info: join date, last login, status

**Navbar** (`frontend/src/components/layout/Navbar.vue`):
- Glass navbar dengan sticky positioning
- User dropdown: avatar, profile link, logout button
- Mobile responsive

#### Design System

**iOS-Inspired Features**:
- Spring physics animations (Motion One)
- Press feedback (scale 0.97 on tap)
- Glass effect (backdrop blur, rgba backgrounds)
- Haptic feedback (Vibration API)
- Gradient theme (Indigo #6366f1 & Fuchsia #d946ef)
- Smooth scrollbar, font smoothing

**Animations** (`frontend/src/style.css`):
- `@keyframes springIn`, `fadeIn`, `shake`, `bounce`
- `.active-scale`, `.glass-card`, `.gradient-bg` utilities

---

## Edge Cases

| Scenario | Handling | Expected Behavior |
|----------|----------|-------------------|
| 5 failed login attempts | Increment counter, set locked_until | Account locked 15 min, error message |
| Token expired during API call | Interceptor catches 401, auto-refresh | Seamless refresh, retry request |
| Refresh token expired | Refresh API returns 401 | Clear auth, redirect to login |
| Multiple concurrent logins | Each creates new session | Multiple sessions tracked |
| Network offline | Axios error, no response | Error message "Tidak dapat terhubung" |
| XSS attempt via input | Numeric validation, parameterized queries | Input sanitized, injection prevented |
| User deleted while logged in | Token validation checks existence | 401 error, auto-logout |

> 🗺️ **User Journeys**: Lihat [Authentication User Journeys](../guides/authentication-user-journeys.md) untuk detailed flow diagrams

---

## Security Considerations

### ✅ Implemented

1. **Password Security**: Bcrypt cost 12, policy enforcement, never returned in responses
2. **Token Security**: JWT HMAC-SHA256, short expiry (15 min), refresh tokens (30 days), SHA256 hashing for storage
3. **Session Security**: IP & user agent tracking, revocation on logout, activity logging
4. **Rate Limiting**: Max 5 failed attempts, 15 min lockout, counter reset on success
5. **Input Validation**: Frontend numeric validation, backend parameterized queries, email format validation
6. **CORS Configuration**: Allowed origins configured, credentials support

### ⚠️ Future Enhancements

- Rate limiting per IP (prevent distributed attacks)
- CSRF protection dengan tokens
- 2FA/MFA untuk admin accounts
- Password history (prevent reuse)
- Session timeout after inactivity
- Device management UI

---

## Configuration

### Backend Environment Variables

**File**: `backend/.env`

```env
# Server
SERVER_PORT=8080
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=sirine_go

# JWT
JWT_SECRET=sirine-go-jwt-secret-key-change-in-production
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=720h

# Security
BCRYPT_COST=12
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=15m

# Frontend
FRONTEND_URL=http://localhost:5173
```

### Frontend Environment Variables

**File**: `frontend/.env`

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_NAME=Sirine Go
VITE_TIMEZONE=Asia/Jakarta
```

---

## Testing

### Unit Tests

**Backend**:
- `backend/tests/unit/services/password_service_test.go` - 8 test cases, 95%+ coverage
- `backend/tests/unit/models/user_test.go` - 7 test cases, 90%+ coverage

**Frontend**:
- `frontend/src/tests/unit/stores/auth.spec.js` - 15 test cases, 95%+ coverage

**Run Commands**:
```bash
# Backend
cd backend && go test -v ./tests/...

# Frontend
cd frontend && yarn test
```

### Manual Testing

> 📋 **Complete Test Plan**: Lihat [AUTH Test Plan](../testing/AUTH-test-plan.md) untuk:
> - Integration tests (4 scenarios)
> - Manual testing checklist (10 cases)
> - Security testing (5 cases)
> - Performance testing (4 benchmarks)
> - Browser compatibility testing
> - Mobile testing (iOS & Android)

**Quick Verification**:
- [ ] Login dengan NIP `99999`, Password `Admin@123`
- [ ] Redirect ke `/dashboard/admin`
- [ ] Refresh page → tetap logged in
- [ ] Logout → redirect ke `/login`
- [ ] 5 failed logins → account locked
- [ ] Token auto-refresh after 15 min

---

## Performance Metrics

| Metric | Target | Notes |
|--------|--------|-------|
| Login Response Time | < 500ms | Including bcrypt (200ms) |
| Dashboard Load Time | < 1s | With stats rendering |
| Token Validation | < 50ms | JWT parsing |
| Frontend Bundle | < 500KB | Gzipped |
| Lighthouse Score | 90+ | Performance & Accessibility |

---

## Dependencies

### Backend
- `github.com/gin-gonic/gin` v1.11.0 - HTTP framework
- `github.com/golang-jwt/jwt/v5` v5.2.0 - JWT
- `golang.org/x/crypto` v0.46.0 - bcrypt
- `gorm.io/gorm` v1.31.1 - ORM

### Frontend
- `vue` ^3.5.24 - Framework
- `pinia` ^2.2.0 - State management
- `vue-router` ^4.4.0 - Routing
- `axios` ^1.13.2 - HTTP client
- `@motionone/vue` ^10.16.4 - Animations
- `tailwindcss` ^4.1.18 - Styling

---

## Known Limitations

1. **Email Service**: Belum diimplementasi (Sprint 3)
2. **Password Change**: Flow belum complete (Sprint 3)
3. **Profile Photo Upload**: Belum ada (Sprint 5)
4. **Real-time Notifications**: Belum implemented (Sprint 4)
5. **Rate Limiting per IP**: Hanya per-user saat ini
6. **Session Management UI**: Belum ada device management view

---

## Related Documentation

- **Test Plan**: [AUTH Test Plan](../testing/AUTH-test-plan.md)
- **User Journeys**: [Authentication User Journeys](../guides/authentication-user-journeys.md)
- **API Documentation**: [Authentication API](../api/authentication.md)
- **Testing Guide**: [TESTING_GUIDE.md](../../TESTING_GUIDE.md)
- **Sprint Implementation**: [SPRINT1_IMPLEMENTATION.md](../../SPRINT1_IMPLEMENTATION.md)
- **Sprint Summary**: [SPRINT1_SUMMARY.md](../../SPRINT1_SUMMARY.md)

---

## Changelog

### Version 1.0.0 (27 Desember 2025)
- ✅ Initial implementation - Sprint 1 complete
- ✅ JWT-based authentication dengan refresh tokens
- ✅ Login/logout flow dengan session management
- ✅ Role-based access control (7 roles)
- ✅ Rate limiting & account lockout (5 attempts, 15 min)
- ✅ iOS-inspired frontend design (glass effect, haptic feedback)
- ✅ Activity logging untuk audit trail
- ✅ Comprehensive testing (unit + integration + manual)

---

**Developer**: Zulfikar Hidayatullah (+62 857-1583-8733)  
**Status**: ✅ Production Ready untuk authentication features  
**Next Sprint**: User Management & Profile (Sprint 2)
