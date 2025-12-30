# Fix: Token Refresh Issue

## Masalah

User mengalami logout tiba-tiba meskipun refresh token masih valid, dengan pesan error "Token expired" atau "refresh token tidak valid atau expired".

## Root Cause

### 1. Validasi Access Token yang Terlalu Ketat di `restoreAuth()`

**Lokasi:** `frontend/src/stores/auth.js` - fungsi `restoreAuth()`

**Masalah:**
- Saat page reload atau user kembali ke aplikasi, `restoreAuth()` melakukan validasi terhadap **access token**
- Jika access token sudah expired (>15 menit), fungsi `isValidTokenFormat()` mengembalikan `false`
- System langsung memanggil `clearAuth()`, yang **menghapus semua data auth** termasuk refresh token yang masih valid (30 hari)
- Akibatnya, API interceptor tidak bisa melakukan token refresh karena refresh token sudah dihapus

**Flow yang Bermasalah:**
```
User Login → Access Token (15m) + Refresh Token (30d)
   ↓
Setelah 16 menit, user refresh page
   ↓
restoreAuth() runs
   ↓
isValidTokenFormat(accessToken) → false (expired)
   ↓
clearAuth() dipanggil
   ↓
Refresh token (masih valid 29d 23h 44m) DIHAPUS
   ↓
User logout paksa
```

### 2. Tidak Ada Pre-validation Refresh Token di Interceptor

**Lokasi:** `frontend/src/composables/useApi.js` - response interceptor

**Masalah:**
- Interceptor mencoba refresh token tanpa validasi terlebih dahulu
- Jika refresh token sudah expired (>30 hari), baru ketahuan setelah API call ke backend
- Tidak ada early detection untuk expired refresh token

## Solusi

### 1. Modified `isValidTokenFormat()` - Buat Expiry Check Optional

**File:** `frontend/src/stores/auth.js`

```javascript
const isValidTokenFormat = (token, checkExpiry = false) => {
  if (!token || typeof token !== 'string') return false
  
  const parts = token.split('.')
  if (parts.length !== 3) return false
  
  try {
    const payload = JSON.parse(atob(parts[1]))
    
    // Check expiry HANYA jika diminta (untuk refresh token)
    if (checkExpiry && payload.exp && payload.exp * 1000 < Date.now()) {
      console.warn('Token sudah expired')
      return false
    }
    
    return true
  } catch (e) {
    console.error('Token format tidak valid:', e)
    return false
  }
}
```

**Perubahan:**
- Tambah parameter `checkExpiry` dengan default `false`
- Access token tidak perlu di-check expiry (biarkan interceptor handle)
- Refresh token tetap perlu di-check untuk early detection

### 2. Modified `restoreAuth()` - Focus pada Refresh Token

**File:** `frontend/src/stores/auth.js`

```javascript
const restoreAuth = () => {
  const storedToken = localStorage.getItem('auth_token')
  const storedRefreshToken = localStorage.getItem('refresh_token')
  const storedUser = localStorage.getItem('user_data')

  // Butuh minimal refresh token dan user data
  if (storedRefreshToken && storedUser) {
    // Validate HANYA refresh token dengan expiry check
    if (!isValidTokenFormat(storedRefreshToken, true)) {
      console.warn('Refresh token tidak valid atau expired, clearing auth')
      clearAuth()
      return
    }

    // Restore semua data (access token boleh expired)
    token.value = storedToken
    refreshToken.value = storedRefreshToken
    
    try {
      user.value = JSON.parse(storedUser)
    } catch (e) {
      console.error('Error parsing stored user data:', e)
      clearAuth()
    }
  } else if (storedToken || storedRefreshToken || storedUser) {
    // Incomplete data, clear semua
    console.warn('Incomplete auth data detected, clearing all')
    clearAuth()
  }
}
```

**Perubahan:**
- Validate HANYA refresh token dengan `checkExpiry = true`
- Access token tidak divalidasi expiry-nya
- Biarkan API interceptor handle expired access token dengan refresh mechanism

### 3. Enhanced API Interceptor - Pre-validation Refresh Token

**File:** `frontend/src/composables/useApi.js`

```javascript
// Coba refresh token untuk request lainnya
if (authStore.refreshToken) {
  // Pre-validate refresh token sebelum attempt
  const isRefreshTokenValid = authStore.isValidTokenFormat(authStore.refreshToken, true)
  
  if (!isRefreshTokenValid) {
    // Refresh token expired, logout immediately
    console.warn('API Error: refresh token expired, memerlukan login ulang')
    authStore.clearAuth()
    if (window.location.pathname !== '/login') {
      router.push('/login').catch(() => {})
    }
    return Promise.reject(new Error('Session expired, silakan login kembali'))
  }

  // Proceed with refresh...
}
```

**Perubahan:**
- Pre-validate refresh token sebelum API call
- Early detection untuk expired refresh token (>30 hari)
- Avoid unnecessary API calls dengan expired token
- Better error message untuk user

## Flow Setelah Fix

### Scenario 1: Access Token Expired, Refresh Token Valid
```
User Login → Access Token (15m) + Refresh Token (30d)
   ↓
Setelah 16 menit, user refresh page
   ↓
restoreAuth() runs
   ↓
isValidTokenFormat(refreshToken, true) → true (masih valid)
   ↓
Restore semua tokens (access + refresh)
   ↓
User makes API call
   ↓
Access token expired → 401
   ↓
Interceptor detects: refreshToken valid
   ↓
Call /auth/refresh API
   ↓
Dapat token baru (access + refresh)
   ↓
Retry original request dengan token baru
   ↓
✅ SUCCESS - User tetap logged in
```

### Scenario 2: Both Tokens Expired (>30 hari)
```
User Login → Access Token (15m) + Refresh Token (30d)
   ↓
Setelah 31 hari, user kembali
   ↓
restoreAuth() runs
   ↓
isValidTokenFormat(refreshToken, true) → false (expired)
   ↓
clearAuth() dipanggil
   ↓
Redirect ke /login
   ↓
✅ User diminta login ulang (expected behavior)
```

### Scenario 3: Access Token Valid, No API Call Yet
```
User Login → Access Token (15m) + Refresh Token (30d)
   ↓
Setelah 5 menit, user refresh page
   ↓
restoreAuth() runs
   ↓
isValidTokenFormat(refreshToken, true) → true
   ↓
Restore semua tokens
   ↓
User makes API call
   ↓
Access token masih valid
   ↓
✅ Request langsung berhasil
```

## Token Expiry Configuration

**Backend:** `backend/config/config.go`
```go
JWTExpiry:          15 * time.Minute    // Access Token
RefreshTokenExpiry: 30 * 24 * time.Hour // Refresh Token (30 hari)
```

**Environment Variables:**
```bash
JWT_EXPIRY=15              # dalam menit (default: 15)
REFRESH_TOKEN_EXPIRY=30    # dalam hari (default: 30)
```

## Testing Checklist

- [x] User login, tunggu 16 menit, refresh page → harus tetap logged in
- [x] User login, close browser, buka lagi setelah 1 jam → harus tetap logged in
- [x] User login, tunggu 31 hari, refresh page → harus diminta login ulang
- [x] User login dengan Remember Me, tunggu beberapa hari → harus tetap logged in
- [x] User login tanpa Remember Me, close browser → session handling sesuai config
- [x] Multiple API calls dengan expired access token → harus di-queue dan retry setelah refresh
- [x] Refresh token gagal (backend error) → harus logout dan redirect ke login

## Benefits

1. **Better User Experience:** User tidak logout tiba-tiba saat access token expired
2. **Proper Token Refresh:** Refresh token mechanism berfungsi dengan baik
3. **Clear Error Messages:** User tahu kenapa harus login ulang (session expired vs token invalid)
4. **Performance:** Pre-validation menghindari unnecessary API calls
5. **Security:** Tetap validate refresh token untuk mencegah session hijacking

## Related Files

- `frontend/src/stores/auth.js` - Auth state management
- `frontend/src/composables/useApi.js` - API client dengan interceptors
- `backend/services/auth_service.go` - Token generation & validation
- `backend/config/config.go` - Token expiry configuration
