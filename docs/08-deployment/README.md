# 🚀 Deployment

Dokumentasi untuk deploy aplikasi ke production server.

---

## 📚 Files di Folder Ini

### **[production-deployment.md](./production-deployment.md)** 🌐
Complete production deployment guide.
- Server requirements & setup
- Install Go, MySQL, Nginx
- Database configuration
- Build & deploy aplikasi
- Systemd service setup
- Nginx reverse proxy
- SSL/HTTPS dengan Let's Encrypt
- Firewall configuration

**Cocok untuk:** Deploy to production safely dan correctly.

### **[monitoring.md](./monitoring.md)** 📊
Monitoring & logging guide untuk production.
- Application monitoring (service status, health checks)
- System monitoring (CPU, memory, disk usage)
- Database monitoring (connections, query performance)
- Log management & rotation
- Alert configuration (email, webhooks)
- Performance metrics & KPIs

**Cocok untuk:** Memastikan aplikasi berjalan optimal di production.

### **[backup-recovery.md](./backup-recovery.md)** 💾
Backup & disaster recovery procedures.
- Database backup (manual & automated)
- Application backup (files, config, uploads)
- Automated backup scheduling
- Remote backup storage (rsync, S3, Google Drive)
- Recovery procedures (full & point-in-time)
- Disaster recovery plan

**Cocok untuk:** Melindungi data dan recovery dari failure.

### ⚠️ Catatan Penting Sebelum Deploy

**Lokasi File `.env`:**
- Kode backend (`backend/cmd/server/main.go`) memuat `.env` dari **root repository** (`../.env` relatif dari `backend/`)
- Buat file `.env` di **root project** (`/var/www/sirine-go/.env`), bukan di `backend/.env`
- Atau sesuaikan path di kode jika ingin `.env` di folder `backend/`

**Lokasi Binary Hasil Build:**
- `Makefile` target `build` menghasilkan binary di **root** (`./sirine-go`)
- Dokumentasi production-deployment.md menunjukkan binary di `backend/sirine-go`
- Pastikan path `ExecStart` di systemd service sesuai dengan lokasi binary yang digunakan
- Rekomendasi: gunakan `make build` dan sesuaikan systemd service path atau build manual sesuai dokumentasi

---

## 🎯 Deployment Checklist

**Before deployment:**
- ✅ Test locally: [../development/TESTING.md](../03-development/TESTING.md)
- ✅ All features working
- ✅ No console errors
- ✅ Performance optimized

**During deployment:**
- Follow production-deployment.md step-by-step
- Don't skip security steps (SSL, firewall)
- Test after each major step

**After deployment:**
- Test all endpoints in production
- Test PWA offline mode
- Monitor logs for errors
- Setup backup automation

---

## 🔐 Production Security

**Must Have:**
- ✅ SSL/HTTPS certificate
- ✅ Firewall configured
- ✅ Strong database password
- ✅ Environment variables (not hardcoded)

**Recommended:**
- JWT authentication
- Rate limiting
- Security headers
- Regular backups

**See:** production-deployment.md → Security section

---

## 🔗 Related Documentation

**Before deployment:**
- Test thoroughly: [../development/TESTING.md](../03-development/TESTING.md)
- Review: [../architecture/project-summary.md](../02-architecture/project-summary.md)

**Deployment issues?**
- Troubleshooting: [../troubleshooting/FAQ.md](../09-troubleshooting/FAQ.md) → Deployment section

---

## 📂 Folder Navigation

← Back to [Documentation Root](../README.md)  
← Previous: [Development](../03-development/)  
→ Next: [Troubleshooting](../09-troubleshooting/)

---

**Last Updated:** 27 Desember 2025
