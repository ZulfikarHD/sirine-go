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

**📖 Detail lengkap:** Lihat [getting-started/QUICKSTART.md](./getting-started/QUICKSTART.md)

---

## 📖 Dokumentasi

Dokumentasi terorganisir dalam folders berdasarkan kategori:

### 📘 **Getting Started** → `getting-started/`

Mulai di sini jika baru pertama kali:

1. **[QUICKSTART.md](./getting-started/QUICKSTART.md)** ⚡  
   Setup dalam 5 menit - Step-by-step paling simple

2. **[SETUP_GUIDE.md](./getting-started/SETUP_GUIDE.md)** 📋  
   Setup lengkap dengan penjelasan detail setiap step

3. **[CHECKLIST.md](./getting-started/CHECKLIST.md)** ✅  
   Checklist untuk verifikasi setup sudah benar

---

### 🏗️ **Architecture** → `architecture/`

Pahami struktur dan design decisions:

4. **[ARCHITECTURE_EXPLAINED.md](./architecture/ARCHITECTURE_EXPLAINED.md)** 🏛️  
   Penjelasan semua package dan kenapa dibutuhkan

5. **[FOLDER_STRUCTURE.md](./architecture/FOLDER_STRUCTURE.md)** 📁  
   Struktur folder dan file organization

6. **[PROJECT_SUMMARY.md](./architecture/PROJECT_SUMMARY.md)** 📊  
   Overview lengkap project (features, decisions, metrics)

---

### 🛠️ **Development** → `development/`

Build dan test fitur baru:

7. **[CUSTOMIZATION_GUIDE.md](./development/CUSTOMIZATION_GUIDE.md)** 🎨  
   Cara menambah model, endpoint, dan component baru

8. **[API_DOCUMENTATION.md](./development/API_DOCUMENTATION.md)** 🔌  
   Complete API reference dengan contoh request/response

9. **[VALIDATION_GUIDE.md](./VALIDATION_GUIDE.md)** ✅  
   Server-side validation dengan Gin (Laravel-style)

10. **[VALIDATION_EXAMPLES.md](./VALIDATION_EXAMPLES.md)** 📋  
    Practical validation examples dan test cases

11. **[TESTING.md](./development/TESTING.md)** 🧪  
    Testing guide (manual & automated testing)

---

### 🚀 **Deployment** → `deployment/`

Ready untuk production:

10. **[DEPLOYMENT.md](./deployment/DEPLOYMENT.md)** 🌐  
    Deployment ke production server dengan Nginx + SSL

---

### ❓ **Troubleshooting** → `troubleshooting/`

Mengalami masalah:

11. **[FAQ.md](./troubleshooting/FAQ.md)** 💡  
    Common questions dan solusi masalah umum

---

## 📋 Prerequisites

Pastikan terinstall:
- ✅ **Go 1.24+** - Backend language
- ✅ **Node.js 18+ & Yarn** - Frontend tooling  
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

**📖 Detail lengkap:** Lihat [FOLDER_STRUCTURE.md](./FOLDER_STRUCTURE.md)

---

## 🔌 API Quick Reference

```http
# Health check
GET /health

# Examples CRUD
GET    /api/examples      # Get all
GET    /api/examples/:id  # Get by ID
POST   /api/examples      # Create
PUT    /api/examples/:id  # Update
DELETE /api/examples/:id  # Delete
```

**📖 Detail lengkap:** Lihat [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

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

**📖 Solusi lengkap:** Lihat [troubleshooting/FAQ.md](./troubleshooting/FAQ.md)

---

## 📚 Documentation Flow

**Rekomendasi urutan baca dokumentasi:**

```
START HERE
    ↓
README.md (Overview) ← You are here
    ↓
📘 getting-started/
    ├─ QUICKSTART.md (5 menit setup)
    ├─ SETUP_GUIDE.md (Detail setup)
    └─ CHECKLIST.md (Verify setup)
    ↓
🏗️ architecture/
    ├─ ARCHITECTURE_EXPLAINED.md (Understand tech)
    ├─ FOLDER_STRUCTURE.md (Understand structure)
    └─ PROJECT_SUMMARY.md (Overview)
    ↓
🛠️ development/
    ├─ CUSTOMIZATION_GUIDE.md (Build features)
    ├─ API_DOCUMENTATION.md (API reference)
    └─ TESTING.md (Test your code)
    ↓
🚀 deployment/
    └─ DEPLOYMENT.md (Deploy to production)
    ↓
❓ troubleshooting/
    └─ FAQ.md (When stuck)
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

## 👨‍💻 Developer Info

**Zulfikar Hidayatullah**
- 📞 Phone: +62 857-1583-8733
- 🌍 Timezone: Asia/Jakarta (WIB)
- 🧠 Personality: INFJ (Professional)

---

## 📄 Version & License

- **Version:** 1.0.0
- **Last Updated:** 27 Desember 2025
- **License:** Private & Proprietary

---

## 🎉 Ready to Start?

Pilih path Anda:

- 🚀 **Quick Setup:** [getting-started/QUICKSTART.md](./getting-started/QUICKSTART.md)
- 📖 **Detail Setup:** [getting-started/SETUP_GUIDE.md](./getting-started/SETUP_GUIDE.md)
- 🏗️ **Understand Architecture:** [architecture/](./architecture/)
- 🎨 **Build Features:** [development/CUSTOMIZATION_GUIDE.md](./development/CUSTOMIZATION_GUIDE.md)
- 🌐 **Deploy:** [deployment/DEPLOYMENT.md](./deployment/DEPLOYMENT.md)

**Happy Coding! 🎯**
