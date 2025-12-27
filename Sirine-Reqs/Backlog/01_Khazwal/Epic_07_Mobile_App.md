# Epic 7: Mobile App (PWA)

**Epic ID:** KHAZWAL-EPIC-07  
**Priority:** 🟡 Medium (Phase 2)  
**Estimated Duration:** 2 Minggu  

---

## 📋 Overview

Epic ini mencakup pengembangan Progressive Web App (PWA) untuk staff Khazanah Awal agar dapat input data langsung di lapangan tanpa harus ke komputer.

---

## 🗄️ Database Reference

### Tables (Same as desktop)
- All Khazwal tables
- `activity_logs` - Track mobile actions

### Additional Storage
- IndexedDB - Offline data cache
- Service Worker - Background sync

---

## 📝 Backlog Items

### US-KW-024: Mobile App untuk Staff Khazanah Awal

| Field | Value |
|-------|-------|
| **ID** | US-KW-024 |
| **Story Points** | 21 |
| **Priority** | 🟡 Medium |
| **Dependencies** | Epic 1, 2, 3 completed |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin akses aplikasi via mobile, sehingga bisa input data langsung di lapangan tanpa ke komputer.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-024-FE-01 | Configure PWA manifest & service worker | 4h | Frontend |
| KW-024-FE-02 | Setup IndexedDB for offline storage | 4h | Frontend |
| KW-024-FE-03 | Create mobile-optimized layout | 4h | Frontend |
| KW-024-FE-04 | Create mobile PO queue view | 3h | Frontend |
| KW-024-FE-05 | Create mobile process forms | 4h | Frontend |
| KW-024-FE-06 | Integrate barcode/QR scanner (camera API) | 4h | Frontend |
| KW-024-FE-07 | Create mobile photo upload | 3h | Frontend |
| KW-024-FE-08 | Create mobile notification view | 2h | Frontend |
| KW-024-FE-09 | Create mobile performance view | 2h | Frontend |
| KW-024-FE-10 | Implement offline mode (basic view) | 4h | Frontend |
| KW-024-FE-11 | Implement background sync | 4h | Frontend |
| KW-024-FE-12 | Optimize bundle size for fast load | 3h | Frontend |
| KW-024-BE-01 | Create sync API endpoint | 3h | Backend |
| KW-024-BE-02 | Handle conflict resolution | 3h | Backend |
| KW-024-BE-03 | Log device type in activity_logs | 1h | Backend |

#### Acceptance Criteria
- [ ] Responsive mobile web (PWA)
- [ ] Fitur mobile:
  - View daftar PO
  - Mulai/selesai proses (penyiapan/penghitungan/pemotongan)
  - Input data hasil
  - Scan barcode/QR (plat, palet, material)
  - Upload foto (kerusakan, palet)
  - View notification
  - View my performance
- [ ] Offline mode (basic view)
- [ ] Sync otomatis saat online
- [ ] Fast & lightweight (< 2 sec load)

#### PWA Requirements
```javascript
// manifest.json
{
  "name": "Sirine - Khazanah Awal",
  "short_name": "Khazwal",
  "start_url": "/khazwal",
  "display": "standalone",
  "background_color": "#ffffff",
  "theme_color": "#1a365d",
  "icons": [
    { "src": "/icons/192.png", "sizes": "192x192", "type": "image/png" },
    { "src": "/icons/512.png", "sizes": "512x512", "type": "image/png" }
  ]
}
```

#### Offline Strategy
```
1. PENDING queue data - Cache First
2. PO details - Network First with Cache Fallback
3. Form submissions - Queue in IndexedDB, sync when online
4. Photos - Compress & queue for upload
```

#### Mobile UI Considerations
```
- Touch targets: min 44x44 px
- Font size: min 16px body
- Input fields: large, easy to tap
- Buttons: full width on mobile
- Forms: max 3 fields per screen
- Loading: skeleton screens
- Feedback: haptic + visual
```

---

## 📊 Epic Summary

| User Story | Story Points | Priority | Phase |
|------------|--------------|----------|-------|
| US-KW-024 | 21 | Medium | Phase 2 |
| **Total** | **21** | - | - |

---

## 🔗 Dependencies Graph

```
Epic 1, 2, 3 (Core Desktop)
    │
    └── US-KW-024 (Mobile PWA)
            │
            ├── Barcode Scanner
            ├── Camera Upload
            ├── Offline Mode
            └── Background Sync
```

---

## 🧪 Testing Checklist

### Unit Tests
- [ ] IndexedDB operations
- [ ] Offline queue management
- [ ] Sync conflict resolution

### Integration Tests
- [ ] Background sync API
- [ ] Photo upload API (compressed)

### Device Tests
- [ ] iOS Safari
- [ ] Android Chrome
- [ ] Install PWA flow
- [ ] Offline behavior
- [ ] Camera access
- [ ] Push notification

---

## 📱 UI/UX Mobile Guidelines

### Touch Targets
- All interactive elements: min 44×44 px
- Spacing between targets: min 8px

### Typography
- Body: 16px minimum
- Headers: 20-24px
- High contrast ratios

### Forms
- Label above field
- Large input areas
- Clear error messages
- Progress indicator

### Navigation
- Bottom navigation bar
- Swipe gestures where appropriate
- Pull-to-refresh

### Performance
- First Contentful Paint: < 1.5s
- Time to Interactive: < 3s
- Lighthouse PWA score: > 90

---

**Last Updated:** 27 December 2025  
**Status:** Ready for Development
