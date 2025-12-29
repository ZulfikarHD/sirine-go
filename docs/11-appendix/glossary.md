# 📖 Glossary - Sirine Go

Comprehensive glossary of technical terms, acronyms, dan domain-specific terminology yang digunakan dalam Sirine Go project.

---

## 🔤 Technical Terms

### A

**API (Application Programming Interface)**  
Interface yang memungkinkan komunikasi antara berbagai komponen software. Dalam Sirine Go, REST API digunakan untuk komunikasi backend-frontend.

**Authentication (Autentikasi)**  
Proses verifikasi identitas user. Sirine Go menggunakan JWT-based authentication.

**Authorization (Otorisasi)**  
Proses menentukan akses permissions user. Sirine Go menggunakan Role-Based Access Control (RBAC).

**Auto-Migration**  
Fitur GORM yang otomatis create/update database schema berdasarkan model definitions.

---

### B

**Backend**  
Server-side application yang menangani business logic, database operations, dan API endpoints. Sirine Go backend menggunakan Go dengan Gin framework.

**Bcrypt**  
Algorithm untuk password hashing yang secure dan slow by design (cost factor 12 di Sirine Go).

**Bulk Operations**  
Operations yang melakukan multiple actions sekaligus (bulk delete, bulk import).

---

### C

**CORS (Cross-Origin Resource Sharing)**  
Security mechanism yang mengontrol cross-domain requests. Configured di backend untuk allow frontend access.

**CSV (Comma-Separated Values)**  
File format untuk bulk data import/export. Sirine Go supports CSV untuk user management.

**Composable**  
Vue 3 pattern untuk reusable logic dengan Composition API (contoh: `useAuth`, `useHaptic`).

---

### D

**Database Migration**  
Process untuk update database schema secara controlled dan reversible.

**Debounce**  
Technique untuk delay function execution hingga period of inactivity (digunakan untuk search input, 300ms delay).

**Dependency Injection**  
Design pattern dimana dependencies di-pass ke component/service rather than created internally.

---

### E

**Endpoint**  
URL path yang accepts HTTP requests (contoh: `/api/users`, `/api/auth/login`).

**Environment Variable**  
Configuration values yang stored di `.env` file untuk different environments (dev, staging, production).

**ETag**  
HTTP header untuk caching dan conditional requests, contains hash of resource content.

---

### F

**Frontend**  
Client-side application yang runs di browser. Sirine Go frontend menggunakan Vue 3.

**Foreign Key**  
Database column yang references primary key di another table untuk maintain referential integrity.

---

### G

**Gin**  
High-performance HTTP web framework untuk Go, digunakan sebagai backend framework Sirine Go.

**GORM (Go Object-Relational Mapping)**  
ORM library untuk Go yang provides type-safe database operations.

**Goroutine**  
Lightweight thread managed by Go runtime, digunakan untuk async operations (contoh: activity logging).

---

### H

**Handler**  
Function yang processes HTTP requests dan returns responses. Di Sirine Go, handlers ada di `handlers/` folder.

**Haptic Feedback**  
Tactile vibration feedback pada mobile devices untuk enhance UX.

**Hash**  
One-way cryptographic function untuk convert data ke fixed-size string (digunakan untuk passwords, tokens).

---

### I

**Index (Database)**  
Data structure yang improves query performance pada specific columns.

**ISO 8601**  
International standard untuk date/time format (contoh: `2025-12-29T10:30:00Z`).

---

### J

**JWT (JSON Web Token)**  
Compact, URL-safe token format untuk transmit information securely. Sirine Go uses JWT untuk authentication.

**JSON (JavaScript Object Notation)**  
Lightweight data interchange format yang human-readable dan machine-parseable.

---

### L

**Lazy Loading**  
Technique untuk defer loading resources hingga needed, improves initial page load performance.

**Load Balancer**  
Distributes incoming traffic across multiple servers untuk improve reliability dan performance.

---

### M

**Middleware**  
Software component yang sits between request dan handler untuk process requests/responses (contoh: auth middleware, logging middleware).

**Migration**  
Script untuk modify database schema dalam controlled manner (create tables, add columns, etc.).

**Model**  
Data structure yang represents database table. Di Sirine Go, models ada di `models/` folder.

**Motion-V**  
Animation library untuk Vue 3 yang provides physics-based animations dengan spring presets.

---

### N

**NIP (Nomor Induk Pegawai)**  
Employee ID number, digunakan sebagai unique identifier dan username untuk login.

**Normalization**  
Database design technique untuk reduce redundancy dan improve data integrity.

---

### O

**ORM (Object-Relational Mapping)**  
Technique untuk convert data between relational database dan object-oriented programming language.

**Optimistic Update**  
Update UI immediately before server response, rollback jika error occurs.

---

### P

**Pagination**  
Technique untuk split large datasets into pages untuk improve performance dan UX.

**Pinia**  
Official state management library untuk Vue 3, modern replacement untuk Vuex.

**Polling**  
Technique untuk periodically check server untuk updates (30-second interval untuk notifications).

**Primary Key**  
Column atau combination of columns yang uniquely identifies each row dalam table.

---

### Q

**Query**  
Request untuk data dari database dengan specific criteria.

**Query String**  
Part of URL yang contains parameters (contoh: `/api/users?page=1&limit=20`).

---

### R

**RBAC (Role-Based Access Control)**  
Authorization system dimana permissions assigned based on user roles (Admin, Manager, Staff).

**REST (Representational State Transfer)**  
Architectural style untuk designing networked applications menggunakan HTTP methods.

**Refresh Token**  
Long-lived token (30 days) yang digunakan untuk obtain new access tokens.

**Route Guard**  
Middleware di Vue Router yang controls navigation access based on authentication/authorization.

---

### S

**Service Layer**  
Architecture layer yang contains business logic, sits between handlers dan database.

**Session**  
Period dimana user authenticated, tracked dengan tokens dan stored di database.

**Skeleton Loading**  
Placeholder UI yang shows shape of content while loading, improves perceived performance.

**Soft Delete**  
Marking record as deleted without actually removing from database (sets `deleted_at` timestamp).

**SQL (Structured Query Language)**  
Standard language untuk interact dengan relational databases.

**SPA (Single Page Application)**  
Web application yang loads single HTML page dan dynamically updates content.

---

### T

**Token**  
String yang represents user session atau authentication credential (JWT, reset token).

**Transaction**  
Group of database operations yang execute as single unit (all succeed atau all rollback).

**Tailwind CSS**  
Utility-first CSS framework untuk rapid UI development.

---

### U

**UUID (Universally Unique Identifier)**  
128-bit identifier yang globally unique, digunakan untuk generated filenames.

**UX (User Experience)**  
Overall experience user has ketika interacting dengan application.

---

### V

**Validation**  
Process untuk check data meets requirements before processing (frontend & backend validation).

**Vue 3**  
Progressive JavaScript framework untuk building user interfaces, menggunakan Composition API.

**Vite**  
Modern build tool yang provides fast development server dan optimized production builds.

---

### W

**WebSocket**  
Protocol untuk bidirectional communication antara client dan server (alternative to polling).

**Webhook**  
HTTP callback yang triggers when specific event occurs (untuk notifications, integrations).

---

## 🏢 Domain-Specific Terms

### Khazwal
**Department** yang menangani warehouse operations, inventory management, dan material preparation.

### Cetak
Process untuk printing atau marking materials dalam production workflow.

### Verifikasi
Quality control step untuk verify materials atau products meet requirements.

### Khazkhir
Final warehouse operations sebelum distribution, ensuring product readiness.

---

## 🎯 Role Definitions

### ADMIN
Super administrator dengan full system access, can manage all users dan configurations.

### MANAGER_KHAZWAL
Manager role untuk Khazwal department dengan elevated permissions untuk department operations.

### STAFF_KHAZWAL
Staff role untuk Khazwal department dengan limited permissions untuk daily operations.

### MANAGER_KEUANGAN
Manager role untuk Finance department dengan financial operations permissions.

### STAFF_KEUANGAN
Staff role untuk Finance department dengan basic financial operations access.

### MANAGER_DISTRIBUSI
Manager role untuk Distribution department dengan distribution management permissions.

### STAFF_DISTRIBUSI
Staff role untuk Distribution department dengan basic distribution operations access.

---

## 🕐 Shift Definitions

### PAGI (Morning Shift)
Typical hours: 06:00 - 14:00

### SORE (Afternoon Shift)
Typical hours: 14:00 - 22:00

### MALAM (Night Shift)
Typical hours: 22:00 - 06:00

---

## 📊 Status Definitions

### Active
User account is enabled dan can login.

### Inactive
User account is disabled, cannot login.

### Locked
User account temporarily locked due to failed login attempts atau security reasons.

---

## 🎮 Gamification Terms

### Achievement
Unlockable badge yang awarded when user meets specific criteria.

### Points
Numeric value accumulated through achievements dan activities.

### Level
Tier based on total points (Bronze, Silver, Gold, Platinum).

### Streak
Consecutive days user performs specific action (login streak).

---

## 🔐 Security Terms

### Access Token
Short-lived token (15 minutes) untuk authenticate API requests.

### Refresh Token
Long-lived token (30 days) untuk obtain new access tokens.

### Password Policy
Set of rules untuk ensure password strength (min 8 chars, uppercase, number, special char).

### Session Revocation
Process untuk invalidate user sessions (on logout, password change).

### Token Hash
SHA256 hash of token stored di database untuk security.

---

## 📞 Support

Untuk questions about terminology atau clarifications:

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733

---

## 🔗 Related Documentation

- [Architecture Overview](../02-architecture/overview.md)
- [API Documentation](../03-development/api-documentation.md)
- [Database Models](../05-guides/database/models.md)
- [Resources](./resources.md)

---

**Last Updated:** 29 Desember 2025  
**Total Terms:** 100+ definitions
