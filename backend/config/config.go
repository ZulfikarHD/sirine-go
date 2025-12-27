package config

import (
	"os"
	"strconv"
	"time"
)

// Config merupakan struktur untuk application configuration
// yang di-load dari environment variables
type Config struct {
	// Server
	ServerPort string
	GinMode    string
	
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	
	// JWT
	JWTSecret           string
	JWTExpiry           time.Duration
	RefreshTokenExpiry  time.Duration
	
	// Security
	BcryptCost          int
	MaxLoginAttempts    int
	LockoutDuration     time.Duration
	
	// Frontend
	FrontendURL         string
	
	// Email (untuk Sprint 3)
	EmailSMTPHost       string
	EmailSMTPPort       int
	EmailUsername       string
	EmailPassword       string
	EmailFromAddress    string
}

// LoadConfig memuat configuration dari environment variables
// dengan default values untuk development
func LoadConfig() *Config {
	return &Config{
		// Server
		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),
		
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "sirine_go"),
		
		// JWT
		JWTSecret:          getEnv("JWT_SECRET", "sirine-go-jwt-secret-key-change-in-production"),
		JWTExpiry:          getDurationEnv("JWT_EXPIRY", 15*time.Minute),
		RefreshTokenExpiry: getDurationEnv("REFRESH_TOKEN_EXPIRY", 30*24*time.Hour), // 30 days
		
		// Security
		BcryptCost:       getIntEnv("BCRYPT_COST", 12),
		MaxLoginAttempts: getIntEnv("MAX_LOGIN_ATTEMPTS", 5),
		LockoutDuration:  getDurationEnv("LOCKOUT_DURATION", 15*time.Minute),
		
		// Frontend
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		
		// Email
		EmailSMTPHost:    getEnv("EMAIL_SMTP_HOST", "localhost"),
		EmailSMTPPort:    getIntEnv("EMAIL_SMTP_PORT", 1025),
		EmailUsername:    getEnv("EMAIL_USERNAME", ""),
		EmailPassword:    getEnv("EMAIL_PASSWORD", ""),
		EmailFromAddress: getEnv("EMAIL_FROM_ADDRESS", "noreply@sirine.local"),
	}
}

// Helper functions untuk read environment variables dengan default values

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}
