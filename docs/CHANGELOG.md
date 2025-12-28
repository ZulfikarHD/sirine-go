# 📝 Changelog

Semua perubahan signifikan pada project Sirine Go akan didokumentasikan di file ini.

Format mengikuti [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), dan versioning mengikuti [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Security Guide**: Panduan keamanan komprehensif.
- **Configuration Guide**: Detail environment variables.
- **Contributing Guide**: Workflow development.
- **Validation Guide**: Dokumentasi validasi input (Laravel-style).
- **Authentication Docs**: Konsolidasi dokumentasi auth sistem.

### Changed
- **Documentation Structure**: Reorganisasi folder `docs/` menjadi struktur yang lebih modular (Laravel-style).
- **Database Guide**: Split menjadi `management`, `models`, dan `migrations`.

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
