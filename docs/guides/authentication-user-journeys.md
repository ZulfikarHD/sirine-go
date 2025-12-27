# 🗺️ Authentication System - User Journeys

**Feature**: Authentication System  
**Sprint**: 1  
**Last Updated**: 27 Desember 2025

---

## 📋 Overview

User Journeys merupakan visual documentation untuk user flows dalam authentication system yang bertujuan untuk memahami complete user experience dari berbagai perspectives, yaitu: admin login, staff login, first-time user, password reset, dan error scenarios.

---

## 🎯 Journey Map Index

| Journey ID | Journey Name | User Type | Complexity | Status |
|------------|--------------|-----------|------------|--------|
| J1 | [Happy Path - Admin Login](#journey-1-happy-path---admin-login) | Admin | Simple | ✅ Complete |
| J2 | [Happy Path - Staff Login](#journey-2-happy-path---staff-login) | Staff | Simple | ✅ Complete |
| J3 | [Error - Invalid Credentials](#journey-3-error---invalid-credentials) | Any | Simple | ✅ Complete |
| J4 | [Error - Account Locked](#journey-4-error---account-locked) | Any | Medium | ✅ Complete |
| J5 | [Session Persistence](#journey-5-session-persistence) | Any | Simple | ✅ Complete |
| J6 | [Token Expiry & Refresh](#journey-6-token-expiry--refresh) | Any | Medium | ✅ Complete |
| J7 | [Logout Flow](#journey-7-logout-flow) | Any | Simple | ✅ Complete |
| J8 | [Protected Route Access](#journey-8-protected-route-access) | Any | Simple | ✅ Complete |

---

## Journey 1: Happy Path - Admin Login

**User Type**: Administrator  
**Goal**: Login dan akses admin dashboard  
**Frequency**: Daily  
**Duration**: ~30 seconds

### Flow Diagram

```
📍 START: User opens application
    │
    ├─▶ Browser loads http://localhost:5173
    │   └─ Frontend checks localStorage for token
    │      └─ No token found
    │         └─ Router guard triggers
    │            └─ Auto-redirect to /login
    │
    ├─▶ LOGIN PAGE
    │   ├─ Glass effect card appears (spring animation)
    │   ├─ Form fields visible:
    │   │  ├─ NIP input (placeholder: "Masukkan NIP (5 digit)")
    │   │  ├─ Password input (with show/hide toggle)
    │   │  └─ Remember Me checkbox
    │   └─ "Masuk" button (gradient: indigo → fuchsia)
    │
    ├─▶ USER INPUT
    │   ├─ Types NIP: "99999"
    │   │  └─ Frontend validation: numeric only, max 5 chars
    │   ├─ Types Password: "Admin@123"
    │   │  └─ Password masked by default
    │   └─ Checks "Remember Me" ✓
    │
    ├─▶ SUBMIT FORM
    │   ├─ Click "Masuk" button
    │   │  └─ Button press animation (scale 0.97)
    │   │     └─ Haptic feedback (200ms vibration)
    │   ├─ Loading spinner appears
    │   └─ POST /api/auth/login
    │      ├─ Request body:
    │      │  {
    │      │    "nip": "99999",
    │      │    "password": "Admin@123",
    │      │    "remember_me": true
    │      │  }
    │      └─ Headers: Content-Type: application/json
    │
    ├─▶ BACKEND PROCESSING
    │   ├─ AuthHandler receives request
    │   ├─ AuthService.Login() called
    │   │  ├─ Find user by NIP in database
    │   │  ├─ Check if account locked
    │   │  ├─ Check if account active
    │   │  ├─ Verify password (bcrypt)
    │   │  ├─ Reset failed_login_attempts to 0
    │   │  ├─ Update last_login_at timestamp
    │   │  ├─ Generate JWT token (15 min expiry)
    │   │  ├─ Generate refresh token (30 days)
    │   │  ├─ Hash tokens (SHA256) for storage
    │   │  ├─ Create user_session record
    │   │  └─ Log LOGIN activity
    │   └─ Return response:
    │      {
    │        "success": true,
    │        "message": "Login berhasil",
    │        "data": {
    │          "token": "eyJhbGc...",
    │          "refresh_token": "eyJhbGc...",
    │          "user": {
    │            "id": 1,
    │            "nip": "99999",
    │            "full_name": "Administrator",
    │            "role": "ADMIN",
    │            "department": "KHAZWAL"
    │          },
    │          "require_password_change": false
    │        }
    │      }
    │
    ├─▶ FRONTEND PROCESSING
    │   ├─ Axios receives 200 OK response
    │   ├─ useAuth composable processes data
    │   ├─ Auth store updated:
    │   │  ├─ user = response.data.user
    │   │  ├─ token = response.data.token
    │   │  └─ refreshToken = response.data.refresh_token
    │   ├─ Data persisted to localStorage:
    │   │  ├─ auth_token
    │   │  ├─ refresh_token
    │   │  └─ user_data (JSON)
    │   └─ Success haptic feedback (200ms pulse)
    │
    ├─▶ NAVIGATION
    │   ├─ getDashboardRoute() called
    │   │  └─ Role = ADMIN → returns "/dashboard/admin"
    │   └─ Router.push('/dashboard/admin')
    │
    ├─▶ ADMIN DASHBOARD
    │   ├─ Page loads with fade-in animation
    │   ├─ Navbar appears:
    │   │  ├─ App logo & title
    │   │  ├─ User avatar (initial: "AD")
    │   │  └─ User dropdown (nama: "Administrator")
    │   ├─ Stats cards render (staggered animation):
    │   │  ├─ Total Users: 24
    │   │  ├─ PO Aktif: 12
    │   │  ├─ Produksi Hari Ini: 850/1000
    │   │  └─ QC Pass Rate: 98.5%
    │   ├─ Quick Actions grid
    │   └─ Recent Activity feed
    │
    └─▶ ✅ END: User successfully logged in
        └─ Session active for 30 days (remember me)
```

### Key Touchpoints

| Step | User Action | System Response | Duration |
|------|-------------|-----------------|----------|
| 1 | Opens app | Auto-redirect to login | < 1s |
| 2 | Views login form | Glass card animation | 0.5s |
| 3 | Inputs credentials | Real-time validation | ~10s |
| 4 | Clicks submit | Loading + API call | 0.3s |
| 5 | Backend auth | Token generation | 0.2s |
| 6 | Frontend update | Store + localStorage | < 0.1s |
| 7 | Navigation | Route to dashboard | 0.5s |
| 8 | Dashboard load | Stats render | 0.8s |

**Total Duration**: ~13 seconds

---

## Journey 2: Happy Path - Staff Login

**User Type**: Staff (any non-admin role)  
**Goal**: Login dan akses staff dashboard  
**Frequency**: Daily  
**Duration**: ~30 seconds

### Flow Diagram

```
📍 START: Staff user opens application
    │
    ├─▶ [Same as Journey 1 until LOGIN PAGE]
    │
    ├─▶ USER INPUT
    │   ├─ Types NIP: "12345" (example staff NIP)
    │   ├─ Types Password: "StaffPass123!"
    │   └─ No "Remember Me" (unchecked)
    │
    ├─▶ SUBMIT & BACKEND PROCESSING
    │   └─ [Same as Journey 1]
    │      └─ Response includes:
    │         "role": "OPERATOR_CETAK",
    │         "department": "CETAK",
    │         "shift": "PAGI"
    │
    ├─▶ NAVIGATION
    │   ├─ getDashboardRoute() called
    │   │  └─ Role = OPERATOR_CETAK → returns "/dashboard/staff"
    │   └─ Router.push('/dashboard/staff')
    │
    ├─▶ STAFF DASHBOARD
    │   ├─ Page loads with fade-in animation
    │   ├─ Navbar shows:
    │   │  └─ User info: "Nama Staff - OPERATOR_CETAK"
    │   ├─ Task Cards render:
    │   │  ├─ PO-2025-045 (In Progress)
    │   │  ├─ QC Inspection (Pending)
    │   │  └─ Verifikasi Data (Completed)
    │   ├─ Performance Metrics:
    │   │  ├─ Target Harian: 850/1000 (85%)
    │   │  ├─ Quality Rate: 98.5%
    │   │  └─ Efficiency: 92%
    │   └─ Notifications panel
    │
    └─▶ ✅ END: Staff successfully logged in
        └─ Session active for 15 minutes (no remember me)
```

### Differences from Admin Journey

| Aspect | Admin | Staff |
|--------|-------|-------|
| Dashboard Route | `/dashboard/admin` | `/dashboard/staff` |
| Dashboard Content | Stats, Quick Actions, Activity | Tasks, Performance, Notifications |
| Session Duration | 30 days (remember me) | 15 min (no remember me) |
| Permissions | Full access | Limited to assigned tasks |

---

## Journey 3: Error - Invalid Credentials

**User Type**: Any  
**Goal**: Attempt login dengan wrong credentials  
**Frequency**: Occasional (user error)  
**Duration**: ~15 seconds

### Flow Diagram

```
📍 START: User at login page
    │
    ├─▶ USER INPUT
    │   ├─ Types NIP: "12345"
    │   └─ Types Password: "wrongpassword"
    │
    ├─▶ SUBMIT FORM
    │   └─ POST /api/auth/login
    │
    ├─▶ BACKEND PROCESSING
    │   ├─ Find user by NIP → User found
    │   ├─ Verify password → ❌ FAILED
    │   ├─ Increment failed_login_attempts (1 → 2)
    │   ├─ Save user record
    │   └─ Return 401 Unauthorized:
    │      {
    │        "success": false,
    │        "message": "NIP atau password salah"
    │      }
    │
    ├─▶ FRONTEND ERROR HANDLING
    │   ├─ Axios catches 401 error
    │   ├─ Error message displayed:
    │   │  └─ Red banner: "NIP atau password salah"
    │   ├─ Card shake animation (400ms)
    │   │  └─ translateX: 0 → -10px → 10px → -10px → 0
    │   ├─ Error haptic feedback
    │   │  └─ Vibration: [100ms, 50ms, 100ms]
    │   └─ Form data retained (tidak clear)
    │
    └─▶ ⚠️ END: User sees error, can retry
        └─ Failed attempts: 2/5
```

### User Experience Notes

- ✅ Error message clear dan actionable
- ✅ Form data retained untuk easy correction
- ✅ Visual feedback (shake) grabs attention
- ✅ Haptic feedback untuk mobile users
- ✅ No sensitive info leaked (generic message)

---

## Journey 4: Error - Account Locked

**User Type**: Any  
**Goal**: Login after 5 failed attempts  
**Frequency**: Rare (security incident)  
**Duration**: ~10 seconds + 15 min wait

### Flow Diagram

```
📍 START: User at login (after 4 failed attempts)
    │
    ├─▶ USER INPUT (5th attempt)
    │   ├─ Types NIP: "12345"
    │   └─ Types Password: "stillwrong"
    │
    ├─▶ SUBMIT FORM
    │   └─ POST /api/auth/login
    │
    ├─▶ BACKEND PROCESSING
    │   ├─ Find user by NIP → User found
    │   ├─ Verify password → ❌ FAILED
    │   ├─ Increment failed_login_attempts (4 → 5)
    │   ├─ Check: attempts >= MAX_LOGIN_ATTEMPTS (5)
    │   │  └─ ✅ TRUE → Trigger lockout
    │   ├─ Calculate locked_until:
    │   │  └─ Now + 15 minutes = 14:45:00
    │   ├─ Update user record:
    │   │  ├─ failed_login_attempts = 5
    │   │  └─ locked_until = 2025-12-27 14:45:00
    │   └─ Return 401 Unauthorized:
    │      {
    │        "success": false,
    │        "message": "Terlalu banyak percobaan login gagal, akun Anda dikunci selama 15 menit"
    │      }
    │
    ├─▶ FRONTEND ERROR HANDLING
    │   ├─ Error banner displayed (red):
    │   │  └─ "Akun terkunci hingga 14:45:00"
    │   ├─ Strong shake animation
    │   └─ Long error haptic (500ms)
    │
    ├─▶ SUBSEQUENT ATTEMPTS (with correct password)
    │   ├─ User tries: NIP "12345", Password "CorrectPass123!"
    │   ├─ Backend checks: user.IsLocked() → TRUE
    │   │  └─ locked_until (14:45) > now (14:32)
    │   └─ Return 401:
    │      "Akun Anda terkunci hingga 14:45:00"
    │
    └─▶ ⛔ END: User must wait 15 minutes
        ├─ Option 1: Wait until 14:45
        │  └─ After 14:45, can login normally
        ├─ Option 2: Contact admin
        │  └─ Admin manually unlocks account
        └─ Security: Prevents brute force attacks
```

### Security Flow

```
Attempt 1: ❌ → failed_login_attempts = 1
Attempt 2: ❌ → failed_login_attempts = 2
Attempt 3: ❌ → failed_login_attempts = 3
Attempt 4: ❌ → failed_login_attempts = 4
Attempt 5: ❌ → failed_login_attempts = 5 → 🔒 LOCKED
           └─ locked_until = now + 15 min
```

---

## Journey 5: Session Persistence

**User Type**: Any authenticated user  
**Goal**: Maintain session across page refreshes  
**Frequency**: Very common  
**Duration**: < 1 second

### Flow Diagram

```
📍 START: User at dashboard (logged in)
    │
    ├─▶ USER ACTION
    │   └─ Presses F5 (refresh page)
    │
    ├─▶ BROWSER RELOAD
    │   ├─ Page unloads
    │   ├─ localStorage persists:
    │   │  ├─ auth_token
    │   │  ├─ refresh_token
    │   │  └─ user_data
    │   └─ Page reloads
    │
    ├─▶ APP INITIALIZATION
    │   ├─ main.js executes
    │   ├─ Pinia store created
    │   ├─ Router initialized
    │   └─ App.vue mounted
    │
    ├─▶ AUTH RESTORATION
    │   ├─ App.vue onMounted() hook
    │   ├─ authStore.restoreAuth() called
    │   │  ├─ Read localStorage.getItem('auth_token')
    │   │  ├─ Read localStorage.getItem('refresh_token')
    │   │  ├─ Read localStorage.getItem('user_data')
    │   │  ├─ Parse JSON (with error handling)
    │   │  └─ Set store state:
    │   │     ├─ token = stored token
    │   │     ├─ refreshToken = stored refresh
    │   │     └─ user = parsed user data
    │   └─ isAuthenticated = true
    │
    ├─▶ ROUTER NAVIGATION
    │   ├─ Router tries to load /dashboard
    │   ├─ beforeEach guard checks:
    │   │  └─ authStore.isAuthenticated → TRUE
    │   └─ Allow navigation
    │
    ├─▶ DASHBOARD RENDERS
    │   ├─ Navbar shows user info (from store)
    │   ├─ Dashboard content loads
    │   └─ No API call needed (data from localStorage)
    │
    └─▶ ✅ END: Seamless experience
        └─ User doesn't notice any interruption
```

### Performance

| Metric | Value |
|--------|-------|
| Restore time | < 50ms |
| API calls | 0 (uses cache) |
| User perception | Instant |

---

## Journey 6: Token Expiry & Refresh

**User Type**: Any authenticated user  
**Goal**: Auto-refresh expired token  
**Frequency**: Every 15 minutes (JWT expiry)  
**Duration**: ~500ms (transparent to user)

### Flow Diagram

```
📍 START: User logged in, browsing app (14 min after login)
    │
    ├─▶ USER ACTION
    │   └─ Clicks "Profile" link (15 min after login)
    │
    ├─▶ NAVIGATION
    │   ├─ Router navigates to /profile
    │   └─ Profile.vue component mounted
    │
    ├─▶ API CALL (with expired token)
    │   ├─ GET /api/auth/me
    │   ├─ Headers: Authorization: Bearer <expired-token>
    │   └─ Backend receives request
    │
    ├─▶ BACKEND VALIDATION
    │   ├─ AuthMiddleware.ValidateToken()
    │   ├─ JWT parse → Token expired (15 min passed)
    │   └─ Return 401 Unauthorized:
    │      {
    │        "success": false,
    │        "message": "Token tidak valid atau sudah expired"
    │      }
    │
    ├─▶ AXIOS RESPONSE INTERCEPTOR
    │   ├─ Catches 401 error
    │   ├─ Checks: originalRequest._retry → false
    │   ├─ Checks: authStore.refreshToken → exists
    │   ├─ Sets: originalRequest._retry = true
    │   └─ Initiates refresh flow
    │
    ├─▶ TOKEN REFRESH
    │   ├─ POST /api/auth/refresh
    │   │  └─ Body: { "refresh_token": "eyJhbGc..." }
    │   ├─ Backend processes:
    │   │  ├─ Validate refresh token
    │   │  ├─ Generate new JWT (15 min)
    │   │  ├─ Generate new refresh token (30 days)
    │   │  └─ Update session record
    │   └─ Return 200 OK:
    │      {
    │        "success": true,
    │        "data": {
    │          "token": "new-jwt-token",
    │          "refresh_token": "new-refresh-token",
    │          "user": { ... }
    │        }
    │      }
    │
    ├─▶ FRONTEND UPDATE
    │   ├─ authStore.setAuth(response.data)
    │   ├─ localStorage updated with new tokens
    │   ├─ Original request headers updated:
    │   │  └─ Authorization: Bearer <new-token>
    │   └─ Retry original request
    │
    ├─▶ RETRY ORIGINAL API CALL
    │   ├─ GET /api/auth/me (with new token)
    │   ├─ Backend validates → ✅ Valid
    │   └─ Return user data
    │
    ├─▶ PROFILE PAGE RENDERS
    │   └─ User sees profile data
    │
    └─▶ ✅ END: Seamless token refresh
        └─ User didn't notice any interruption
```

### Timeline

```
T+0:00  → User logs in (token expires at T+15:00)
T+14:59 → User clicks Profile
T+15:00 → API call with expired token
T+15:00 → 401 response
T+15:00 → Auto refresh token (200ms)
T+15:00 → Retry API call with new token
T+15:00 → Profile loads
```

**User Experience**: Completely transparent, no interruption

---

## Journey 7: Logout Flow

**User Type**: Any authenticated user  
**Goal**: Securely logout dari aplikasi  
**Frequency**: Daily (end of session)  
**Duration**: ~2 seconds

### Flow Diagram

```
📍 START: User at dashboard (logged in)
    │
    ├─▶ USER ACTION
    │   ├─ Clicks user avatar (top right navbar)
    │   └─ Dropdown menu appears (slide-down animation)
    │
    ├─▶ DROPDOWN MENU
    │   ├─ User info displayed:
    │   │  ├─ Avatar with initial
    │   │  ├─ Full name
    │   │  └─ Role
    │   ├─ Menu items:
    │   │  ├─ 👤 Profile
    │   │  └─ 🚪 Keluar (red text)
    │   └─ User clicks "Keluar"
    │
    ├─▶ FRONTEND LOGOUT
    │   ├─ useAuth.logout() called
    │   ├─ Loading state set
    │   └─ POST /api/auth/logout
    │      └─ Headers: Authorization: Bearer <token>
    │
    ├─▶ BACKEND PROCESSING
    │   ├─ AuthHandler.Logout()
    │   ├─ Extract user_id from context
    │   ├─ Hash token (SHA256)
    │   ├─ Find session in database
    │   ├─ Update session:
    │   │  └─ is_revoked = true
    │   ├─ Log LOGOUT activity
    │   └─ Return 200 OK:
    │      {
    │        "success": true,
    │        "message": "Logout berhasil"
    │      }
    │
    ├─▶ FRONTEND CLEANUP
    │   ├─ authStore.clearAuth() called
    │   │  ├─ user = null
    │   │  ├─ token = null
    │   │  └─ refreshToken = null
    │   ├─ localStorage.removeItem('auth_token')
    │   ├─ localStorage.removeItem('refresh_token')
    │   ├─ localStorage.removeItem('user_data')
    │   └─ Success haptic feedback (200ms)
    │
    ├─▶ NAVIGATION
    │   └─ router.push('/login')
    │
    ├─▶ LOGIN PAGE
    │   ├─ Page loads with fade-in
    │   ├─ Form empty (no pre-filled data)
    │   └─ Success message (optional):
    │      "Anda telah logout. Silakan login kembali."
    │
    └─▶ ✅ END: User logged out securely
        ├─ Session revoked di database
        ├─ Token tidak bisa digunakan lagi
        └─ localStorage cleared
```

### Security Checklist

- ✅ Session revoked di backend (is_revoked = true)
- ✅ Token hash removed dari active sessions
- ✅ localStorage cleared completely
- ✅ Pinia store reset
- ✅ Activity logged untuk audit
- ✅ User redirected to login
- ✅ Old token cannot be reused

---

## Journey 8: Protected Route Access

**User Type**: Unauthenticated visitor  
**Goal**: Try access protected route without login  
**Frequency**: Common (new visitors)  
**Duration**: < 1 second

### Flow Diagram

```
📍 START: User types URL directly
    │
    ├─▶ USER ACTION
    │   └─ Types: http://localhost:5173/dashboard
    │      └─ Presses Enter
    │
    ├─▶ ROUTER NAVIGATION
    │   ├─ Router tries to navigate to /dashboard
    │   └─ beforeEach guard executes
    │
    ├─▶ AUTHENTICATION CHECK
    │   ├─ Check: authStore.user → null
    │   ├─ Check: authStore.token → null
    │   ├─ Attempt: authStore.restoreAuth()
    │   │  └─ localStorage.getItem('auth_token') → null
    │   └─ Result: isAuthenticated = false
    │
    ├─▶ ROUTE PROTECTION
    │   ├─ Check: route.meta.requiresAuth → true
    │   ├─ Check: isAuthenticated → false
    │   └─ Guard blocks navigation
    │
    ├─▶ REDIRECT
    │   ├─ next({ path: '/login', query: { redirect: '/dashboard' } })
    │   └─ Browser URL changes to:
    │      http://localhost:5173/login?redirect=/dashboard
    │
    ├─▶ LOGIN PAGE
    │   ├─ Page loads
    │   ├─ Form displayed
    │   └─ Redirect param stored
    │
    ├─▶ [User logs in successfully]
    │   └─ After login:
    │      ├─ Check: route.query.redirect → '/dashboard'
    │      └─ router.push('/dashboard')
    │
    └─▶ ✅ END: User redirected back to intended page
        └─ Seamless experience after login
```

### URL Flow

```
User types:     /dashboard
Router checks:  Not authenticated
Redirects to:   /login?redirect=/dashboard
User logs in:   Success
Router checks:  redirect param exists
Navigates to:   /dashboard (original destination)
```

---

## 🎯 Journey Success Metrics

| Journey | Success Rate Target | Current | Status |
|---------|---------------------|---------|--------|
| J1: Admin Login | > 95% | - | 🔄 Pending |
| J2: Staff Login | > 95% | - | 🔄 Pending |
| J3: Invalid Credentials | 100% (error shown) | - | 🔄 Pending |
| J4: Account Locked | 100% (locked) | - | 🔄 Pending |
| J5: Session Persistence | > 99% | - | 🔄 Pending |
| J6: Token Refresh | > 98% | - | 🔄 Pending |
| J7: Logout | > 99% | - | 🔄 Pending |
| J8: Protected Route | 100% (redirected) | - | 🔄 Pending |

---

## 📊 User Experience Metrics

| Metric | Target | Current | Notes |
|--------|--------|---------|-------|
| Login Time (P50) | < 3s | - | From form submit to dashboard |
| Login Time (P95) | < 5s | - | Including slow networks |
| Session Restore | < 100ms | - | Page refresh experience |
| Token Refresh | < 500ms | - | Transparent to user |
| Logout Time | < 2s | - | Complete cleanup |
| Error Recovery | < 5s | - | User can retry immediately |

---

## 🔗 Related Documentation

- **Feature Documentation**: [Authentication System](../features/AUTHENTICATION.md)
- **Test Plan**: [AUTH Test Plan](../testing/AUTH-test-plan.md)
- **API Documentation**: [Authentication API](../api/authentication.md)
- **Testing Guide**: [TESTING_GUIDE.md](../../TESTING_GUIDE.md)

---

**Last Updated**: 27 Desember 2025  
**Version**: 1.0.0 - Sprint 1  
**Status**: ✅ User Journeys Complete
