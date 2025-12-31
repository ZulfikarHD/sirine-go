# API Documentation: Counting

## Overview

API untuk mengelola penghitungan hasil cetak (counting) yang mencakup queue management, input hasil penghitungan, breakdown kerusakan, dan finalisasi dengan tracking yang akurat untuk memastikan data quality dan audit trail.

**Base URL:** `https://api.sirinegolocalhost:8080/api`  
**API Prefix:** `/khazwal/counting`

## Authentication

Semua endpoint memerlukan header:
```
Authorization: Bearer <access_token>
```

**Required Roles:** `STAFF_KHAZWAL`, `ADMIN`, `MANAGER`

---

## Endpoints

### 1. Get Counting Queue

Mengambil list PO yang menunggu penghitungan dengan sorting FIFO berdasarkan print completion time.

**Endpoint:** `GET /khazwal/counting/queue`

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| machine_id | integer | No | - | Filter by machine ID |
| date_from | string | No | - | Filter by date (YYYY-MM-DD) |
| date_to | string | No | - | Filter by date (YYYY-MM-DD) |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "po_id": 123,
      "po_number": 1234567890,
      "obc_number": "OBC123456",
      "target_quantity": 500,
      "print_completed_at": "2025-12-30T14:30:00+07:00",
      "waiting_minutes": 125,
      "is_overdue": false,
      "machine": {
        "id": 5,
        "name": "MC-01",
        "code": "MC-01"
      },
      "operator": {
        "id": 10,
        "name": "John Doe",
        "nip": "12345"
      }
    }
  ],
  "meta": {
    "total": 5,
    "overdue_count": 2
  }
}
```

**Business Logic:**
- Query PO dengan `current_status = 'WAITING_COUNTING'`
- Join dengan `print_job_summaries` untuk data mesin & operator
- Calculate `waiting_minutes` = NOW - finalized_at
- Flag `is_overdue` jika waiting > 120 menit
- Sort by `finalized_at` ASC (FIFO)

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 500 | INTERNAL_ERROR | Gagal mengambil counting queue |

---

### 2. Get Counting Detail

Mengambil detail counting record dengan relasi PO dan print info untuk display atau resume.

**Endpoint:** `GET /khazwal/counting/:id`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Counting record ID |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": 456,
    "production_order_id": 123,
    "status": "IN_PROGRESS",
    "started_at": "2025-12-30T15:00:00+07:00",
    "quantity_good": 480,
    "quantity_defect": 15,
    "total_counted": 495,
    "variance_from_target": -5,
    "percentage_good": 96.97,
    "percentage_defect": 3.03,
    "defect_breakdown": [
      {"type": "Warna pudar", "quantity": 10},
      {"type": "Tinta blobor", "quantity": 5}
    ],
    "variance_reason": "Tumpah saat transport",
    "counted_by": {
      "id": 20,
      "name": "Jane Doe",
      "nip": "54321"
    },
    "po": {
      "po_number": 1234567890,
      "obc_number": "OBC123456",
      "target_quantity": 500
    },
    "print_info": {
      "machine_name": "MC-01",
      "operator_name": "John Doe",
      "finalized_at": "2025-12-30T14:30:00+07:00"
    }
  }
}
```

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 400 | INVALID_ID | ID tidak valid |
| 404 | NOT_FOUND | Counting record tidak ditemukan |

---

### 3. Start Counting

Memulai proses penghitungan untuk PO tertentu dengan create counting record dan update PO status.

**Endpoint:** `POST /khazwal/counting/:po_id/start`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| po_id | integer | Production Order ID |

**Request Body:**

```json
{}
```
*Empty body - user ID dari auth token*

**Response:** `201 Created`

```json
{
  "success": true,
  "message": "Penghitungan berhasil dimulai",
  "data": {
    "id": 456,
    "production_order_id": 123,
    "status": "IN_PROGRESS",
    "started_at": "2025-12-30T15:00:00+07:00",
    "counted_by": 20
  }
}
```

**Business Logic:**
1. Validate PO exists dan status = `WAITING_COUNTING`
2. Check tidak ada counting record lain dengan status IN_PROGRESS untuk PO ini
3. Insert ke `khazwal_counting_results` dengan status = 'IN_PROGRESS'
4. Update `production_orders.current_status = 'SEDANG_DIHITUNG'`
5. Update `po_stage_tracking` set `started_at`
6. Log to `activity_logs`

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 400 | INVALID_PO_ID | PO ID tidak valid |
| 400 | PO_NOT_READY | PO belum siap untuk penghitungan |
| 409 | COUNTING_EXISTS | Counting untuk PO ini sudah ada |
| 500 | INTERNAL_ERROR | Gagal memulai penghitungan |

---

### 4. Update Counting Result

Update hasil penghitungan (dapat dipanggil multiple times sebelum finalize) dengan auto-calculation.

**Endpoint:** `PATCH /khazwal/counting/:id/result`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Counting record ID |

**Request Body:**

```json
{
  "quantity_good": 480,
  "quantity_defect": 15,
  "defect_breakdown": [
    {"type": "Warna pudar", "quantity": 10},
    {"type": "Tinta blobor", "quantity": 5}
  ],
  "variance_reason": "Tumpah saat transport"
}
```

**Validation Rules:**

| Field | Rule | Error Message |
|-------|------|---------------|
| quantity_good | >= 0 | Jumlah baik harus >= 0 |
| quantity_defect | >= 0 | Jumlah rusak harus >= 0 |
| defect_breakdown | Required jika defect% > 5% | defect_breakdown wajib diisi karena persentase rusak > 5% |
| defect_breakdown sum | Must equal quantity_defect | Total defect_breakdown harus sama dengan quantity_defect |
| variance_reason | Required jika variance != 0 | variance_reason wajib diisi karena ada selisih dari target |
| status | Must be IN_PROGRESS | Counting tidak dalam status IN_PROGRESS |

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Hasil penghitungan berhasil disimpan",
  "data": {
    "id": 456,
    "quantity_good": 480,
    "quantity_defect": 15,
    "total_counted": 495,
    "variance_from_target": -5,
    "percentage_good": 96.97,
    "percentage_defect": 3.03,
    "defect_breakdown": [
      {"type": "Warna pudar", "quantity": 10},
      {"type": "Tinta blobor", "quantity": 5}
    ],
    "variance_reason": "Tumpah saat transport"
  }
}
```

**Auto-Calculated Fields:**
- `total_counted = quantity_good + quantity_defect`
- `variance_from_target = total_counted - target_quantity`
- `percentage_good = (quantity_good / total_counted) * 100`
- `percentage_defect = (quantity_defect / total_counted) * 100`

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 400 | INVALID_ID | ID tidak valid |
| 400 | COUNTING_NOT_IN_PROGRESS | Counting tidak dalam status IN_PROGRESS |
| 400 | COUNTING_COMPLETED | Counting sudah selesai dan tidak bisa diubah |
| 400 | INVALID_QUANTITY | quantity_good dan quantity_defect harus >= 0 |
| 422 | BREAKDOWN_REQUIRED | defect_breakdown wajib diisi karena persentase rusak > 5% |
| 422 | BREAKDOWN_SUM_MISMATCH | Total defect_breakdown harus sama dengan quantity_defect |
| 422 | VARIANCE_REASON_REQUIRED | variance_reason wajib diisi karena ada selisih dari target |
| 500 | INTERNAL_ERROR | Gagal menyimpan hasil penghitungan |

---

### 5. Finalize Counting

Menyelesaikan penghitungan dengan lock data dan advance PO ke next stage (SIAP_POTONG).

**Endpoint:** `POST /khazwal/counting/:id/finalize`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Counting record ID |

**Request Body:**

```json
{}
```
*Empty body - all data sudah ada di counting record*

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Penghitungan berhasil diselesaikan. PO siap untuk pemotongan.",
  "data": {
    "id": 456,
    "status": "COMPLETED",
    "completed_at": "2025-12-30T15:30:00+07:00",
    "duration_minutes": 30
  }
}
```

**Business Logic:**
1. Validate counting record exists dan status = IN_PROGRESS
2. Validate all required fields filled:
   - `quantity_good` and `quantity_defect` not NULL
   - If defect > 5%, `defect_breakdown` filled
   - If variance != 0, `variance_reason` filled
3. Update `khazwal_counting_results`:
   - `status = 'COMPLETED'`
   - `completed_at = NOW()`
   - `duration_minutes = (completed_at - started_at) in minutes`
4. Update `production_orders`:
   - `current_stage = 'KHAZWAL_CUTTING'`
   - `current_status = 'SIAP_POTONG'`
5. Update `po_stage_tracking` set `completed_at`
6. Create immutable `activity_logs` entry
7. Return success

**Transaction:** All steps wrapped in database transaction untuk atomicity.

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 400 | INVALID_ID | ID tidak valid |
| 400 | COUNTING_NOT_IN_PROGRESS | Counting tidak dalam status IN_PROGRESS |
| 400 | REQUIRED_FIELDS_MISSING | Field quantity_good dan quantity_defect wajib diisi |
| 400 | BREAKDOWN_REQUIRED | defect_breakdown wajib diisi karena persentase rusak > 5% |
| 400 | VARIANCE_REASON_REQUIRED | variance_reason wajib diisi karena ada selisih dari target |
| 422 | BREAKDOWN_SUM_MISMATCH | Total defect_breakdown harus sama dengan quantity_defect |
| 500 | INTERNAL_ERROR | Gagal menyelesaikan penghitungan |

---

## Data Structures

### DefectBreakdownItem

```typescript
interface DefectBreakdownItem {
  type: string;      // Jenis kerusakan
  quantity: number;  // Jumlah lembar
}
```

**Predefined Types:**
- Warna pudar
- Tinta blobor
- Kertas sobek
- Register tidak pas
- Lainnya

### Counting Status Flow

```
PENDING → IN_PROGRESS → COMPLETED
         ↑           ↓
         └─ (editable) ─┘
```

**State Rules:**
- `PENDING`: Initial state (not used in current flow)
- `IN_PROGRESS`: Data dapat di-update multiple times
- `COMPLETED`: Data immutable, tidak bisa di-update

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| INVALID_ID | 400 | ID tidak valid |
| INVALID_PO_ID | 400 | PO ID tidak valid |
| INVALID_QUANTITY | 400 | Quantity harus >= 0 |
| PO_NOT_READY | 400 | PO belum siap untuk penghitungan |
| COUNTING_NOT_IN_PROGRESS | 400 | Counting tidak dalam status IN_PROGRESS |
| COUNTING_COMPLETED | 400 | Counting sudah selesai |
| REQUIRED_FIELDS_MISSING | 400 | Field wajib belum diisi |
| COUNTING_EXISTS | 409 | Counting sudah ada |
| NOT_FOUND | 404 | Resource tidak ditemukan |
| BREAKDOWN_REQUIRED | 422 | Breakdown wajib diisi |
| BREAKDOWN_SUM_MISMATCH | 422 | Sum breakdown tidak match |
| VARIANCE_REASON_REQUIRED | 422 | Variance reason wajib diisi |
| INTERNAL_ERROR | 500 | Server error |

---

## Rate Limiting

- **Queue endpoint:** 10 requests/minute per user
- **Update result:** 30 requests/minute per user (untuk auto-save)
- **Other endpoints:** 20 requests/minute per user

---

## Webhooks (Future Enhancement)

*Not implemented in P0*

Potential webhooks untuk future enhancement:
- `counting.started` - Saat counting dimulai
- `counting.updated` - Saat result di-update
- `counting.finalized` - Saat counting diselesaikan
- `counting.overdue` - Saat waiting time > threshold

---

*Last Updated: 30 Desember 2025*
