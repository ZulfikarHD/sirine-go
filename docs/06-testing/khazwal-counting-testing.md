# Testing Guide: Khazwal Counting

**Feature:** Epic 2 - Penghitungan (Counting)  
**Version:** 1.0.0  
**Last Updated:** 2025-12-30

---

## Overview

Testing guide untuk feature Khazwal Counting yang bertujuan untuk memastikan semua functionality bekerja dengan benar, yaitu: queue management, input penghitungan, validation rules, dan finalisasi.

---

## Test Coverage Summary

| Category | Test Cases | Priority |
|----------|------------|----------|
| Backend API | 15 | High |
| Frontend Components | 12 | High |
| Integration Flow | 8 | Critical |
| Edge Cases | 10 | Medium |
| **Total** | **45** | - |

---

## Backend API Tests

### 1. Queue Endpoint Tests

| Test ID | Scenario | Expected Result |
|---------|----------|-----------------|
| BE-Q-001 | GET queue tanpa filter | Return semua PO dengan status WAITING_COUNTING |
| BE-Q-002 | GET queue dengan machine_id filter | Return hanya PO dari mesin tertentu |
| BE-Q-003 | GET queue dengan date range | Return PO dalam range tanggal |
| BE-Q-004 | Queue sorting FIFO | PO sorted by print_completed_at ASC |
| BE-Q-005 | Overdue flag calculation | is_overdue = true jika waiting > 120 menit |

### 2. Start Counting Tests

| Test ID | Scenario | Expected Result |
|---------|----------|-----------------|
| BE-S-001 | Start counting untuk PO valid | Status 201, counting record created |
| BE-S-002 | Start dengan PO tidak ready | Status 400, error message |
| BE-S-003 | Start dengan counting sudah ada | Status 409, conflict error |
| BE-S-004 | Start updates PO status | PO status = SEDANG_DIHITUNG |
| BE-S-005 | Start logs activity | Activity log created |

### 3. Update Result Tests

| Test ID | Scenario | Expected Result |
|---------|----------|-----------------|
| BE-U-001 | Update dengan data valid | Status 200, calculations correct |
| BE-U-002 | Update dengan quantity negatif | Status 422, validation error |
| BE-U-003 | Defect > 5% tanpa breakdown | Status 422, breakdown required |
| BE-U-004 | Breakdown sum tidak match | Status 422, sum mismatch error |
| BE-U-005 | Variance tanpa reason | Status 422, reason required |
| BE-U-006 | Update setelah completed | Status 400, already finalized |
| BE-U-007 | Multiple updates before finalize | All updates accepted |

### 4. Finalize Tests

| Test ID | Scenario | Expected Result |
|---------|----------|-----------------|
| BE-F-001 | Finalize dengan data lengkap | Status 200, status = COMPLETED |
| BE-F-002 | Finalize tanpa required fields | Status 400, validation error |
| BE-F-003 | Finalize updates PO | PO status = SIAP_POTONG, stage = KHAZWAL_CUTTING |
| BE-F-004 | Finalize locks data | Subsequent updates rejected |
| BE-F-005 | Duration calculation | duration_minutes calculated correctly |

### 5. Calculation Tests

| Test ID | Scenario | Input | Expected Output |
|---------|----------|-------|-----------------|
| BE-C-001 | Total counted | Good: 480, Defect: 15 | total_counted: 495 |
| BE-C-002 | Variance | Total: 495, Target: 500 | variance: -5 |
| BE-C-003 | Percentage good | Good: 480, Total: 495 | percentage_good: 96.97 |
| BE-C-004 | Percentage defect | Defect: 15, Total: 495 | percentage_defect: 3.03 |

---

## Frontend Component Tests

### 1. CountingQueueCard Component

| Test ID | Scenario | Expected Behavior |
|---------|----------|-------------------|
| FE-Q-001 | Display PO info | Shows PO number, OBC, target |
| FE-Q-002 | Overdue warning | Red indicator jika is_overdue = true |
| FE-Q-003 | Waiting time format | Display formatted waiting time |
| FE-Q-004 | Click navigation | Navigate to CountingWorkPage |

### 2. CountingInputForm Component

| Test ID | Scenario | Expected Behavior |
|---------|----------|-------------------|
| FE-I-001 | Input validation | Tidak accept nilai negatif |
| FE-I-002 | Number keyboard | inputmode="numeric" pada mobile |
| FE-I-003 | Debounced save | Save triggered 1s setelah input stop |
| FE-I-004 | Loading state | Show loading saat saving |

### 3. CountingCalculation Component

| Test ID | Scenario | Expected Behavior |
|---------|----------|-------------------|
| FE-C-001 | Real-time calculation | Update otomatis saat input berubah |
| FE-C-002 | Variance indicator | Color based on positive/negative/zero |
| FE-C-003 | Tolerance warning | Show warning jika variance > 2% |
| FE-C-004 | Defect threshold | Show alert jika defect > 5% |

### 4. DefectBreakdownForm Component

| Test ID | Scenario | Expected Behavior |
|---------|----------|-------------------|
| FE-D-001 | Sum validation | Show error jika sum != quantity_defect |
| FE-D-002 | Progress bar | Visual indicator untuk sum progress |
| FE-D-003 | Valid state | Green indicator jika sum = quantity_defect |

---

## Integration Flow Tests

### 1. Happy Path Flow

| Test ID | Step | Action | Expected Result |
|---------|------|--------|-----------------|
| INT-001 | 1 | Open queue page | List PO ditampilkan |
| INT-002 | 2 | Click PO card | Navigate ke work page |
| INT-003 | 3 | Click "Mulai Penghitungan" | Counting started, form shown |
| INT-004 | 4 | Input quantity good & defect | Real-time calc updated |
| INT-005 | 5 | Click "Selesai" | Finalize modal shown |
| INT-006 | 6 | Confirm finalize | Success, redirect to queue |

### 2. High Defect Flow (>5%)

| Test ID | Step | Action | Expected Result |
|---------|------|--------|-----------------|
| INT-007 | 1-3 | (Same as happy path) | Counting started |
| INT-008 | 4 | Input defect 10% | Breakdown form appears |
| INT-009 | 5 | Fill breakdown | Sum validation active |
| INT-010 | 6 | Finalize | Success only if sum valid |

---

## Edge Cases Testing

| Test ID | Scenario | Expected Behavior |
|---------|----------|-------------------|
| EDGE-001 | Network error saat save | Retry mechanism, error message |
| EDGE-002 | Browser refresh mid-counting | Data restored from server |
| EDGE-003 | Concurrent counting (2 staff) | Second staff gets conflict error |
| EDGE-004 | Zero total counted | Allowed (can finalize with 0+0) |
| EDGE-005 | Variance exactly 0 | No variance reason required |
| EDGE-006 | Defect exactly 5% | Breakdown NOT required (>5% rule) |
| EDGE-007 | Very large numbers | UI handles large quantities |
| EDGE-008 | Decimal input | Rounded to integer |
| EDGE-009 | Empty queue | Show empty state message |
| EDGE-010 | Finalize modal closed | Data NOT finalized |

---

## Manual QA Checklist

### Desktop Testing

- [ ] Queue page loads with correct data
- [ ] Card click navigation works
- [ ] Start counting creates record
- [ ] Input form validation works
- [ ] Real-time calculation accurate
- [ ] Defect breakdown shows when needed
- [ ] Variance reason shows when needed
- [ ] Auto-save working (debounced)
- [ ] Finalize modal shows summary
- [ ] Finalize success redirects
- [ ] Error messages user-friendly
- [ ] Loading states displayed
- [ ] Responsive layout (1920x1080, 1366x768)

### Mobile Testing

- [ ] Queue page mobile responsive
- [ ] Cards stack vertically
- [ ] Number keyboard appears on input
- [ ] Touch targets minimum 44x44px
- [ ] Input form scrollable
- [ ] Modal fits screen
- [ ] Navigation smooth
- [ ] No horizontal scroll
- [ ] Text readable (min 14px)
- [ ] Buttons accessible
- [ ] Test on iOS Safari
- [ ] Test on Android Chrome

### User Experience

- [ ] Haptic feedback on actions (mobile)
- [ ] Press feedback (scale on tap)
- [ ] Smooth animations
- [ ] No layout shift
- [ ] Loading skeletons shown
- [ ] Empty states clear
- [ ] Error states recoverable
- [ ] Success feedback visible
- [ ] Navigation intuitive
- [ ] Labels descriptive

### Data Integrity

- [ ] Queue shows correct PO
- [ ] Calculations accurate
- [ ] Data persisted on save
- [ ] Data locked after finalize
- [ ] Cannot modify after complete
- [ ] PO status updated correctly
- [ ] Activity logs created
- [ ] No data loss on refresh
- [ ] Concurrent updates handled
- [ ] Transaction rollback on error

---

## Performance Tests

| Test ID | Metric | Target | How to Test |
|---------|--------|--------|-------------|
| PERF-001 | Queue load time | < 2s | Chrome DevTools Network |
| PERF-002 | Start counting response | < 1s | API response time |
| PERF-003 | Auto-save debounce | 1s | Input then wait 1s |
| PERF-004 | Finalize transaction | < 3s | API response time |
| PERF-005 | Animation frame rate | 60fps | Chrome Performance tab |

---

## Security Tests

| Test ID | Scenario | Expected Behavior |
|---------|----------|-------------------|
| SEC-001 | Unauthenticated access | 401 Unauthorized |
| SEC-002 | Wrong role access | 403 Forbidden |
| SEC-003 | SQL injection attempt | Sanitized, no execution |
| SEC-004 | XSS in variance_reason | Escaped, no script execution |
| SEC-005 | Token expired | Auto-refresh or redirect login |

---

## Regression Test Suite

Run setiap kali ada changes pada counting feature:

1. **Smoke Tests** (5 min)
   - Login as STAFF_KHAZWAL
   - Open counting queue
   - Start counting
   - Input and save
   - Finalize successfully

2. **Critical Path** (15 min)
   - Happy path flow
   - High defect flow
   - Edit before finalize
   - Mobile responsive

3. **Full Regression** (30 min)
   - All backend API tests
   - All frontend component tests
   - All edge cases
   - Cross-browser testing

---

## Bug Reporting Template

```markdown
**Bug ID:** BUG-COUNTING-XXX
**Severity:** [Critical/High/Medium/Low]
**Environment:** [Dev/Staging/Production]

**Steps to Reproduce:**
1. 
2. 
3. 

**Expected Result:**

**Actual Result:**

**Screenshots:**

**Browser/Device:**

**Additional Notes:**
```

---

## Test Data Setup

### Required Test Data

```sql
-- PO dengan status WAITING_COUNTING
INSERT INTO production_orders ...

-- Print job summary finalized
INSERT INTO print_job_summaries ...

-- Test users dengan role STAFF_KHAZWAL
INSERT INTO users ...
```

---

## Related Documentation

- **API Reference:** [Khazwal Counting API](../04-api-reference/khazwal-counting.md)
- **Sprint Documentation:** [Sprint Counting Implementation](../10-sprints/sprint-counting-implementation.md)
- **User Journeys:** [Counting User Journeys](../07-user-journeys/khazwal-counting/)

---

*Last Updated: 2025-12-30*
