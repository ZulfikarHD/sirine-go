# 🔄 Session Management

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
