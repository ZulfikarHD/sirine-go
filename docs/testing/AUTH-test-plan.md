# 🧪 AUTH - Authentication System Test Plan

**Feature Code**: AUTH  
**Feature**: Authentication System dengan JWT  
**Sprint**: 1  
**Last Updated**: 27 Desember 2025

---

## 📋 Test Plan Overview

Test Plan merupakan comprehensive QA checklist untuk authentication system yang bertujuan untuk memastikan semua functionality, security, dan edge cases ter-cover dengan baik, yaitu: unit tests, integration tests, manual testing, security testing, dan performance testing.

---

## 🎯 Test Scope

### In Scope
- ✅ Login/Logout functionality
- ✅ JWT token generation dan validation
- ✅ Session management
- ✅ Rate limiting & account lockout
- ✅ Password hashing & validation
- ✅ Role-based access control
- ✅ Token refresh mechanism
- ✅ Frontend state management
- ✅ Router guards

### Out of Scope (Future Sprints)
- ❌ Password reset via email (Sprint 3)
- ❌ Two-factor authentication (Backlog)
- ❌ Social login (Backlog)
- ❌ Biometric authentication (Backlog)

---

## 🔧 Unit Tests

### Backend Unit Tests

#### Password Service Tests
**Location**: `backend/tests/unit/services/password_service_test.go`

| Test Case | Description | Status |
|-----------|-------------|--------|
| `TestHashPassword` | Password hashing dengan bcrypt cost 12 | ✅ Pass |
| `TestHashPasswordEmpty` | Error handling untuk empty password | ✅ Pass |
| `TestVerifyPassword` | Password verification dengan 4 scenarios | ✅ Pass |
| `TestValidatePasswordPolicy` | Policy enforcement (8 char, uppercase, number, special) | ✅ Pass |
| `TestGetPasswordStrength` | Strength calculation (0-4 scale) | ✅ Pass |
| `TestValidateAndHash` | Combined validation dan hashing | ✅ Pass |
| `BenchmarkHashPassword` | Performance benchmark (~200ms) | ✅ Pass |
| `BenchmarkVerifyPassword` | Verification performance | ✅ Pass |

**Run Command**:
```bash
cd backend
go test -v ./tests/unit/services
```

**Coverage**: 95%+

#### User Model Tests
**Location**: `backend/tests/unit/models/user_test.go`

| Test Case | Description | Status |
|-----------|-------------|--------|
| `TestUserIsLocked` | Lock status checking (3 scenarios) | ✅ Pass |
| `TestUserIsActive` | Active status checking (4 scenarios) | ✅ Pass |
| `TestUserHasRole` | Role checking (4 scenarios) | ✅ Pass |
| `TestUserIsAdmin` | Admin role checking (4 scenarios) | ✅ Pass |
| `TestUserToSafeUser` | Safe user conversion (no password leak) | ✅ Pass |
| `TestUserSessionIsValid` | Session validation (4 scenarios) | ✅ Pass |
| `TestPasswordResetTokenIsValid` | Token validation (4 scenarios) | ✅ Pass |

**Run Command**:
```bash
cd backend
go test -v ./tests/unit/models
```

**Coverage**: 90%+

### Frontend Unit Tests

#### Auth Store Tests
**Location**: `frontend/src/tests/unit/stores/auth.spec.js`

| Test Case | Description | Status |
|-----------|-------------|--------|
| Initial State | User null, not authenticated | ✅ Pass |
| `setAuth` | Set auth data + localStorage | ✅ Pass |
| `clearAuth` | Clear auth data + localStorage | ✅ Pass |
| `restoreAuth` | Restore dari localStorage | ✅ Pass |
| `restoreAuth` (invalid JSON) | Error handling | ✅ Pass |
| `hasRole` (single) | Role checking | ✅ Pass |
| `hasRole` (multiple) | Multiple roles checking | ✅ Pass |
| `hasRole` (no user) | Null user handling | ✅ Pass |
| `isAdmin` (ADMIN) | Admin role detection | ✅ Pass |
| `isAdmin` (MANAGER) | Manager role detection | ✅ Pass |
| `isAdmin` (STAFF) | Non-admin rejection | ✅ Pass |
| `userRole` computed | Role getter | ✅ Pass |
| `userDepartment` computed | Department getter | ✅ Pass |
| `requirePasswordChange` computed | Password change flag | ✅ Pass |
| localStorage integration | Persistence testing | ✅ Pass |

**Run Command**:
```bash
cd frontend
yarn test
```

**Coverage**: 95%+

---

## 🔗 Integration Tests

### API Integration Tests

#### Test Case 1: Complete Login Flow
**Objective**: Verify end-to-end login process

**Steps**:
1. POST `/api/auth/login` dengan valid credentials
2. Verify response contains token, refresh_token, user data
3. Extract token dari response
4. GET `/api/auth/me` dengan token di Authorization header
5. Verify user data returned correctly

**Expected Result**:
- ✅ Login returns 200 OK
- ✅ Token dan refresh_token valid
- ✅ User data complete dan accurate
- ✅ `/api/auth/me` returns same user data

**Command**:
```bash
# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"nip":"99999","password":"Admin@123"}' \
  | jq -r '.data.token')

# 2. Get current user
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/auth/me
```

#### Test Case 2: Token Refresh Flow
**Objective**: Verify token refresh mechanism

**Steps**:
1. Login dan get refresh_token
2. Wait for access token to expire (atau set JWT_EXPIRY=1m)
3. POST `/api/auth/refresh` dengan refresh_token
4. Verify new tokens generated
5. Use new token untuk API call

**Expected Result**:
- ✅ Refresh returns new access token
- ✅ New refresh token generated
- ✅ Old token invalid
- ✅ New token works for API calls

#### Test Case 3: Logout Flow
**Objective**: Verify session revocation

**Steps**:
1. Login dan get token
2. POST `/api/auth/logout` dengan token
3. Try use same token untuk `/api/auth/me`

**Expected Result**:
- ✅ Logout returns 200 OK
- ✅ Token revoked di database
- ✅ Subsequent requests dengan old token return 401

#### Test Case 4: Rate Limiting
**Objective**: Verify account lockout mechanism

**Steps**:
1. Login dengan wrong password 5x
2. Check user record di database
3. Try login dengan correct password

**Expected Result**:
- ✅ After 5 failures, account locked
- ✅ `locked_until` timestamp set (15 min)
- ✅ `failed_login_attempts` = 5
- ✅ Correct password rejected until lockout expires

---

## 🧑‍💻 Manual Testing

### Test Case M1: Login Success
**Priority**: P0 (Critical)

**Prerequisites**:
- Backend running di `localhost:8080`
- Frontend running di `localhost:5173`
- Database seeded dengan admin user

**Steps**:
1. Buka `http://localhost:5173`
2. Verify auto-redirect ke `/login`
3. Input NIP: `99999`
4. Input Password: `Admin@123`
5. Check "Remember Me" checkbox
6. Click "Masuk" button

**Expected Results**:
- [ ] Form validation passes
- [ ] Loading spinner appears
- [ ] Redirect ke `/dashboard/admin`
- [ ] Navbar shows "Administrator"
- [ ] User dropdown functional
- [ ] Token saved di localStorage
- [ ] Haptic feedback (vibration) triggered
- [ ] Spring animation on button press

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case M2: Login Failed
**Priority**: P0 (Critical)

**Steps**:
1. Buka login page
2. Input NIP: `12345`
3. Input Password: `wrongpassword`
4. Click "Masuk"

**Expected Results**:
- [ ] Error message: "NIP atau password salah"
- [ ] Card shake animation triggered
- [ ] Error haptic feedback (vibration)
- [ ] Form tidak clear (data tetap ada)
- [ ] No redirect
- [ ] No token saved

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case M3: Rate Limiting
**Priority**: P0 (Critical)

**Steps**:
1. Login dengan wrong password (attempt 1)
2. Login dengan wrong password (attempt 2)
3. Login dengan wrong password (attempt 3)
4. Login dengan wrong password (attempt 4)
5. Login dengan wrong password (attempt 5)
6. Try login dengan correct password

**Expected Results**:
- [ ] First 4 attempts: "NIP atau password salah"
- [ ] 5th attempt: "Akun terkunci selama 15 menit"
- [ ] Correct password rejected
- [ ] Error message shows lockout time
- [ ] Database: `locked_until` set
- [ ] Database: `failed_login_attempts` = 5

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case M4: Session Persistence
**Priority**: P0 (Critical)

**Steps**:
1. Login successfully
2. Verify di dashboard
3. Refresh page (F5)
4. Check localStorage di DevTools

**Expected Results**:
- [ ] After refresh, tetap logged in
- [ ] Tidak redirect ke login
- [ ] User data masih ada
- [ ] Navbar masih shows user info
- [ ] localStorage contains: `auth_token`, `refresh_token`, `user_data`

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case M5: Protected Routes
**Priority**: P0 (Critical)

**Steps**:
1. Logout dari aplikasi
2. Manually navigate ke `http://localhost:5173/dashboard`
3. Check URL

**Expected Results**:
- [ ] Auto-redirect ke `/login`
- [ ] Query param added: `?redirect=/dashboard`
- [ ] After login, redirect back to `/dashboard`

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case M6: Logout Flow
**Priority**: P0 (Critical)

**Steps**:
1. Login successfully
2. Click user avatar di navbar
3. Click "Keluar" button
4. Check localStorage

**Expected Results**:
- [ ] Dropdown menu appears
- [ ] Redirect ke `/login`
- [ ] localStorage cleared (no auth_token)
- [ ] Haptic feedback triggered
- [ ] Cannot access protected routes

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case M7: Token Expiry & Auto Refresh
**Priority**: P1 (High)

**Prerequisites**:
- Set `JWT_EXPIRY=1m` di backend .env untuk testing

**Steps**:
1. Login successfully
2. Wait 1+ minutes
3. Navigate ke Profile page (triggers API call)
4. Check Network tab di DevTools

**Expected Results**:
- [ ] API call returns 401 initially
- [ ] Refresh token API called automatically
- [ ] New token saved
- [ ] Original API call retried
- [ ] Page loads successfully
- [ ] User tidak notice interruption

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case M8: Role-Based Dashboard
**Priority**: P1 (High)

**Steps**:
1. Login as ADMIN (NIP: 99999)
2. Check redirect destination
3. Logout
4. Login as STAFF (create test user first)
5. Check redirect destination

**Expected Results**:
- [ ] ADMIN redirects to `/dashboard/admin`
- [ ] STAFF redirects to `/dashboard/staff`
- [ ] Admin dashboard shows: stats, quick actions, activity
- [ ] Staff dashboard shows: tasks, performance, notifications

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case M9: Profile Page
**Priority**: P1 (High)

**Steps**:
1. Login successfully
2. Click user dropdown → "Profile"
3. Verify profile data

**Expected Results**:
- [ ] Redirect ke `/profile`
- [ ] Avatar shows user initial (2 letters)
- [ ] NIP displayed correctly
- [ ] Full name displayed
- [ ] Email displayed
- [ ] Phone displayed (if exists)
- [ ] Role badge dengan correct color
- [ ] Department badge
- [ ] Shift badge
- [ ] Join date formatted correctly
- [ ] Last login timestamp

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case M10: Mobile Responsive
**Priority**: P1 (High)

**Steps**:
1. Open DevTools → Toggle device toolbar
2. Test di iPhone SE (375x667)
3. Test di iPad (768x1024)
4. Test login flow
5. Test dashboard navigation

**Expected Results**:
- [ ] Login form responsive
- [ ] Buttons accessible
- [ ] Navbar hamburger menu works
- [ ] Dashboard cards stack vertically
- [ ] No horizontal scroll
- [ ] Touch targets ≥ 44x44px
- [ ] Glass effect visible
- [ ] Animations smooth

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

## 🔐 Security Testing

### Test Case S1: XSS Prevention
**Priority**: P0 (Critical)

**Steps**:
1. Try input `<script>alert('xss')</script>` di NIP field
2. Try input `<img src=x onerror=alert('xss')>` di password
3. Submit form

**Expected Results**:
- [ ] NIP field only accepts numbers
- [ ] Script tidak executed
- [ ] No alert popup
- [ ] Input sanitized

**Status**: [ ] Pass [ ] Fail

---

### Test Case S2: SQL Injection Prevention
**Priority**: P0 (Critical)

**Steps**:
1. Try input `' OR '1'='1` di NIP field
2. Try input `admin'--` di NIP field
3. Submit login

**Expected Results**:
- [ ] Login fails safely
- [ ] No database error exposed
- [ ] GORM parameterized queries prevent injection
- [ ] Error message generic

**Status**: [ ] Pass [ ] Fail

---

### Test Case S3: Token Tampering
**Priority**: P0 (Critical)

**Steps**:
1. Login dan get token
2. Open DevTools → Application → Local Storage
3. Modify `auth_token` value manually
4. Try access `/dashboard`

**Expected Results**:
- [ ] 401 Unauthorized response
- [ ] Auto-logout triggered
- [ ] Redirect ke `/login`
- [ ] Error message: "Token tidak valid"

**Status**: [ ] Pass [ ] Fail

---

### Test Case S4: Password Visibility
**Priority**: P1 (High)

**Steps**:
1. Input password
2. Check Network tab → Payload
3. Check localStorage
4. Check database

**Expected Results**:
- [ ] Password sent via HTTPS (production)
- [ ] Password not visible di localStorage
- [ ] Database stores bcrypt hash only
- [ ] Hash different from plaintext
- [ ] No password in API responses

**Status**: [ ] Pass [ ] Fail

---

### Test Case S5: CORS Protection
**Priority**: P1 (High)

**Steps**:
1. Try API call dari different origin (e.g., `http://evil.com`)
2. Check browser console

**Expected Results**:
- [ ] CORS error in console
- [ ] Request blocked by browser
- [ ] Only `http://localhost:5173` allowed in dev
- [ ] Production domain configured

**Status**: [ ] Pass [ ] Fail

---

## ⚡ Performance Testing

### Test Case P1: Password Hashing Performance
**Objective**: Verify bcrypt performance acceptable

**Command**:
```bash
cd backend
go test -bench=BenchmarkHashPassword ./tests/unit/services
```

**Expected Results**:
- [ ] Hashing time: 150-250ms per operation
- [ ] Memory usage: < 20KB per operation
- [ ] Consistent performance across runs

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case P2: Login Response Time
**Objective**: Verify login endpoint performance

**Command**:
```bash
ab -n 100 -c 10 -p login.json -T application/json \
  http://localhost:8080/api/auth/login
```

**Expected Results**:
- [ ] Average response time: < 500ms
- [ ] 95th percentile: < 1000ms
- [ ] No failed requests
- [ ] Throughput: > 20 req/s

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case P3: Frontend Bundle Size
**Objective**: Verify bundle size acceptable

**Command**:
```bash
cd frontend
yarn build
du -sh dist/assets/*.js
```

**Expected Results**:
- [ ] Total JS bundle: < 500KB gzipped
- [ ] Main chunk: < 200KB
- [ ] Vendor chunk: < 300KB
- [ ] Lazy-loaded routes working

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

### Test Case P4: Dashboard Load Time
**Objective**: Verify dashboard loads quickly

**Steps**:
1. Login successfully
2. Open DevTools → Performance tab
3. Record dashboard load
4. Check metrics

**Expected Results**:
- [ ] First Contentful Paint: < 1s
- [ ] Time to Interactive: < 2s
- [ ] Total load time: < 3s
- [ ] No layout shifts

**Actual Result**: _______________

**Status**: [ ] Pass [ ] Fail

---

## 🌐 Browser Compatibility

### Test Case B1: Cross-Browser Testing

| Browser | Version | Login | Dashboard | Animations | Status |
|---------|---------|-------|-----------|------------|--------|
| Chrome | Latest | [ ] | [ ] | [ ] | [ ] Pass [ ] Fail |
| Firefox | Latest | [ ] | [ ] | [ ] | [ ] Pass [ ] Fail |
| Safari | Latest | [ ] | [ ] | [ ] | [ ] Pass [ ] Fail |
| Edge | Latest | [ ] | [ ] | [ ] | [ ] Pass [ ] Fail |

---

## 📱 Mobile Testing

### Test Case MB1: iOS Testing

| Device | iOS Version | Login | Navigation | Haptic | Status |
|--------|-------------|-------|------------|--------|--------|
| iPhone 12 | 16+ | [ ] | [ ] | [ ] | [ ] Pass [ ] Fail |
| iPhone SE | 15+ | [ ] | [ ] | [ ] | [ ] Pass [ ] Fail |
| iPad Pro | 16+ | [ ] | [ ] | [ ] | [ ] Pass [ ] Fail |

### Test Case MB2: Android Testing

| Device | Android Version | Login | Navigation | Vibration | Status |
|--------|-----------------|-------|------------|-----------|--------|
| Pixel 6 | 13+ | [ ] | [ ] | [ ] | [ ] Pass [ ] Fail |
| Samsung S21 | 12+ | [ ] | [ ] | [ ] | [ ] Pass [ ] Fail |
| OnePlus 9 | 11+ | [ ] | [ ] | [ ] | [ ] Pass [ ] Fail |

---

## 📊 Test Summary

### Coverage Summary

| Category | Total Tests | Passed | Failed | Coverage |
|----------|-------------|--------|--------|----------|
| Backend Unit | 18 | 18 | 0 | 95%+ |
| Frontend Unit | 15 | 15 | 0 | 95%+ |
| Integration | 4 | - | - | - |
| Manual | 10 | - | - | - |
| Security | 5 | - | - | - |
| Performance | 4 | - | - | - |
| Browser | 4 | - | - | - |
| Mobile | 6 | - | - | - |
| **TOTAL** | **66** | **33** | **0** | **~50%** |

### Test Execution Status

- ✅ **Unit Tests**: Complete (33/33 passing)
- 🔄 **Integration Tests**: Ready to run
- 🔄 **Manual Tests**: Checklist ready
- 🔄 **Security Tests**: Checklist ready
- 🔄 **Performance Tests**: Checklist ready
- 🔄 **Browser Tests**: Pending
- 🔄 **Mobile Tests**: Pending

---

## 🚀 Test Execution Commands

```bash
# Run all backend tests
cd backend && go test -v ./tests/...

# Run all frontend tests
cd frontend && yarn test

# Run with coverage
cd backend && go test -cover ./tests/...
cd frontend && yarn test:coverage

# Run integration tests
# (Manual execution via curl commands above)

# Run performance tests
cd backend && go test -bench=. ./tests/...
ab -n 100 -c 10 -p login.json -T application/json http://localhost:8080/api/auth/login
```

---

## 📋 QA Sign-off

### Sprint 1 Acceptance

- [ ] All unit tests passing (33/33)
- [ ] Integration tests executed successfully
- [ ] Manual testing checklist completed
- [ ] Security tests passed
- [ ] Performance benchmarks met
- [ ] Cross-browser testing completed
- [ ] Mobile testing completed
- [ ] No critical bugs found
- [ ] Documentation updated

**QA Engineer**: _______________  
**Date**: _______________  
**Sign-off**: [ ] Approved [ ] Rejected

**Notes**: _______________

---

## 🔗 Related Documentation

- **Feature Documentation**: [Authentication System](../features/AUTHENTICATION.md)
- **User Journeys**: [Authentication User Journeys](../guides/authentication-user-journeys.md)
- **Testing Guide**: [TESTING_GUIDE.md](../../TESTING_GUIDE.md)
- **API Documentation**: [Authentication API](../api/authentication.md)

---

**Last Updated**: 27 Desember 2025  
**Version**: 1.0.0 - Sprint 1  
**Status**: ✅ Test Plan Complete
