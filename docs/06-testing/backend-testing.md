# 🔧 Backend Testing Guide

Panduan lengkap untuk testing backend Go aplikasi dengan unit tests dan integration tests.

---

## 📋 Daftar Isi

1. [Setup Testing Environment](#setup-testing-environment)
2. [Unit Testing Services](#unit-testing-services)
3. [Unit Testing Handlers](#unit-testing-handlers)
4. [Running Tests](#running-tests)
5. [Best Practices](#best-practices)

---

## ⚙️ Setup Testing Environment

### **Required Dependencies**

Testing framework sudah included dalam Go standard library, namun kami menggunakan testify untuk assertions yang lebih readable:

```bash
cd backend

# Install testify untuk better assertions
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
```

### **Test Database Setup**

Create separate test database untuk isolate test data:

```sql
-- Login ke MySQL
mysql -u root -p

-- Create test database
CREATE DATABASE sirine_go_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON sirine_go_test.* TO 'your_user'@'localhost';
FLUSH PRIVILEGES;
```

**Update `.env.test`:**

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=sirine_go_test
```

---

## 🧪 Unit Testing Services

Service layer testing focuses pada business logic isolation tanpa HTTP context.

### **Service Test Structure**

Create test file dengan nama `*_test.go` di folder yang sama dengan file yang di-test:

```
backend/services/
├── user_service.go
├── user_service_test.go      # Test file
├── auth_service.go
└── auth_service_test.go       # Test file
```

### **Example: User Service Test**

```go
// backend/services/user_service_test.go
package services

import (
    "testing"
    "sirine-go/backend/database"
    "sirine-go/backend/models"
    "github.com/stretchr/testify/assert"
)

// Setup function untuk initialize test database
func setupTestDB(t *testing.T) {
    // Connect ke test database
    database.ConnectTest()
    
    // Run migrations
    database.DB.AutoMigrate(&models.User{})
    
    // Cleanup setelah test selesai
    t.Cleanup(func() {
        database.DB.Migrator().DropTable(&models.User{})
        database.CloseTest()
    })
}

func TestUserService_Create(t *testing.T) {
    setupTestDB(t)
    
    service := NewUserService()
    
    // Test data
    user := &models.User{
        Name:       "Test User",
        Email:      "test@example.com",
        Password:   "hashedpassword",
        Role:       "user",
        Department: "IT",
        IsActive:   true,
    }
    
    // Test create
    err := service.Create(user)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    assert.Equal(t, "Test User", user.Name)
    assert.Equal(t, "test@example.com", user.Email)
}

func TestUserService_GetByID(t *testing.T) {
    setupTestDB(t)
    
    service := NewUserService()
    
    // Create test user first
    user := &models.User{
        Name:  "Test User",
        Email: "test@example.com",
        Role:  "user",
    }
    service.Create(user)
    
    // Test get by ID
    found, err := service.GetByID(user.ID)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, found)
    assert.Equal(t, user.ID, found.ID)
    assert.Equal(t, "Test User", found.Name)
}

func TestUserService_GetByEmail(t *testing.T) {
    setupTestDB(t)
    
    service := NewUserService()
    
    // Create test user
    user := &models.User{
        Name:  "Test User",
        Email: "test@example.com",
        Role:  "user",
    }
    service.Create(user)
    
    // Test get by email
    found, err := service.GetByEmail("test@example.com")
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, found)
    assert.Equal(t, "test@example.com", found.Email)
}

func TestUserService_Update(t *testing.T) {
    setupTestDB(t)
    
    service := NewUserService()
    
    // Create test user
    user := &models.User{
        Name:  "Old Name",
        Email: "test@example.com",
        Role:  "user",
    }
    service.Create(user)
    
    // Update user
    user.Name = "New Name"
    err := service.Update(user)
    
    // Verify update
    updated, _ := service.GetByID(user.ID)
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "New Name", updated.Name)
}

func TestUserService_Delete(t *testing.T) {
    setupTestDB(t)
    
    service := NewUserService()
    
    // Create test user
    user := &models.User{
        Name:  "Test User",
        Email: "test@example.com",
        Role:  "user",
    }
    service.Create(user)
    
    // Delete user (soft delete)
    err := service.Delete(user.ID)
    
    // Try to get deleted user
    found, getErr := service.GetByID(user.ID)
    
    // Assertions
    assert.NoError(t, err)
    assert.Error(t, getErr) // Should error karena soft deleted
    assert.Nil(t, found)
}

func TestUserService_List_WithFilters(t *testing.T) {
    setupTestDB(t)
    
    service := NewUserService()
    
    // Create multiple test users
    users := []*models.User{
        {Name: "Admin User", Email: "admin@test.com", Role: "admin", Department: "IT"},
        {Name: "Manager User", Email: "manager@test.com", Role: "manager", Department: "IT"},
        {Name: "Regular User", Email: "user@test.com", Role: "user", Department: "HR"},
    }
    
    for _, u := range users {
        service.Create(u)
    }
    
    // Test filter by role
    filters := map[string]interface{}{
        "role": "admin",
    }
    result, err := service.List(filters, 1, 10)
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 1, len(result))
    assert.Equal(t, "admin", result[0].Role)
    
    // Test filter by department
    filters = map[string]interface{}{
        "department": "IT",
    }
    result, err = service.List(filters, 1, 10)
    
    assert.NoError(t, err)
    assert.Equal(t, 2, len(result))
}
```

---

## 🌐 Unit Testing Handlers

Handler testing focuses pada HTTP request/response handling dengan Gin test mode.

### **Handler Test Structure**

```
backend/handlers/
├── user_handler.go
├── user_handler_test.go      # Test file
├── auth_handler.go
└── auth_handler_test.go       # Test file
```

### **Example: User Handler Test**

```go
// backend/handlers/user_handler_test.go
package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "bytes"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "sirine-go/backend/services"
    "sirine-go/backend/database"
    "sirine-go/backend/models"
)

func setupTestRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    return r
}

func setupTestHandler(t *testing.T) (*UserHandler, *gin.Engine) {
    // Setup database
    database.ConnectTest()
    database.DB.AutoMigrate(&models.User{})
    
    t.Cleanup(func() {
        database.DB.Migrator().DropTable(&models.User{})
        database.CloseTest()
    })
    
    // Create handler dengan service
    service := services.NewUserService()
    handler := NewUserHandler(service)
    router := setupTestRouter()
    
    return handler, router
}

func TestUserHandler_GetAll(t *testing.T) {
    handler, router := setupTestHandler(t)
    
    // Register route
    router.GET("/api/users", handler.GetAll)
    
    // Create test request
    req, _ := http.NewRequest("GET", "/api/users", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "data")
    assert.Contains(t, w.Body.String(), "meta")
}

func TestUserHandler_GetByID(t *testing.T) {
    handler, router := setupTestHandler(t)
    
    // Create test user first
    service := services.NewUserService()
    user := &models.User{
        Name:  "Test User",
        Email: "test@example.com",
        Role:  "user",
    }
    service.Create(user)
    
    // Register route
    router.GET("/api/users/:id", handler.GetByID)
    
    // Test request
    req, _ := http.NewRequest("GET", "/api/users/"+string(user.ID), nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    
    data := response["data"].(map[string]interface{})
    assert.Equal(t, "Test User", data["name"])
}

func TestUserHandler_Create(t *testing.T) {
    handler, router := setupTestHandler(t)
    
    // Register route
    router.POST("/api/users", handler.Create)
    
    // Test data
    data := map[string]interface{}{
        "name":       "New User",
        "email":      "newuser@test.com",
        "password":   "password123",
        "role":       "user",
        "department": "IT",
        "is_active":  true,
    }
    jsonData, _ := json.Marshal(data)
    
    // Test request
    req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusCreated, w.Code)
    assert.Contains(t, w.Body.String(), "berhasil dibuat")
    assert.Contains(t, w.Body.String(), "newuser@test.com")
}

func TestUserHandler_Create_ValidationError(t *testing.T) {
    handler, router := setupTestHandler(t)
    router.POST("/api/users", handler.Create)
    
    // Invalid data (missing required fields)
    data := map[string]interface{}{
        "name": "Incomplete User",
        // Missing email, password, role
    }
    jsonData, _ := json.Marshal(data)
    
    // Test request
    req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusBadRequest, w.Code)
    assert.Contains(t, w.Body.String(), "error")
}

func TestUserHandler_Update(t *testing.T) {
    handler, router := setupTestHandler(t)
    
    // Create test user
    service := services.NewUserService()
    user := &models.User{
        Name:  "Old Name",
        Email: "test@example.com",
        Role:  "user",
    }
    service.Create(user)
    
    // Register route
    router.PUT("/api/users/:id", handler.Update)
    
    // Update data
    data := map[string]interface{}{
        "name": "Updated Name",
    }
    jsonData, _ := json.Marshal(data)
    
    // Test request
    req, _ := http.NewRequest("PUT", "/api/users/"+string(user.ID), bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "berhasil diupdate")
}

func TestUserHandler_Delete(t *testing.T) {
    handler, router := setupTestHandler(t)
    
    // Create test user
    service := services.NewUserService()
    user := &models.User{
        Name:  "To Be Deleted",
        Email: "delete@test.com",
        Role:  "user",
    }
    service.Create(user)
    
    // Register route
    router.DELETE("/api/users/:id", handler.Delete)
    
    // Test request
    req, _ := http.NewRequest("DELETE", "/api/users/"+string(user.ID), nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "berhasil dihapus")
}
```

---

## ▶️ Running Tests

### **Run All Tests**

```bash
cd backend

# Run all tests
go test ./...

# Run with verbose output
go test -v ./...
```

### **Run Specific Package**

```bash
# Test services only
go test ./services

# Test handlers only
go test ./handlers

# Test specific file
go test ./services/user_service_test.go
```

### **Run Specific Test Function**

```bash
# Run specific test
go test -run TestUserService_Create ./services

# Run tests matching pattern
go test -run TestUser ./...
```

### **Test Coverage**

```bash
# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out
```

**Expected Output:**

```
?       sirine-go/backend/cmd/server    [no test files]
ok      sirine-go/backend/handlers      2.134s  coverage: 73.5% of statements
ok      sirine-go/backend/services      1.876s  coverage: 81.2% of statements
ok      sirine-go/backend/models        0.543s  coverage: 65.8% of statements
```

### **Run Tests with Race Detector**

Detect race conditions dalam concurrent code:

```bash
go test -race ./...
```

---

## ✅ Best Practices

### **1. Test Isolation**

Setiap test harus independent dan tidak depend pada test lain:

```go
func TestExample(t *testing.T) {
    // Setup - create fresh state
    setupTestDB(t)
    
    // Test logic
    // ...
    
    // Cleanup - handled by t.Cleanup()
}
```

### **2. Table-Driven Tests**

Gunakan table-driven approach untuk test multiple scenarios:

```go
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   *models.User
        wantErr bool
    }{
        {
            name: "valid user",
            input: &models.User{
                Name:  "Test",
                Email: "test@example.com",
                Role:  "user",
            },
            wantErr: false,
        },
        {
            name: "missing email",
            input: &models.User{
                Name: "Test",
                Role: "user",
            },
            wantErr: true,
        },
        {
            name: "invalid role",
            input: &models.User{
                Name:  "Test",
                Email: "test@example.com",
                Role:  "invalid",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateUser(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### **3. Mock External Dependencies**

Gunakan mocks untuk external services (email, third-party APIs):

```go
// Create mock interface
type MockEmailService struct {
    mock.Mock
}

func (m *MockEmailService) SendEmail(to string, subject string, body string) error {
    args := m.Called(to, subject, body)
    return args.Error(0)
}

// Use in test
func TestPasswordReset(t *testing.T) {
    mockEmail := new(MockEmailService)
    mockEmail.On("SendEmail", "user@test.com", mock.Anything, mock.Anything).Return(nil)
    
    service := NewPasswordService(mockEmail)
    err := service.RequestReset("user@test.com")
    
    assert.NoError(t, err)
    mockEmail.AssertExpectations(t)
}
```

### **4. Test Error Cases**

Always test error scenarios, bukan hanya happy path:

```go
func TestUserService_Create_DuplicateEmail(t *testing.T) {
    setupTestDB(t)
    service := NewUserService()
    
    // Create first user
    user1 := &models.User{Email: "test@example.com"}
    service.Create(user1)
    
    // Try to create duplicate
    user2 := &models.User{Email: "test@example.com"}
    err := service.Create(user2)
    
    // Should error karena duplicate email
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "sudah terdaftar")
}
```

### **5. Descriptive Test Names**

Gunakan descriptive names yang explain what & why:

```go
// ❌ Bad
func TestCreate(t *testing.T) { }

// ✅ Good
func TestUserService_Create_ValidUser_Success(t *testing.T) { }
func TestUserService_Create_DuplicateEmail_ReturnsError(t *testing.T) { }
```

---

## 📚 Related Documentation

- [overview.md](./overview.md) - Testing strategy & manual testing
- [api-testing.md](./api-testing.md) - API & integration testing
- [performance-testing.md](./performance-testing.md) - Load & performance testing

---

## 📞 Support

Jika ada pertanyaan tentang backend testing:
- Developer: Zulfikar Hidayatullah
- Phone: +62 857-1583-8733

---

**Last Updated:** 28 Desember 2025  
**Version:** 2.0.0 (Phase 2 Restructure)
