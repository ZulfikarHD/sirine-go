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
- Monitoring & logging
- Backup automation
- Troubleshooting deployment issues

**Cocok untuk:** Deploy to production safely dan correctly.

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
- ✅ Test locally: [../development/TESTING.md](../development/TESTING.md)
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
- Test thoroughly: [../development/TESTING.md](../development/TESTING.md)
- Review: [../architecture/project-summary.md](../architecture/project-summary.md)

**Deployment issues?**
- Troubleshooting: [../troubleshooting/FAQ.md](../troubleshooting/FAQ.md) → Deployment section

---

## 📂 Folder Navigation

← Back to [Documentation Root](../README.md)  
← Previous: [Development](../development/)  
→ Next: [Troubleshooting](../troubleshooting/)

---

**Last Updated:** 27 Desember 2025
