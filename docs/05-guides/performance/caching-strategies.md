# 💾 Caching Strategies - Sirine Go

Panduan strategi caching untuk meningkatkan performa aplikasi Sirine Go di semua layer (backend, frontend, database).

## 📋 Daftar Isi

1. [Application-Level Caching](#application-level-caching)
2. [HTTP Response Caching](#http-response-caching)
3. [Static Asset Caching](#static-asset-caching)
4. [Browser Caching](#browser-caching)
5. [Cache Invalidation](#cache-invalidation-strategies)

---

## 💻 Application-Level Caching

### **In-Memory Cache (Go)**

```go
// File: backend/utils/cache.go

package utils

import (
    "sync"
    "time"
)

type CacheItem struct {
    Value      interface{}
    Expiration time.Time
}

type Cache struct {
    items map[string]CacheItem
    mu    sync.RWMutex
}

func NewCache() *Cache {
    cache := &Cache{
        items: make(map[string]CacheItem),
    }
    
    // Clean expired items every 5 minutes
    go cache.cleanupExpired()
    
    return cache
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.items[key] = CacheItem{
        Value:      value,
        Expiration: time.Now().Add(duration),
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    item, found := c.items[key]
    if !found {
        return nil, false
    }
    
    // Check expiration
    if time.Now().After(item.Expiration) {
        return nil, false
    }
    
    return item.Value, true
}

func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}

func (c *Cache) cleanupExpired() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for key, item := range c.items {
            if now.After(item.Expiration) {
                delete(c.items, key)
            }
        }
        c.mu.Unlock()
    }
}
```

### **Usage Example**

```go
// File: backend/services/user_service.go

var userCache = utils.NewCache()

func GetUserByID(id int) (*models.User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    
    // Check cache first
    if cached, found := userCache.Get(cacheKey); found {
        return cached.(*models.User), nil
    }
    
    // Query database
    var user models.User
    err := db.First(&user, id).Error
    if err != nil {
        return nil, err
    }
    
    // Cache for 5 minutes
    userCache.Set(cacheKey, &user, 5*time.Minute)
    
    return &user, nil
}
```

### **Cache User Profile (Example)**

```go
// File: backend/handlers/profile_handler.go

func GetProfile(c *gin.Context) {
    userID := c.GetInt("user_id")
    cacheKey := fmt.Sprintf("profile:%d", userID)
    
    // Try cache
    if cached, found := cache.Get(cacheKey); found {
        c.JSON(200, gin.H{"data": cached})
        return
    }
    
    // Query from database
    var profile models.Profile
    if err := db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
        c.JSON(404, gin.H{"error": "Profile not found"})
        return
    }
    
    // Cache for 10 minutes
    cache.Set(cacheKey, profile, 10*time.Minute)
    
    c.JSON(200, gin.H{"data": profile})
}
```

---

## 🌐 HTTP Response Caching

### **Cache-Control Headers**

```go
// File: backend/middleware/cache_headers.go

func CacheHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        path := c.Request.URL.Path
        
        // No cache for API endpoints
        if strings.HasPrefix(path, "/api") {
            c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
            c.Header("Pragma", "no-cache")
            c.Header("Expires", "0")
        }
        
        // Cache static assets
        if strings.HasPrefix(path, "/assets") {
            c.Header("Cache-Control", "public, max-age=31536000, immutable")
        }
        
        c.Next()
    }
}
```

### **ETag Support**

```go
func GetUserHandler(c *gin.Context) {
    user := getUserFromDB()
    
    // Generate ETag from user data
    etag := fmt.Sprintf(`"%x"`, md5.Sum([]byte(user.UpdatedAt.String())))
    
    // Check if client has cached version
    if c.GetHeader("If-None-Match") == etag {
        c.Status(304) // Not Modified
        return
    }
    
    // Send response with ETag
    c.Header("ETag", etag)
    c.JSON(200, user)
}
```

---

## 🎨 Static Asset Caching

### **Nginx Configuration**

```nginx
# File: /etc/nginx/sites-available/sirine-go

server {
    # Frontend static files
    location / {
        root /var/www/sirine-go/frontend/dist;
        try_files $uri $uri/ /index.html;
        
        # Cache HTML files for short duration
        location ~* \.html$ {
            expires 1h;
            add_header Cache-Control "public, max-age=3600";
        }
        
        # Cache CSS/JS with hash (1 year)
        location ~* \.(js|css)$ {
            expires 1y;
            add_header Cache-Control "public, max-age=31536000, immutable";
        }
        
        # Cache images (1 week)
        location ~* \.(jpg|jpeg|png|gif|ico|svg|webp)$ {
            expires 7d;
            add_header Cache-Control "public, max-age=604800";
        }
        
        # Cache fonts (1 year)
        location ~* \.(woff|woff2|ttf|otf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, max-age=31536000, immutable";
        }
    }
    
    # API endpoints - no cache
    location /api {
        proxy_pass http://localhost:8080;
        add_header Cache-Control "no-cache, no-store, must-revalidate";
    }
}
```

### **Vite Build with Hashing**

```javascript
// File: frontend/vite.config.js

export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        // Add hash to filenames for cache busting
        entryFileNames: 'assets/[name].[hash].js',
        chunkFileNames: 'assets/[name].[hash].js',
        assetFileNames: 'assets/[name].[hash].[ext]'
      }
    }
  }
})
```

**Result:**
```
app.a3b2c1d4.js      ← Hash changes when content changes
app.a3b2c1d4.css
logo.5e6f7g8h.svg
```

---

## 🖥️ Browser Caching

### **Service Worker Cache**

```javascript
// File: frontend/public/sw.js

const CACHE_NAME = 'sirine-go-v1'
const CACHE_URLS = [
  '/',
  '/index.html',
  '/assets/app.js',
  '/assets/app.css',
  '/logo.svg'
]

// Install: Cache static assets
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(CACHE_URLS)
    })
  )
})

// Fetch: Cache-first strategy for static assets
self.addEventListener('fetch', (event) => {
  const { request } = event
  
  // Network-first for API calls
  if (request.url.includes('/api')) {
    event.respondWith(networkFirst(request))
    return
  }
  
  // Cache-first for static assets
  event.respondWith(cacheFirst(request))
})

async function cacheFirst(request) {
  const cached = await caches.match(request)
  return cached || fetch(request)
}

async function networkFirst(request) {
  try {
    const response = await fetch(request)
    return response
  } catch (error) {
    return caches.match(request)
  }
}
```

### **LocalStorage Caching (Vue)**

```javascript
// File: frontend/src/services/storageCache.js

export const storageCache = {
  set(key, value, expirationMinutes = 5) {
    const item = {
      value,
      expiration: Date.now() + (expirationMinutes * 60 * 1000)
    }
    localStorage.setItem(key, JSON.stringify(item))
  },
  
  get(key) {
    const itemStr = localStorage.getItem(key)
    if (!itemStr) return null
    
    const item = JSON.parse(itemStr)
    
    // Check expiration
    if (Date.now() > item.expiration) {
      localStorage.removeItem(key)
      return null
    }
    
    return item.value
  },
  
  remove(key) {
    localStorage.removeItem(key)
  },
  
  clear() {
    localStorage.clear()
  }
}
```

### **Usage in Vue Component**

```vue
<script setup>
import { ref, onMounted } from 'vue'
import { storageCache } from '@/services/storageCache'
import api from '@/services/api'

const userProfile = ref(null)

onMounted(async () => {
  // Try cache first
  const cached = storageCache.get('user-profile')
  if (cached) {
    userProfile.value = cached
    return
  }
  
  // Fetch from API
  const response = await api.get('/api/profile')
  userProfile.value = response.data
  
  // Cache for 10 minutes
  storageCache.set('user-profile', response.data, 10)
})
</script>
```

---

## 🔄 Cache Invalidation Strategies

### **Time-Based Expiration (TTL)**

```go
// Simple: Cache expires after fixed time
cache.Set("users:list", users, 5*time.Minute)
```

**Use when:**
- Data changes infrequently
- Stale data is acceptable for short periods
- Simple to implement

### **Event-Based Invalidation**

```go
// File: backend/services/user_service.go

func UpdateUser(id int, updates map[string]interface{}) error {
    // Update database
    err := db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
    if err != nil {
        return err
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("user:%d", id)
    cache.Delete(cacheKey)
    
    // Also invalidate user list
    cache.Delete("users:list")
    
    return nil
}
```

**Use when:**
- Data must be fresh immediately
- Updates are infrequent
- Cache invalidation logic is simple

### **Write-Through Cache**

```go
func UpdateUserProfile(userID int, profile models.Profile) error {
    // 1. Update database
    err := db.Save(&profile).Error
    if err != nil {
        return err
    }
    
    // 2. Update cache immediately
    cacheKey := fmt.Sprintf("profile:%d", userID)
    cache.Set(cacheKey, profile, 10*time.Minute)
    
    return nil
}
```

**Use when:**
- Read operations far exceed writes
- Consistency is important
- Cache misses are expensive

### **Cache Tagging**

```go
type CacheWithTags struct {
    cache *Cache
    tags  map[string][]string // tag -> keys
    mu    sync.RWMutex
}

func (c *CacheWithTags) SetWithTags(key string, value interface{}, duration time.Duration, tags []string) {
    c.cache.Set(key, value, duration)
    
    c.mu.Lock()
    defer c.mu.Unlock()
    
    for _, tag := range tags {
        c.tags[tag] = append(c.tags[tag], key)
    }
}

func (c *CacheWithTags) InvalidateTag(tag string) {
    c.mu.Lock()
    keys := c.tags[tag]
    delete(c.tags, tag)
    c.mu.Unlock()
    
    for _, key := range keys {
        c.cache.Delete(key)
    }
}

// Usage:
cache.SetWithTags("user:1", user, 5*time.Minute, []string{"users", "profile"})
cache.InvalidateTag("users") // Invalidates all user-related cache
```

---

## ✅ Caching Best Practices

### **General Guidelines**

1. **Cache Expensive Operations**
   - ✅ Database queries
   - ✅ External API calls
   - ✅ Heavy computations
   - ❌ Simple variable lookups

2. **Set Appropriate TTL**
   - User data: 5-10 minutes
   - Static content: 1 year
   - API responses: 1-5 minutes
   - Real-time data: Don't cache

3. **Cache Keys**
   - Use descriptive names: `user:123` not `u123`
   - Include versioning: `profile:v2:123`
   - Use consistent format

4. **Monitor Cache Performance**
   - Track hit/miss ratio
   - Monitor memory usage
   - Log cache operations
   - Measure impact on response time

### **Common Pitfalls**

- ❌ Caching everything (wastes memory)
- ❌ Too long TTL (stale data)
- ❌ Too short TTL (no benefit)
- ❌ No cache invalidation strategy
- ❌ Not monitoring cache performance

---

## 📊 Cache Performance Metrics

### **Target Metrics**

- **Cache Hit Ratio:** > 80%
- **Memory Usage:** < 500MB for cache
- **Cache Response Time:** < 1ms
- **Invalidation Time:** < 10ms

### **Monitoring Cache**

```go
type CacheStats struct {
    Hits     int64
    Misses   int64
    Sets     int64
    Deletes  int64
    Size     int
}

func (c *Cache) GetStats() CacheStats {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    return CacheStats{
        Hits:    c.hits,
        Misses:  c.misses,
        Sets:    c.sets,
        Deletes: c.deletes,
        Size:    len(c.items),
    }
}

// Expose via API
func CacheStatsHandler(c *gin.Context) {
    stats := cache.GetStats()
    hitRatio := float64(stats.Hits) / float64(stats.Hits+stats.Misses) * 100
    
    c.JSON(200, gin.H{
        "stats": stats,
        "hit_ratio": fmt.Sprintf("%.2f%%", hitRatio),
    })
}
```

---

## 📞 Support

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733

## 📖 Related Documentation

- [Backend Optimization](./backend-optimization.md)
- [Frontend Optimization](./frontend-optimization.md)
- [Database Optimization](./database-optimization.md)
- [Performance Guide](./README.md)

---

**Last Updated:** 29 Desember 2025  
**Version:** 1.0.0
