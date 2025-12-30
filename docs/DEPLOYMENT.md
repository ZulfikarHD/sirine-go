# Deployment Guide - Sirine Go

**Developer**: Zulfikar Hidayatullah  
**Version**: 1.0.0  
**Last Updated**: 2025-01-29

Panduan lengkap untuk deployment Sistem Produksi Pita Cukai ke production environment.

---

## Table of Contents

- [Pre-Deployment Checklist](#pre-deployment-checklist)
- [Server Requirements](#server-requirements)
- [Database Setup](#database-setup)
- [Backend Deployment](#backend-deployment)
- [Frontend Deployment](#frontend-deployment)
- [Nginx Configuration](#nginx-configuration)
- [SSL/HTTPS Setup](#sslhttps-setup)
- [Monitoring & Logging](#monitoring--logging)
- [Backup Strategy](#backup-strategy)
- [Troubleshooting](#troubleshooting)

---

## Pre-Deployment Checklist

### Code Readiness
- [ ] Semua tests passing (unit, integration, E2E)
- [ ] Code review completed
- [ ] No critical linter errors
- [ ] Performance benchmarks met (login < 500ms, dashboard < 1s)
- [ ] Security audit passed
- [ ] All TODOs resolved atau documented

### Environment Configuration
- [ ] Production `.env` files prepared
- [ ] Strong JWT_SECRET generated
- [ ] Strong database password set
- [ ] SMTP email configured dan tested
- [ ] CORS origins set ke production domain
- [ ] Rate limiting configured
- [ ] Security headers enabled

### Database
- [ ] Backup current database (jika existing)
- [ ] Migration scripts tested
- [ ] Seed data prepared
- [ ] Indexes optimized
- [ ] Database user permissions configured

### Infrastructure
- [ ] Server provisioned (minimal specs met)
- [ ] Domain name configured
- [ ] SSL certificate obtained
- [ ] Firewall rules configured
- [ ] Monitoring tools setup

---

## Server Requirements

### Minimum Specifications

**Backend Server:**
- CPU: 2 cores
- RAM: 4GB
- Storage: 20GB SSD
- OS: Ubuntu 20.04 LTS atau newer

**Database Server:** (dapat same server untuk small-medium deployment)
- CPU: 2 cores
- RAM: 4GB
- Storage: 50GB SSD (adjust based on data volume)

**Frontend:** (dapat serve dari backend server atau CDN)
- Static files hosting (Nginx)
- CDN recommended untuk better performance

### Software Requirements

**Backend:**
- Go 1.21 atau newer
- MySQL 8.0 atau newer
- Nginx (sebagai reverse proxy)
- Systemd (untuk service management)

**Frontend:**
- Node.js 18+ dan Yarn (untuk build)
- Nginx (untuk serving static files)

---

## Database Setup

### 1. Install MySQL

```bash
# Update package list
sudo apt update

# Install MySQL Server
sudo apt install mysql-server -y

# Secure installation
sudo mysql_secure_installation
```

### 2. Create Database dan User

```sql
-- Login ke MySQL sebagai root
sudo mysql -u root -p

-- Create database
CREATE DATABASE sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create dedicated user
CREATE USER 'sirine_user'@'localhost' IDENTIFIED BY 'YOUR_STRONG_PASSWORD';

-- Grant permissions
GRANT ALL PRIVILEGES ON sirine_go.* TO 'sirine_user'@'localhost';
FLUSH PRIVILEGES;

EXIT;
```

### 3. Run Database Schema

```bash
# Run setup script
mysql -u sirine_user -p sirine_go < backend/database/setup.sql

# Run optimization script
mysql -u sirine_user -p sirine_go < backend/database/optimization.sql
```

### 4. Verify Database Setup

```bash
# Check tables created
mysql -u sirine_user -p sirine_go -e "SHOW TABLES;"

# Check indexes
mysql -u sirine_user -p sirine_go -e "
SELECT TABLE_NAME, INDEX_NAME, COLUMN_NAME 
FROM information_schema.STATISTICS 
WHERE TABLE_SCHEMA = 'sirine_go' 
ORDER BY TABLE_NAME, INDEX_NAME;
"
```

### 5. Create Initial Admin User

```sql
-- Login ke database
mysql -u sirine_user -p sirine_go

-- Insert admin user (password: Admin@123)
INSERT INTO users (nip, full_name, email, phone, password_hash, role, department, status, must_change_password) 
VALUES (
    'admin', 
    'Administrator', 
    'admin@sirine-go.com',
    '08123456789',
    '$2a$12$LqTYjX.OO4dBjPbDxz0Nse7JxYZxYxYxYxYxYxYxYxYxYxYxYxY', -- Hash dari Admin@123
    'ADMIN', 
    'KHAZWAL', 
    'ACTIVE',
    TRUE
);
```

---

## Backend Deployment

### 1. Prepare Server

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
cd /tmp
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# Set Go path
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

### 2. Deploy Application

```bash
# Create app directory
sudo mkdir -p /var/www/sirine-go/backend
sudo chown $USER:$USER /var/www/sirine-go

# Clone atau copy source code
cd /var/www/sirine-go/backend
# (Upload your code here via git, scp, rsync, etc.)

# Copy environment file
cp env.example.txt .env

# Edit .env dengan production values
nano .env
```

**Important `.env` changes untuk production:**

```env
APP_ENV=production
APP_PORT=8080
APP_HOST=127.0.0.1
FRONTEND_URL=https://yourdomain.com

DB_HOST=localhost
DB_USER=sirine_user
DB_PASSWORD=YOUR_STRONG_PASSWORD
DB_NAME=sirine_go

JWT_SECRET=YOUR_STRONG_JWT_SECRET_GENERATE_WITH_OPENSSL
JWT_EXPIRY=15
REFRESH_TOKEN_EXPIRY=30

BCRYPT_COST=12
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=15

# Email settings
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password

CORS_ALLOWED_ORIGINS=https://yourdomain.com
DEBUG=false
LOG_LEVEL=warn
```

### 3. Build Application

```bash
# Install dependencies
cd /var/www/sirine-go/backend
go mod download

# Build binary
go build -o sirine-go cmd/main.go

# Test run
./sirine-go
# Press Ctrl+C setelah verify no errors
```

### 4. Create Systemd Service

```bash
# Create service file
sudo nano /etc/systemd/system/sirine-go.service
```

**Service configuration:**

```ini
[Unit]
Description=Sirine Go Backend API
After=network.target mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/var/www/sirine-go/backend
ExecStart=/var/www/sirine-go/backend/sirine-go
Restart=always
RestartSec=5
StandardOutput=append:/var/log/sirine-go/app.log
StandardError=append:/var/log/sirine-go/error.log

# Environment
Environment="GIN_MODE=release"

# Security
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

### 5. Start Service

```bash
# Create log directory
sudo mkdir -p /var/log/sirine-go
sudo chown www-data:www-data /var/log/sirine-go

# Reload systemd
sudo systemctl daemon-reload

# Enable service (auto-start on boot)
sudo systemctl enable sirine-go

# Start service
sudo systemctl start sirine-go

# Check status
sudo systemctl status sirine-go

# View logs
sudo journalctl -u sirine-go -f
```

---

## Frontend Deployment

### 1. Build Production Bundle

**Di local machine (atau CI/CD):**

```bash
cd frontend

# Install dependencies
yarn install

# Copy environment file
cp env.example.txt .env

# Edit .env untuk production
nano .env
```

**Production `.env`:**

```env
VITE_API_BASE_URL=https://api.yourdomain.com/api
VITE_APP_NAME=Sistem Produksi Pita Cukai
VITE_APP_ENV=production
VITE_TIMEZONE=Asia/Jakarta

VITE_ENABLE_PWA=true
VITE_ENABLE_OFFLINE=true
VITE_ENABLE_DEVTOOLS=false
VITE_DEBUG=false
```

```bash
# Build for production
yarn build

# Hasil build ada di folder dist/
# Upload folder dist/ ke server
```

### 2. Deploy ke Server

```bash
# Di server, create directory untuk frontend
sudo mkdir -p /var/www/sirine-go/frontend
sudo chown www-data:www-data /var/www/sirine-go/frontend

# Upload dist/ folder
# (gunakan scp, rsync, atau git)
# Example dengan rsync:
rsync -avz dist/ user@server:/var/www/sirine-go/frontend/
```

---

## Nginx Configuration

### 1. Install Nginx

```bash
sudo apt install nginx -y
sudo systemctl enable nginx
sudo systemctl start nginx
```

### 2. Configure Server Block

```bash
sudo nano /etc/nginx/sites-available/sirine-go
```

**Nginx configuration:**

```nginx
# Redirect HTTP ke HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name yourdomain.com www.yourdomain.com;
    
    # Redirect semua HTTP requests ke HTTPS
    return 301 https://$server_name$request_uri;
}

# HTTPS server
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Root directory untuk frontend
    root /var/www/sirine-go/frontend;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/json application/javascript application/xml+rss application/rss+xml font/truetype font/opentype application/vnd.ms-fontobject image/svg+xml;

    # Frontend - serve static files
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API proxy ke backend
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # Buffering
        proxy_buffering off;
    }

    # Cache static assets
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Service worker (no cache)
    location /sw.js {
        add_header Cache-Control "no-cache";
        proxy_cache_bypass $http_pragma;
        proxy_cache_revalidate on;
        expires off;
        access_log off;
    }

    # Manifest (no cache)
    location /manifest.webmanifest {
        add_header Cache-Control "no-cache";
        expires off;
    }
}
```

### 3. Enable Site

```bash
# Create symbolic link
sudo ln -s /etc/nginx/sites-available/sirine-go /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

---

## SSL/HTTPS Setup

### Option 1: Let's Encrypt (Free, Recommended)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain certificate
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Test auto-renewal
sudo certbot renew --dry-run

# Certificate akan auto-renew via cron
```

### Option 2: Manual SSL Certificate

Jika menggunakan SSL certificate dari provider lain:

```bash
# Copy certificate files ke server
sudo mkdir -p /etc/ssl/sirine-go
sudo cp your-certificate.crt /etc/ssl/sirine-go/
sudo cp your-private-key.key /etc/ssl/sirine-go/
sudo chmod 600 /etc/ssl/sirine-go/your-private-key.key

# Update Nginx config dengan path ke certificates
```

---

## Monitoring & Logging

### 1. Application Logs

```bash
# View backend logs
sudo journalctl -u sirine-go -f

# View specific log file
sudo tail -f /var/log/sirine-go/app.log
sudo tail -f /var/log/sirine-go/error.log

# View Nginx access logs
sudo tail -f /var/nginx/access.log

# View Nginx error logs
sudo tail -f /var/nginx/error.log
```

### 2. Database Monitoring

```sql
-- Check active connections
SHOW PROCESSLIST;

-- Check slow queries
SELECT * FROM mysql.slow_log ORDER BY start_time DESC LIMIT 10;

-- Check table sizes
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
FROM information_schema.TABLES
WHERE table_schema = 'sirine_go'
ORDER BY (data_length + index_length) DESC;
```

### 3. System Monitoring

```bash
# CPU dan Memory usage
htop

# Disk usage
df -h

# Network usage
nethogs

# Check backend service health
curl http://localhost:8080/api/health

# Check frontend
curl https://yourdomain.com
```

### 4. Setup Automated Monitoring (Optional)

**Install monitoring tools:**

```bash
# Install Prometheus (for metrics)
# Install Grafana (for dashboards)
# Install Loki (for log aggregation)
# Configure alerts via email/Slack
```

---

## Backup Strategy

### 1. Database Backup

**Create backup script:**

```bash
# Create backup directory
sudo mkdir -p /var/backups/sirine-go
sudo chmod 700 /var/backups/sirine-go

# Create backup script
sudo nano /usr/local/bin/backup-sirine-db.sh
```

**Backup script content:**

```bash
#!/bin/bash

BACKUP_DIR="/var/backups/sirine-go"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="sirine_go"
DB_USER="sirine_user"
DB_PASS="YOUR_DB_PASSWORD"

# Create backup
mysqldump -u $DB_USER -p$DB_PASS $DB_NAME | gzip > $BACKUP_DIR/sirine_go_$DATE.sql.gz

# Delete backups older than 30 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

# Log
echo "Backup completed: sirine_go_$DATE.sql.gz"
```

```bash
# Make executable
sudo chmod +x /usr/local/bin/backup-sirine-db.sh

# Test backup
sudo /usr/local/bin/backup-sirine-db.sh
```

**Schedule daily backup via cron:**

```bash
# Edit crontab
sudo crontab -e

# Add daily backup at 2 AM
0 2 * * * /usr/local/bin/backup-sirine-db.sh >> /var/log/sirine-go/backup.log 2>&1
```

### 2. Application Backup

```bash
# Backup entire application directory
sudo tar -czf /var/backups/sirine-go/app_$(date +%Y%m%d).tar.gz /var/www/sirine-go
```

### 3. Restore dari Backup

```bash
# Restore database
gunzip < /var/backups/sirine-go/sirine_go_20250129_020000.sql.gz | mysql -u sirine_user -p sirine_go

# Restore application
sudo tar -xzf /var/backups/sirine-go/app_20250129.tar.gz -C /
```

---

## Troubleshooting

### Backend Issues

**Service tidak start:**

```bash
# Check logs
sudo journalctl -u sirine-go -n 50

# Check process
ps aux | grep sirine-go

# Check port
sudo netstat -tlnp | grep 8080

# Test manually
cd /var/www/sirine-go/backend
./sirine-go
```

**Database connection error:**

```bash
# Test database connection
mysql -u sirine_user -p sirine_go

# Check .env configuration
cat /var/www/sirine-go/backend/.env | grep DB_

# Check MySQL service
sudo systemctl status mysql
```

### Frontend Issues

**404 errors:**

```bash
# Check files exist
ls -la /var/www/sirine-go/frontend/

# Check Nginx config
sudo nginx -t

# Check Nginx logs
sudo tail -f /var/log/nginx/error.log
```

**API calls failing:**

```bash
# Check CORS settings di backend .env
# Check Nginx proxy configuration
# Check backend service running
sudo systemctl status sirine-go

# Test API directly
curl http://localhost:8080/api/health
```

### Performance Issues

**Slow queries:**

```sql
-- Enable slow query log
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;

-- Check slow queries
SELECT * FROM mysql.slow_log ORDER BY query_time DESC LIMIT 10;
```

**High memory usage:**

```bash
# Check memory
free -m

# Check processes
top -o %MEM

# Restart backend service jika perlu
sudo systemctl restart sirine-go
```

### SSL Issues

**Certificate errors:**

```bash
# Test SSL
openssl s_client -connect yourdomain.com:443

# Check certificate expiry
sudo certbot certificates

# Renew certificate manually
sudo certbot renew
```

---

## Post-Deployment Verification

### Checklist

- [ ] Backend API accessible via `https://yourdomain.com/api/health`
- [ ] Frontend loads correctly di `https://yourdomain.com`
- [ ] Login works dengan admin credentials
- [ ] HTTPS enforced (HTTP redirects ke HTTPS)
- [ ] SSL certificate valid
- [ ] PWA installable pada mobile devices
- [ ] Database queries performing well (< 100ms average)
- [ ] Logs being written correctly
- [ ] Backups running automatically
- [ ] Email notifications working
- [ ] Rate limiting protecting endpoints
- [ ] Security headers present (check di browser DevTools)

### Testing

```bash
# Test API health
curl https://yourdomain.com/api/health

# Test login
curl -X POST https://yourdomain.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"nip":"99999","password":"Admin@123"}'

# Test rate limiting
for i in {1..10}; do curl https://yourdomain.com/api/auth/login; done

# Check security headers
curl -I https://yourdomain.com
```

---

## Maintenance

### Regular Tasks

**Daily:**
- Monitor application logs
- Check system resources (CPU, RAM, Disk)
- Verify backups completed

**Weekly:**
- Review slow query logs
- Check security logs
- Update dependencies jika ada security patches

**Monthly:**
- Run database OPTIMIZE TABLE
- Clean up old logs (> 30 days)
- Review and test backup restore procedure

**Quarterly:**
- Security audit
- Performance review
- Capacity planning review

---

## Support & Contact

**Developer**: Zulfikar Hidayatullah  
**Phone**: +62 857-1583-8733  
**Email**: (add your email here)

**Documentation**: `/docs` directory  
**Issues**: (add issue tracker URL if applicable)

---

**Deployment Date**: _____________  
**Deployed By**: _____________  
**Production URL**: _____________  
**Database Version**: _____________  
**Backend Version**: _____________  
**Frontend Version**: _____________
