# 🚀 Backend Getting Started

Panduan awal untuk memulai backend development dalam Sirine Go App, yaitu setup environment, instalasi dependencies, dan menjalankan aplikasi.

---

## 📋 Prerequisites

### Required Software

Sebelum memulai development, pastikan sudah terinstal:

| Software | Version | Check Command | Purpose |
|----------|---------|---------------|---------|
| **Go** | 1.21+ | `go version` | Backend programming language |
| **MySQL** | 8.0+ | `mysql --version` | Database server |
| **Git** | Latest | `git --version` | Version control |

### Optional Tools

| Tool | Purpose | Installation |
|------|---------|--------------|
| **Air** | Hot reload untuk development | `go install github.com/air-verse/air@latest` |
| **MySQL Workbench** | GUI untuk database management | [Download](https://dev.mysql.com/downloads/workbench/) |

### Knowledge Prerequisites

Disarankan memiliki pengetahuan dasar tentang:
- Go programming language syntax
- HTTP/REST API concepts
- SQL dan relational database
- JSON format

---

## 🔧 Setup Development Environment

### Step 1: Clone Repository

```bash
# Clone project
git clone <repository-url>
cd sirine-go

# Navigate ke backend folder
cd backend
```

### Step 2: Install Dependencies

```bash
# Download Go modules
go mod download

# Verify dependencies
go mod verify
```

### Step 3: Setup Database

```bash
# Login ke MySQL
mysql -u root -p

# Create database
CREATE DATABASE sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# Create database user (optional, untuk security)
CREATE USER 'sirine_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON sirine_go.* TO 'sirine_user'@'localhost';
FLUSH PRIVILEGES;

# Exit MySQL
EXIT;
```

### Step 4: Configure Environment

```bash
# Copy environment template
cp .env.example .env

# Edit dengan credentials Anda
nano .env
```

**Konfigurasi minimal `.env`:**
```env
# Server
SERVER_PORT=8080
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=sirine_go

# JWT Secret (generate dengan: openssl rand -base64 32)
JWT_SECRET=your-secret-key-min-32-chars

# Frontend URL untuk CORS
FRONTEND_URL=http://localhost:5173
```

### Step 5: Run Migrations

```bash
# Jalankan database migrations
go run cmd/migrate/main.go
```

**Output yang diharapkan:**
```
🔄 Running migrations...
✅ Migration completed: 001_create_users_table
✅ Migration completed: 002_create_achievements_table
✅ Migration completed: 003_create_notifications_table
✅ All migrations completed successfully!
```

### Step 6: Seed Database (Optional)

```bash
# Populate database dengan sample data
go run cmd/seed/main.go
```

**Sample data yang dibuat:**
- Admin user: `admin` / `Admin@123`
- Sample departments dan shifts
- Initial achievements

### Step 7: Start Development Server

```bash
# Dengan hot reload (recommended)
air

# Atau tanpa hot reload
go run cmd/server/main.go
```

**Output yang diharapkan:**
```
[GIN-debug] Listening and serving HTTP on :8080
✅ Server running on http://localhost:8080
```

---

## 🎯 Development Commands

### Server Commands

```bash
# Start server dengan hot reload
air

# Start server tanpa hot reload
go run cmd/server/main.go

# Build binary untuk production
go build -o bin/server cmd/server/main.go

# Run production binary
./bin/server
```

### Database Commands

```bash
# Run migrations
go run cmd/migrate/main.go

# Seed database
go run cmd/seed/main.go

# Reset database (drop & recreate)
mysql -u root -p sirine_go < migrations/reset.sql
go run cmd/migrate/main.go
go run cmd/seed/main.go
```

### Testing Commands

```bash
# Run all tests
go test ./... -v

# Run tests dengan coverage
go test ./... -v -cover

# Run specific package tests
go test ./internal/services -v

# Run specific test function
go test ./internal/services -v -run TestCreateUser
```

### Code Quality Commands

```bash
# Format code
go fmt ./...

# Vet code untuk common mistakes
go vet ./...

# Run linter (install golangci-lint first)
golangci-lint run
```

---

## 📁 Project Structure

```
backend/
├── cmd/                        # Command-line executables
│   ├── server/                 # Main server entry point
│   │   └── main.go
│   ├── migrate/                # Database migration runner
│   │   └── main.go
│   └── seed/                   # Database seeder
│       └── main.go
│
├── internal/                   # Private application code
│   ├── config/                 # Configuration management
│   │   └── config.go           # Load .env dan app config
│   │
│   ├── database/               # Database connection
│   │   └── database.go         # GORM setup dan connection
│   │
│   ├── middleware/             # HTTP middlewares
│   │   ├── auth.go             # JWT authentication
│   │   ├── cors.go             # CORS handling
│   │   └── logger.go           # Request logging
│   │
│   ├── models/                 # Database models (GORM)
│   │   ├── user.go
│   │   ├── achievement.go
│   │   ├── notification.go
│   │   └── activity_log.go
│   │
│   ├── handlers/               # HTTP handlers (controllers)
│   │   ├── auth_handler.go     # Login, logout, refresh
│   │   ├── user_handler.go     # User CRUD
│   │   └── profile_handler.go  # Profile management
│   │
│   ├── services/               # Business logic layer
│   │   ├── auth_service.go     # Auth business logic
│   │   ├── user_service.go     # User business logic
│   │   └── profile_service.go  # Profile business logic
│   │
│   ├── repositories/           # Data access layer (optional)
│   │   └── user_repository.go  # Complex queries
│   │
│   ├── utils/                  # Utility functions
│   │   ├── jwt.go              # JWT generation & validation
│   │   ├── hash.go             # Password hashing
│   │   ├── response.go         # Response helpers
│   │   └── validator.go        # Input validation
│   │
│   └── routes/                 # Route definitions
│       └── routes.go           # API route setup
│
├── migrations/                 # SQL migration files
│   ├── 001_create_users_table.sql
│   ├── 002_create_achievements_table.sql
│   └── 003_create_notifications_table.sql
│
├── public/                     # Static files & uploads
│   └── uploads/
│       └── profiles/           # Profile photos
│
├── .env                        # Environment variables (git-ignored)
├── .env.example                # Environment template
├── .air.toml                   # Air configuration untuk hot reload
├── go.mod                      # Go dependencies
├── go.sum                      # Dependencies checksums
└── README.md                   # Backend documentation
```

### Folder Structure Explanation

**`cmd/`** - Entry points aplikasi, yaitu:
- Setiap subfolder berisi `main.go` yang dapat dijalankan
- Minimal logic, hanya setup dan start aplikasi

**`internal/`** - Core application code yang private:
- Code di folder ini tidak bisa di-import oleh external projects
- Berisi semua business logic dan infrastructure code

**`migrations/`** - Database schema definitions:
- SQL files untuk create/alter tables
- Numbered prefix untuk execution order (001, 002, etc.)

---

## 🔍 Testing Your Setup

### 1. Check Server Health

```bash
# Test jika server running
curl http://localhost:8080/health
```

**Expected response:**
```json
{
  "status": "ok",
  "timestamp": "2025-12-28T10:30:00Z"
}
```

### 2. Test Database Connection

```bash
# Login ke MySQL dan check tables
mysql -u root -p sirine_go

# List tables
SHOW TABLES;
```

**Expected output:**
```
+----------------------+
| Tables_in_sirine_go  |
+----------------------+
| users                |
| achievements         |
| notifications        |
| user_achievements    |
| activity_logs        |
+----------------------+
```

### 3. Test API Endpoint

```bash
# Test login dengan seeded admin account
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "admin",
    "password": "Admin@123"
  }'
```

**Expected response:**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "user": {
      "id": 1,
      "nip": "admin",
      "full_name": "Administrator",
      "role": "SUPER_ADMIN"
    }
  }
}
```

---

## ⚠️ Troubleshooting

### Port Already in Use

```bash
# Error: address already in use
# Solution: Kill process pada port 8080
lsof -ti:8080 | xargs kill -9
```

### Database Connection Failed

```bash
# Error: Error 1045: Access denied
# Solution: Check credentials di .env
# Verify dengan:
mysql -u your_user -p -h localhost
```

### Missing Dependencies

```bash
# Error: cannot find package
# Solution: Re-download dependencies
go mod download
go mod tidy
```

### Air Not Found

```bash
# Error: air: command not found
# Solution: Install Air dan pastikan GOPATH di PATH
go install github.com/air-verse/air@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## 📚 Next Steps

Setelah setup berhasil, lanjutkan dengan:

1. **[Service Pattern Guide](./service-pattern.md)** - Pelajari arsitektur aplikasi
2. **[Middleware Guide](./middleware.md)** - Pahami authentication flow
3. **[API Reference](../../04-api-reference/README.md)** - Explore available endpoints
4. **[Database Models](../../05-guides/database/models.md)** - Understand data structure

---

**Last Updated:** 28 Desember 2025  
**Status:** ✅ Production Ready
