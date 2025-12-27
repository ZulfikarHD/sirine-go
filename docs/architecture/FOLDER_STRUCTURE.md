# 📁 Struktur Folder - Sirine Go App

## 🌳 Tree Structure

```
sirine-go/
│
├── 📂 backend/                         # Backend (Go + Gin Framework)
│   ├── 📂 cmd/
│   │   └── 📂 server/
│   │       └── main.go                 # 🚀 Entry point aplikasi
│   │
│   ├── 📂 config/
│   │   └── config.go                   # ⚙️ Environment configuration
│   │
│   ├── 📂 database/
│   │   ├── database.go                 # 🗄️ Database connection & migration
│   │   └── setup.sql                   # 📝 SQL setup script
│   │
│   ├── 📂 handlers/
│   │   └── example_handler.go          # 🎯 HTTP request handlers
│   │
│   ├── 📂 middleware/
│   │   └── cors.go                     # 🔒 CORS middleware
│   │
│   ├── 📂 models/
│   │   └── example.go                  # 📊 Database models (GORM)
│   │
│   ├── 📂 routes/
│   │   └── routes.go                   # 🛣️ API route definitions
│   │
│   ├── 📂 services/
│   │   └── example_service.go          # 💼 Business logic layer
│   │
│   ├── .env                            # 🔐 Environment variables
│   ├── go.mod                          # 📦 Go dependencies
│   └── go.sum                          # 🔒 Go dependencies lock
│
├── 📂 frontend/                        # Frontend (Vue 3 + Vite)
│   ├── 📂 src/
│   │   ├── 📂 components/              # 🧩 Vue components
│   │   │   ├── ExampleCard.vue         # Card component dengan animasi
│   │   │   └── ExampleForm.vue         # Form component
│   │   │
│   │   ├── 📂 composables/             # 🔧 Reusable logic
│   │   │   ├── useApi.js               # API client (axios)
│   │   │   └── useExamples.js          # Business logic untuk examples
│   │   │
│   │   ├── 📂 views/                   # 📄 Page views
│   │   │   └── Home.vue                # Main page
│   │   │
│   │   ├── App.vue                     # 🏠 Root component
│   │   ├── main.js                     # 🚀 Entry point + PWA registration
│   │   └── style.css                   # 🎨 Tailwind CSS + custom styles
│   │
│   ├── 📂 public/                      # Static assets
│   │
│   ├── index.html                      # 📄 HTML template + PWA meta tags
│   ├── vite.config.js                  # ⚡ Vite + PWA configuration
│   ├── tailwind.config.js              # 🎨 Tailwind configuration
│   ├── postcss.config.js               # 🔧 PostCSS configuration
│   ├── package.json                    # 📦 Dependencies & scripts
│   └── yarn.lock                       # 🔒 Dependencies lock
│
├── 📄 Makefile                         # 🛠️ Development commands
│
├── 📚 Documentation Files:
│   ├── README.md                       # Main documentation
│   ├── QUICKSTART.md                   # Quick start (5 menit)
│   ├── SETUP_GUIDE.md                  # Detailed setup guide
│   ├── API_DOCUMENTATION.md            # Complete API docs
│   ├── DEPLOYMENT.md                   # Production deployment
│   ├── CHECKLIST.md                    # Setup checklist
│   ├── PROJECT_SUMMARY.md              # Project summary
│   ├── ARCHITECTURE_EXPLAINED.md       # Package explanations
│   └── FOLDER_STRUCTURE.md             # This file
│
├── .env                                # 🔐 Root environment (optional)
└── .gitignore                          # 🚫 Git ignore rules
```

---

## 📂 Backend Folder Details

### `cmd/server/`
**Purpose:** Entry point aplikasi  
**Files:** `main.go`  
**Responsibility:**
- Initialize configuration
- Connect to database
- Setup routes
- Start server

### `config/`
**Purpose:** Configuration management  
**Files:** `config.go`  
**Responsibility:**
- Load environment variables
- Provide config struct
- Default values

### `database/`
**Purpose:** Database setup & connection  
**Files:** `database.go`, `setup.sql`  
**Responsibility:**
- Connect to MySQL
- Auto migration
- Database utilities

### `handlers/`
**Purpose:** HTTP request handlers  
**Files:** `example_handler.go`  
**Responsibility:**
- Handle HTTP requests
- Validate input
- Call services
- Return responses

### `middleware/`
**Purpose:** HTTP middleware  
**Files:** `cors.go`  
**Responsibility:**
- CORS handling
- Authentication (future)
- Logging (future)

### `models/`
**Purpose:** Database models  
**Files:** `example.go`  
**Responsibility:**
- Define database structure
- GORM models
- Validation rules

### `routes/`
**Purpose:** Route definitions  
**Files:** `routes.go`  
**Responsibility:**
- Define API endpoints
- Group routes
- Apply middleware

### `services/`
**Purpose:** Business logic  
**Files:** `example_service.go`  
**Responsibility:**
- Business logic
- Data processing
- Database operations

---

## 📂 Frontend Folder Details

### `src/components/`
**Purpose:** Reusable Vue components  
**Files:** `ExampleCard.vue`, `ExampleForm.vue`  
**Responsibility:**
- UI components
- Reusable across pages
- Props & events

### `src/composables/`
**Purpose:** Reusable logic (Composition API)  
**Files:** `useApi.js`, `useExamples.js`  
**Responsibility:**
- API calls
- State management
- Business logic
- Reusable functions

### `src/views/`
**Purpose:** Page views  
**Files:** `Home.vue`  
**Responsibility:**
- Full page components
- Compose multiple components
- Page-level logic

### `src/App.vue`
**Purpose:** Root component  
**Responsibility:**
- App wrapper
- Global layout
- Route rendering

### `src/main.js`
**Purpose:** Application entry point  
**Responsibility:**
- Create Vue app
- Register plugins
- Register Service Worker (PWA)
- Mount app

### `src/style.css`
**Purpose:** Global styles  
**Responsibility:**
- Tailwind directives
- Custom utility classes
- Global CSS

---

## 🔄 Data Flow

### Request Flow (Backend):
```
HTTP Request
    ↓
routes.go (routing)
    ↓
example_handler.go (HTTP handling)
    ↓
example_service.go (business logic)
    ↓
example.go (model)
    ↓
database.go (database)
    ↓
MySQL
```

### Response Flow (Backend):
```
MySQL
    ↓
database.go
    ↓
example.go (model)
    ↓
example_service.go (process data)
    ↓
example_handler.go (format response)
    ↓
routes.go
    ↓
HTTP Response (JSON)
```

### Frontend Flow:
```
User Interaction
    ↓
Home.vue (view)
    ↓
ExampleForm.vue / ExampleCard.vue (components)
    ↓
useExamples.js (composable)
    ↓
useApi.js (API client)
    ↓
axios (HTTP request)
    ↓
Backend API
```

---

## 📊 File Count Summary

### Backend:
- **Go Files:** 8 files
- **Config Files:** 2 files (go.mod, .env)
- **Total:** 10 files

### Frontend:
- **Vue Files:** 4 files (.vue)
- **JavaScript Files:** 3 files (.js)
- **Config Files:** 5 files (vite.config, tailwind.config, etc)
- **Total:** 12+ files

### Documentation:
- **Markdown Files:** 9 files
- **Total:** 9 files

### Grand Total: ~31 files

---

## 🎯 Folder Organization Benefits

### ✅ Clear Separation
- Backend dan Frontend terpisah jelas
- Mudah di-navigate
- Mudah di-maintain

### ✅ Scalability
- Mudah tambah fitur baru
- Mudah tambah module
- Mudah refactor

### ✅ Team Collaboration
- Backend developer fokus di folder `backend/`
- Frontend developer fokus di folder `frontend/`
- Tidak bentrok

### ✅ Deployment
- Backend bisa di-deploy terpisah
- Frontend bisa di-deploy terpisah
- Atau deploy together

---

## 🚀 Quick Navigation

### Want to add new API endpoint?
```
1. Create model: backend/models/your_model.go
2. Create service: backend/services/your_service.go
3. Create handler: backend/handlers/your_handler.go
4. Add routes: backend/routes/routes.go
```

### Want to add new page?
```
1. Create view: frontend/src/views/YourPage.vue
2. Create components (if needed): frontend/src/components/
3. Create composable (if needed): frontend/src/composables/
```

### Want to modify styling?
```
1. Global styles: frontend/src/style.css
2. Component styles: Inside .vue files
3. Tailwind config: frontend/tailwind.config.js
```

---

## 📝 File Naming Conventions

### Backend (Go):
- **Files:** snake_case (example_handler.go)
- **Packages:** lowercase (handlers, services)
- **Structs:** PascalCase (ExampleHandler)
- **Functions:** PascalCase (NewExampleHandler)

### Frontend (Vue/JS):
- **Components:** PascalCase (ExampleCard.vue)
- **Composables:** camelCase with 'use' prefix (useApi.js)
- **Views:** PascalCase (Home.vue)
- **Variables:** camelCase (isLoading)

---

## 🔍 Find Files Quickly

### Backend:
```bash
# Find all Go files
find backend -name "*.go"

# Find specific file
find backend -name "*handler*"
```

### Frontend:
```bash
# Find all Vue files
find frontend/src -name "*.vue"

# Find all JS files
find frontend/src -name "*.js"
```

---

## 📚 Related Documentation

**Understanding the project:**
- **[ARCHITECTURE_EXPLAINED.md](./ARCHITECTURE_EXPLAINED.md)** - Why each package exists
- **[PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md)** - Complete overview

**Working with the structure:**
- **[CUSTOMIZATION_GUIDE.md](./CUSTOMIZATION_GUIDE.md)** - Add new files/features
- **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)** - API endpoints

**Setup:**
- **[QUICKSTART.md](./QUICKSTART.md)** - Quick setup
- **[SETUP_GUIDE.md](./SETUP_GUIDE.md)** - Detailed setup

---

**Developer:** Zulfikar Hidayatullah  
**Date:** 27 Desember 2025  
**Version:** 1.0.0
