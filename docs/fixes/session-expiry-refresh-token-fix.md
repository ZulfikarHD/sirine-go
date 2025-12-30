# Fix: Session Expiry Refresh Token Issue

## Masalah

User mengalami logout otomatis setiap X menit (sekitar 15 menit) meskipun masih aktif menggunakan aplikasi, dengan error message:

```
API Error: gagal refresh token - refresh token expired atau sudah di-revoke
API Error: refresh token expired atau sudah di-revoke
```

**Impact:** UX sangat buruk karena user harus login ulang setiap 15 menit.

## Root Cause

### Bug di `Login()` Function - Session Expiry Logic

**Lokasi:** `backend/services/auth_service.go` lines 128-134

**Kode Bermasalah:**
```go
// Calculate expiry based on RememberMe
var expiresAt time.Time
if req.RememberMe {
    expiresAt = time.Now().Add(s.config.RefreshTokenExpiry) // 30 days
} else {
    expiresAt = time.Now().Add(s.config.JWTExpiry) // 15 minutes ⚠️
}
```

**Penjelasan Masalah:**

1. **Access Token (JWT):** Expiry 15 menit (untuk security)
2. **Refresh Token:** Expiry 30 hari (untuk UX)
3. **Session di Database:** Menyimpan `expires_at` untuk validasi

**Flow yang Salah:**
```
User Login (tanpa Remember Me checkbox)
   ↓
Backend creates session dengan expires_at = NOW + 15 minutes
   ↓
User aktif selama 10 menit, access token expired
   ↓
Frontend interceptor coba refresh token
   ↓
Backend: POST /api/auth/refresh
   ↓
Backend validate session di database
   ↓
session.expires_at < NOW → EXPIRED ⚠️
   ↓
Response 401: "refresh token expired atau sudah di-revoke"
   ↓
Frontend logout user
   ↓
User harus login ulang setiap 15 menit
```

**Kesalahan Konsep:**

- **Session expiry TIDAK BOLEH mengikuti Access Token expiry**
- Session harus selalu mengikuti **Refresh Token expiry** (30 hari)
- `RememberMe` flag seharusnya untuk **extend** session, bukan untuk shorten it

## Solusi

### 1. Fix Session Expiry Logic

**File:** `backend/services/auth_service.go`

**Before:**
```go
var expiresAt time.Time
if req.RememberMe {
    expiresAt = time.Now().Add(s.config.RefreshTokenExpiry) // 30 days
} else {
    expiresAt = time.Now().Add(s.config.JWTExpiry) // 15 minutes ⚠️ BUG!
}
```

**After:**
```go
// Session expiry harus selalu mengikuti refresh token expiry (30 hari)
// untuk memastikan refresh token mechanism berfungsi dengan baik.
// RememberMe dapat digunakan untuk extend expiry lebih lama jika diperlukan.
var expiresAt time.Time
if req.RememberMe {
    // Extend session untuk remember me (misalnya 90 hari)
    expiresAt = time.Now().Add(s.config.RefreshTokenExpiry * 3)
} else {
    // Default session expiry mengikuti refresh token (30 hari)
    expiresAt = time.Now().Add(s.config.RefreshTokenExpiry)
}
```

### 2. Tambah Dokumentasi di `RefreshAuthToken()`

**File:** `backend/services/auth_service.go`

Menambahkan comment untuk clarity bahwa session expiry harus mengikuti refresh token expiry:

```go
// Update session dengan token dan refresh token baru
// Session expiry tetap mengikuti refresh token expiry (30 hari)
newTokenHash := hashToken(newToken)
newRefreshTokenHash := hashToken(newRefreshToken)

s.db.Model(&models.UserSession{}).
    Where("refresh_token_hash = ?", refreshTokenHash).
    Updates(map[string]interface{}{
        "token_hash":         newTokenHash,
        "refresh_token_hash": newRefreshTokenHash,
        "expires_at":         time.Now().Add(s.config.RefreshTokenExpiry),
    })
```

## Flow Setelah Fix

### Scenario 1: Normal Login (tanpa Remember Me)

```
User Login → RememberMe = false
   ↓
Backend creates session:
   - token_hash (access token hash)
   - refresh_token_hash
   - expires_at = NOW + 30 days ✅
   ↓
User aktif, access token expired setelah 15 menit
   ↓
Frontend interceptor detects 401, calls /api/auth/refresh
   ↓
Backend validates session:
   - session.expires_at > NOW → VALID ✅
   - session.is_revoked = false → VALID ✅
   ↓
Backend generates new tokens (access + refresh)
   ↓
Backend updates session:
   - new token_hash
   - new refresh_token_hash
   - expires_at = NOW + 30 days (reset)
   ↓
Frontend receives new tokens
   ↓
Frontend retries original request dengan token baru
   ↓
✅ User tetap logged in selama 30 hari (selama aktif)
```

### Scenario 2: Remember Me Login

```
User Login → RememberMe = true
   ↓
Backend creates session:
   - token_hash
   - refresh_token_hash
   - expires_at = NOW + 90 days ✅ (extended)
   ↓
Session valid selama 90 hari
   ↓
✅ Better UX untuk user yang sering menggunakan aplikasi
```

### Scenario 3: Inactive User (>30 hari)

```
User Login → Last activity 31 days ago
   ↓
User kembali ke aplikasi
   ↓
Frontend restoreAuth() validates refresh token
   ↓
Refresh token JWT expired (>30 hari)
   ↓
Frontend: clearAuth() → redirect to /login
   ↓
✅ User diminta login ulang (expected behavior)
```

## Token Configuration

**Backend:** `backend/config/config.go`

```go
JWTExpiry:          15 * time.Minute      // Access Token
RefreshTokenExpiry: 30 * 24 * time.Hour  // Refresh Token (30 hari)
```

**Environment Variables:**

```bash
JWT_EXPIRY=15m              # Access token expiry
REFRESH_TOKEN_EXPIRY=720h   # Refresh token expiry (30 days = 720 hours)
```

## Session Expiry Strategy

| Login Type | Access Token | Refresh Token | Session DB Expiry | Behavior |
|------------|--------------|---------------|-------------------|----------|
| Normal (RememberMe = false) | 15 min | 30 days | 30 days ✅ | Auto-refresh every 15 min, session valid 30 days |
| Remember Me (RememberMe = true) | 15 min | 30 days | 90 days ✅ | Auto-refresh every 15 min, session valid 90 days |

**Key Points:**
- Access Token selalu 15 menit (security)
- Refresh Token selalu 30 hari (generated dengan JWT)
- Session DB expiry minimum 30 hari (untuk support refresh mechanism)
- Remember Me hanya extend session DB expiry (opsional)

## Testing Checklist

### Manual Testing

- [ ] **Test 1:** Login tanpa Remember Me, tunggu 20 menit, buat API call
  - **Expected:** Token auto-refresh, user tetap logged in
  
- [ ] **Test 2:** Login tanpa Remember Me, close browser, buka lagi setelah 1 jam
  - **Expected:** User tetap logged in, token refresh saat API call pertama
  
- [ ] **Test 3:** Login dengan Remember Me, tunggu beberapa hari
  - **Expected:** User tetap logged in selama <90 hari
  
- [ ] **Test 4:** Login, tidak aktif selama 31 hari, refresh page
  - **Expected:** User diminta login ulang (refresh token expired)
  
- [ ] **Test 5:** Multiple tabs, access token expired
  - **Expected:** Hanya 1 refresh token call, other requests queued

### Database Verification

```sql
-- Check session expiry
SELECT 
    user_id,
    ip_address,
    created_at,
    expires_at,
    TIMESTAMPDIFF(DAY, NOW(), expires_at) as days_until_expiry,
    is_revoked
FROM user_sessions
WHERE user_id = <USER_ID>
ORDER BY created_at DESC
LIMIT 1;
```

**Expected Result:**
- `days_until_expiry` ≈ 30 untuk normal login
- `days_until_expiry` ≈ 90 untuk remember me login

### API Testing

```bash
# 1. Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "123456",
    "password": "password",
    "remember_me": false
  }'

# Save access_token dan refresh_token

# 2. Wait 16 minutes atau manual set expired token

# 3. Try refresh token
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<REFRESH_TOKEN>"
  }'

# Expected: 200 OK dengan new tokens
```

## Benefits

### 1. User Experience
- ✅ User tidak logout setiap 15 menit
- ✅ Seamless token refresh di background
- ✅ Session persist selama 30 hari (normal) atau 90 hari (remember me)
- ✅ No interruption saat bekerja

### 2. Security
- ✅ Access token tetap 15 menit (short-lived)
- ✅ Refresh token rotation (new token setiap refresh)
- ✅ Session tracking di database
- ✅ Token revocation support

### 3. Development
- ✅ Clear separation: Access Token vs Refresh Token vs Session
- ✅ Proper token refresh mechanism
- ✅ Better error messages
- ✅ Easier debugging

## Breaking Changes

**None.** Ini adalah bug fix yang tidak memerlukan frontend changes.

## Migration Notes

**Existing Sessions:** Sessions yang sudah ada dengan `expires_at` lama (15 menit) akan:
- Tetap invalid sampai expired
- User akan logout dan login ulang
- Session baru akan menggunakan logic yang benar (30 hari)

**Recommendation:** Clear existing sessions untuk immediate effect:

```sql
-- Optional: Clear existing sessions untuk immediate fix
-- Users akan diminta login ulang
UPDATE user_sessions SET is_revoked = true WHERE expires_at < NOW() + INTERVAL 30 DAY;
```

## Related Files

- `backend/services/auth_service.go` - Token generation & session management
- `backend/config/config.go` - Token expiry configuration
- `frontend/src/composables/useApi.js` - API interceptor dengan refresh mechanism
- `frontend/src/stores/auth.js` - Auth state management
- `docs/fixes/token-refresh-issue-fix.md` - Related frontend fix

## References

- [JWT Best Practices](https://datatracker.ietf.org/doc/html/rfc8725)
- [OAuth 2.0 Refresh Token](https://www.rfc-editor.org/rfc/rfc6749#section-1.5)
- Sprint 01 - Authentication Implementation

---

**Fixed on:** December 30, 2025  
**Developer:** Zulfikar Hidayatullah  
**Impact:** Critical (UX-breaking bug)  
**Status:** ✅ Fixed & Documented
