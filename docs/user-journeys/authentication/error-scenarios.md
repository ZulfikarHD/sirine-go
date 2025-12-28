# ⚠️ Error Scenarios

## Journey 3: Error - Invalid Credentials

**User Type**: Any  
**Goal**: Attempt login dengan wrong credentials  
**Frequency**: Occasional (user error)  
**Duration**: ~15 seconds

### Flow Diagram

```
📍 START: User at login page
    │
    ├─▶ USER INPUT
    │   ├─ Types NIP: "12345"
    │   └─ Types Password: "wrongpassword"
    │
    ├─▶ SUBMIT FORM
    │   └─ POST /api/auth/login
    │
    ├─▶ BACKEND PROCESSING
    │   ├─ Find user by NIP → User found
    │   ├─ Verify password → ❌ FAILED
    │   ├─ Increment failed_login_attempts (1 → 2)
    │   ├─ Save user record
    │   └─ Return 401 Unauthorized:
    │      {
    │        "success": false,
    │        "message": "NIP atau password salah"
    │      }
    │
    ├─▶ FRONTEND ERROR HANDLING
    │   ├─ Axios catches 401 error
    │   ├─ Error message displayed:
    │   │  └─ Red banner: "NIP atau password salah"
    │   ├─ Card shake animation (400ms)
    │   │  └─ translateX: 0 → -10px → 10px → -10px → 0
    │   ├─ Error haptic feedback
    │   │  └─ Vibration: [100ms, 50ms, 100ms]
    │   └─ Form data retained (tidak clear)
    │
    └─▶ ⚠️ END: User sees error, can retry
        └─ Failed attempts: 2/5
```

### User Experience Notes

- ✅ Error message clear dan actionable
- ✅ Form data retained untuk easy correction
- ✅ Visual feedback (shake) grabs attention
- ✅ Haptic feedback untuk mobile users
- ✅ No sensitive info leaked (generic message)

---

## Journey 4: Error - Account Locked

**User Type**: Any  
**Goal**: Login after 5 failed attempts  
**Frequency**: Rare (security incident)  
**Duration**: ~10 seconds + 15 min wait

### Flow Diagram

```
📍 START: User at login (after 4 failed attempts)
    │
    ├─▶ USER INPUT (5th attempt)
    │   ├─ Types NIP: "12345"
    │   └─ Types Password: "stillwrong"
    │
    ├─▶ SUBMIT FORM
    │   └─ POST /api/auth/login
    │
    ├─▶ BACKEND PROCESSING
    │   ├─ Find user by NIP → User found
    │   ├─ Verify password → ❌ FAILED
    │   ├─ Increment failed_login_attempts (4 → 5)
    │   ├─ Check: attempts >= MAX_LOGIN_ATTEMPTS (5)
    │   │  └─ ✅ TRUE → Trigger lockout
    │   ├─ Calculate locked_until:
    │   │  └─ Now + 15 minutes = 14:45:00
    │   ├─ Update user record:
    │   │  ├─ failed_login_attempts = 5
    │   │  └─ locked_until = 2025-12-27 14:45:00
    │   └─ Return 401 Unauthorized:
    │      {
    │        "success": false,
    │        "message": "Terlalu banyak percobaan login gagal, akun Anda dikunci selama 15 menit"
    │      }
    │
    ├─▶ FRONTEND ERROR HANDLING
    │   ├─ Error banner displayed (red):
    │   │  └─ "Akun terkunci hingga 14:45:00"
    │   ├─ Strong shake animation
    │   └─ Long error haptic (500ms)
    │
    ├─▶ SUBSEQUENT ATTEMPTS (with correct password)
    │   ├─ User tries: NIP "12345", Password "CorrectPass123!"
    │   ├─ Backend checks: user.IsLocked() → TRUE
    │   │  └─ locked_until (14:45) > now (14:32)
    │   └─ Return 401:
    │      "Akun Anda terkunci hingga 14:45:00"
    │
    └─▶ ⛔ END: User must wait 15 minutes
        ├─ Option 1: Wait until 14:45
        │  └─ After 14:45, can login normally
        ├─ Option 2: Contact admin
        │  └─ Admin manually unlocks account
        └─ Security: Prevents brute force attacks
```

### Security Flow

```
Attempt 1: ❌ → failed_login_attempts = 1
Attempt 2: ❌ → failed_login_attempts = 2
Attempt 3: ❌ → failed_login_attempts = 3
Attempt 4: ❌ → failed_login_attempts = 4
Attempt 5: ❌ → failed_login_attempts = 5 → 🔒 LOCKED
           └─ locked_until = now + 15 min
```
