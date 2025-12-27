# Adding New Models Guide

## Overview

Guide ini menjelaskan cara menambahkan model/table baru ke database dengan **Registry Pattern** yang flexible dan automatic, yaitu: Anda hanya perlu menambahkan model di **satu tempat** untuk auto-migrate ke semua command.

## Quick Start

### Langkah-langkah Menambahkan Model Baru

1. **Buat file model** di `backend/models/`
2. **Register model** di `backend/database/models_registry.go`
3. **Selesai!** 🎉 Model otomatis ter-migrate di semua command

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
// yang digunakan untuk mengelola data produk dalam sistem
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

// IsAvailable memeriksa apakah produk tersedia untuk dijual
func (p *Product) IsAvailable() bool {
	return p.Status == ProductStatusPublished && p.Stock > 0
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
	registry.Register(&models.UserSession{}, "user_sessions")
	registry.Register(&models.PasswordResetToken{}, "password_reset_tokens")
	registry.Register(&models.ActivityLog{}, "activity_logs")
	registry.Register(&models.Notification{}, "notifications")

	// ✨ MODEL BARU - Tambahkan di sini!
	registry.Register(&models.Product{}, "products")

	return registry
}
```

**Catatan Penting:**
- **Order matters!** Taruh model dengan foreign key **setelah** parent model
- Contoh: `Product` depends on `User`, jadi `Product` harus setelah `User`

### 3. Test Migration

Jalankan migration untuk test:

```bash
# Fresh migration dengan model baru
make db-fresh

# Expected output:
# ✅ Migrations completed! (6 tables migrated)
```

### 4. Verify Database

Check apakah table berhasil dibuat:

```bash
mysql -u root -p sirine_go -e "SHOW TABLES;"
```

Expected output:
```
+---------------------+
| Tables_in_sirine_go |
+---------------------+
| activity_logs       |
| notifications       |
| password_reset_tokens |
| products            | ← Table baru!
| user_sessions       |
| users               |
+---------------------+
```

### 5. (Optional) Add Seeder

Jika ingin seed data awal untuk model baru, edit `backend/cmd/seed/main.go`:

```go
func main() {
	// ... existing code ...

	log.Println("🌱 Starting database seeding...")

	seedAdminUser()
	seedDemoUsers()
	seedProducts() // ← Function baru untuk seed products

	log.Println("✅ Database seeding completed!")
}

// seedProducts membuat sample products untuk testing
func seedProducts() {
	db := database.GetDB()

	products := []models.Product{
		{
			SKU:         "PROD-001",
			Name:        "Sample Product 1",
			Description: "This is a sample product",
			Price:       99.99,
			Stock:       100,
			Status:      models.ProductStatusPublished,
			CreatedBy:   1, // Admin user ID
		},
		// ... more products
	}

	for _, product := range products {
		var existing models.Product
		if err := db.Where("sku = ?", product.SKU).First(&existing).Error; err == nil {
			log.Printf("ℹ️  Product %s sudah ada, skip", product.SKU)
			continue
		}

		if err := db.Create(&product).Error; err != nil {
			log.Printf("Warning: Failed to seed product %s: %v", product.SKU, err)
			continue
		}

		log.Printf("✅ Product seeded: %s - %s", product.SKU, product.Name)
	}
}
```

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

3. **Enum types** untuk status fields:
   ```go
   type Status string
   const (
       StatusActive Status = "ACTIVE"
       StatusInactive Status = "INACTIVE"
   )
   ```

4. **Soft delete** dengan `DeletedAt`:
   ```go
   DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
   ```

5. **Relationships** dengan proper foreign keys:
   ```go
   UserID uint64 `gorm:"not null" json:"user_id"`
   User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
   ```

### Registry Order

**Urutan registration penting!** Follow these rules:

```go
// ✅ BENAR - Parent models dulu
registry.Register(&models.User{}, "users")           // Parent
registry.Register(&models.Product{}, "products")     // Child (depends on User)
registry.Register(&models.OrderItem{}, "order_items") // Child (depends on Product)

// ❌ SALAH - Child sebelum parent
registry.Register(&models.Product{}, "products")     // Error! User belum ada
registry.Register(&models.User{}, "users")
```

**Dependency Order:**
1. Independent tables (no foreign keys)
2. Tables dengan foreign keys ke #1
3. Tables dengan foreign keys ke #2
4. And so on...

### Naming Conventions

1. **Model name**: Singular, PascalCase
   ```go
   type Product struct { }
   ```

2. **Table name**: Plural, snake_case
   ```go
   func (Product) TableName() string {
       return "products"
   }
   ```

3. **Field names**: PascalCase in struct, snake_case in DB
   ```go
   CreatedAt time.Time // → created_at in database
   ```

4. **Enum values**: UPPER_SNAKE_CASE
   ```go
   const ProductStatusActive ProductStatus = "ACTIVE"
   ```

## Advanced: Complex Relationships

### One-to-Many

```go
type User struct {
	ID       uint64    `gorm:"primaryKey"`
	Products []Product `gorm:"foreignKey:UserID"` // One user has many products
}

type Product struct {
	ID     uint64 `gorm:"primaryKey"`
	UserID uint64 `gorm:"not null"`
	User   User   `gorm:"foreignKey:UserID"`
}
```

### Many-to-Many

```go
type Product struct {
	ID         uint64     `gorm:"primaryKey"`
	Categories []Category `gorm:"many2many:product_categories;"`
}

type Category struct {
	ID       uint64    `gorm:"primaryKey"`
	Products []Product `gorm:"many2many:product_categories;"`
}

// Join table akan dibuat otomatis: product_categories
```

### Self-Referencing

```go
type Category struct {
	ID       uint64     `gorm:"primaryKey"`
	ParentID *uint64    `gorm:"index"`
	Parent   *Category  `gorm:"foreignKey:ParentID"`
	Children []Category `gorm:"foreignKey:ParentID"`
}
```

## Troubleshooting

### Error: Foreign key constraint fails

**Problem:**
```
Error 1452: Cannot add or update a child row: a foreign key constraint fails
```

**Solution:**
Pastikan parent model di-register **sebelum** child model di registry:
```go
registry.Register(&models.User{}, "users")      // Parent dulu
registry.Register(&models.Product{}, "products") // Child setelahnya
```

### Error: Table already exists

**Problem:**
```
Error 1050: Table 'products' already exists
```

**Solution:**
```bash
# Option 1: Rollback dan migrate ulang
make db-rollback
make db-migrate

# Option 2: Fresh migration
make db-fresh
```

### Error: Unknown column

**Problem:**
```
Error 1054: Unknown column 'product_s_k_u' in 'where clause'
```

**Solution:**
Gunakan explicit column tag:
```go
SKU string `gorm:"column:sku;..." json:"sku"` // Explicit column name
```

### Migration tidak detect perubahan

**Problem:**
GORM AutoMigrate tidak detect perubahan field type.

**Solution:**
GORM AutoMigrate **tidak mengubah** column type yang sudah ada. Untuk perubahan schema:
```bash
# Drop dan recreate (development only!)
make db-fresh
```

Atau manual alter table di production:
```sql
ALTER TABLE products MODIFY COLUMN price DECIMAL(12,2);
```

## Examples

### Complete Model Example

```go
package models

import (
	"time"
	"gorm.io/gorm"
)

// OrderStatus merupakan enum untuk status order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// Order merupakan model untuk transaksi order
type Order struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNumber   string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	UserID        uint64         `gorm:"not null;index" json:"user_id"`
	TotalAmount   float64        `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status        OrderStatus    `gorm:"type:enum('PENDING','CONFIRMED','SHIPPED','DELIVERED','CANCELLED');default:'PENDING'" json:"status"`
	OrderDate     time.Time      `gorm:"not null" json:"order_date"`
	ShippedDate   *time.Time     `gorm:"type:timestamp null" json:"shipped_date"`
	DeliveredDate *time.Time     `gorm:"type:timestamp null" json:"delivered_date"`
	Notes         string         `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User       User        `gorm:"foreignKey:UserID;constraint:OnDelete:RESTRICT" json:"-"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// TableName menentukan nama tabel
func (Order) TableName() string {
	return "orders"
}

// OrderItem merupakan model untuk item dalam order
type OrderItem struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   uint64    `gorm:"not null;index" json:"order_id"`
	ProductID uint64    `gorm:"not null;index" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	UnitPrice float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal  float64   `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Order   Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"-"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT" json:"-"`
}

// TableName menentukan nama tabel
func (OrderItem) TableName() string {
	return "order_items"
}
```

**Registration di registry:**
```go
registry.Register(&models.User{}, "users")           // 1. Independent
registry.Register(&models.Product{}, "products")     // 2. Depends on User
registry.Register(&models.Order{}, "orders")         // 3. Depends on User
registry.Register(&models.OrderItem{}, "order_items") // 4. Depends on Order & Product
```

## Summary

### ✅ Keuntungan Registry Pattern

1. **Single Source of Truth** - Model hanya di-register di satu tempat
2. **Automatic Migration** - Semua command (server, migrate, rollback) auto-sync
3. **Easy Maintenance** - Tambah model baru hanya 2 langkah
4. **Type Safety** - Compile-time checking untuk models
5. **Flexible** - Mudah menambah/hapus models tanpa edit banyak file

### 📋 Checklist Menambah Model Baru

- [ ] Buat file model di `backend/models/`
- [ ] Define struct dengan GORM tags yang tepat
- [ ] Implement `TableName()` method
- [ ] Register di `backend/database/models_registry.go`
- [ ] Pastikan dependency order benar
- [ ] Test dengan `make db-fresh`
- [ ] Verify table dengan `SHOW TABLES`
- [ ] (Optional) Add seeder di `cmd/seed/main.go`
- [ ] Commit changes

## Support

Untuk pertanyaan atau issues, hubungi:
- Developer: Zulfikar Hidayatullah (+62 857-1583-8733)
