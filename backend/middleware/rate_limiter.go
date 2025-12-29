package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter merupakan middleware untuk rate limiting requests
// yang mencegah abuse dan brute force attacks
type RateLimiter struct {
	requests map[string]*requestInfo
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// requestInfo merupakan informasi tracking untuk each IP/user
type requestInfo struct {
	count      int
	firstRequest time.Time
	blocked    bool
	blockedUntil time.Time
}

// NewRateLimiter membuat instance baru dari RateLimiter
// limit: maximum requests allowed
// window: time window untuk counting requests (e.g., 15 minutes)
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*requestInfo),
		limit:    limit,
		window:   window,
	}

	// Cleanup goroutine untuk remove expired entries
	go rl.cleanupExpiredEntries()

	return rl
}

// RateLimit merupakan middleware handler untuk apply rate limiting
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get identifier (IP address atau user ID)
		identifier := rl.getIdentifier(c)

		// Check dan update rate limit
		allowed, resetTime := rl.allowRequest(identifier)

		if !allowed {
			c.Header("X-RateLimit-Limit", string(rune(rl.limit)))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Terlalu banyak permintaan. Silakan coba lagi nanti.",
				"retry_after": resetTime.Format("15:04:05"),
			})
			c.Abort()
			return
		}

		// Get remaining requests
		rl.mu.RLock()
		info := rl.requests[identifier]
		remaining := rl.limit - info.count
		rl.mu.RUnlock()

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", string(rune(rl.limit)))
		c.Header("X-RateLimit-Remaining", string(rune(remaining)))
		c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

		c.Next()
	}
}

// getIdentifier mendapatkan unique identifier untuk client (IP atau user ID)
func (rl *RateLimiter) getIdentifier(c *gin.Context) string {
	// Try to get user ID dari context (untuk authenticated requests)
	if userID, exists := c.Get("user_id"); exists {
		return "user_" + userID.(string)
	}

	// Fallback ke IP address untuk unauthenticated requests
	ip := c.ClientIP()
	return "ip_" + ip
}

// allowRequest memeriksa apakah request diizinkan berdasarkan rate limit
func (rl *RateLimiter) allowRequest(identifier string) (bool, time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Get atau create request info
	info, exists := rl.requests[identifier]
	if !exists {
		info = &requestInfo{
			count:      1,
			firstRequest: now,
		}
		rl.requests[identifier] = info
		return true, now.Add(rl.window)
	}

	// Check jika masih blocked
	if info.blocked && now.Before(info.blockedUntil) {
		return false, info.blockedUntil
	}

	// Reset block jika sudah expired
	if info.blocked && now.After(info.blockedUntil) {
		info.blocked = false
		info.count = 0
		info.firstRequest = now
	}

	// Check jika window sudah expired, reset counter
	if now.Sub(info.firstRequest) > rl.window {
		info.count = 1
		info.firstRequest = now
		return true, now.Add(rl.window)
	}

	// Increment counter
	info.count++

	// Check jika melebihi limit
	if info.count > rl.limit {
		info.blocked = true
		info.blockedUntil = info.firstRequest.Add(rl.window)
		return false, info.blockedUntil
	}

	return true, info.firstRequest.Add(rl.window)
}

// cleanupExpiredEntries membersihkan entries yang sudah expired
func (rl *RateLimiter) cleanupExpiredEntries() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()

		for identifier, info := range rl.requests {
			// Remove entries yang sudah tidak aktif > 2x window duration
			if now.Sub(info.firstRequest) > (rl.window * 2) {
				delete(rl.requests, identifier)
			}
		}

		rl.mu.Unlock()
	}
}

// LoginRateLimiter membuat rate limiter khusus untuk login endpoint
// dengan limit lebih ketat (5 requests per 15 minutes)
func LoginRateLimiter() gin.HandlerFunc {
	limiter := NewRateLimiter(5, 15*time.Minute)
	return limiter.RateLimit()
}

// APIRateLimiter membuat rate limiter untuk general API endpoints
// dengan limit lebih longgar (100 requests per minute)
func APIRateLimiter() gin.HandlerFunc {
	limiter := NewRateLimiter(100, 1*time.Minute)
	return limiter.RateLimit()
}

// StrictRateLimiter membuat rate limiter untuk sensitive operations
// dengan limit sangat ketat (3 requests per hour)
func StrictRateLimiter() gin.HandlerFunc {
	limiter := NewRateLimiter(3, 1*time.Hour)
	return limiter.RateLimit()
}

// IPWhitelist middleware untuk whitelist certain IPs (optional)
func IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	whitelist := make(map[string]bool)
	for _, ip := range allowedIPs {
		whitelist[ip] = true
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Allow localhost untuk development
		if clientIP == "127.0.0.1" || clientIP == "::1" {
			c.Next()
			return
		}

		// Check whitelist
		if !whitelist[clientIP] && len(whitelist) > 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access denied dari IP address ini",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityHeaders middleware untuk menambahkan security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Enforce HTTPS (uncomment untuk production dengan HTTPS)
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';")
		
		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Permissions Policy
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		c.Next()
	}
}

// ValidateContentType middleware untuk memvalidasi Content-Type header
func ValidateContentType(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip validation untuk GET dan DELETE requests
		if c.Request.Method == "GET" || c.Request.Method == "DELETE" {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")
		
		// Check jika Content-Type diizinkan
		allowed := false
		for _, allowedType := range allowedTypes {
			if contentType == allowedType || 
			   len(contentType) > len(allowedType) && contentType[:len(allowedType)] == allowedType {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"success": false,
				"message": "Content-Type tidak didukung",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
