# 📖 Panduan Setup Lengkap - Sirine Go App

Panduan step-by-step untuk setup dan menjalankan aplikasi dari awal.

## 📋 Daftar Isi

1. [Prerequisites](#prerequisites)
2. [Setup Database](#setup-database)
3. [Konfigurasi Environment](#konfigurasi-environment)
4. [Install Dependencies](#install-dependencies)
5. [Menjalankan Aplikasi](#menjalankan-aplikasi)
6. [Testing](#testing)
7. [Troubleshooting](#troubleshooting)

---

## 📋 Prerequisites

Pastikan semua ini sudah terinstall di sistem Anda:

### **Required Software:**

#### 1. **Go 1.24+**
```bash
# Check version
go version

# Jika belum terinstall (Ubuntu/Debian)
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### 2. **Node.js 18+ & Yarn**
```bash
# Check version
node --version
yarn --version

# Install Node.js (Ubuntu/Debian)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Install Yarn
npm install -g yarn
```

#### 3. **MySQL 8.0+**
```bash
# Check version
mysql --version

# Install MySQL (Ubuntu/Debian)
sudo apt update
sudo apt install mysql-server -y

# Secure installation
sudo mysql_secure_installation
```

#### 4. **Git**
```bash
# Check version
git --version

# Install Git (Ubuntu/Debian)
sudo apt install git -y
```

### **Verifikasi Prerequisites:**
```bash
go version      # Should show: go version go1.24.x
node --version  # Should show: v18.x.x or higher
yarn --version  # Should show: 1.22.x or higher
mysql --version # Should show: mysql Ver 8.0.x
git --version   # Should show: git version 2.x.x
```

✅ **Jika semua command di atas berhasil, Anda siap melanjutkan!**

---

## 🗄️ Setup Database

### **Step 1: Start MySQL Service**

```bash
# Check status
sudo systemctl status mysql

# Start jika belum running
sudo systemctl start mysql

# Enable auto-start on boot
sudo systemctl enable mysql
```

### **Step 2: Create Database**

**Option 1 - Via MySQL CLI:**
```bash
# Login ke MySQL
mysql -u root -p
# Masukkan password MySQL Anda
```

```sql
-- Buat database
CREATE DATABASE sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Verifikasi
SHOW DATABASES LIKE 'sirine_go';

-- Exit
EXIT;
```

**Option 2 - One-liner:**
```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

**Option 3 - Via SQL File:**
```bash
mysql -u root -p < backend/database/setup.sql
```

### **Step 3: Verifikasi Database**

```bash
mysql -u root -p -e "USE sirine_go; SHOW TABLES;"
```

Output (saat ini kosong, table akan dibuat otomatis oleh GORM):
```
Empty set (0.00 sec)
```

✅ **Database siap digunakan!**

---

## ⚙️ Konfigurasi Environment

### **Step 1: Copy Environment File**

File `.env` sudah ada di `backend/.env`. Jika tidak ada, buat:

```bash
cd backend
cat > .env << 'EOF'
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=sirine_go

# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# Timezone
TZ=Asia/Jakarta
EOF
```

### **Step 2: Edit Configuration**

```bash
nano backend/.env
```

**Ubah hanya bagian ini:**
```env
DB_PASSWORD=your_mysql_password_here
```

**Contoh:**
```env
DB_PASSWORD=mySecretPassword123
```

**Save:** `Ctrl+O`, Enter, `Ctrl+X`

### **Step 3: Verifikasi Configuration**

```bash
cat backend/.env
```

Pastikan semua value sudah benar, terutama `DB_PASSWORD`.

✅ **Environment configured!**

---

## 📦 Install Dependencies

### **Option 1 - Menggunakan Makefile (Recommended):**

Dari root folder project:
```bash
make install
```

Output yang diharapkan:
```
Installing Go dependencies...
go: downloading github.com/gin-gonic/gin v1.11.0
go: downloading gorm.io/gorm v1.31.1
...
Installing frontend dependencies...
[1/4] Resolving packages...
[2/4] Fetching packages...
[3/4] Linking dependencies...
[4/4] Building fresh packages...
✅ Dependencies installed!
```

### **Option 2 - Manual:**

**Backend:**
```bash
cd backend
go mod download
go mod tidy
cd ..
```

**Frontend:**
```bash
cd frontend
yarn install
cd ..
```

### **Verifikasi Installation:**

**Backend:**
```bash
cd backend
go list -m all | head -5
```

Output:
```
sirine-go/backend
github.com/gin-gonic/gin v1.11.0
gorm.io/gorm v1.31.1
...
```

**Frontend:**
```bash
cd frontend
yarn list --depth=0 | grep -E "vue|vite|tailwind"
```

Output:
```
├─ vue@3.5.24
├─ vite@7.2.4
├─ tailwindcss@4.1.18
```

✅ **Dependencies installed successfully!**

---

## 🚀 Menjalankan Aplikasi

### **Development Mode (Recommended)**

Buka **2 terminal** di root folder project.

#### **Terminal 1 - Backend:**

```bash
make dev-backend
```

**Expected Output:**
```
Starting backend server...
Database connected successfully

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] GET    /health                   --> main.main.func1 (4 handlers)
[GIN-debug] GET    /api/examples             --> sirine-go/backend/handlers.(*ExampleHandler).GetAll-fm (4 handlers)
[GIN-debug] POST   /api/examples             --> sirine-go/backend/handlers.(*ExampleHandler).Create-fm (4 handlers)
[GIN-debug] PUT    /api/examples/:id         --> sirine-go/backend/handlers.(*ExampleHandler).Update-fm (4 handlers)
[GIN-debug] DELETE /api/examples/:id         --> sirine-go/backend/handlers.(*ExampleHandler).Delete-fm (4 handlers)

Server berjalan di port 8080
```

✅ **Backend is running!**

#### **Terminal 2 - Frontend:**

```bash
make dev-frontend
```

**Expected Output:**
```
Starting frontend server...

VITE v7.2.4  ready in 523 ms

➜  Local:   http://localhost:5173/
➜  Network: use --host to expose
➜  press h + enter to show help
```

✅ **Frontend is running!**

#### **Buka Browser:**

```
http://localhost:5173
```

Anda akan melihat aplikasi dengan:
- Header "Sirine Go App"
- Status indicator (Online/Offline)
- Button "Tambah Data Baru"
- Empty state atau list data

🎉 **Aplikasi berhasil berjalan!**

---

## 🧪 Testing

### **1. Test Backend API**

Buka terminal baru:

**Health Check:**
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok","message":"Server berjalan dengan baik"}
```

**Get All Examples:**
```bash
curl http://localhost:8080/api/examples
```

Expected response:
```json
{"data":[]}
```

**Create Example:**
```bash
curl -X POST http://localhost:8080/api/examples \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Example",
    "content": "Ini adalah test content",
    "is_active": true
  }'
```

Expected response:
```json
{
  "message": "Data berhasil dibuat",
  "data": {
    "id": 1,
    "title": "Test Example",
    "content": "Ini adalah test content",
    "is_active": true,
    "created_at": "2025-12-27T11:00:00+07:00",
    "updated_at": "2025-12-27T11:00:00+07:00"
  }
}
```

✅ **Backend API working!**

### **2. Test Frontend UI**

Di browser (`http://localhost:5173`):

1. **Test Create:**
   - Klik "Tambah Data Baru"
   - Isi form:
     - Judul: "Test dari UI"
     - Konten: "Konten test"
     - Centang "Aktif"
   - Klik "Simpan"
   - ✅ Data muncul di list

2. **Test Edit:**
   - Klik button "Edit" di card
   - Ubah judul jadi "Test Updated"
   - Klik "Perbarui"
   - ✅ Data terupdate

3. **Test Delete:**
   - Klik button "Hapus"
   - Konfirmasi di dialog
   - ✅ Data terhapus

### **3. Test Offline Mode (PWA)**

1. Buka DevTools (F12)
2. Go to **Network** tab
3. Set throttling to **"Offline"**
4. Refresh page
5. ✅ App tetap berfungsi!
6. ✅ Status indicator berubah jadi "Offline"

### **4. Test Responsive Design**

1. Buka DevTools (F12)
2. Toggle device toolbar (Ctrl+Shift+M)
3. Test berbagai screen sizes:
   - Mobile (375px)
   - Tablet (768px)
   - Desktop (1920px)
4. ✅ Layout responsive di semua ukuran

### **5. Test PWA Installation**

**Desktop:**
1. Lihat icon install di address bar (⊕)
2. Klik icon
3. Klik "Install"
4. ✅ App terinstall sebagai native app

**Mobile:**
1. Buka di Chrome mobile
2. Menu (⋮) → "Add to Home Screen"
3. Tap "Add"
4. ✅ App terinstall di home screen

---

## 🐛 Troubleshooting

### **Problem 1: Database Connection Failed**

**Error:**
```
Failed to connect to database: Error 1045: Access denied for user 'root'@'localhost'
```

**Solution:**
```bash
# Check MySQL running
sudo systemctl status mysql

# Check password di .env
cat backend/.env | grep DB_PASSWORD

# Test connection manual
mysql -u root -p -e "USE sirine_go;"

# Jika password salah, reset:
sudo mysql
ALTER USER 'root'@'localhost' IDENTIFIED BY 'new_password';
FLUSH PRIVILEGES;
EXIT;

# Update backend/.env dengan password baru
```

---

### **Problem 2: Port Already in Use**

**Error (Backend):**
```
listen tcp :8080: bind: address already in use
```

**Solution:**
```bash
# Find process using port 8080
sudo lsof -i :8080

# Kill process
sudo kill -9 <PID>

# Atau ubah port di backend/.env
echo "SERVER_PORT=8081" >> backend/.env
```

**Error (Frontend):**
```
Port 5173 is in use
```

**Solution:**
```bash
# Find process using port 5173
sudo lsof -i :5173

# Kill process
sudo kill -9 <PID>
```

---

### **Problem 3: Go Module Errors**

**Error:**
```
go: module not found
```

**Solution:**
```bash
cd backend
go clean -modcache
go mod download
go mod tidy
```

---

### **Problem 4: Frontend Build Errors**

**Error:**
```
Error: Cannot find module 'vite'
```

**Solution:**
```bash
cd frontend
rm -rf node_modules yarn.lock
yarn install
```

---

### **Problem 5: CORS Errors**

**Error di Browser Console:**
```
Access to fetch at 'http://localhost:8080/api/examples' has been blocked by CORS policy
```

**Solution:**
```bash
# Pastikan backend running
# Restart backend
# Terminal 1
Ctrl+C
make dev-backend
```

---

### **Problem 6: Database Tables Not Created**

**Error:**
```
Error 1146: Table 'sirine_go.examples' doesn't exist
```

**Solution:**
```bash
# GORM auto-migrate akan buat table saat backend start
# Restart backend untuk trigger migration
Ctrl+C
make dev-backend

# Verifikasi table dibuat
mysql -u root -p sirine_go -e "SHOW TABLES;"
```

Expected output:
```
+---------------------+
| Tables_in_sirine_go |
+---------------------+
| examples            |
+---------------------+
```

---

### **Problem 7: Yarn Command Not Found**

**Error:**
```
bash: yarn: command not found
```

**Solution:**
```bash
# Install yarn globally
npm install -g yarn

# Verifikasi
yarn --version
```

---

## 📊 Verification Checklist

Setelah setup, verifikasi semua berjalan dengan baik:

### **Backend:**
- [ ] MySQL service running
- [ ] Database `sirine_go` exists
- [ ] `backend/.env` configured correctly
- [ ] Go dependencies installed
- [ ] Backend starts without errors
- [ ] Listening on port 8080
- [ ] Health check returns 200 OK
- [ ] Database tables auto-created

### **Frontend:**
- [ ] Node.js & Yarn installed
- [ ] Frontend dependencies installed
- [ ] Frontend starts without errors
- [ ] Accessible at http://localhost:5173
- [ ] No console errors
- [ ] Can create/read/update/delete data
- [ ] Responsive design works
- [ ] PWA features work

### **Integration:**
- [ ] Frontend can call backend API
- [ ] No CORS errors
- [ ] Data persists in database
- [ ] Offline mode works
- [ ] PWA installable

---

## 🎯 Next Steps

Setelah setup berhasil:

1. **Customize Models:**
   - Edit `backend/models/example.go`
   - Tambah fields sesuai kebutuhan
   - Restart backend untuk migration

2. **Add New Endpoints:**
   - Create service di `backend/services/`
   - Create handler di `backend/handlers/`
   - Add routes di `backend/routes/routes.go`

3. **Customize UI:**
   - Edit components di `frontend/src/components/`
   - Edit views di `frontend/src/views/`
   - Customize styles di `frontend/src/style.css`

4. **Add Features:**
   - Authentication
   - File upload
   - Search & filter
   - Pagination
   - Real-time updates

5. **Deploy to Production:**
   - See `DEPLOYMENT.md` for deployment guide

---

## 📚 Related Documentation

**After successful setup:**
- ✅ **[CHECKLIST.md](./CHECKLIST.md)** - Verify setup
- 📚 **[ARCHITECTURE_EXPLAINED.md](./ARCHITECTURE_EXPLAINED.md)** - Understand tech stack
- 📁 **[FOLDER_STRUCTURE.md](./FOLDER_STRUCTURE.md)** - Project structure
- 🔌 **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)** - API reference

**Next steps:**
- 🛠️ **[CUSTOMIZATION_GUIDE.md](./CUSTOMIZATION_GUIDE.md)** - Build features
- 🧪 **[TESTING.md](./TESTING.md)** - Test your code
- 🚀 **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Deploy to production

**Need help?**
- ❓ **[FAQ.md](./FAQ.md)** - Common questions & solutions

---

## 📞 Need Help?

Jika masih ada masalah:

1. Check error messages di terminal
2. Check browser console (F12)
3. Lihat troubleshooting section di atas
4. Contact developer:
   - Zulfikar Hidayatullah
   - +62 857-1583-8733

---

## ✅ Summary

Setup complete checklist:

```bash
# 1. Install prerequisites
✅ Go 1.24+
✅ Node.js 18+ & Yarn
✅ MySQL 8.0+
✅ Git

# 2. Setup database
✅ MySQL running
✅ Database created

# 3. Configure environment
✅ backend/.env configured

# 4. Install dependencies
✅ make install

# 5. Run application
✅ make dev-backend (Terminal 1)
✅ make dev-frontend (Terminal 2)
✅ Open http://localhost:5173

# 6. Test
✅ Backend API works
✅ Frontend UI works
✅ Offline mode works
✅ PWA installable
```

**🎉 Selamat! Setup berhasil! Happy coding! 🚀**

---

**Last Updated:** 27 Desember 2025  
**Version:** 1.0.0
