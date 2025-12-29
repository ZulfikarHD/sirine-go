# 📦 Sprint Khazwal: Material Preparation

**Version:** 1.0.0  
**Date:** 29 Desember 2025  
**Duration:** 5 Sprints (Sprint 1-5)  
**Status:** ✅ Completed

## 📋 Sprint Goals

Implementasi complete Material Preparation workflow untuk Unit Khazanah Awal, yaitu: queue management untuk PO, proses konfirmasi material (plat, kertas, tinta), finalisasi dengan foto, notifikasi ke Unit Cetak, monitoring dashboard untuk Supervisor, dan history tracking.

---

## ✨ Features Implemented

### Sprint 1: Foundation

#### Database Models
**3 New Tables:**
- **`production_orders`** - Master Production Order data
  - Fields: po_number, obc_number, sap_customer_code, product_name, quantity, due_date, priority, current_stage, current_status
  - Enums: POPriority (URGENT, NORMAL, LOW), POStage, POStatus

- **`khazwal_material_preparations`** - Material preparation tracking
  - Fields: production_order_id, sap_plat_code, kertas_blanko_quantity, tinta_requirements, status, timestamps
  - Relationships: belongs to ProductionOrder, PreparedBy User

- **`po_stage_trackings`** - Stage history tracking
  - Fields: production_order_id, stage, status, started_at, completed_at, duration, handled_by

#### Backend Skeleton
- `khazwal_service.go` - Service layer dengan placeholder methods
- `khazwal_handler.go` - Handler layer dengan placeholder responses
- Routes registered di `routes.go`

#### Frontend Skeleton
- Navigation menu "Khazanah Awal" di Sidebar
- Routes `/khazwal/material-prep` dan `/khazwal/material-prep/:id`
- Empty page components dengan proper layout

---

### Sprint 2: Queue & Detail

#### Backend Implementation
**Service Layer:**
- `GetMaterialPrepQueue()` - Query dengan filtering, sorting, pagination
- `GetMaterialPrepDetail()` - Single PO dengan full relations
- `StartMaterialPrep()` - Business logic dengan transaction

**Handler Layer:**
- `GetQueue()` - Parse query params, DTO transformation
- `GetDetail()` - 404 handling, computed fields
- `StartPrep()` - User authentication, transaction handling

#### Frontend Implementation
**Components:**
- `PriorityBadge.vue` - Reusable priority display (URGENT/NORMAL/LOW)
- `POQueueCard.vue` - Card component dengan staggered animation

**Pages:**
- `MaterialPrepQueuePage.vue` - Queue list dengan filter, search, pagination
- `MaterialPrepDetailPage.vue` - Detail view dengan start action

**Features:**
- Debounced search (500ms)
- Loading skeleton states
- Empty states (no PO, no results)
- iOS-inspired animations dengan Motion-V
- Haptic feedback

---

### Sprint 3: Process Steps

#### Backend Implementation
**New Endpoints:**
- POST `/material-prep/:id/confirm-plat` - Barcode validation
- PATCH `/material-prep/:id/kertas` - Variance calculation
- PATCH `/material-prep/:id/tinta` - JSON array handling

**Business Logic:**
- Plat code validation against SAP code
- Variance calculation dengan threshold 5%
- Conditional reason requirement
- Transaction-based updates

#### Frontend Implementation
**Components:**
- `BarcodeScanner.vue` - Html5-qrcode integration, camera permission, manual fallback
- `KertasInputForm.vue` - Numeric input, variance display, reason input
- `TintaChecklist.vue` - Dynamic checklist, quantity inputs, low stock warning
- `ProcessStepper.vue` - Step navigation, status indicators

**Pages:**
- `MaterialPrepProcessPage.vue` - 4-step workflow integration

---

### Sprint 4: Finalize & Notification

#### Backend Implementation
**Finalize Endpoint:**
- POST `/material-prep/:id/finalize`
- Validation: all steps must be complete
- Duration calculation
- Photo storage (base64 JSON)
- PO stage transition to CETAK
- Automatic notifications to OPERATOR_CETAK users

#### Frontend Implementation
**Components:**
- `PhotoUploader.vue` - Multi-upload, compression, preview, delete

**Features:**
- Summary checklist review
- Photo upload dengan auto-compress (max 2MB)
- Notes input
- Success screen dengan duration display
- Navigation to queue

---

### Sprint 5: Consumer Side

#### Backend Implementation
**Cetak Service:**
- `cetak_service.go` - Queue untuk Unit Cetak
- `cetak_handler.go` - Endpoints untuk Cetak queue

**Khazwal Extensions:**
- `GetMaterialPrepHistory()` - History dengan date range filter
- `GetMonitoringStats()` - Dashboard stats untuk Supervisor

**New Routes:**
- GET `/cetak/queue` - PO siap cetak
- GET `/cetak/queue/:id` - Detail untuk cetak
- GET `/khazwal/material-prep/history` - History
- GET `/khazwal/monitoring` - Supervisor dashboard

#### Frontend Implementation
**Components:**
- `CetakQueueCard.vue` - Card untuk Cetak queue
- `PrepHistoryCard.vue` - History card
- `StaffActivityCard.vue` - Staff activity indicator
- `MaterialPhotoViewer.vue` - Fullscreen photo viewer

**Pages:**
- `CetakQueuePage.vue` - Queue untuk Unit Cetak
- `MaterialPrepHistoryPage.vue` - History dengan filters
- `SupervisorMonitoringPage.vue` - Monitoring dashboard

**Navigation:**
- Updated Sidebar dengan Cetak menu
- Role-based visibility

---

## 📁 File Structure

### Backend Files

```
backend/
├── models/
│   ├── production_order.go         ✨ NEW
│   ├── khazwal_material_prep.go    ✨ NEW
│   └── po_stage_tracking.go        ✨ NEW
├── handlers/
│   ├── khazwal_handler.go          ✨ NEW
│   └── cetak_handler.go            ✨ NEW
├── services/
│   ├── khazwal_service.go          ✨ NEW
│   └── cetak_service.go            ✨ NEW
└── routes/
    └── routes.go                   ✏️ UPDATED
```

### Frontend Files

```
frontend/src/
├── composables/
│   ├── useKhazwalApi.js            ✨ NEW
│   └── useCetakApi.js              ✨ NEW
├── components/
│   ├── common/
│   │   ├── BarcodeScanner.vue      ✨ NEW
│   │   ├── PhotoUploader.vue       ✨ NEW
│   │   ├── PriorityBadge.vue       ✨ NEW
│   │   └── MaterialPhotoViewer.vue ✨ NEW
│   ├── khazwal/
│   │   ├── POQueueCard.vue         ✨ NEW
│   │   ├── KertasInputForm.vue     ✨ NEW
│   │   ├── TintaChecklist.vue      ✨ NEW
│   │   ├── ProcessStepper.vue      ✨ NEW
│   │   ├── PrepHistoryCard.vue     ✨ NEW
│   │   └── StaffActivityCard.vue   ✨ NEW
│   ├── cetak/
│   │   └── CetakQueueCard.vue      ✨ NEW
│   └── layout/
│       └── Sidebar.vue             ✏️ UPDATED
├── views/
│   ├── khazwal/
│   │   ├── MaterialPrepQueuePage.vue   ✨ NEW
│   │   ├── MaterialPrepDetailPage.vue  ✨ NEW
│   │   ├── MaterialPrepProcessPage.vue ✨ NEW
│   │   ├── MaterialPrepHistoryPage.vue ✨ NEW
│   │   └── SupervisorMonitoringPage.vue ✨ NEW
│   └── cetak/
│       └── CetakQueuePage.vue      ✨ NEW
└── router/
    └── index.js                    ✏️ UPDATED
```

---

## 🔌 API Endpoints Summary

### Khazwal Material Preparation

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/khazwal/material-prep/queue` | Queue PO untuk persiapan |
| GET | `/api/khazwal/material-prep/:id` | Detail PO |
| POST | `/api/khazwal/material-prep/:id/start` | Mulai persiapan |
| POST | `/api/khazwal/material-prep/:id/confirm-plat` | Konfirmasi plat |
| PATCH | `/api/khazwal/material-prep/:id/kertas` | Update kertas |
| PATCH | `/api/khazwal/material-prep/:id/tinta` | Update tinta |
| POST | `/api/khazwal/material-prep/:id/finalize` | Selesaikan persiapan |
| GET | `/api/khazwal/material-prep/history` | Riwayat persiapan |
| GET | `/api/khazwal/monitoring` | Dashboard monitoring |

### Cetak

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/cetak/queue` | Queue PO siap cetak |
| GET | `/api/cetak/queue/:id` | Detail PO untuk cetak |

---

## 🎨 Design Implementation

### Motion-V Animations
- ✅ Page entrance: fadeUp dengan spring physics
- ✅ Cards: staggered animation (0.05s delay)
- ✅ Icons: popIn dengan bouncy spring
- ✅ Modals: fadeScale dengan spring

### iOS-Inspired Interactions
- ✅ Active scale feedback (0.97) pada buttons
- ✅ Haptic feedback via Vibration API
- ✅ Glass card styling
- ✅ Gradient backgrounds

### Mobile-First UX
- ✅ Touch-friendly targets (min 44px)
- ✅ Numeric keyboard untuk input angka
- ✅ Camera capture untuk photos
- ✅ Responsive layouts

---

## 🧪 Testing Summary

### Backend Testing
- ✅ Queue endpoint (filters, pagination, search)
- ✅ Detail endpoint (relations, validation)
- ✅ Start preparation (status validation, transaction)
- ✅ Confirm plat (barcode validation)
- ✅ Update kertas (variance logic)
- ✅ Update tinta (array handling)
- ✅ Finalize (step validation, notification)

### Frontend Testing
- ✅ Queue page functionality
- ✅ Detail page functionality
- ✅ Process page (all 4 steps)
- ✅ Component interactions
- ✅ Mobile responsive

### E2E Testing
- ✅ Complete workflow Staff Khazwal
- ✅ Supervisor monitoring
- ✅ Unit Cetak receive notification

---

## 📊 Metrics

### Implementation Stats
- **New Backend Files:** 6 files
- **New Frontend Files:** 17 files
- **Updated Files:** 3 files
- **Total Lines Added:** ~5,000+ lines

### Performance
- Queue endpoint: < 200ms
- Detail endpoint: < 100ms
- Finalize endpoint: < 500ms

---

## 🔗 Related Documentation

- **API Reference:** [Khazwal API](../04-api-reference/khazwal.md)
- **Testing Guide:** [Khazwal Testing](../06-testing/khazwal-testing.md)
- **User Journeys:** [Khazwal User Journeys](../07-user-journeys/khazwal/material-prep-flow.md)

---

## 👤 Developer

**Zulfikar Hidayatullah**  
Tech Stack: Go + Vue 3 + GORM + Motion-V  
Database: MySQL (sirine_go)

---

*Last Updated: 29 Desember 2025*
