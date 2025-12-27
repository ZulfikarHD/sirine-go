# ✅ Setup Checklist - Sirine Go App

Gunakan checklist ini untuk memastikan semua setup sudah benar.

> **📖 Setup guide:** [QUICKSTART.md](./QUICKSTART.md) atau [SETUP_GUIDE.md](./SETUP_GUIDE.md)

---

## 📋 Prerequisites

- [ ] Go 1.24+ terinstall (`go version`)
- [ ] Node.js 18+ terinstall (`node --version`)
- [ ] Yarn terinstall (`yarn --version`)
- [ ] MySQL 8.0+ terinstall dan berjalan (`mysql --version`)
- [ ] Git terinstall (`git --version`)

**Verify:**
```bash
go version && node --version && yarn --version && mysql --version && git --version
```

---

## 🗄️ Database Setup

- [ ] MySQL service berjalan
- [ ] Database `sirine_go` sudah dibuat
- [ ] User MySQL sudah dikonfigurasi
- [ ] Koneksi database berhasil di-test

**Commands:**
```bash
# Check MySQL
sudo systemctl status mysql

# Create database
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# Test connection
mysql -u root -p sirine_go -e "SELECT 1;"
```

---

## ⚙️ Configuration

- [ ] File `backend/.env` sudah ada
- [ ] `DB_USER` sudah diset
- [ ] `DB_PASSWORD` sudah diset
- [ ] `DB_NAME` sudah diset (sirine_go)
- [ ] `SERVER_PORT` sudah diset (8080)
- [ ] `TZ` sudah diset (Asia/Jakarta)

**Verify:**
```bash
cat backend/.env
```

---

## 📦 Dependencies

### Backend (Go)
- [ ] Go modules initialized (`backend/go.mod` exists)
- [ ] Dependencies downloaded
- [ ] No module errors

**Commands:**
```bash
cd backend
go mod download
go mod tidy
go list -m all | head -5
```

### Frontend (Node.js)
- [ ] `frontend/package.json` exists
- [ ] `frontend/node_modules` folder created
- [ ] `frontend/yarn.lock` generated
- [ ] No dependency errors

**Commands:**
```bash
cd frontend
yarn install
yarn list --depth=0 | grep -E "vue|vite|tailwind"
```

---

## 🔨 Build Test

### Backend
- [ ] Backend compiles without errors
- [ ] No syntax errors
- [ ] All imports resolved

**Commands:**
```bash
cd backend
go build -o test-build cmd/server/main.go
rm test-build
echo "✅ Backend compilation successful!"
```

### Frontend
- [ ] Frontend builds without errors
- [ ] Tailwind CSS configured
- [ ] Vite config valid
- [ ] PWA plugin configured

**Commands:**
```bash
cd frontend
yarn build
ls dist/  # Should show index.html and assets/
```

---

## 🚀 Runtime Test

### Backend Server
- [ ] Backend starts without errors
- [ ] Database connection successful
- [ ] Auto migration runs
- [ ] Server listening on port 8080
- [ ] No runtime errors in logs

**Commands:**
```bash
# Terminal 1
make dev-backend
```

**Expected Output:**
```
Database connected successfully
Server berjalan di port 8080
```

### Frontend Server
- [ ] Frontend starts without errors
- [ ] Vite dev server running
- [ ] Hot reload working
- [ ] No console errors
- [ ] Accessible at localhost:5173

**Commands:**
```bash
# Terminal 2
make dev-frontend
```

**Expected Output:**
```
VITE v7.2.4 ready in xxx ms
➜ Local: http://localhost:5173/
```

---

## 🧪 API Testing

### Health Check
- [ ] Health endpoint responds
- [ ] Returns 200 OK
- [ ] JSON response valid

**Test:**
```bash
curl http://localhost:8080/health
```

**Expected:**
```json
{"status":"ok","message":"Server berjalan dengan baik"}
```

### CRUD Operations
- [ ] GET /api/examples works
- [ ] POST /api/examples works
- [ ] GET /api/examples/:id works
- [ ] PUT /api/examples/:id works
- [ ] DELETE /api/examples/:id works

**Test:**
```bash
# Create
curl -X POST http://localhost:8080/api/examples \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","content":"Test content","is_active":true}'

# Read all
curl http://localhost:8080/api/examples

# Read one
curl http://localhost:8080/api/examples/1

# Update
curl -X PUT http://localhost:8080/api/examples/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated","content":"Updated content","is_active":false}'

# Delete
curl -X DELETE http://localhost:8080/api/examples/1
```

---

## 🎨 Frontend Testing

### UI Components
- [ ] Page loads without errors
- [ ] Tailwind CSS styles applied
- [ ] Animations working (Motion-v)
- [ ] Responsive design works
- [ ] Online/Offline indicator visible

### Functionality
- [ ] "Tambah Data Baru" button works
- [ ] Form appears when clicked
- [ ] Form validation works
- [ ] Submit creates data
- [ ] Data appears in list
- [ ] Edit button works
- [ ] Delete button works
- [ ] Confirmation dialog appears

### Browser Console
- [ ] No JavaScript errors
- [ ] No network errors
- [ ] Service Worker registered
- [ ] API calls successful

**Check in DevTools:**
1. Open DevTools (F12)
2. Console tab → no errors
3. Network tab → API calls 200 OK
4. Application → Service Workers → activated

---

## 🌐 PWA Testing

### Service Worker
- [ ] Service Worker registered
- [ ] Status: "activated and running"
- [ ] Cache storage created
- [ ] Static assets cached

**Check:**
DevTools → Application → Service Workers

### Offline Mode
- [ ] App works offline
- [ ] Cached data accessible
- [ ] Status indicator shows "Offline"
- [ ] Page loads from cache

**Test:**
1. Load app normally
2. DevTools → Network → Offline
3. Refresh page
4. App should still work ✅

### PWA Installation
- [ ] Install prompt available
- [ ] App can be installed
- [ ] Installed app works
- [ ] Manifest valid

**Test:**
- Desktop: Install icon in address bar
- Mobile: "Add to Home Screen" available

---

## 📱 Mobile Testing

### Responsive Design
- [ ] Mobile layout works
- [ ] Touch interactions work
- [ ] Buttons are touch-friendly
- [ ] Forms are mobile-friendly
- [ ] No horizontal scroll

**Test:**
1. DevTools → Toggle device toolbar
2. Test different screen sizes
3. Test on real mobile device

### Mobile Browser
- [ ] Works on Chrome Android
- [ ] Works on Safari iOS
- [ ] PWA installable on mobile
- [ ] Offline mode works on mobile

---

## 🔧 Development Tools

### Makefile
- [ ] `make help` shows commands
- [ ] `make install` works
- [ ] `make dev-backend` works
- [ ] `make dev-frontend` works
- [ ] `make build` works
- [ ] `make clean` works

**Test:**
```bash
make help
```

### Hot Reload
- [ ] Backend auto-restarts on code change (if using air)
- [ ] Frontend auto-reloads on code change (Vite)
- [ ] Changes reflect immediately

---

## 📊 Database Verification

### Tables Created
- [ ] `examples` table exists
- [ ] Columns correct (id, title, content, is_active, timestamps)
- [ ] Indexes created
- [ ] Soft delete working (deleted_at)

**Check:**
```bash
mysql -u root -p sirine_go -e "DESCRIBE examples;"
mysql -u root -p sirine_go -e "SELECT * FROM examples;"
```

### Data Integrity
- [ ] Create inserts data
- [ ] Update modifies data
- [ ] Delete soft-deletes data
- [ ] Timestamps auto-update

---

## 🐛 Common Issues Resolved

- [ ] No "database connection failed" error
- [ ] No "port already in use" error
- [ ] No "module not found" error
- [ ] No "command not found" error
- [ ] No CORS errors
- [ ] No 404 errors on API calls

---

## 📝 Documentation

- [ ] README.md reviewed
- [ ] QUICKSTART.md reviewed
- [ ] SETUP_GUIDE.md reviewed
- [ ] API_DOCUMENTATION.md reviewed
- [ ] ARCHITECTURE_EXPLAINED.md reviewed
- [ ] FOLDER_STRUCTURE.md reviewed
- [ ] DEPLOYMENT.md reviewed

---

## 🎯 Final Verification

### Backend
- [ ] ✅ Server starts successfully
- [ ] ✅ Database connected
- [ ] ✅ All endpoints working
- [ ] ✅ Error handling works
- [ ] ✅ CORS configured

### Frontend
- [ ] ✅ App loads successfully
- [ ] ✅ UI renders correctly
- [ ] ✅ Animations smooth
- [ ] ✅ API integration works
- [ ] ✅ PWA features work
- [ ] ✅ Offline mode works
- [ ] ✅ Responsive design works

### Overall
- [ ] ✅ No errors in backend logs
- [ ] ✅ No errors in frontend console
- [ ] ✅ No errors in database
- [ ] ✅ All features working
- [ ] ✅ Documentation complete
- [ ] ✅ Ready for development
- [ ] ✅ Ready for production (after security hardening)

---

## 🎉 Completion

Jika semua checklist di atas sudah ✅:

### **🎊 SELAMAT! Setup Anda 100% BERHASIL! 🎊**

### Next Steps:
1. Start developing your features
2. Customize models and UI
3. Add authentication if needed
4. Deploy to production when ready

### Quick Commands:
```bash
# Start development
make dev-backend    # Terminal 1
make dev-frontend   # Terminal 2

# Open browser
# http://localhost:5173

# Start coding! 🚀
```

---

## 📞 Need Help?

Jika ada yang tidak ✅:

1. **Check error messages** di terminal/console
2. **Troubleshooting guides:**
   - [FAQ.md](./FAQ.md) - Common solutions
   - [SETUP_GUIDE.md](./SETUP_GUIDE.md) - Detailed troubleshooting
3. **Contact developer:**
   - Zulfikar Hidayatullah
   - +62 857-1583-8733

---

## 📚 Next Steps After Checklist

Semua ✅? Great! Lanjut ke:
- **[ARCHITECTURE_EXPLAINED.md](./ARCHITECTURE_EXPLAINED.md)** - Understand tech
- **[CUSTOMIZATION_GUIDE.md](./CUSTOMIZATION_GUIDE.md)** - Build features
- **[TESTING.md](./TESTING.md)** - Test your code

---

**Date:** _______________  
**Checked by:** _______________  
**Status:** ☐ In Progress  ☐ Complete  

---

**Last Updated:** 27 Desember 2025  
**Version:** 1.0.0
