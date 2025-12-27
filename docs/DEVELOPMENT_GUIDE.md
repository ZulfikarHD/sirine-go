# Development Guide - Sirine Go

## Development Mode Setup

Dalam development mode, aplikasi berjalan dengan dua server terpisah:

### 1. Backend Server (Go/Gin)
- **Port**: `8080`
- **URL**: `http://10.30.11.65:8080` atau `http://localhost:8080`
- **Fungsi**: Menyediakan REST API endpoints

### 2. Frontend Server (Vite Dev Server)
- **Port**: `5173`
- **URL**: `http://10.30.11.65:5173` atau `http://localhost:5173`
- **Fungsi**: Serve Vue 3 application dengan hot-reload

## Cara Menjalankan

### Menjalankan Kedua Server

```bash
# Terminal 1 - Backend
make dev-backend

# Terminal 2 - Frontend
make dev-frontend
```

### Atau Jalankan Bersamaan (Background)

```bash
# Jalankan backend di background
make dev-backend &

# Jalankan frontend di background
make dev-frontend &
```

## Akses Aplikasi

### ✅ Development Mode (Saat ini)

| Service | URL | Keterangan |
|---------|-----|------------|
| **Frontend UI** | `http://10.30.11.65:5173` | Akses aplikasi Vue melalui Vite dev server |
| **Backend API** | `http://10.30.11.65:8080/api/*` | REST API endpoints |
| **Health Check** | `http://10.30.11.65:8080/health` | Cek status backend |
| **API Info** | `http://10.30.11.65:8080/` | Informasi API dan development mode |

### ❌ Kesalahan Umum

**JANGAN akses `http://10.30.11.65:8080/` untuk melihat UI!**

Dalam development mode:
- Port `8080` = Backend API saja (JSON responses)
- Port `5173` = Frontend UI dengan hot-reload

## Frontend Configuration

Frontend sudah dikonfigurasi untuk berkomunikasi dengan backend API:

```javascript
// frontend/src/config/api.js
const API_BASE_URL = 'http://10.30.11.65:8080/api'
```

CORS sudah diaktifkan di backend untuk menerima request dari Vite dev server.

## Testing Endpoints

### Test Backend API

```bash
# Health check
curl http://10.30.11.65:8080/health

# API info
curl http://10.30.11.65:8080/

# Login endpoint
curl -X POST http://10.30.11.65:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"nip":"12345","password":"password"}'
```

### Test Frontend

Buka browser dan akses:
```
http://10.30.11.65:5173
```

## Production Mode

Untuk production, build frontend terlebih dahulu:

```bash
# Build frontend
make build-frontend

# Jalankan backend (akan serve static files dari dist/)
make run-backend
```

Dalam production mode:
- Hanya satu server (port 8080)
- Backend serve static files + API
- Akses UI dan API dari port yang sama

## Troubleshooting

### Error: "404 Not Found" saat akses port 8080

**Penyebab**: Mencoba akses UI dari port backend

**Solusi**: Akses frontend melalui `http://10.30.11.65:5173`

### Error: "Failed to run dependency scan"

**Penyebab**: Package motion-v belum terinstall atau salah package

**Solusi**:
```bash
cd frontend
yarn remove @motionone/vue
yarn add motion-v
```

### Error: "Address already in use"

**Penyebab**: Port sudah digunakan oleh process lain

**Solusi**:
```bash
# Kill process di port 8080
lsof -ti:8080 | xargs kill -9

# Kill process di port 5173
lsof -ti:5173 | xargs kill -9
```

### Backend tidak bisa connect ke database

**Cek konfigurasi `.env`**:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=sirine_go
```

## Environment Variables

Pastikan file `.env` ada di root project dengan konfigurasi yang benar:

```env
# Server
SERVER_PORT=8080
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=sirine_go

# JWT
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_EXPIRY=15m
JWT_REFRESH_EXPIRY=720h

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://10.30.11.65:5173
```

## Hot Reload

### Frontend (Vite)
- ✅ Otomatis hot-reload saat file `.vue`, `.js`, `.css` berubah
- ✅ Instant update tanpa refresh browser

### Backend (Go)
- ❌ Tidak ada hot-reload bawaan
- Manual restart: `Ctrl+C` lalu `make dev-backend` lagi
- Atau gunakan tools seperti `air` untuk auto-reload:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Jalankan dengan air
cd backend && air
```

## Network Access

Untuk akses dari device lain di network yang sama:

1. Pastikan firewall allow port 8080 dan 5173
2. Gunakan IP address mesin development (contoh: `10.30.11.65`)
3. Update CORS di backend jika perlu

## Next Steps

- [ ] Setup database seeding untuk development data
- [ ] Implement hot-reload untuk backend dengan Air
- [ ] Setup Docker untuk consistent development environment
- [ ] Add development tools (debugging, profiling)
