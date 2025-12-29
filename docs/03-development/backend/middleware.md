# 🛡️ Middleware Development

Panduan untuk membuat dan menggunakan middleware dalam Sirine Go App, yaitu functions yang mengintercept HTTP requests untuk authentication, logging, CORS, dan custom processing.

---

## 📋 Overview

Middleware merupakan functions yang execute sebelum atau sesudah HTTP handler, yang bertujuan untuk:
- Authentication dan authorization
- Request/response logging
- CORS handling
- Rate limiting
- Request validation
- Error recovery

### Middleware Flow

```
HTTP Request
    ↓
[CORS Middleware]
    ↓
[Logger Middleware]
    ↓
[Auth Middleware]
    ↓
[Role Authorization Middleware]
    ↓
Handler
    ↓
Response
```

---

## 🔐 Authentication Middleware

### JWT Token Validation

```go
// internal/middleware/auth.go
package middleware

import (
    "strings"
    "net/http"
    "sirine-go/internal/utils"
    "github.com/gin-gonic/gin"
)

// AuthMiddleware untuk validasi JWT token pada protected routes
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token dari Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "Token tidak ditemukan",
            })
            c.Abort()  // Stop request processing
            return
        }
        
        // Token format: "Bearer <token>"
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            // Tidak ada "Bearer " prefix
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "Format token tidak valid",
            })
            c.Abort()
            return
        }
        
        // Validate token dengan JWT utility
        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "Token tidak valid atau expired",
            })
            c.Abort()
            return
        }
        
        // Set user info ke context untuk diakses oleh handler
        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Set("nip", claims.NIP)
        
        // Continue ke next handler
        c.Next()
    }
}

// GetUserID untuk mendapatkan user ID dari context
func GetUserID(c *gin.Context) (uint, bool) {
    userID, exists := c.Get("user_id")
    if !exists {
        return 0, false
    }
    
    id, ok := userID.(uint)
    return id, ok
}

// GetUserRole untuk mendapatkan user role dari context
func GetUserRole(c *gin.Context) (string, bool) {
    role, exists := c.Get("role")
    if !exists {
        return "", false
    }
    
    roleStr, ok := role.(string)
    return roleStr, ok
}
```

### JWT Utility Functions

```go
// internal/utils/jwt.go
package utils

import (
    "os"
    "time"
    "errors"
    "github.com/golang-jwt/jwt/v5"
)

// JWTClaims untuk custom claims dalam JWT token
type JWTClaims struct {
    UserID uint   `json:"user_id"`
    NIP    string `json:"nip"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// GenerateToken untuk membuat access token dengan expiry 15 menit
func GenerateToken(userID uint, nip, role string) (string, error) {
    claims := JWTClaims{
        UserID: userID,
        NIP:    nip,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "sirine-go",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// GenerateRefreshToken untuk membuat refresh token dengan expiry 30 hari
func GenerateRefreshToken(userID uint) (string, error) {
    claims := JWTClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(720 * time.Hour)), // 30 days
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "sirine-go",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// ValidateToken untuk validasi JWT token dan extract claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    // Extract claims
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        // Check jika token expired
        if claims.ExpiresAt.Before(time.Now()) {
            return nil, errors.New("token expired")
        }
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}
```

---

## 🔑 Role-Based Authorization

### Role Authorization Middleware

```go
// internal/middleware/auth.go (continued)

// RequireRole untuk membatasi akses berdasarkan role
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get user role dari context (set by AuthMiddleware)
        userRole, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "User role tidak ditemukan",
            })
            c.Abort()
            return
        }
        
        roleStr, ok := userRole.(string)
        if !ok {
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "error":   "Invalid role type",
            })
            c.Abort()
            return
        }
        
        // Check jika user role ada dalam allowed roles
        authorized := false
        for _, allowedRole := range allowedRoles {
            if roleStr == allowedRole {
                authorized = true
                break
            }
        }
        
        if !authorized {
            c.JSON(http.StatusForbidden, gin.H{
                "success": false,
                "error":   "Anda tidak memiliki akses untuk resource ini",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// RequireAdminRole untuk shortcut admin-only routes
func RequireAdminRole() gin.HandlerFunc {
    return RequireRole("SUPER_ADMIN", "ADMIN_YAKES")
}
```

### Usage in Routes

```go
// internal/routes/routes.go
package routes

import (
    "sirine-go/internal/handlers"
    "sirine-go/internal/middleware"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    api := r.Group("/api")
    
    // Public routes (tidak perlu authentication)
    {
        api.POST("/auth/login", authHandler.Login)
        api.POST("/auth/refresh", authHandler.RefreshToken)
    }
    
    // Protected routes (perlu authentication)
    protected := api.Group("")
    protected.Use(middleware.AuthMiddleware())
    {
        // All authenticated users
        protected.GET("/profile", profileHandler.GetProfile)
        protected.PUT("/profile", profileHandler.UpdateProfile)
        
        // Admin-only routes
        admin := protected.Group("")
        admin.Use(middleware.RequireRole("SUPER_ADMIN", "ADMIN_YAKES"))
        {
            admin.POST("/users", userHandler.CreateUser)
            admin.GET("/users", userHandler.ListUsers)
            admin.DELETE("/users/:id", userHandler.DeleteUser)
        }
        
        // Staff-only routes
        staff := protected.Group("")
        staff.Use(middleware.RequireRole("STAFF_KHAZWAL", "STAFF_FISIOTERAPI"))
        {
            staff.GET("/my-activities", activityHandler.GetMyActivities)
        }
    }
}
```

---

## 📝 Logger Middleware

### Request/Response Logging

```go
// internal/middleware/logger.go
package middleware

import (
    "time"
    "log"
    "github.com/gin-gonic/gin"
)

// LoggerMiddleware untuk logging setiap HTTP request
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        startTime := time.Now()
        
        // Get request info
        method := c.Request.Method
        path := c.Request.URL.Path
        clientIP := c.ClientIP()
        
        // Process request
        c.Next()
        
        // Calculate latency
        latency := time.Since(startTime)
        statusCode := c.Writer.Status()
        
        // Log format: [timestamp] method path status latency ip
        log.Printf("[%s] %s %s %d %v %s",
            time.Now().Format("2006-01-02 15:04:05"),
            method,
            path,
            statusCode,
            latency,
            clientIP,
        )
        
        // Log errors jika ada
        if len(c.Errors) > 0 {
            log.Printf("Errors: %v", c.Errors.String())
        }
    }
}

// ActivityLogMiddleware untuk logging user activities ke database
func ActivityLogMiddleware(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip logging untuk non-modifying requests
        if c.Request.Method == "GET" {
            c.Next()
            return
        }
        
        // Get user info dari context
        userID, exists := c.Get("user_id")
        if !exists {
            c.Next()
            return
        }
        
        // Process request
        c.Next()
        
        // Create activity log jika request success
        if c.Writer.Status() < 400 {
            activityLog := models.ActivityLog{
                UserID:     userID.(uint),
                Action:     c.Request.Method,
                Resource:   c.Request.URL.Path,
                IPAddress:  c.ClientIP(),
                UserAgent:  c.Request.UserAgent(),
                StatusCode: c.Writer.Status(),
            }
            
            // Save ke database (non-blocking)
            go func() {
                if err := db.Create(&activityLog).Error; err != nil {
                    log.Printf("Failed to save activity log: %v", err)
                }
            }()
        }
    }
}
```

---

## 🌐 CORS Middleware

### Cross-Origin Resource Sharing

```go
// internal/middleware/cors.go
package middleware

import (
    "os"
    "github.com/gin-gonic/gin"
)

// CORSMiddleware untuk handle cross-origin requests dari frontend
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        frontendURL := os.Getenv("FRONTEND_URL")
        if frontendURL == "" {
            frontendURL = "http://localhost:5173"
        }
        
        // Set CORS headers
        c.Writer.Header().Set("Access-Control-Allow-Origin", frontendURL)
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
        
        // Handle preflight requests
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

---

## 🚨 Error Recovery Middleware

### Panic Recovery

```go
// internal/middleware/recovery.go
package middleware

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
)

// RecoveryMiddleware untuk handle panic dan prevent server crash
func RecoveryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // Log error dengan stack trace
                log.Printf("PANIC: %v", err)
                log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
                
                // Return error response
                c.JSON(http.StatusInternalServerError, gin.H{
                    "success": false,
                    "error":   "Terjadi kesalahan pada server",
                })
                
                c.Abort()
            }
        }()
        
        c.Next()
    }
}
```

---

## ⚡ Rate Limiting Middleware

### Request Rate Limiter

```go
// internal/middleware/rate_limiter.go
package middleware

import (
    "sync"
    "time"
    "net/http"
    "github.com/gin-gonic/gin"
)

// RateLimiter untuk limit request per IP address
type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.Mutex
    limit    int           // Max requests
    window   time.Duration // Time window
}

// NewRateLimiter untuk create new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
}

// RateLimitMiddleware untuk apply rate limiting
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        
        rl.mu.Lock()
        defer rl.mu.Unlock()
        
        now := time.Now()
        
        // Clean old requests
        if timestamps, exists := rl.requests[ip]; exists {
            var validRequests []time.Time
            for _, t := range timestamps {
                if now.Sub(t) < rl.window {
                    validRequests = append(validRequests, t)
                }
            }
            rl.requests[ip] = validRequests
        }
        
        // Check rate limit
        if len(rl.requests[ip]) >= rl.limit {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "success": false,
                "error":   "Terlalu banyak request. Silakan coba lagi nanti.",
            })
            c.Abort()
            return
        }
        
        // Add current request
        rl.requests[ip] = append(rl.requests[ip], now)
        
        c.Next()
    }
}
```

**Usage:**
```go
// Limit: 100 requests per minute
limiter := middleware.NewRateLimiter(100, time.Minute)
r.Use(limiter.RateLimitMiddleware())
```

---

## 🔧 Custom Middleware Example

### Content-Type Validation

```go
// internal/middleware/validate_content_type.go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
)

// ValidateJSONMiddleware untuk ensure request content-type adalah JSON
func ValidateJSONMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip untuk GET, DELETE, OPTIONS
        if c.Request.Method == "GET" || c.Request.Method == "DELETE" || c.Request.Method == "OPTIONS" {
            c.Next()
            return
        }
        
        contentType := c.GetHeader("Content-Type")
        
        // Check jika Content-Type adalah application/json
        if !strings.Contains(contentType, "application/json") {
            c.JSON(http.StatusBadRequest, gin.H{
                "success": false,
                "error":   "Content-Type harus application/json",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

---

## 📦 Middleware Registration

### Main Setup

```go
// cmd/server/main.go
package main

import (
    "time"
    "sirine-go/internal/middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.New()  // Create router tanpa default middleware
    
    // Global middleware (applied ke semua routes)
    r.Use(middleware.RecoveryMiddleware())      // Panic recovery
    r.Use(middleware.CORSMiddleware())          // CORS handling
    r.Use(middleware.LoggerMiddleware())        // Request logging
    
    // Rate limiting (100 req/min)
    limiter := middleware.NewRateLimiter(100, time.Minute)
    r.Use(limiter.RateLimitMiddleware())
    
    // Setup routes dengan route-specific middleware
    setupRoutes(r)
    
    r.Run(":8080")
}
```

### Middleware Order

**Important:** Middleware execute dalam order yang didefinisikan

```
Recovery     (1st - catch all panics)
    ↓
CORS         (2nd - handle preflight)
    ↓
Logger       (3rd - log requests)
    ↓
Rate Limiter (4th - prevent abuse)
    ↓
Auth         (5th - route-specific)
    ↓
Role Check   (6th - route-specific)
    ↓
Handler      (last - business logic)
```

---

## ✅ Best Practices

### 1. Use c.Abort()

```go
// Do: Stop processing dengan Abort
if unauthorized {
    c.JSON(401, gin.H{"error": "Unauthorized"})
    c.Abort()
    return
}

// Don't: Lupa Abort bisa cause handler tetap execute
```

### 2. Set Context Values Safely

```go
// Do: Type assertion dengan check
if userID, ok := c.Get("user_id"); ok {
    id := userID.(uint)
}

// Don't: Assume value exists
id := c.MustGet("user_id").(uint)  // Panic jika tidak exists
```

### 3. Error Logging

```go
// Do: Log errors untuk debugging
if err != nil {
    log.Printf("ERROR: %v", err)
    c.JSON(500, gin.H{"error": "Internal error"})
}
```

### 4. Non-Blocking Operations

```go
// Do: Background tasks dengan goroutine
go func() {
    saveActivityLog(...)
}()

// Don't: Block request processing
saveActivityLog(...)  // Slow operation
```

---

## 📚 Related Documentation

- [Getting Started Guide](./getting-started.md) - Setup development environment
- [Service Pattern Guide](./service-pattern.md) - Architecture overview
- [Security Guide](../../05-guides/security/overview.md) - Security best practices
- [API Reference](../../04-api-reference/README.md) - API endpoints

---

**Last Updated:** 28 Desember 2025  
**Status:** ✅ Production Ready
