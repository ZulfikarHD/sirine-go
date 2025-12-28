# 📝 Documentation Maintenance Guide

## Overview

Dokumentasi yang up-to-date sama pentingnya dengan code yang bersih. Guide ini menjelaskan **kapan** Anda harus membuat atau mengupdate dokumentasi dan **apa saja** yang perlu diupdate berdasarkan jenis perubahan yang Anda lakukan.

---

## 🚦 When to Update Documentation?

Update dokumentasi diperlukan setiap kali ada perubahan yang mempengaruhi:
1. **Fungsionalitas User** (Fitur baru, perubahan flow)
2. **Developer Experience** (Setup, config, API, commands)
3. **System Architecture** (Database, infrastructure, security)

### Quick Checklist Matrix

| Jenis Perubahan | Files yang Wajib Diupdate |
|-----------------|---------------------------|
| ✨ **New Feature** | `CHANGELOG.md`, `README.md` (Features), `api-documentation.md` |
| 🗄️ **Database** | `models.md`, `migrations.md` |
| ⚙️ **Config/.env** | `.env.example`, `configuration.md` |
| 🔌 **API Endpoint** | `api-documentation.md` |
| 🔐 **Security** | `security.md` |
| 🐛 **Bug Fix** | `CHANGELOG.md` (jika signifikan) |

---

## 📋 Detailed Workflows

### 1. Menambahkan Fitur Baru (New Feature)

Jika Anda membuat fitur baru (contoh: "Export Laporan PDF"):

1. **Update `CHANGELOG.md`**: Tambahkan entry di bagian `Added`.
2. **Update `README.md`**: Tambahkan poin di section "Key Features" jika fitur major.
3. **Update API Docs**: Jika ada endpoint baru, dokumentasikan di `docs/development/api-documentation.md`.
   - Method, URL
   - Request Body/Params
   - Response Success & Error
4. **Create/Update User Journey**: Jika flow user berubah, update/buat file di `docs/user-journeys/`.
   - Gunakan diagram alur (text-based graph).

### 2. Perubahan Database (Database Changes)

Jika Anda menambah tabel atau kolom:

1. **Update Models Guide**: Edit `docs/guides/database/models.md` jika ada pattern baru.
2. **Update Migrations Guide**: Edit `docs/guides/database/migrations.md` jika ada workflow migrasi khusus.
3. **Seeders**: Jika perlu data awal, pastikan `make db-seed` diperbarui dan didokumentasikan.

### 3. Perubahan Konfigurasi (Configuration Changes)

Jika Anda menambah variable baru di `.env`:

1. **Update `.env.example`**: Backend dan Frontend.
2. **Update Configuration Guide**: Edit `docs/guides/configuration.md`.
   - Tambahkan nama variable.
   - Jelaskan fungsinya.
   - Berikan default value dan contoh.

### 4. Perubahan API (API Changes)

Jika Anda mengubah contract API (Request/Response):

1. **Update API Documentation**: `docs/development/api-documentation.md`.
   - **PENTING**: Tandai jika ada Breaking Change.
2. **Update Frontend Service**: Pastikan dokumentasi frontend (jika ada) sinkron.

---

## ✍️ Writing Style Guide

### Bahasa & Tone
- **Formal & Professional**: Gunakan Bahasa Indonesia baku untuk penjelasan.
- **English for Tech Terms**: Gunakan istilah teknis dalam Bahasa Inggris (e.g., "Request", "Response", "Middleware", "Endpoint").
- **Personality**: Professional, informatif, dan membantu (INFJ-style).

### Formatting
- **Code Blocks**: Selalu gunakan syntax highlighting (```go, ```json).
- **Icons**: Gunakan emoji yang relevan untuk visual hierarchy (e.g., 🚀, ⚠️, 📝).
- **Links**: Gunakan relative links (`./file.md`) agar bisa dinavigasi di text editor.

### Diagram
Gunakan format text-based diagram untuk User Journeys agar mudah di-maintain (tidak perlu tool external).

```
[User] -> (Login Page) -> [API: /auth/login] -> (Dashboard)
```

---

## 🛠️ Maintenance Tools

### 1. Linting Documentation
Pastikan tidak ada broken links. Anda bisa cek manual dengan navigasi di VS Code / Cursor.

### 2. Verifikasi Langkah
Sebelum commit dokumentasi "How-to" (seperti `quickstart.md`), **COBA DULU** langkah-langkahnya dari awal (clean slate).
- Apakah command-nya jalan?
- Apakah ada step yang terlewat?

---

## 🔄 Lifecycle Example

**Skenario**: Anda menambahkan fitur "Forgot Password".

1. **Code**: Implementasi Backend & Frontend.
2. **Env**: Menambah `MAIL_HOST`, `MAIL_PORT` di `.env`.
   - 👉 Update `backend/.env.example`.
   - 👉 Update `docs/guides/configuration.md`.
3. **DB**: Menambah tabel `password_resets`.
   - 👉 (Optional) Update `docs/guides/database/models.md` jika ada relasi kompleks.
4. **API**: Endpoint `POST /forgot-password`.
   - 👉 Update `docs/development/api-documentation.md`.
5. **Docs**:
   - 👉 Tambah entry di `docs/CHANGELOG.md`.
   - 👉 Buat `docs/user-journeys/authentication/forgot-password.md` (Flow diagram).

---

## 🎯 Definition of Done (DoD) for Docs

Pull Request dianggap selesai jika:
- [ ] Fitur berjalan sesuai spec.
- [ ] Unit test passed.
- [ ] **Dokumentasi terkait sudah diupdate** (sesuai checklist di atas).
