# 🧪 Testing Guide - Sirine Go Authentication System

**Last Updated**: 27 Desember 2025  
**Sprint**: 1 - Authentication System

---

## 📋 Overview

Testing Guide merupakan comprehensive documentation untuk testing authentication system, yaitu: unit tests, integration tests, manual testing, dan performance testing. Panduan ini mencakup setup testing tools, test cases, dan best practices untuk memastikan code quality dan reliability.

---

## 🎯 Testing Strategy

### Testing Pyramid

```
           /\
          /  \
         / E2E \
        /--------\
       /          \
      / Integration \
     /--------------\
    /                \
   /   Unit Tests     \
  /____________________\
```

1. **Unit Tests** (70%): Test individual functions/methods
2. **Integration Tests** (20%): Test component interactions
3. **E2E Tests** (10%): Test complete user flows

---

## 🔧 Backend Testing (Go)

### Setup

Testing tools sudah built-in di Go, tidak perlu install tambahan.

### Running Tests

```bash
# Run all tests
cd backend
go test ./...

# Run tests dengan verbose output
go test -v ./...

# Run tests dengan coverage
go test -cover ./...

# Run tests dengan detailed coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test file
go test -v ./services/password_service_test.go

# Run specific test function
go test -v -run TestHashPassword ./services

# Run benchmarks
go test -bench=. ./services
```

### Test Files Created

#### 1. Password Service Tests

**File**: `backend/services/password_service_test.go`

**Test Cases**:
- ✅ `TestHashPassword` - Password hashing dengan bcrypt
- ✅ `TestHashPasswordEmpty` - Error handling untuk empty password
- ✅ `TestVerifyPassword` - Password verification dengan berbagai cases
- ✅ `TestValidatePasswordPolicy` - Password policy enforcement
- ✅ `TestGetPasswordStrength` - Password strength calculation
- ✅ `TestValidateAndHash` - Combined validation dan hashing
- ✅ `BenchmarkHashPassword` - Performance benchmarking
- ✅ `BenchmarkVerifyPassword` - Verification performance

**Run**:
```bash
cd backend
go test -v ./services/password_service_test.go ./services/password_service.go
```

**Expected Output**:
```
=== RUN   TestHashPassword
--- PASS: TestHashPassword (0.20s)
=== RUN   TestVerifyPassword
--- PASS: TestVerifyPassword (0.40s)
=== RUN   TestValidatePasswordPolicy
--- PASS: TestValidatePasswordPolicy (0.00s)
...
PASS
ok      sirine-go/backend/services      0.650s
```

#### 2. User Model Tests

**File**: `backend/models/user_test.go`

**Test Cases**:
- ✅ `TestUserIsLocked` - User lock status checking
- ✅ `TestUserIsActive` - User active status checking
- ✅ `TestUserHasRole` - Role checking dengan multiple roles
- ✅ `TestUserIsAdmin` - Admin role checking
- ✅ `TestUserToSafeUser` - Safe user conversion
- ✅ `TestUserSessionIsValid` - Session validation
- ✅ `TestPasswordResetTokenIsValid` - Reset token validation

**Run**:
```bash
cd backend
go test -v ./models/user_test.go ./models/*.go
```

### Writing New Backend Tests

**Template**:
```go
package services

import (
    "testing"
)

func TestYourFunction(t *testing.T) {
    // Arrange - Setup test data
    service := NewYourService()
    input := "test-input"
    
    // Act - Execute function
    result, err := service.YourFunction(input)
    
    // Assert - Verify results
    if err != nil {
        t.Fatalf("YourFunction gagal: %v", err)
    }
    
    if result != "expected" {
        t.Errorf("Expected 'expected', got '%s'", result)
    }
}

// Table-driven test
func TestYourFunctionCases(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"Valid input", "test", "result", false},
        {"Invalid input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := YourFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Error = %v, wantErr %v", err, tt.wantErr)
            }
            if result != tt.expected {
                t.Errorf("Got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

---

## 🎨 Frontend Testing (Vue + Vitest)

### Setup

```bash
cd frontend
yarn add -D vitest @vitest/ui @vue/test-utils jsdom happy-dom
```

### Running Tests

```bash
# Run all tests
yarn test

# Run tests dengan watch mode
yarn test --watch

# Run tests dengan UI
yarn test:ui

# Run tests dengan coverage
yarn test:coverage

# Run specific test file
yarn test auth.spec.js
```

### Test Files Created

#### 1. Auth Store Tests

**File**: `frontend/src/stores/auth.spec.js`

**Test Cases**:
- ✅ Initial state (user null, not authenticated)
- ✅ `setAuth` - Set auth data dan localStorage
- ✅ `clearAuth` - Clear auth data dan localStorage
- ✅ `restoreAuth` - Restore dari localStorage
- ✅ `hasRole` - Role checking
- ✅ `isAdmin` - Admin checking
- ✅ Computed properties (userRole, userDepartment, requirePasswordChange)

**Run**:
```bash
cd frontend
yarn test auth.spec.js
```

**Expected Output**:
```
 ✓ src/stores/auth.spec.js (15 tests) 150ms
   ✓ Auth Store
     ✓ Initial State
       ✓ harus memiliki user null di awal
       ✓ harus tidak authenticated di awal
     ✓ setAuth
       ✓ harus set auth data dengan benar
       ✓ harus menyimpan ke localStorage
     ...

Test Files  1 passed (1)
     Tests  15 passed (15)
  Start at  14:30:00
  Duration  1.2s
```

### Configuration Files

#### vitest.config.js
```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/tests/setup.js'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
    },
  },
})
```

#### src/tests/setup.js
```javascript
import { vi } from 'vitest'

// Mock localStorage
global.localStorage = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}

// Mock navigator.vibrate
global.navigator.vibrate = vi.fn()
```

### Writing New Frontend Tests

**Template**:
```javascript
import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import YourComponent from './YourComponent.vue'

describe('YourComponent', () => {
  let wrapper

  beforeEach(() => {
    wrapper = mount(YourComponent, {
      props: {
        // your props
      },
    })
  })

  it('should render correctly', () => {
    expect(wrapper.exists()).toBe(true)
    expect(wrapper.text()).toContain('Expected Text')
  })

  it('should handle click event', async () => {
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted()).toHaveProperty('click')
  })
})
```

---

## 🔗 Integration Testing

### API Integration Tests

**Manual Testing dengan curl**:

```bash
# 1. Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "99999",
    "password": "Admin@123",
    "remember_me": false
  }'

# Save token dari response
TOKEN="eyJhbGc..."

# 2. Get current user
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $TOKEN"

# 3. Logout
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer $TOKEN"

# 4. Try access protected route tanpa token (should fail)
curl http://localhost:8080/api/auth/me
```

### Database Integration

**Test dengan Go**:
```go
// backend/tests/integration/auth_integration_test.go
func TestLoginIntegration(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    defer cleanupTestDB(db)
    
    // Create test user
    user := createTestUser(db)
    
    // Test login
    authService := services.NewAuthService(db, cfg)
    response, err := authService.Login(services.LoginRequest{
        NIP:      user.NIP,
        Password: "TestPassword123!",
    }, "127.0.0.1", "test-agent")
    
    if err != nil {
        t.Fatalf("Login failed: %v", err)
    }
    
    if response.Token == "" {
        t.Error("Token should not be empty")
    }
}
```

---

## 🌐 Manual Testing

### Authentication Flow Testing

#### Test Case 1: Login Success
1. Buka `http://localhost:5173`
2. Input NIP: `99999`, Password: `Admin@123`
3. Klik "Masuk"
4. **Expected**: 
   - ✅ Redirect ke `/dashboard/admin`
   - ✅ Navbar shows "Administrator"
   - ✅ Token saved di localStorage
   - ✅ Haptic feedback (vibration)

#### Test Case 2: Login Failed
1. Input wrong credentials
2. Klik "Masuk"
3. **Expected**:
   - ✅ Error message "NIP atau password salah"
   - ✅ Card shake animation
   - ✅ Error haptic feedback
   - ✅ Form tidak clear

#### Test Case 3: Rate Limiting
1. Login dengan wrong password 5x
2. **Expected**:
   - ✅ Error "Akun terkunci selama 15 menit"
   - ✅ Cannot login even dengan correct password
   - ✅ `locked_until` timestamp set di database

#### Test Case 4: Session Persistence
1. Login successfully
2. Refresh page (F5)
3. **Expected**:
   - ✅ Tetap logged in
   - ✅ Tidak redirect ke login
   - ✅ User data masih ada

#### Test Case 5: Protected Routes
1. Logout
2. Try access `/dashboard` via URL
3. **Expected**:
   - ✅ Auto-redirect ke `/login`
   - ✅ Query param `?redirect=/dashboard` added

#### Test Case 6: Token Expiry
1. Login
2. Wait 15+ minutes (atau set JWT_EXPIRY=1m)
3. Make API call (e.g., visit profile)
4. **Expected**:
   - ✅ Auto token refresh
   - ✅ Seamless experience
   - ✅ New token saved

### Security Testing

#### Test Case 7: XSS Prevention
1. Try input `<script>alert('xss')</script>` di NIP field
2. **Expected**:
   - ✅ Only numbers allowed (validation)
   - ✅ Script tidak executed

#### Test Case 8: SQL Injection Prevention
1. Try input `' OR '1'='1` di NIP field
2. **Expected**:
   - ✅ GORM parameterized queries prevent injection
   - ✅ Login fails safely

#### Test Case 9: Token Tampering
1. Login dan get token
2. Modify token manually di localStorage
3. Try access protected route
4. **Expected**:
   - ✅ 401 Unauthorized
   - ✅ Auto-logout
   - ✅ Redirect ke login

---

## 📊 Performance Testing

### Backend Performance

```bash
# Benchmark password hashing
cd backend
go test -bench=BenchmarkHashPassword -benchmem ./services

# Expected output:
# BenchmarkHashPassword-8    5   200000000 ns/op   15000 B/op   10 allocs/op
# ~200ms per hash (bcrypt cost 12)
```

### Load Testing dengan Apache Bench

```bash
# Test login endpoint
ab -n 100 -c 10 -p login.json -T application/json \
  http://localhost:8080/api/auth/login

# login.json content:
# {"nip":"99999","password":"Admin@123","remember_me":false}

# Expected:
# Requests per second: 50-100 req/s
# Time per request: 10-20ms (avg)
```

### Frontend Performance

```bash
# Build production
cd frontend
yarn build

# Check bundle size
du -sh dist/assets/*.js

# Expected:
# < 500KB gzipped total
```

---

## ✅ Testing Checklist

### Before Commit

- [ ] All unit tests pass (`go test ./...`)
- [ ] Frontend tests pass (`yarn test`)
- [ ] No linter errors
- [ ] Code coverage > 70%
- [ ] Manual testing completed
- [ ] Performance benchmarks acceptable

### Before Deploy

- [ ] Integration tests pass
- [ ] Security tests pass
- [ ] Load testing completed
- [ ] Browser compatibility verified
- [ ] Mobile testing completed
- [ ] Documentation updated

---

## 🎯 Coverage Goals

### Current Coverage (Sprint 1)

**Backend**:
- Password Service: 95%+ (11 test cases)
- User Models: 90%+ (7 test cases)
- Auth Service: 0% (TODO: Sprint 2)
- Handlers: 0% (TODO: Sprint 2)

**Frontend**:
- Auth Store: 95%+ (15 test cases)
- Composables: 0% (TODO: Sprint 2)
- Components: 0% (TODO: Sprint 2)

**Target for Sprint 2**:
- Backend: 80%+ coverage
- Frontend: 70%+ coverage

---

## 🚀 CI/CD Integration

### GitHub Actions Workflow

**File**: `.github/workflows/test.yml`

```yaml
name: Tests

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      - name: Run tests
        run: |
          cd backend
          go test -v -cover ./...

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        run: |
          cd frontend
          yarn install
      - name: Run tests
        run: |
          cd frontend
          yarn test
```

---

## 📝 Best Practices

### Backend Testing

1. **Use Table-Driven Tests**: Test multiple cases efficiently
2. **Mock External Dependencies**: Database, HTTP clients, etc.
3. **Test Edge Cases**: Empty inputs, nil pointers, errors
4. **Benchmark Performance**: Critical functions like hashing
5. **Use Subtests**: Organize related tests with `t.Run()`

### Frontend Testing

1. **Test User Behavior**: Not implementation details
2. **Mock API Calls**: Use `vi.fn()` untuk axios
3. **Test Accessibility**: ARIA labels, keyboard navigation
4. **Test Error States**: Loading, error messages
5. **Snapshot Testing**: For UI components (optional)

### General

1. **Write Tests First** (TDD): Red → Green → Refactor
2. **Keep Tests Simple**: One assertion per test (ideally)
3. **Use Descriptive Names**: Test name = documentation
4. **Clean Up**: `beforeEach`, `afterEach` untuk cleanup
5. **Run Tests Often**: Before commit, before push

---

## 🐛 Debugging Tests

### Go Tests

```bash
# Run with verbose + print statements
go test -v ./services -run TestHashPassword

# Run with debugger (Delve)
dlv test ./services -- -test.run TestHashPassword
```

### Frontend Tests

```bash
# Run with UI untuk debugging
yarn test:ui

# Run specific test dengan watch
yarn test --watch auth.spec.js

# Add console.log di test
it('should work', () => {
  console.log('Debug:', wrapper.html())
  expect(true).toBe(true)
})
```

---

## 📚 Resources

### Go Testing
- [Go Testing Package](https://pkg.go.dev/testing)
- [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Go Test Examples](https://go.dev/doc/tutorial/add-a-test)

### Vitest
- [Vitest Documentation](https://vitest.dev/)
- [Vue Test Utils](https://test-utils.vuejs.org/)
- [Testing Library](https://testing-library.com/docs/vue-testing-library/intro/)

---

## 🎓 Next Steps

### Sprint 2 Testing Goals

1. **Auth Service Tests**: JWT generation, validation, refresh
2. **Handler Tests**: HTTP request/response testing
3. **Middleware Tests**: Auth middleware, role middleware
4. **Component Tests**: Login form, Navbar, Dashboard
5. **E2E Tests**: Complete user flows dengan Playwright

---

## 📞 Support

**Developer**: Zulfikar Hidayatullah  
**Phone**: +62 857-1583-8733  
**Testing Questions**: Refer to this guide atau contact developer

---

**Last Updated**: 27 Desember 2025  
**Version**: 1.0.0 - Sprint 1  
**Status**: ✅ Testing Framework Ready
