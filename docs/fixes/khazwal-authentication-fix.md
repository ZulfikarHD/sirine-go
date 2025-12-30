# Fix: Khazwal Authentication Error

**Date:** 30 Desember 2025  
**Issue:** User authentication error saat mencoba memulai persiapan material  
**Status:** ✅ Fixed

---

## Problem

User mendapatkan error "User tidak terautentikasi" saat klik "Mulai Persiapan" di halaman Material Prep Detail, meskipun sudah login dengan akun superadmin.

### Error Message
```
Gagal Memulai Persiapan
User tidak terautentikasi
```

---

## Root Cause

**Case sensitivity mismatch** antara auth middleware dan handler:

- **Auth Middleware** (`backend/middleware/auth_middleware.go:66`):
  - Menyimpan data dengan key: `user_id` (lowercase dengan underscore)
  
  ```go
  c.Set("user_id", user.ID)
  c.Set("user_role", user.Role)
  ```

- **Khazwal Handler** (`backend/handlers/khazwal_handler.go`):
  - Mencoba mengambil data dengan key: `userID` (camelCase)
  
  ```go
  userIDInterface, exists := c.Get("userID")  // ❌ SALAH
  ```

Karena key tidak cocok, `exists` akan selalu `false`, sehingga handler mengembalikan error "User tidak terautentikasi".

---

## Solution

Mengubah semua `c.Get("userID")` menjadi `c.Get("user_id")` di file `khazwal_handler.go`:

### Files Changed

**`backend/handlers/khazwal_handler.go`**

Fixed 5 functions yang terpengaruh:

1. **StartPrep** (line ~175)
2. **ConfirmPlat** (line ~255)
3. **UpdateKertas** (line ~334)
4. **UpdateTinta** (line ~413)
5. **Finalize** (line ~492)

### Change Pattern

**Before:**
```go
userIDInterface, exists := c.Get("userID")  // ❌
```

**After:**
```go
userIDInterface, exists := c.Get("user_id")  // ✅
```

---

## Verification

### 1. Restart Backend Server

```bash
cd backend
go run cmd/server/main.go
```

### 2. Test Flow

1. Login dengan akun superadmin atau staff_khazwal
2. Navigate ke `/khazwal/material-prep`
3. Pilih salah satu PO dari queue
4. Klik "Mulai Persiapan"
5. Confirm di dialog
6. Pastikan berhasil tanpa error authentication

### Expected Result

- ✅ Status berubah dari "Menunggu Persiapan" → "Sedang Dipersiapkan"
- ✅ Success alert muncul: "Persiapan Material Dimulai"
- ✅ Timeline tracking tercatat dengan user yang benar
- ✅ Button "Lanjutkan Proses" muncul

---

## Prevention

### Consistency Rules

1. **Auth Middleware** menggunakan **snake_case** untuk context keys:
   - `user_id`
   - `user_role`
   - `user`
   - `claims`

2. **Handlers** harus menggunakan key yang **sama persis**:
   ```go
   // ✅ BENAR
   c.Get("user_id")
   c.Get("user_role")
   
   // ❌ SALAH
   c.Get("userID")
   c.Get("userRole")
   ```

3. **Code Review Checklist**:
   - [ ] Pastikan case consistency antara middleware dan handler
   - [ ] Test authentication untuk semua protected endpoints
   - [ ] Verifikasi error handling untuk unauthenticated requests

---

## Related Files

- `backend/middleware/auth_middleware.go`
- `backend/handlers/khazwal_handler.go`
- `backend/routes/routes.go`
- `frontend/src/views/khazwal/MaterialPrepDetailPage.vue`

---

## Notes

- Issue ini hanya mempengaruhi Khazwal handlers
- Handlers lain (user, profile, achievement, dll) sudah menggunakan `user_id` yang benar
- Tidak ada perubahan pada frontend atau middleware
- Tidak perlu migration atau database changes
