# Sirine Go App

Web application offline-capable menggunakan Gin (Go), Vue 3, Vite, dan MySQL dengan Tailwind CSS dan Motion-v untuk animasi.

## 🚀 Tech Stack

### Backend
- **Go** dengan **Gin Framework**
- **GORM** untuk ORM
- **MySQL** sebagai database
- **CORS** middleware
- Service pattern architecture

### Frontend
- **Vue 3** dengan Composition API
- **Vite** sebagai build tool
- **Tailwind CSS** untuk styling
- **Motion-v** (@motionone/vue) untuk animasi
- **Axios** untuk HTTP requests
- **VueUse** untuk composables utilities
- **PWA** dengan Vite Plugin PWA untuk offline capabilities

## 📋 Prerequisites

- Go 1.24+ terinstall
- Node.js 18+ dan Yarn terinstall
- MySQL 8.0+ terinstall dan berjalan
- Git

## 🛠️ Setup & Installation

### 1. Clone Repository

```bash
cd /home/sirinedev/WebApp/Developement/sirine-go
```

### 2. Setup Database

Buat database MySQL:

```sql
CREATE DATABASE sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. Konfigurasi Environment

File `.env` sudah dibuat. Sesuaikan dengan konfigurasi MySQL Anda:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=sirine_go

SERVER_PORT=8080
GIN_MODE=debug

TZ=Asia/Jakarta
```

### 4. Install Dependencies Backend

```bash
go mod download
```

### 5. Install Dependencies Frontend

```bash
cd frontend
yarn install
```

## 🏃 Running the Application

### Development Mode

#### 1. Jalankan Backend (Terminal 1)

```bash
go run cmd/server/main.go
```

Backend akan berjalan di `http://localhost:8080`

#### 2. Jalankan Frontend (Terminal 2)

```bash
cd frontend
yarn dev
```

Frontend akan berjalan di `http://localhost:5173`

### Production Mode

#### 1. Build Frontend

```bash
cd frontend
yarn build
```

#### 2. Build & Run Backend

```bash
go build -o sirine-go cmd/server/main.go
./sirine-go
```

Aplikasi akan berjalan di `http://localhost:8080` dengan frontend yang sudah di-build.

## 📁 Struktur Project

```
sirine-go/
├── cmd/
│   └── server/
│       └── main.go           # Entry point aplikasi
├── config/
│   └── config.go             # Konfigurasi aplikasi
├── database/
│   └── database.go           # Database connection & migration
├── models/
│   └── example.go            # Database models
├── services/
│   └── example_service.go    # Business logic layer
├── handlers/
│   └── example_handler.go    # HTTP handlers
├── middleware/
│   └── cors.go               # CORS middleware
├── routes/
│   └── routes.go             # Route definitions
├── frontend/
│   ├── src/
│   │   ├── components/       # Vue components
│   │   ├── composables/      # Composable functions
│   │   ├── views/            # Page views
│   │   ├── App.vue           # Root component
│   │   ├── main.js           # Entry point
│   │   └── style.css         # Global styles (Tailwind)
│   ├── public/               # Static assets
│   ├── index.html
│   ├── vite.config.js        # Vite configuration
│   ├── tailwind.config.js    # Tailwind configuration
│   └── package.json
├── .env                      # Environment variables
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 🎨 Features

### Backend Features
- ✅ RESTful API dengan Gin
- ✅ CRUD operations dengan GORM
- ✅ Service pattern untuk business logic
- ✅ CORS middleware
- ✅ Environment-based configuration
- ✅ Auto migration database
- ✅ Timezone Asia/Jakarta (WIB)
- ✅ Pesan error dalam Bahasa Indonesia

### Frontend Features
- ✅ Modern UI dengan Tailwind CSS
- ✅ Smooth animations dengan Motion-v
- ✅ Responsive design (mobile-first)
- ✅ Offline capabilities dengan PWA
- ✅ Online/Offline status indicator
- ✅ API caching untuk offline access
- ✅ Composable pattern untuk reusable logic
- ✅ Form validation
- ✅ Loading & error states

## 🔌 API Endpoints

### Health Check
```
GET /health
```

### Examples
```
GET    /api/examples      # Get all examples
GET    /api/examples/:id  # Get example by ID
POST   /api/examples      # Create new example
PUT    /api/examples/:id  # Update example
DELETE /api/examples/:id  # Delete example
```

### Example Request Body (POST/PUT)
```json
{
  "title": "Judul Example",
  "content": "Konten example",
  "is_active": true
}
```

## 🌐 Offline Capabilities

Aplikasi ini dapat berfungsi 100% offline dengan fitur:

1. **Service Worker**: Mengcache semua assets (JS, CSS, HTML, images)
2. **API Caching**: Request API di-cache dengan strategi NetworkFirst
3. **PWA Manifest**: Aplikasi dapat diinstall sebagai Progressive Web App
4. **Online/Offline Indicator**: Menampilkan status koneksi real-time

### Cara Kerja Offline:
- Saat online: Data di-fetch dari server dan di-cache
- Saat offline: Data diambil dari cache
- Saat kembali online: Data otomatis sync dengan server

## 🎯 Best Practices

### Backend
- Service pattern untuk separation of concerns
- Error handling dengan pesan Bahasa Indonesia
- Consistent API response format
- Environment-based configuration

### Frontend
- Composition API untuk better code organization
- Composables untuk reusable logic
- Component-based architecture
- Mobile-first responsive design
- Smooth animations untuk better UX
- Loading & error states untuk user feedback

## 📱 Mobile Optimization

Aplikasi ini dioptimasi untuk pengalaman mobile:
- Responsive grid layout
- Touch-friendly buttons
- Smooth scroll behavior
- Mobile-first CSS
- PWA installable di mobile devices

## 🔧 Customization

### Menambah Model Baru
1. Buat model di `models/`
2. Buat service di `services/`
3. Buat handler di `handlers/`
4. Tambahkan routes di `routes/routes.go`
5. Tambahkan migration di `cmd/server/main.go`

### Menambah Component Vue
1. Buat component di `frontend/src/components/`
2. Buat composable jika perlu di `frontend/src/composables/`
3. Import dan gunakan di views

## 🐛 Troubleshooting

### Database Connection Error
- Pastikan MySQL berjalan
- Check credentials di `.env`
- Pastikan database sudah dibuat

### Frontend Build Error
- Hapus `node_modules` dan `yarn.lock`
- Run `yarn install` lagi
- Check Node.js version (minimal 18+)

### Port Already in Use
- Backend: Ubah `SERVER_PORT` di `.env`
- Frontend: Ubah `server.port` di `vite.config.js`

## 👨‍💻 Developer

**Zulfikar Hidayatullah**
- Phone: +62 857-1583-8733
- Timezone: Asia/Jakarta (WIB)

## 📄 License

This project is private and proprietary.
