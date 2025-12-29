# Database Migrations Guide

## Overview

Sirine Go menggunakan automatic migrations yang terintegrasi dengan **Registry Pattern**. Tidak perlu membuat file migration manual (seperti `2024_01_01_create_users_table.php`), cukup update model dan registry.

---

## 🚀 Migration Workflow

Workflow migrasi berjalan otomatis melalui Makefile commands:

### 1. Test Migration (Development)
Untuk mengetes apakah model baru berhasil dimigrate:

```bash
make db-migrate
```

### 2. Fresh Migration (Reset)
Untuk menghapus semua data dan migrate ulang dari awal (HANYA DEV):

```bash
make db-fresh
```

### 3. Rollback
Untuk drop semua tables:

```bash
make db-rollback
```

---

## ✅ Verifikasi Database

Setelah menjalankan migration, verifikasi bahwa table telah dibuat dengan benar:

### Command Line Check
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
| products            |
| user_sessions       |
| users               |
+---------------------+
```

### Column Check
Untuk memastikan column attributes (null, default, type) benar:
```bash
mysql -u root -p sirine_go -e "DESCRIBE users;"
```

---

## 🌱 Database Seeding

Seeding digunakan untuk mengisi data awal (seperti admin user, default categories, dll).

### Menjalankan Seeder
```bash
make db-seed
```

### Menambahkan Seeder Baru
Edit file `backend/cmd/seed/main.go`:

```go
func main() {
	// ... existing code ...

	log.Println("🌱 Starting database seeding...")

	seedAdminUser()
	seedDemoUsers()
	seedProducts() // ← Function baru
	
	log.Println("✅ Database seeding completed!")
}

// Function seeder baru
func seedProducts() {
	db := database.GetDB()

	products := []models.Product{
		{
			SKU:         "PROD-001",
			Name:        "Sample Product",
			Price:       99.99,
			Status:      models.ProductStatusPublished,
			CreatedBy:   1,
		},
	}

	for _, product := range products {
		if err := db.Where("sku = ?", product.SKU).First(&models.Product{}).Error; err == nil {
			continue // Skip if exists
		}
		db.Create(&product)
		log.Printf("✅ Product seeded: %s", product.SKU)
	}
}
```

---

## ❓ Troubleshooting

### Error: Table already exists
```
Error 1050: Table 'products' already exists
```
**Solution:**
Gunakan `make db-fresh` untuk reset database development.

### Error: Foreign key constraint fails
```
Error 1452: Cannot add or update a child row
```
**Solution:**
Pastikan urutan registry benar (Parent sebelum Child). Lihat [Models Guide](./models.md#registry-order).

### Migration tidak detect perubahan
**Problem:** GORM AutoMigrate tidak mengubah column type yang sudah ada.
**Solution:**
- Development: `make db-fresh`
- Production: Manual SQL `ALTER TABLE ...`

---

## 🔗 Related Documentation

- [Creating New Models](./models.md)
- [Database Management](./management.md)
