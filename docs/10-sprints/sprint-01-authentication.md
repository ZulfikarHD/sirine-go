# 🔐 Sprint 1: Foundation & Authentication

**Version:** 1.1.0  
**Date:** 27 Desember 2025  
**Duration:** 1 week  
**Status:** ✅ Completed

## 📋 Sprint Goals

Membangun foundation sistem dan implementasi complete authentication system dengan JWT dan role-based access control.

---

## ✨ Features Implemented

### 1. Authentication System

#### JWT-Based Authentication
- Access token dengan 15-minute expiry
- Refresh token dengan 30-day expiry
- Auto token refresh on 401 Unauthorized
- Token persistence di localStorage
- Secure token storage dengan SHA256 hashing

#### Login/Logout Flow
- Login dengan NIP & password
- Session tracking dengan IP address dan User Agent
- Logout dengan session revocation
- Auto-redirect to login jika unauthorized

#### Security Features
- Password hashing dengan bcrypt cost 12
- Rate limiting (5 failed attempts → 15 min lockout)
- Account lockout mechanism
- Input validation (frontend & backend)

### 2. Authorization System

#### Role-Based Access Control (RBAC)
**7 Roles:**
- `ADMIN` - Super Admin
- `MANAGER_KHAZWAL` - Manager Khazwal
- `STAFF_KHAZWAL` - Staff Khazwal
- `MANAGER_KEUANGAN` - Manager Keuangan
- `STAFF_KEUANGAN` - Staff Keuangan
- `MANAGER_DISTRIBUSI` - Manager Distribusi
- `STAFF_DISTRIBUSI` - Staff Distribusi

**4 Departments:**
- KHAZWAL
- KEUANGAN
- DISTRIBUSI
- ADMIN

**3 Shifts:**
- PAGI
- SORE
- MALAM

### 3. Backend Architecture

#### Service Pattern
```
backend/
├── services/
│   ├── auth_service.go
│   ├── password_service.go
│   └── user_service.go
├── handlers/
│   ├── auth_handler.go
│   └── user_handler.go
└── middleware/
    ├── auth_middleware.go
    ├── role_middleware.go
    └── activity_logger.go
```

#### Middleware System
- **AuthMiddleware** - JWT validation
- **RoleMiddleware** - RBAC enforcement
- **ActivityLogger** - Audit trail
- **CORS** - Cross-origin requests

### 4. Database Schema

#### Users Table
```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nip VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    phone VARCHAR(20),
    role ENUM('ADMIN', 'MANAGER_KHAZWAL', 'STAFF_KHAZWAL', ...),
    department ENUM('ADMIN', 'KHAZWAL', 'KEUANGAN', 'DISTRIBUSI'),
    shift ENUM('PAGI', 'SORE', 'MALAM'),
    is_active BOOLEAN DEFAULT TRUE,
    require_password_change BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### User Sessions Table
```sql
CREATE TABLE user_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

---

## 🎨 Frontend Implementation

### State Management (Pinia)

```javascript
// stores/auth.js
export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    accessToken: null,
    isAuthenticated: false
  }),
  
  actions: {
    async login(nip, password) { ... },
    async logout() { ... },
    async refreshToken() { ... },
    async fetchCurrentUser() { ... }
  }
})
```

### Router Guards

```javascript
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.meta.roles && !to.meta.roles.includes(authStore.user?.role)) {
    next('/unauthorized')
  } else {
    next()
  }
})
```

### Components

1. **Login Page** - Form login dengan validation
2. **Dashboard** - Role-based dashboards (Admin/Staff)
3. **Navbar** - User dropdown dengan profile, logout
4. **Sidebar** - Role-based menu visibility

---

## 🔌 API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/auth/login` | Login user | No |
| POST | `/api/auth/logout` | Logout user | Yes |
| GET | `/api/auth/me` | Get current user | Yes |
| POST | `/api/auth/refresh` | Refresh token | Yes |
| GET | `/health` | Health check | No |

---

## 🎨 Design System

### Apple-Inspired Design
- **Colors:** Indigo (#6366f1) & Fuchsia (#d946ef) gradient
- **Glass Effect:** Cards dengan backdrop blur
- **Spring Animations:** Motion-V dengan physics-based animations
- **iOS-style Feedback:** Scale 0.97 on active
- **Haptic Feedback:** For mobile devices

### Animation Presets

```javascript
// Spring presets (iOS-like)
export const springPresets = {
  default: { type: 'spring', stiffness: 400, damping: 30, mass: 0.8 },
  snappy: { type: 'spring', stiffness: 500, damping: 35, mass: 0.6 },
  gentle: { type: 'spring', stiffness: 300, damping: 25, mass: 1 }
}
```

---

## 📦 Tech Stack

### Backend
- **Framework:** Gin (Go web framework)
- **ORM:** GORM
- **Database:** MySQL 8.0+
- **Authentication:** JWT (golang-jwt/jwt)
- **Password:** bcrypt

### Frontend
- **Framework:** Vue 3 (Composition API)
- **State:** Pinia
- **Router:** Vue Router 4
- **Styling:** Tailwind CSS 4
- **Animations:** Motion-V
- **Icons:** Lucide Vue Next
- **Build:** Vite

---

## 🧪 Testing

### Test Scenarios

✅ **Login Flow**
- Login dengan credentials valid
- Login dengan credentials invalid
- Auto-redirect ke dashboard setelah login
- Session persistence setelah refresh

✅ **Logout Flow**
- Logout menghapus session
- Redirect ke login page
- Token revoked di backend

✅ **Protected Routes**
- Unauthorized access redirect ke login
- Role-based access control works
- Token refresh otomatis

✅ **Rate Limiting**
- 5 failed attempts → lockout
- Lockout duration 15 minutes
- Successful login resets counter

---

## 🔐 Security Measures

1. **Password Security**
   - Bcrypt hashing dengan cost 12
   - Strong password policy enforced

2. **Token Security**
   - JWT signed dengan HMAC-SHA256
   - Short-lived access tokens (15 min)
   - Refresh token rotation
   - Token hash storage di database

3. **Session Security**
   - IP address tracking
   - User agent tracking
   - Session revocation on logout
   - Activity logging

4. **Input Validation**
   - Frontend validation dengan Zod/Vuelidate
   - Backend validation dengan custom validators
   - SQL injection protection (GORM parameterized queries)

---

## 📊 Sprint Metrics

### Development Stats
- **Backend Files:** 12 files created
- **Frontend Files:** 15 files created
- **API Endpoints:** 5 endpoints
- **Database Tables:** 3 tables
- **Components:** 4 major components
- **Test Scenarios:** 10+ scenarios

### Code Quality
- **Backend:** Clean architecture dengan service pattern
- **Frontend:** Composition API dengan composables
- **Type Safety:** TypeScript untuk frontend
- **Error Handling:** Comprehensive error handling

---

## 🔄 Lessons Learned

### What Went Well ✅
- Service pattern memudahkan maintenance
- JWT + refresh token approach solid
- Motion-V animations smooth & performant
- RBAC implementation flexible

### Challenges 🎯
- Token refresh edge cases membutuhkan careful handling
- Rate limiting implementation butuh testing menyeluruh
- Mobile responsiveness butuh extra attention

### Improvements for Next Sprint 💡
- Add more comprehensive unit tests
- Implement E2E testing
- Add logging middleware untuk debugging
- Consider Redis untuk session storage (scalability)

---

## 📚 Documentation

### Files Created
- `docs/features/AUTHENTICATION.md` - Feature documentation
- `SPRINT1_IMPLEMENTATION.md` - Implementation guide
- API documentation updated
- README.md updated

---

## 🎯 Next Steps (Sprint 2)

1. **User Management System**
   - CRUD operations untuk users
   - Search & filtering
   - Pagination
   - Bulk operations

2. **Profile Management**
   - View profile
   - Edit profile (self-service)
   - Change password

3. **Enhanced UI/UX**
   - User list table
   - User form modal
   - Role badges
   - Advanced filters

---

## 📞 Support

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733  
**Sprint Lead:** Zulfikar Hidayatullah

---

## 🔗 Related Documentation

- [Architecture Overview](../02-architecture/overview.md)
- [API Documentation](../03-development/api-documentation.md)
- [Authentication Guide](../05-guides/authentication/README.md)
- [Security Guide](../05-guides/security.md)

---

**Sprint Status:** ✅ Completed  
**Next Sprint:** [Sprint 2: User Management](./sprint-02-user-management.md)  
**Last Updated:** 29 Desember 2025
