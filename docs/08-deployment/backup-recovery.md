# 💾 Backup & Recovery - Sirine Go

Panduan lengkap untuk backup data dan disaster recovery procedures untuk aplikasi Sirine Go.

## 📋 Daftar Isi

1. [Database Backup](#database-backup)
2. [Application Backup](#application-backup)
3. [Automated Backup](#automated-backup)
4. [Backup Storage](#backup-storage)
5. [Recovery Procedures](#recovery-procedures)
6. [Disaster Recovery](#disaster-recovery-plan)

---

## 🗄️ Database Backup

### **Manual Database Backup**

```bash
# Full database backup
mysqldump -u sirine_user -p sirine_go > sirine_go_backup_$(date +%Y%m%d_%H%M%S).sql

# Backup with compression
mysqldump -u sirine_user -p sirine_go | gzip > sirine_go_backup_$(date +%Y%m%d_%H%M%S).sql.gz

# Backup specific tables
mysqldump -u sirine_user -p sirine_go users user_sessions > users_backup.sql

# Backup without data (schema only)
mysqldump -u sirine_user -p --no-data sirine_go > schema_backup.sql
```

### **Database Restore**

```bash
# Restore from backup
mysql -u sirine_user -p sirine_go < sirine_go_backup_20250129.sql

# Restore from compressed backup
gunzip < sirine_go_backup_20250129.sql.gz | mysql -u sirine_user -p sirine_go

# Restore specific tables
mysql -u sirine_user -p sirine_go < users_backup.sql
```

### **Verify Backup Integrity**

```bash
# Test restore to temporary database
mysql -u root -p -e "CREATE DATABASE sirine_go_test;"
mysql -u root -p sirine_go_test < sirine_go_backup_20250129.sql

# Check table count
mysql -u root -p -e "USE sirine_go_test; SHOW TABLES;" | wc -l

# Check record counts
mysql -u root -p sirine_go_test -e "
    SELECT 'users' as table_name, COUNT(*) as count FROM users
    UNION ALL
    SELECT 'notifications', COUNT(*) FROM notifications
    UNION ALL
    SELECT 'activity_logs', COUNT(*) FROM activity_logs;
"

# Drop test database
mysql -u root -p -e "DROP DATABASE sirine_go_test;"
```

---

## 📦 Application Backup

### **Backup Application Files**

```bash
#!/bin/bash
# File: /var/scripts/backup-app.sh

BACKUP_DIR="/var/backups/sirine-go"
APP_DIR="/var/www/sirine-go"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Backup application files (excluding node_modules, vendor, etc.)
tar -czf "$BACKUP_DIR/app_$DATE.tar.gz" \
    --exclude="node_modules" \
    --exclude="vendor" \
    --exclude=".git" \
    --exclude="frontend/dist" \
    -C /var/www sirine-go

echo "✅ Application backup completed: app_$DATE.tar.gz"
```

### **Backup Configuration Files**

```bash
#!/bin/bash
# File: /var/scripts/backup-config.sh

BACKUP_DIR="/var/backups/sirine-go/config"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p "$BACKUP_DIR"

# Backup important config files
tar -czf "$BACKUP_DIR/config_$DATE.tar.gz" \
    /var/www/sirine-go/backend/.env \
    /etc/nginx/sites-available/sirine-go \
    /etc/systemd/system/sirine-go.service \
    /etc/logrotate.d/sirine-go

echo "✅ Configuration backup completed: config_$DATE.tar.gz"
```

### **Backup Uploaded Files**

```bash
#!/bin/bash
# Backup user-uploaded files (photos, documents, etc.)

BACKUP_DIR="/var/backups/sirine-go/uploads"
UPLOAD_DIR="/var/www/sirine-go/backend/uploads"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p "$BACKUP_DIR"

# Backup with rsync (incremental)
rsync -avz --delete "$UPLOAD_DIR/" "$BACKUP_DIR/latest/"

# Create dated snapshot
tar -czf "$BACKUP_DIR/uploads_$DATE.tar.gz" -C "$BACKUP_DIR" latest

echo "✅ Uploads backup completed: uploads_$DATE.tar.gz"
```

---

## ⚙️ Automated Backup

### **Complete Backup Script**

```bash
#!/bin/bash
# File: /var/scripts/full-backup.sh

set -e  # Exit on error

BACKUP_ROOT="/var/backups/sirine-go"
DATE=$(date +%Y%m%d_%H%M%S)
LOG_FILE="/var/log/backup-sirine-go.log"

# Database credentials
DB_USER="sirine_user"
DB_PASS="YOUR_PASSWORD"  # Or load from .env
DB_NAME="sirine_go"

# Create backup directories
mkdir -p "$BACKUP_ROOT/database"
mkdir -p "$BACKUP_ROOT/application"
mkdir -p "$BACKUP_ROOT/config"
mkdir -p "$BACKUP_ROOT/uploads"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

log "========== Starting Full Backup =========="

# 1. Database Backup
log "Backing up database..."
mysqldump -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" | \
    gzip > "$BACKUP_ROOT/database/db_$DATE.sql.gz"
log "✅ Database backup completed"

# 2. Application Backup
log "Backing up application..."
tar -czf "$BACKUP_ROOT/application/app_$DATE.tar.gz" \
    --exclude="node_modules" \
    --exclude="vendor" \
    --exclude=".git" \
    --exclude="frontend/dist" \
    -C /var/www sirine-go
log "✅ Application backup completed"

# 3. Configuration Backup
log "Backing up configuration..."
tar -czf "$BACKUP_ROOT/config/config_$DATE.tar.gz" \
    /var/www/sirine-go/backend/.env \
    /etc/nginx/sites-available/sirine-go \
    /etc/systemd/system/sirine-go.service 2>/dev/null
log "✅ Configuration backup completed"

# 4. Uploads Backup
if [ -d "/var/www/sirine-go/backend/uploads" ]; then
    log "Backing up uploads..."
    tar -czf "$BACKUP_ROOT/uploads/uploads_$DATE.tar.gz" \
        -C /var/www/sirine-go/backend uploads
    log "✅ Uploads backup completed"
fi

# 5. Create manifest
log "Creating backup manifest..."
cat > "$BACKUP_ROOT/manifest_$DATE.txt" <<EOF
Sirine Go Backup Manifest
Date: $(date)
Hostname: $(hostname)

Files:
- database/db_$DATE.sql.gz
- application/app_$DATE.tar.gz
- config/config_$DATE.tar.gz
- uploads/uploads_$DATE.tar.gz

Sizes:
$(du -sh "$BACKUP_ROOT/database/db_$DATE.sql.gz")
$(du -sh "$BACKUP_ROOT/application/app_$DATE.tar.gz")
$(du -sh "$BACKUP_ROOT/config/config_$DATE.tar.gz")
$(du -sh "$BACKUP_ROOT/uploads/uploads_$DATE.tar.gz" 2>/dev/null || echo "No uploads")

Total Size: $(du -sh "$BACKUP_ROOT" | awk '{print $1}')
EOF

log "✅ Backup manifest created"

# 6. Cleanup old backups (keep last 7 days)
log "Cleaning up old backups..."
find "$BACKUP_ROOT/database" -name "db_*.sql.gz" -mtime +7 -delete
find "$BACKUP_ROOT/application" -name "app_*.tar.gz" -mtime +7 -delete
find "$BACKUP_ROOT/config" -name "config_*.tar.gz" -mtime +7 -delete
find "$BACKUP_ROOT/uploads" -name "uploads_*.tar.gz" -mtime +7 -delete
find "$BACKUP_ROOT" -name "manifest_*.txt" -mtime +7 -delete
log "✅ Old backups cleaned up"

log "========== Backup Completed Successfully =========="
log "Total backup size: $(du -sh $BACKUP_ROOT | awk '{print $1}')"
```

### **Schedule Automated Backups**

```bash
# Edit crontab
crontab -e
```

**Add backup schedule:**
```cron
# Full backup every day at 2 AM
0 2 * * * /var/scripts/full-backup.sh >> /var/log/backup-sirine-go.log 2>&1

# Database backup every 6 hours
0 */6 * * * mysqldump -u sirine_user -pPASSWORD sirine_go | gzip > /var/backups/sirine-go/database/db_$(date +\%Y\%m\%d_\%H\%M).sql.gz

# Weekly full backup with remote upload (Sunday 3 AM)
0 3 * * 0 /var/scripts/full-backup.sh && /var/scripts/upload-to-remote.sh
```

### **Email Notification on Backup**

```bash
# Add to end of backup script

# Send email notification
RECIPIENT="admin@example.com"
SUBJECT="Sirine Go Backup Report - $(date +%Y-%m-%d)"
BODY="Backup completed successfully at $(date)\n\nManifest:\n$(cat $BACKUP_ROOT/manifest_$DATE.txt)"

echo -e "$BODY" | mail -s "$SUBJECT" "$RECIPIENT"
```

---

## ☁️ Backup Storage

### **Local Storage Structure**

```
/var/backups/sirine-go/
├── database/
│   ├── db_20250129_020000.sql.gz
│   ├── db_20250128_020000.sql.gz
│   └── ...
├── application/
│   ├── app_20250129_020000.tar.gz
│   └── ...
├── config/
│   ├── config_20250129_020000.tar.gz
│   └── ...
├── uploads/
│   ├── uploads_20250129_020000.tar.gz
│   └── ...
└── manifest_20250129_020000.txt
```

### **Remote Backup (Cloud Storage)**

**Option 1: rsync to Remote Server**
```bash
#!/bin/bash
# File: /var/scripts/upload-to-remote.sh

REMOTE_USER="backup"
REMOTE_HOST="backup.example.com"
REMOTE_DIR="/backups/sirine-go"
LOCAL_DIR="/var/backups/sirine-go"

# Upload via rsync over SSH
rsync -avz --delete \
    -e "ssh -i /root/.ssh/backup_key" \
    "$LOCAL_DIR/" \
    "$REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/"

echo "✅ Backup uploaded to remote server"
```

**Option 2: AWS S3**
```bash
#!/bin/bash
# Install AWS CLI: sudo apt install awscli

# Configure AWS credentials first
# aws configure

BUCKET="s3://my-sirine-backups"
LOCAL_DIR="/var/backups/sirine-go"

# Upload to S3
aws s3 sync "$LOCAL_DIR" "$BUCKET/$(date +%Y/%m)/" \
    --storage-class STANDARD_IA

echo "✅ Backup uploaded to S3"
```

**Option 3: Google Drive (rclone)**
```bash
#!/bin/bash
# Install rclone: curl https://rclone.org/install.sh | sudo bash

# Configure rclone first
# rclone config

LOCAL_DIR="/var/backups/sirine-go"
REMOTE_NAME="gdrive"  # From rclone config

# Upload to Google Drive
rclone sync "$LOCAL_DIR" "$REMOTE_NAME:Sirine-Go-Backups/$(date +%Y/%m)/"

echo "✅ Backup uploaded to Google Drive"
```

---

## 🔄 Recovery Procedures

### **Full System Recovery**

```bash
#!/bin/bash
# File: /var/scripts/full-recovery.sh

BACKUP_DIR="/var/backups/sirine-go"
RESTORE_DATE="20250129_020000"  # Specify backup date

set -e

echo "========== Starting Full Recovery =========="
echo "⚠️ WARNING: This will overwrite existing data!"
read -p "Continue? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "Recovery cancelled"
    exit 1
fi

# 1. Stop application
echo "Stopping application..."
sudo systemctl stop sirine-go
sudo systemctl stop nginx

# 2. Restore database
echo "Restoring database..."
mysql -u sirine_user -p sirine_go < <(gunzip -c "$BACKUP_DIR/database/db_$RESTORE_DATE.sql.gz")

# 3. Restore application files
echo "Restoring application..."
cd /var/www
sudo tar -xzf "$BACKUP_DIR/application/app_$RESTORE_DATE.tar.gz"

# 4. Restore configuration
echo "Restoring configuration..."
sudo tar -xzf "$BACKUP_DIR/config/config_$RESTORE_DATE.tar.gz" -C /

# 5. Restore uploads
if [ -f "$BACKUP_DIR/uploads/uploads_$RESTORE_DATE.tar.gz" ]; then
    echo "Restoring uploads..."
    sudo tar -xzf "$BACKUP_DIR/uploads/uploads_$RESTORE_DATE.tar.gz" \
        -C /var/www/sirine-go/backend
fi

# 6. Fix permissions
echo "Fixing permissions..."
sudo chown -R www-data:www-data /var/www/sirine-go

# 7. Start services
echo "Starting services..."
sudo systemctl start sirine-go
sudo systemctl start nginx

# 8. Verify
echo "Verifying recovery..."
sleep 5
curl -s http://localhost:8080/health

echo "========== Recovery Completed =========="
```

### **Database-Only Recovery**

```bash
# Quick database restore
gunzip < /var/backups/sirine-go/database/db_20250129_020000.sql.gz | \
    mysql -u sirine_user -p sirine_go

# Verify restore
mysql -u sirine_user -p sirine_go -e "SHOW TABLES;"
```

### **Point-in-Time Recovery**

```bash
# Restore from specific backup time
BACKUP_TIME="20250129_020000"

# 1. Restore base backup
gunzip < /var/backups/sirine-go/database/db_$BACKUP_TIME.sql.gz | \
    mysql -u sirine_user -p sirine_go

# 2. Apply binary logs (if enabled)
# mysqlbinlog --start-datetime="2025-01-29 02:00:00" \
#     --stop-datetime="2025-01-29 12:00:00" \
#     /var/log/mysql/mysql-bin.* | \
#     mysql -u sirine_user -p sirine_go
```

---

## 🚨 Disaster Recovery Plan

### **Recovery Time Objective (RTO)**

| Component | Target RTO | Priority |
|-----------|------------|----------|
| Database | 30 minutes | Critical |
| Application | 1 hour | High |
| Uploads | 2 hours | Medium |

### **Recovery Point Objective (RPO)**

| Backup Type | Frequency | Max Data Loss |
|-------------|-----------|---------------|
| Database | Every 6 hours | 6 hours |
| Full backup | Daily | 24 hours |
| Remote backup | Weekly | 7 days |

### **Emergency Contact List**

```text
Primary Administrator:
- Name: Zulfikar Hidayatullah
- Phone: +62 857-1583-8733
- Email: admin@example.com

Hosting Provider:
- Provider: [Your Provider]
- Phone: [Support Number]
- Account: [Account ID]

Database Administrator:
- Name: [DBA Name]
- Phone: [DBA Phone]
```

### **Disaster Scenarios & Procedures**

**Scenario 1: Database Corruption**
```bash
1. Stop application
   sudo systemctl stop sirine-go

2. Restore from last backup
   gunzip < latest_backup.sql.gz | mysql -u sirine_user -p sirine_go

3. Verify data integrity
   mysql sirine_go -e "CHECK TABLE users, notifications, activity_logs;"

4. Start application
   sudo systemctl start sirine-go
```

**Scenario 2: Server Failure**
```bash
1. Provision new server
2. Install dependencies (Go, MySQL, Nginx)
3. Download backup from remote storage
4. Run full-recovery.sh script
5. Update DNS to new server IP
6. Verify application functionality
```

**Scenario 3: Accidental Data Deletion**
```bash
1. Identify deletion time
2. Restore database to point before deletion
3. Export affected records
4. Import to production database
5. Verify data integrity
```

---

## ✅ Backup Checklist

### Daily
- [ ] Verify automated backup completed
- [ ] Check backup log for errors
- [ ] Verify backup file sizes

### Weekly
- [ ] Test database restore to staging
- [ ] Upload backups to remote storage
- [ ] Review backup storage usage
- [ ] Verify backup integrity

### Monthly
- [ ] Full disaster recovery test
- [ ] Review and update recovery procedures
- [ ] Clean up old backups (>30 days)
- [ ] Test restore from remote backup

---

## 📊 Backup Statistics

### Monitor Backup Health

```bash
#!/bin/bash
# File: /var/scripts/backup-stats.sh

BACKUP_DIR="/var/backups/sirine-go"

echo "===== Backup Statistics ====="
echo ""
echo "Latest Backups:"
echo "Database: $(ls -lht $BACKUP_DIR/database/*.sql.gz | head -1 | awk '{print $6, $7, $8, $9}')"
echo "Application: $(ls -lht $BACKUP_DIR/application/*.tar.gz | head -1 | awk '{print $6, $7, $8, $9}')"
echo ""
echo "Backup Count (last 7 days):"
echo "Database: $(find $BACKUP_DIR/database -name '*.sql.gz' -mtime -7 | wc -l)"
echo "Application: $(find $BACKUP_DIR/application -name '*.tar.gz' -mtime -7 | wc -l)"
echo ""
echo "Total Backup Size:"
du -sh "$BACKUP_DIR"
```

---

## 📞 Support

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733

## 📖 Related Documentation

- [Production Deployment](./production-deployment.md)
- [Monitoring & Logging](./monitoring.md)
- [Database Management](../05-guides/database/management.md)
- [Troubleshooting](../09-troubleshooting/README.md)

---

**Last Updated:** 29 Desember 2025  
**Version:** 1.0.0
