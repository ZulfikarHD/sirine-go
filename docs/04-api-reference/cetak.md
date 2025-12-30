# Cetak Queue API Documentation

**Feature**: Cetak Queue Management System  
**Version**: 1.0.0  
**Sprint**: Sprint 6  
**Last Updated**: 30 Desember 2025

---

## Overview

Cetak Queue API merupakan RESTful endpoints untuk Unit Cetak dalam mengakses antrian PO yang telah siap untuk dicetak setelah material preparation selesai di Unit Khazanah Awal, yaitu: queue retrieval dengan OBC Master information, detail PO dengan material photos, dan filtering berdasarkan priority.

**Base URL**: `http://localhost:8080/api`

---

## Authentication

Semua endpoint memerlukan authentication header:

```http
Authorization: Bearer <access_token>
```

**Required Roles:**
- Operator Cetak
- Supervisor Cetak
- Admin
- Manager

---

## Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/cetak/queue` | Retrieve queue PO siap cetak dengan pagination |
| GET | `/cetak/queue/:id` | Retrieve detail PO untuk cetak |

---

## API Endpoints

### GET /api/cetak/queue

Mengambil list Production Orders dengan status `READY_FOR_CETAK` yang telah selesai material preparation dan siap untuk proses cetak.

**Endpoint:** `GET /api/cetak/queue`

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | integer | No | 1 | Nomor halaman |
| `per_page` | integer | No | 20 | Jumlah item per halaman (max: 100) |
| `search` | string | No | - | Pencarian by PO number, OBC number, atau product name |
| `priority` | string | No | - | Filter by priority: URGENT, NORMAL, LOW |

**Success Response: `200 OK`**

```json
{
  "success": true,
  "message": "Queue cetak berhasil diambil",
  "data": {
    "items": [
      {
        "po_id": 1,
        "po_number": 123456,
        "obc_number": "OBC-2025-001",
        "product_name": "Buku Pemilik Kendaraan Bermotor",
        "priority": "URGENT",
        "priority_score": 100,
        "quantity": 5000,
        "due_date": "2025-01-15",
        "days_until_due": 15,
        "is_past_due": false,
        "material_ready_at": "2025-12-30 10:30:00",
        "prepared_by_id": 5,
        "prepared_by_name": "John Doe",
        "material_photos": [
          "https://storage.example.com/materials/photo1.jpg",
          "https://storage.example.com/materials/photo2.jpg"
        ],
        "notes": "Material sudah lengkap, siap cetak",
        "prep_id": 10,
        "obc_master": {
          "id": 1,
          "obc_number": "OBC-2025-001",
          "material": "BPKB",
          "material_description": "Buku Pemilik Kendaraan Bermotor",
          "seri": "Seri A",
          "warna": "Biru",
          "factory_code": "F001",
          "plat_number": "P001",
          "personalization": "Perso"
        }
      }
    ],
    "total": 45,
    "page": 1,
    "per_page": 20,
    "total_pages": 3
  }
}
```

**Error Response: `500 Internal Server Error`**

```json
{
  "success": false,
  "message": "Gagal mengambil queue cetak",
  "error": "database connection error"
}
```

**Example Request:**

```bash
# Get first page dengan default pagination
curl -X GET "http://localhost:8080/api/cetak/queue" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Search by OBC number
curl -X GET "http://localhost:8080/api/cetak/queue?search=OBC-2025-001" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Filter urgent priority only
curl -X GET "http://localhost:8080/api/cetak/queue?priority=URGENT" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Custom pagination
curl -X GET "http://localhost:8080/api/cetak/queue?page=2&per_page=50" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Business Rules:**
1. Hanya menampilkan PO dengan status `READY_FOR_CETAK`
2. Default sorting by `priority_score DESC, due_date ASC`
3. Past due items tetap muncul di queue (tidak auto-expire)
4. Material photos array bisa kosong jika staff tidak upload
5. OBC Master data preloaded untuk performance

---

### GET /api/cetak/queue/:id

Mengambil detail lengkap Production Order untuk proses cetak, termasuk full material preparation info, OBC Master specifications, dan material photos.

**Endpoint:** `GET /api/cetak/queue/:id`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Production Order ID |

**Success Response: `200 OK`**

```json
{
  "success": true,
  "message": "Detail PO berhasil diambil",
  "data": {
    "po_id": 1,
    "po_number": 123456,
    "obc_number": "OBC-2025-001",
    "sap_customer_code": "CUST001",
    "sap_product_code": "PROD-BPKB-001",
    "product_name": "Buku Pemilik Kendaraan Bermotor",
    "product_specifications": {
      "size": "A4",
      "pages": 20,
      "color": "Full Color"
    },
    "quantity_ordered": 5000,
    "quantity_target_lembar_besar": 625,
    "estimated_rims": 7,
    "order_date": "2025-12-15",
    "due_date": "2025-01-15",
    "priority": "URGENT",
    "priority_score": 100,
    "days_until_due": 15,
    "is_past_due": false,
    "current_status": "READY_FOR_CETAK",
    "notes": "",
    "obc_master": {
      "id": 1,
      "obc_number": "OBC-2025-001",
      "material": "BPKB",
      "material_description": "Buku Pemilik Kendaraan Bermotor",
      "seri": "Seri A",
      "warna": "Biru",
      "factory_code": "F001",
      "plat_number": "P001",
      "personalization": "Perso"
    },
    "material_prep": {
      "prep_id": 10,
      "status": "COMPLETED",
      "started_at": "2025-12-30 09:00:00",
      "completed_at": "2025-12-30 10:30:00",
      "duration_minutes": 90,
      "prepared_by_id": 5,
      "prepared_by_name": "John Doe",
      "sap_plat_code": "PLT-2025-001",
      "plat_retrieved_at": "2025-12-30 09:15:00",
      "kertas_blanko_quantity": 625,
      "kertas_blanko_actual": 630,
      "kertas_blanko_variance": 5,
      "tinta_requirements": {
        "cyan": 2,
        "magenta": 2,
        "yellow": 2,
        "black": 3
      },
      "tinta_actual": {
        "cyan": 2,
        "magenta": 2,
        "yellow": 2,
        "black": 3
      },
      "material_photos": [
        "https://storage.example.com/materials/photo1.jpg",
        "https://storage.example.com/materials/photo2.jpg",
        "https://storage.example.com/materials/photo3.jpg"
      ],
      "notes": "Semua material sudah dicek, siap cetak"
    }
  }
}
```

**Error Response: `404 Not Found`**

```json
{
  "success": false,
  "message": "Production Order tidak ditemukan"
}
```

**Error Response: `400 Bad Request`**

```json
{
  "success": false,
  "message": "PO tidak dalam status siap cetak"
}
```

**Example Request:**

```bash
# Get detail PO by ID
curl -X GET "http://localhost:8080/api/cetak/queue/1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Business Rules:**
1. PO harus dalam status `READY_FOR_CETAK`, jika tidak return 400
2. Material prep data harus exist dan completed
3. Material photos bisa empty array jika tidak ada upload
4. Variance kertas bisa positif (excess) atau negatif (shortage)
5. OBC Master specifications included untuk reference operator cetak

---

## Data Structures

### OBCMasterInfo

OBC Master information yang included di response untuk display specifications.

```typescript
interface OBCMasterInfo {
  id: number;
  obc_number: string;
  material: string;                  // e.g., "BPKB", "STNK"
  material_description: string;      // Full description
  seri: string;                      // e.g., "Seri A", "Seri B"
  warna: string;                     // e.g., "Biru", "Merah"
  factory_code: string;              // e.g., "F001"
  plat_number: string;               // e.g., "P001"
  personalization: string;           // "Perso", "Non-Perso", "-"
}
```

### MaterialPrepInfo

Material preparation information untuk reference operator cetak.

```typescript
interface MaterialPrepInfo {
  prep_id: number;
  status: string;                    // "COMPLETED"
  started_at: string;                // ISO 8601 datetime
  completed_at: string;              // ISO 8601 datetime
  duration_minutes: number;
  prepared_by_id: number;
  prepared_by_name: string;
  sap_plat_code: string;
  plat_retrieved_at: string;
  kertas_blanko_quantity: number;
  kertas_blanko_actual: number;
  kertas_blanko_variance: number;    // Can be positive (excess) or negative (shortage)
  tinta_requirements: Record<string, number>;
  tinta_actual: Record<string, number>;
  material_photos: string[];         // Array of photo URLs
  notes: string;
}
```

---

## Filtering & Sorting

### Default Sorting

Queue secara default di-sort berdasarkan:
1. **Priority Score** (DESC) - Urgent items first
2. **Due Date** (ASC) - Earlier deadlines first

### Priority Levels

| Priority | Score | Description |
|----------|-------|-------------|
| URGENT | 100 | High priority, biasanya near deadline atau critical |
| NORMAL | 50 | Standard priority |
| LOW | 25 | Low priority, flexible deadline |

### Search Behavior

Search parameter melakukan pencarian pada:
- PO Number (partial match)
- OBC Number (partial match)
- Product Name (partial match)

Case insensitive, menggunakan LIKE query.

---

## Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | INVALID_PARAMS | Query parameters tidak valid |
| 400 | NOT_READY_FOR_CETAK | PO tidak dalam status siap cetak |
| 401 | UNAUTHORIZED | Token tidak valid atau expired |
| 403 | FORBIDDEN | User tidak punya akses ke resource |
| 404 | NOT_FOUND | Production Order tidak ditemukan |
| 500 | INTERNAL_ERROR | Server error, database issue |

---

## Performance Considerations

### Database Optimization

1. **Preload Strategy:**
   - OBC Master data di-preload via `Preload("OBCMaster")`
   - Material prep data di-preload via `Preload("KhazwalMaterialPrep.PreparedByUser")`

2. **Indexing:**
   - Index pada `current_status` column untuk fast filtering
   - Index pada `priority_score` dan `due_date` untuk sorting

3. **Pagination:**
   - Default 20 items per page
   - Maximum 100 items per page untuk prevent large queries

### Response Size

- Average response size: ~5-10KB per item
- With 20 items: ~100-200KB per request
- Material photos tidak di-include di queue list, hanya URLs

---

## Related Documentation

- **Khazwal API:** [Khazwal Material Preparation API](./khazwal.md)
- **OBC Master API:** [OBC Master API](./obc-master.md)
- **Testing Guide:** [Cetak Testing](../06-testing/khazwal-sprint6-testing.md)
- **User Journeys:** [Cetak User Flow](../07-user-journeys/khazwal/cetak-queue-flow.md)
- **Sprint Documentation:** [Sprint 6](../10-sprints/sprint-khazwal-sprint6.md)

---

*Last Updated: 30 Desember 2025*
*Version: 1.0.0*
*Sprint: Sprint 6 - Consumer Side & Polish*
