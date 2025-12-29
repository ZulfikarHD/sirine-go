# 🔧 Backend Development Guide

Panduan comprehensive untuk backend development dalam Sirine Go App menggunakan Go, Gin framework, dan MySQL.

**Tech Stack:**
- Go 1.21+
- Gin (Web framework)
- GORM (ORM)
- MySQL 8.0+
- JWT (Authentication)
- Bcrypt (Password hashing)

**Last Updated:** 28 Desember 2025

---

## 📋 Overview

Backend Sirine Go App dibangun dengan Go dan Gin framework untuk fast, reliable, dan scalable API server, dengan focus pada:
- RESTful API design
- Service pattern architecture
- JWT-based authentication
- Comprehensive error handling
- Activity logging untuk audit trail

---

## 🚀 Quick Start

### Prerequisites

**Required:**
- Go 1.21+ (`go version`)
- MySQL 8.0+ running
- Air (hot reload tool)
- Basic knowledge Go dan Gin framework

### Setup Development Environment

```bash
# Navigate to backend folder
cd backend

# Install Air untuk hot reload (optional)
go install github.com/air-verse/air@latest

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Edit .env dengan database credentials
vim .env

# Run migrations
go run cmd/migrate/main.go

# Seed database
go run cmd/seed/main.go

# Start server
air
# Or without hot reload:
go run cmd/server/main.go
```

### Development Commands

```bash
# Run server (with hot reload)
air

# Run server (without hot reload)
go run cmd/server/main.go

# Run tests
go test ./... -v

# Run specific test
go test ./internal/services -v

# Build binary
go build -o bin/server cmd/server/main.go

# Run migrations
go run cmd/migrate/main.go

# Seed database
go run cmd/seed/main.go
```

---

## 📁 Project Structure

```
backend/
├── cmd/
│   ├── server/          # Main server entry point
│   │   └── main.go
│   ├── migrate/         # Database migrations
│   │   └── main.go
│   └── seed/            # Database seeding
│       └── main.go
├── internal/
│   ├── config/          # Configuration management
│   │   └── config.go
│   ├── database/        # Database connection & setup
│   │   └── database.go
│   ├── middleware/      # HTTP middlewares
│   │   ├── auth.go      # JWT authentication
│   │   ├── cors.go      # CORS handling
│   │   └── logger.go    # Request logging
│   ├── models/          # Database models (GORM)
│   │   ├── user.go
│   │   ├── achievement.go
│   │   └── notification.go
│   ├── handlers/        # HTTP handlers (controllers)
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   └── profile_handler.go
│   ├── services/        # Business logic layer
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   └── profile_service.go
│   ├── repositories/    # Data access layer (optional)
│   │   └── user_repository.go
│   ├── utils/           # Utility functions
│   │   ├── jwt.go       # JWT utilities
│   │   ├── hash.go      # Password hashing
│   │   ├── response.go  # Response helpers
│   │   └── validator.go # Input validation
│   └── routes/          # Route definitions
│       └── routes.go
├── migrations/          # SQL migration files
│   ├── 001_create_users_table.sql
│   └── 002_create_achievements_table.sql
├── .env                 # Environment variables
├── .env.example         # Environment template
├── go.mod               # Go dependencies
└── go.sum               # Dependencies checksums
```

---

## 🏗️ Architecture Pattern

### Service Pattern

Sirine Go menggunakan **Service Pattern** untuk separation of concerns:

```
Request → Handler → Service → Model/DB → Response
```

**Layers:**
1. **Handler (Controller):** Handle HTTP requests/responses
2. **Service:** Business logic & validation
3. **Model/Repository:** Data access & persistence

### Example Flow

```go
// 1. Handler receives request
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    c.ShouldBindJSON(&req)
    
    // 2. Pass to service
    user, err := h.service.CreateUser(&req)
    if err != nil {
        utils.ErrorResponse(c, 500, err.Error())
        return
    }
    
    // 3. Return response
    utils.SuccessResponse(c, 201, "User created", user)
}

// 4. Service handles business logic
func (s *UserService) CreateUser(req *CreateUserRequest) (*models.User, error) {
    // Validate
    if err := s.validateUser(req); err != nil {
        return nil, err
    }
    
    // Hash password
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
    
    // Create user
    user := &models.User{
        NIP:      req.NIP,
        Password: string(hashedPassword),
        // ...
    }
    
    // 5. Save to database
    if err := s.db.Create(user).Error; err != nil {
        return nil, err
    }
    
    return user, nil
}
```

---

## 🔐 Authentication & Authorization

### JWT Implementation

```go
// internal/utils/jwt.go
package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string) (string, error) {
    claims := JWTClaims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}
```

### Auth Middleware

```go
// internal/middleware/auth.go
package middleware

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token from header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"success": false, "error": "Token tidak ditemukan"})
            c.Abort()
            return
        }
        
        // Validate token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
            c.JSON(401, gin.H{"success": false, "error": "Token tidak valid"})
            c.Abort()
            return
        }
        
        // Set user info to context
        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}

// Role-based authorization
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("role")
        
        for _, role := range roles {
            if userRole == role {
                c.Next()
                return
            }
        }
        
        c.JSON(403, gin.H{"success": false, "error": "Akses ditolak"})
        c.Abort()
    }
}
```

---

## 🗄️ Database & Models

### Model Definition (GORM)

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
    Password     string         `gorm:"not null" json:"-"` // Hidden from JSON
    Role         string         `gorm:"not null" json:"role"`
    Department   string         `gorm:"not null" json:"department"`
    Shift        string         `gorm:"not null" json:"shift"`
    Status       string         `gorm:"default:ACTIVE" json:"status"`
    TotalPoints  int            `gorm:"default:0" json:"total_points"`
    Level        string         `gorm:"default:Bronze" json:"level"`
    ProfilePhoto string         `json:"profile_photo"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
    
    // Relationships
    Achievements []Achievement `gorm:"many2many:user_achievements;" json:"achievements,omitempty"`
    Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
}

// Table name override
func (User) TableName() string {
    return "users"
}

// Hooks
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Auto-calculate level based on points
    u.Level = calculateLevel(u.TotalPoints)
    return nil
}

func calculateLevel(points int) string {
    if points >= 1000 {
        return "Platinum"
    } else if points >= 500 {
        return "Gold"
    } else if points >= 100 {
        return "Silver"
    }
    return "Bronze"
}
```

### Database Queries

```go
// Simple query
var user models.User
db.First(&user, 1) // Find by ID

// Query with conditions
db.Where("role = ?", "ADMIN").Find(&users)

// Query with multiple conditions
db.Where("role = ? AND status = ?", "ADMIN", "ACTIVE").Find(&users)

// Pagination
db.Offset(offset).Limit(limit).Find(&users)

// Preload relationships
db.Preload("Achievements").Find(&users)

// Count
var count int64
db.Model(&models.User{}).Count(&count)

// Create
user := models.User{NIP: "12345", FullName: "John Doe"}
db.Create(&user)

// Update
db.Model(&user).Updates(map[string]interface{}{
    "full_name": "Jane Doe",
    "email": "jane@example.com",
})

// Soft delete
db.Delete(&user, 1)

// Permanent delete
db.Unscoped().Delete(&user, 1)
```

---

## 🛠️ Common Patterns

### Service Implementation

```go
// internal/services/user_service.go
package services

type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

// List users with pagination
func (s *UserService) ListUsers(page, pageSize int, filters map[string]interface{}) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    
    query := s.db.Model(&models.User{})
    
    // Apply filters
    if role, ok := filters["role"].(string); ok && role != "" {
        query = query.Where("role = ?", role)
    }
    
    if department, ok := filters["department"].(string); ok && department != "" {
        query = query.Where("department = ?", department)
    }
    
    // Count total
    query.Count(&total)
    
    // Pagination
    offset := (page - 1) * pageSize
    err := query.Offset(offset).Limit(pageSize).Find(&users).Error
    
    return users, total, err
}

// Create user
func (s *UserService) CreateUser(req *CreateUserRequest) (*models.User, error) {
    // Check if NIP exists
    var existingUser models.User
    if err := s.db.Where("nip = ?", req.NIP).First(&existingUser).Error; err == nil {
        return nil, errors.New("NIP sudah terdaftar")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
    if err != nil {
        return nil, err
    }
    
    // Create user
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
    
    if err := s.db.Create(user).Error; err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### Handler Implementation

```go
// internal/handlers/user_handler.go
package handlers

type UserHandler struct {
    service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
    return &UserHandler{service: service}
}

// GET /api/users
func (h *UserHandler) ListUsers(c *gin.Context) {
    // Parse query params
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    
    filters := map[string]interface{}{
        "role":       c.Query("role"),
        "department": c.Query("department"),
        "status":     c.Query("status"),
    }
    
    // Call service
    users, total, err := h.service.ListUsers(page, pageSize, filters)
    if err != nil {
        utils.ErrorResponse(c, 500, "Gagal mengambil data users")
        return
    }
    
    // Return response
    utils.PaginatedResponse(c, 200, users, total, page, pageSize)
}

// POST /api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ValidationError(c, err)
        return
    }
    
    user, err := h.service.CreateUser(&req)
    if err != nil {
        utils.ErrorResponse(c, 400, err.Error())
        return
    }
    
    utils.SuccessResponse(c, 201, "User berhasil dibuat", user)
}
```

### Response Helpers

```go
// internal/utils/response.go
package utils

func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
    c.JSON(status, gin.H{
        "success": true,
        "message": message,
        "data":    data,
    })
}

func ErrorResponse(c *gin.Context, status int, message string) {
    c.JSON(status, gin.H{
        "success": false,
        "error":   message,
    })
}

func ValidationError(c *gin.Context, err error) {
    c.JSON(400, gin.H{
        "success": false,
        "error":   "Validasi gagal",
        "details": err.Error(),
    })
}

func PaginatedResponse(c *gin.Context, status int, data interface{}, total int64, page, pageSize int) {
    c.JSON(status, gin.H{
        "success": true,
        "data": gin.H{
            "items":      data,
            "total":      total,
            "page":       page,
            "page_size":  pageSize,
            "total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
        },
    })
}
```

---

## 🧪 Testing

### Unit Test Example

```go
// internal/services/user_service_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    service := NewUserService(db)
    
    // Test case
    req := &CreateUserRequest{
        NIP:        "12345",
        FullName:   "Test User",
        Email:      "test@example.com",
        Phone:      "081234567890",
        Password:   "Test@123",
        Role:       "STAFF_KHAZWAL",
        Department: "KHAZWAL",
        Shift:      "PAGI",
    }
    
    user, err := service.CreateUser(req)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "12345", user.NIP)
    assert.NotEmpty(t, user.Password) // Password should be hashed
    assert.NotEqual(t, "Test@123", user.Password) // Should not be plain text
}
```

---

## ✅ Best Practices

### 1. Error Handling

```go
// Always check errors
if err != nil {
    log.Printf("Error: %v", err)
    return nil, err
}

// Use descriptive error messages
if user == nil {
    return nil, errors.New("user tidak ditemukan")
}

// Don't expose sensitive info in errors
// Bad: return errors.New(fmt.Sprintf("Database error: %v", dbErr))
// Good: return errors.New("terjadi kesalahan pada server")
```

### 2. Security

```go
// Always hash passwords
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

// Use parameterized queries (GORM handles this)
db.Where("nip = ?", nip).First(&user) // Safe from SQL injection

// Validate input
if len(req.Password) < 8 {
    return errors.New("password minimal 8 karakter")
}

// Hide sensitive fields from JSON
type User struct {
    Password string `json:"-"` // Won't be included in JSON response
}
```

### 3. Database Transactions

```go
func (s *Service) CreateUserWithAchievement(userData, achievementData) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Create user
        if err := tx.Create(&user).Error; err != nil {
            return err // Rollback
        }
        
        // Award achievement
        if err := tx.Create(&userAchievement).Error; err != nil {
            return err // Rollback
        }
        
        return nil // Commit
    })
}
```

### 4. Logging

```go
import "log"

// Log important events
log.Printf("User %s logged in successfully", user.NIP)

// Log errors for debugging
log.Printf("ERROR: Failed to create user: %v", err)

// Don't log sensitive data
// Bad: log.Printf("Password: %s", password)
// Good: log.Printf("Password validation failed for user %s", userNIP)
```

---

## 📚 Additional Resources

### Documentation
- [Go Docs](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [JWT Go](https://github.com/golang-jwt/jwt)

### Related Guides
- [Frontend Guide](./frontend-guide.md)
- [API Documentation](../04-api-reference/README.md)
- [Database Models](../05-guides/database/models.md)
- [Configuration Guide](../05-guides/configuration.md)

---

**Last Updated:** 28 Desember 2025  
**Status:** ✅ Production Ready
