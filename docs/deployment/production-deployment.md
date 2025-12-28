# 🚀 Panduan Deployment - Sirine Go App

Panduan untuk deploy aplikasi ke production server.

## 📋 Daftar Isi

1. [Prerequisites Production](#prerequisites-production)
2. [Setup Server](#setup-server)
3. [Deploy Aplikasi](#deploy-aplikasi)
4. [Setup Systemd Service](#setup-systemd-service)
5. [Setup Nginx](#setup-nginx)
6. [SSL/HTTPS](#ssl-https)
7. [Monitoring](#monitoring)
8. [Maintenance](#maintenance)

---

## 📋 Prerequisites Production

### **Server Requirements:**
- Ubuntu 20.04+ / Debian 11+ / CentOS 8+
- RAM: Minimal 1GB (Recommended 2GB+)
- Storage: Minimal 10GB
- CPU: 1 core (Recommended 2+ cores)

### **Software Requirements:**
- Go 1.24+
- MySQL 8.0+
- Nginx
- Git

---

## 🔧 Setup Server

### **1. Update System**

```bash
sudo apt update && sudo apt upgrade -y
```

### **2. Install Go**

```bash
# Download Go
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz

# Extract
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

### **3. Install MySQL**

```bash
# Install
sudo apt install mysql-server -y

# Secure installation
sudo mysql_secure_installation

# Start service
sudo systemctl start mysql
sudo systemctl enable mysql
```

### **4. Install Nginx**

```bash
sudo apt install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
```

### **5. Install Git**

```bash
sudo apt install git -y
```

---

## 🗄️ Setup Database Production

### **1. Login MySQL**

```bash
sudo mysql -u root -p
```

### **2. Create Database & User**

```sql
-- Create database
CREATE DATABASE sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create user
CREATE USER 'sirine_user'@'localhost' IDENTIFIED BY 'StrongPassword123!';

-- Grant privileges
GRANT ALL PRIVILEGES ON sirine_go.* TO 'sirine_user'@'localhost';
FLUSH PRIVILEGES;

-- Verify
SHOW GRANTS FOR 'sirine_user'@'localhost';

EXIT;
```

### **3. Test Connection**

```bash
mysql -u sirine_user -p sirine_go -e "SELECT 1;"
```

---

## 📦 Deploy Aplikasi

### **1. Create App Directory**

```bash
sudo mkdir -p /var/www/sirine-go
sudo chown $USER:$USER /var/www/sirine-go
```

### **2. Clone Repository**

```bash
cd /var/www/sirine-go
git clone <your-repo-url> .

# Atau upload via SCP/SFTP
```

### **3. Configure Environment**

```bash
nano backend/.env
```

**Production Configuration:**
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=sirine_user
DB_PASSWORD=StrongPassword123!
DB_NAME=sirine_go

# Server Configuration
SERVER_PORT=8080
GIN_MODE=release

# Timezone
TZ=Asia/Jakarta
```

### **4. Install Dependencies**

```bash
# Backend
cd backend
go mod download
go mod tidy

# Frontend (jika build di server)
cd ../frontend
npm install -g yarn
yarn install
```

### **5. Build Frontend**

```bash
cd frontend
yarn build
```

Output: `frontend/dist/` folder

### **6. Build Backend**

```bash
cd backend
go build -ldflags="-s -w" -o sirine-go cmd/server/main.go
```

Output: `backend/sirine-go` binary

### **7. Test Run**

```bash
cd backend
./sirine-go
```

Expected:
```
Database connected successfully
Server berjalan di port 8080
```

Ctrl+C untuk stop.

---

## 🔄 Setup Systemd Service

### **1. Create Service File**

```bash
sudo nano /etc/systemd/system/sirine-go.service
```

**Content:**
```ini
[Unit]
Description=Sirine Go App
After=network.target mysql.service
Wants=mysql.service

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

[Install]
WantedBy=multi-user.target
```

### **2. Create Log Directory**

```bash
sudo mkdir -p /var/log/sirine-go
sudo chown www-data:www-data /var/log/sirine-go
```

### **3. Set Permissions**

```bash
sudo chown -R www-data:www-data /var/www/sirine-go
sudo chmod +x /var/www/sirine-go/backend/sirine-go
```

### **4. Enable & Start Service**

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable auto-start
sudo systemctl enable sirine-go

# Start service
sudo systemctl start sirine-go

# Check status
sudo systemctl status sirine-go
```

Expected output:
```
● sirine-go.service - Sirine Go App
   Loaded: loaded (/etc/systemd/system/sirine-go.service; enabled)
   Active: active (running) since ...
```

### **5. View Logs**

```bash
# Real-time logs
sudo journalctl -u sirine-go -f

# Last 50 lines
sudo journalctl -u sirine-go -n 50

# App logs
sudo tail -f /var/log/sirine-go/app.log
```

---

## 🌐 Setup Nginx

### **1. Create Nginx Config**

```bash
sudo nano /etc/nginx/sites-available/sirine-go
```

**Content (HTTP Only):**
```nginx
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;

    # Logging
    access_log /var/log/nginx/sirine-go-access.log;
    error_log /var/log/nginx/sirine-go-error.log;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript 
               application/x-javascript application/xml+rss 
               application/json application/javascript;

    # Frontend static files
    location / {
        root /var/www/sirine-go/frontend/dist;
        try_files $uri $uri/ /index.html;
        
        # Cache static assets
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # API proxy to backend
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check
    location /health {
        proxy_pass http://localhost:8080;
        access_log off;
    }

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
}
```

### **2. Enable Site**

```bash
# Create symbolic link
sudo ln -s /etc/nginx/sites-available/sirine-go /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

### **3. Test**

```bash
curl http://your-domain.com/health
```

Expected: `{"status":"ok","message":"Server berjalan dengan baik"}`

---

## 🔒 SSL/HTTPS Setup

### **1. Install Certbot**

```bash
sudo apt install certbot python3-certbot-nginx -y
```

### **2. Obtain SSL Certificate**

```bash
sudo certbot --nginx -d your-domain.com -d www.your-domain.com
```

Follow prompts:
1. Enter email
2. Agree to terms
3. Choose redirect option (2 - Redirect HTTP to HTTPS)

### **3. Verify SSL**

```bash
curl https://your-domain.com/health
```

### **4. Auto-Renewal**

```bash
# Test renewal
sudo certbot renew --dry-run

# Certbot auto-renewal is enabled by default
sudo systemctl status certbot.timer
```

### **5. Updated Nginx Config (HTTPS)**

Certbot will auto-update config. Verify:

```bash
sudo nano /etc/nginx/sites-available/sirine-go
```

Should have:
```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com www.your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    
    # ... rest of config
}

server {
    listen 80;
    server_name your-domain.com www.your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

---

## 🔥 Firewall Setup

### **1. Install UFW**

```bash
sudo apt install ufw -y
```

### **2. Configure Rules**

```bash
# Allow SSH (IMPORTANT!)
sudo ufw allow 22/tcp

# Allow HTTP
sudo ufw allow 80/tcp

# Allow HTTPS
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status
```

---

## 📊 Monitoring

### **1. Check Service Status**

```bash
# App status
sudo systemctl status sirine-go

# Nginx status
sudo systemctl status nginx

# MySQL status
sudo systemctl status mysql
```

### **2. View Logs**

```bash
# App logs
sudo tail -f /var/log/sirine-go/app.log
sudo tail -f /var/log/sirine-go/error.log

# Nginx logs
sudo tail -f /var/log/nginx/sirine-go-access.log
sudo tail -f /var/log/nginx/sirine-go-error.log

# System logs
sudo journalctl -u sirine-go -f
```

### **3. Resource Monitoring**

```bash
# CPU & Memory
htop

# Disk usage
df -h

# Network connections
sudo netstat -tulpn | grep :8080
```

---

## 🔄 Maintenance

### **1. Update Aplikasi**

```bash
cd /var/www/sirine-go

# Pull latest code
git pull origin main

# Rebuild frontend
cd frontend
yarn install
yarn build

# Rebuild backend
cd ../backend
go build -ldflags="-s -w" -o sirine-go cmd/server/main.go

# Restart service
sudo systemctl restart sirine-go

# Check status
sudo systemctl status sirine-go
```

### **2. Database Backup**

**Manual Backup:**
```bash
mysqldump -u sirine_user -p sirine_go > backup_$(date +%Y%m%d_%H%M%S).sql
```

**Automated Backup (Cron):**
```bash
# Create backup directory
sudo mkdir -p /var/backups/sirine-go
sudo chown $USER:$USER /var/backups/sirine-go

# Edit crontab
crontab -e
```

Add:
```cron
# Backup database every day at 2 AM
0 2 * * * mysqldump -u sirine_user -pYOUR_PASSWORD sirine_go > /var/backups/sirine-go/db_$(date +\%Y\%m\%d).sql

# Delete backups older than 7 days
0 3 * * * find /var/backups/sirine-go -name "db_*.sql" -mtime +7 -delete
```

### **3. Log Rotation**

```bash
sudo nano /etc/logrotate.d/sirine-go
```

Content:
```
/var/log/sirine-go/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 www-data www-data
    sharedscripts
    postrotate
        systemctl reload sirine-go > /dev/null 2>&1 || true
    endscript
}
```

---

## 🐛 Production Troubleshooting

### **Service Won't Start**

```bash
# Check logs
sudo journalctl -u sirine-go -n 50

# Check permissions
ls -la /var/www/sirine-go/backend/sirine-go

# Check port
sudo netstat -tulpn | grep :8080
```

### **502 Bad Gateway**

```bash
# Check app running
sudo systemctl status sirine-go

# Check Nginx config
sudo nginx -t

# Check logs
sudo tail -f /var/log/nginx/error.log
```

### **Database Connection Error**

```bash
# Check MySQL running
sudo systemctl status mysql

# Test connection
mysql -u sirine_user -p sirine_go -e "SELECT 1;"

# Check .env
cat /var/www/sirine-go/backend/.env
```

---

## ✅ Deployment Checklist

- [ ] Server setup complete
- [ ] Go, MySQL, Nginx installed
- [ ] Database created
- [ ] User & permissions configured
- [ ] Application deployed
- [ ] Environment configured
- [ ] Frontend built
- [ ] Backend built
- [ ] Systemd service created
- [ ] Service enabled & running
- [ ] Nginx configured
- [ ] SSL certificate installed
- [ ] Firewall configured
- [ ] Logs accessible
- [ ] Monitoring setup
- [ ] Backup automation configured
- [ ] Test all endpoints
- [ ] Test PWA offline mode

---

## 📚 Related Documentation

**Before deployment:**
- **[SETUP_GUIDE.md](./SETUP_GUIDE.md)** - Local setup
- **[TESTING.md](./TESTING.md)** - Test before deploy
- **[CHECKLIST.md](./CHECKLIST.md)** - Verify everything works

**After deployment:**
- **[FAQ.md](./FAQ.md)** - Deployment troubleshooting
- **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)** - API reference

---

## 📞 Support

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733

**Deployment issues?** Check [FAQ.md](./FAQ.md) → Deployment section

---

**Last Updated:** 27 Desember 2025  
**Version:** 1.0.0
