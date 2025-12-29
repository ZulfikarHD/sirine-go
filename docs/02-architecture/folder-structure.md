# 📁 Struktur Folder - Sirine Go App

## 🌳 Tree Structure

```
sirine-go/
│
├── 📂 backend/                         # Backend (Go + Gin Framework)
│   ├── 📂 cmd/                         # Entry points & Utility CLI tools
│   │   ├── 📂 server/                  # 🚀 Main API Server
│   │   ├── 📂 migrate/                 # Database Migration Tool
│   │   ├── 📂 seed/                    # Database Seeder Tool
│   │   ├── 📂 hash/                    # Password Hash Utility
│   │   ├── 📂 genhash/                 # Hash Generator Utility
│   │   ├── 📂 checkuser/               # User Verification Utility
│   │   ├── 📂 testdb/                  # Database Connection Tester
│   │   └── 📂 testpass/                # Password Verification Tester
│   │
│   ├── 📂 config/
│   │   └── config.go                   # ⚙️ Environment configuration
│   │
│   ├── 📂 database/
│   │   ├── database.go                 # 🗄️ Database connection & GORM setup
│   │   ├── models_registry.go          # 📝 Model auto-migration registry
│   │   └── setup.sql                   # 📝 Raw SQL setup script
│   │
│   ├── 📂 handlers/
│   │   ├── auth_handler.go             # 🔐 Login, Logout, Refresh Token
│   │   └── ...                         # Other request handlers
│   │
│   ├── 📂 middleware/
│   │   ├── auth_middleware.go          # 🔐 JWT Validation & Context Setting
│   │   ├── role_middleware.go          # 👮 RBAC (Role Based Access Control)
│   │   └── cors.go                     # 🔒 CORS Configuration
│   │
│   ├── 📂 models/
│   │   ├── user.go                     # 👤 User Entity & Validation
│   │   ├── user_session.go             # 🎫 Session/Device Tracking
│   │   ├── password_reset_token.go     # 🔑 Password Reset Logic
│   │   ├── activity_log.go             # 📝 User Activity Logging
│   │   └── notification.go             # 🔔 Notification System
│   │
│   ├── 📂 routes/
│   │   └── routes.go                   # 🛣️ API Route Definitions & Grouping
│   │
│   ├── 📂 services/
│   │   ├── auth_service.go             # 🔐 Auth Business Logic (Login/Register)
│   │   ├── password_service.go         # 🔑 Password Hashing & Comparison
│   │   └── user_service_requests.go    # 💼 User Management Requests
│   │
│   ├── 📂 utils/
│   │   └── validation.go               # 🛠️ Common Validation Helpers
│   │
│   ├── 📂 tests/
│   │   ├── 📂 unit/                    # 🧪 Unit Tests (Isolated)
│   │   └── 📂 integration/             # 🔄 Integration Tests (API + DB)
│   │
│   ├── .env                            # 🔐 Environment variables (GitIgnored)
│   ├── go.mod                          # 📦 Go dependencies definition
│   └── go.sum                          # 🔒 Go dependencies checksums
│
├── 📂 frontend/                        # Frontend (Vue 3 + Vite)
│   ├── 📂 src/
│   │   ├── 📂 assets/                  # 🖼️ Static assets (Images, SVG)
│   │   │   └── vue.svg
│   │   │
│   │   ├── 📂 components/              # 🧩 Reusable Vue Components
│   │   │   ├── 📂 layout/              # 📐 Layout Components (Sidebar, Navbar)
│   │   │   ├── ExampleCard.vue         # UI Component Example
│   │   │   └── ExampleForm.vue         # Form Component Example
│   │   │
│   │   ├── 📂 composables/             # 🔧 Composition API (Hooks)
│   │   │   ├── useApi.js               # 🌐 Axios Wrapper
│   │   │   ├── useAuth.js              # 🔐 Authentication Logic
│   │   │   └── useExamples.js          # 📝 Example Logic
│   │   │
│   │   ├── 📂 stores/                  # 🏪 State Management (Pinia)
│   │   │   └── auth.js                 # 👤 Auth Store (User & Token)
│   │   │
│   │   ├── 📂 router/                  # 🚦 Client-side Routing
│   │   │   └── index.js                # Route Definitions & Guards
│   │   │
│   │   ├── 📂 views/                   # 📄 Page Views
│   │   │   ├── 📂 auth/                # 🔐 Login Pages
│   │   │   ├── 📂 dashboards/          # 📊 Admin/User Dashboards
│   │   │   ├── 📂 profile/             # 👤 User Profile Pages
│   │   │   ├── Home.vue                # 🏠 Landing Page
│   │   │   └── NotFound.vue            # 🚫 404 Page
│   │   │
│   │   ├── 📂 tests/                   # 🧪 Frontend Tests (Vitest)
│   │   │   ├── 📂 unit/                # Unit Tests
│   │   │   ├── 📂 integration/         # Component Integration Tests
│   │   │   └── setup.js                # Test Environment Setup
│   │   │
│   │   ├── App.vue                     # 🏠 Root Application Component
│   │   ├── main.js                     # 🚀 Entry Point & PWA Setup
│   │   └── style.css                   # 🎨 Tailwind & Global CSS
│   │
│   ├── 📂 public/                      # 🌍 Public Static Files (Favicon, Manifest)
│   │
│   ├── index.html                      # 📄 HTML Entry Point
│   ├── vite.config.js                  # ⚡ Vite Build Config
│   ├── tailwind.config.js              # 🎨 Tailwind CSS Config
│   ├── postcss.config.js               # 🔧 CSS Post-Processing Config
│   ├── package.json                    # 📦 NPM Dependencies
│   └── yarn.lock                       # 🔒 Dependency Lock File
│
├── 📂 docs/                            # 📚 Documentation
│   ├── 📂 api/                         # 🔌 API Specifications
│   ├── 📂 architecture/                # 🏗️ System Architecture & Design
│   ├── 📂 deployment/                  # 🚀 Deployment Guides
│   ├── 📂 development/                 # 💻 Dev Guides & Workflow
│   ├── 📂 getting-started/             # 🏁 Setup & Installation
│   ├── 📂 guides/                      # 📖 General Guides & Maintenance
│   ├── 📂 testing/                     # 🧪 Testing Strategies
│   ├── 📂 troubleshooting/             # 🔧 FAQ & Debugging
│   ├── 📂 user-journeys/               # 🗺️ User Flow Documentation
│   ├── CHANGELOG.md                    # 📜 Version History
│   ├── CONTRIBUTING.md                 # 🤝 Contribution Guidelines
│   └── README.md                       # 📘 Documentation Home
│
├── 📂 Sirine-Reqs/                     # 📋 Requirements & Specs
├── 📄 Makefile                         # 🛠️ Task Runner (Build, Test, Run)
├── 📄 README.md                        # 📖 Project Root Documentation
├── 📄 SPRINT1_SUMMARY.md               # 📊 Sprint 1 Progress Report
└── .gitignore                          # 🚫 Git Ignore Rules
```

---

## 📂 Backend Folder Details

### `cmd/`
**Purpose:** Entry points untuk aplikasi dan utility tools.
**Files:**
- `server/main.go`: Menjalankan API Server utama.
- `migrate/main.go`: Menjalankan database migration.
- `seed/main.go`: Mengisi database dengan data awal (seeding).
- `hash/`, `genhash/`: Utility untuk generate/check password hash manual.
- `testdb/`: Utility untuk cek koneksi database.

### `middleware/`
**Purpose:** HTTP Middleware untuk intercept request.
**Files:**
- `auth_middleware.go`: Validasi JWT Token di header `Authorization`.
- `role_middleware.go`: Membatasi akses berdasarkan role user (RBAC).
- `cors.go`: Konfigurasi Cross-Origin Resource Sharing.

### `models/`
**Purpose:** Definisi struktur data (GORM Models).
**Files:**
- `user.go`: Struktur tabel user & validasi.
- `user_session.go`: Tracking session login active.
- `activity_log.go`: Mencatat aktivitas user.

### `utils/`
**Purpose:** Helper functions umum.
**Files:**
- `validation.go`: Fungsi validasi input yang reusable.

---

## 📂 Frontend Folder Details

### `src/components/layout/`
**Purpose:** Komponen layout yang digunakan di banyak halaman.
**Responsibility:** Sidebar, Navbar, Footer, Container layouts.

### `src/stores/`
**Purpose:** Global State Management menggunakan Pinia.
**Files:** `auth.js`
**Responsibility:** Menyimpan status login, token, dan user profile secara global agar bisa diakses component manapun.

### `src/views/`
**Purpose:** Halaman utama aplikasi (Pages).
**Structure:**
- `auth/`: Halaman Login, Register, Forgot Password.
- `dashboards/`: Halaman dashboard setelah login (Admin, User).
- `profile/`: Halaman edit profile user.

### `src/assets/`
**Purpose:** Static assets yang di-import di code.
**Files:** Images, SVGs, global icons.

---

## 🔍 Find Files Quickly

### Backend:
```bash
# Find all utility tools in cmd
ls -R backend/cmd

# Find all Unit Tests
ls -R backend/tests/unit
```

### Frontend:
```bash
# Find all Layout Components
ls frontend/src/components/layout

# Find all Page Views
ls -R frontend/src/views
```

---

**Developer:** Zulfikar Hidayatullah  
**Date:** 28 Desember 2025  
**Version:** 1.0.3 (Updated with full docs structure)
