# 📚 Dokumentasi Sirine Go App

> Modern full-stack web application dengan 100% offline capability menggunakan Go (Gin) dan Vue 3.

## 🎯 Tentang Sirine Go App

**Sirine Go App** adalah aplikasi web modern yang dapat berfungsi sepenuhnya offline, dibangun dengan teknologi terkini untuk memberikan performa optimal dan user experience terbaik.

**Tech Stack:**
- **Backend:** Go (Gin) + MySQL + GORM
- **Frontend:** Vue 3 + Vite + Tailwind CSS 4.1.18 + Motion-v
- **Architecture:** RESTful API dengan Service Pattern
- **Special Features:** Progressive Web App (PWA) untuk offline capabilities

---

## 🚀 Quick Start

Untuk setup cepat (5 menit), gunakan command berikut:

```bash
# 1. Setup database
mysql -u root -p -e "CREATE DATABASE sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 2. Edit .env
nano backend/.env  # Set DB_PASSWORD

# 3. Install dependencies
make install

# 4. Run (buka 2 terminal)
make dev-backend   # Terminal 1
make dev-frontend  # Terminal 2

# 5. Buka browser: http://localhost:5173
```

**📖 Detail lengkap:** Lihat [getting-started/quickstart.md](./01-getting-started/quickstart.md)

---

## 🆕 Recent Updates

### Sprint OBC Master: Excel Import & PO Generation (2024-12-30)

**Status**: ✅ Complete

**Quick Links:**
- 📋 **[sprint-obc-master.md](./10-sprints/sprint-obc-master.md)** - Complete sprint summary
- 📡 **[obc-master.md](./04-api-reference/obc-master.md)** - API documentation
- 🧪 **[obc-import-testing.md](./06-testing/obc-import-testing.md)** - Testing guide

**What's Done:**
- ✅ OBC Master model dengan 39 fields
- ✅ Excel import dari SAP dengan upsert logic
- ✅ Auto PO generation dengan buffer 6%
- ✅ 4 REST API endpoints (import, list, detail, generate)
- ✅ PPIC role & authorization
- ✅ Complete documentation & test data

---

### Sprint 6: Testing, Optimization & Deployment (2025-01-29)

**Status**: 🟡 In Progress - Testing Phase

**Quick Links:**
- 📊 **[SPRINT6_SUMMARY.md](./SPRINT6_SUMMARY.md)** - Complete sprint summary (updated with actual results)
- 🧪 **[06-testing/SPRINT6_TESTING.md](./06-testing/SPRINT6_TESTING.md)** - Honest test execution report
- 🚀 **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Production deployment guide
- 🔐 **[AUTHENTICATION.md](./AUTHENTICATION.md)** - Authentication system documentation

**What's Done:**
- ✅ Comprehensive test suite created (70+ test cases)
- ✅ Database optimization (composite indexes, stored procedures)
- ✅ Frontend optimization (code splitting, 300KB bundle)
- ✅ PWA configured (offline support, installable)
- ✅ Security hardening (rate limiting, security headers)
- ✅ Complete deployment documentation

**What's In Progress:**
- ⚠️ Test execution (~50% run, 94% passing of those run)
- ⚠️ Fixing test mocks (axios, localStorage)
- 🔄 Running complete test suite

**Honest Note**: Tests were created but not all executed during initial sprint. Documentation now reflects actual testing reality with transparent reporting.

---

## 📖 Dokumentasi

Dokumentasi terorganisir dalam folders berdasarkan kategori:

### 📘 **Getting Started** → `getting-started/`

Mulai di sini jika baru pertama kali:

1. **[quickstart.md](./01-getting-started/quickstart.md)** ⚡  
   Setup dalam 5 menit - Step-by-step paling simple

2. **[installation.md](./01-getting-started/installation.md)** 📋  
   Setup lengkap dengan penjelasan detail setiap step

3. **[checklist.md](./01-getting-started/checklist.md)** ✅  
   Checklist untuk verifikasi setup sudah benar

---

### 🏗️ **Architecture** → `architecture/`

Pahami struktur dan design decisions:

4. **[overview.md](./02-architecture/overview.md)** 🏛️  
   Penjelasan semua package dan kenapa dibutuhkan

5. **[folder-structure.md](./02-architecture/folder-structure.md)** 📁  
   Struktur folder dan file organization

6. **[project-summary.md](./02-architecture/project-summary.md)** 📊  
   Overview lengkap project (features, decisions, metrics)

---

### 🛠️ **Development** → `development/`

Build dan test fitur baru:

7. **[customization-guide.md](./03-development/customization-guide.md)** 🎨  
   Cara menambah model, endpoint, dan component baru

8. **[api-documentation.md](./03-development/api-documentation.md)** 🔌  
   Complete API reference dengan contoh request/response

9. **[testing.md](./03-development/testing.md)** 🧪  
   Testing guide (manual & automated testing)

---

### 📖 **Guides** → `guides/`

Panduan mendalam untuk topik spesifik:

10. **[guides/authentication/README.md](./05-guides/authentication/README.md)** 🔐  
    Sistem autentikasi, security flows, dan user journeys

11. **[guides/database/management.md](./05-guides/database/management.md)** 🗄️  
    Manajemen database, backup, dan maintenance

12. **[guides/database/models.md](./05-guides/database/models.md)** 📝  
    Panduan membuat model baru dengan Registry Pattern

13. **[guides/database/migrations.md](./05-guides/database/migrations.md)** 🔄  
    Panduan migrasi dan database seeding

14. **[guides/validation/guide.md](./05-guides/validation/guide.md)** ✅  
    Server-side validation dengan Gin (Laravel-style)

15. **[guides/security.md](./05-guides/security.md)** 🛡️  
    Security best practices dan protection mechanisms

16. **[guides/configuration.md](./05-guides/configuration.md)** ⚙️  
    Environment variables dan configuration guide

17. **[guides/documentation-maintenance.md](./05-guides/documentation-maintenance.md)** 📝  
    Panduan maintenance dokumentasi (When & How to update)

18. **[guides/utilities/hash-commands.md](./05-guides/utilities/hash-commands.md)** 🔧  
    Utility commands untuk hashing dan security

---

### 🔌 **API Reference** → `api/`

Complete API documentation untuk semua endpoints:

19. **[api/README.md](./04-api-reference/README.md)** 📡  
    API hub dengan quick reference dan conventions

20. **[api/user-management.md](./04-api-reference/user-management.md)** 👥  
    User Management & Profile API (Sprint 2)

21. **[api/khazwal.md](./04-api-reference/khazwal.md)** 📦  
    Khazwal Material Preparation API (Sprint Khazwal)

22. **[api/obc-master.md](./04-api-reference/obc-master.md)** 📋  
    OBC Master Excel Import & PO Generation API (Sprint OBC Master)

---

### 🗺️ **User Journeys** → `user-journeys/`

Detailed user flows dan interaction patterns:

21. **[user-journeys/user-management/admin-user-management.md](./07-user-journeys/user-management/admin-user-management.md)** 👨‍💼  
    Admin user management journey dengan iOS-inspired UX

22. **[user-journeys/user-management/user-profile-management.md](./07-user-journeys/user-management/user-profile-management.md)** 👤  
    Self-service profile management journey

23. **[user-journeys/khazwal/material-prep-flow.md](./07-user-journeys/khazwal/material-prep-flow.md)** 📦  
    Khazwal material preparation workflow journey

---

### 🧪 **Testing** → `testing/`

Comprehensive testing guides dan test scenarios:

24. **[testing/user-management-testing.md](./06-testing/user-management-testing.md)** ✅  
    Complete testing guide untuk User Management & Profile

25. **[testing/khazwal-testing.md](./06-testing/khazwal-testing.md)** 📦  
    Complete testing guide untuk Khazwal Material Preparation

26. **[testing/obc-import-testing.md](./06-testing/obc-import-testing.md)** 📋  
    Complete testing guide untuk OBC Master Excel Import

---

### 🚀 **Deployment** → `deployment/`

Ready untuk production:

24. **[production-deployment.md](./08-deployment/production-deployment.md)** 🌐  
    Deployment ke production server dengan Nginx + SSL

---

### ❓ **Troubleshooting** → `troubleshooting/`

Mengalami masalah:

25. **[faq.md](./09-troubleshooting/faq.md)** 💡  
    Common questions dan solusi masalah umum

---

## 📋 Prerequisites

Pastikan terinstall:
- ✅ **Go 1.24+** - Backend language
- ✅ **Node.js 18+** & **Yarn** - Frontend tooling  
- ✅ **MySQL 8.0+** - Database
- ✅ **Git** - Version control

**Verify:**
```bash
go version && node --version && yarn --version && mysql --version
```

---

## 🎨 Key Features

### **🔐 Authentication & Security (Sprint 1-3)**
- JWT-based authentication dengan refresh token
- Role-based access control (7 roles, 4 departments)
- Password management (change, forgot, reset)
- Session tracking dengan IP & User Agent
- Rate limiting & account lockout
- Password strength enforcement
- Bcrypt cost 12 untuk password hashing

### **👥 User Management (Sprint 2)**
- Admin CRUD users dengan comprehensive features
- Search & filters (role, department, status)
- Bulk operations (delete, update status)
- Auto-generated secure passwords
- Self-service profile management
- Activity logging untuk audit trail
- Pagination support

### **🔔 Notifications & Audit (Sprint 4)**
- In-app notification system
- Real-time updates via polling (30s)
- Notification center dengan filtering
- Activity logs viewer untuk Admin
- Before/after change tracking
- Comprehensive audit trail
- Advanced filters dan statistics

### **🎮 Gamification (Sprint 5)**
- Achievement system dengan 6 achievements
- Points tracking dan level system (Bronze → Platinum)
- Profile photo upload dengan auto-resize
- CSV bulk import/export users
- Haptic feedback (7 patterns)
- Loading skeletons untuk better UX
- Animated points counter

### **📦 Khazwal Material Preparation (Sprint Khazwal)**
- Queue management untuk Production Orders
- Multi-step workflow (Plat, Kertas, Tinta)
- Barcode scanner dengan camera integration
- Real-time variance calculation
- Photo upload untuk material evidence
- Automatic notification ke Unit Cetak
- Supervisor monitoring dashboard
- History tracking dengan date filters

### **📋 OBC Master & PO Generation (Sprint OBC Master)**
- Excel import dari SAP export (39 fields)
- Auto-generate Production Orders dengan buffer 6%
- Intelligent PO splitting (max 40k per PO)
- Upsert logic untuk update existing data
- List & filter dengan composite indexes
- PPIC role untuk access control
- Relationship tracking ke Production Orders
- Complete audit trail

### **🎨 Modern UI/UX**
- Apple-inspired design (Indigo & Fuchsia gradient)
- Glass effect cards dengan backdrop blur
- Spring physics animations (Motion-V)
- iOS-style press feedback
- Responsive design (mobile-first)
- Loading states & empty states
- Staggered list animations
- Drag & drop file upload

### **🚀 Performance**
- Vite build tool (10-100x faster than Webpack)
- GORM ORM (type-safe & optimized)
- Gin framework (40,000+ req/s)
- Optimistic updates untuk instant UI feedback
- Debounced search (300ms)
- Background goroutines untuk non-blocking operations
- Asset caching & optimization
- Image optimization (auto-resize, JPEG 90%)

### **🌐 100% Offline Capable**
- Service Worker untuk caching
- PWA installable (Desktop & Mobile)
- NetworkFirst caching strategy untuk API
- Online/Offline status indicator real-time

---

## 🛠️ Development Commands

```bash
make help              # Show all available commands
make install           # Install all dependencies
make dev-backend       # Run backend development server
make dev-frontend      # Run frontend development server
make build             # Build for production
make clean             # Clean build files
```

---

## 📂 Project Structure

```
sirine-go/
├── backend/           # Go + Gin backend
│   ├── cmd/          # Entry point
│   ├── config/       # Configuration
│   ├── database/     # Database setup
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # Middleware
│   ├── models/       # Database models
│   ├── routes/       # API routes
│   └── services/     # Business logic
│
├── frontend/         # Vue 3 + Vite frontend
│   ├── src/
│   │   ├── components/  # Vue components
│   │   ├── composables/ # Reusable logic
│   │   └── views/       # Page views
│   └── public/       # Static assets
│
└── docs/             # Documentation (you are here!)
```

**📖 Detail lengkap:** Lihat [folder-structure.md](./02-architecture/folder-structure.md)

---

## 🔌 API Quick Reference

```http
# Health check
GET /health

# Authentication (Sprint 1 & 3)
POST   /api/auth/login              # Login
POST   /api/auth/logout             # Logout
GET    /api/auth/me                 # Current user
POST   /api/auth/refresh            # Refresh token
POST   /api/auth/forgot-password    # Request reset link
POST   /api/auth/reset-password     # Reset dengan token

# User Management - Admin (Sprint 2 & 5)
GET    /api/users                   # List users dengan filters
GET    /api/users/search            # Search users
GET    /api/users/:id               # Get user detail
POST   /api/users                   # Create user
PUT    /api/users/:id               # Update user
DELETE /api/users/:id               # Delete user
POST   /api/users/bulk-delete       # Bulk delete
POST   /api/users/bulk-update-status # Bulk update status
POST   /api/users/import            # Import dari CSV
GET    /api/users/export            # Export ke CSV
POST   /api/users/:id/reset-password # Admin force reset

# Profile - Self Service (Sprint 2, 3, 5)
GET    /api/profile                 # Get profile
PUT    /api/profile                 # Update profile
PUT    /api/profile/password        # Change password
POST   /api/profile/photo           # Upload photo
DELETE /api/profile/photo           # Delete photo
GET    /api/profile/activity        # Get activity logs

# Notifications (Sprint 4)
GET    /api/notifications           # List notifications
GET    /api/notifications/unread-count # Badge count
GET    /api/notifications/recent    # Recent notifications
PUT    /api/notifications/:id/read  # Mark as read
PUT    /api/notifications/read-all  # Mark all as read
DELETE /api/notifications/:id       # Delete notification

# Activity Logs - Admin (Sprint 4)
GET    /api/admin/activity-logs     # List logs dengan filters
GET    /api/admin/activity-logs/:id # Log detail
GET    /api/admin/activity-logs/user/:id # User activity
GET    /api/admin/activity-logs/stats    # Statistics

# Achievements & Gamification (Sprint 5)
GET    /api/achievements            # List achievements
GET    /api/profile/achievements    # User achievements
GET    /api/profile/stats           # User gamification stats
POST   /api/admin/achievements/award # Award achievement
GET    /api/admin/users/:id/achievements # User achievements (Admin)

# Khazwal Material Preparation (Sprint Khazwal)
GET    /api/khazwal/material-prep/queue         # Queue PO
GET    /api/khazwal/material-prep/:id           # Detail PO
POST   /api/khazwal/material-prep/:id/start     # Start preparation
POST   /api/khazwal/material-prep/:id/confirm-plat # Confirm plat
PATCH  /api/khazwal/material-prep/:id/kertas    # Update kertas
PATCH  /api/khazwal/material-prep/:id/tinta     # Update tinta
POST   /api/khazwal/material-prep/:id/finalize  # Finalize preparation
GET    /api/khazwal/material-prep/history       # History
GET    /api/khazwal/monitoring                  # Supervisor monitoring

# Cetak Queue (Sprint Khazwal)
GET    /api/cetak/queue             # Queue PO siap cetak
GET    /api/cetak/queue/:id         # Detail PO untuk cetak

# OBC Master (Sprint OBC Master)
POST   /api/obc/import              # Import Excel dari SAP
GET    /api/obc                     # List OBC dengan filters
GET    /api/obc/:id                 # Detail OBC + Production Orders
POST   /api/obc/:id/generate-po     # Generate POs manually
```

**📖 Detail lengkap:** Lihat [development/api-documentation.md](./03-development/api-documentation.md)

---

## 🐛 Troubleshooting Quick Fix

### Database Error?
```bash
sudo systemctl start mysql
mysql -u root -p -e "USE sirine_go;"
```

### Port Already in Use?
```bash
sudo lsof -i :8080 && sudo kill -9 <PID>  # Backend
sudo lsof -i :5173 && sudo kill -9 <PID>  # Frontend
```

### Dependencies Error?
```bash
make clean && make install
```

**📖 Solusi lengkap:** Lihat [faq.md](./09-troubleshooting/faq.md)

---

## 📚 Documentation Flow

**Rekomendasi urutan baca dokumentasi:**

```
START HERE
    ↓
README.md (Overview) ← You are here
    ↓
📘 getting-started/
    ├─ quickstart.md (5 menit setup)
    ├─ installation.md (Detail setup)
    └─ checklist.md (Verify setup)
    ↓
🏗️ architecture/
    ├─ overview.md (Understand tech)
    ├─ folder-structure.md (Understand structure)
    └─ project-summary.md (Overview)
    ↓
🛠️ development/
    ├─ customization-guide.md (Build features)
    ├─ api-documentation.md (API reference)
    └─ testing.md (Test your code)
    ↓
📖 guides/
    ├─ authentication/ (Security & Flows)
    ├─ database/ (Models & Migrations)
    ├─ security.md (Best Practices)
    ├─ configuration.md (Env Vars)
    ├─ documentation-maintenance.md (How-to Doc)
    ├─ validation/ (Rules & Examples)
    └─ utilities/ (Helpers)
    ↓
🔌 api/
    ├─ README.md (API Hub)
    └─ user-management.md (User APIs)
    ↓
🗺️ user-journeys/
    └─ user-management/ (User flows & UX)
    ↓
🧪 testing/
    └─ user-management-testing.md (Test scenarios)
    ↓
🚀 deployment/
    └─ production-deployment.md (Deploy to production)
    ↓
❓ troubleshooting/
    └─ faq.md (When stuck)
```

---

## 🎯 Best Practices

### Backend
- ✅ Service Pattern untuk separation of concerns
- ✅ Error handling dalam Bahasa Indonesia
- ✅ Consistent API response format
- ✅ Environment-based configuration

### Frontend
- ✅ Composition API untuk better organization
- ✅ Composables untuk reusable logic
- ✅ Component-based architecture
- ✅ Mobile-first responsive design
- ✅ Smooth animations untuk better UX

---

## 🤝 Contributing

Ingin berkontribusi? Silakan baca **[CONTRIBUTING.md](./CONTRIBUTING.md)** untuk panduan development workflow dan pull request standards.

---

## 👨‍💻 Developer Info

**Zulfikar Hidayatullah**
- 📞 Phone: +62 857-1583-8733
- 🌍 Timezone: Asia/Jakarta (WIB)
- 🧠 Personality: INFJ (Professional)

---

## 📄 Version & License

- **Version:** 1.7.0 (Sprint OBC Master Complete)
- **Last Updated:** 30 Desember 2024
- **License:** Private & Proprietary
- **Changelog:** Lihat **[CHANGELOG.md](../CHANGELOG.md)** untuk complete changelog Sprint 1-5

### Sprint Status

| Sprint | Status | Features |
|--------|--------|----------|
| **Sprint 1** | ✅ Complete | Authentication, JWT, RBAC, Rate Limiting |
| **Sprint 2** | ✅ Complete | User Management, Profile, Search & Filters |
| **Sprint 3** | ✅ Complete | Password Management, Forgot/Reset, Force Change |
| **Sprint 4** | ✅ Complete | Notifications, Activity Logs, Audit Trail |
| **Sprint 5** | ✅ Complete | Gamification, Photo Upload, CSV Import/Export |
| **Sprint Khazwal** | ✅ Complete | Material Preparation, Workflow, Monitoring |
| **Sprint 6** | 🚧 Next | Testing, Optimization, Production Deployment |

---

## 🎉 Ready to Start?

Pilih path Anda:

- 🚀 **Quick Setup:** [getting-started/quickstart.md](./01-getting-started/quickstart.md)
- 📖 **Detail Setup:** [getting-started/installation.md](./01-getting-started/installation.md)
- 🏗️ **Understand Architecture:** [architecture/](./02-architecture/)
- 🎨 **Build Features:** [development/customization-guide.md](./03-development/customization-guide.md)
- 🌐 **Deploy:** [deployment/production-deployment.md](./08-deployment/production-deployment.md)

**Happy Coding! 🎯**
