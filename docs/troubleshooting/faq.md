# ❓ FAQ - Frequently Asked Questions

Kumpulan pertanyaan umum dan solusinya untuk Sirine Go App.

---

## 📋 Daftar Isi

- [Setup & Installation](#setup--installation)
- [Database Issues](#database-issues)
- [Backend Issues](#backend-issues)
- [Frontend Issues](#frontend-issues)
- [Development Issues](#development-issues)
- [PWA & Offline](#pwa--offline)
- [Performance](#performance)
- [Deployment](#deployment)

---

## 🔧 Setup & Installation

### Q: Berapa lama waktu setup aplikasi ini?

**A:** Sekitar 5-10 menit jika semua prerequisites sudah terinstall. Lihat [QUICKSTART.md](./QUICKSTART.md).

---

### Q: Apa saja yang harus diinstall sebelum mulai?

**A:** Prerequisites yang dibutuhkan:
- Go 1.24+
- Node.js 18+
- Yarn
- MySQL 8.0+
- Git

Verify dengan:
```bash
go version && node --version && yarn --version && mysql --version
```

---

### Q: Apakah bisa menggunakan npm instead of yarn?

**A:** Bisa, tapi kami recommend yarn karena:
- Lebih cepat
- Lock file lebih reliable
- Better workspace support

Jika tetap ingin npm:
```bash
cd frontend
rm yarn.lock
npm install
npm run dev
```

---

### Q: Error "make: command not found"

**A:** Install make:

**Ubuntu/Debian:**
```bash
sudo apt install make
```

**macOS:**
```bash
xcode-select --install
```

**Alternatif tanpa make:**
```bash
# Backend
cd backend && go run cmd/server/main.go

# Frontend
cd frontend && yarn dev
```

---

## 🗄️ Database Issues

### Q: Error "Access denied for user 'root'@'localhost'"

**A:** Password MySQL salah atau user tidak ada.

**Solusi:**
```bash
# 1. Check password di .env
cat backend/.env | grep DB_PASSWORD

# 2. Test connection
mysql -u root -p -e "SELECT 1;"

# 3. Jika lupa password, reset:
sudo mysql
ALTER USER 'root'@'localhost' IDENTIFIED BY 'new_password';
FLUSH PRIVILEGES;
EXIT;

# 4. Update backend/.env dengan password baru
```

---

### Q: Error "Unknown database 'sirine_go'"

**A:** Database belum dibuat.

**Solusi:**
```bash
mysql -u root -p -e "CREATE DATABASE sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

---

### Q: Error "Error 1045: Access denied"

**A:** GORM tidak bisa connect ke MySQL.

**Solusi:**
1. Check MySQL running: `sudo systemctl status mysql`
2. Check credentials di `backend/.env`
3. Test manual: `mysql -u root -p sirine_go -e "SELECT 1;"`
4. Check port: Default 3306

---

### Q: Table tidak terbuat otomatis

**A:** Auto-migration tidak jalan.

**Solusi:**
```bash
# 1. Restart backend
Ctrl+C
make dev-backend

# 2. Check logs untuk error
# 3. Verify table:
mysql -u root -p sirine_go -e "SHOW TABLES;"
```

---

### Q: Bagaimana cara reset database?

**A:** 
```bash
# Drop dan create ulang
mysql -u root -p -e "DROP DATABASE sirine_go;"
mysql -u root -p -e "CREATE DATABASE sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# Restart backend untuk migration
make dev-backend
```

---

## 🔧 Backend Issues

### Q: Error "listen tcp :8080: bind: address already in use"

**A:** Port 8080 sudah digunakan.

**Solusi 1 - Kill process:**
```bash
sudo lsof -i :8080
sudo kill -9 <PID>
```

**Solusi 2 - Ganti port:**
```bash
# Edit backend/.env
SERVER_PORT=8081
```

---

### Q: Backend crash dengan error "panic: runtime error"

**A:** Ada bug di code.

**Solusi:**
1. Check terminal logs untuk detail error
2. Check line number yang error
3. Debug dengan menambah log:
```go
fmt.Printf("Debug: %+v\n", variable)
```

---

### Q: Hot reload tidak jalan di backend

**A:** Go tidak punya built-in hot reload.

**Solusi - Install air:**
```bash
go install github.com/cosmtrek/air@latest

# Run dengan air
cd backend
air
```

---

### Q: Error "go: module not found"

**A:** Dependencies tidak terinstall.

**Solusi:**
```bash
cd backend
go clean -modcache
go mod download
go mod tidy
```

---

### Q: CORS error di browser

**A:** Frontend tidak bisa akses backend.

**Solusi:**
1. Pastikan backend running
2. Check middleware CORS di `backend/middleware/cors.go`
3. Restart backend:
```bash
Ctrl+C
make dev-backend
```

---

## 🎨 Frontend Issues

### Q: Error "Port 5173 is in use"

**A:** Port frontend sudah digunakan.

**Solusi:**
```bash
sudo lsof -i :5173
sudo kill -9 <PID>
```

---

### Q: Error "Cannot find module 'vite'"

**A:** Dependencies tidak terinstall.

**Solusi:**
```bash
cd frontend
rm -rf node_modules yarn.lock
yarn install
```

---

### Q: Tailwind CSS tidak bekerja

**A:** PostCSS configuration issue.

**Solusi:**
```bash
# 1. Check files exist:
ls frontend/postcss.config.js
ls frontend/tailwind.config.js

# 2. Reinstall:
cd frontend
yarn add -D tailwindcss postcss autoprefixer
```

---

### Q: Build error "out of memory"

**A:** Node.js kehabisan memory.

**Solusi:**
```bash
# Increase memory limit
NODE_OPTIONS=--max_old_space_size=4096 yarn build
```

---

### Q: Hot reload tidak jalan di frontend

**A:** Vite HMR issue.

**Solusi:**
1. Clear cache: `rm -rf frontend/node_modules/.vite`
2. Restart: `Ctrl+C` → `yarn dev`
3. Hard refresh browser: `Ctrl+Shift+R`

---

### Q: Component tidak muncul di browser

**A:** Check beberapa hal:

**Solusi:**
1. Check browser console (F12) untuk error
2. Check component import path
3. Check component registered di parent
4. Verify API response di Network tab

---

## 🛠️ Development Issues

### Q: Bagaimana cara debug backend?

**A:** Gunakan print debugging:

```go
// handlers/example_handler.go
func (h *ExampleHandler) GetAll(c *gin.Context) {
    fmt.Println("=== Debug GetAll ===")
    
    examples, err := h.service.GetAll()
    fmt.Printf("Examples: %+v\n", examples)
    fmt.Printf("Error: %v\n", err)
    
    // ... rest of code
}
```

Atau gunakan Delve debugger:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug cmd/server/main.go
```

---

### Q: Bagaimana cara debug frontend?

**A:** Gunakan browser DevTools:

```javascript
// frontend/src/composables/useExamples.js
export const useExamples = () => {
  const fetchExamples = async () => {
    console.log('=== Fetching examples ===')
    
    const response = await api.get('/api/examples')
    console.log('Response:', response.data)
    
    // ... rest of code
  }
}
```

**DevTools:**
- Console: `F12` → Console tab
- Network: `F12` → Network tab
- Vue DevTools: Install extension

---

### Q: Bagaimana cara test API tanpa frontend?

**A:** Gunakan cURL atau Postman.

**cURL:**
```bash
# GET
curl http://localhost:8080/api/examples

# POST
curl -X POST http://localhost:8080/api/examples \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","content":"Content","is_active":true}'

# PUT
curl -X PUT http://localhost:8080/api/examples/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated"}'

# DELETE
curl -X DELETE http://localhost:8080/api/examples/1
```

**Postman:**
1. Download Postman
2. Import collection dari API docs
3. Test endpoints

---

### Q: Code saya tidak berubah setelah edit

**A:** Cache issue.

**Solusi:**
```bash
# Backend - restart
Ctrl+C
make dev-backend

# Frontend - clear cache
rm -rf frontend/node_modules/.vite
Ctrl+C
make dev-frontend

# Browser - hard refresh
Ctrl+Shift+R
```

---

## 🌐 PWA & Offline

### Q: Service Worker tidak register

**A:** Check beberapa hal:

**Solusi:**
1. Pastikan production build: `cd frontend && yarn build`
2. Serve production: Backend harus serve `frontend/dist`
3. Check browser console (F12)
4. Check DevTools → Application → Service Workers

---

### Q: Offline mode tidak bekerja

**A:** Service Worker belum active atau cache kosong.

**Solusi:**
```bash
# 1. Build production
cd frontend && yarn build

# 2. Clear cache
# DevTools (F12) → Application → Clear storage → Clear site data

# 3. Reload page
# 4. Go offline: DevTools → Network → Offline
# 5. Reload - should work from cache
```

---

### Q: PWA tidak bisa diinstall

**A:** Manifest atau Service Worker issue.

**Solusi:**
1. Check manifest: DevTools → Application → Manifest
2. Check Service Worker: DevTools → Application → Service Workers
3. Verify HTTPS (production) atau localhost (dev)
4. Check console untuk errors

---

### Q: Cache tidak update dengan data baru

**A:** Service Worker menggunakan stale cache.

**Solusi:**
1. DevTools → Application → Service Workers
2. Check "Update on reload"
3. Atau clear cache: Storage → Clear storage
4. Reload page

---

## 🚀 Performance

### Q: Backend lambat

**A:** Check beberapa hal:

**Solusi:**
1. **Database query optimization:**
   ```go
   // Bad - N+1 query
   db.Find(&examples)
   
   // Good - Preload
   db.Preload("Relation").Find(&examples)
   ```

2. **Add indexes:**
   ```go
   type Example struct {
       Title string `gorm:"index"` // Add index
   }
   ```

3. **Use caching** (Redis, in-memory)

---

### Q: Frontend lambat saat load

**A:** Bundle size terlalu besar.

**Solusi:**
1. **Lazy loading components:**
   ```javascript
   // Before
   import HeavyComponent from './HeavyComponent.vue'
   
   // After
   const HeavyComponent = defineAsyncComponent(() =>
     import('./HeavyComponent.vue')
   )
   ```

2. **Check bundle size:**
   ```bash
   yarn build
   # Check dist/assets/*.js size
   ```

3. **Optimize images:** Use WebP format

---

### Q: Animasi tidak smooth

**A:** Performance issue atau animation config.

**Solusi:**
1. Reduce animation complexity
2. Use `transform` instead of `top/left`
3. Add `will-change: transform`
4. Check browser performance: DevTools → Performance

---

## 🌍 Deployment

### Q: Bagaimana cara deploy ke production?

**A:** Lihat [DEPLOYMENT.md](./DEPLOYMENT.md) untuk panduan lengkap.

**Quick steps:**
1. Setup server (Ubuntu + Nginx + MySQL)
2. Build aplikasi: `make build`
3. Upload ke server
4. Configure systemd service
5. Setup Nginx reverse proxy
6. Install SSL certificate

---

### Q: Error 502 Bad Gateway setelah deploy

**A:** Backend tidak running atau Nginx config salah.

**Solusi:**
```bash
# Check backend running
sudo systemctl status sirine-go

# Start jika mati
sudo systemctl start sirine-go

# Check Nginx config
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx

# Check logs
sudo journalctl -u sirine-go -f
sudo tail -f /var/log/nginx/error.log
```

---

### Q: SSL certificate error

**A:** Certificate tidak valid atau expired.

**Solusi:**
```bash
# Renew certificate
sudo certbot renew

# Force renew
sudo certbot renew --force-renewal

# Check certificate
sudo certbot certificates
```

---

## 🔐 Security

### Q: Bagaimana cara menambah authentication?

**A:** Implement JWT authentication.

**Quick guide:**
1. Install JWT package: `go get github.com/golang-jwt/jwt/v5`
2. Create auth middleware
3. Protect routes dengan middleware
4. Frontend: Store token di localStorage
5. Add token ke axios headers

Lihat [CUSTOMIZATION_GUIDE.md](./CUSTOMIZATION_GUIDE.md) untuk contoh lengkap.

---

### Q: API terbuka untuk umum, bagaimana cara protect?

**A:** Beberapa opsi:

**Production security:**
1. **JWT Authentication** - Recommended
2. **API Key** - Simple tapi less secure
3. **Rate Limiting** - Prevent abuse
4. **CORS** - Already configured
5. **HTTPS** - Must have for production

---

## 📊 General

### Q: Bagaimana cara contribute ke project ini?

**A:** Project ini private, tapi Anda bisa:
1. Fork untuk penggunaan personal
2. Contact developer untuk collaboration
3. Report bugs via issue tracker (if available)

---

### Q: Apakah ada video tutorial?

**A:** Belum ada saat ini. Dokumentasi tertulis sudah sangat lengkap:
- [QUICKSTART.md](./QUICKSTART.md) - 5 menit setup
- [SETUP_GUIDE.md](./SETUP_GUIDE.md) - Detail guide
- [ARCHITECTURE_EXPLAINED.md](./ARCHITECTURE_EXPLAINED.md) - Penjelasan tech

---

### Q: Apakah support bahasa lain selain Indonesia?

**A:** Saat ini hanya Bahasa Indonesia. Untuk multi-language:
1. Install i18n package
2. Create language files
3. Update components dengan translations
4. Add language switcher

---

### Q: Pertanyaan saya tidak ada di sini

**A:** Contact developer:
- **Developer:** Zulfikar Hidayatullah
- **Phone:** +62 857-1583-8733
- **Email:** Check dengan developer

Atau check dokumentasi lain:
- [SETUP_GUIDE.md](./SETUP_GUIDE.md) - Setup troubleshooting
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Deployment troubleshooting
- [PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md) - Project overview

---

## 📞 Masih Butuh Bantuan?

Jika masih ada masalah:

1. ✅ Check error messages di terminal/console
2. ✅ Check dokumentasi related:
   - Setup issues → [SETUP_GUIDE.md](./SETUP_GUIDE.md)
   - API issues → [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
   - Deploy issues → [DEPLOYMENT.md](./DEPLOYMENT.md)
3. ✅ Google error message specific
4. ✅ Contact developer: +62 857-1583-8733

---

**Last Updated:** 27 Desember 2025  
**Version:** 1.0.0
