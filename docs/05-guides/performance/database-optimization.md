# 🗄️ Database Performance Optimization - MySQL

Panduan optimasi performa database MySQL untuk aplikasi Sirine Go dengan fokus pada query optimization dan index strategy.

## 📋 Daftar Isi

1. [Query Optimization](#query-optimization)
2. [Index Strategy](#index-strategy)
3. [Table Design](#table-design-optimization)
4. [Connection Management](#connection-management)
5. [Caching](#database-caching)
6. [Monitoring & Analysis](#monitoring--analysis)
7. [Maintenance](#database-maintenance)

---

## 🎯 Query Optimization

### **Use EXPLAIN to Analyze Queries**

```sql
-- Check query execution plan
EXPLAIN SELECT * FROM users 
WHERE email = 'test@example.com';

-- Extended analysis
EXPLAIN ANALYZE SELECT * FROM users 
WHERE email = 'test@example.com';
```

**Key Metrics:**
- `type: ALL` = ❌ Full table scan (BAD)
- `type: ref` = ✅ Index used (GOOD)
- `type: const` = ✅ Primary key lookup (BEST)
- `Extra: Using filesort` = ⚠️ Needs optimization
- `Extra: Using temporary` = ⚠️ Needs optimization

### **Avoid SELECT * Queries**

```sql
-- ❌ BAD: Fetches all columns (slow & wasteful)
SELECT * FROM users WHERE id = 1;

-- ✅ GOOD: Only fetch needed columns
SELECT id, name, email, role FROM users WHERE id = 1;
```

**Benefits:**
- Reduced network bandwidth
- Faster query execution
- Lower memory usage
- Better index utilization

### **Use LIMIT for Large Result Sets**

```sql
-- ❌ BAD: No limit (could return millions of rows)
SELECT * FROM activity_logs 
WHERE created_at > '2025-01-01';

-- ✅ GOOD: Paginated results
SELECT id, user_id, action, created_at 
FROM activity_logs 
WHERE created_at > '2025-01-01'
ORDER BY created_at DESC
LIMIT 20 OFFSET 0;
```

### **Optimize JOIN Queries**

```sql
-- ❌ BAD: Join without proper indexes
SELECT u.name, p.bio 
FROM users u
JOIN profiles p ON u.id = p.user_id
WHERE u.status = 'active';

-- ✅ GOOD: Add indexes + select specific columns
-- First, add indexes:
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_profiles_user_id ON profiles(user_id);

-- Then query:
SELECT u.name, p.bio 
FROM users u
INNER JOIN profiles p ON u.id = p.user_id
WHERE u.status = 'active';
```

### **Use IN Instead of Multiple OR**

```sql
-- ❌ BAD: Multiple OR conditions
SELECT * FROM users 
WHERE role = 'admin' 
   OR role = 'moderator' 
   OR role = 'manager';

-- ✅ GOOD: Use IN clause
SELECT id, name, email, role FROM users 
WHERE role IN ('admin', 'moderator', 'manager');
```

### **Avoid Functions on Indexed Columns**

```sql
-- ❌ BAD: Function prevents index usage
SELECT * FROM users 
WHERE LOWER(email) = 'test@example.com';

-- ✅ GOOD: Store email in lowercase, search as-is
SELECT id, name, email FROM users 
WHERE email = 'test@example.com';
```

---

## 📇 Index Strategy

### **Current Indexes (Sirine Go)**

```sql
-- Users table
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created ON users(created_at);

-- User sessions
CREATE INDEX idx_sessions_token ON user_sessions(token);
CREATE INDEX idx_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_sessions_expires ON user_sessions(expires_at);

-- Notifications
CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_read ON notifications(is_read);
CREATE INDEX idx_notifications_created ON notifications(created_at);
CREATE INDEX idx_notifications_user_read ON notifications(user_id, is_read);

-- Activity logs
CREATE INDEX idx_activity_user ON activity_logs(user_id);
CREATE INDEX idx_activity_created ON activity_logs(created_at);
CREATE INDEX idx_activity_user_created ON activity_logs(user_id, created_at);

-- Achievements
CREATE INDEX idx_user_achievements_user ON user_achievements(user_id);
CREATE INDEX idx_user_achievements_achievement ON user_achievements(achievement_id);
```

### **Index Guidelines**

**When to Create Index:**
- ✅ Columns used in WHERE clauses
- ✅ Columns used in JOIN conditions
- ✅ Columns used in ORDER BY
- ✅ Foreign key columns
- ✅ Columns frequently searched

**When NOT to Create Index:**
- ❌ Small tables (< 1000 rows)
- ❌ Columns with low cardinality (few distinct values)
- ❌ Columns frequently updated
- ❌ Too many indexes (slows INSERT/UPDATE)

### **Composite Indexes**

```sql
-- Scenario: Search by user_id and filter by is_read
-- ✅ GOOD: Composite index
CREATE INDEX idx_notifications_user_read 
ON notifications(user_id, is_read);

-- Now this query uses index efficiently:
SELECT * FROM notifications 
WHERE user_id = 123 AND is_read = 0;

-- Note: Order matters!
-- This index helps: user_id only
-- This index helps: user_id + is_read
-- This index does NOT help: is_read only
```

### **Index Maintenance**

```sql
-- Check index usage
SELECT 
    table_name,
    index_name,
    column_name,
    cardinality
FROM information_schema.statistics
WHERE table_schema = 'sirine_go';

-- Find unused indexes
SELECT 
    OBJECT_SCHEMA,
    OBJECT_NAME,
    INDEX_NAME
FROM performance_schema.table_io_waits_summary_by_index_usage
WHERE INDEX_NAME IS NOT NULL
AND COUNT_STAR = 0;

-- Rebuild indexes
ANALYZE TABLE users;
OPTIMIZE TABLE users;
```

---

## 🏗️ Table Design Optimization

### **Use Appropriate Data Types**

```sql
-- ❌ BAD: Oversized data types
CREATE TABLE users (
    id VARCHAR(255),           -- Too large!
    name VARCHAR(500),         -- Too large!
    age VARCHAR(100),          -- Should be INT!
    is_active VARCHAR(10)      -- Should be BOOLEAN!
);

-- ✅ GOOD: Appropriate data types
CREATE TABLE users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    age TINYINT UNSIGNED,      -- 0-255
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### **Normalization vs Denormalization**

**Normalized (Sirine Go Current):**
```sql
-- ✅ GOOD: Avoid data duplication
CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(150)
);

CREATE TABLE profiles (
    id INT PRIMARY KEY,
    user_id INT,
    bio TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

**Selective Denormalization (for performance):**
```sql
-- ✅ GOOD: Cache frequently accessed data
CREATE TABLE notifications (
    id INT PRIMARY KEY,
    user_id INT,
    user_name VARCHAR(100),    -- Denormalized for faster display
    message TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### **Partitioning Large Tables**

```sql
-- Partition activity_logs by month
CREATE TABLE activity_logs (
    id INT NOT NULL,
    user_id INT,
    action VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id, created_at)
)
PARTITION BY RANGE (YEAR(created_at) * 100 + MONTH(created_at)) (
    PARTITION p202501 VALUES LESS THAN (202502),
    PARTITION p202502 VALUES LESS THAN (202503),
    PARTITION p202503 VALUES LESS THAN (202504),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);
```

---

## 🔌 Connection Management

### **Optimal Connection Pool Settings**

```go
// File: backend/database/database.go

sqlDB.SetMaxOpenConns(25)        // Max open connections
sqlDB.SetMaxIdleConns(10)        // Max idle connections
sqlDB.SetConnMaxLifetime(time.Hour)      // Connection lifetime
sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Idle timeout
```

**Guidelines by server size:**

| Server Spec | MaxOpenConns | MaxIdleConns | Justification |
|-------------|--------------|--------------|---------------|
| Small (1 CPU, 1GB RAM) | 10 | 5 | Avoid overload |
| Medium (2-4 CPU, 4GB RAM) | 25 | 10 | ✅ Recommended |
| Large (4+ CPU, 8GB RAM) | 50 | 20 | High concurrency |

### **Monitor Connection Usage**

```sql
-- Check current connections
SHOW PROCESSLIST;

-- Check connection stats
SHOW STATUS LIKE 'Threads_connected';
SHOW STATUS LIKE 'Max_used_connections';

-- Check wait timeout
SHOW VARIABLES LIKE 'wait_timeout';
```

### **Prepared Statements (GORM Default)**

```go
// ✅ GORM automatically uses prepared statements
db.Where("email = ?", email).First(&user)

// SQL executed:
// PREPARE stmt FROM 'SELECT * FROM users WHERE email = ?'
// EXECUTE stmt USING 'test@example.com'
```

**Benefits:**
- SQL injection prevention
- Query plan caching
- Faster repeated queries

---

## 💾 Database Caching

### **Query Result Caching (Application Level)**

```go
// File: backend/services/user_service.go

var userCache = make(map[int]*models.User)
var cacheMutex sync.RWMutex

func GetUserByID(id int) (*models.User, error) {
    // Check cache first
    cacheMutex.RLock()
    if user, found := userCache[id]; found {
        cacheMutex.RUnlock()
        return user, nil
    }
    cacheMutex.RUnlock()
    
    // Query database
    var user models.User
    err := db.First(&user, id).Error
    if err != nil {
        return nil, err
    }
    
    // Cache result
    cacheMutex.Lock()
    userCache[id] = &user
    cacheMutex.Unlock()
    
    return &user, nil
}
```

### **MySQL Query Cache (Legacy)**

```sql
-- ⚠️ Note: Query cache removed in MySQL 8.0+
-- Use application-level caching instead
```

### **InnoDB Buffer Pool**

```sql
-- Check buffer pool size
SHOW VARIABLES LIKE 'innodb_buffer_pool_size';

-- Recommended: 50-70% of available RAM
-- For 4GB RAM server:
SET GLOBAL innodb_buffer_pool_size = 2147483648; -- 2GB

-- Add to /etc/mysql/my.cnf:
[mysqld]
innodb_buffer_pool_size = 2G
innodb_buffer_pool_instances = 4  -- For 2GB+ buffer pool
```

---

## 📊 Monitoring & Analysis

### **Enable Slow Query Log**

```sql
-- Enable slow query log
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;  -- Log queries > 1 second
SET GLOBAL log_queries_not_using_indexes = 'ON';

-- Check configuration
SHOW VARIABLES LIKE 'slow_query%';
SHOW VARIABLES LIKE 'long_query%';
```

**View slow queries:**
```bash
sudo tail -f /var/log/mysql/mysql-slow.log
```

### **Performance Schema**

```sql
-- Enable performance schema
SET GLOBAL performance_schema = ON;

-- View top slow queries
SELECT 
    DIGEST_TEXT,
    COUNT_STAR,
    AVG_TIMER_WAIT / 1000000000 AS avg_ms,
    SUM_TIMER_WAIT / 1000000000 AS total_ms
FROM performance_schema.events_statements_summary_by_digest
ORDER BY SUM_TIMER_WAIT DESC
LIMIT 10;
```

### **Table Statistics**

```sql
-- Table sizes
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb,
    table_rows
FROM information_schema.tables
WHERE table_schema = 'sirine_go'
ORDER BY (data_length + index_length) DESC;

-- Index sizes
SELECT 
    table_name,
    index_name,
    ROUND(stat_value * @@innodb_page_size / 1024 / 1024, 2) AS size_mb
FROM mysql.innodb_index_stats
WHERE database_name = 'sirine_go'
AND stat_name = 'size';
```

---

## 🔧 Database Maintenance

### **Regular Optimization**

```sql
-- Analyze table (update statistics)
ANALYZE TABLE users;

-- Optimize table (defragment & rebuild indexes)
OPTIMIZE TABLE users;

-- Check for corruption
CHECK TABLE users;

-- Repair table (if needed)
REPAIR TABLE users;
```

### **Automated Maintenance Script**

```bash
#!/bin/bash
# File: /var/scripts/mysql-maintenance.sh

MYSQL="mysql -u sirine_user -p${DB_PASSWORD} sirine_go"

# Optimize all tables
for table in users profiles notifications activity_logs user_achievements; do
    echo "Optimizing $table..."
    $MYSQL -e "OPTIMIZE TABLE $table;"
    $MYSQL -e "ANALYZE TABLE $table;"
done

echo "Maintenance complete: $(date)"
```

**Schedule with cron:**
```bash
# Run every Sunday at 3 AM
0 3 * * 0 /var/scripts/mysql-maintenance.sh >> /var/log/mysql-maintenance.log 2>&1
```

### **Archive Old Data**

```sql
-- Move old activity logs to archive table
CREATE TABLE activity_logs_archive LIKE activity_logs;

-- Move logs older than 6 months
INSERT INTO activity_logs_archive
SELECT * FROM activity_logs
WHERE created_at < DATE_SUB(NOW(), INTERVAL 6 MONTH);

-- Delete from main table
DELETE FROM activity_logs
WHERE created_at < DATE_SUB(NOW(), INTERVAL 6 MONTH);

-- Optimize table
OPTIMIZE TABLE activity_logs;
```

---

## ✅ Optimization Checklist

### Query Optimization
- [ ] All queries use indexes (check with EXPLAIN)
- [ ] No SELECT * queries
- [ ] LIMIT used for large result sets
- [ ] JOIN queries optimized with proper indexes
- [ ] No functions on indexed columns

### Index Strategy
- [ ] Foreign key columns indexed
- [ ] WHERE clause columns indexed
- [ ] ORDER BY columns indexed
- [ ] Composite indexes for multi-column queries
- [ ] Unused indexes removed

### Table Design
- [ ] Appropriate data types used
- [ ] Proper normalization
- [ ] Selective denormalization where beneficial
- [ ] Large tables partitioned (if > 1M rows)

### Configuration
- [ ] Connection pool optimized
- [ ] InnoDB buffer pool = 50-70% RAM
- [ ] Slow query log enabled
- [ ] Performance schema enabled
- [ ] Regular maintenance scheduled

---

## 📊 Performance Targets

### Query Performance
- Simple SELECT: < 10ms
- JOIN query: < 50ms
- Complex query: < 100ms
- INSERT/UPDATE: < 20ms

### Resource Usage
- Buffer pool hit ratio: > 99%
- Slow queries: < 1% of total
- Table locks: Minimal
- Connections: < 70% of max

### Maintenance
- Optimize tables: Weekly
- Analyze statistics: Weekly
- Archive old data: Monthly
- Backup database: Daily

---

## 📞 Support

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733

## 📖 Related Documentation

- [Database Management Guide](../database/management.md)
- [Database Models Guide](../database/models.md)
- [Backend Optimization](./backend-optimization.md)
- [Migrations Guide](../database/migrations.md)

---

**Last Updated:** 29 Desember 2025  
**Version:** 1.0.0
