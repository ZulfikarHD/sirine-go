# 🔐 Login Flows

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
