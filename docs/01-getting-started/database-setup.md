# 🗄️ Database Setup Guide

Panduan untuk setup MySQL database untuk Sirine Go App.

---

## 📋 Prerequisites

Pastikan MySQL 8.0+ sudah terinstall. Jika belum, lihat [prerequisites.md](./prerequisites.md).

---

## 🚀 Quick Setup

### **One-Command Setup:**

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

---

## 📝 Step-by-Step Setup

### **Step 1: Start MySQL Service**

```bash
# Check if MySQL is running
sudo systemctl status mysql

# Start MySQL if not running
sudo systemctl start mysql

# Enable auto-start on boot
sudo systemctl enable mysql
```

**Expected output:**
```
● mysql.service - MySQL Community Server
     Loaded: loaded
     Active: active (running)
```

---

### **Step 2: Login to MySQL**

```bash
mysql -u root -p
```

Enter your MySQL root password when prompted.

You should see the MySQL prompt:
```
mysql>
```

---

### **Step 3: Create Database**

```sql
-- Create database dengan UTF-8 support
CREATE DATABASE sirine_go 
  CHARACTER SET utf8mb4 
  COLLATE utf8mb4_unicode_ci;

-- Verify database created
SHOW DATABASES LIKE 'sirine_go';

-- Use database
USE sirine_go;

-- Show current database
SELECT DATABASE();
```

**Expected output:**
```
+-------------------+
| Database          |
+-------------------+
| sirine_go         |
+-------------------+
1 row in set (0.00 sec)
```

---

### **Step 4: Create Database User** (Optional but Recommended)

Untuk security, buat dedicated user untuk aplikasi instead of menggunakan root:

```sql
-- Create application user
CREATE USER 'sirine_user'@'localhost' 
  IDENTIFIED BY 'your_secure_password_here';

-- Grant privileges
GRANT ALL PRIVILEGES ON sirine_go.* 
  TO 'sirine_user'@'localhost';

-- Apply privileges
FLUSH PRIVILEGES;

-- Verify user created
SELECT User, Host FROM mysql.user WHERE User='sirine_user';

-- Exit MySQL
EXIT;
```

**Test new user:**
```bash
mysql -u sirine_user -p sirine_go
# Enter password

# You should be logged in successfully
mysql> SELECT DATABASE();
+------------+
| DATABASE() |
+------------+
| sirine_go  |
+------------+

mysql> EXIT;
```

---

### **Step 5: Verify Database Setup**

```bash
# Check database exists
mysql -u root -p -e "SHOW DATABASES LIKE 'sirine_go';"

# Check database character set
mysql -u root -p -e "SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME='sirine_go';"
```

**Expected output:**
```
+----------------------------+------------------------+
| DEFAULT_CHARACTER_SET_NAME | DEFAULT_COLLATION_NAME |
+----------------------------+------------------------+
| utf8mb4                    | utf8mb4_unicode_ci     |
+----------------------------+------------------------+
```

---

## ⚙️ Database Configuration

### **Character Set & Collation**

Aplikasi menggunakan `utf8mb4` untuk full Unicode support (termasuk emoji 😊).

**Why utf8mb4?**
- ✅ Support semua Unicode characters (emoji, Asian languages, etc.)
- ✅ Standard MySQL recommendation untuk UTF-8
- ✅ Better than `utf8` (which is actually utf8mb3)

**Collation: utf8mb4_unicode_ci**
- `unicode` = Unicode-aware sorting
- `ci` = Case Insensitive (A = a)

---

## 🔧 Database Tables

Tables akan dibuat otomatis oleh GORM Auto-Migration saat backend pertama kali running.

**Expected tables setelah backend start:**

```sql
-- Check tables
mysql -u root -p sirine_go -e "SHOW TABLES;"
```

**Output (setelah backend running):**

```
+---------------------+
| Tables_in_sirine_go |
+---------------------+
| users               |
| sessions            |
| activity_logs       |
| notifications       |
| achievements        |
| user_achievements   |
| password_resets     |
+---------------------+
```

### **View Table Structure:**

```bash
# Describe users table
mysql -u root -p sirine_go -e "DESCRIBE users;"

# Show create statement
mysql -u root -p sirine_go -e "SHOW CREATE TABLE users\G"
```

---

## 🔄 Database Migrations

GORM Auto-Migration handles schema changes automatically.

### **How Auto-Migration Works:**

1. Backend starts
2. GORM compares model structs dengan database schema
3. Creates missing tables
4. Adds missing columns
5. **TIDAK** delete columns (safe)
6. **TIDAK** modify column types (safe)

### **Manual Migrations** (Advanced)

Jika butuh complex migrations:

```sql
-- Example: Add new column
ALTER TABLE users ADD COLUMN phone VARCHAR(20) AFTER email;

-- Example: Add index
CREATE INDEX idx_users_email ON users(email);

-- Example: Add foreign key
ALTER TABLE activity_logs 
  ADD CONSTRAINT fk_user_id 
  FOREIGN KEY (user_id) REFERENCES users(id);
```

---

## 💾 Database Backup & Restore

### **Backup Database:**

```bash
# Full backup
mysqldump -u root -p sirine_go > backup_$(date +%Y%m%d).sql

# Backup with compression
mysqldump -u root -p sirine_go | gzip > backup_$(date +%Y%m%d).sql.gz

# Backup specific tables
mysqldump -u root -p sirine_go users sessions > users_backup.sql
```

### **Restore Database:**

```bash
# Restore from backup
mysql -u root -p sirine_go < backup_20251228.sql

# Restore from compressed backup
gunzip < backup_20251228.sql.gz | mysql -u root -p sirine_go
```

### **Automated Backup Script:**

```bash
#!/bin/bash
# Save as: backup_db.sh

BACKUP_DIR="/home/user/db_backups"
DB_NAME="sirine_go"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR
mysqldump -u root -p$DB_PASSWORD $DB_NAME | gzip > $BACKUP_DIR/backup_$DATE.sql.gz

# Keep only last 7 days backups
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete

echo "Backup completed: backup_$DATE.sql.gz"
```

**Schedule backup dengan cron:**
```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /path/to/backup_db.sh
```

---

## 🐛 Troubleshooting

### **Problem: Can't connect to MySQL server**

```bash
# Check if MySQL is running
sudo systemctl status mysql

# Start MySQL
sudo systemctl start mysql

# Check MySQL is listening
sudo netstat -tlnp | grep mysql
```

### **Problem: Access denied for user 'root'@'localhost'**

**Solution 1: Reset root password**
```bash
# Stop MySQL
sudo systemctl stop mysql

# Start MySQL in safe mode
sudo mysqld_safe --skip-grant-tables &

# Login without password
mysql -u root

# Reset password
mysql> FLUSH PRIVILEGES;
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY 'new_password';
mysql> EXIT;

# Restart MySQL normally
sudo systemctl restart mysql
```

**Solution 2: Use sudo to login**
```bash
sudo mysql -u root
```

### **Problem: Database charset not utf8mb4**

```sql
-- Check current charset
SELECT DEFAULT_CHARACTER_SET_NAME 
FROM information_schema.SCHEMATA 
WHERE SCHEMA_NAME='sirine_go';

-- Change database charset
ALTER DATABASE sirine_go 
  CHARACTER SET utf8mb4 
  COLLATE utf8mb4_unicode_ci;
```

### **Problem: Tables not created after backend start**

**Solution:**
```bash
# Check backend logs untuk error messages
# Backend akan log migration errors

# Manual check backend connection
mysql -u root -p sirine_go -e "SELECT 1;"

# Verify .env database credentials correct
cat backend/.env | grep DB_
```

---

## ✅ Verification Checklist

- [ ] MySQL service running (`sudo systemctl status mysql`)
- [ ] Database `sirine_go` created (`SHOW DATABASES;`)
- [ ] Database using utf8mb4 charset
- [ ] Can login to database (`mysql -u root -p sirine_go`)
- [ ] Application user created (if using dedicated user)
- [ ] Privileges granted correctly

---

## 🎯 Next Steps

Setelah database setup:

1. ✅ Database ready
2. ➡️ Configure backend environment: [backend-setup.md](./backend-setup.md)
3. ➡️ Or continue with complete guide: [installation.md](./installation.md)

---

## 📚 Related Documentation

- [prerequisites.md](./prerequisites.md) - System requirements
- [backend-setup.md](./backend-setup.md) - Backend configuration
- [../05-guides/database/management.md](../05-guides/database/management.md) - Database management
- [../05-guides/database/models.md](../05-guides/database/models.md) - Database models

---

## 📞 Need Help?

Jika ada masalah dengan database setup:
- **Developer:** Zulfikar Hidayatullah
- **Phone:** +62 857-1583-8733

---

**Last Updated:** 28 Desember 2025  
**Version:** 2.0.0 (Phase 2 Restructure)
