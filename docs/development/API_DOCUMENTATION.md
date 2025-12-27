# 📚 API Documentation - Sirine Go App

Complete API reference untuk Sirine Go App.

**Base URL:** `http://localhost:8080` (development)

> **📖 Belajar customize API?** Lihat [CUSTOMIZATION_GUIDE.md](./CUSTOMIZATION_GUIDE.md)

---

## 🔍 Overview

API ini menggunakan RESTful architecture dengan JSON sebagai format data.

### Response Format

**Success Response:**
```json
{
  "data": {...},
  "message": "Pesan sukses (opsional)"
}
```

**Error Response:**
```json
{
  "error": "Pesan error dalam Bahasa Indonesia"
}
```

## 🏥 Health Check

### Check Server Status

```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "message": "Server berjalan dengan baik"
}
```

**Status Codes:**
- `200 OK` - Server berjalan normal

---

## 📝 Examples API

### Get All Examples

Mengambil semua data examples.

```http
GET /api/examples
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Example 1",
      "content": "Konten example 1",
      "is_active": true,
      "created_at": "2025-12-27T11:00:00+07:00",
      "updated_at": "2025-12-27T11:00:00+07:00"
    },
    {
      "id": 2,
      "title": "Example 2",
      "content": "Konten example 2",
      "is_active": false,
      "created_at": "2025-12-27T11:05:00+07:00",
      "updated_at": "2025-12-27T11:05:00+07:00"
    }
  ]
}
```

**Status Codes:**
- `200 OK` - Berhasil mengambil data
- `500 Internal Server Error` - Gagal mengambil data

---

### Get Example by ID

Mengambil detail example berdasarkan ID.

```http
GET /api/examples/:id
```

**Parameters:**
- `id` (path) - ID example (integer)

**Example Request:**
```bash
curl http://localhost:8080/api/examples/1
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "title": "Example 1",
    "content": "Konten example 1",
    "is_active": true,
    "created_at": "2025-12-27T11:00:00+07:00",
    "updated_at": "2025-12-27T11:00:00+07:00"
  }
}
```

**Status Codes:**
- `200 OK` - Berhasil mengambil data
- `400 Bad Request` - ID tidak valid
- `404 Not Found` - Data tidak ditemukan
- `500 Internal Server Error` - Gagal mengambil data

---

### Create Example

Membuat example baru.

```http
POST /api/examples
```

**Request Body:**
```json
{
  "title": "Example Baru",
  "content": "Ini adalah konten example baru",
  "is_active": true
}
```

**Field Descriptions:**
- `title` (string, required) - Judul example (max 255 karakter)
- `content` (string, optional) - Konten example (text)
- `is_active` (boolean, optional) - Status aktif (default: true)

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/examples \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Example Baru",
    "content": "Ini adalah konten example baru",
    "is_active": true
  }'
```

**Response:**
```json
{
  "message": "Data berhasil dibuat",
  "data": {
    "id": 3,
    "title": "Example Baru",
    "content": "Ini adalah konten example baru",
    "is_active": true,
    "created_at": "2025-12-27T11:10:00+07:00",
    "updated_at": "2025-12-27T11:10:00+07:00"
  }
}
```

**Status Codes:**
- `201 Created` - Berhasil membuat data
- `400 Bad Request` - Data tidak valid
- `500 Internal Server Error` - Gagal menyimpan data

---

### Update Example

Memperbarui example yang sudah ada.

```http
PUT /api/examples/:id
```

**Parameters:**
- `id` (path) - ID example (integer)

**Request Body:**
```json
{
  "title": "Example Updated",
  "content": "Konten yang sudah diperbarui",
  "is_active": false
}
```

**Example Request:**
```bash
curl -X PUT http://localhost:8080/api/examples/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Example Updated",
    "content": "Konten yang sudah diperbarui",
    "is_active": false
  }'
```

**Response:**
```json
{
  "message": "Data berhasil diperbarui"
}
```

**Status Codes:**
- `200 OK` - Berhasil memperbarui data
- `400 Bad Request` - ID atau data tidak valid
- `500 Internal Server Error` - Gagal memperbarui data

---

### Delete Example

Menghapus example (soft delete).

```http
DELETE /api/examples/:id
```

**Parameters:**
- `id` (path) - ID example (integer)

**Example Request:**
```bash
curl -X DELETE http://localhost:8080/api/examples/1
```

**Response:**
```json
{
  "message": "Data berhasil dihapus"
}
```

**Status Codes:**
- `200 OK` - Berhasil menghapus data
- `400 Bad Request` - ID tidak valid
- `500 Internal Server Error` - Gagal menghapus data

---

## 🔐 Authentication (Coming Soon)

Saat ini API belum menggunakan authentication. Untuk production, disarankan menambahkan:
- JWT authentication
- API key
- Rate limiting

---

## 📊 Error Codes

| Status Code | Meaning |
|------------|---------|
| 200 | OK - Request berhasil |
| 201 | Created - Resource berhasil dibuat |
| 400 | Bad Request - Request tidak valid |
| 401 | Unauthorized - Authentication required |
| 404 | Not Found - Resource tidak ditemukan |
| 500 | Internal Server Error - Server error |

---

## 🌐 CORS

API mendukung CORS untuk origins:
- `http://localhost:5173` (Development)
- `http://localhost:8080` (Production)

Allowed Methods: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`

---

## 📝 Notes

### Timezone
Semua timestamp menggunakan timezone **Asia/Jakarta (WIB, UTC+7)**.

### Soft Delete
Delete operation menggunakan soft delete. Data tidak benar-benar dihapus dari database, hanya ditandai dengan `deleted_at`.

### Validation
- `title` wajib diisi dan maksimal 255 karakter
- `content` opsional, tipe text
- `is_active` opsional, default `true`

### Pagination (Coming Soon)
Untuk implementasi pagination di masa depan:
```http
GET /api/examples?page=1&limit=10
```

### Filtering (Coming Soon)
Untuk implementasi filtering di masa depan:
```http
GET /api/examples?is_active=true
GET /api/examples?search=keyword
```

---

## 🧪 Testing dengan cURL

### Create, Read, Update, Delete Flow

```bash
# 1. Create
curl -X POST http://localhost:8080/api/examples \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","content":"Test content","is_active":true}'

# 2. Read All
curl http://localhost:8080/api/examples

# 3. Read One
curl http://localhost:8080/api/examples/1

# 4. Update
curl -X PUT http://localhost:8080/api/examples/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated","content":"Updated content","is_active":false}'

# 5. Delete
curl -X DELETE http://localhost:8080/api/examples/1
```

---

## 🔧 Development Tips

### Pretty Print JSON Response
```bash
curl http://localhost:8080/api/examples | jq
```

### Save Response to File
```bash
curl http://localhost:8080/api/examples > response.json
```

### Check Response Time
```bash
curl -w "\nTime: %{time_total}s\n" http://localhost:8080/api/examples
```

---

## 📚 Related Documentation

**Using the API:**
- **[CUSTOMIZATION_GUIDE.md](./CUSTOMIZATION_GUIDE.md)** - Add new endpoints
- **[TESTING.md](./TESTING.md)** - Test API endpoints
- **[FAQ.md](./FAQ.md)** - API troubleshooting

**Backend Architecture:**
- **[ARCHITECTURE_EXPLAINED.md](./ARCHITECTURE_EXPLAINED.md)** - Understand tech stack
- **[FOLDER_STRUCTURE.md](./FOLDER_STRUCTURE.md)** - Project structure

---

## 📞 Support

Jika ada pertanyaan tentang API:
- Developer: Zulfikar Hidayatullah
- Phone: +62 857-1583-8733
- Also check: [FAQ.md](./FAQ.md)

---

**Last Updated:** 27 Desember 2025  
**Version:** 1.0.0
