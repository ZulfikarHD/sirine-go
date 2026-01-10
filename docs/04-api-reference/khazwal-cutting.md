# Khazwal Cutting API Reference

## Overview

API Khazwal Cutting merupakan endpoint untuk mengelola proses pemotongan lembar besar menjadi lembar kirim (sisiran kiri dan kanan) yang mencakup queue management, start process, input hasil, dan finalisasi dengan waste tracking.

**Base Path:** `/api/khazwal/cutting`  
**Authentication:** Required (JWT Token)  
**Roles:** STAFF_KHAZWAL, ADMIN, MANAGER  

---

## Endpoints Summary

| Method | Endpoint | Description | Sprint |
|--------|----------|-------------|--------|
| GET | `/queue` | Get cutting queue dengan filters | Sprint 1 ✅ |
| GET | `/:id` | Get cutting detail | Sprint 1 ✅ |
| POST | `/po/:po_id/start` | Start cutting process | Sprint 1 ✅ |
| PATCH | `/:id/result` | Update cutting result | Sprint 2 📋 |
| POST | `/:id/finalize` | Finalize cutting | Sprint 3 📋 |

---

## 1. Get Cutting Queue

Mengambil list PO yang siap untuk dipotong dengan filtering dan sorting.

### Request

```http
GET /api/khazwal/cutting/queue
Authorization: Bearer {token}
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `priority` | string | No | Filter by priority: URGENT, HIGH, NORMAL, LOW |
| `date_from` | string | No | Filter start date (YYYY-MM-DD) |
| `date_to` | string | No | Filter end date (YYYY-MM-DD) |
| `sort_by` | string | No | Sort field: `priority` (default), `date` |
| `sort_order` | string | No | Sort order: `asc`, `desc` |

### Response

**Success (200 OK):**

```json
{
  "data": [
    {
      "po_id": 123,
      "po_number": 2024001,
      "obc_number": "OBC-2024-001",
      "priority": "URGENT",
      "input_lembar_besar": 15000,
      "estimated_output": 30000,
      "counting_completed_at": "2024-01-10T10:30:00Z",
      "waiting_minutes": 45,
      "is_overdue": false
    }
  ],
  "meta": {
    "total": 5,
    "urgent_count": 2,
    "normal_count": 3
  }
}
```

**Business Logic:**
- PO dengan status `SIAP_POTONG` otomatis masuk queue
- Default sorting: Priority (URGENT first) + FIFO (oldest counting completion)
- `estimated_output` = `input_lembar_besar` × 2
- `is_overdue` = true jika waiting > 60 menit

---

## 2. Get Cutting Detail

Mengambil detail cutting record berdasarkan ID.

### Request

```http
GET /api/khazwal/cutting/:id
Authorization: Bearer {token}
```

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Cutting record ID |

### Response

**Success (200 OK):**

```json
{
  "id": 456,
  "production_order_id": 123,
  "status": "IN_PROGRESS",
  "started_at": "2024-01-10T11:00:00Z",
  "completed_at": null,
  "input_lembar_besar": 15000,
  "expected_output": 30000,
  "output_sisiran_kiri": null,
  "output_sisiran_kanan": null,
  "total_output": 0,
  "waste_quantity": 0,
  "waste_percentage": null,
  "waste_reason": "",
  "waste_photo_url": "",
  "cutting_machine": "Mesin A",
  "cut_by": {
    "id": 789,
    "name": "Siti Aminah",
    "nip": "12345"
  },
  "po": {
    "po_number": 2024001,
    "obc_number": "OBC-2024-001",
    "priority": "URGENT",
    "target_quantity": 15000
  }
}
```

**Error Responses:**

```json
// 404 Not Found
{
  "error": "Cutting record not found"
}
```

---

## 3. Start Cutting Process

Memulai proses pemotongan untuk PO tertentu.

### Request

```http
POST /api/khazwal/cutting/po/:po_id/start
Authorization: Bearer {token}
Content-Type: application/json
```

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `po_id` | integer | Production Order ID |

**Request Body:**

```json
{
  "cutting_machine": "Mesin A"
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `cutting_machine` | string | Yes | Must not be empty |

### Response

**Success (200 OK):**

```json
{
  "id": 456,
  "production_order_id": 123,
  "input_lembar_besar": 15000,
  "expected_output": 30000,
  "cutting_machine": "Mesin A",
  "status": "IN_PROGRESS",
  "started_at": "2024-01-10T11:00:00Z",
  "cut_by": 789
}
```

**Error Responses:**

```json
// 400 Bad Request - Invalid input
{
  "error": "Invalid request body",
  "details": "cutting_machine is required"
}

// 400 Bad Request - PO not ready
{
  "error": "PO is not ready for cutting"
}

// 404 Not Found
{
  "error": "Production order not found"
}

// 409 Conflict - Already started
{
  "error": "Cutting already started for this PO"
}

// 500 Internal Server Error
{
  "error": "Failed to start cutting",
  "details": "database connection error"
}
```

**Side Effects:**
1. Creates new record in `khazwal_cutting_results`
2. Updates `production_orders.current_status` = `'SEDANG_DIPOTONG'`
3. Updates `po_stage_tracking.started_at`
4. Creates activity log entry

---

## 4. Update Cutting Result (Sprint 2)

Mengupdate hasil pemotongan dengan input sisiran kiri dan kanan.

### Request

```http
PATCH /api/khazwal/cutting/:id/result
Authorization: Bearer {token}
Content-Type: application/json
```

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Cutting record ID |

**Request Body:**

```json
{
  "output_sisiran_kiri": 14800,
  "output_sisiran_kanan": 14900,
  "waste_reason": "Kertas robek saat proses cutting",
  "waste_photo_url": "/uploads/waste/photo123.jpg"
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `output_sisiran_kiri` | integer | Yes | Must be >= 0 |
| `output_sisiran_kanan` | integer | Yes | Must be >= 0 |
| `waste_reason` | string | Conditional | Required if waste > 2% |
| `waste_photo_url` | string | Conditional | Required if waste > 2% |

### Response

**Success (200 OK):**

```json
{
  "id": 456,
  "output_sisiran_kiri": 14800,
  "output_sisiran_kanan": 14900,
  "total_output": 29700,
  "expected_output": 30000,
  "waste_quantity": 300,
  "waste_percentage": 1.0,
  "waste_reason": "Kertas robek saat proses cutting",
  "waste_photo_url": "/uploads/waste/photo123.jpg"
}
```

**Auto-Calculations:**
- `total_output` = `output_sisiran_kiri` + `output_sisiran_kanan`
- `waste_quantity` = `expected_output` - `total_output`
- `waste_percentage` = (`waste_quantity` / `expected_output`) × 100

**Business Rules:**
- Waste threshold: 2%
- If waste > 2%: `waste_reason` and `waste_photo_url` REQUIRED

---

## 5. Finalize Cutting (Sprint 3)

Menyelesaikan proses pemotongan dan generate verification labels.

### Request

```http
POST /api/khazwal/cutting/:id/finalize
Authorization: Bearer {token}
```

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Cutting record ID |

**Request Body:** Empty (all data dari existing record)

### Response

**Success (200 OK):**

```json
{
  "id": 456,
  "status": "COMPLETED",
  "completed_at": "2024-01-10T13:30:00Z",
  "duration_minutes": 150,
  "labels_generated": 60
}
```

**Side Effects:**
1. Updates cutting status to `'COMPLETED'`
2. Calculates `duration_minutes`
3. Generates verification labels (1 label per 500 lembar)
4. Updates `production_orders` stage to `'VERIFIKASI'`
5. Updates `po_stage_tracking.completed_at`
6. Creates notification to Tim Verifikasi
7. Creates activity log entry

**Label Generation Logic:**
```
Total Labels = CEIL(total_output / 500)
Label Quantity = 500 (atau sisa untuk label terakhir)
Sisiran Pattern = Alternating (KIRI, KANAN, KIRI, KANAN, ...)
```

**Validation:**
- Cutting must be `IN_PROGRESS`
- `output_sisiran_kiri` and `output_sisiran_kanan` must be filled
- If waste > 2%: `waste_reason` and `waste_photo_url` must be filled

---

## Error Codes Reference

| Code | Message | Description |
|------|---------|-------------|
| 400 | Invalid request body | Request validation failed |
| 400 | PO is not ready for cutting | PO status bukan SIAP_POTONG |
| 400 | Cutting is not in progress | Trying to update/finalize non-active cutting |
| 400 | Output data must be filled | Finalize tanpa input hasil |
| 400 | Waste > 2% requires reason and photo | Missing waste documentation |
| 401 | User not authenticated | Token invalid/expired |
| 403 | Forbidden | User tidak punya akses (role tidak sesuai) |
| 404 | Production order not found | PO ID tidak ada |
| 404 | Cutting record not found | Cutting ID tidak ada |
| 409 | Cutting already started | Duplicate start attempt |
| 500 | Internal server error | Database/system error |

---

## Common Request Examples

### Example 1: Get URGENT priority queue

```bash
curl -X GET "http://localhost:8080/api/khazwal/cutting/queue?priority=URGENT" \
  -H "Authorization: Bearer eyJhbGc..."
```

### Example 2: Start cutting with machine selection

```bash
curl -X POST "http://localhost:8080/api/khazwal/cutting/po/123/start" \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{
    "cutting_machine": "Mesin B"
  }'
```

### Example 3: Get all PO completed today

```bash
curl -X GET "http://localhost:8080/api/khazwal/cutting/queue?date_from=2024-01-10&date_to=2024-01-10" \
  -H "Authorization: Bearer eyJhbGc..."
```

---

## Database Schema Reference

### Table: `khazwal_cutting_results`

```sql
CREATE TABLE khazwal_cutting_results (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    production_order_id BIGINT UNSIGNED NOT NULL UNIQUE,
    input_lembar_besar INT NOT NULL DEFAULT 0,
    expected_output INT NOT NULL DEFAULT 0,
    output_sisiran_kiri INT NULL,
    output_sisiran_kanan INT NULL,
    total_output INT NOT NULL DEFAULT 0,
    waste_quantity INT NOT NULL DEFAULT 0,
    waste_percentage DECIMAL(5,2) NULL,
    waste_reason TEXT,
    waste_photo_url VARCHAR(500),
    cutting_machine VARCHAR(100),
    cut_by BIGINT UNSIGNED NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    duration_minutes INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (production_order_id) REFERENCES production_orders(id),
    FOREIGN KEY (cut_by) REFERENCES users(id)
);
```

---

## Related Documentation

- **Sprint Documentation:** [Sprint Cutting Phase 1](../10-sprints/sprint-cutting-phase1.md)
- **Testing Guide:** [Cutting Testing](../06-testing/khazwal-cutting-testing.md)
- **User Manual:** TBD

---

**Last Updated:** 2026-01-10  
**API Version:** 1.0  
**Status:** Phase 1 Complete (Queue & Start), Phase 2-3 Planned
