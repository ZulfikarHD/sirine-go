# 📊 Monitoring & Logging - Sirine Go

Panduan monitoring aplikasi Sirine Go di production environment untuk memastikan sistem berjalan optimal dan identifikasi masalah lebih cepat.

## 📋 Daftar Isi

1. [Application Monitoring](#application-monitoring)
2. [System Monitoring](#system-monitoring)
3. [Database Monitoring](#database-monitoring)
4. [Log Management](#log-management)
5. [Alert Configuration](#alert-configuration)
6. [Performance Metrics](#performance-metrics)

---

## 📱 Application Monitoring

### **Service Status Monitoring**

```bash
# Check application status
sudo systemctl status sirine-go

# View last 50 log lines
sudo journalctl -u sirine-go -n 50

# Follow logs in real-time
sudo journalctl -u sirine-go -f

# View logs from specific time
sudo journalctl -u sirine-go --since "2025-01-01 00:00:00"
```

### **Health Check Endpoint**

```go
// File: backend/routes/routes.go

func SetupRouter(db *gorm.DB) *gin.Engine {
    router := gin.Default()
    
    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        // Check database connection
        sqlDB, err := db.DB()
        if err != nil {
            c.JSON(500, gin.H{
                "status": "error",
                "message": "Database connection failed",
            })
            return
        }
        
        if err := sqlDB.Ping(); err != nil {
            c.JSON(500, gin.H{
                "status": "error",
                "message": "Database ping failed",
            })
            return
        }
        
        c.JSON(200, gin.H{
            "status": "ok",
            "message": "Server berjalan dengan baik",
            "timestamp": time.Now(),
        })
    })
    
    return router
}
```

### **Monitoring Script**

```bash
#!/bin/bash
# File: /var/scripts/monitor-sirine-go.sh

# Check if service is running
if ! systemctl is-active --quiet sirine-go; then
    echo "❌ Service is DOWN at $(date)"
    # Send alert (email/SMS/webhook)
    systemctl restart sirine-go
    exit 1
fi

# Check health endpoint
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
if [ "$HTTP_CODE" != "200" ]; then
    echo "⚠️ Health check failed with code $HTTP_CODE at $(date)"
    exit 1
fi

echo "✅ Service is healthy at $(date)"
```

**Schedule with cron:**
```bash
# Check every 5 minutes
*/5 * * * * /var/scripts/monitor-sirine-go.sh >> /var/log/sirine-go-monitor.log 2>&1
```

### **Response Time Monitoring**

```go
// File: backend/middleware/metrics.go

func ResponseTimeMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        path := c.Request.URL.Path
        status := c.Writer.Status()
        
        // Log slow requests (>500ms)
        if duration > 500*time.Millisecond {
            log.Printf("SLOW REQUEST: %s %s - %d - %v\n", 
                c.Request.Method, path, status, duration)
        }
        
        // Log to metrics (optional)
        // metrics.RecordResponseTime(path, duration)
    }
}
```

---

## 🖥️ System Monitoring

### **Resource Usage Monitoring**

```bash
# CPU & Memory usage
htop

# Or using top
top -p $(pgrep -f sirine-go)

# Disk usage
df -h

# Specific directory
du -sh /var/www/sirine-go
du -sh /var/log/sirine-go

# Network usage
sudo nethogs

# Check open connections
sudo netstat -tulpn | grep :8080
```

### **Automated System Monitoring**

```bash
#!/bin/bash
# File: /var/scripts/system-monitor.sh

THRESHOLD_CPU=80
THRESHOLD_MEM=80
THRESHOLD_DISK=85

# Check CPU usage
CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
if (( $(echo "$CPU_USAGE > $THRESHOLD_CPU" | bc -l) )); then
    echo "⚠️ HIGH CPU: ${CPU_USAGE}% at $(date)"
fi

# Check memory usage
MEM_USAGE=$(free | grep Mem | awk '{printf("%.0f", $3/$2 * 100)}')
if [ "$MEM_USAGE" -gt "$THRESHOLD_MEM" ]; then
    echo "⚠️ HIGH MEMORY: ${MEM_USAGE}% at $(date)"
fi

# Check disk usage
DISK_USAGE=$(df -h / | awk 'NR==2 {print $5}' | cut -d'%' -f1)
if [ "$DISK_USAGE" -gt "$THRESHOLD_DISK" ]; then
    echo "⚠️ HIGH DISK: ${DISK_USAGE}% at $(date)"
fi
```

**Schedule with cron:**
```bash
# Check every 15 minutes
*/15 * * * * /var/scripts/system-monitor.sh >> /var/log/system-monitor.log 2>&1
```

### **Process Monitoring**

```bash
# Check process details
ps aux | grep sirine-go

# Memory usage of process
pmap -x $(pgrep -f sirine-go)

# Thread count
ps -Lf $(pgrep -f sirine-go) | wc -l

# Open files
sudo lsof -p $(pgrep -f sirine-go)
```

---

## 🗄️ Database Monitoring

### **Connection Monitoring**

```sql
-- Check active connections
SHOW PROCESSLIST;

-- Connection statistics
SHOW STATUS LIKE 'Threads_connected';
SHOW STATUS LIKE 'Threads_running';
SHOW STATUS LIKE 'Max_used_connections';

-- Detailed connection info
SELECT 
    ID,
    USER,
    HOST,
    DB,
    COMMAND,
    TIME,
    STATE,
    INFO
FROM information_schema.PROCESSLIST
WHERE USER = 'sirine_user';
```

### **Query Performance**

```sql
-- Enable slow query log
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;

-- View slow queries
SELECT 
    DIGEST_TEXT,
    COUNT_STAR AS executions,
    AVG_TIMER_WAIT / 1000000000 AS avg_ms,
    MAX_TIMER_WAIT / 1000000000 AS max_ms
FROM performance_schema.events_statements_summary_by_digest
ORDER BY AVG_TIMER_WAIT DESC
LIMIT 10;
```

### **Database Size Monitoring**

```sql
-- Database size
SELECT 
    table_schema AS 'Database',
    ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
FROM information_schema.tables
WHERE table_schema = 'sirine_go'
GROUP BY table_schema;

-- Table sizes
SELECT 
    table_name AS 'Table',
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)',
    table_rows AS 'Rows'
FROM information_schema.tables
WHERE table_schema = 'sirine_go'
ORDER BY (data_length + index_length) DESC;
```

### **InnoDB Buffer Pool**

```sql
-- Buffer pool usage
SHOW STATUS LIKE 'Innodb_buffer_pool_%';

-- Hit ratio (should be > 99%)
SELECT 
    CONCAT(
        ROUND(
            (1 - (VARIABLE_VALUE / (SELECT VARIABLE_VALUE 
                FROM information_schema.GLOBAL_STATUS 
                WHERE VARIABLE_NAME='Innodb_buffer_pool_read_requests'))) * 100, 
            2
        ), 
        '%'
    ) AS buffer_pool_hit_ratio
FROM information_schema.GLOBAL_STATUS
WHERE VARIABLE_NAME='Innodb_buffer_pool_reads';
```

---

## 📝 Log Management

### **Application Logs**

**Log Locations:**
```bash
# Application logs (via systemd)
/var/log/sirine-go/app.log
/var/log/sirine-go/error.log

# Systemd journal
sudo journalctl -u sirine-go

# Nginx logs
/var/log/nginx/sirine-go-access.log
/var/log/nginx/sirine-go-error.log
```

### **Log Rotation Configuration**

```bash
# File: /etc/logrotate.d/sirine-go

/var/log/sirine-go/*.log {
    daily                    # Rotate daily
    rotate 14                # Keep 14 days
    compress                 # Compress old logs
    delaycompress            # Delay compression for 1 cycle
    notifempty               # Don't rotate if empty
    create 0640 www-data www-data
    sharedscripts
    postrotate
        systemctl reload sirine-go > /dev/null 2>&1 || true
    endscript
}

/var/log/nginx/sirine-go-*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 www-data adm
    sharedscripts
    postrotate
        systemctl reload nginx > /dev/null 2>&1 || true
    endscript
}
```

**Test log rotation:**
```bash
sudo logrotate -f /etc/logrotate.d/sirine-go
```

### **Log Analysis**

```bash
# Count errors
grep -i "error" /var/log/sirine-go/app.log | wc -l

# Find specific error
grep -i "database connection" /var/log/sirine-go/error.log

# Top error messages
awk '/ERROR/ {print $0}' /var/log/sirine-go/error.log | \
    sort | uniq -c | sort -rn | head -10

# Request by status code (Nginx)
awk '{print $9}' /var/log/nginx/sirine-go-access.log | \
    sort | uniq -c | sort -rn

# Top requested endpoints
awk '{print $7}' /var/log/nginx/sirine-go-access.log | \
    sort | uniq -c | sort -rn | head -10

# Average response time
awk '{sum+=$NF; count++} END {print sum/count}' \
    /var/log/nginx/sirine-go-access.log
```

### **Structured Logging (Go)**

```go
// File: backend/utils/logger.go

package utils

import (
    "encoding/json"
    "log"
    "os"
)

type LogEntry struct {
    Level     string      `json:"level"`
    Message   string      `json:"message"`
    Timestamp string      `json:"timestamp"`
    Context   interface{} `json:"context,omitempty"`
}

func LogInfo(message string, context interface{}) {
    entry := LogEntry{
        Level:     "INFO",
        Message:   message,
        Timestamp: time.Now().Format(time.RFC3339),
        Context:   context,
    }
    logJSON(entry)
}

func LogError(message string, context interface{}) {
    entry := LogEntry{
        Level:     "ERROR",
        Message:   message,
        Timestamp: time.Now().Format(time.RFC3339),
        Context:   context,
    }
    logJSON(entry)
}

func logJSON(entry LogEntry) {
    jsonData, _ := json.Marshal(entry)
    log.Println(string(jsonData))
}
```

---

## 🚨 Alert Configuration

### **Email Alerts**

```bash
#!/bin/bash
# File: /var/scripts/send-alert.sh

RECIPIENT="admin@example.com"
SUBJECT="Sirine Go Alert"
MESSAGE=$1

echo "$MESSAGE" | mail -s "$SUBJECT" "$RECIPIENT"
```

**Install mail utility:**
```bash
sudo apt install mailutils -y
```

### **Webhook Alerts (Slack/Discord)**

```bash
#!/bin/bash
# File: /var/scripts/send-webhook.sh

WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
MESSAGE=$1

curl -X POST "$WEBHOOK_URL" \
    -H 'Content-Type: application/json' \
    -d "{\"text\": \"⚠️ Sirine Go Alert: $MESSAGE\"}"
```

### **Alert Conditions**

```bash
#!/bin/bash
# File: /var/scripts/alert-monitor.sh

# Service down
if ! systemctl is-active --quiet sirine-go; then
    /var/scripts/send-webhook.sh "Service is DOWN!"
fi

# High error rate
ERROR_COUNT=$(grep -c "ERROR" /var/log/sirine-go/error.log | tail -100)
if [ "$ERROR_COUNT" -gt 50 ]; then
    /var/scripts/send-webhook.sh "High error rate: $ERROR_COUNT errors in last 100 lines"
fi

# Database connection issues
if grep -q "database connection failed" /var/log/sirine-go/error.log; then
    /var/scripts/send-webhook.sh "Database connection issues detected"
fi
```

---

## 📊 Performance Metrics

### **Key Performance Indicators (KPIs)**

**Application Metrics:**
```bash
# Requests per second
tail -10000 /var/log/nginx/sirine-go-access.log | \
    awk '{print $4}' | cut -d: -f1-2 | \
    sort | uniq -c

# Average response time
awk '{sum+=$NF; count++} END {print "Average:", sum/count, "ms"}' \
    /var/log/nginx/sirine-go-access.log

# 95th percentile response time
awk '{print $NF}' /var/log/nginx/sirine-go-access.log | \
    sort -n | \
    awk 'BEGIN{c=0} {total[c]=$1; c++} END{print total[int(c*0.95)]}'
```

**System Metrics:**
```bash
# Memory usage
free -h

# CPU load
uptime

# Disk I/O
iostat -x 1 5
```

### **Metrics Dashboard (Simple)**

```bash
#!/bin/bash
# File: /var/scripts/metrics-dashboard.sh

clear
echo "================================"
echo "  Sirine Go Metrics Dashboard"
echo "================================"
echo ""

# Service status
echo "🔧 SERVICE STATUS"
systemctl status sirine-go --no-pager -l | head -3
echo ""

# Resource usage
echo "💻 RESOURCE USAGE"
echo "CPU: $(top -bn1 | grep "Cpu(s)" | awk '{print $2}')"
echo "Memory: $(free | grep Mem | awk '{printf("%.1f%%", $3/$2 * 100)}')"
echo "Disk: $(df -h / | awk 'NR==2 {print $5}')"
echo ""

# Database connections
echo "🗄️ DATABASE"
echo "Connections: $(mysql -u sirine_user -p$DB_PASSWORD -e "SHOW STATUS LIKE 'Threads_connected';" | tail -1 | awk '{print $2}')"
echo ""

# Recent errors
echo "❌ RECENT ERRORS (last hour)"
echo "Count: $(grep -c "ERROR" /var/log/sirine-go/error.log || echo "0")"
```

---

## ✅ Monitoring Checklist

### Daily Tasks
- [ ] Check application logs for errors
- [ ] Verify service is running
- [ ] Check disk space
- [ ] Review slow query log

### Weekly Tasks
- [ ] Analyze error trends
- [ ] Review performance metrics
- [ ] Check log file sizes
- [ ] Verify backup completion

### Monthly Tasks
- [ ] Review all monitoring scripts
- [ ] Update alert thresholds
- [ ] Clean old logs (>30 days)
- [ ] Performance audit

---

## 🎯 Target Metrics

### Application
- **Uptime:** > 99.9%
- **Response Time:** < 100ms (avg)
- **Error Rate:** < 0.1%
- **Memory Usage:** < 300MB

### System
- **CPU Usage:** < 70% (avg)
- **Memory Usage:** < 80%
- **Disk Usage:** < 85%
- **Load Average:** < CPU cores

### Database
- **Query Time:** < 50ms (avg)
- **Connection Pool:** < 70% utilized
- **Buffer Pool Hit Ratio:** > 99%
- **Slow Queries:** < 1%

---

## 📞 Support

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733

## 📖 Related Documentation

- [Production Deployment](./production-deployment.md)
- [Backup & Recovery](./backup-recovery.md)
- [Troubleshooting](../09-troubleshooting/README.md)
- [Performance Optimization](../05-guides/performance/README.md)

---

**Last Updated:** 29 Desember 2025  
**Version:** 1.0.0
