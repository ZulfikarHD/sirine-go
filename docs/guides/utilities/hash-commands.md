# Hash Commands Guide

## Overview

Hash Commands merupakan utility tools untuk password management yang memudahkan generate dan verify bcrypt password hashes, yaitu: command-line interface yang simple dan powerful untuk development dan maintenance.

## Available Commands

### `hash:make` - Generate Password Hash

Generate bcrypt hash untuk password baru dengan cost factor 12 (recommended untuk production).

#### Usage via Makefile

```bash
make hash-make PASSWORD="YourPassword"
```

#### Usage Direct (Go Run)

```bash
cd backend && go run cmd/hash/main.go make "YourPassword"
```

#### Examples

**Generate hash untuk admin password:**
```bash
make hash-make PASSWORD="Admin@123"
```

**Output:**
```
✅ Hash generated successfully!

Password: Admin@123
Hash:     $2a$12$nNag4FTJB0fiX/f22aINOuYcuP8cUOWOyCZub6tBCom0Evxv4ahTK

Copy hash di atas untuk digunakan di seeder atau database.
```

**Generate hash untuk demo user:**
```bash
make hash-make PASSWORD="Demo@123"
```

### `hash:check` - Verify Password

Verify apakah plaintext password match dengan bcrypt hash yang ada.

#### Usage Direct (Recommended)

```bash
cd backend && go run cmd/hash/main.go check "Password" '$2a$12$...'
```

**Note:** Gunakan **single quotes** untuk hash karena bcrypt hash mengandung karakter `$` yang perlu di-escape di shell.

#### Examples

**Verify password match:**
```bash
cd backend && go run cmd/hash/main.go check "Admin@123" '$2a$12$nNag4FTJB0fiX/f22aINOuYcuP8cUOWOyCZub6tBCom0Evxv4ahTK'
```

**Output (Success):**
```
✅ Password MATCH!

Password: Admin@123
Hash:     $2a$12$nNag4FTJB0fiX/f22aINOuYcuP8cUOWOyCZub6tBCom0Evxv4ahTK

Verifikasi berhasil - password sesuai dengan hash.
```

**Output (Failed):**
```
❌ Password TIDAK MATCH!

Password: WrongPassword
Hash:     $2a$12$nNag4FTJB0fiX/f22aINOuYcuP8cUOWOyCZub6tBCom0Evxv4ahTK

Verifikasi gagal - password tidak sesuai dengan hash.
```

## Use Cases

### 1. Development & Testing

Saat development, generate hash untuk test users:

```bash
# Generate hash untuk multiple test users
make hash-make PASSWORD="Admin@123"
make hash-make PASSWORD="Manager@123"
make hash-make PASSWORD="Staff@123"
```

Copy hashes yang di-generate ke seeder atau test fixtures.

### 2. Debug Authentication Issues

Jika ada masalah login, verify hash yang tersimpan di database:

```bash
# Get hash dari database (via script atau manual)
# Then verify dengan password yang digunakan

cd backend && go run cmd/hash/main.go check "Admin@123" '$2a$12$...'
```

### 3. Manual User Creation

Saat perlu create user secara manual di database:

```bash
# 1. Generate hash
make hash-make PASSWORD="NewUser@123"

# 2. Copy hash ke SQL INSERT statement
INSERT INTO users (nip, full_name, email, password_hash, ...)
VALUES ('12345', 'New User', 'user@example.com', '$2a$12$...', ...);
```

### 4. Password Policy Testing

Test berbagai password formats untuk validate policy:

```bash
make hash-make PASSWORD="Short123"        # Test min length
make hash-make PASSWORD="NOLOWERCASE123"  # Test lowercase requirement
make hash-make PASSWORD="nouppercase123"  # Test uppercase requirement
make hash-make PASSWORD="NoSpecial123"    # Test special char requirement
```

## Technical Details

### Bcrypt Cost Factor

Commands menggunakan **cost factor 12** yang merupakan:
- **Recommended** untuk production environments
- **Balance** antara security dan performance
- Generate hash membutuhkan ~250-300ms per password

### Hash Format

Bcrypt hash memiliki format:
```
$2a$12$SSSSSSSSSSSSSSSSSSSSSSHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH
 │  │  │                     │
 │  │  └─ 22-char salt       └─ 31-char hash value
 │  └─ Cost factor (12)
 └─ Algorithm version (2a)
```

### Security Considerations

1. **Salt**: Setiap hash memiliki unique salt yang di-generate otomatis
2. **Deterministic**: Same password + same salt = same hash
3. **One-way**: Tidak bisa di-reverse (hanya bisa verify)
4. **Slow by design**: Cost factor membuat brute force attacks expensive

## Integration with Seeder

Hash commands terintegrasi dengan seeding system:

```go
// backend/cmd/seed/main.go
passwordService := services.NewPasswordService()

// Generate hash secara dynamic
adminHash, _ := passwordService.HashPassword("Admin@123")

adminUser := models.User{
    NIP:          "99999",
    PasswordHash: adminHash,  // Dynamic hash
    // ...
}
```

**Benefits:**
- Fresh hash setiap seed (different salt)
- No hardcoded hashes di code
- Easy to update passwords

## Troubleshooting

### Error: Password tidak diberikan

```bash
❌ Error: PASSWORD tidak diberikan

Usage: make hash-make PASSWORD="YourPassword"
```

**Solution:** Provide PASSWORD parameter:
```bash
make hash-make PASSWORD="YourPassword"
```

### Error: Hash escaping di Makefile

Bcrypt hash mengandung `$` yang di-interpret sebagai variable oleh shell.

**Solution:** Use direct go run dengan single quotes:
```bash
cd backend && go run cmd/hash/main.go check "Password" '$2a$12$...'
```

### Error: Failed to generate hash

Jika generate hash gagal, check:
1. Password tidak kosong
2. Bcrypt library ter-install (`go mod download`)
3. Tidak ada memory issues (cost 12 butuh ~10MB RAM per hash)

## Best Practices

### Development

1. **Generate hashes fresh** untuk setiap environment
   ```bash
   make hash-make PASSWORD="DevPassword@123"
   ```

2. **Don't commit** plain passwords ke git
   ```bash
   # Good: Store in .env (gitignored)
   ADMIN_PASSWORD=Admin@123
   
   # Bad: Hardcode di code
   password := "Admin@123"  // Don't do this!
   ```

3. **Test hash verification** sebelum deploy
   ```bash
   cd backend && go run cmd/hash/main.go check "Password" '$hash'
   ```

### Production

1. **Generate unique hashes** untuk production
   ```bash
   make hash-make PASSWORD="StrongProductionPassword!@#$"
   ```

2. **Store securely** - never log or print hashes
   ```go
   // Bad
   log.Printf("User hash: %s", user.PasswordHash)
   
   // Good
   log.Printf("User %s authenticated", user.Email)
   ```

3. **Rotate regularly** - change passwords periodically
   ```bash
   # Generate new hash
   make hash-make PASSWORD="NewRotatedPassword123"
   
   # Update di database atau seeder
   ```

## Performance Considerations

### Cost Factor Impact

| Cost | Time per Hash | Security Level |
|------|---------------|----------------|
| 10   | ~80ms        | Minimum        |
| **12**   | **~250ms**       | **Recommended**   |
| 14   | ~1000ms      | High Security  |
| 16   | ~4000ms      | Maximum        |

**Our choice: Cost 12**
- Good balance untuk web applications
- Fast enough untuk user experience
- Slow enough untuk security

### Batch Hashing

Untuk generate banyak hashes:

```bash
# Create script
for password in "User1@123" "User2@123" "User3@123"; do
    echo "Generating hash for $password"
    cd backend && go run cmd/hash/main.go make "$password"
done
```

## Quick Reference

```bash
# Generate hash
make hash-make PASSWORD="YourPassword"

# Verify hash (direct)
cd backend && go run cmd/hash/main.go check "Password" '$hash'

# Show help
cd backend && go run cmd/hash/main.go

# All commands
make help | grep hash
```

## Related Documentation

- [Database Management Guide](./DATABASE_MANAGEMENT.md) - Database seeding dengan hashed passwords
- [Password Service](../backend/services/password_service.go) - Password hashing implementation
- [User Model](../backend/models/user.go) - User authentication model

## Support

Untuk pertanyaan atau issues, hubungi:
- Developer: Zulfikar Hidayatullah (+62 857-1583-8733)
