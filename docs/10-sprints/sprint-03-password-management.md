# 🔑 Sprint 3: Password Management & Security

**Version:** 1.3.0  
**Date:** 28 Desember 2025  
**Duration:** 1 week  
**Status:** ✅ Completed

## 📋 Sprint Goals

Implementasi comprehensive password management system dengan self-service password change, forgot/reset password flow, dan security enhancements.

---

## ✨ Features Implemented

### 1. Change Password (Self-Service)

#### Features
- **Current Password Verification** untuk security
- **New Password Validation** dengan strength requirements
- **Password Confirmation** untuk avoid typos
- **Real-time Strength Indicator** dengan visual feedback
- **Requirements Checklist** (✅/❌ untuk setiap requirement)
- **Session Revocation** setelah password change (force re-login)
- **Success Feedback** dengan auto-redirect

#### Password Strength Levels
- **Weak** (0-1): Red color
- **Fair** (2): Orange color
- **Good** (3): Yellow color
- **Strong** (4): Green color

### 2. Forgot Password Flow

#### Request Reset
- **Input Field:** NIP atau Email
- **Email Validation** (jika email provided)
- **Token Generation** (32-byte secure random)
- **Token Hashing** dengan SHA256
- **Email Sending** dengan reset link
- **1-Hour Expiry** untuk security
- **Generic Response** (prevent email enumeration)

#### Reset Password
- **Token Validation** (exist, not used, not expired)
- **New Password Input** dengan strength validation
- **Password Confirmation** field
- **Single-Use Tokens** (marked as used after reset)
- **Session Revocation** setelah reset
- **Success Redirect** ke login page

### 3. Force Change Password

#### First-Time Login
- **Fullscreen Blocking Modal** yang tidak bisa ditutup
- **Mandatory Password Change** sebelum access app
- **New Password Input** dengan validation
- **Password Confirmation** field
- **Requirements Checklist** display
- **Session Refresh** setelah change

#### Admin Force Reset
- **Admin can force reset** user password
- **User flagged** dengan `require_password_change = true`
- **User redirected** to force change modal on next login
- **Audit Log** created untuk transparency

### 4. Session Expired Modal

#### Features
- **Auto-detect** token expiration
- **Fullscreen Modal** dengan glass effect
- **Clear Message** dalam Bahasa Indonesia
- **Login Button** untuk re-authenticate
- **No Close Button** (force re-login)
- **Preserve Current Page** untuk redirect after login

---

## 🎨 New Components

### 1. PasswordStrength.vue
```vue
Features:
- Visual progress bar
- Color-coded strength (red → green)
- Requirements checklist dengan icons
- Real-time validation
- Smooth animations
```

### 2. SessionExpired.vue
```vue
Features:
- Fullscreen blocking modal
- Cannot be dismissed
- Clear error message
- Login redirect button
- Glass effect background
```

### 3. ChangePassword.vue (Page)
```vue
Features:
- Current password field
- New password field
- Confirm password field
- Password strength indicator
- Submit button dengan loading state
- Success/error handling
```

### 4. ForgotPassword.vue (Page)
```vue
Features:
- NIP or Email input
- Email validation
- Submit button
- Success message dengan instructions
- Link to login page
```

### 5. ResetPassword.vue (Page)
```vue
Features:
- Token extraction dari URL
- New password field
- Confirm password field
- Strength indicator
- Submit button
- Error handling (invalid/expired token)
```

### 6. ForceChangePassword.vue (Modal)
```vue
Features:
- Fullscreen modal (cannot close)
- Explanation message
- New password field
- Confirm password field
- Strength indicator
- Submit button
```

---

## 🔌 API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| PUT | `/api/profile/password` | Change own password | Yes |
| POST | `/api/auth/forgot-password` | Request reset link | No |
| POST | `/api/auth/reset-password` | Reset dengan token | No |
| POST | `/api/users/:id/reset-password` | Admin force reset | Yes (Admin) |

---

## 💾 Backend Services

### PasswordService Methods

```go
type PasswordService struct {
    db *gorm.DB
}

// Password Management
func (s *PasswordService) ChangePassword(userID int, currentPassword, newPassword string) error
func (s *PasswordService) ValidatePasswordPolicy(password string) error
func (s *PasswordService) GetPasswordStrength(password string) int

// Reset Token Management
func (s *PasswordService) GenerateResetToken() (string, string, error) // token, hash
func (s *PasswordService) RequestPasswordReset(nipOrEmail string) error
func (s *PasswordService) ResetPassword(token, newPassword string) error
func (s *PasswordService) SendResetEmail(email, token string) error

// Admin Functions
func (s *PasswordService) AdminForceReset(adminID, targetUserID int) error
```

### Password Policy

```go
type PasswordPolicy struct {
    MinLength       int   // 8
    RequireUpper    bool  // true
    RequireLower    bool  // true
    RequireNumber   bool  // true
    RequireSpecial  bool  // true
}

// Validation
func ValidatePasswordPolicy(password string) error {
    // Min 8 characters
    // At least 1 uppercase letter
    // At least 1 lowercase letter
    // At least 1 digit
    // At least 1 special character (!@#$%^&*()_+-=[]{}|;:,.<>?)
}
```

### Password Strength Calculator

```go
func GetPasswordStrength(password string) int {
    score := 0
    
    if len(password) >= 8 { score++ }
    if hasUppercase(password) && hasLowercase(password) { score++ }
    if hasNumber(password) { score++ }
    if hasSpecialChar(password) { score++ }
    
    return score // 0-4
}
```

---

## 🗄️ Database Changes

### Password Reset Tokens Table

```sql
CREATE TABLE password_reset_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_token_hash (token_hash),
    INDEX idx_expires_at (expires_at)
);
```

### Users Table Updates

```sql
ALTER TABLE users ADD COLUMN require_password_change BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN password_changed_at TIMESTAMP NULL;
```

---

## 🔐 Security Measures

### Token Security
- **32-byte Random Token** generated via crypto/rand
- **SHA256 Hashing** untuk storage
- **1-Hour Expiry** untuk minimize attack window
- **Single-Use Tokens** (cannot be reused)
- **Auto-Cleanup** expired tokens (background job)

### Password Security
- **Bcrypt Hashing** dengan cost 12
- **Password Policy Enforcement** (frontend & backend)
- **Strength Validation** sebelum accept
- **No Password in Logs** atau error messages

### Email Security
- **Generic Responses** (prevent email enumeration)
  - "If account exists, email sent" (never confirm/deny)
- **Rate Limiting** on reset requests (prevent abuse)
- **Token in URL** (not in email body)

### Session Security
- **Session Revocation** after password change
- **Force Re-Login** after password change
- **All Sessions Invalidated** (not just current)
- **Audit Log** created untuk password changes

---

## 📱 Frontend Implementation

### Router Integration

```javascript
// Force password change check
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (authStore.user?.require_password_change && to.path !== '/force-change-password') {
    next('/force-change-password')
  } else {
    next()
  }
})
```

### Validation Composable

```javascript
// composables/usePasswordValidation.js
export function usePasswordValidation() {
  const requirements = ref({
    minLength: false,
    hasUpper: false,
    hasLower: false,
    hasNumber: false,
    hasSpecial: false
  })
  
  const strength = computed(() => {
    const score = Object.values(requirements.value).filter(Boolean).length
    return {
      score,
      label: ['Weak', 'Fair', 'Good', 'Strong'][score - 1] || 'Weak',
      color: ['red', 'orange', 'yellow', 'green'][score - 1] || 'red'
    }
  })
  
  const validate = (password) => {
    requirements.value.minLength = password.length >= 8
    requirements.value.hasUpper = /[A-Z]/.test(password)
    requirements.value.hasLower = /[a-z]/.test(password)
    requirements.value.hasNumber = /\d/.test(password)
    requirements.value.hasSpecial = /[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]/.test(password)
  }
  
  return { requirements, strength, validate }
}
```

---

## 🧪 Testing

### Test Scenarios

✅ **Change Password**
- Validate current password correctly
- Reject wrong current password
- Validate new password strength
- Confirm password matches
- Session revoked after change
- Auto-redirect to login

✅ **Forgot Password**
- Accept NIP or Email
- Send email dengan reset link
- Generic response (no enumeration)
- Token expires after 1 hour
- Rate limiting works

✅ **Reset Password**
- Token validation works
- Invalid token rejected
- Expired token rejected
- Used token rejected
- Password strength validated
- Success redirect to login

✅ **Force Change Password**
- Modal shows on first login
- Modal blocks all interactions
- Cannot be dismissed
- Password validation works
- Success removes flag
- User can access app after change

✅ **Session Expired**
- Modal shows on token expiry
- Cannot be dismissed
- Login redirect works
- Current page preserved

---

## 🔄 User Flows

### Change Password Flow
```
1. User clicks "Ganti Password" di Navbar
2. Navigate to /change-password
3. Enter current password
4. Enter new password (strength indicator shows)
5. Confirm new password
6. Submit
7. Session revoked
8. Redirect to login
9. Login dengan new password
```

### Forgot Password Flow
```
1. User clicks "Lupa Password?" di Login
2. Navigate to /forgot-password
3. Enter NIP atau Email
4. Submit
5. Generic success message
6. Email sent dengan reset link
7. User clicks link in email
8. Navigate to /reset-password?token=xxx
9. Enter new password
10. Confirm password
11. Submit
12. Redirect to login
13. Login dengan new password
```

### Force Change Password Flow
```
1. Admin force reset user password
2. User login dengan old password
3. Force change modal shows (fullscreen)
4. User enter new password
5. Confirm password
6. Submit
7. Flag removed
8. Redirect to dashboard
```

---

## 📊 Sprint Metrics

### Development Stats
- **API Endpoints:** 4 new endpoints
- **Components:** 6 new Vue components
- **Services:** Enhanced PasswordService
- **Database Tables:** 1 new table
- **Test Scenarios:** 20+ scenarios

### Security Improvements
- **Password Policy:** Enforced
- **Token Security:** Enhanced
- **Session Management:** Improved
- **Email Security:** Hardened

---

## 🔄 Lessons Learned

### What Went Well ✅
- Password strength indicator intuitive dan helpful
- Token-based reset flow secure dan user-friendly
- Force change modal effective untuk mandatory changes
- Session revocation prevents unauthorized access

### Challenges 🎯
- Email sending membutuhkan SMTP configuration
- Token expiry edge cases butuh careful testing
- Force change modal UX butuh multiple iterations

### Improvements for Next Sprint 💡
- Add email verification
- Implement 2FA (Two-Factor Authentication)
- Add password history (prevent reuse)
- Add account recovery options

---

## 📚 Documentation

### Files Updated
- API documentation (4 new endpoints)
- Security guide updated
- Password policy documented
- User guide updated

---

## 🎯 Next Steps (Sprint 4)

1. **Notification System**
   - In-app notifications
   - Real-time updates
   - Notification center
   - Mark as read/delete

2. **Activity Logs & Audit**
   - Enhanced activity logging
   - Activity log viewer (Admin)
   - Before/after comparison
   - Filtering & search

3. **Profile Enhancements**
   - Profile photo upload
   - Additional fields
   - Profile stats

---

## 📞 Support

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733  
**Sprint Lead:** Zulfikar Hidayatullah

---

## 🔗 Related Documentation

- [Security Guide](../05-guides/security.md)
- [Authentication Guide](../05-guides/authentication/README.md)
- [API Documentation](../03-development/api-documentation.md)
- [Testing Guide](../06-testing/README.md)

---

**Sprint Status:** ✅ Completed  
**Previous Sprint:** [Sprint 2: User Management](./sprint-02-user-management.md)  
**Next Sprint:** [Sprint 4: Notifications & Audit](./sprint-04-notifications-audit.md)  
**Last Updated:** 29 Desember 2025
