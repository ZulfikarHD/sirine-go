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

### Email Configuration (Sprint 3 - Password Reset)
| Variable | Description | Example | Required |
|----------|-------------|---------|----------|
| `EMAIL_SMTP_HOST` | SMTP server host | `smtp.gmail.com` | Yes* |
| `EMAIL_SMTP_PORT` | SMTP server port | `587` | Yes* |
| `EMAIL_USERNAME` | Email account username | `your-email@gmail.com` | Yes* |
| `EMAIL_PASSWORD` | Email account password atau app password | `your-app-password` | Yes* |
| `EMAIL_FROM_ADDRESS` | Email sender address | `noreply@sirine.local` | Yes* |

**\*Required** jika fitur Forgot Password diaktifkan. Untuk development, bisa gunakan [Mailtrap](https://mailtrap.io/) atau [MailHog](https://github.com/mailhog/MailHog).

### File Upload Configuration (Sprint 5)
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `UPLOAD_DIR` | Directory untuk simpan uploaded files | `./public/uploads/profiles` | No |
| `MAX_UPLOAD_SIZE` | Max file size dalam bytes | `5242880` (5MB) | No |
| `ALLOWED_FORMATS` | Allowed file extensions | `.jpg,.jpeg,.png,.webp` | No |

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

### File Upload (Sprint 5)
| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_MAX_PHOTO_SIZE` | Max photo upload size (MB) | `5` |

### Feature Flags (Optional)
| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_ENABLE_GAMIFICATION` | Enable achievement system | `true` |
| `VITE_ENABLE_NOTIFICATIONS` | Enable notification system | `true` |
| `VITE_POLLING_INTERVAL` | Notification polling interval (ms) | `30000` (30s) |

## 🚀 Environment Setup

### Development

**Backend `.env`:**
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

# JWT
JWT_SECRET=your-dev-secret-key-min-32-chars
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=720h
BCRYPT_COST=12

# Security
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=15m

# Email (Development - Mailtrap)
EMAIL_SMTP_HOST=smtp.mailtrap.io
EMAIL_SMTP_PORT=2525
EMAIL_USERNAME=your_mailtrap_username
EMAIL_PASSWORD=your_mailtrap_password
EMAIL_FROM_ADDRESS=noreply@sirine.local

# File Upload
UPLOAD_DIR=./public/uploads/profiles
MAX_UPLOAD_SIZE=5242880
ALLOWED_FORMATS=.jpg,.jpeg,.png,.webp

# CORS
FRONTEND_URL=http://localhost:5173
```

**Frontend `.env`:**
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_NAME=Sirine Go
VITE_TIMEZONE=Asia/Jakarta
VITE_MAX_PHOTO_SIZE=5
VITE_ENABLE_GAMIFICATION=true
VITE_ENABLE_NOTIFICATIONS=true
VITE_POLLING_INTERVAL=30000
```

---

### Production

**Backend `.env` (Production):**
```env
# Server
SERVER_PORT=8080
GIN_MODE=release
TRUSTED_PROXIES=your_proxy_ip

# Database
DB_HOST=your_db_host
DB_PORT=3306
DB_USER=sirine_user
DB_PASSWORD=strong_random_password_here
DB_NAME=sirine_go_production

# JWT (IMPORTANT: Generate secure random strings)
JWT_SECRET=super-long-random-secret-key-at-least-64-chars-for-production-security
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=720h
BCRYPT_COST=12

# Security
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=15m

# Email (Production - Gmail/SendGrid)
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_USERNAME=noreply@yourdomain.com
EMAIL_PASSWORD=your_app_specific_password
EMAIL_FROM_ADDRESS=noreply@yourdomain.com

# File Upload
UPLOAD_DIR=/var/www/sirine-go/public/uploads/profiles
MAX_UPLOAD_SIZE=5242880
ALLOWED_FORMATS=.jpg,.jpeg,.png,.webp

# CORS
FRONTEND_URL=https://yourdomain.com
```

**Frontend `.env.production`:**
```env
VITE_API_BASE_URL=https://api.yourdomain.com
VITE_APP_NAME=Sirine Go
VITE_TIMEZONE=Asia/Jakarta
VITE_MAX_PHOTO_SIZE=5
VITE_ENABLE_GAMIFICATION=true
VITE_ENABLE_NOTIFICATIONS=true
VITE_POLLING_INTERVAL=30000
```

### Important Production Notes

1. **JWT_SECRET**: Gunakan strong random string minimal 64 karakter
   ```bash
   # Generate dengan openssl
   openssl rand -base64 64
   ```

2. **Database Password**: Gunakan password yang kuat dan unique

3. **Email Configuration**: 
   - Gmail: Gunakan [App Passwords](https://support.google.com/accounts/answer/185833)
   - SendGrid/Mailgun: Lebih reliable untuk production
   - Setup SPF/DKIM records untuk prevent spam filtering

4. **File Permissions**:
   ```bash
   mkdir -p /var/www/sirine-go/public/uploads/profiles
   chown -R www-data:www-data /var/www/sirine-go/public/uploads
   chmod -R 755 /var/www/sirine-go/public/uploads
   ```

5. **CORS**: Pastikan `FRONTEND_URL` match dengan domain production Anda

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
