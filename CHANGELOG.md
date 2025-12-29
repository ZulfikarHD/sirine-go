# 📋 Changelog

Semua perubahan signifikan pada project ini akan didokumentasikan dalam file ini.

Format berdasarkan [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), dan project ini mengikuti [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

---

## [1.5.0] - Sprint 5: Enhancements & Gamification - 2025-12-28

### ✨ Added

#### Gamification System
- **Achievement System** dengan 6 initial achievements (First Login, Week Streak, Month Streak, Profile Complete, Early Bird, Night Owl)
- **Points Tracking** dengan level system (Bronze, Silver, Gold, Platinum)
- **Level Calculation** berdasarkan total points (Bronze: 0-99, Silver: 100-499, Gold: 500-999, Platinum: 1000+)
- **Achievement Badges** dengan locked/unlocked states dan unlock animations
- **Points Display** dengan animated counter dan progress bar ke next level
- **User Stats API** untuk comprehensive gamification statistics

#### Profile Enhancements
- **Profile Photo Upload** dengan auto-resize ke 200x200px
- **Image Optimization** dengan JPEG quality 90%
- **Format Validation** (JPG, PNG, WebP dengan max 5MB)
- **Photo Preview** sebelum upload
- **Delete Photo** functionality dengan confirmation
- **Static File Serving** untuk uploaded photos

#### Bulk Operations
- **CSV Import** untuk bulk user creation dengan validation per row
- **CSV Export** dengan filters (role, department, status)
- **Import Result** dengan success/error breakdown
- **Error Tracking** untuk failed imports dengan row details
- **CsvImport Component** dengan drag & drop dan format guide

#### UX Enhancements
- **Haptic Feedback** composable dengan multiple patterns (light, medium, heavy, success, error, achievement)
- **Loading Skeletons** untuk table, card, list, grid, dan profile
- **Context-specific Haptic** untuk button press, login success, form submit, toggles
- **Celebratory Patterns** untuk achievement unlocks
- **Smooth Animations** dengan number tween untuk points counter

#### New Components
- `AchievementBadge.vue` - Visual achievement badge dengan icon dan unlock status
- `PointsDisplay.vue` - Animated points counter dengan level progress
- `PhotoUpload.vue` - Drag & drop photo upload dengan preview
- `CsvImport.vue` - CSV import dengan validation dan result display
- `LoadingSkeleton.vue` - Multiple skeleton types untuk loading states

#### Backend Services
- `AchievementService` - Achievement management dan auto-check logic
- `FileService` - Photo upload dengan auto-resize dan validation
- CSV methods di `UserService` (BulkImportUsersFromCSV, ExportUsersToCSV)

#### Database Schema
- `achievements` table (id, code, name, description, icon, points, category, criteria)
- `user_achievements` table (user_id, achievement_id, unlocked_at)
- Added `total_points` dan `level` columns ke `users` table

#### API Endpoints
- `GET /api/achievements` - List all achievements
- `GET /api/profile/achievements` - Get user achievements dengan unlock status
- `GET /api/profile/stats` - Get user gamification stats
- `POST /api/admin/achievements/award` - Award achievement (Admin only)
- `GET /api/admin/users/:id/achievements` - Get user achievements (Admin)
- `POST /api/profile/photo` - Upload profile photo
- `DELETE /api/profile/photo` - Delete profile photo
- `POST /api/users/import` - Import users from CSV (Admin)
- `GET /api/users/export` - Export users to CSV (Admin/Manager)

### 🎨 Changed
- Enhanced Profile page dengan points display dan recent achievements
- Updated Achievements page dengan category filtering
- Enhanced UserList page dengan CSV import/export buttons
- Updated auth store dengan updateUserField() method
- Improved haptic feedback integration across app

### 🔧 Technical
- Added image resize library (github.com/nfnt/resize)
- Implemented auto-resize untuk profile photos (200x200px)
- CSV parsing dengan row-by-row validation
- Transaction support untuk atomic updates
- Optimistic rendering untuk instant UI feedback

---

## [1.4.0] - Sprint 4: Notifications & Audit - 2025-12-28

### ✨ Added

#### Notification System
- **In-App Notifications** dengan real-time updates via polling (30s interval)
- **Notification Bell** component dengan animated badge count
- **Notification Center** page dengan tab switching (Semua / Belum Dibaca)
- **Notification Dropdown** dengan recent 5 notifications preview
- **Mark as Read** functionality (single dan bulk)
- **Delete Notification** functionality
- **Optimistic Updates** untuk instant UI feedback
- **Notification Store** (Pinia) dengan polling management
- **Notification Types** dengan color-coded icons (SUCCESS, INFO, WARNING, ERROR)

#### Audit & Activity Logs
- **Activity Logs Viewer** untuk Admin/Manager dengan comprehensive filters
- **ActivityLogTable** component dengan expandable rows
- **JsonDiff** component untuk before/after comparison
- **Activity Stats** endpoint dengan breakdown by action
- **User Activity** tracking dengan IP address dan User Agent
- **Change Tracking** dengan JSON format untuk before/after values
- **Filter Support** (action, entity_type, user, date range, search)
- **Pagination** untuk activity logs dengan configurable page size

#### New Components
- `NotificationBell.vue` - Bell icon dengan badge dan dropdown
- `ActivityLogTable.vue` - Responsive table dengan expandable rows
- `JsonDiff.vue` - Side-by-side before/after comparison dengan field-by-field diff

#### Backend Services
- `NotificationService` - Complete notification management
- `ActivityLogService` - Activity log queries dengan filters dan pagination

#### API Endpoints
- `GET /api/notifications` - List user notifications
- `GET /api/notifications/unread-count` - Get unread badge count
- `GET /api/notifications/recent` - Get recent notifications
- `PUT /api/notifications/:id/read` - Mark notification as read
- `PUT /api/notifications/read-all` - Mark all as read
- `DELETE /api/notifications/:id` - Delete notification
- `POST /api/notifications` - Create notification (testing/admin)
- `GET /api/admin/activity-logs` - List logs dengan filters
- `GET /api/admin/activity-logs/:id` - Get log detail
- `GET /api/admin/activity-logs/user/:id` - Get user activity
- `GET /api/admin/activity-logs/stats` - Get activity statistics
- `GET /api/profile/activity` - Get own activity logs

### 🎨 Changed
- Updated Navbar dengan NotificationBell integration
- Updated Sidebar dengan "Notifikasi" dan "Audit Logs" menu items
- Enhanced activity logging middleware dengan background execution

### 🔧 Technical
- Real-time notification updates via 30-second polling
- Optimistic updates dengan automatic rollback on errors
- Background goroutines untuk activity logging (non-blocking)
- Preloading relationships di GORM untuk optimized queries
- Debounced search (300ms) untuk efficient API calls

---

## [1.3.0] - Sprint 3: Password Management & Security - 2025-12-28

### ✨ Added

#### Password Management
- **Change Password** flow untuk self-service password update
- **Forgot Password** flow dengan email reset link
- **Reset Password** flow dengan secure token validation
- **Force Change Password** modal untuk first-time login users
- **Password Strength Indicator** dengan real-time validation dan requirements checklist
- **Session Expiration** modal untuk token expired scenarios

#### Security Features
- **Password Policy Enforcement** (min 8 chars, uppercase, number, special char)
- **Token Security** dengan SHA256 hashing dan 1-hour expiry
- **Single-Use Tokens** dengan used_at tracking
- **Session Revocation** setelah password change
- **Email Security** dengan prevention untuk email enumeration

#### New Components
- `PasswordStrength.vue` - Visual progress bar dengan color coding dan requirements checklist
- `SessionExpired.vue` - Fullscreen blocking modal untuk expired sessions

#### New Pages
- `ForgotPassword.vue` - Request reset link dengan NIP atau Email
- `ResetPassword.vue` - Reset password dengan token dari email
- `ChangePassword.vue` - Self-service password change
- `ForceChangePassword.vue` - Fullscreen modal untuk mandatory password change

#### Backend Services
- Enhanced `PasswordService` dengan methods:
  - `ChangePassword()` - Change password dengan current password validation
  - `GenerateResetToken()` - Generate secure 32-byte token
  - `ResetPassword()` - Reset dengan token validation
  - `SendResetEmail()` - Send email dengan reset link
  - `RequestPasswordReset()` - Combined flow untuk forgot password
  - `ValidatePasswordPolicy()` - Enforce password requirements
  - `GetPasswordStrength()` - Calculate strength score (0-4)

#### API Endpoints
- `PUT /api/profile/password` - Change own password
- `POST /api/auth/forgot-password` - Request reset link
- `POST /api/auth/reset-password` - Reset dengan token
- `POST /api/users/:id/reset-password` - Admin force reset

### 🎨 Changed
- Updated Login page dengan "Lupa Password?" link
- Enhanced Navbar dengan "Ganti Password" menu item
- Updated Router dengan navigation guards untuk force password check
- Enhanced auth response dengan `require_password_change` flag

### 🔧 Technical
- Bcrypt cost 12 untuk password hashing
- SHA256 untuk token storage
- Auto-invalidate old tokens saat generate new
- Force re-login setelah password change

### 🔐 Security
- Password complexity requirements enforced
- Token expiration (1 hour)
- Session revocation on password change
- Email enumeration prevention

---

## [1.2.0] - Sprint 2: User Management & Profile - 2025-12-28

### ✨ Added

#### User Management (Admin)
- **User List Page** dengan search, filters, dan pagination
- **Create User** modal dengan auto-generated password (12 chars, strong policy)
- **Edit User** functionality untuk Admin/Manager
- **Delete User** (soft delete) dengan confirmation
- **Bulk Operations** (bulk delete, bulk update status)
- **Search Users** by NIP atau Full Name
- **Filter Users** by Role, Department, Status
- **Pagination** dengan configurable page size (default: 20)
- **User Form Modal** dengan slide-up animation (iOS-style)
- **Password Display** dengan copy button untuk generated passwords

#### Profile Management (Self-Service)
- **Profile View** page dengan complete user information
- **Edit Profile** page untuk self-service updates
- **Profile Update** untuk Full Name, Email, Phone
- **Read-Only Fields** untuk NIP, Role, Department, Shift (greyed out)
- **Success Flow** dengan auto-redirect dan success message

#### Components
- `RoleBadge.vue` - Color-coded role badges dengan icons
- `UserFormModal.vue` - Create/Edit user modal dengan validation

#### Backend Services
- `UserService` dengan comprehensive methods:
  - `GetAllUsers()` - List dengan filters dan pagination
  - `GetUserByID()`, `GetUserByNIP()`, `GetUserByEmail()`
  - `CreateUser()` - Auto-generate secure password
  - `UpdateUser()` - Admin user update
  - `DeleteUser()` - Soft delete
  - `SearchUsers()` - Search by NIP/name/email
  - `UpdateProfile()` - Self-service profile update
  - `BulkDeleteUsers()`, `BulkUpdateStatus()`

#### Middleware
- `ActivityLogger` - Auto-log CREATE/UPDATE/DELETE operations dengan before/after changes

#### API Endpoints
- `GET /api/users` - List users dengan filters
- `GET /api/users/search` - Search users
- `GET /api/users/:id` - Get user detail
- `POST /api/users` - Create user (Admin only)
- `PUT /api/users/:id` - Update user (Admin only)
- `DELETE /api/users/:id` - Delete user (Admin only)
- `POST /api/users/bulk-delete` - Bulk delete (Admin only)
- `POST /api/users/bulk-update-status` - Bulk update status (Admin only)
- `GET /api/profile` - Get own profile
- `PUT /api/profile` - Update own profile

### 🎨 Changed
- Updated Sidebar dengan "Manajemen User" menu (Admin/Manager only)
- Updated auth store dengan `fetchCurrentUser()` method
- Enhanced Profile page dengan navigation ke EditProfile

### 🔧 Technical
- Auto-generated password: 12 chars dengan strong policy (1 uppercase, 1 lowercase, 1 digit, 1 special char)
- Debounced search: 300ms delay untuk optimize API calls
- Activity logging via middleware dengan background goroutine
- Role-based authorization dengan RequireRole middleware

---

## [1.1.0] - Sprint 1: Foundation & Authentication - 2025-12-27

### ✨ Added

#### Authentication System
- **JWT-based Authentication** dengan 15-minute access token expiry
- **Refresh Token** mechanism dengan 30-day expiry
- **Login/Logout** flow dengan session tracking
- **Auto Token Refresh** on 401 Unauthorized
- **Token Persistence** di localStorage
- **Session Tracking** dengan IP address dan User Agent
- **Rate Limiting** (5 failed attempts → 15 min lockout)
- **Account Lockout** mechanism untuk security

#### Authorization
- **Role-Based Access Control (RBAC)** dengan 7 roles:
  - ADMIN (Super Admin)
  - MANAGER_KHAZWAL (Manager Khazwal)
  - STAFF_KHAZWAL (Staff Khazwal)
  - MANAGER_KEUANGAN (Manager Keuangan)
  - STAFF_KEUANGAN (Staff Keuangan)
  - MANAGER_DISTRIBUSI (Manager Distribusi)
  - STAFF_DISTRIBUSI (Staff Distribusi)
- **4 Departments**: KHAZWAL, KEUANGAN, DISTRIBUSI, ADMIN
- **3 Shifts**: PAGI, SORE, MALAM

#### Backend Architecture
- **Service Pattern** untuk business logic separation
- **Middleware System** (Auth, Role, Activity Logger, Rate Limiting)
- **Error Handling** dengan descriptive messages dalam Bahasa Indonesia
- **GORM** untuk type-safe database operations
- **Auto Migration** untuk database schema

#### Database Schema
- `users` table dengan 7 roles, 4 departments, 3 shifts
- `user_sessions` table untuk token tracking
- `password_reset_tokens` table untuk password reset flow
- `activity_logs` table untuk audit trail

#### Backend Services
- `AuthService` - Authentication logic (login, logout, token management)
- `PasswordService` - Password hashing dengan bcrypt cost 12
- `SessionService` - Session tracking dan revocation

#### Frontend Architecture
- **Pinia Store** untuk state management (auth store)
- **Vue Router** dengan navigation guards
- **Composables** untuk reusable logic (useAuth, useApi)
- **Auto Token Injection** untuk authenticated requests

#### UI Components & Pages
- **Login Page** dengan iOS-inspired design
- **Admin Dashboard** (role: ADMIN, MANAGER)
- **Staff Dashboard** (role: STAFF)
- **Profile Page** dengan user information
- **Navbar** dengan user dropdown (profile, change password, logout)
- **Sidebar** dengan role-based menu visibility

#### Design System
- **Apple-inspired Design** dengan Indigo (#6366f1) & Fuchsia (#d946ef) gradient
- **Glass Effect Cards** dengan backdrop blur
- **Spring Physics Animations** dengan Motion-V (bukan CSS @keyframes)
- **iOS-style Press Feedback** (scale 0.97 on active)
- **Haptic Feedback** untuk mobile devices
- **Mobile-First** responsive design

#### API Endpoints
- `POST /api/auth/login` - Login dengan NIP & password
- `POST /api/auth/logout` - Logout dan revoke session
- `GET /api/auth/me` - Get current user info
- `POST /api/auth/refresh` - Refresh JWT token
- `GET /health` - Health check

### 🔐 Security
- **Password Hashing** dengan bcrypt cost 12
- **JWT Security** dengan HMAC-SHA256 signing
- **Token Hash Storage** dengan SHA256 di database
- **Session Revocation** on logout
- **Activity Logging** untuk audit trail
- **Input Validation** (frontend & backend)
- **SQL Injection Protection** dengan GORM parameterized queries

### 🎨 Design
- Motion-V animations dengan spring presets (stiffness: 400-500, damping: 25-35)
- Staggered list animations (0.05s delay per item)
- Entrance animations (fade-up, fade-scale)
- Modal animations dengan spring physics
- Active-scale class untuk button feedback
- Backdrop blur hanya pada navbar (performance optimization)

### 📦 Dependencies
#### Backend
- `github.com/gin-gonic/gin` - Web framework
- `gorm.io/gorm` - ORM
- `gorm.io/driver/mysql` - MySQL driver
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `golang.org/x/crypto` - Bcrypt hashing
- `github.com/joho/godotenv` - Environment variables

#### Frontend
- `vue` ^3.5.13 - Frontend framework
- `pinia` ^2.2.0 - State management
- `vue-router` ^4.4.0 - Routing
- `motion-v` - Animations (iOS-style)
- `tailwindcss` ^4.1.18 - Styling
- `vite` ^7.2.4 - Build tool
- `lucide-vue-next` - Icons

### 🧪 Testing
- Manual testing completed untuk all acceptance criteria
- Login/logout flow tested
- Protected routes tested
- Rate limiting tested
- Session persistence tested

### 📝 Documentation
- Complete feature documentation (`docs/features/AUTHENTICATION.md`)
- Sprint implementation guide (`SPRINT1_IMPLEMENTATION.md`)
- Updated README.md dengan Sprint 1 status
- Environment example files (`.env.example`)

---

## [1.0.0] - Initial Project Setup - 2025-12-20

### ✨ Added
- Initial project structure (backend + frontend)
- Database setup script (`backend/database/setup.sql`)
- Development environment configuration
- Makefile untuk development commands
- Documentation structure di `docs/` folder
- Git repository initialization

### 🔧 Technical
- Go module initialization
- Yarn workspace setup
- MySQL database creation
- Environment variables setup
- CORS configuration

---

## Legend

- ✨ **Added** - Fitur baru
- 🎨 **Changed** - Perubahan pada fitur existing
- 🐛 **Fixed** - Bug fixes
- 🗑️ **Deprecated** - Features yang akan dihapus
- 🔒 **Removed** - Features yang dihapus
- 🔐 **Security** - Security improvements
- 🔧 **Technical** - Technical changes/improvements

---

## Version Format

Format: `[MAJOR.MINOR.PATCH]`

- **MAJOR** - Breaking changes
- **MINOR** - New features (backwards compatible)
- **PATCH** - Bug fixes (backwards compatible)

---

**Last Updated**: 28 Desember 2025  
**Current Version**: 1.5.0  
**Status**: Sprint 5 Complete - Ready for Sprint 6 (Testing & Deployment)
