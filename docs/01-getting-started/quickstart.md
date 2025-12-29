# ⚡ Quick Start - Sirine Go App

Setup dan jalankan aplikasi dalam **5 menit**!

> **📖 Butuh detail lengkap?** Lihat [SETUP_GUIDE.md](./SETUP_GUIDE.md)

## 🎯 TL;DR

```bash
# 1. Setup database
mysql -u root -p -e "CREATE DATABASE sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 2. Edit .env (sesuaikan DB_PASSWORD)
nano backend/.env

# 3. Install dependencies
make install

# 4. Jalankan backend (Terminal 1)
make dev-backend

# 5. Jalankan frontend (Terminal 2)
make dev-frontend

# 6. Buka browser
# http://localhost:5173
```

---

## 📋 Prerequisites Check

Pastikan sudah terinstall:

```bash
go version      # Go 1.24+
node --version  # Node.js 18+
yarn --version  # Yarn 1.22+
mysql --version # MySQL 8.0+
```

Jika belum, lihat [SETUP_GUIDE.md](SETUP_GUIDE.md#prerequisites) untuk instalasi.

---

## 🚀 Step-by-Step (5 Menit)

### **Step 1: Database (30 detik)**

```bash
mysql -u root -p
```

```sql
CREATE DATABASE sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

### **Step 2: Environment (30 detik)**

```bash
nano backend/.env
```

Ubah **hanya** ini:
```env
DB_PASSWORD=your_mysql_password
```

Save: `Ctrl+O`, Enter, `Ctrl+X`

### **Step 3: Install (2 menit)**

```bash
make install
```

Tunggu sampai selesai...

### **Step 4: Run Backend (30 detik)**

**Terminal 1:**
```bash
make dev-backend
```

Tunggu sampai muncul:
```
Database connected successfully
Server berjalan di port 8080
```

✅ **Backend Ready!**

### **Step 5: Run Frontend (30 detik)**

**Terminal 2:**
```bash
make dev-frontend
```

Tunggu sampai muncul:
```
➜  Local:   http://localhost:5173/
```

✅ **Frontend Ready!**

### **Step 6: Test (30 detik)**

Buka browser: **http://localhost:5173**

1. Klik "Tambah Data Baru"
2. Isi form
3. Klik "Simpan"
4. Data muncul ✅

**🎉 Done! Aplikasi berjalan!**

---

## 🧪 Quick Test

### **Test API:**
```bash
curl http://localhost:8080/health
```

Expected: `{"status":"ok","message":"Server berjalan dengan baik"}`

### **Test Offline:**
1. F12 → Network tab
2. Set "Offline"
3. Refresh page
4. App tetap jalan ✅

---

## 🐛 Quick Troubleshooting

### **Database Error?**
```bash
# Check MySQL running
sudo systemctl status mysql

# Start if not running
sudo systemctl start mysql
```

### **Port Already in Use?**
```bash
# Backend (8080)
sudo lsof -i :8080
sudo kill -9 <PID>

# Frontend (5173)
sudo lsof -i :5173
sudo kill -9 <PID>
```

### **Dependencies Error?**
```bash
# Backend
cd backend && go mod download && go mod tidy

# Frontend
cd frontend && rm -rf node_modules && yarn install
```

---

## 🎯 Commands Cheat Sheet

```bash
make help              # Show all commands
make install           # Install dependencies
make dev-backend       # Run backend
make dev-frontend      # Run frontend
make build             # Build for production
make clean             # Clean build files
```

---

## 🎉 Success! Apa Selanjutnya?

Setup berhasil! Sekarang Anda bisa:

### 📚 **Belajar Lebih Dalam:**
- [ARCHITECTURE_EXPLAINED.md](./ARCHITECTURE_EXPLAINED.md) - Pahami tech stack
- [FOLDER_STRUCTURE.md](./FOLDER_STRUCTURE.md) - Struktur project
- [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) - API reference

### 🛠️ **Mulai Develop:**
- [CUSTOMIZATION_GUIDE.md](./CUSTOMIZATION_GUIDE.md) - Tambah fitur baru
- [TESTING.md](./TESTING.md) - Testing guide

### 🚀 **Deploy:**
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Deploy ke production

### ❓ **Troubleshooting:**
- [FAQ.md](./FAQ.md) - Common questions & solutions
- [CHECKLIST.md](./CHECKLIST.md) - Verify setup

---

## 📖 Documentation Flow

```
QUICKSTART.md (You are here)
    ↓
SETUP_GUIDE.md (Detail setup)
    ↓
CHECKLIST.md (Verify setup)
    ↓
ARCHITECTURE_EXPLAINED.md (Understand tech)
    ↓
CUSTOMIZATION_GUIDE.md (Build features)
    ↓
TESTING.md (Test your code)
    ↓
DEPLOYMENT.md (Deploy to production)
```

---

## 📞 Need Help?

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733

**Check:**
- [FAQ.md](./FAQ.md) - Common issues
- [SETUP_GUIDE.md](./SETUP_GUIDE.md) - Detailed troubleshooting

---

**Happy Coding! 🚀**
