# 🚪 Logout Flow

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
