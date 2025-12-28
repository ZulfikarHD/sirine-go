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

## 🔗 Related Documentation

- [Migration Guide](./migrations.md) - Cara menjalankan migrasi dan seeding
- [Database Management](./management.md)
