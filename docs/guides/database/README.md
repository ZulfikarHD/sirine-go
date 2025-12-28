# 🗄️ Database Documentation

## Overview

Sirine Go menggunakan MySQL 8.0+ dengan GORM sebagai ORM. Dokumentasi ini mencakup manajemen database, modeling, dan migrasi.

## Guides

### 1. [Management & Maintenance](./management.md)
Panduan operasional database.
- Backup & Restore
- Performance Tuning
- Troubleshooting

### 2. [Creating New Models](./models.md)
Panduan membuat model baru menggunakan **Registry Pattern**.
- Model Structure
- Relationships
- Best Practices

### 3. [Migrations & Seeding](./migrations.md)
Workflow migrasi database dan seeding data.
- Running Migrations (`make db-migrate`)
- Verifying Tables
- Seeding Data

---

## Related Commands

```bash
make db-fresh   # Reset database (Development only)
make db-migrate # Run pending migrations
make db-seed    # Seed initial data
```
