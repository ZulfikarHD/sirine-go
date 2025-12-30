# Sprint 6 Testing Guide: Cetak Queue & Polish

**Feature**: Sprint 6 - Consumer Side & Polish  
**Version**: 1.0.0  
**Last Updated**: 30 Desember 2025

---

## Overview

Testing guide untuk Sprint 6 implementation yang mencakup Cetak Queue dengan OBC Master integration, enhanced Supervisor Monitoring, mobile enhancements (pull-to-refresh), dan UI polish dengan empty states.

---

## Test Environment Setup

### Prerequisites

```bash
# Backend running
make run-backend

# Frontend running
cd frontend && yarn dev

# Database seeded dengan test data
make seed-data
```

### Test Data Requirements

Untuk comprehensive testing, pastikan database memiliki:

- **Minimum 20 Production Orders** dengan status `READY_FOR_CETAK`
- **Mix of priorities:** URGENT (5), NORMAL (10), LOW (5)
- **Mix of due dates:** Past due (3), Due soon (5), Normal (12)
- **OBC Master data** terkait dengan setiap PO
- **Material photos** untuk beberapa PO (optional)
- **Staff users** untuk testing monitoring

---

## 1. Backend API Testing

### 1.1 Cetak Queue API

#### Test Case: CQ-001 - Get Queue Default Pagination

**Endpoint:** `GET /api/cetak/queue`

**Request:**
```bash
curl -X GET "http://localhost:8080/api/cetak/queue" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Queue cetak berhasil diambil",
  "data": {
    "items": [...],  // Array of 20 items (default)
    "total": 45,
    "page": 1,
    "per_page": 20,
    "total_pages": 3
  }
}
```

**Assertions:**
- [ ] Status code: 200
- [ ] Response structure matches schema
- [ ] Items array length ≤ per_page
- [ ] Each item has obc_master object
- [ ] OBC Master fields present: material, seri, warna, personalization
- [ ] Sorted by priority_score DESC, due_date ASC

---

#### Test Case: CQ-002 - Search by OBC Number

**Request:**
```bash
curl -X GET "http://localhost:8080/api/cetak/queue?search=OBC-2025" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected Behavior:**
- [ ] Returns only items matching "OBC-2025" in OBC number
- [ ] Case insensitive search
- [ ] Partial match supported

---

#### Test Case: CQ-003 - Filter by Priority

**Request:**
```bash
curl -X GET "http://localhost:8080/api/cetak/queue?priority=URGENT" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected Behavior:**
- [ ] Returns only URGENT priority items
- [ ] Priority field in response matches filter
- [ ] Sorted by due_date ASC within same priority

---

#### Test Case: CQ-004 - Custom Pagination

**Request:**
```bash
curl -X GET "http://localhost:8080/api/cetak/queue?page=2&per_page=50" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected Behavior:**
- [ ] Returns page 2 dengan 50 items max
- [ ] Pagination metadata correct
- [ ] No duplicate items across pages

---

#### Test Case: CQ-005 - Get Detail by ID

**Endpoint:** `GET /api/cetak/queue/:id`

**Request:**
```bash
curl -X GET "http://localhost:8080/api/cetak/queue/1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "po_id": 1,
    "obc_master": { ... },
    "material_prep": { ... }
  }
}
```

**Assertions:**
- [ ] Status code: 200
- [ ] obc_master object complete
- [ ] material_prep object complete
- [ ] Material photos array present (can be empty)
- [ ] All timestamps in ISO 8601 format

---

#### Test Case: CQ-006 - Detail Invalid ID

**Request:**
```bash
curl -X GET "http://localhost:8080/api/cetak/queue/99999" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected Response:**
```json
{
  "success": false,
  "message": "Production Order tidak ditemukan"
}
```

**Assertions:**
- [ ] Status code: 404
- [ ] Error message descriptive

---

#### Test Case: CQ-007 - Detail Not Ready for Cetak

**Scenario:** PO exists tapi status bukan READY_FOR_CETAK

**Expected Response:**
```json
{
  "success": false,
  "message": "PO tidak dalam status siap cetak"
}
```

**Assertions:**
- [ ] Status code: 400
- [ ] Clear error message

---

### 1.2 Authorization Testing

#### Test Case: AUTH-001 - Unauthorized Access

**Request:**
```bash
curl -X GET "http://localhost:8080/api/cetak/queue"
# No Authorization header
```

**Expected:**
- [ ] Status code: 401
- [ ] Error message about missing token

---

#### Test Case: AUTH-002 - Invalid Token

**Request:**
```bash
curl -X GET "http://localhost:8080/api/cetak/queue" \
  -H "Authorization: Bearer INVALID_TOKEN"
```

**Expected:**
- [ ] Status code: 401
- [ ] Error message about invalid token

---

#### Test Case: AUTH-003 - Insufficient Role

**Scenario:** User dengan role selain Operator/Supervisor Cetak

**Expected:**
- [ ] Status code: 403
- [ ] Error message about insufficient permissions

---

## 2. Frontend Component Testing

### 2.1 CetakQueueCard Component

#### Test Case: FC-001 - OBC Master Display

**Scenario:** Card displays OBC Master badges correctly

**Steps:**
1. Navigate to `/cetak/queue`
2. Observe first queue card

**Expected:**
- [ ] Material badge visible dengan correct text
- [ ] Seri badge visible
- [ ] Warna badge visible
- [ ] Personalization badge shows jika "Perso"
- [ ] Plat Number badge shows
- [ ] Badges wrapped properly on mobile

---

#### Test Case: FC-002 - Priority Badge

**Expected:**
- [ ] URGENT: Red background
- [ ] NORMAL: Gray background
- [ ] LOW: Blue background
- [ ] Badge text readable

---

#### Test Case: FC-003 - Card Click

**Steps:**
1. Click any queue card

**Expected:**
- [ ] Haptic feedback triggers (vibration on mobile)
- [ ] Detail modal opens
- [ ] Loading indicator shows while fetching
- [ ] Detail data loads correctly

---

### 2.2 CetakQueuePage

#### Test Case: FP-001 - Page Load

**Steps:**
1. Navigate to `/cetak/queue`

**Expected:**
- [ ] Page loads within 2 seconds
- [ ] Loading skeleton shows during fetch
- [ ] Queue items rendered with stagger animation
- [ ] Breadcrumbs show correct path
- [ ] Search input visible
- [ ] Priority filter visible

---

#### Test Case: FP-002 - Search Functionality

**Steps:**
1. Type "OBC-2025" in search input
2. Wait 300ms (debounce)

**Expected:**
- [ ] Loading indicator shows
- [ ] Filtered results display
- [ ] Results match search query
- [ ] Empty state shows if no results

---

#### Test Case: FP-003 - Priority Filter

**Steps:**
1. Click priority filter dropdown
2. Select "Urgent"

**Expected:**
- [ ] Dropdown opens dengan animation
- [ ] Options visible
- [ ] Selection updates filter
- [ ] Queue refreshes
- [ ] Only URGENT items shown

---

#### Test Case: FP-004 - Pagination

**Steps:**
1. Scroll to bottom
2. Click next page

**Expected:**
- [ ] Page number updates
- [ ] New items load
- [ ] Scroll to top animated
- [ ] No duplicate items

---

### 2.3 Pull-to-Refresh (Mobile)

#### Test Case: PTR-001 - Pull-to-Refresh Trigger

**Device:** Mobile or Chrome DevTools mobile mode

**Steps:**
1. Scroll to top (scrollY = 0)
2. Pull down slowly > 80px
3. Release

**Expected:**
- [ ] Pull indicator appears
- [ ] Text changes: "Tarik untuk refresh..." → "Lepas untuk refresh..."
- [ ] Icon spins when > 80px
- [ ] Haptic feedback on release (20ms vibration)
- [ ] Queue refreshes
- [ ] Indicator disappears after refresh

---

#### Test Case: PTR-002 - Pull-to-Refresh Not at Top

**Steps:**
1. Scroll down
2. Try pull-to-refresh gesture

**Expected:**
- [ ] Pull-to-refresh NOT triggered
- [ ] Normal scroll behavior
- [ ] No indicator shows

---

#### Test Case: PTR-003 - Pull-to-Refresh During Loading

**Steps:**
1. Trigger refresh
2. Try pull-to-refresh again while loading

**Expected:**
- [ ] Second pull ignored
- [ ] No conflict with ongoing refresh

---

### 2.4 Empty States

#### Test Case: ES-001 - Empty Queue (No Items)

**Setup:** Database dengan 0 READY_FOR_CETAK items

**Expected:**
- [ ] Empty state card shows
- [ ] Icon dengan gradient background
- [ ] Icon pop-in animation
- [ ] Title: "Antrian Kosong"
- [ ] Descriptive message about waiting for material prep
- [ ] Glass card styling

---

#### Test Case: ES-002 - Empty Search Results

**Steps:**
1. Search for non-existent OBC
2. No results returned

**Expected:**
- [ ] Empty state shows
- [ ] Title: "Tidak Ditemukan"
- [ ] Message suggests changing search criteria

---

#### Test Case: ES-003 - Empty Staff Activity (Monitoring)

**Page:** `/khazwal/monitoring`

**Setup:** No active staff

**Expected:**
- [ ] Empty state with UserX icon
- [ ] Title: "Belum Ada Aktivitas"
- [ ] Descriptive message

---

### 2.5 Supervisor Monitoring Page

#### Test Case: SM-001 - Stats Cards

**Steps:**
1. Navigate to `/khazwal/monitoring`

**Expected:**
- [ ] 4 stats cards visible
- [ ] Total in queue count correct
- [ ] Total in progress count correct
- [ ] Total completed today count correct
- [ ] Average duration formatted (e.g., "90m" or "1j 30m")

---

#### Test Case: SM-002 - Staff Activity Cards

**Expected:**
- [ ] Active staff cards show
- [ ] Each card shows current PO
- [ ] OBC Number visible dengan bold styling
- [ ] OBC Master badges show (Material, Seri, Warna)
- [ ] Duration running shown
- [ ] Status badge (Aktif/Idle)

---

#### Test Case: SM-003 - Auto-Refresh

**Steps:**
1. Wait 30 seconds
2. Observe stats update

**Expected:**
- [ ] Stats refresh automatically
- [ ] Countdown indicator shows "Auto-refresh: Xs"
- [ ] Loading state tidak intrusive
- [ ] No console errors

---

#### Test Case: SM-004 - Manual Refresh

**Steps:**
1. Click refresh button

**Expected:**
- [ ] Refresh icon spins
- [ ] Data refreshes immediately
- [ ] Countdown resets to 30s
- [ ] Haptic feedback (10ms)

---

## 3. Store Integration Testing

### 3.1 Cetak Store

#### Test Case: ST-001 - Queue Fetch

**Action:** `cetakStore.getCetakQueue()`

**Expected:**
- [ ] `queueLoading` true during fetch
- [ ] `queue` array populated on success
- [ ] `queuePagination` updated
- [ ] `queueLoading` false after complete
- [ ] `queueError` null on success

---

#### Test Case: ST-002 - Detail Fetch

**Action:** `cetakStore.getCetakDetail(poId)`

**Expected:**
- [ ] `detailLoading` true during fetch
- [ ] `currentDetail` populated on success
- [ ] `detailLoading` false after complete
- [ ] `detailError` null on success

---

#### Test Case: ST-003 - Store Getters

**Test:** `urgentQueue` getter

**Expected:**
- [ ] Returns only URGENT priority items
- [ ] Reactive to state changes

**Test:** `pastDueQueue` getter

**Expected:**
- [ ] Returns only past due items
- [ ] is_past_due = true

---

#### Test Case: ST-004 - Clear Methods

**Actions:**
1. `cetakStore.clearQueue()`
2. `cetakStore.clearDetail()`

**Expected:**
- [ ] Queue array empty
- [ ] Pagination reset
- [ ] Detail null
- [ ] Errors cleared

---

### 3.2 Khazwal Store (Monitoring)

#### Test Case: ST-005 - Monitoring Stats Fetch

**Action:** `khazwalStore.getMonitoringStats()`

**Expected:**
- [ ] `monitoringLoading` true during fetch
- [ ] `monitoringStats` populated
- [ ] Includes: stats, staff_active, recent_completions
- [ ] `monitoringLoading` false after complete

---

## 4. Cross-Browser Testing

### Browsers to Test

- [ ] Chrome (latest)
- [ ] Firefox (latest)
- [ ] Safari (latest)
- [ ] Edge (latest)
- [ ] Mobile Safari (iOS)
- [ ] Mobile Chrome (Android)

### Features to Verify

- [ ] Pull-to-refresh works on mobile browsers
- [ ] Haptic feedback works (where supported)
- [ ] Animations smooth (60fps)
- [ ] Glass morphism renders correctly
- [ ] Touch gestures work
- [ ] Responsive layouts adapt

---

## 5. Performance Testing

### 5.1 Load Time

**Metrics:**

| Metric | Target | Actual |
|--------|--------|--------|
| Initial page load | < 2s | ___ |
| Queue fetch | < 1s | ___ |
| Detail modal open | < 500ms | ___ |
| Pull-to-refresh | < 1.5s | ___ |
| Empty state render | < 100ms | ___ |

---

### 5.2 Animation Performance

**Test:** Monitor FPS during animations

**Tools:** Chrome DevTools > Performance

**Target:** 60 FPS maintained

**Animations to test:**
- [ ] Page entrance animations
- [ ] Card stagger animations
- [ ] Empty state icon pop-in
- [ ] Pull-to-refresh indicator

---

### 5.3 Memory Leaks

**Steps:**
1. Open page
2. Navigate between pages 10 times
3. Monitor memory usage

**Expected:**
- [ ] Memory usage stable
- [ ] No continuous increase
- [ ] Stores cleaned up on unmount

---

## 6. Accessibility Testing

### 6.1 Keyboard Navigation

- [ ] Tab through all interactive elements
- [ ] Focus visible on all elements
- [ ] Enter/Space activates buttons
- [ ] Escape closes modals

---

### 6.2 Screen Reader

- [ ] Page title announced
- [ ] Cards have descriptive labels
- [ ] Loading states announced
- [ ] Error messages announced

---

## 7. Error Scenarios

### 7.1 Network Errors

**Test:** Simulate network failure

**Expected:**
- [ ] Error message displayed
- [ ] Retry option available
- [ ] No blank/broken UI

---

### 7.2 API Errors

**Scenarios:**
- [ ] 401 Unauthorized → Redirect to login
- [ ] 403 Forbidden → Permission error message
- [ ] 404 Not Found → Item not found message
- [ ] 500 Server Error → Generic error message

---

## 8. Regression Testing

### Existing Features to Verify

- [ ] Material Prep Queue still works
- [ ] Material Prep Process still works
- [ ] Khazwal History still works
- [ ] Authentication still works
- [ ] Navigation still works
- [ ] No breaking changes

---

## 9. Test Execution Summary

### Manual Test Execution

| Test Suite | Total Tests | Passed | Failed | Skipped |
|------------|-------------|--------|--------|---------|
| Backend API | 10 | ___ | ___ | ___ |
| Frontend Components | 15 | ___ | ___ | ___ |
| Store Integration | 6 | ___ | ___ | ___ |
| Mobile UX | 5 | ___ | ___ | ___ |
| Performance | 3 | ___ | ___ | ___ |
| Accessibility | 2 | ___ | ___ | ___ |
| Error Handling | 2 | ___ | ___ | ___ |
| Regression | 6 | ___ | ___ | ___ |
| **TOTAL** | **49** | ___ | ___ | ___ |

---

## 10. Bug Report Template

**Bug ID:** [BUG-XXX]  
**Severity:** [Critical / High / Medium / Low]  
**Component:** [Backend / Frontend / Store / Mobile]

**Description:**
[Clear description of the bug]

**Steps to Reproduce:**
1. [Step 1]
2. [Step 2]
3. [Step 3]

**Expected Behavior:**
[What should happen]

**Actual Behavior:**
[What actually happens]

**Environment:**
- Browser: [Chrome 120]
- OS: [Windows 11]
- Screen: [Desktop / Mobile]

**Screenshots:**
[Attach screenshots if applicable]

---

## Related Documentation

- **API Reference:** [Cetak API](../04-api-reference/cetak.md)
- **Sprint Documentation:** [Sprint 6](../10-sprints/sprint-khazwal-sprint6.md)
- **User Journeys:** [Cetak Queue Flow](../07-user-journeys/khazwal/cetak-queue-flow.md)

---

*Last Updated: 30 Desember 2025*  
*Version: 1.0.0*  
*Status: Ready for Testing*
