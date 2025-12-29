# 🏗️ Arsitektur & Penjelasan Package - Sirine Go App

Dokumen ini menjelaskan **SEMUA** package yang digunakan dan **KENAPA** diperlukan.

## 📁 Struktur Folder (Baru)

Lihat [folder-structure.md](./folder-structure.md) untuk detail struktur folder yang akurat.

---

## 🔧 Backend Packages (Go)

### 1. **github.com/gin-gonic/gin** - Web Framework

**Apa itu?**
- Framework web untuk Go, seperti Express.js di Node.js atau Laravel di PHP

**Kenapa dibutuhkan?**
- Handle HTTP requests/responses
- Routing (GET /api/users, POST /api/users, dll)
- Middleware support
- JSON parsing otomatis
- Performance tinggi (bisa handle 40,000+ requests/second)

### 2. **gorm.io/gorm** + **gorm.io/driver/mysql** - ORM

**Apa itu?**
- ORM = Object-Relational Mapping
- Convert Go struct ↔ MySQL table

**Kenapa dibutuhkan?**
- **Auto Migration:** Buat table otomatis dari struct
- **Type Safety:** Compile-time error checking
- **SQL Injection Protection:** Query aman otomatis
- **Relationship:** Handle foreign keys, joins, dll

### 3. **github.com/golang-jwt/jwt/v5** - Authentication

**Apa itu?**
- Library untuk generate dan validate JSON Web Tokens (JWT)

**Kenapa dibutuhkan?**
- **Stateless Auth:** Server tidak perlu simpan session di memory
- **Security:** Standard industri untuk secure API authentication
- **Claims:** Bisa simpan data user (ID, Role) dalam token yang terenkripsi

### 4. **github.com/gin-contrib/cors** - CORS Middleware

**Apa itu?**
- CORS = Cross-Origin Resource Sharing

**Kenapa dibutuhkan?**
- Browser BLOCK request dari domain berbeda (security policy)
- Mengizinkan frontend (`localhost:5173`) mengakses backend (`localhost:8080`)

### 5. **github.com/joho/godotenv** - Environment Variables

**Apa itu?**
- Load file `.env` ke environment variables

**Kenapa dibutuhkan?**
- **Security:** Credentials (DB password, JWT Secret) tidak di-commit ke Git
- **Config:** Mudah ganti config untuk Dev/Prod

---

## 🎨 Frontend Packages (Node.js/Yarn)

### 1. **vue@^3.5.24** - JavaScript Framework

**Apa itu?**
- Framework untuk build user interface

**Kenapa dibutuhkan?**
- **Reactive:** Data berubah → UI auto update
- **Component-Based:** Reusable components
- **Composition API:** Better code organization

### 2. **pinia@^2.2.0** - State Management

**Apa itu?**
- Store library resmi untuk Vue (pengganti Vuex)

**Kenapa dibutuhkan?**
- **Global State:** Share data (User auth, Theme) antar component dengan mudah
- **DevTools:** Mudah debug state changes
- **Modular:** Store bisa dipecah per fitur (Auth Store, UI Store)

### 3. **vite@^7.2.4** - Build Tool

**Apa itu?**
- Modern build tool untuk frontend

**Kenapa dibutuhkan?**
- **Super Fast:** 10-100x lebih cepat dari Webpack
- **HMR:** Instant refresh saat coding
- **Optimized Build:** Bundle optimization untuk production

### 4. **tailwindcss@^4.1.18** - CSS Framework (NEWEST!)

**Apa itu?**
- Utility-first CSS framework

**Kenapa dibutuhkan?**
- **Fast Styling:** Tidak perlu tulis CSS manual
- **Consistent Design:** Design system built-in
- **Responsive:** Mobile-first responsive design

### 5. **motion-v@^1.7.4** - Animation Library

**Apa itu?**
- Library animasi performa tinggi untuk Vue (based on Motion One)

**Kenapa dibutuhkan?**
- **Better UX:** Animasi membuat UI terasa smooth & modern
- **Performance:** Hardware accelerated animations
- **Simple API:** Declarative animations (`<Motion :initial="..." :animate="...">`)

### 6. **axios@^1.13.2** - HTTP Client

**Apa itu?**
- Library untuk HTTP requests

**Kenapa dibutuhkan?**
- **Interceptors:** Auto add Bearer Token ke setiap request
- **Global Error Handling:** Handle 401 Unauthorized (auto logout) di satu tempat

### 7. **vite-plugin-pwa** + **workbox-window** - PWA

**Apa itu?**
- Plugin untuk Progressive Web App

**Kenapa dibutuhkan?**
- **Offline Capabilities:** App jalan tanpa internet
- **Installable:** Install sebagai native app di Desktop/Mobile
- **Caching:** Cache assets & API responses

### 8. **vitest** - Testing Framework

**Apa itu?**
- Unit testing framework powered by Vite

**Kenapa dibutuhkan?**
- **Fast:** Menggunakan konfigurasi Vite yang sama
- **Quality Assurance:** Memastikan logic frontend berjalan benar

---

## 🏗️ Arsitektur Backend (Service Pattern)

### Flow Request:
```
HTTP Request (routes)
    ↓
Middleware (Auth/CORS)
    ↓
Handler (HTTP Layer/Validation)
    ↓
Service (Business Logic)
    ↓
Database (Models/GORM)
    ↓
Response
```

### Keuntungan:
1.  **Separation of Concerns:** Logic terpisah dari HTTP handling.
2.  **Testable:** Service bisa di-test tanpa menjalankan server HTTP.
3.  **Reusable:** Service bisa dipanggil dari Handler, CLI command, atau service lain.

---

## 🌐 PWA & Offline - Cara Kerja Detail

### 1. Service Worker Registration
Terjadi di `frontend/src/main.js` menggunakan `workbox-window`.

### 2. Caching Strategy
- **Static Assets (JS/CSS/Images):** CacheFirst (Ambil cache dulu, kalau tidak ada baru download).
- **API Calls:** NetworkFirst (Coba internet dulu, kalau offline ambil dari cache).

### 3. Offline Flow
User offline → Buka App → Service Worker serve HTML/JS dari Cache → App Jalan → API call gagal → Ambil data terakhir dari Cache.

---

## 📊 Perbandingan: Dengan vs Tanpa Package

| Fitur | Tanpa Package | Dengan Package |
|-------|---------------|----------------|
| **State** | Prop drilling (ribet) | **Pinia** (Global Store) |
| **Auth** | Manual session | **JWT** (Stateless & Secure) |
| **Styling** | Manual CSS | **Tailwind** (Utility classes) |
| **Animation** | Complex CSS/JS | **Motion-v** (Declarative) |
| **Offline** | Implement manual (hard) | **PWA Plugin** (Auto) |

---

## 🚀 Kesimpulan

Stack teknologi ini dipilih untuk memaksimalkan:
1.  **Performance** (Go + Vite)
2.  **User Experience** (Vue + Motion-v + PWA)
3.  **Developer Experience** (Tailwind + Pinia + GORM)
4.  **Security** (JWT + Middleware)

**Total Package:** ~15 core packages
**Total Benefit:** Production ready, secure, & scalable system.

---

## 📚 Related Documentation

**Understand the code:**
- **[folder-structure.md](./folder-structure.md)** - Where everything is
- **[api-documentation.md](../03-development/api-documentation.md)** - API reference
- **[project-summary.md](./project-summary.md)** - Complete overview

**Start coding:**
- **[customization-guide.md](../03-development/customization-guide.md)** - Add features
- **[testing.md](../03-development/testing.md)** - Test your code

---

**Developer:** Zulfikar Hidayatullah
**Date:** 28 Desember 2025
**Version:** 1.0.1
