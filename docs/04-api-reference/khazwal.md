# Khazwal Material Preparation API Documentation

**Feature**: Khazwal Material Preparation System  
**Version**: 1.0.0  
**Last Updated**: 29 Desember 2025

---

## Overview

Khazwal Material Preparation API merupakan RESTful endpoints untuk mengelola workflow penyiapan material di Unit Khazanah Awal, yaitu: queue management, proses konfirmasi plat/kertas/tinta, finalisasi persiapan, serta monitoring dan history tracking.

**Base URL**: `http://localhost:8080/api`

---

## Endpoints Summary

### Material Preparation (Staff Khazwal)

| Method | Endpoint | Auth Required | Roles | Description |
|--------|----------|---------------|-------|-------------|
| GET | `/khazwal/material-prep/queue` | ✅ Yes | Staff Khazwal, Admin, Manager | Queue PO menunggu persiapan |
| GET | `/khazwal/material-prep/:id` | ✅ Yes | Staff Khazwal, Admin, Manager | Detail PO untuk persiapan |
| POST | `/khazwal/material-prep/:id/start` | ✅ Yes | Staff Khazwal, Admin, Manager | Mulai proses persiapan |
| POST | `/khazwal/material-prep/:id/confirm-plat` | ✅ Yes | Staff Khazwal, Admin, Manager | Konfirmasi pengambilan plat |
| PATCH | `/khazwal/material-prep/:id/kertas` | ✅ Yes | Staff Khazwal, Admin, Manager | Update jumlah kertas blanko |
| PATCH | `/khazwal/material-prep/:id/tinta` | ✅ Yes | Staff Khazwal, Admin, Manager | Update checklist tinta |
| POST | `/khazwal/material-prep/:id/finalize` | ✅ Yes | Staff Khazwal, Admin, Manager | Selesaikan persiapan |
| GET | `/khazwal/material-prep/history` | ✅ Yes | Staff Khazwal, Admin, Manager | Riwayat persiapan |

### Monitoring (Supervisor)

| Method | Endpoint | Auth Required | Roles | Description |
|--------|----------|---------------|-------|-------------|
| GET | `/khazwal/monitoring` | ✅ Yes | Supervisor Khazwal, Admin, Manager | Dashboard monitoring |

### Cetak Queue (Unit Cetak)

| Method | Endpoint | Auth Required | Roles | Description |
|--------|----------|---------------|-------|-------------|
| GET | `/cetak/queue` | ✅ Yes | Operator Cetak, Supervisor Cetak, Admin, Manager | Queue PO siap cetak |
| GET | `/cetak/queue/:id` | ✅ Yes | Operator Cetak, Supervisor Cetak, Admin, Manager | Detail PO untuk cetak |

---

## Material Preparation APIs

### GET /api/khazwal/material-prep/queue

**Description**: Mendapatkan daftar Production Order (PO) yang menunggu atau sedang dalam proses persiapan material.

**Authentication**: Required (Bearer token)  
**Authorization**: STAFF_KHAZWAL, ADMIN, MANAGER

#### Request

**Headers**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters**:

| Field | Type | Required | Description | Default |
|-------|------|----------|-------------|---------|
| `page` | integer | ❌ No | Nomor halaman | 1 |
| `per_page` | integer | ❌ No | Items per halaman (max 50) | 10 |
| `status` | string | ❌ No | Filter status (WAITING_MATERIAL_PREP, MATERIAL_PREP_IN_PROGRESS) | - |
| `priority` | string | ❌ No | Filter prioritas (URGENT, NORMAL, LOW) | - |
| `search` | string | ❌ No | Search by PO number, OBC number, atau product name | - |

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "po_number": 1001234567,
        "obc_number": "ABC123456",
        "product_name": "Kertas HVS A4 80gsm",
        "quantity_ordered": 1000,
        "quantity_target_lembar_besar": 500,
        "priority": "URGENT",
        "priority_score": 95,
        "due_date": "2025-12-30",
        "days_until_due": 1,
        "current_status": "WAITING_MATERIAL_PREP",
        "material_prep": {
          "id": 1,
          "status": "PENDING",
          "sap_plat_code": "PLAT-001",
          "kertas_blanko_quantity": 500
        }
      }
    ],
    "total": 8,
    "page": 1,
    "per_page": 10,
    "total_pages": 1
  }
}
```

---

### GET /api/khazwal/material-prep/:id

**Description**: Mendapatkan detail lengkap PO untuk proses persiapan material.

**Authentication**: Required (Bearer token)  
**Authorization**: STAFF_KHAZWAL, ADMIN, MANAGER

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "po_number": 1001234567,
    "obc_number": "ABC123456",
    "sap_customer_code": "CUST-001",
    "sap_product_code": "PROD-001",
    "product_name": "Kertas HVS A4 80gsm",
    "product_specifications": {
      "paper_type": "HVS",
      "gsm": 80,
      "size": "A4"
    },
    "quantity_ordered": 1000,
    "quantity_target_lembar_besar": 500,
    "estimated_rims": 10,
    "order_date": "2025-12-20",
    "due_date": "2025-12-30",
    "priority": "URGENT",
    "priority_score": 95,
    "current_stage": "KHAZWAL_MATERIAL_PREP",
    "current_status": "WAITING_MATERIAL_PREP",
    "notes": "Urgent order untuk customer VIP",
    "material_prep": {
      "id": 1,
      "status": "PENDING",
      "sap_plat_code": "PLAT-001",
      "kertas_blanko_quantity": 500,
      "tinta_requirements": [
        {"color": "cyan", "quantity": 5},
        {"color": "magenta", "quantity": 5},
        {"color": "yellow", "quantity": 5},
        {"color": "black", "quantity": 10}
      ],
      "plat_retrieved_at": null,
      "kertas_blanko_actual": null,
      "tinta_actual": null,
      "started_at": null,
      "completed_at": null,
      "prepared_by": null
    },
    "stage_tracking": []
  }
}
```

**Error (404 Not Found)**:
```json
{
  "success": false,
  "message": "Production Order tidak ditemukan"
}
```

---

### POST /api/khazwal/material-prep/:id/start

**Description**: Memulai proses persiapan material untuk PO tertentu.

**Authentication**: Required (Bearer token)  
**Authorization**: STAFF_KHAZWAL, ADMIN, MANAGER

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "message": "Persiapan material berhasil dimulai",
  "data": {
    "id": 1,
    "po_number": 1001234567,
    "current_status": "MATERIAL_PREP_IN_PROGRESS",
    "material_prep": {
      "status": "IN_PROGRESS",
      "started_at": "2025-12-29T10:00:00+07:00",
      "prepared_by": 5
    }
  }
}
```

**Error (400 Bad Request)**:
```json
{
  "success": false,
  "message": "PO tidak dapat diproses. Status harus WAITING_MATERIAL_PREP"
}
```

---

### POST /api/khazwal/material-prep/:id/confirm-plat

**Description**: Konfirmasi pengambilan plat cetak dengan validasi barcode.

**Authentication**: Required (Bearer token)  
**Authorization**: STAFF_KHAZWAL, ADMIN, MANAGER

#### Request Body

```json
{
  "plat_code": "PLAT-001"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `plat_code` | string | ✅ Yes | Kode plat hasil scan barcode |

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "message": "Plat berhasil dikonfirmasi",
  "data": {
    "plat_retrieved_at": "2025-12-29T10:05:00+07:00"
  }
}
```

**Error (400 Bad Request - Kode tidak cocok)**:
```json
{
  "success": false,
  "message": "Kode plat tidak sesuai. Expected: PLAT-001, Got: PLAT-002"
}
```

---

### PATCH /api/khazwal/material-prep/:id/kertas

**Description**: Update jumlah kertas blanko aktual yang disiapkan.

**Authentication**: Required (Bearer token)  
**Authorization**: STAFF_KHAZWAL, ADMIN, MANAGER

#### Request Body

```json
{
  "actual_qty": 480,
  "variance_reason": "Beberapa lembar rusak saat pengecekan"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `actual_qty` | integer | ✅ Yes | Jumlah kertas aktual (min: 0) |
| `variance_reason` | string | Conditional | Wajib jika variance > 5% |

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "message": "Jumlah kertas berhasil diupdate",
  "data": {
    "kertas_blanko_actual": 480,
    "kertas_blanko_variance": -20,
    "variance_percentage": -4.0
  }
}
```

**Error (400 Bad Request - Variance tanpa alasan)**:
```json
{
  "success": false,
  "message": "Variance lebih dari 5%, alasan wajib diisi"
}
```

---

### PATCH /api/khazwal/material-prep/:id/tinta

**Description**: Update checklist dan jumlah tinta yang disiapkan.

**Authentication**: Required (Bearer token)  
**Authorization**: STAFF_KHAZWAL, ADMIN, MANAGER

#### Request Body

```json
{
  "tinta_actual": [
    {"color": "cyan", "quantity": 5.5, "checked": true},
    {"color": "magenta", "quantity": 5.0, "checked": true},
    {"color": "yellow", "quantity": 5.0, "checked": true},
    {"color": "black", "quantity": 10.5, "checked": true}
  ]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `tinta_actual` | array | ✅ Yes | Array of tinta items |
| `tinta_actual[].color` | string | ✅ Yes | Nama warna tinta |
| `tinta_actual[].quantity` | float | ✅ Yes | Jumlah dalam kg (min: 0) |
| `tinta_actual[].checked` | boolean | ✅ Yes | Status checklist |

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "message": "Checklist tinta berhasil diupdate"
}
```

---

### POST /api/khazwal/material-prep/:id/finalize

**Description**: Menyelesaikan proses persiapan material dan mengirim notifikasi ke Unit Cetak.

**Authentication**: Required (Bearer token)  
**Authorization**: STAFF_KHAZWAL, ADMIN, MANAGER

#### Request Body

```json
{
  "photos": ["base64_string_1", "base64_string_2"],
  "notes": "Semua material sudah disiapkan dengan baik"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `photos` | array | ❌ No | Array base64 string foto (max 5) |
| `notes` | string | ❌ No | Catatan tambahan |

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "message": "Material preparation berhasil diselesaikan",
  "data": {
    "prep_id": 1,
    "po_number": 1001234567,
    "obc_number": "ABC123456",
    "duration_minutes": 45,
    "completed_at": "2025-12-29T10:45:00+07:00",
    "prepared_by_name": "John Doe",
    "photos_count": 2
  }
}
```

**Error (400 Bad Request - Langkah belum selesai)**:
```json
{
  "success": false,
  "message": "Material preparation tidak dapat diselesaikan. Pastikan semua langkah (plat, kertas, tinta) sudah selesai."
}
```

---

### GET /api/khazwal/material-prep/history

**Description**: Mendapatkan riwayat persiapan material yang sudah selesai.

**Authentication**: Required (Bearer token)  
**Authorization**: STAFF_KHAZWAL, ADMIN, MANAGER

#### Query Parameters

| Field | Type | Required | Description | Default |
|-------|------|----------|-------------|---------|
| `page` | integer | ❌ No | Nomor halaman | 1 |
| `per_page` | integer | ❌ No | Items per halaman | 10 |
| `start_date` | string | ❌ No | Filter tanggal mulai (YYYY-MM-DD) | - |
| `end_date` | string | ❌ No | Filter tanggal akhir (YYYY-MM-DD) | - |
| `search` | string | ❌ No | Search by PO/OBC number | - |

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "po_number": 1001234567,
        "obc_number": "ABC123456",
        "product_name": "Kertas HVS A4 80gsm",
        "completed_at": "2025-12-29T10:45:00+07:00",
        "duration_minutes": 45,
        "prepared_by_name": "John Doe"
      }
    ],
    "total": 15,
    "page": 1,
    "per_page": 10,
    "total_pages": 2
  }
}
```

---

### GET /api/khazwal/monitoring

**Description**: Dashboard monitoring untuk Supervisor Khazwal.

**Authentication**: Required (Bearer token)  
**Authorization**: SUPERVISOR_KHAZWAL, ADMIN, MANAGER

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "data": {
    "stats": {
      "queue_count": 5,
      "in_progress_count": 2,
      "completed_today": 8,
      "avg_duration_minutes": 42
    },
    "staff_activity": [
      {
        "user_id": 5,
        "user_name": "John Doe",
        "status": "active",
        "current_po": "ABC123456",
        "started_at": "2025-12-29T10:00:00+07:00"
      }
    ],
    "recent_completions": [
      {
        "po_number": 1001234567,
        "obc_number": "ABC123456",
        "completed_at": "2025-12-29T10:45:00+07:00",
        "prepared_by": "John Doe",
        "duration_minutes": 45
      }
    ]
  }
}
```

---

## Cetak Queue APIs

### GET /api/cetak/queue

**Description**: Queue PO yang sudah siap untuk proses cetak.

**Authentication**: Required (Bearer token)  
**Authorization**: OPERATOR_CETAK, SUPERVISOR_CETAK, ADMIN, MANAGER

#### Response

**Success (200 OK)**:
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "po_number": 1001234567,
        "obc_number": "ABC123456",
        "product_name": "Kertas HVS A4 80gsm",
        "priority": "URGENT",
        "current_status": "READY_FOR_CETAK",
        "material_ready_at": "2025-12-29T10:45:00+07:00",
        "prepared_by_name": "John Doe",
        "photos_count": 2
      }
    ],
    "total": 3,
    "page": 1,
    "per_page": 10,
    "total_pages": 1
  }
}
```

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| VALIDATION_ERROR | 400 | Request body tidak valid |
| UNAUTHORIZED | 401 | Token tidak valid atau expired |
| FORBIDDEN | 403 | Tidak punya akses ke resource |
| NOT_FOUND | 404 | Resource tidak ditemukan |
| INVALID_STATUS | 400 | Status PO tidak sesuai untuk aksi |
| INCOMPLETE_STEPS | 400 | Langkah persiapan belum lengkap |
| INTERNAL_ERROR | 500 | Server error |

---

## Related Documentation

- **Testing**: [Khazwal Testing Guide](../06-testing/khazwal-testing.md)
- **User Journeys**: [Khazwal User Journeys](../07-user-journeys/khazwal/material-prep-flow.md)
- **Sprint Documentation**: [Sprint Khazwal](../10-sprints/sprint-khazwal-material-prep.md)

---

*Last Updated: 29 Desember 2025*
