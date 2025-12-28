# Sprint 3: Password Management - Quick Reference

Referensi cepat untuk Sprint 3 implementation.

---

## 🎯 Sprint 3 Overview

**Goal**: Password management system dengan change, reset, dan force change password  
**Status**: ✅ **COMPLETE**  
**Files Changed**: 15 files (7 backend, 8 frontend)

---

## 📁 Files Created/Modified

### Backend (7 files)

```
backend/
├── services/
│   └── password_service.go          ✨ Enhanced (6 new methods)
├── handlers/
│   ├── auth_handler.go               ✨ Enhanced (2 new endpoints)
│   └── password_handler.go           ✅ NEW
└── routes/
    └── routes.go                     ✨ Enhanced (4 new routes)
```

### Frontend (8 files)

```
frontend/src/
├── views/
│   ├── auth/
│   │   ├── Login.vue                 ✨ Enhanced (forgot password link)
│   │   ├── ForgotPassword.vue        ✅ NEW
│   │   ├── ResetPassword.vue         ✅ NEW
│   │   └── ForceChangePassword.vue   ✅ NEW
│   └── profile/
│       └── ChangePassword.vue        ✅ NEW
├── components/
│   ├── auth/
│   │   ├── PasswordStrength.vue      ✅ NEW
│   │   └── SessionExpired.vue        ✅ NEW
│   └── layout/
│       └── Navbar.vue                ✨ Enhanced (change password menu)
├── composables/
│   └── useAuth.js                    ✨ Enhanced (3 new methods)
└── router/
    └── index.js                      ✨ Enhanced (4 new routes + guard)
```

---

## 🔌 API Endpoints

### Public Endpoints

```http
POST /api/auth/forgot-password
Content-Type: application/json

{
  "nip_or_email": "99999"
}

Response: 200 OK
{
  "success": true,
  "message": "Jika NIP/Email terdaftar, link reset password telah dikirim..."
}
```

```http
POST /api/auth/reset-password
Content-Type: application/json

{
  "token": "64-char-hex-token",
  "new_password": "NewPassword@123"
}

Response: 200 OK
{
  "success": true,
  "message": "Password berhasil direset. Silakan login dengan password baru Anda."
}
```

### Protected Endpoints

```http
PUT /api/profile/password
Authorization: Bearer <token>
Content-Type: application/json

{
  "current_password": "OldPassword@123",
  "new_password": "NewPassword@123"
}

Response: 200 OK
{
  "success": true,
  "message": "Password berhasil diubah. Silakan login kembali dengan password baru Anda."
}
```

### Admin Endpoints

```http
POST /api/users/:id/reset-password
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "new_password": "TempPassword@123"
}

Response: 200 OK
{
  "success": true,
  "message": "Password berhasil direset",
  "password": "TempPassword@123",
  "note": "User harus mengubah password saat login pertama kali"
}
```

---

## 🛣️ Frontend Routes

```javascript
// Public Routes
/login                      // Login page (+ forgot password link)
/forgot-password            // Request reset link
/reset-password?token=xxx   // Reset password dengan token

// Protected Routes
/force-change-password      // First-time login (blocking modal)
/profile/change-password    // Self-service change password
```

---

## 🔐 Password Policy

```javascript
// Requirements
- Minimal 8 karakter
- Minimal 1 huruf besar (A-Z)
- Minimal 1 angka (0-9)
- Minimal 1 karakter spesial (!@#$%^&*, dll)

// Validation
✅ "Abc123!@"      // Valid
❌ "abc123!@"      // No uppercase
❌ "Abcdefgh"      // No number
❌ "Abc12345"      // No special char
❌ "Abc123!"       // < 8 chars (actually valid, 8 chars)
```

---

## 🎨 Components Usage

### PasswordStrength Component

```vue
<template>
  <PasswordStrength
    :password="form.password"
    :show-requirements="true"
  />
</template>

<script setup>
import PasswordStrength from '@/components/auth/PasswordStrength.vue'
import { ref } from 'vue'

const form = ref({ password: '' })
</script>
```

**Features:**
- Real-time strength calculation
- Color-coded progress bar (red → yellow → green → emerald)
- Requirements checklist dengan checkmarks
- Exposed `allRequirementsMet` untuk validation

### SessionExpired Component

```vue
<template>
  <SessionExpired
    :show="isSessionExpired"
    :allow-dismiss="false"
    @login="handleLogin"
  />
</template>

<script setup>
import SessionExpired from '@/components/auth/SessionExpired.vue'
import { ref } from 'vue'

const isSessionExpired = ref(false)

const handleLogin = () => {
  // Handle login redirect
}
</script>
```

---

## 🔧 Composable Methods

### useAuth Composable

```javascript
import { useAuth } from '@/composables/useAuth'

const {
  forgotPassword,    // Request reset link
  resetPassword,     // Reset dengan token
  changePassword,    // Change own password
  isLoading,         // Loading state
  error,             // Error message
} = useAuth()

// Forgot Password
await forgotPassword('99999')
// or
await forgotPassword('admin@sirine.local')

// Reset Password
await resetPassword('token-from-email', 'NewPassword@123')

// Change Password
await changePassword('OldPassword@123', 'NewPassword@123')
// Note: Auto-logout setelah 2 detik
```

---

## 🗄️ Database Schema

### password_reset_tokens Table

```sql
CREATE TABLE password_reset_tokens (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  token_hash VARCHAR(255) NOT NULL UNIQUE,
  expires_at TIMESTAMP NOT NULL,
  used_at TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  INDEX idx_user_id (user_id),
  INDEX idx_expires_at (expires_at),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### users Table (relevant fields)

```sql
ALTER TABLE users ADD COLUMN must_change_password BOOLEAN DEFAULT FALSE;
```

---

## ⚙️ Configuration

### Environment Variables (.env)

```env
# Email Configuration
EMAIL_SMTP_HOST=smtp.mailtrap.io
EMAIL_SMTP_PORT=2525
EMAIL_USERNAME=your-username
EMAIL_PASSWORD=your-password
EMAIL_FROM_ADDRESS=noreply@sirine.local

# JWT Configuration
JWT_SECRET=your-secure-secret-key
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=720h

# Security
BCRYPT_COST=12
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=15m

# Frontend
FRONTEND_URL=http://localhost:5173
```

---

## 🧪 Testing Commands

### Manual Testing Flow

```bash
# 1. Start backend
cd backend
make run

# 2. Start frontend
cd frontend
yarn dev

# 3. Setup email service (MailHog)
go install github.com/mailhog/MailHog@latest
MailHog
# Access: http://localhost:8025

# 4. Test forgot password
# - Go to http://localhost:5173/login
# - Click "Lupa Password?"
# - Enter NIP/Email
# - Check MailHog for email
# - Click reset link
# - Enter new password

# 5. Test change password
# - Login as user
# - Click avatar → "Ganti Password"
# - Fill form
# - Verify auto-logout

# 6. Test force change password
# - Create user dengan must_change_password = true
# - Login → verify redirect ke force change
# - Change password
# - Verify access granted
```

---

## 🐛 Common Issues & Solutions

### Issue: Email tidak terkirim

**Solution:**
```bash
# Check SMTP configuration
echo $EMAIL_SMTP_HOST
echo $EMAIL_SMTP_PORT

# Test SMTP connection
telnet smtp.mailtrap.io 2525

# Check backend logs
tail -f backend/logs/app.log
```

### Issue: Token invalid

**Solution:**
```sql
-- Check token in database
SELECT * FROM password_reset_tokens 
WHERE user_id = <user_id>
ORDER BY created_at DESC;

-- Check if expired
SELECT *, 
  CASE WHEN expires_at < NOW() THEN 'EXPIRED' ELSE 'VALID' END as status
FROM password_reset_tokens;
```

### Issue: Session tidak revoked

**Solution:**
```sql
-- Check user sessions
SELECT * FROM user_sessions 
WHERE user_id = <user_id>
ORDER BY created_at DESC;

-- Manual revoke (for testing)
UPDATE user_sessions 
SET is_revoked = true 
WHERE user_id = <user_id>;
```

---

## 📊 Performance Benchmarks

```
Operation                  | Expected Time
---------------------------|---------------
Password hashing (bcrypt)  | ~200ms
Token generation           | <10ms
Email sending (async)      | Non-blocking
Password strength calc     | <5ms
Page load (lazy)           | <500ms
Animation frame rate       | 60fps
```

---

## 🎯 User Flows

### Flow 1: Forgot Password

```
User → Login Page
  ↓
Click "Lupa Password?"
  ↓
Enter NIP/Email → Submit
  ↓
Success Message
  ↓
Check Email
  ↓
Click Reset Link
  ↓
Reset Password Page
  ↓
Enter New Password → Submit
  ↓
Success → Redirect to Login
  ↓
Login dengan Password Baru
```

### Flow 2: Change Password

```
User → Dashboard
  ↓
Click Avatar → "Ganti Password"
  ↓
Change Password Page
  ↓
Enter Current + New Password → Submit
  ↓
Success Message
  ↓
Auto-Logout (2s)
  ↓
Redirect to Login
  ↓
Login dengan Password Baru
```

### Flow 3: Force Change Password

```
New User → Login
  ↓
Check must_change_password = true
  ↓
Redirect to Force Change Password Modal
  ↓
Blocking Modal (cannot dismiss)
  ↓
Enter Temporary + New Password → Submit
  ↓
Set must_change_password = false
  ↓
Redirect to Dashboard
  ↓
Normal Access Granted
```

---

## 🔗 Related Documentation

- [Sprint 3 Summary](./SPRINT3_SUMMARY.md) - Complete implementation details
- [Sprint 3 Testing Guide](./SPRINT3_TESTING_GUIDE.md) - Comprehensive testing scenarios
- [Sprint Plan](./docs/sprint_plan_authentication.md) - Original sprint planning

---

## 📞 Support

**Developer**: Zulfikar Hidayatullah  
**Phone**: +62 857-1583-8733  
**Timezone**: Asia/Jakarta (WIB)

---

**Last Updated**: December 28, 2025  
**Sprint Status**: ✅ **COMPLETE**
