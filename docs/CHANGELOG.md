# 📝 Changelog

Semua perubahan signifikan pada project Sirine Go akan didokumentasikan di file ini.

Format mengikuti [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), dan versioning mengikuti [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- **Sprint 3**: Password Management & Security
- **Sprint 4**: Notifications & Audit Logs
- **Sprint 5**: Enhancements & Gamification
- **Sprint 6**: Testing & Production Deployment

---

## [1.1.0] - 2025-12-28

### Added
- **User Management System** (Sprint 2)
  - Admin panel untuk CRUD users dengan role-based access
  - Search dan filter users (by role, department, status, name/NIP)
  - Bulk operations (delete, update status)
  - Auto-generated passwords (12 characters dengan strong policy)
  - Pagination (20 users per page)
  - Activity logging untuk audit trail
- **Profile Management** (Sprint 2)
  - Self-service profile view dan edit
  - User dapat update: full_name, email, phone
  - Read-only fields: NIP, role, department, shift
  - Form validation dengan error handling
- **Activity Logger Middleware**
  - Auto-log CREATE/UPDATE/DELETE operations
  - Before/after changes dalam JSON format
  - IP address dan user agent tracking
  - Background processing untuk performance
- **Motion-v Animations**
  - iOS-style spring physics animations
  - Staggered table row animations
  - Modal slide-up dengan spring bounce
  - Page entrance animations
- **Frontend Components**
  - UserList page dengan search/filter/pagination
  - UserFormModal untuk create/edit users
  - RoleBadge component dengan color coding
  - EditProfile page untuk self-service

### Fixed
- **CORS Issues**: Changed from raw axios to configured apiClient
- **API Base URL**: Properly configured with `/api` prefix
- **Authentication Flow**: Fixed endpoint paths untuk consistency
- **Temporal Dead Zone**: Fixed function hoisting issues
- **Linter Warnings**: Fixed gradient class naming

### Changed
- **API Client**: All stores now use configured apiClient dengan auto auth headers
- **User Store**: Implemented proper state management dengan Pinia
- **Auth Composable**: Updated all endpoint paths untuk consistency
- **Sidebar Navigation**: Added role-based menu visibility

### Documentation
- **Reorganized Structure**: Moved API docs ke `docs/api/` folder
- **User Journeys**: Added comprehensive user journey diagrams dengan Mermaid
  - Admin User Management journey dengan iOS UX patterns
  - Self-service Profile Management journey
- **Testing Guide**: Consolidated testing documentation
  - Backend API testing scenarios
  - Frontend component testing
  - Integration & E2E testing
  - Performance & security testing
  - Accessibility testing
- **API Reference**: Moved `user-management-api.md` ke `docs/api/user-management.md`
- **Updated Links**: Fixed all documentation cross-references

---

## [1.0.0] - 2025-12-27

### Added
- **Authentication System** (Sprint 1)
  - JWT Authentication dengan Access & Refresh Tokens.
  - Role-Based Access Control (RBAC).
  - Rate Limiting & Account Lockout.
  - Session Management.
- **Database**
  - Registry Pattern untuk automatic migrations.
  - Seeding mechanism.
- **Frontend UI**
  - iOS-inspired design system.
  - Motion-v animations.
  - Glassmorphism components.
- **Backend Architecture**
  - Service Repository Pattern dengan Gin framework.
  - Centralized error handling.

### Fixed
- Initial setup bugs.
- Database connection retry mechanism.

### Security
- Bcrypt password hashing.
- Token hashing in database.
- Secure HTTP headers.
