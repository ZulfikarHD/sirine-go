# Khazwal Material Preparation Testing Guide

**Feature**: Khazwal Material Preparation System  
**Version**: 1.0.0  
**Last Updated**: 29 Desember 2025

---

## Overview

Testing guide ini mencakup comprehensive test scenarios untuk Khazwal Material Preparation features, yaitu: backend API testing, frontend integration testing, end-to-end user flows, dan workflow validation.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Backend API Testing](#backend-api-testing)
3. [Frontend Component Testing](#frontend-component-testing)
4. [End-to-End Testing](#end-to-end-testing)
5. [Workflow Testing](#workflow-testing)
6. [Performance Testing](#performance-testing)

---

## Prerequisites

### Environment Setup

**Backend**:
```bash
cd backend
go run cmd/server/main.go
```

**Frontend**:
```bash
cd frontend
yarn dev
```

### Test Data

**Staff Khazwal Account**:
- NIP: `20001`
- Password: `Demo@123`
- Role: STAFF_KHAZWAL

**Supervisor Khazwal Account**:
- NIP: `30001`
- Password: `Demo@123`
- Role: SUPERVISOR_KHAZWAL

**Operator Cetak Account**:
- NIP: `20002`
- Password: `Demo@123`
- Role: OPERATOR_CETAK

### Seed Data

Jalankan seeder untuk sample Production Orders:
```bash
cd backend && go run cmd/seed/main.go
```

---

## Backend API Testing

### 1. Queue Endpoint

#### 1.1 GET /api/khazwal/material-prep/queue - List Queue

**Test Case 1: Get queue dengan default pagination**

```bash
# Setup - Login sebagai Staff Khazwal
TOKEN=$(curl -s http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"nip":"20001","password":"Demo@123"}' \
  | jq -r '.data.access_token')

# Test
curl http://localhost:8080/api/khazwal/material-prep/queue \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected Result**:
```json
{
  "success": true,
  "data": {
    "items": [...],
    "total": 8,
    "page": 1,
    "per_page": 10,
    "total_pages": 1
  }
}
```

**Test Case 2: Filter by priority**

```bash
curl "http://localhost:8080/api/khazwal/material-prep/queue?priority=URGENT" \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Test Case 3: Search by OBC number**

```bash
curl "http://localhost:8080/api/khazwal/material-prep/queue?search=ABC123" \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Test Case 4: Unauthorized access**

```bash
curl http://localhost:8080/api/khazwal/material-prep/queue \
  | jq
```

**Expected**: 401 Unauthorized

---

### 2. Start Preparation

#### 2.1 POST /api/khazwal/material-prep/:id/start

**Test Case 1: Start preparation (valid)**

```bash
curl -X POST http://localhost:8080/api/khazwal/material-prep/1/start \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected Result**:
```json
{
  "success": true,
  "message": "Persiapan material berhasil dimulai"
}
```

**Test Case 2: Start preparation (already in progress)**

```bash
# Jalankan request yang sama lagi
curl -X POST http://localhost:8080/api/khazwal/material-prep/1/start \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected**: 400 Bad Request - Status tidak valid

**Test Case 3: Start preparation (invalid ID)**

```bash
curl -X POST http://localhost:8080/api/khazwal/material-prep/9999/start \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected**: 404 Not Found

---

### 3. Confirm Plat

#### 3.1 POST /api/khazwal/material-prep/:id/confirm-plat

**Test Case 1: Confirm plat (valid code)**

```bash
curl -X POST http://localhost:8080/api/khazwal/material-prep/1/confirm-plat \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"plat_code":"PLAT-001"}' \
  | jq
```

**Test Case 2: Confirm plat (invalid code)**

```bash
curl -X POST http://localhost:8080/api/khazwal/material-prep/1/confirm-plat \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"plat_code":"WRONG-CODE"}' \
  | jq
```

**Expected**: 400 Bad Request - Kode tidak sesuai

---

### 4. Update Kertas

#### 4.1 PATCH /api/khazwal/material-prep/:id/kertas

**Test Case 1: Update kertas (no variance)**

```bash
curl -X PATCH http://localhost:8080/api/khazwal/material-prep/1/kertas \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"actual_qty":500}' \
  | jq
```

**Test Case 2: Update kertas (variance > 5% with reason)**

```bash
curl -X PATCH http://localhost:8080/api/khazwal/material-prep/1/kertas \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"actual_qty":450,"variance_reason":"Beberapa lembar rusak"}' \
  | jq
```

**Test Case 3: Update kertas (variance > 5% without reason)**

```bash
curl -X PATCH http://localhost:8080/api/khazwal/material-prep/1/kertas \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"actual_qty":400}' \
  | jq
```

**Expected**: 400 Bad Request - Alasan wajib diisi

---

### 5. Update Tinta

#### 5.1 PATCH /api/khazwal/material-prep/:id/tinta

**Test Case 1: Update tinta (all checked)**

```bash
curl -X PATCH http://localhost:8080/api/khazwal/material-prep/1/tinta \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "tinta_actual": [
      {"color":"cyan","quantity":5.5,"checked":true},
      {"color":"magenta","quantity":5.0,"checked":true},
      {"color":"yellow","quantity":5.0,"checked":true},
      {"color":"black","quantity":10.5,"checked":true}
    ]
  }' \
  | jq
```

---

### 6. Finalize

#### 6.1 POST /api/khazwal/material-prep/:id/finalize

**Test Case 1: Finalize (all steps complete)**

```bash
curl -X POST http://localhost:8080/api/khazwal/material-prep/1/finalize \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"photos":[],"notes":"Semua material siap"}' \
  | jq
```

**Expected Result**:
```json
{
  "success": true,
  "message": "Material preparation berhasil diselesaikan",
  "data": {
    "prep_id": 1,
    "duration_minutes": 45,
    "completed_at": "..."
  }
}
```

**Test Case 2: Finalize (steps incomplete)**

```bash
# Gunakan PO baru yang belum diproses
curl -X POST http://localhost:8080/api/khazwal/material-prep/2/finalize \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"photos":[],"notes":""}' \
  | jq
```

**Expected**: 400 Bad Request - Langkah belum lengkap

---

## Frontend Component Testing

### 1. MaterialPrepQueuePage

| Test ID | Scenario | Expected Result | Status |
|---------|----------|-----------------|--------|
| FE-Q-01 | Load page | Queue list displayed dengan PO cards | ☐ |
| FE-Q-02 | Search by PO number | Filter results correctly | ☐ |
| FE-Q-03 | Filter by priority | Show only matching priority | ☐ |
| FE-Q-04 | Click PO card | Navigate to detail page | ☐ |
| FE-Q-05 | Empty state | Show "Belum Ada PO" message | ☐ |
| FE-Q-06 | Loading state | Show skeleton cards | ☐ |
| FE-Q-07 | Pagination | Navigate between pages | ☐ |
| FE-Q-08 | Mobile responsive | Layout adapts correctly | ☐ |

### 2. MaterialPrepDetailPage

| Test ID | Scenario | Expected Result | Status |
|---------|----------|-----------------|--------|
| FE-D-01 | Load detail | PO info displayed correctly | ☐ |
| FE-D-02 | Start preparation | Show confirmation dialog | ☐ |
| FE-D-03 | Confirm start | Status changes, navigate to process | ☐ |
| FE-D-04 | Due date past | Red color indicator shown | ☐ |
| FE-D-05 | Due date near | Yellow color indicator shown | ☐ |
| FE-D-06 | Back button | Return to queue page | ☐ |

### 3. MaterialPrepProcessPage

| Test ID | Scenario | Expected Result | Status |
|---------|----------|-----------------|--------|
| FE-P-01 | Load process | Stepper shows 4 steps | ☐ |
| FE-P-02 | Step 1: Scan valid barcode | Green checkmark, advance to step 2 | ☐ |
| FE-P-03 | Step 1: Scan invalid barcode | Show mismatch warning | ☐ |
| FE-P-04 | Step 1: Manual input | Accept manual code entry | ☐ |
| FE-P-05 | Step 2: Input kertas | Real-time variance calculation | ☐ |
| FE-P-06 | Step 2: Variance > 5% | Require reason input | ☐ |
| FE-P-07 | Step 3: Check all tinta | Enable next button | ☐ |
| FE-P-08 | Step 3: Low stock warning | Show alert for < 10kg | ☐ |
| FE-P-09 | Step 4: Upload photo | Preview thumbnails shown | ☐ |
| FE-P-10 | Step 4: Finalize | Success screen with duration | ☐ |
| FE-P-11 | Navigation back | Previous step accessible | ☐ |

### 4. SupervisorMonitoringPage

| Test ID | Scenario | Expected Result | Status |
|---------|----------|-----------------|--------|
| FE-M-01 | Load dashboard | Stats cards displayed | ☐ |
| FE-M-02 | Staff activity | Active staff shown | ☐ |
| FE-M-03 | Recent completions | Timeline displayed | ☐ |
| FE-M-04 | Auto-refresh | Data updates every 30s | ☐ |

### 5. CetakQueuePage

| Test ID | Scenario | Expected Result | Status |
|---------|----------|-----------------|--------|
| FE-C-01 | Load queue | Ready POs displayed | ☐ |
| FE-C-02 | View photos | Photo viewer opens | ☐ |
| FE-C-03 | Filter/search | Results filtered correctly | ☐ |

---

## End-to-End Testing

### Complete Workflow Test

**Scenario**: Staff Khazwal menyelesaikan 1 PO dari queue sampai finalize

**Steps**:

1. **Login sebagai Staff Khazwal**
   - Navigate to `/login`
   - Enter NIP: `20001`, Password: `Demo@123`
   - ✅ Expected: Redirect to dashboard

2. **Navigate to Material Prep Queue**
   - Click menu "Khazanah Awal" → "Persiapan Material"
   - ✅ Expected: Queue page dengan PO list

3. **Select PO dari Queue**
   - Click PO card dengan priority URGENT
   - ✅ Expected: Detail page dengan PO info

4. **Start Preparation**
   - Click "Mulai Persiapan"
   - Confirm dialog
   - ✅ Expected: Navigate to process page

5. **Step 1: Confirm Plat**
   - Click "Scan Plat"
   - Scan atau manual input kode plat
   - ✅ Expected: Green checkmark, auto advance

6. **Step 2: Input Kertas**
   - Input actual quantity
   - (Optional) Input reason jika variance > 5%
   - Click "Simpan"
   - ✅ Expected: Auto advance to step 3

7. **Step 3: Checklist Tinta**
   - Check semua warna tinta
   - Input quantity per warna
   - Click "Simpan"
   - ✅ Expected: Auto advance to step 4

8. **Step 4: Finalize**
   - (Optional) Upload photos
   - (Optional) Add notes
   - Click "Selesai & Kirim ke Unit Cetak"
   - ✅ Expected: Success screen dengan duration

9. **Verify Notification**
   - Login sebagai Operator Cetak
   - Check notification bell
   - ✅ Expected: "Material Siap - PO #..." notification

10. **Verify Cetak Queue**
    - Navigate to "Antrian Cetak"
    - ✅ Expected: PO muncul dengan status READY_FOR_CETAK

---

## Workflow Testing

### Negative Test Cases

| Test ID | Scenario | Expected Result |
|---------|----------|-----------------|
| NEG-01 | Start prep tanpa login | 401 Unauthorized |
| NEG-02 | Start prep dengan role USER | 403 Forbidden |
| NEG-03 | Start prep PO yang sudah started | 400 Bad Request |
| NEG-04 | Confirm plat dengan kode salah | 400 Bad Request |
| NEG-05 | Update kertas dengan variance tanpa alasan | 400 Bad Request |
| NEG-06 | Finalize tanpa complete steps | 400 Bad Request |
| NEG-07 | Access monitoring sebagai Staff | 403 Forbidden |

### Concurrent Access Test

| Test ID | Scenario | Expected Result |
|---------|----------|-----------------|
| CON-01 | 2 staff start same PO | First wins, second gets error |
| CON-02 | Staff A processing, Staff B tries to start | Error: sedang diproses |

---

## Performance Testing

### Load Test Scenarios

```bash
# Install k6
# brew install k6

# Run load test untuk queue endpoint
k6 run -u 10 -d 30s scripts/load-test-queue.js
```

### Performance Benchmarks

| Endpoint | Target | Acceptable |
|----------|--------|------------|
| GET /queue | < 200ms | < 500ms |
| GET /detail | < 100ms | < 300ms |
| POST /start | < 300ms | < 500ms |
| POST /finalize | < 500ms | < 1000ms |

---

## Checklist Summary

### Backend Testing ☐

- [ ] Queue endpoint (filters, pagination, search)
- [ ] Detail endpoint (relations, validation)
- [ ] Start preparation (status validation, transaction)
- [ ] Confirm plat (barcode validation)
- [ ] Update kertas (variance logic)
- [ ] Update tinta (array handling)
- [ ] Finalize (step validation, notification)
- [ ] History endpoint
- [ ] Monitoring endpoint (role-based)
- [ ] Cetak queue endpoint

### Frontend Testing ☐

- [ ] MaterialPrepQueuePage (all scenarios)
- [ ] MaterialPrepDetailPage (all scenarios)
- [ ] MaterialPrepProcessPage (all steps)
- [ ] BarcodeScanner component
- [ ] KertasInputForm component
- [ ] TintaChecklist component
- [ ] PhotoUploader component
- [ ] SupervisorMonitoringPage
- [ ] CetakQueuePage
- [ ] Mobile responsive semua pages

### E2E Testing ☐

- [ ] Complete workflow Staff Khazwal
- [ ] Supervisor monitoring flow
- [ ] Operator Cetak receive notification flow
- [ ] Error scenarios
- [ ] Concurrent access scenarios

---

## Related Documentation

- **API Reference**: [Khazwal API](../04-api-reference/khazwal.md)
- **User Journeys**: [Khazwal User Journeys](../07-user-journeys/khazwal/material-prep-flow.md)
- **Sprint Documentation**: [Sprint Khazwal](../10-sprints/sprint-khazwal-material-prep.md)

---

*Last Updated: 29 Desember 2025*
