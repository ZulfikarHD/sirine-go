# Khazwal Cutting Testing Guide

## Overview

Panduan testing untuk fitur Khazwal Cutting yang mencakup manual testing, API testing, dan regression testing untuk memastikan proses pemotongan lembar besar berfungsi dengan benar dari queue management hingga finalisasi.

**Feature:** Khazwal Cutting (Epic 03)  
**Scope:** Sprint 1 (Queue & Start Process)  
**Status:** ✅ Implementation Complete, Testing Ready  

---

## Prerequisites

### Data Setup

Sebelum testing, pastikan memiliki:

```sql
-- 1. User dengan role STAFF_KHAZWAL
INSERT INTO users (nip, full_name, email, password_hash, role, department, status) VALUES
('12345', 'Test Staff Khazwal', 'staff@test.local', '$2a$12$...', 'STAFF_KHAZWAL', 'KHAZWAL', 'ACTIVE');

-- 2. PO dengan counting completed
-- PO harus di stage KHAZWAL_CUTTING dengan status SIAP_POTONG
UPDATE production_orders 
SET current_stage = 'KHAZWAL_CUTTING', current_status = 'SIAP_POTONG'
WHERE id = 123;

-- 3. Counting result untuk PO tersebut
INSERT INTO khazwal_counting_results 
(production_order_id, quantity_good, quantity_defect, status, completed_at, counted_by)
VALUES (123, 15000, 100, 'COMPLETED', NOW(), 1);
```

### Environment

- Backend running: `http://localhost:8080`
- Frontend running: `http://localhost:5173`
- Database accessible
- Test user authenticated

---

## Test Plan Summary

| Test Type | Total Cases | Status |
|-----------|-------------|--------|
| **Manual UI Tests** | 15 | 📋 Ready |
| **API Tests** | 12 | 📋 Ready |
| **Edge Cases** | 8 | 📋 Ready |
| **Performance Tests** | 3 | 📋 Ready |
| **Total** | **38** | 📋 Ready |

---

## 1. Manual UI Testing

### Test Suite 1.1: Cutting Queue Page

| ID | Test Case | Expected Result | Status |
|----|-----------|-----------------|--------|
| UI-01 | Navigate to Pemotongan from sidebar | Page loads, shows queue | ⏳ |
| UI-02 | View queue with no PO | Empty state displays with message | ⏳ |
| UI-03 | View queue with 5 PO | All 5 cards displayed with correct data | ⏳ |
| UI-04 | Check priority badge colors | URGENT=red, NORMAL=blue | ⏳ |
| UI-05 | Check estimated output calculation | Shows input × 2 correctly | ⏳ |
| UI-06 | Check waiting time display | Shows minutes correctly | ⏳ |
| UI-07 | Check overdue indicator | Red warning if > 60 min | ⏳ |
| UI-08 | Filter by URGENT priority | Only URGENT PO displayed | ⏳ |
| UI-09 | Filter by date range | PO within range displayed | ⏳ |
| UI-10 | Reset filter button | All PO displayed again | ⏳ |
| UI-11 | Refresh button | Data refreshes, loading shown | ⏳ |
| UI-12 | Stats cards | Total, Urgent, Normal counts correct | ⏳ |
| UI-13 | Mobile responsiveness | Cards stack vertically, buttons accessible | ⏳ |
| UI-14 | Click "Mulai Pemotongan" | Navigates to start page | ⏳ |
| UI-15 | Loading skeleton | Shows while fetching queue | ⏳ |

### Test Suite 1.2: Cutting Start Page

| ID | Test Case | Expected Result | Status |
|----|-----------|-----------------|--------|
| UI-16 | Navigate from queue card | Page loads with PO details | ⏳ |
| UI-17 | Check PO info display | Number, OBC, priority shown | ⏳ |
| UI-18 | Check input display | Shows lembar besar from counting | ⏳ |
| UI-19 | Check estimated output | Shows input × 2 | ⏳ |
| UI-20 | Operator auto-fill | Shows logged-in user name | ⏳ |
| UI-21 | Machine selector empty | Submit button disabled | ⏳ |
| UI-22 | Select machine | Submit button enabled | ⏳ |
| UI-23 | Click cancel | Returns to queue page | ⏳ |
| UI-24 | Submit valid form | Success, redirects to queue | ⏳ |
| UI-25 | Breadcrumb navigation | Shows "Pemotongan > PO {number}" | ⏳ |

---

## 2. API Testing

### Test Suite 2.1: GET /api/khazwal/cutting/queue

**Base URL:** `http://localhost:8080/api/khazwal/cutting/queue`

| ID | Test Case | Request | Expected Response | Status |
|----|-----------|---------|-------------------|--------|
| API-01 | Get queue without filters | `GET /queue` | 200, array of PO | ⏳ |
| API-02 | Filter by URGENT priority | `GET /queue?priority=URGENT` | 200, only URGENT PO | ⏳ |
| API-03 | Filter by date range | `GET /queue?date_from=2024-01-01&date_to=2024-01-31` | 200, PO in range | ⏳ |
| API-04 | Sort by date ascending | `GET /queue?sort_by=date&sort_order=asc` | 200, oldest first | ⏳ |
| API-05 | Invalid priority value | `GET /queue?priority=INVALID` | 200, empty array | ⏳ |
| API-06 | Unauthorized request | No auth token | 401 Unauthorized | ⏳ |
| API-07 | Wrong role | Token with OPERATOR_CETAK | 403 Forbidden | ⏳ |

**Example Request:**

```bash
curl -X GET "http://localhost:8080/api/khazwal/cutting/queue?priority=URGENT" \
  -H "Authorization: Bearer {token}"
```

**Expected Response:**

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
    "total": 1,
    "urgent_count": 1,
    "normal_count": 0
  }
}
```

### Test Suite 2.2: POST /api/khazwal/cutting/po/:po_id/start

**Base URL:** `http://localhost:8080/api/khazwal/cutting/po/{po_id}/start`

| ID | Test Case | Request Body | Expected Response | Status |
|----|-----------|--------------|-------------------|--------|
| API-08 | Valid start request | `{"cutting_machine": "Mesin A"}` | 200, cutting created | ⏳ |
| API-09 | Missing machine | `{}` | 400, validation error | ⏳ |
| API-10 | PO not found | Valid body, invalid PO ID | 404, PO not found | ⏳ |
| API-11 | Already started | Same PO twice | 409, already started | ⏳ |
| API-12 | Counting not completed | PO without counting | 400, counting not ready | ⏳ |

**Example Request:**

```bash
curl -X POST "http://localhost:8080/api/khazwal/cutting/po/123/start" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "cutting_machine": "Mesin A"
  }'
```

**Expected Response:**

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

---

## 3. Edge Cases & Error Scenarios

| ID | Scenario | Steps | Expected Behavior | Status |
|----|----------|-------|-------------------|--------|
| EDGE-01 | No PO in queue | Navigate to queue | Empty state with message | ⏳ |
| EDGE-02 | Concurrent start | 2 users start same PO | One succeeds, one gets 409 | ⏳ |
| EDGE-03 | PO deleted during start | Start PO that gets deleted | 404 error, graceful handling | ⏳ |
| EDGE-04 | Network timeout | Slow API response | Loading state, timeout after 30s | ⏳ |
| EDGE-05 | Large queue (50+ PO) | View queue with many PO | Pagination/scroll works | ⏳ |
| EDGE-06 | Invalid JWT token | Expired/malformed token | 401, redirect to login | ⏳ |
| EDGE-07 | Machine with special chars | Machine name with emoji | Stores correctly | ⏳ |
| EDGE-08 | Counting result = 0 | PO with 0 good quantity | estimated_output = 0 | ⏳ |

---

## 4. Performance Testing

| ID | Test Case | Target | Measurement | Status |
|----|-----------|--------|-------------|--------|
| PERF-01 | Queue load time | < 2 seconds | Time from click to data display | ⏳ |
| PERF-02 | Queue with 100 PO | < 3 seconds | API response time | ⏳ |
| PERF-03 | Start process time | < 1 second | API response + DB update | ⏳ |

**Performance Targets:**
- Queue API: < 500ms
- Start API: < 300ms
- Frontend render: < 100ms
- Total user wait: < 2s

---

## 5. Database Verification

### After Start Cutting

```sql
-- Verify cutting record created
SELECT * FROM khazwal_cutting_results 
WHERE production_order_id = 123;
-- Expected: 1 row, status = 'IN_PROGRESS'

-- Verify PO status updated
SELECT current_stage, current_status 
FROM production_orders 
WHERE id = 123;
-- Expected: current_stage = 'KHAZWAL_CUTTING', current_status = 'SEDANG_DIPOTONG'

-- Verify stage tracking
SELECT started_at 
FROM po_stage_tracking 
WHERE production_order_id = 123 AND stage = 'KHAZWAL_CUTTING';
-- Expected: timestamp = NOW()

-- Verify activity log
SELECT * FROM activity_logs 
WHERE entity_type = 'CUTTING' AND entity_id = (SELECT id FROM khazwal_cutting_results WHERE production_order_id = 123)
ORDER BY created_at DESC LIMIT 1;
-- Expected: 1 row, action = 'CREATE' or 'START'
```

---

## 6. Regression Testing

### Checklist untuk setiap deployment:

- [ ] Login as STAFF_KHAZWAL works
- [ ] Sidebar shows "Pemotongan" menu
- [ ] Queue page loads without errors
- [ ] Can filter queue by priority
- [ ] Can start cutting process
- [ ] PO status updates correctly
- [ ] No console errors
- [ ] Mobile responsive
- [ ] Counting feature still works
- [ ] Material prep feature still works

---

## 7. Bug Reporting Template

Jika menemukan bug, report dengan format:

```markdown
## Bug Report

**ID:** BUG-CUTTING-001
**Severity:** High / Medium / Low
**Status:** Open / In Progress / Closed

### Description
Brief description of the bug

### Steps to Reproduce
1. Step 1
2. Step 2
3. Step 3

### Expected Behavior
What should happen

### Actual Behavior
What actually happens

### Screenshots
[Attach if applicable]

### Environment
- Browser: Chrome 120
- OS: Windows 11
- Backend: localhost:8080
- Frontend: localhost:5173

### Additional Context
Any other relevant information
```

---

## 8. QA Sign-off Checklist

Before marking as "Production Ready":

### Functional
- [ ] All UI tests passed
- [ ] All API tests passed
- [ ] All edge cases handled
- [ ] Error messages are user-friendly
- [ ] Success feedback is clear

### Non-Functional
- [ ] Performance targets met
- [ ] Mobile responsive tested
- [ ] Cross-browser compatible (Chrome, Firefox, Safari)
- [ ] Accessibility checked (keyboard navigation)
- [ ] Security validated (auth, role checks)

### Documentation
- [ ] API docs complete
- [ ] User manual updated
- [ ] Sprint doc reviewed
- [ ] Known issues documented

### Deployment
- [ ] Database migration tested
- [ ] Rollback plan ready
- [ ] Monitoring alerts configured
- [ ] Production data backup done

---

## 9. Test Data Generator

### Script untuk generate test data:

```sql
-- Generate 10 test PO ready for cutting
INSERT INTO production_orders (po_number, obc_number, priority, target_quantity, current_stage, current_status, created_at)
VALUES
(2024001, 'OBC-2024-001', 'URGENT', 15000, 'KHAZWAL_CUTTING', 'SIAP_POTONG', NOW() - INTERVAL 2 HOUR),
(2024002, 'OBC-2024-002', 'NORMAL', 20000, 'KHAZWAL_CUTTING', 'SIAP_POTONG', NOW() - INTERVAL 1 HOUR),
(2024003, 'OBC-2024-003', 'URGENT', 10000, 'KHAZWAL_CUTTING', 'SIAP_POTONG', NOW() - INTERVAL 3 HOUR);

-- Generate counting results
INSERT INTO khazwal_counting_results (production_order_id, quantity_good, quantity_defect, status, completed_at, counted_by)
SELECT id, target_quantity * 0.98, target_quantity * 0.02, 'COMPLETED', NOW() - INTERVAL 30 MINUTE, 1
FROM production_orders
WHERE current_stage = 'KHAZWAL_CUTTING' AND current_status = 'SIAP_POTONG';
```

---

## 10. Related Documentation

- **API Reference:** [Khazwal Cutting API](../04-api-reference/khazwal-cutting.md)
- **Sprint Documentation:** [Sprint Cutting Phase 1](../10-sprints/sprint-cutting-phase1.md)
- **User Manual:** TBD

---

**Last Updated:** 2026-01-10  
**Tested By:** TBD  
**Status:** Testing Ready  
**Next Test Date:** TBD
