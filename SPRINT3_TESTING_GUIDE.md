# Sprint 3: Password Management - Testing Guide

Quick reference untuk testing Sprint 3 features.

---

## 🚀 Quick Start

### 1. Setup Email Service (Development)

**Option A: Mailtrap (Recommended)**
```env
EMAIL_SMTP_HOST=smtp.mailtrap.io
EMAIL_SMTP_PORT=2525
EMAIL_USERNAME=your-mailtrap-username
EMAIL_PASSWORD=your-mailtrap-password
EMAIL_FROM_ADDRESS=noreply@sirine.local
```

**Option B: MailHog (Local)**
```bash
# Install MailHog
go install github.com/mailhog/MailHog@latest

# Run MailHog
MailHog

# Configure .env
EMAIL_SMTP_HOST=localhost
EMAIL_SMTP_PORT=1025
EMAIL_USERNAME=
EMAIL_PASSWORD=
EMAIL_FROM_ADDRESS=noreply@sirine.local

# Access web UI: http://localhost:8025
```

### 2. Start Services

```bash
# Terminal 1: Backend
cd backend
make run

# Terminal 2: Frontend
cd frontend
yarn dev
```

### 3. Access Application
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- MailHog UI (if using): http://localhost:8025

---

## 🧪 Test Scenarios

### Scenario 1: Forgot Password Flow

**Steps:**
1. Go to http://localhost:5173/login
2. Click "Lupa Password?"
3. Enter NIP: `99999` atau Email: `admin@sirine.local`
4. Click "Kirim Link Reset Password"
5. Check email inbox (Mailtrap/MailHog)
6. Copy reset link dari email
7. Paste link ke browser
8. Enter new password: `NewAdmin@123`
9. Confirm password: `NewAdmin@123`
10. Click "Reset Password"
11. Verify redirect ke login
12. Login dengan password baru

**Expected Results:**
- ✅ Success message: "Email terkirim"
- ✅ Email received dengan reset link
- ✅ Reset page loaded dengan token
- ✅ Password strength indicator works
- ✅ Success message: "Password berhasil direset"
- ✅ Auto-redirect ke login (3s)
- ✅ Login success dengan new password
- ✅ Old password tidak bisa digunakan

---

### Scenario 2: Change Password (Self-Service)

**Steps:**
1. Login sebagai user (NIP: 99999)
2. Click user avatar (kanan atas)
3. Click "Ganti Password"
4. Enter current password: `Admin@123`
5. Enter new password: `NewAdmin@456`
6. Confirm password: `NewAdmin@456`
7. Click "Simpan Password Baru"
8. Wait for auto-logout (2s)
9. Login dengan password baru

**Expected Results:**
- ✅ Password strength indicator shows real-time
- ✅ Requirements checklist updates
- ✅ Success message: "Password berhasil diubah"
- ✅ Auto-logout dan redirect ke login
- ✅ All sessions revoked
- ✅ Login success dengan new password

---

### Scenario 3: Force Change Password (First-Time Login)

**Preparation:**
```sql
-- Create test user dengan must_change_password = true
INSERT INTO users (nip, full_name, email, password_hash, role, department, must_change_password, status)
VALUES ('12345', 'Test User', 'test@sirine.local', '$2a$12$...', 'STAFF_KHAZWAL', 'KHAZWAL', true, 'ACTIVE');
```

**Steps:**
1. Login dengan NIP: `12345`, Password: `Test@123`
2. Verify redirect ke force change password modal
3. Try access dashboard → blocked
4. Enter temporary password: `Test@123`
5. Enter new password: `NewTest@456`
6. Confirm password: `NewTest@456`
7. Click "Ubah Password & Lanjutkan"
8. Verify redirect ke dashboard

**Expected Results:**
- ✅ Fullscreen modal muncul (tidak bisa di-dismiss)
- ✅ Info alert: "Ini adalah login pertama Anda..."
- ✅ Cannot access dashboard tanpa change password
- ✅ Password strength indicator works
- ✅ Success change password
- ✅ `must_change_password` set to false
- ✅ Redirect ke dashboard
- ✅ Can access dashboard normally

---

### Scenario 4: Admin Force Reset User Password

**Steps:**
1. Login sebagai Admin (NIP: 99999)
2. Go to "Manajemen User"
3. Find user "Test User"
4. Click "..." menu → "Reset Password"
5. Enter new password: `Reset@123`
6. Click "Reset Password"
7. Copy generated password
8. Logout
9. Login sebagai Test User dengan password baru

**Expected Results:**
- ✅ Admin can reset user password
- ✅ New password displayed once
- ✅ Copy button works
- ✅ User `must_change_password` set to true
- ✅ User forced to change password on next login

---

## 🔍 Validation Testing

### Password Policy Validation

Test dengan passwords berikut:

| Password | Expected Result | Reason |
|----------|----------------|---------|
| `abc` | ❌ Error | < 8 characters |
| `abcdefgh` | ❌ Error | No uppercase |
| `Abcdefgh` | ❌ Error | No number |
| `Abcdefg1` | ❌ Error | No special char |
| `Abc123!@` | ✅ Valid | Meets all requirements |

### Password Strength Indicator

| Password | Strength | Color | Progress |
|----------|----------|-------|----------|
| `abc` | Lemah | Red | 25% |
| `Abc123` | Sedang | Yellow | 50% |
| `Abc123!` | Kuat | Green | 75% |
| `Abc123!@#$%^` | Sangat Kuat | Emerald | 100% |

---

## 🐛 Error Testing

### Token Expiry Test

1. Request forgot password
2. Get reset link
3. Wait 1 hour (atau modify token expiry untuk testing)
4. Try reset password
5. Expected: "Token reset sudah kadaluarsa"

### Token Reuse Test

1. Request forgot password
2. Reset password successfully
3. Try use same link again
4. Expected: "Token reset sudah digunakan"

### Wrong Current Password

1. Go to change password
2. Enter wrong current password
3. Expected: "Password saat ini tidak valid"

### Password Mismatch

1. Enter new password: `Abc123!@`
2. Enter confirm password: `Abc123!@#`
3. Expected: "Password tidak cocok"

---

## 📧 Email Testing

### Email Content Verification

Check email contains:
- ✅ User's full name
- ✅ Reset link dengan token
- ✅ Expiry information (1 hour)
- ✅ Security notice
- ✅ Ignore message if not requested

### Email Link Format

```
http://localhost:5173/reset-password?token=<64-char-hex-token>
```

---

## 🔒 Security Testing

### Session Revocation Test

1. Login dari 2 browsers (Chrome & Firefox)
2. Change password dari Chrome
3. Try access dashboard dari Firefox
4. Expected: Session expired modal → redirect to login

### Email Enumeration Prevention

1. Request forgot password dengan email tidak terdaftar
2. Expected: Success message (same as valid email)
3. No error revealing email tidak exist

### Rate Limiting (Future)

Currently tidak ada rate limiting untuk forgot password.
Recommendation: Add rate limiting (max 3 requests per 15 minutes).

---

## 📱 Mobile Testing

### Responsive Design

Test pada screen sizes:
- 📱 Mobile: 375px (iPhone SE)
- 📱 Mobile: 414px (iPhone 12 Pro)
- 📱 Tablet: 768px (iPad)
- 💻 Desktop: 1024px+

### Touch Interactions

- ✅ Button press feedback (scale 0.97)
- ✅ Input focus states
- ✅ Password toggle button
- ✅ Dropdown menus

### Haptic Feedback

Test pada mobile device:
- ✅ Success: Single pulse (200ms)
- ✅ Error: Double pulse (100ms, 50ms, 100ms)

---

## 🎨 UI/UX Testing

### Animations

- ✅ Page entrance: Fade + scale
- ✅ Icon pop-in: Scale + rotate
- ✅ Success/error messages: Fade + slide
- ✅ Progress bar: Smooth width transition
- ✅ Modal backdrop: Fade in/out

### Loading States

- ✅ Button spinner during submit
- ✅ Disabled state dengan opacity
- ✅ Loading text: "Memproses..."

### Error Messages

All error messages dalam Bahasa Indonesia:
- ✅ "Password minimal 8 karakter"
- ✅ "Password harus mengandung minimal 1 huruf besar"
- ✅ "Password tidak cocok"
- ✅ "Token reset sudah kadaluarsa"

---

## 🚨 Edge Cases

### Empty Form Submission

1. Submit form tanpa isi fields
2. Expected: Validation errors untuk required fields

### Special Characters in Password

1. Test dengan password: `P@ssw0rd!@#$%^&*()`
2. Expected: Valid dan accepted

### Very Long Password

1. Test dengan password 50+ characters
2. Expected: Accepted (no max length)

### Concurrent Password Changes

1. Open 2 tabs
2. Start change password di both tabs
3. Submit dari tab 1 → success
4. Submit dari tab 2 → error (current password wrong)

---

## 📊 Performance Testing

### Backend Response Times

Expected response times:
- Login: < 300ms
- Forgot password: < 200ms
- Reset password: < 300ms
- Change password: < 300ms

### Frontend Rendering

- Initial page load: < 500ms
- Password strength calculation: < 5ms (real-time)
- Animation frame rate: 60fps

---

## ✅ Acceptance Criteria Checklist

### Forgot Password
- [ ] User dapat request reset link
- [ ] Email terkirim dengan link valid
- [ ] Link expired setelah 1 hour
- [ ] Token single-use (tidak bisa reuse)
- [ ] Success message user-friendly

### Reset Password
- [ ] User dapat reset dengan token valid
- [ ] Password policy enforced
- [ ] Password strength indicator works
- [ ] Success → redirect ke login
- [ ] Can login dengan new password

### Change Password
- [ ] Current password validated
- [ ] New password ≠ current password
- [ ] Password policy enforced
- [ ] All sessions revoked
- [ ] Auto-logout setelah change

### Force Change Password
- [ ] Triggered untuk first-time login
- [ ] Fullscreen blocking modal
- [ ] Cannot access dashboard tanpa change
- [ ] Success → `must_change_password = false`
- [ ] Can access dashboard setelah change

---

## 🔧 Troubleshooting

### Email Tidak Terkirim

**Check:**
1. SMTP credentials correct?
2. SMTP host/port reachable?
3. Check backend logs untuk error
4. Verify email service running (MailHog)

**Solution:**
```bash
# Test SMTP connection
telnet smtp.mailtrap.io 2525

# Check backend logs
tail -f backend/logs/app.log
```

### Token Invalid Error

**Check:**
1. Token format correct (64 hex chars)?
2. Token expired?
3. Token already used?

**Solution:**
```sql
-- Check token in database
SELECT * FROM password_reset_tokens WHERE token_hash = '<hash>';

-- Check expiry
SELECT *, expires_at < NOW() as is_expired FROM password_reset_tokens;
```

### Session Not Revoked

**Check:**
1. Database connection OK?
2. User sessions table updated?

**Solution:**
```sql
-- Check user sessions
SELECT * FROM user_sessions WHERE user_id = <user_id>;

-- Verify is_revoked flag
SELECT is_revoked FROM user_sessions WHERE user_id = <user_id>;
```

---

## 📝 Test Report Template

```markdown
## Sprint 3 Testing Report

**Tester**: [Your Name]
**Date**: [Date]
**Environment**: Development

### Test Results

#### Forgot Password Flow
- [ ] Request reset link: PASS / FAIL
- [ ] Email received: PASS / FAIL
- [ ] Reset password: PASS / FAIL
- [ ] Login dengan new password: PASS / FAIL

#### Change Password Flow
- [ ] Validation works: PASS / FAIL
- [ ] Password changed: PASS / FAIL
- [ ] Sessions revoked: PASS / FAIL
- [ ] Auto-logout: PASS / FAIL

#### Force Change Password
- [ ] Modal triggered: PASS / FAIL
- [ ] Dashboard blocked: PASS / FAIL
- [ ] Password changed: PASS / FAIL
- [ ] Access granted: PASS / FAIL

### Issues Found
1. [Issue description]
2. [Issue description]

### Recommendations
1. [Recommendation]
2. [Recommendation]
```

---

**Happy Testing! 🎉**

Jika ada issues atau questions, contact developer:
- **Name**: Zulfikar Hidayatullah
- **Phone**: +62 857-1583-8733
