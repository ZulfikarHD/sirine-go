# API Reference: Khazwal Counting

**Version:** 1.0.0  
**Base URL:** `/api/khazwal/counting`  
**Authentication:** Required (Bearer Token)  
**Roles:** `STAFF_KHAZWAL`, `ADMIN`, `MANAGER`

---

## Overview

API untuk mengelola proses penghitungan hasil cetak (Epic 2: Counting) yang bertujuan untuk tracking jumlah lembar besar hasil cetak dengan detail kerusakan, yaitu: queue management, input hasil penghitungan, breakdown kerusakan, dan finalisasi untuk advance ke stage pemotongan.

---

## Endpoints

### 1. Get Counting Queue

Mengambil list PO yang menunggu penghitungan dengan sorting FIFO berdasarkan waktu selesai cetak.

**Endpoint:** `GET /api/khazwal/counting/queue`

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `machine_id` | integer | No | Filter by mesin cetak |
| `date_from` | date | No | Filter dari tanggal (YYYY-MM-DD) |
| `date_to` | date | No | Filter sampai tanggal (YYYY-MM-DD) |

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
- Sort by `print_completed_at` ASC (FIFO)
- Flag `is_overdue` jika waiting > 120 menit
- Calculate `waiting_minutes` = NOW - print_completed_at

---

### 2. Get Counting Detail

Mengambil detail counting record dengan relasi PO dan print info.

**Endpoint:** `GET /api/khazwal/counting/:id`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Counting record ID |

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

**Error Response:** `404 Not Found`

```json
{
  "success": false,
  "message": "Counting record tidak ditemukan"
}
```

---

### 3. Start Counting

Memulai proses penghitungan untuk PO tertentu dengan create counting record dan update PO status.

**Endpoint:** `POST /api/khazwal/counting/:po_id/start`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `po_id` | integer | Production Order ID |

**Request Body:** Empty `{}`

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

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 400 | Bad Request | PO belum siap untuk penghitungan |
| 409 | Conflict | Counting sudah ada untuk PO ini |

**Business Logic:**
1. Validate PO status = `WAITING_COUNTING`
2. Check tidak ada counting IN_PROGRESS lain untuk PO ini
3. Create counting record dengan status `IN_PROGRESS`
4. Update PO status ke `SEDANG_DIHITUNG`
5. Update `po_stage_tracking.started_at` untuk stage `KHAZWAL_COUNTING`
6. Log activity

---

### 4. Update Counting Result

Update hasil penghitungan (dapat dipanggil multiple times sebelum finalize).

**Endpoint:** `PATCH /api/khazwal/counting/:id/result`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Counting record ID |

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

**Field Validation:**

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `quantity_good` | integer | Yes | >= 0 |
| `quantity_defect` | integer | Yes | >= 0 |
| `defect_breakdown` | array | Conditional | Required jika defect > 5%, sum must equal quantity_defect |
| `variance_reason` | string | Conditional | Required jika variance != 0 |

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

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 400 | Bad Request | Counting tidak dalam status IN_PROGRESS |
| 422 | Unprocessable Entity | Defect breakdown required (defect > 5%) |
| 422 | Unprocessable Entity | Defect breakdown sum mismatch |
| 422 | Unprocessable Entity | Variance reason required |

**Business Logic:**
- Status must be `IN_PROGRESS` (locked if `COMPLETED`)
- Auto-calculate: `total_counted`, `variance_from_target`, `percentage_good`, `percentage_defect`
- Validate defect breakdown requirement (> 5%)
- Validate variance reason requirement (variance != 0)

---

### 5. Finalize Counting

Finalize counting dengan lock data dan advance PO ke next stage (KHAZWAL_CUTTING).

**Endpoint:** `POST /api/khazwal/counting/:id/finalize`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Counting record ID |

**Request Body:** Empty `{}`

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

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 400 | Bad Request | Required fields not filled |
| 400 | Bad Request | Status not IN_PROGRESS |
| 422 | Unprocessable Entity | Defect breakdown required |
| 422 | Unprocessable Entity | Variance reason required |

**Business Logic:**
1. Validate status = `IN_PROGRESS`
2. Validate all required fields filled
3. Update counting: status = `COMPLETED`, set `completed_at`, calculate `duration_minutes`
4. Update PO: `current_stage` = `KHAZWAL_CUTTING`, `current_status` = `SIAP_POTONG`
5. Update `po_stage_tracking.completed_at` untuk stage `KHAZWAL_COUNTING`
6. Log activity (immutable)
7. Data locked (cannot be modified after finalize)

---

## Data Structures

### DefectBreakdownItem

```typescript
interface DefectBreakdownItem {
  type: string;      // Jenis kerusakan
  quantity: number;  // Jumlah lembar rusak
}
```

**Predefined Types:**
- "Warna pudar"
- "Tinta blobor"
- "Kertas sobek"
- "Register tidak pas"
- "Lainnya"

---

## Business Rules

| Rule ID | Description | Implementation |
|---------|-------------|----------------|
| BR-001 | Defect breakdown wajib jika rusak > 5% | Validate pada PATCH result dan POST finalize |
| BR-002 | Variance reason wajib jika variance != 0 | Validate pada PATCH result dan POST finalize |
| BR-003 | Defect breakdown sum must equal quantity_defect | Validate pada PATCH result |
| BR-004 | Data locked setelah finalize | Status check pada all update operations |
| BR-005 | FIFO queue sorting | Sort by print_completed_at ASC |
| BR-006 | Overdue flag jika > 120 menit | Flag pada queue response |
| BR-007 | Multiple PATCH allowed before finalize | No restriction pada PATCH calls |
| BR-008 | Single counting per PO | Check existing IN_PROGRESS on start |

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Request body tidak valid |
| `UNAUTHORIZED` | 401 | Token tidak valid atau expired |
| `FORBIDDEN` | 403 | Tidak punya akses ke resource |
| `RESOURCE_NOT_FOUND` | 404 | Counting record tidak ditemukan |
| `CONFLICT` | 409 | Counting already exists |
| `UNPROCESSABLE_ENTITY` | 422 | Business validation failed |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Related Documentation

- **Sprint Documentation:** [Sprint Counting Implementation](../10-sprints/sprint-counting-implementation.md)
- **Testing Guide:** [Counting Testing](../06-testing/khazwal-counting-testing.md)
- **User Journeys:** [Counting User Journeys](../07-user-journeys/khazwal-counting/)

---

*Last Updated: 2025-12-30*
