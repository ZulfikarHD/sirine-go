# ⚙️ Configuration Guide

## Overview

Configuration Guide ini menjelaskan cara mengonfigurasi aplikasi Sirine Go untuk berbagai environment (Development, Staging, Production). Konfigurasi utama dikelola melalui environment variables (`.env`).

## 📂 Configuration Files

### Backend
Lokasi: `backend/.env`
Template: `backend/.env.example`

### Frontend
Lokasi: `frontend/.env`
Template: `frontend/.env.example`

## 📝 Backend Configuration

### Server
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_PORT` | Port tempat backend server berjalan | `8080` | No |
| `GIN_MODE` | Mode operasi Gin (`debug`, `release`, `test`) | `debug` | No |
| `TRUSTED_PROXIES` | IP proxy yang dipercaya (untuk production) | `127.0.0.1` | No |

### Database
| Variable | Description | Example | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | Host database MySQL | `localhost` | Yes |
| `DB_PORT` | Port database MySQL | `3306` | Yes |
| `DB_USER` | Username database | `root` | Yes |
| `DB_PASSWORD` | Password database | `secret` | Yes |
| `DB_NAME` | Nama database | `sirine_go` | Yes |

### Authentication (JWT)
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `JWT_SECRET` | Secret key untuk signing JWT | (random string) | Yes |
| `JWT_EXPIRY` | Masa berlaku access token | `15m` | No |
| `REFRESH_TOKEN_EXPIRY` | Masa berlaku refresh token | `720h` (30 days) | No |
| `BCRYPT_COST` | Cost factor untuk password hashing | `12` | No |

### Security
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `MAX_LOGIN_ATTEMPTS` | Maksimal gagal login sebelum lockout | `5` | No |
| `LOCKOUT_DURATION` | Durasi akun terkunci | `15m` | No |

### Frontend URL (CORS)
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `FRONTEND_URL` | URL frontend yang diizinkan akses API | `http://localhost:5173` | Yes |

---

## 💻 Frontend Configuration

### API Connection
| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_BASE_URL` | URL base endpoint backend API | `http://localhost:8080` |

### App Settings
| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_APP_NAME` | Nama aplikasi | `Sirine Go` |
| `VITE_TIMEZONE` | Timezone default aplikasi | `Asia/Jakarta` |

## 🚀 Environment Setup

### Development
Gunakan `debug` mode untuk melihat logs detail dan error stack traces.
```env
GIN_MODE=debug
```

### Production
Wajib gunakan `release` mode untuk performa optimal.
```env
GIN_MODE=release
JWT_SECRET=(gunakan string panjang acak dan aman)
DB_PASSWORD=(password kuat)
```

## 🔄 Loading Configuration

Backend menggunakan package `godotenv` untuk load file `.env` saat startup. Code konfigurasi ada di `backend/config/config.go`.

```go
// Contoh akses config di code
dbHost := config.GetEnv("DB_HOST", "localhost")
```

## ⚠️ Important Notes

1. **Never commit `.env` files**: Pastikan `.env` ada di `.gitignore`.
2. **Sync `.env.example`**: Jika menambah variable baru, update juga `.env.example`.
3. **Restart Required**: Perubahan pada `.env` memerlukan restart aplikasi agar efektif.
