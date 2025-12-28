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

**📖 Detail lengkap:** Lihat [getting-started/quickstart.md](./getting-started/quickstart.md)

---

## 📖 Dokumentasi

Dokumentasi terorganisir dalam folders berdasarkan kategori:

### 📘 **Getting Started** → `getting-started/`

Mulai di sini jika baru pertama kali:

1. **[quickstart.md](./getting-started/quickstart.md)** ⚡  
   Setup dalam 5 menit - Step-by-step paling simple

2. **[installation.md](./getting-started/installation.md)** 📋  
   Setup lengkap dengan penjelasan detail setiap step

3. **[checklist.md](./getting-started/checklist.md)** ✅  
   Checklist untuk verifikasi setup sudah benar

---

### 🏗️ **Architecture** → `architecture/`

Pahami struktur dan design decisions:

4. **[overview.md](./architecture/overview.md)** 🏛️  
   Penjelasan semua package dan kenapa dibutuhkan

5. **[folder-structure.md](./architecture/folder-structure.md)** 📁  
   Struktur folder dan file organization

6. **[project-summary.md](./architecture/project-summary.md)** 📊  
   Overview lengkap project (features, decisions, metrics)

---

### 🛠️ **Development** → `development/`

Build dan test fitur baru:

7. **[customization-guide.md](./development/customization-guide.md)** 🎨  
   Cara menambah model, endpoint, dan component baru

8. **[api-documentation.md](./development/api-documentation.md)** 🔌  
   Complete API reference dengan contoh request/response

9. **[testing.md](./development/testing.md)** 🧪  
   Testing guide (manual & automated testing)

---

### 📖 **Guides** → `guides/`

Panduan mendalam untuk topik spesifik:

10. **[guides/authentication/README.md](./guides/authentication/README.md)** 🔐  
    Sistem autentikasi, security flows, dan user journeys

11. **[guides/database/management.md](./guides/database/management.md)** 🗄️  
    Manajemen database, backup, dan maintenance

12. **[guides/database/models.md](./guides/database/models.md)** 📝  
    Panduan membuat model baru dengan Registry Pattern

13. **[guides/database/migrations.md](./guides/database/migrations.md)** 🔄  
    Panduan migrasi dan database seeding

14. **[guides/validation/guide.md](./guides/validation/guide.md)** ✅  
    Server-side validation dengan Gin (Laravel-style)

15. **[guides/security.md](./guides/security.md)** 🛡️  
    Security best practices dan protection mechanisms

16. **[guides/configuration.md](./guides/configuration.md)** ⚙️  
    Environment variables dan configuration guide

17. **[guides/documentation-maintenance.md](./guides/documentation-maintenance.md)** 📝  
    Panduan maintenance dokumentasi (When & How to update)

18. **[guides/utilities/hash-commands.md](./guides/utilities/hash-commands.md)** 🔧  
    Utility commands untuk hashing dan security

---

### 🔌 **API Reference** → `api/`

Complete API documentation untuk semua endpoints:

19. **[api/README.md](./api/README.md)** 📡  
    API hub dengan quick reference dan conventions

20. **[api/user-management.md](./api/user-management.md)** 👥  
    User Management & Profile API (Sprint 2)

---

### 🗺️ **User Journeys** → `user-journeys/`

Detailed user flows dan interaction patterns:

21. **[user-journeys/user-management/admin-user-management.md](./user-journeys/user-management/admin-user-management.md)** 👨‍💼  
    Admin user management journey dengan iOS-inspired UX

22. **[user-journeys/user-management/user-profile-management.md](./user-journeys/user-management/user-profile-management.md)** 👤  
    Self-service profile management journey

---

### 🧪 **Testing** → `testing/`

Comprehensive testing guides dan test scenarios:

23. **[testing/user-management-testing.md](./testing/user-management-testing.md)** ✅  
    Complete testing guide untuk User Management & Profile

---

### 🚀 **Deployment** → `deployment/`

Ready untuk production:

24. **[production-deployment.md](./deployment/production-deployment.md)** 🌐  
    Deployment ke production server dengan Nginx + SSL

---

### ❓ **Troubleshooting** → `troubleshooting/`

Mengalami masalah:

25. **[faq.md](./troubleshooting/faq.md)** 💡  
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

### **🌐 100% Offline Capable**
- Service Worker untuk caching
- PWA installable (Desktop & Mobile)
- NetworkFirst caching strategy untuk API
- Online/Offline status indicator real-time

### **🎨 Modern UI/UX**
- Smooth animations dengan Motion-v
- Tailwind CSS (newest version 4.1.18)
- Responsive design (mobile-first)
- Loading & error states
- Professional personality (INFJ)

### **🚀 Performance**
- Vite build tool (10-100x faster)
- GORM ORM (type-safe & optimized)
- Gin framework (40,000+ req/s)
- Asset caching & optimization

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

**📖 Detail lengkap:** Lihat [folder-structure.md](./architecture/folder-structure.md)

---

## 🔌 API Quick Reference

```http
# Health check
GET /health

# Authentication
POST   /api/auth/login    # Login
POST   /api/auth/logout   # Logout
GET    /api/auth/me       # Current user

# User Management (Admin)
GET    /api/users         # List users
POST   /api/users         # Create user
PUT    /api/users/:id     # Update user
DELETE /api/users/:id     # Delete user

# Profile (Self-Service)
GET    /api/profile       # Get profile
PUT    /api/profile       # Update profile
```

**📖 Detail lengkap:** Lihat [api/README.md](./api/README.md)

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

**📖 Solusi lengkap:** Lihat [faq.md](./troubleshooting/faq.md)

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

- **Version:** 1.1.0 (Sprint 2 Complete)
- **Last Updated:** 28 Desember 2025
- **License:** Private & Proprietary
- **Changelog:** Lihat **[CHANGELOG.md](./CHANGELOG.md)** untuk riwayat perubahan.

---

## 🎉 Ready to Start?

Pilih path Anda:

- 🚀 **Quick Setup:** [getting-started/quickstart.md](./getting-started/quickstart.md)
- 📖 **Detail Setup:** [getting-started/installation.md](./getting-started/installation.md)
- 🏗️ **Understand Architecture:** [architecture/](./architecture/)
- 🎨 **Build Features:** [development/customization-guide.md](./development/customization-guide.md)
- 🌐 **Deploy:** [deployment/production-deployment.md](./deployment/production-deployment.md)

**Happy Coding! 🎯**
