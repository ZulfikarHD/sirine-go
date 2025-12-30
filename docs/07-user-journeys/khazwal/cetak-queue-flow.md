# User Journey: Cetak Queue Flow

**Feature**: Cetak Queue (Consumer Side)  
**Sprint**: Sprint 6  
**User Role**: Operator Cetak, Supervisor Cetak  
**Last Updated**: 30 Desember 2025

---

## Overview

User journey untuk Operator/Supervisor Cetak dalam mengakses antrian PO yang siap untuk dicetak setelah material preparation selesai di Unit Khazanah Awal, yaitu: browse queue dengan filtering, view detail PO dengan OBC Master info dan material photos, serta mobile enhancements dengan pull-to-refresh.

---

## User Persona

**Name:** Budi Santoso  
**Role:** Operator Cetak  
**Experience:** 2 tahun di Unit Cetak  
**Device:** Tablet Samsung (untuk portability di ruang produksi)  
**Goals:**
- Melihat antrian PO yang siap dicetak dengan cepat
- Memahami specifications dari OBC Master
- Memprioritaskan PO urgent dan past due
- Akses material photos untuk referral

---

## Journey Map

### Scenario 1: Budi Checks Cetak Queue (Normal Flow)

```
📍 START: Operator Login Dashboard
    │
    ├─▶ Navigate to: Cetak > Antrian Cetak
    │   └─ Path: /cetak/queue
    │
    ├─▶ Page Loads (< 2s)
    │   ├─ Loading skeleton appears
    │   ├─ Queue items fetched from API
    │   └─ Items rendered dengan stagger animation
    │
    ├─▶ Browse Queue
    │   ├─ See 20 items (default pagination)
    │   ├─ Each card shows:
    │   │   ├─ PO Number
    │   │   ├─ OBC Number (bold indigo text)
    │   │   ├─ Priority badge (color-coded)
    │   │   ├─ OBC Master badges (Material, Seri, Warna, Personalization)
    │   │   ├─ Quantity
    │   │   ├─ Due date with urgency indicator
    │   │   ├─ Material ready timestamp
    │   │   └─ Prepared by name with avatar
    │   └─ Queue sorted by: Priority → Due Date
    │
    ├─▶ Click Card to View Detail
    │   ├─ Haptic feedback (10ms vibration)
    │   ├─ Detail modal opens dengan animation
    │   ├─ Loading indicator while fetching
    │   └─ Detail displayed:
    │       ├─ Full PO information
    │       ├─ OBC Master specifications
    │       ├─ Material prep info (duration, prepared by, notes)
    │       ├─ Kertas & Tinta details
    │       └─ Material photos gallery (if any)
    │
    ├─▶ Close Detail Modal
    │   └─ Back to queue list
    │
    └─▶ Proceed to Cetak Process
        └─ (Out of scope - Sprint 7+)
```

**Touchpoints:**
1. Login → Dashboard
2. Sidebar Navigation
3. Queue List Page
4. Card Interaction
5. Detail Modal
6. Back to Queue

**Pain Points Addressed:**
- ✅ No need to remember PO details - all info visible
- ✅ OBC Master specs immediately visible
- ✅ Priority clear dengan color coding
- ✅ Material photos available untuk reference

---

### Scenario 2: Budi Uses Search to Find Specific OBC

```
📍 START: Cetak Queue Page
    │
    ├─▶ Type in Search Box: "OBC-2025-001"
    │   ├─ Input debounced (300ms delay)
    │   └─ Search icon animates
    │
    ├─▶ Queue Filters Automatically
    │   ├─ Loading indicator shows
    │   ├─ API called with search param
    │   └─ Results filtered
    │
    ├─▶ View Filtered Results
    │   ├─ Only matching items shown
    │   └─ Result count updated: "1 PO ditemukan"
    │
    ├─▶ Clear Search (if needed)
    │   ├─ Click X icon atau clear input
    │   └─ Full queue restored
    │
    └─▶ 📍 END
```

**Time:** ~30 seconds  
**Success Rate:** High (instant feedback)

---

### Scenario 3: Budi Filters by Priority

```
📍 START: Cetak Queue Page
    │
    ├─▶ Click Priority Filter Button
    │   └─ Dropdown opens dengan animation
    │
    ├─▶ Select "Urgent"
    │   ├─ Haptic feedback (10ms)
    │   ├─ Dropdown closes
    │   └─ Filter applied
    │
    ├─▶ Queue Refreshes
    │   ├─ Loading indicator
    │   ├─ Only URGENT items shown
    │   └─ Filter badge shows: "Urgent"
    │
    ├─▶ Work on Urgent Items
    │   └─ Can see all urgent items prioritized
    │
    ├─▶ Clear Filter
    │   ├─ Select "Semua Prioritas"
    │   └─ Full queue restored
    │
    └─▶ 📍 END
```

**Time:** ~15 seconds  
**Satisfaction:** High (easy filtering)

---

### Scenario 4: Budi Uses Pull-to-Refresh on Mobile

```
📍 START: Cetak Queue Page (on Mobile/Tablet)
    │
    ├─▶ User at Top of Page (scrollY = 0)
    │   └─ Sees queue items
    │
    ├─▶ Pull Down Gesture
    │   ├─ Touchstart detected
    │   ├─ Pull indicator appears
    │   └─ Text: "Tarik untuk refresh..."
    │
    ├─▶ Pull > 80px
    │   ├─ Indicator changes: "Lepas untuk refresh..."
    │   └─ Icon starts spinning
    │
    ├─▶ Release Touch
    │   ├─ Haptic feedback (20ms vibration)
    │   ├─ Queue refresh triggered
    │   ├─ Loading indicator shows
    │   └─ Indicator disappears after refresh
    │
    ├─▶ New Data Loaded
    │   ├─ Queue updated with latest data
    │   ├─ New items highlighted (if any)
    │   └─ Smooth animation
    │
    └─▶ 📍 END
```

**Time:** ~2-3 seconds (including animation)  
**Delight Factor:** High (iOS-like experience)

---

### Scenario 5: Budi Views Empty Queue

```
📍 START: Cetak Queue Page
    │
    ├─▶ Page Loads
    │   └─ No items dengan status READY_FOR_CETAK
    │
    ├─▶ Empty State Displayed
    │   ├─ Glass card with gradient icon
    │   ├─ Icon pop-in animation (spring physics)
    │   ├─ Title: "Antrian Kosong"
    │   └─ Message: "Belum ada PO yang siap untuk dicetak..."
    │
    ├─▶ User Understands Context
    │   └─ Knows to wait for Khazwal to finish prep
    │
    ├─▶ User Refreshes (Pull-to-refresh atau button)
    │   └─ Check for new items
    │
    └─▶ 📍 END
```

**Time:** Immediate understanding  
**Clarity:** High (clear messaging)

---

## Interaction Details

### Card Interactions

**States:**
1. **Default:** Glass card dengan shadow
2. **Hover (Desktop):** Shadow increases
3. **Active (Touch):** Scale 0.97 untuk press feedback
4. **Loading:** Skeleton animation

**Visual Hierarchy:**
```
┌─────────────────────────────────────┐
│ PO #123456          [URGENT badge]  │  ← Primary info
│ OBC: OBC-2025-001                   │  ← Bold indigo
├─────────────────────────────────────┤
│ Buku Pemilik Kendaraan Bermotor    │  ← Product name
├─────────────────────────────────────┤
│ [Material: BPKB] [Seri: A]         │  ← OBC badges
│ [Warna: Biru] [Perso]              │
├─────────────────────────────────────┤
│ ✓ Material Siap | 📦 5,000         │  ← Status + Qty
│ 📅 15 Jan 2025  | 15 hari lagi     │  ← Due date
├─────────────────────────────────────┤
│ [JD] John Doe   | 📸 3 photos      │  ← Prepared by
└─────────────────────────────────────┘
```

---

### Modal Interactions

**Opening:**
- Backdrop fade-in (0.2s)
- Modal slide-up + scale (spring physics)
- Focus trap activated

**Closing:**
- Click backdrop
- Click X button
- Press Escape key
- Swipe down (mobile)

---

### Pull-to-Refresh Behavior

**Trigger Conditions:**
- ✅ Must be at scrollY = 0
- ✅ Pull distance > 80px
- ✅ Not already loading
- ✅ Touch device

**Visual Feedback:**
```
Pull < 80px:  "Tarik untuk refresh..."  [static icon]
Pull > 80px:  "Lepas untuk refresh..."  [spinning icon]
Released:     [Fade out, refresh starts]
```

---

## Error Handling

### Network Error

```
📍 User Tries to Load Queue
    │
    ├─▶ Network Request Fails
    │   └─ Error caught by store
    │
    ├─▶ Error Alert Displayed
    │   ├─ Title: "Gagal memuat antrian cetak"
    │   ├─ Detail: "Silakan periksa koneksi Anda"
    │   └─ Action: Retry button
    │
    ├─▶ User Clicks Retry
    │   └─ Retry fetch
    │
    └─▶ 📍 END
```

---

### Unauthorized Access

```
📍 User Token Expired
    │
    ├─▶ API Returns 401
    │   └─ Interceptor catches
    │
    ├─▶ Auto Redirect to Login
    │   ├─ Auth store cleared
    │   └─ Redirect: /login
    │
    ├─▶ User Re-authenticates
    │   └─ Return to original page
    │
    └─▶ 📍 END
```

---

## Success Metrics

### Performance

| Metric | Target | Sprint 6 Actual |
|--------|--------|-----------------|
| Page Load Time | < 2s | ~1.5s |
| Queue Fetch | < 1s | ~800ms |
| Detail Modal Open | < 500ms | ~300ms |
| Pull-to-Refresh | < 1.5s | ~1.2s |
| Search Response | < 300ms | ~250ms |

### User Satisfaction

- **Task Completion Rate:** 95%+ (easy to find & view PO)
- **Error Rate:** < 2% (robust error handling)
- **Time to Find PO:** < 30s (with search)
- **Mobile Satisfaction:** High (pull-to-refresh delight)

---

## Accessibility Considerations

### Keyboard Navigation

- Tab: Navigate between cards
- Enter/Space: Open detail
- Escape: Close modal
- Arrow keys: Navigate pagination

### Screen Reader

- Card labels: "PO number, OBC number, Priority, Due date"
- Loading states: "Loading queue items"
- Empty states: "Queue is empty, waiting for materials"
- Errors: "Error loading queue, retry button available"

---

## Mobile UX Optimizations

### Touch Targets

- Minimum 44x44px (iOS guideline)
- Cards: Full card tappable
- Buttons: Adequate spacing

### Gestures

- ✅ Pull-to-refresh (custom implementation)
- ✅ Swipe to navigate (browser default)
- ✅ Pinch to zoom photos (in modal)

### Performance

- GPU-accelerated animations (transform, opacity)
- Debounced search input
- Lazy load images (material photos)
- Stagger delay minimal (0.05s max)

---

## Related Documentation

- **API Reference:** [Cetak API](../../04-api-reference/cetak.md)
- **Sprint Documentation:** [Sprint 6](../../10-sprints/sprint-khazwal-sprint6.md)
- **Testing Guide:** [Sprint 6 Testing](../../06-testing/khazwal-sprint6-testing.md)
- **Design Standards:** [Design Standard Rules](../../../.cursor/rules/design-standard.mdc)

---

*Last Updated: 30 Desember 2025*  
*Version: 1.0.0*  
*Sprint: Sprint 6 - Consumer Side & Polish*
