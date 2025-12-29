# ⚡ Backend Performance Optimization - Go/Gin

Panduan optimasi performa backend Go dengan Gin framework untuk aplikasi Sirine Go.

## 📋 Daftar Isi

1. [Go Runtime Optimization](#go-runtime-optimization)
2. [Database Connection Pooling](#database-connection-pooling)
3. [Middleware Optimization](#middleware-optimization)
4. [Request Handling](#request-handling-optimization)
5. [Memory Management](#memory-management)
6. [Concurrency Patterns](#concurrency-patterns)
7. [Profiling & Debugging](#profiling--debugging)

---

## 🔧 Go Runtime Optimization

### **Environment Variables**

**Production Configuration:**
```bash
# Enable release mode
export GIN_MODE=release

# Set GOMAXPROCS (default: number of CPUs)
export GOMAXPROCS=4

# Garbage collector tuning
export GOGC=100  # Default, adjust based on memory usage
```

### **Build Optimization**

```bash
# Build with optimizations
go build -ldflags="-s -w" -o sirine-go cmd/server/main.go

# Flags explanation:
# -s : Omit symbol table
# -w : Omit DWARF debug info
# Result: ~30-40% smaller binary
```

### **Runtime Configuration**

```go
import "runtime"

func init() {
    // Set max CPUs (optional, defaults to all available)
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    // Tune GC (optional, experiment with values)
    debug.SetGCPercent(100) // Default
}
```

---

## 🗄️ Database Connection Pooling

### **Optimal Configuration**

```go
// File: backend/database/database.go

func InitDatabase() (*gorm.DB, error) {
    // ... connection setup ...
    
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    // Connection pool settings
    sqlDB.SetMaxOpenConns(25)      // Maximum open connections
    sqlDB.SetMaxIdleConns(10)      // Maximum idle connections
    sqlDB.SetConnMaxLifetime(time.Hour)     // Connection lifetime
    sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Idle timeout
    
    return db, nil
}
```

### **Connection Pool Guidelines**

**For different server specs:**

```go
// Small Server (1-2 CPU, 1-2GB RAM)
sqlDB.SetMaxOpenConns(10)
sqlDB.SetMaxIdleConns(5)

// Medium Server (2-4 CPU, 4GB RAM)
sqlDB.SetMaxOpenConns(25)  // ✅ Recommended untuk Sirine Go
sqlDB.SetMaxIdleConns(10)

// Large Server (4+ CPU, 8GB+ RAM)
sqlDB.SetMaxOpenConns(50)
sqlDB.SetMaxIdleConns(20)
```

### **Monitoring Connection Pool**

```go
func GetDBStats(c *gin.Context) {
    sqlDB, _ := database.DB.DB()
    stats := sqlDB.Stats()
    
    c.JSON(200, gin.H{
        "max_open_connections": stats.MaxOpenConnections,
        "open_connections":     stats.OpenConnections,
        "in_use":              stats.InUse,
        "idle":                stats.Idle,
        "wait_count":          stats.WaitCount,
        "wait_duration":       stats.WaitDuration.String(),
    })
}
```

---

## ⚙️ Middleware Optimization

### **Middleware Order (Penting!)**

```go
// File: backend/routes/routes.go

func SetupRouter(db *gorm.DB) *gin.Engine {
    router := gin.New()
    
    // 1. Recovery first (panic handler)
    router.Use(gin.Recovery())
    
    // 2. CORS (filter early)
    router.Use(middleware.CORSMiddleware())
    
    // 3. Gzip compression (reduce bandwidth)
    router.Use(gzip.Gzip(gzip.DefaultCompression))
    
    // 4. Logger (after compression for accurate size)
    router.Use(gin.Logger())
    
    // 5. Auth middleware (only on protected routes)
    // Apply per route group, NOT globally
    
    return router
}
```

### **Conditional Middleware**

```go
// ❌ BAD: Apply logger to ALL routes
router.Use(gin.Logger())

// ✅ GOOD: Skip logger for health checks
router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
    SkipPaths: []string{"/health", "/metrics"},
}))
```

### **Custom Activity Logger Optimization**

```go
// File: backend/middleware/activity_logger.go

func ActivityLogger(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() // Process request first
        
        // Log asynchronously (non-blocking)
        go func() {
            log := models.ActivityLog{
                UserID:    GetUserID(c),
                Action:    c.Request.Method + " " + c.Request.URL.Path,
                Timestamp: time.Now(),
            }
            db.Create(&log) // Fire and forget
        }()
    }
}
```

---

## 🎯 Request Handling Optimization

### **Reduce JSON Serialization Overhead**

```go
// ❌ SLOW: Create new struct every time
func GetUsers(c *gin.Context) {
    var users []models.User
    db.Find(&users)
    
    // This creates intermediate structs
    c.JSON(200, gin.H{
        "status": "success",
        "data": users,
    })
}

// ✅ FAST: Reuse response struct
type UserResponse struct {
    Status string        `json:"status"`
    Data   []models.User `json:"data"`
}

var userResponsePool = sync.Pool{
    New: func() interface{} {
        return &UserResponse{}
    },
}

func GetUsers(c *gin.Context) {
    var users []models.User
    db.Find(&users)
    
    resp := userResponsePool.Get().(*UserResponse)
    defer userResponsePool.Put(resp)
    
    resp.Status = "success"
    resp.Data = users
    
    c.JSON(200, resp)
}
```

### **Pagination for Large Datasets**

```go
// File: backend/handlers/user_handler.go

func ListUsers(c *gin.Context) {
    // Parse pagination params
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    
    // Limit max results (prevent abuse)
    if limit > 100 {
        limit = 100
    }
    
    offset := (page - 1) * limit
    
    var users []models.User
    var total int64
    
    // Count total (cache this if possible)
    db.Model(&models.User{}).Count(&total)
    
    // Get page
    db.Limit(limit).Offset(offset).Find(&users)
    
    c.JSON(200, gin.H{
        "data": users,
        "pagination": gin.H{
            "page":  page,
            "limit": limit,
            "total": total,
        },
    })
}
```

### **Avoid N+1 Queries**

```go
// ❌ BAD: N+1 queries
func GetUsersWithAchievements(c *gin.Context) {
    var users []models.User
    db.Find(&users)
    
    // This triggers 1 query per user!
    for i := range users {
        db.Model(&users[i]).Association("Achievements").Find(&users[i].Achievements)
    }
    
    c.JSON(200, users)
}

// ✅ GOOD: Single query with Preload
func GetUsersWithAchievements(c *gin.Context) {
    var users []models.User
    db.Preload("Achievements").Find(&users)
    
    c.JSON(200, users)
}
```

---

## 💾 Memory Management

### **Proper Resource Cleanup**

```go
func UploadFile(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": "File required"})
        return
    }
    
    // Open file
    src, err := file.Open()
    if err != nil {
        c.JSON(500, gin.H{"error": "Cannot open file"})
        return
    }
    defer src.Close() // ✅ Always close resources
    
    // Process file...
}
```

### **Limit Request Body Size**

```go
// File: cmd/server/main.go

func main() {
    router := gin.Default()
    
    // Limit request body (prevent memory exhaustion)
    router.MaxMultipartMemory = 8 << 20  // 8 MiB
    
    // For JSON
    router.Use(func(c *gin.Context) {
        c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // 10MB
        c.Next()
    })
}
```

### **Avoid Memory Leaks**

```go
// ❌ BAD: Goroutine leak
func NotifyUsers(c *gin.Context) {
    users := getUserList()
    
    for _, user := range users {
        go sendNotification(user) // No way to cancel!
    }
}

// ✅ GOOD: Context-aware goroutines
func NotifyUsers(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
    defer cancel()
    
    users := getUserList()
    
    for _, user := range users {
        go func(u User) {
            select {
            case <-ctx.Done():
                return // Clean exit
            default:
                sendNotification(u)
            }
        }(user)
    }
}
```

---

## 🔄 Concurrency Patterns

### **Worker Pool for Batch Processing**

```go
func ProcessNotifications(notifications []Notification) {
    jobs := make(chan Notification, len(notifications))
    results := make(chan bool, len(notifications))
    
    // Create worker pool (4 workers)
    for w := 1; w <= 4; w++ {
        go worker(jobs, results)
    }
    
    // Send jobs
    for _, n := range notifications {
        jobs <- n
    }
    close(jobs)
    
    // Wait for results
    for i := 0; i < len(notifications); i++ {
        <-results
    }
}

func worker(jobs <-chan Notification, results chan<- bool) {
    for n := range jobs {
        sendNotification(n)
        results <- true
    }
}
```

### **Rate Limiting**

```go
import "golang.org/x/time/rate"

var limiter = rate.NewLimiter(100, 200) // 100 req/sec, burst 200

func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

---

## 🔍 Profiling & Debugging

### **Enable pprof (Development Only)**

```go
import _ "net/http/pprof"

func main() {
    if gin.Mode() == gin.DebugMode {
        go func() {
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
    }
    
    // ... rest of app ...
}
```

### **CPU Profiling**

```bash
# Start profiling
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Commands in pprof:
# top10     - Top 10 functions by CPU
# list main - Show source code
# web       - Open graph in browser
```

### **Memory Profiling**

```bash
# Heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Allocation profile
go tool pprof http://localhost:6060/debug/pprof/allocs
```

### **Benchmark Tests**

```go
// File: backend/handlers/user_handler_test.go

func BenchmarkGetUsers(b *testing.B) {
    router := setupTestRouter()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/api/users", nil)
        router.ServeHTTP(w, req)
    }
}
```

Run:
```bash
go test -bench=. -benchmem ./handlers/
```

---

## ✅ Performance Checklist

### Production Deployment
- [ ] `GIN_MODE=release` enabled
- [ ] Build with `-ldflags="-s -w"`
- [ ] Database connection pool configured
- [ ] Gzip compression enabled
- [ ] Request body size limits set
- [ ] Activity logger is asynchronous
- [ ] N+1 queries eliminated
- [ ] Pagination implemented
- [ ] Resource cleanup with `defer`
- [ ] Rate limiting configured

### Monitoring
- [ ] Log slow database queries (>100ms)
- [ ] Monitor connection pool stats
- [ ] Track memory usage
- [ ] Monitor goroutine count
- [ ] Set up alerts for high response times

---

## 📊 Expected Performance Metrics

### Response Times (Production)
- Health check: < 5ms
- Simple query (by ID): < 20ms
- List with pagination: < 50ms
- Complex query (with joins): < 100ms
- File upload: < 500ms (depends on size)

### Resource Usage
- Memory (idle): 50-100MB
- Memory (under load): 150-300MB
- CPU (idle): 0-5%
- CPU (under load): 30-60%
- Goroutines: < 100

---

## 📞 Support

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733

## 📖 Related Documentation

- [Backend Development Guide](../../03-development/backend/getting-started.md)
- [Service Pattern](../../03-development/backend/service-pattern.md)
- [Database Optimization](./database-optimization.md)
- [Middleware Guide](../../03-development/backend/middleware.md)

---

**Last Updated:** 29 Desember 2025  
**Version:** 1.0.0
