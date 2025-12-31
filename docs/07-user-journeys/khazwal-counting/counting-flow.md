# User Journey: Khazwal Counting Flow

**Feature:** Epic 2 - Penghitungan  
**Actor:** Staff Khazwal  
**Version:** 1.0.0

---

## Journey Overview

Journey ini menjelaskan bagaimana Staff Khazwal melakukan penghitungan hasil cetak dari awal hingga finalisasi, yaitu: melihat queue PO, memulai penghitungan, input hasil, breakdown kerusakan (jika > 5%), input alasan variance (jika ada selisih), dan finalisasi untuk advance PO ke stage pemotongan.

---

## Journey 1: Standard Counting (Happy Path)

**Precondition:** 
- Staff Khazwal sudah login
- Ada PO dengan status WAITING_COUNTING
- Print job sudah finalized

### Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│  COUNTING FLOW - HAPPY PATH                                 │
└─────────────────────────────────────────────────────────────┘

📱 START: Staff Login
    │
    ├─▶ Navigate to: Khazwal > Penghitungan
    │   └─ GET /api/khazwal/counting/queue
    │   └─ See: List PO menunggu (sorted FIFO)
    │   └─ Display: PO number, target, waiting time
    │
    ├─▶ Select: Click PO Card (PO 1234567890)
    │   └─ Navigate: /khazwal/counting/:poId
    │   └─ See: PO info, print info, "Mulai" button
    │
    ├─▶ Action: Click "Mulai Penghitungan"
    │   └─ POST /api/khazwal/counting/:po_id/start
    │   └─ Create counting record (status: IN_PROGRESS)
    │   └─ Update PO status: SEDANG_DIHITUNG
    │   └─ Show: Input form appears
    │
    ├─▶ Input: Jumlah Baik = 495
    │   └─ Trigger: Input debounce (1s)
    │   └─ PATCH /api/khazwal/counting/:id/result
    │   └─ Calculate: total, variance, percentages
    │   └─ Display: Real-time calculation
    │
    ├─▶ Input: Jumlah Rusak = 5
    │   └─ Trigger: Input debounce (1s)
    │   └─ PATCH /api/khazwal/counting/:id/result
    │   └─ Calculate: Defect = 1% (< 5%, no breakdown)
    │   └─ Display: Updated calculation
    │   └─ Note: No variance (500 total = 500 target)
    │
    ├─▶ Action: Click "Selesai Penghitungan"
    │   └─ Show: Finalize modal with summary
    │   └─ Display: Good 495, Defect 5, Variance 0
    │
    ├─▶ Confirm: Click "Selesaikan Penghitungan"
    │   └─ POST /api/khazwal/counting/:id/finalize
    │   └─ Update counting: status = COMPLETED
    │   └─ Update PO: stage = KHAZWAL_CUTTING, status = SIAP_POTONG
    │   └─ Log activity
    │   └─ Show: Success toast
    │   └─ Redirect: Back to queue page
    │
    ✅ END: PO ready for cutting, removed from queue
```

**Duration:** ~3-5 menit  
**Touchpoints:** 2 pages (Queue, Work)  
**API Calls:** 4 (queue, start, update x2, finalize)

---

## Journey 2: High Defect Counting (>5%)

**Scenario:** Staff menemukan kerusakan tinggi yang memerlukan breakdown jenis kerusakan.

### Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│  COUNTING FLOW - HIGH DEFECT (>5%)                          │
└─────────────────────────────────────────────────────────────┘

📱 START: (Steps 1-3 same as Journey 1)
    │
    ├─▶ Input: Jumlah Baik = 450
    │   └─ Auto-save after 1s
    │
    ├─▶ Input: Jumlah Rusak = 50
    │   └─ Calculate: Defect = 10% (> 5%)
    │   └─ Trigger: Defect breakdown form appears
    │   └─ Display: Red warning "Breakdown wajib diisi"
    │
    ├─▶ Fill Breakdown:
    │   ├─ Warna pudar: 30
    │   ├─ Tinta blobor: 15
    │   └─ Kertas sobek: 5
    │   └─ Validate: Sum (30+15+5=50) ✓ Match!
    │   └─ Display: Green "Breakdown valid"
    │   └─ PATCH /api/khazwal/counting/:id/result
    │
    ├─▶ Action: Click "Selesai"
    │   └─ Show: Modal with breakdown detail
    │   └─ Display: 
    │       - Baik: 450 (90%)
    │       - Rusak: 50 (10%)
    │       - Breakdown:
    │         • Warna pudar: 30
    │         • Tinta blobor: 15
    │         • Kertas sobek: 5
    │
    ├─▶ Confirm: Finalize
    │   └─ POST /api/khazwal/counting/:id/finalize
    │   └─ Success with breakdown recorded
    │
    ✅ END: Complete dengan breakdown detail
```

**Duration:** ~5-7 menit  
**Extra Steps:** Breakdown input & validation

---

## Journey 3: Counting dengan Variance

**Scenario:** Total hasil penghitungan tidak sama dengan target (ada selisih).

### Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│  COUNTING FLOW - WITH VARIANCE                              │
└─────────────────────────────────────────────────────────────┘

📱 START: (Steps 1-3 same as Journey 1)
    │
    ├─▶ Input: Jumlah Baik = 480
    │   └─ Auto-save
    │
    ├─▶ Input: Jumlah Rusak = 15
    │   └─ Calculate:
    │       - Total: 495
    │       - Target: 500
    │       - Variance: -5 (-1%)
    │   └─ Display: Orange warning "Selisih dari target"
    │   └─ Trigger: "Keterangan Selisih" field appears
    │   └─ Show: "Wajib diisi karena ada selisih"
    │
    ├─▶ Input Alasan: "5 lembar jatuh saat transport"
    │   └─ PATCH /api/khazwal/counting/:id/result
    │   └─ Enable: Finalize button
    │
    ├─▶ Action: Click "Selesai"
    │   └─ Modal shows:
    │       - Total: 495
    │       - Variance: -5
    │       - Alasan: "5 lembar jatuh saat transport"
    │
    ├─▶ Confirm: Finalize
    │   └─ Success with variance reason logged
    │
    ✅ END: Complete dengan alasan variance
```

**Duration:** ~4-6 menit  
**Extra Steps:** Variance reason input

---

## Journey 4: Edit Before Finalize

**Scenario:** Staff menyadari ada kesalahan input dan perlu koreksi sebelum finalize.

### Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│  COUNTING FLOW - EDIT BEFORE FINALIZE                       │
└─────────────────────────────────────────────────────────────┘

📱 START: Counting IN_PROGRESS
    │
    ├─▶ Current State:
    │   └─ Good: 480, Defect: 15
    │   └─ Already saved via PATCH
    │
    ├─▶ Realize Mistake: "Seharusnya 485 baik, 10 rusak"
    │
    ├─▶ Edit: Change Jumlah Baik → 485
    │   └─ Wait 1s debounce
    │   └─ PATCH /api/khazwal/counting/:id/result (2nd call)
    │
    ├─▶ Edit: Change Jumlah Rusak → 10
    │   └─ Wait 1s debounce
    │   └─ PATCH /api/khazwal/counting/:id/result (3rd call)
    │   └─ Calculate: Defect now 2% (< 5%)
    │   └─ Remove: Breakdown form (no longer needed)
    │
    ├─▶ Verify: Check calculation correct
    │   └─ Good: 485 (98%)
    │   └─ Defect: 10 (2%)
    │   └─ Variance: -5
    │
    ├─▶ Finalize: Proceed to complete
    │
    ✅ END: Complete dengan data corrected
```

**Key Point:** Multiple PATCH allowed before finalize, data updated setiap kali.

---

## Journey 5: Overdue PO (Priority)

**Scenario:** PO sudah menunggu > 2 jam (overdue), perlu diprioritaskan.

### Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│  COUNTING FLOW - OVERDUE PO                                 │
└─────────────────────────────────────────────────────────────┘

📱 START: Staff opens queue
    │
    ├─▶ See Queue:
    │   ├─ PO 111 (⚠️ Red badge: Menunggu 2.5 jam) ← OVERDUE
    │   ├─ PO 222 (Normal: Menunggu 45 menit)
    │   └─ PO 333 (Normal: Menunggu 30 menit)
    │
    ├─▶ Prioritize: Click overdue PO first
    │   └─ Visual: Red border pada card
    │   └─ Alert: "Segera lakukan penghitungan"
    │
    ├─▶ Process: (Same as standard flow)
    │
    ├─▶ Complete: Overdue cleared
    │
    ✅ END: Overdue resolved, next PO processed
```

**Visual Indicators:**
- Red border pada card
- Warning icon
- Overdue badge
- Alert message

---

## Journey 6: Empty Queue

**Scenario:** Tidak ada PO yang menunggu penghitungan.

### Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│  COUNTING FLOW - EMPTY QUEUE                                │
└─────────────────────────────────────────────────────────────┘

📱 START: Staff opens queue
    │
    ├─▶ GET /api/khazwal/counting/queue
    │   └─ Response: data = [] (empty)
    │
    ├─▶ Display Empty State:
    │   ├─ Icon: Document with checkmark
    │   ├─ Title: "Tidak Ada PO Menunggu"
    │   └─ Message: "Semua PO sudah dihitung atau belum ada yang selesai cetak"
    │
    ✅ END: Staff can wait or check back later
```

---

## Error Scenarios

### Error 1: Network Issue During Save

```
Staff input → Auto-save triggered → Network error
    │
    ├─▶ Display: Error toast "Gagal menyimpan, coba lagi"
    ├─▶ Keep: Data in form (not lost)
    ├─▶ Action: Staff can retry manually
    └─▶ Recovery: Auto-retry on next input
```

### Error 2: Concurrent Counting Attempt

```
Staff A starts → Staff B tries same PO
    │
    ├─▶ POST /api/khazwal/counting/:po_id/start
    ├─▶ Response: 409 Conflict
    ├─▶ Display: "Counting sudah dilakukan oleh user lain"
    └─▶ Redirect: Back to queue
```

### Error 3: Session Expired

```
Staff counting → Session expires → Try to save
    │
    ├─▶ PATCH /api/khazwal/counting/:id/result
    ├─▶ Response: 401 Unauthorized
    ├─▶ Trigger: Auto-refresh token
    ├─▶ Retry: Save request
    └─▶ Success: Seamless recovery
```

---

## Key Touchpoints

| Touchpoint | Purpose | Duration |
|------------|---------|----------|
| Queue Page | View & select PO | 30s - 1min |
| Work Page (Before Start) | Review PO info | 10-20s |
| Work Page (Input) | Enter counting results | 2-4min |
| Finalize Modal | Verify & confirm | 20-30s |
| Success Toast | Feedback | 2-3s |

---

## Mobile Considerations

### Mobile-Specific Behaviors

1. **Number Keyboard**
   - Input dengan `inputmode="numeric"`
   - Keyboard muncul otomatis saat focus

2. **Touch Targets**
   - Minimum 44x44px untuk semua buttons
   - Card dengan padding adequate

3. **Haptic Feedback**
   - Light vibration (10ms) saat click
   - Medium vibration (20ms) saat finalize

4. **Scroll Behavior**
   - Queue cards stackable
   - Form scrollable tanpa stuck header

---

## Related Documentation

- **API Reference:** [Khazwal Counting API](../../04-api-reference/khazwal-counting.md)
- **Testing Guide:** [Counting Testing](../../06-testing/khazwal-counting-testing.md)
- **Sprint Documentation:** [Sprint Counting Implementation](../../10-sprints/sprint-counting-implementation.md)

---

*Last Updated: 2025-12-30*
