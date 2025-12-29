# 🏗️ Service Pattern Architecture

Panduan implementasi Service Pattern dalam Sirine Go App, yaitu arsitektur yang memisahkan concerns menjadi layers untuk maintainability dan testability.

---

## 📋 Overview

Service Pattern merupakan architectural pattern yang digunakan dalam Sirine Go App untuk separation of concerns, dimana aplikasi dibagi menjadi layers yang memiliki responsibility berbeda.

### Why Service Pattern?

**Benefits:**
- **Separation of Concerns** - Setiap layer punya tanggung jawab jelas
- **Testability** - Mudah untuk unit test setiap layer
- **Maintainability** - Perubahan pada satu layer tidak affect layer lain
- **Reusability** - Business logic dapat digunakan oleh multiple handlers
- **Scalability** - Mudah untuk add new features

---

## 🔄 Architecture Layers

Service Pattern dalam Sirine Go terdiri dari 3 layers utama:

```
┌─────────────────────────────────────────┐
│         HTTP Request                    │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│  Handler Layer (Controllers)            │
│  - Parse request                        │
│  - Validate input                       │
│  - Call service                         │
│  - Return response                      │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│  Service Layer (Business Logic)         │
│  - Business rules                       │
│  - Data processing                      │
│  - Validation                           │
│  - Orchestrate operations               │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│  Model/Repository Layer (Data Access)   │
│  - Database queries                     │
│  - Data persistence                     │
│  - Relationships                        │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│         Database (MySQL)                │
└─────────────────────────────────────────┘
```

### Request Flow

```
Client Request
    ↓
Router (routes.go)
    ↓
Handler (user_handler.go)
    ↓
Service (user_service.go)
    ↓
Model/GORM (user.go)
    ↓
Database
    ↓
Response back through layers
```

---

## 📦 Layer Responsibilities

### 1. Handler Layer (Controllers)

**Location:** `internal/handlers/`

**Responsibilities:**
- Receive HTTP requests
- Parse request body dan query parameters
- Validate request format
- Call appropriate service methods
- Format dan return HTTP responses
- Handle HTTP-specific errors

**What NOT to do:**
- ❌ Business logic
- ❌ Direct database access
- ❌ Complex data processing

### 2. Service Layer (Business Logic)

**Location:** `internal/services/`

**Responsibilities:**
- Implement business rules
- Data validation dan processing
- Orchestrate multiple operations
- Transaction management
- Error handling dengan business context
- Call repository/model methods

**What NOT to do:**
- ❌ HTTP response formatting
- ❌ Request parsing
- ❌ Direct database queries (use models)

### 3. Model/Repository Layer

**Location:** `internal/models/` dan `internal/repositories/`

**Responsibilities:**
- Define data structures (GORM models)
- Database operations (CRUD)
- Complex queries
- Data relationships
- Database constraints

**What NOT to do:**
- ❌ Business logic
- ❌ HTTP handling
- ❌ Data validation (except DB constraints)

---

## 💻 Implementation Examples

### Complete Feature: Create User

#### 1. Model Definition

```go
// internal/models/user.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID           uint           `gorm:"primarykey" json:"id"`
    NIP          string         `gorm:"unique;not null" json:"nip"`
    FullName     string         `gorm:"not null" json:"full_name"`
    Email        string         `gorm:"unique;not null" json:"email"`
    Phone        string         `gorm:"not null" json:"phone"`
    Password     string         `gorm:"not null" json:"-"` // Hidden dari JSON
    Role         string         `gorm:"not null" json:"role"`
    Department   string         `gorm:"not null" json:"department"`
    Shift        string         `gorm:"not null" json:"shift"`
    Status       string         `gorm:"default:ACTIVE" json:"status"`
    TotalPoints  int            `gorm:"default:0" json:"total_points"`
    Level        string         `gorm:"default:Bronze" json:"level"`
    ProfilePhoto string         `json:"profile_photo"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
    
    // Relationships
    Achievements []Achievement `gorm:"many2many:user_achievements;" json:"achievements,omitempty"`
}

// TableName override untuk specify table name
func (User) TableName() string {
    return "users"
}

// BeforeCreate hook untuk auto-calculate level
func (u *User) BeforeCreate(tx *gorm.DB) error {
    u.Level = calculateLevel(u.TotalPoints)
    return nil
}
```

#### 2. Service Implementation

```go
// internal/services/user_service.go
package services

import (
    "errors"
    "sirine-go/internal/models"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type UserService struct {
    db *gorm.DB
}

// NewUserService untuk initialize service dengan dependency injection
func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

// CreateUserRequest untuk input validation
type CreateUserRequest struct {
    NIP        string `json:"nip" binding:"required"`
    FullName   string `json:"full_name" binding:"required"`
    Email      string `json:"email" binding:"required,email"`
    Phone      string `json:"phone" binding:"required"`
    Password   string `json:"password" binding:"required,min=8"`
    Role       string `json:"role" binding:"required"`
    Department string `json:"department" binding:"required"`
    Shift      string `json:"shift" binding:"required"`
}

// CreateUser untuk membuat user baru dengan validation dan security
func (s *UserService) CreateUser(req *CreateUserRequest) (*models.User, error) {
    // Business validation: Check jika NIP sudah exists
    var existingUser models.User
    if err := s.db.Where("nip = ?", req.NIP).First(&existingUser).Error; err == nil {
        return nil, errors.New("NIP sudah terdaftar")
    }
    
    // Check jika email sudah exists
    if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        return nil, errors.New("Email sudah terdaftar")
    }
    
    // Validate role
    validRoles := []string{"SUPER_ADMIN", "ADMIN_YAKES", "STAFF_KHAZWAL", "STAFF_FISIOTERAPI"}
    if !contains(validRoles, req.Role) {
        return nil, errors.New("Role tidak valid")
    }
    
    // Hash password untuk security
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
    if err != nil {
        return nil, errors.New("Gagal memproses password")
    }
    
    // Create user model
    user := &models.User{
        NIP:        req.NIP,
        FullName:   req.FullName,
        Email:      req.Email,
        Phone:      req.Phone,
        Password:   string(hashedPassword),
        Role:       req.Role,
        Department: req.Department,
        Shift:      req.Shift,
        Status:     "ACTIVE",
    }
    
    // Save ke database
    if err := s.db.Create(user).Error; err != nil {
        return nil, errors.New("Gagal menyimpan user ke database")
    }
    
    return user, nil
}

// ListUsers untuk mendapatkan list users dengan pagination dan filtering
func (s *UserService) ListUsers(page, pageSize int, filters map[string]interface{}) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    
    // Build query dengan GORM
    query := s.db.Model(&models.User{})
    
    // Apply filters
    if role, ok := filters["role"].(string); ok && role != "" {
        query = query.Where("role = ?", role)
    }
    
    if department, ok := filters["department"].(string); ok && department != "" {
        query = query.Where("department = ?", department)
    }
    
    if status, ok := filters["status"].(string); ok && status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Count total untuk pagination
    query.Count(&total)
    
    // Apply pagination
    offset := (page - 1) * pageSize
    err := query.Offset(offset).Limit(pageSize).Find(&users).Error
    
    return users, total, err
}

// GetUserByID untuk mendapatkan user detail berdasarkan ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
    var user models.User
    
    // Preload relationships untuk mendapatkan complete data
    err := s.db.Preload("Achievements").First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("User tidak ditemukan")
        }
        return nil, err
    }
    
    return &user, nil
}

// UpdateUser untuk update user data
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) (*models.User, error) {
    var user models.User
    
    // Check jika user exists
    if err := s.db.First(&user, id).Error; err != nil {
        return nil, errors.New("User tidak ditemukan")
    }
    
    // Update dengan transaction untuk data integrity
    err := s.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&user).Updates(updates).Error; err != nil {
            return err
        }
        return nil
    })
    
    if err != nil {
        return nil, errors.New("Gagal mengupdate user")
    }
    
    return &user, nil
}

// DeleteUser untuk soft delete user
func (s *UserService) DeleteUser(id uint) error {
    var user models.User
    
    // Check jika user exists
    if err := s.db.First(&user, id).Error; err != nil {
        return errors.New("User tidak ditemukan")
    }
    
    // Soft delete (sets deleted_at timestamp)
    if err := s.db.Delete(&user).Error; err != nil {
        return errors.New("Gagal menghapus user")
    }
    
    return nil
}

// Helper function
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

#### 3. Handler Implementation

```go
// internal/handlers/user_handler.go
package handlers

import (
    "net/http"
    "strconv"
    "sirine-go/internal/services"
    "sirine-go/internal/utils"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    service *services.UserService
}

// NewUserHandler untuk initialize handler dengan dependency injection
func NewUserHandler(service *services.UserService) *UserHandler {
    return &UserHandler{service: service}
}

// CreateUser handler untuk POST /api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req services.CreateUserRequest
    
    // Parse dan validate request body
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ValidationError(c, err)
        return
    }
    
    // Call service
    user, err := h.service.CreateUser(&req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    // Return success response
    utils.SuccessResponse(c, http.StatusCreated, "User berhasil dibuat", user)
}

// ListUsers handler untuk GET /api/users
func (h *UserHandler) ListUsers(c *gin.Context) {
    // Parse query parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    
    // Build filters dari query params
    filters := map[string]interface{}{
        "role":       c.Query("role"),
        "department": c.Query("department"),
        "status":     c.Query("status"),
    }
    
    // Call service
    users, total, err := h.service.ListUsers(page, pageSize, filters)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data users")
        return
    }
    
    // Return paginated response
    utils.PaginatedResponse(c, http.StatusOK, users, total, page, pageSize)
}

// GetUser handler untuk GET /api/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
    // Parse ID dari URL parameter
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
        return
    }
    
    // Call service
    user, err := h.service.GetUserByID(uint(id))
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, err.Error())
        return
    }
    
    // Return response
    utils.SuccessResponse(c, http.StatusOK, "User ditemukan", user)
}

// UpdateUser handler untuk PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
    // Parse ID
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
        return
    }
    
    // Parse update data
    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        utils.ValidationError(c, err)
        return
    }
    
    // Call service
    user, err := h.service.UpdateUser(uint(id), updates)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    // Return response
    utils.SuccessResponse(c, http.StatusOK, "User berhasil diupdate", user)
}

// DeleteUser handler untuk DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
    // Parse ID
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
        return
    }
    
    // Call service
    if err := h.service.DeleteUser(uint(id)); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    // Return response
    utils.SuccessResponse(c, http.StatusOK, "User berhasil dihapus", nil)
}
```

#### 4. Route Registration

```go
// internal/routes/routes.go
package routes

import (
    "sirine-go/internal/handlers"
    "sirine-go/internal/middleware"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handlers.UserHandler) {
    // Public routes
    api := r.Group("/api")
    
    // Protected routes dengan authentication
    protected := api.Group("")
    protected.Use(middleware.AuthMiddleware())
    {
        // Admin-only routes
        admin := protected.Group("")
        admin.Use(middleware.RequireRole("SUPER_ADMIN", "ADMIN_YAKES"))
        {
            admin.POST("/users", userHandler.CreateUser)
            admin.GET("/users", userHandler.ListUsers)
            admin.GET("/users/:id", userHandler.GetUser)
            admin.PUT("/users/:id", userHandler.UpdateUser)
            admin.DELETE("/users/:id", userHandler.DeleteUser)
        }
    }
}
```

#### 5. Dependency Injection (Main)

```go
// cmd/server/main.go
package main

import (
    "sirine-go/internal/config"
    "sirine-go/internal/database"
    "sirine-go/internal/handlers"
    "sirine-go/internal/services"
    "sirine-go/internal/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize config
    cfg := config.LoadConfig()
    
    // Initialize database
    db := database.Connect(cfg)
    
    // Initialize services dengan dependency injection
    userService := services.NewUserService(db)
    
    // Initialize handlers dengan dependency injection
    userHandler := handlers.NewUserHandler(userService)
    
    // Setup Gin router
    r := gin.Default()
    
    // Setup routes
    routes.SetupRoutes(r, userHandler)
    
    // Start server
    r.Run(":8080")
}
```

---

## ✅ Best Practices

### 1. Dependency Injection

**Do:** Inject dependencies melalui constructor
```go
func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}
```

**Don't:** Use global variables
```go
// ❌ Bad
var globalDB *gorm.DB

func CreateUser() {
    globalDB.Create(...)
}
```

### 2. Error Handling

**Do:** Return descriptive errors dengan context
```go
if err != nil {
    return nil, errors.New("NIP sudah terdaftar")
}
```

**Don't:** Expose internal errors
```go
// ❌ Bad
return nil, err  // Might expose database errors
```

### 3. Transaction Management

**Do:** Use transactions untuk multiple operations
```go
err := s.db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err  // Rollback
    }
    if err := tx.Create(&log).Error; err != nil {
        return err  // Rollback
    }
    return nil  // Commit
})
```

### 4. Input Validation

**Do:** Validate di service layer
```go
if req.NIP == "" {
    return nil, errors.New("NIP tidak boleh kosong")
}
```

### 5. Response Formatting

**Do:** Use helper functions untuk consistent responses
```go
utils.SuccessResponse(c, 200, "Success", data)
utils.ErrorResponse(c, 400, "Error message")
```

---

## 📚 Related Documentation

- [Getting Started Guide](./getting-started.md) - Setup development environment
- [Middleware Guide](./middleware.md) - Authentication dan authorization
- [Database Models](../../05-guides/database/models.md) - GORM models reference
- [API Reference](../../04-api-reference/README.md) - Complete API documentation

---

**Last Updated:** 28 Desember 2025  
**Status:** ✅ Production Ready
