# 🚀 Deployment

Dokumentasi untuk deploy aplikasi ke production server.

---

## 📚 Files di Folder Ini

### **[DEPLOYMENT.md](./DEPLOYMENT.md)** 🌐
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

---

## 🎯 Deployment Checklist

**Before deployment:**
- ✅ Test locally: [../development/TESTING.md](../development/TESTING.md)
- ✅ All features working
- ✅ No console errors
- ✅ Performance optimized

**During deployment:**
- Follow DEPLOYMENT.md step-by-step
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

**See:** DEPLOYMENT.md → Security section

---

## 🔗 Related Documentation

**Before deployment:**
- Test thoroughly: [../development/TESTING.md](../development/TESTING.md)
- Review: [../architecture/PROJECT_SUMMARY.md](../architecture/PROJECT_SUMMARY.md)

**Deployment issues?**
- Troubleshooting: [../troubleshooting/FAQ.md](../troubleshooting/FAQ.md) → Deployment section

---

## 📂 Folder Navigation

← Back to [Documentation Root](../README.md)  
← Previous: [Development](../development/)  
→ Next: [Troubleshooting](../troubleshooting/)

---

**Last Updated:** 27 Desember 2025
