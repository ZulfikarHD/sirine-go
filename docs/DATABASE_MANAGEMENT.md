# Database Management Guide

## Overview

Database Management merupakan sistem pengelolaan database untuk Sirine Go App yang dijalankan melalui Makefile commands, yaitu: command-command yang menggunakan Go untuk handle migrations dan seeding tanpa memerlukan MySQL CLI.

## Prerequisites

Komponen yang diperlukan, antara lain:
- Go 1.21 atau lebih tinggi
- MySQL 8.0 atau lebih tinggi yang berjalan di localhost:3306
- File `.env` dengan konfigurasi database yang benar

## Quick Start

### 1. Setup Awal

Untuk setup database pertama kali:

```bash
# Buat database
make db-create

# Jalankan migrations
make db-migrate

# Seed dengan data awal
make db-seed
```

Atau gunakan command singkat:

```bash
# Reset lengkap (fresh + seed)
make db-reset
```

## Available Commands

### `make db-create`
Membuat database `sirine_go` dengan character set utf8mb4 dan collation utf8mb4_unicode_ci.

```bash
make db-create
```

**Output:**
```
✅ Database 'sirine_go' created successfully!
```

### `make db-drop`
Menghapus database `sirine_go`. **Peringatan:** Semua data akan hilang!

```bash
make db-drop
```

**Output:**
```
✅ Database 'sirine_go' dropped successfully!
```

### `make db-migrate`
Menjalankan GORM AutoMigrate untuk semua models, dimana akan membuat atau mengupdate table structure sesuai dengan model definitions.

```bash
make db-migrate
```

**Tables yang di-migrate:**
- `users` - User accounts dengan authentication
- `user_sessions` - JWT session tracking
- `password_reset_tokens` - Password reset tokens
- `activity_logs` - Audit trail untuk semua activities
- `notifications` - In-app notifications

**Output:**
```
Running migrations...
✅ Migrations completed successfully!
```

### `make db-rollback`
Menghapus semua tables (rollback migrations). Tables dihapus dalam urutan yang benar untuk menghindari foreign key constraints.

```bash
make db-rollback
```

**Output:**
```
Rolling back migrations...
Dropped table: notifications
Dropped table: activity_logs
Dropped table: password_reset_tokens
Dropped table: user_sessions
Dropped table: users
✅ Rollback completed successfully!
```

### `make db-fresh`
Drop database, create ulang, dan migrate. Command ini berguna untuk reset database structure tanpa data.

```bash
make db-fresh
```

**Workflow:**
1. Drop database
2. Create database
3. Run migrations

**Output:**
```
🔄 Starting fresh migration...
✅ Database 'sirine_go' dropped successfully!
✅ Database 'sirine_go' created successfully!
Running migrations...
✅ Migrations completed successfully!
✅ Fresh migration completed!
```

### `make db-seed`
Menjalankan database seeder untuk populate database dengan initial data.

```bash
make db-seed
```

**Data yang di-seed:**

1. **Admin User**
   - NIP: `99999`
   - Email: `admin@sirine.local`
   - Password: `Admin@123`
   - Role: ADMIN
   - Status: ACTIVE

2. **Demo Users** (untuk testing):
   - `10001` - Manager Produksi (MANAGER)
   - `20001` - Staff Khazanah Awal (STAFF_KHAZWAL)
   - `30001` - Operator Cetak (OPERATOR_CETAK)
   - `40001` - QC Inspector (QC_INSPECTOR)
   - `50001` - Verifikator (VERIFIKATOR)
   - `60001` - Staff Khazanah Akhir (STAFF_KHAZKHIR)
   - **Password untuk semua demo users:** `Demo@123`

**Output:**
```
🌱 Starting database seeding...
✅ Admin user seeded successfully
   NIP: 99999
   Email: admin@sirine.local
   Password: Admin@123
✅ Demo user seeded: 10001 - Manager Produksi
✅ Demo user seeded: 20001 - Staff Khazanah Awal
...
📝 Demo users credentials:
   Password untuk semua demo users: Demo@123
✅ Database seeding completed!
```

### `make db-reset`
Fresh migration + seed. Command paling sering digunakan untuk development.

```bash
make db-reset
```

**Workflow:**
1. Drop database
2. Create database
3. Run migrations
4. Seed data

**Output:**
```
🔄 Resetting database...
[output dari db-fresh]
[output dari db-seed]
✅ Database reset complete!
```

## Development Workflow

### Daily Development

Untuk development sehari-hari, gunakan `db-reset` untuk mendapatkan clean database dengan data fresh:

```bash
make db-reset
```

### Testing New Migrations

Jika Anda menambahkan field baru ke model:

```bash
# Update database structure
make db-migrate
```

### Rollback dan Fresh Start

Jika ada masalah dengan migrations:

```bash
# Rollback semua migrations
make db-rollback

# Kemudian migrate ulang
make db-migrate
```

Atau gunakan shortcut:

```bash
make db-fresh
```

## Architecture

### Migration System

Migration system menggunakan **Registry Pattern** dengan GORM AutoMigrate yang:
- **Centralized model registration** - Semua models di-register di satu tempat
- **Automatic migration** - Semua command auto-sync dengan registry
- Membuat tables otomatis dari struct definitions
- Menambahkan kolom baru yang belum ada
- Menambahkan indexes dan foreign keys
- **Tidak menghapus** kolom yang tidak ada di struct (safe)

**Registry Location:** `backend/database/models_registry.go`

Untuk menambahkan model/table baru, lihat: [ADDING_NEW_MODELS.md](./ADDING_NEW_MODELS.md)

### Seeding System

Seeding system menggunakan GORM untuk insert data dengan:
- Duplicate checking (tidak seed jika data sudah ada)
- Proper password hashing dengan bcrypt cost 12
- Foreign key handling yang benar

### File Structure

```
backend/
├── cmd/
│   ├── migrate/
│   │   └── main.go         # Migration commands
│   ├── seed/
│   │   └── main.go         # Seeding logic
│   └── genhash/
│       └── main.go         # Password hash generator
├── models/
│   ├── user.go            # User model
│   ├── user_session.go    # Session model
│   ├── password_reset_token.go
│   ├── activity_log.go
│   └── notification.go
└── database/
    ├── database.go        # Database connection
    └── setup.sql          # Legacy SQL setup (reference)
```

## Configuration

Database configuration dibaca dari environment variables dengan fallback ke default values untuk development:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=sirine_go
```

Default values (jika tidak ada di .env):
- Host: `localhost`
- Port: `3306`
- User: `root`
- Password: `` (empty)
- Database: `sirine_go`

## Utility Commands

### Generate Password Hash

Untuk generate bcrypt hash untuk password baru:

```bash
cd backend && go run cmd/genhash/main.go "YourPassword123"
```

**Output:**
```
Password: YourPassword123
Hash: $2a$12$...
```

## Troubleshooting

### Error: Cannot connect to database

**Problem:**
```
Failed to connect to MySQL: ...
```

**Solution:**
1. Pastikan MySQL server berjalan
2. Check credentials di file `.env`
3. Pastikan database user memiliki permissions yang cukup

### Error: Table already exists

**Problem:**
```
Error 1050: Table 'users' already exists
```

**Solution:**
Gunakan `db-fresh` untuk drop dan recreate tables:
```bash
make db-fresh
```

### Error: Foreign key constraint fails

**Problem:**
```
Error 1451: Cannot delete or update a parent row
```

**Solution:**
Gunakan `db-rollback` yang menghapus tables dalam urutan yang benar:
```bash
make db-rollback
```

### Error: Unknown column 'nip' in 'WHERE'

**Problem:**
Column name tidak sesuai dengan GORM naming convention.

**Solution:**
Sudah diperbaiki dengan menambahkan explicit column tag:
```go
NIP string `gorm:"column:nip;..."`
```

## Best Practices

### Development

1. **Selalu gunakan `make db-reset`** sebelum testing untuk clean state
2. **Jangan edit migration files** yang sudah di-commit
3. **Test migrations** sebelum push ke repository
4. **Backup data production** sebelum run migrations

### Production

1. **JANGAN gunakan `db-fresh` atau `db-reset`** di production
2. **Gunakan `db-migrate`** untuk update schema
3. **Backup database** sebelum migrate
4. **Test di staging** environment dulu
5. **Monitor logs** saat migration berjalan

### Security

1. **Ganti password default** untuk admin user di production
2. **Hapus demo users** di production
3. **Set strong password** dengan minimum requirements
4. **Enable SSL** untuk database connections di production

## Testing

Untuk verify bahwa database setup benar:

```bash
# Reset database
make db-reset

# Start backend server
make dev-backend

# Test login dengan admin credentials
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@sirine.local","password":"Admin@123"}'
```

Expected response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": { ... },
    "tokens": { ... }
  }
}
```

## Additional Resources

- [GORM Documentation](https://gorm.io/docs/)
- [MySQL Documentation](https://dev.mysql.com/doc/)
- [Bcrypt Package](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

## Support

Untuk pertanyaan atau issues, hubungi:
- Developer: Zulfikar Hidayatullah (+62 857-1583-8733)
- Email: admin@sirine.local
