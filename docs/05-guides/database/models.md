# Adding New Models

## Overview

Guide ini menjelaskan cara menambahkan model/table baru ke database dengan **Registry Pattern** yang flexible dan automatic, yaitu: Anda hanya perlu menambahkan model di **satu tempat** untuk auto-migrate ke semua command.

## Quick Start

### Langkah-langkah Menambahkan Model Baru

1. **Buat file model** di `backend/models/`
2. **Register model** di `backend/database/models_registry.go`
3. **Migrate** dengan `make db-migrate`

---

## Step-by-Step Example

### 1. Buat Model Baru

Buat file model baru, contoh: `backend/models/product.go`

```go
package models

import (
	"time"
	"gorm.io/gorm"
)

// ProductStatus merupakan enum untuk status produk
type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "DRAFT"
	ProductStatusPublished ProductStatus = "PUBLISHED"
	ProductStatusArchived  ProductStatus = "ARCHIVED"
)

// Product merupakan model untuk entitas produk
type Product struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	SKU         string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"sku"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int            `gorm:"default:0" json:"stock"`
	Status      ProductStatus  `gorm:"type:enum('DRAFT','PUBLISHED','ARCHIVED');default:'DRAFT'" json:"status"`
	CreatedBy   uint64         `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Creator User `gorm:"foreignKey:CreatedBy;constraint:OnDelete:RESTRICT" json:"-"`
}

// TableName menentukan nama tabel di database
func (Product) TableName() string {
	return "products"
}
```

### 2. Register Model di Registry

Edit file `backend/database/models_registry.go`, tambahkan model baru:

```go
func NewModelsRegistry() *ModelsRegistry {
	registry := &ModelsRegistry{
		models: make([]interface{}, 0),
		tables: make([]string, 0),
	}

	// Existing models
	registry.Register(&models.User{}, "users")
	
	// ✨ MODEL BARU - Tambahkan di sini!
	// Note: Taruh model dengan foreign key SETELAH parent model
	registry.Register(&models.Product{}, "products")

	return registry
}
```

---

## Best Practices

### Model Design

1. **Gunakan GORM tags** yang tepat untuk definisi column:
   ```go
   Field string `gorm:"type:varchar(255);not null;uniqueIndex" json:"field_name"`
   ```

2. **Explicit column names** untuk fields yang mungkin bermasalah dengan naming convention:
   ```go
   NIP string `gorm:"column:nip;..." json:"nip"`
   ```

3. **Enum types** untuk status fields.

4. **Soft delete** dengan `DeletedAt`.

5. **Relationships** dengan proper foreign keys.

### Registry Order

**Urutan registration penting!**

```go
// ✅ BENAR - Parent models dulu
registry.Register(&models.User{}, "users")           // Parent
registry.Register(&models.Product{}, "products")     // Child (depends on User)

// ❌ SALAH - Child sebelum parent
registry.Register(&models.Product{}, "products")     // Error! User belum ada
registry.Register(&models.User{}, "users")
```

---

## Advanced: Complex Relationships

### One-to-Many
```go
type User struct {
	Products []Product `gorm:"foreignKey:UserID"`
}
type Product struct {
	UserID uint64
	User   User `gorm:"foreignKey:UserID"`
}
```

### Many-to-Many
```go
type Product struct {
	Categories []Category `gorm:"many2many:product_categories;"`
}
// Join table otomatis: product_categories
```

---

---

## 📊 Existing Models (Sprint 1-5)

Berikut adalah model-model yang sudah ada di sistem:

### Core Models (Sprint 1)

#### 1. User Model
**File:** `backend/models/user.go`  
**Table:** `users`

Model utama untuk user authentication dan authorization.

**Fields:**
- `id` (uint64, PK, auto-increment)
- `nip` (string, unique, max 5 chars)
- `full_name` (string, 255 chars)
- `email` (string, unique)
- `phone` (string, 15 chars)
- `password_hash` (string, stored hash)
- `role` (enum: ADMIN, MANAGER_KHAZWAL, STAFF_KHAZWAL, dll)
- `department` (enum: KHAZWAL, KEUANGAN, DISTRIBUSI, ADMIN)
- `shift` (enum: PAGI, SORE, MALAM)
- `status` (enum: ACTIVE, INACTIVE, LOCKED)
- `must_change_password` (bool, default: false)
- `last_login_at` (timestamp, nullable)
- `failed_login_attempts` (int, default: 0)
- `locked_until` (timestamp, nullable)
- `total_points` (int, default: 0) - **Sprint 5**
- `level` (string, calculated from points) - **Sprint 5**
- `profile_photo` (string, file path) - **Sprint 5**
- `created_at`, `updated_at`, `deleted_at` (timestamps)

**Relationships:**
- Has Many: `UserSessions`, `UserAchievements`, `ActivityLogs`, `Notifications`

---

#### 2. UserSession Model
**File:** `backend/models/user_session.go`  
**Table:** `user_sessions`

Untuk tracking active sessions dan token management.

**Fields:**
- `id` (uint64, PK)
- `user_id` (uint64, FK → users.id)
- `token_hash` (string, SHA256 hash of refresh token)
- `ip_address` (string)
- `user_agent` (string)
- `expires_at` (timestamp)
- `revoked_at` (timestamp, nullable)
- `created_at` (timestamp)

**Relationships:**
- Belongs To: `User`

---

#### 3. PasswordResetToken Model
**File:** `backend/models/password_reset_token.go`  
**Table:** `password_reset_tokens`

Untuk forgot password flow (Sprint 3).

**Fields:**
- `id` (uint64, PK)
- `user_id` (uint64, FK → users.id)
- `token_hash` (string, SHA256 hash)
- `expires_at` (timestamp, 1 hour from creation)
- `used_at` (timestamp, nullable)
- `created_at` (timestamp)

**Relationships:**
- Belongs To: `User`

**Methods:**
- `IsValid()` - Check if token valid (not used, not expired)
- `IsExpired()` - Check expiration
- `IsUsed()` - Check if already used

---

#### 4. ActivityLog Model (Sprint 1 & 4)
**File:** `backend/models/activity_log.go`  
**Table:** `activity_logs`

Untuk comprehensive audit trail.

**Fields:**
- `id` (uint64, PK)
- `user_id` (uint64, FK → users.id)
- `action` (enum: CREATE, UPDATE, DELETE, LOGIN, LOGOUT, dll)
- `entity_type` (string, e.g., "users", "profile")
- `entity_id` (string)
- `changes` (JSON, before/after values)
- `ip_address` (string)
- `user_agent` (string)
- `created_at` (timestamp)

**Relationships:**
- Belongs To: `User`

**Methods:**
- `SetChanges(before, after interface{})` - Store changes as JSON
- `GetChanges()` - Parse changes from JSON

---

### Notification System (Sprint 4)

#### 5. Notification Model
**File:** `backend/models/notification.go`  
**Table:** `notifications`

Untuk in-app notification system.

**Fields:**
- `id` (uint64, PK)
- `user_id` (uint64, FK → users.id)
- `title` (string)
- `message` (text)
- `type` (enum: INFO, SUCCESS, WARNING, ERROR)
- `is_read` (bool, default: false)
- `read_at` (timestamp, nullable)
- `created_at`, `updated_at` (timestamps)

**Relationships:**
- Belongs To: `User`

**Methods:**
- `MarkAsRead()` - Mark notification as read dengan timestamp

---

### Gamification System (Sprint 5)

#### 6. Achievement Model
**File:** `backend/models/achievement.go`  
**Table:** `achievements`

Master table untuk achievement definitions.

**Fields:**
- `id` (uint64, PK)
- `code` (string, unique, e.g., "FIRST_LOGIN")
- `name` (string)
- `description` (text)
- `icon` (string, emoji atau icon class)
- `points` (int)
- `category` (enum: LOGIN, PRODUCTIVITY, QUALITY, MILESTONE)
- `criteria` (JSON, achievement criteria)
- `is_active` (bool, default: true)
- `created_at`, `updated_at` (timestamps)

**Relationships:**
- Has Many: `UserAchievements` (through many-to-many)

---

#### 7. UserAchievement Model
**File:** `backend/models/user_achievement.go`  
**Table:** `user_achievements`

Junction table untuk tracking user achievement unlocks.

**Fields:**
- `id` (uint64, PK)
- `user_id` (uint64, FK → users.id)
- `achievement_id` (uint64, FK → achievements.id)
- `unlocked_at` (timestamp)

**Relationships:**
- Belongs To: `User`
- Belongs To: `Achievement`

**Indexes:**
- Unique composite index: (user_id, achievement_id)

---

## 🔍 Model Registry

Semua models di-register di `backend/database/models_registry.go` dengan urutan:

```go
func NewModelsRegistry() *ModelsRegistry {
    registry := &ModelsRegistry{
        models: make([]interface{}, 0),
        tables: make([]string, 0),
    }

    // Core models
    registry.Register(&models.User{}, "users")
    registry.Register(&models.UserSession{}, "user_sessions")
    registry.Register(&models.PasswordResetToken{}, "password_reset_tokens")
    registry.Register(&models.ActivityLog{}, "activity_logs")
    
    // Notification system
    registry.Register(&models.Notification{}, "notifications")
    
    // Gamification system
    registry.Register(&models.Achievement{}, "achievements")
    registry.Register(&models.UserAchievement{}, "user_achievements")

    return registry
}
```

**Urutan penting:** Parent models harus di-register sebelum child models dengan foreign keys.

---

## 📝 Model Conventions

### Naming
- Model struct: PascalCase singular (`User`, `Achievement`)
- Table name: snake_case plural (`users`, `achievements`)
- Fields: camelCase di struct, snake_case di database

### Standard Fields
Semua models sebaiknya include:
```go
CreatedAt   time.Time      `gorm:"autoCreateTime"`
UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
```

### Enums
Define sebagai custom type dengan constants:
```go
type UserRole string

const (
    RoleAdmin    UserRole = "ADMIN"
    RoleManager  UserRole = "MANAGER"
    RoleStaff    UserRole = "STAFF"
)
```

### JSON Tags
- Include `json` tag untuk API responses
- Use `json:"-"` untuk fields yang tidak boleh ter-expose (password, deleted_at)

---

## 🔗 Related Documentation

- [Migration Guide](./migrations.md) - Cara menjalankan migrasi dan seeding
- [Database Management](./management.md) - Database maintenance dan backup
- [API Documentation](../../development/api-documentation.md) - API endpoints untuk semua models
