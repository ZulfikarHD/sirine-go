# Development Strategy - Epic 01: Penyiapan Material

**Epic ID:** KHAZWAL-EPIC-01  
**Source:** `Sirine-Reqs/Backlog/01_Khazwal/Epic_01_Material_Preparation.md`  
**Created:** 29 December 2025  
**Updated:** 30 December 2025 - Phase 6 Integration  
**Developer:** Zulfikar Hidayatullah

---

## 🔔 IMPORTANT UPDATE: Phase 6 OBC Master Integration

**Phase 6 (Implementation) telah diupdate** berdasarkan OBC Master architecture yang telah diimplementasi (Phases 1-3). 

**📋 Lihat dokumentasi lengkap:**
- **[Phase 6 Full Plan](./01_Material_Preparation_Phase6.md)** - Complete implementation dengan OBC Master integration
- **[Integration Summary](./INTEGRATION_SUMMARY.md)** - Quick reference untuk key changes
- **[Before/After Comparison](./BEFORE_AFTER_COMPARISON.md)** - Visual comparison perubahan

**Key Changes:**
- ✅ ProductionOrder sekarang memiliki `obc_master_id` FK (bukan lagi direct OBC fields)
- ✅ Semua queries harus preload `OBCMaster` relationship
- ✅ Access pattern berubah: `po.obc_number` → `po.obc_master.obc_number`
- ✅ 39 OBC fields sekarang tersedia via relationship (vs 5 old fields)
- ✅ Frontend components harus handle `obc_master` nested object

**Prerequisites sebelum Phase 6:**
- ⏳ OBC Master Phase 4 (database migration) harus selesai
- ⏳ OBC Master Phase 5 (service layer fixes) harus selesai

**Phases 1-5 di bawah ini tetap relevan untuk understanding dan planning, namun Phase 6 implementation mengacu ke dokumen terpisah di atas.**

---

## 📊 PHASE 1: FEATURE UNDERSTANDING

### 1.1 Data yang Dikelola

| Data | Deskripsi | Source |
|------|-----------|--------|
| Production Orders (PO) | Master PO dengan informasi OBC, quantity, due date | Database lokal + SAP |
| Material Preparations | Record proses penyiapan material per PO | Database lokal |
| Plat Cetak | Kode plat untuk cetak | SAP API |
| Kertas Blanko | Jumlah kertas yang disiapkan | SAP Inventory |
| Tinta | Jenis dan jumlah tinta per warna | SAP Inventory |
| Notifications | Alert ke Unit Cetak | Database lokal |

### 1.2 Owner & Consumer Mapping

| Feature | Owner (Siapa Membuat) | Consumer (Siapa Melihat) | Data Flow |
|---------|----------------------|-------------------------|-----------|
| **Daftar PO Queue** | System (auto dari PO) | Staff Khazwal | PO Created → Queue → Display |
| **Start Preparation** | Staff Khazwal | Staff Khazwal, Supervisor | Staff Starts → Record Created → Tracking |
| **Confirm Plat** | Staff Khazwal (via scan) | Supervisor, Manager | Staff Scans → Validated → Recorded |
| **Input Kertas** | Staff Khazwal | Supervisor, Inventory System | Staff Input → Validated → SAP Updated |
| **Confirm Tinta** | Staff Khazwal | Supervisor, Inventory System | Staff Checks → Validated → SAP Updated |
| **Finalize Prep** | Staff Khazwal | Unit Cetak, Supervisor, Manager | Staff Finalizes → Notification → Cetak Queue |

### 1.3 Primary User Goals

**Staff Khazwal:**
- Melihat prioritas pekerjaan hari ini
- Memproses PO secara sistematis
- Mencatat progress setiap langkah
- Menyelesaikan persiapan material dengan benar

**Supervisor Khazwal:**
- Monitor progress tim
- Review variance/masalah
- Approve jika ada exceptional cases

**Unit Cetak:**
- Menerima notifikasi material siap
- Melihat PO yang ready untuk cetak

---

## 📊 PHASE 2: CROSS-FRONTEND IMPACT MAPPING

### 2.1 Frontend Applications Affected

| Frontend | Role | Impact Level |
|----------|------|--------------|
| **Web (Staff)** | Staff Khazwal operasional | 🔴 HIGH - Primary User |
| **Web (Supervisor)** | Monitoring & approval | 🟡 MEDIUM |
| **Web (Admin/Manager)** | Dashboard overview | 🟢 LOW |
| **Mobile/PWA** | Field operations | 🔴 HIGH - Primary User |

### 2.2 Detailed Impact Matrix

| Feature | Web Staff | Web Supervisor | Web Admin | Mobile PWA |
|---------|-----------|----------------|-----------|------------|
| **US-KW-001: Queue List** | ✅ Full | ✅ Read-only | ✅ Dashboard widget | ✅ Full |
| **US-KW-002: Start Prep** | ✅ Full | ✅ Monitor | ❌ N/A | ✅ Full |
| **US-KW-003: Confirm Plat** | ✅ Full (+ Scanner) | ✅ Monitor | ❌ N/A | ✅ Full + Camera |
| **US-KW-004: Input Kertas** | ✅ Full | ✅ Monitor | ❌ N/A | ✅ Full |
| **US-KW-005: Confirm Tinta** | ✅ Full | ✅ Monitor | ❌ N/A | ✅ Full |
| **US-KW-006: Finalize** | ✅ Full | ✅ Monitor | ❌ N/A | ✅ Full |

---

## 📊 PHASE 3: MISSING IMPLEMENTATION DETECTION

### 3.1 Owner Side (Data Creation) Checklist

#### US-KW-001: Queue List
| Item | Status | Notes |
|------|--------|-------|
| ✅ UI list untuk melihat data | Specified | `MaterialPrepQueuePage.vue` |
| ✅ Filter UI | Specified | Status, tanggal, prioritas |
| ✅ Search UI | Specified | Nomor PO/OBC |
| ⚠️ Empty state (no PO) | **MISSING** | Perlu design untuk state kosong |
| ⚠️ Pull-to-refresh (mobile) | **MISSING** | Best practice untuk mobile |
| ✅ Skeleton loading | Specified | Task KW-001-FE-06 |

#### US-KW-002: Start Preparation
| Item | Status | Notes |
|------|--------|-------|
| ✅ Start button | Specified | Dengan confirmation |
| ✅ Material display | Specified | Plat, kertas, tinta |
| ⚠️ Cancel/abort preparation | **MISSING** | Jika staff salah pilih PO |
| ⚠️ Re-assign to other staff | **MISSING** | Supervisor capability |
| ✅ Loading & feedback | Specified | Task KW-002-FE-04 |

#### US-KW-003: Confirm Plat
| Item | Status | Notes |
|------|--------|-------|
| ✅ Barcode/QR scanner | Specified | Task KW-003-FE-01 |
| ✅ Validation result UI | Specified | Match/mismatch |
| ⚠️ Manual input fallback | **MISSING** | Jika scanner error |
| ⚠️ Camera permission handling | **MISSING** | Mobile/PWA requirement |
| ⚠️ Offline scan queue | **MISSING** | Network reliability |

#### US-KW-004: Input Kertas
| Item | Status | Notes |
|------|--------|-------|
| ✅ Input form | Specified | Jumlah kertas |
| ✅ Target vs actual comparison | Specified | Real-time |
| ✅ Variance warning | Specified | > 5% |
| ✅ Reason field | Specified | Conditional |
| ⚠️ Numeric keypad (mobile) | **MISSING** | Better UX |
| ⚠️ Unit conversion helper | **MISSING** | Lembar besar explanation |

#### US-KW-005: Confirm Tinta
| Item | Status | Notes |
|------|--------|-------|
| ✅ Checklist form | Specified | Per warna |
| ✅ Quantity input per warna | Specified | Dalam kg |
| ✅ Low stock warning | Specified | < 10kg |
| ⚠️ Color visual indicator | **MISSING** | Better UX with color swatches |
| ⚠️ Stock level real-time | **MISSING** | Fetch from SAP |

#### US-KW-006: Finalize
| Item | Status | Notes |
|------|--------|-------|
| ✅ Summary checklist | Specified | All steps |
| ✅ Photo upload | Specified | Optional |
| ✅ Success feedback | Specified | With duration |
| ⚠️ Print summary/label | **MISSING** | Physical label untuk palet |
| ⚠️ Undo/correction period | **MISSING** | Grace period jika salah |

### 3.2 Consumer Side (Data Display) Checklist

| Item | Status | Notes |
|------|--------|-------|
| ⚠️ **Unit Cetak queue view** | **MISSING** | Where Cetak sees ready materials |
| ⚠️ **Supervisor monitoring dashboard** | **MISSING** | Progress all staff |
| ⚠️ **Manager overview dashboard** | **MISSING** | Aggregate metrics |
| ⚠️ **Notification center integration** | **PARTIAL** | NotificationCenter exists, needs category |
| ⚠️ **Real-time status updates** | **MISSING** | WebSocket/polling |

### 3.3 Integration Points Checklist

| Item | Status | Notes |
|------|--------|-------|
| ⚠️ Backend models | **MISSING** | Need Go models for khazwal tables |
| ⚠️ API endpoints | **MISSING** | All 6 endpoints need creation |
| ⚠️ Database migration | **MISSING** | khazwal_material_preparations table |
| ⚠️ SAP API integration | **MISSING** | Plat, kertas, tinta data |
| ⚠️ Navigation menu update | **MISSING** | Khazwal menu item |
| ⚠️ Role-based routing | **MISSING** | STAFF_KHAZWAL role |

---

## 📊 PHASE 4: GAP ANALYSIS

### 4.1 Critical Gaps (⚠️ Feature Broken Without This)

| Gap | Impact | Severity |
|-----|--------|----------|
| **No Unit Cetak queue view** | Cetak cannot see ready materials | 🔴 CRITICAL |
| **No khazwal menu/routes** | Staff cannot access feature | 🔴 CRITICAL |
| **No khazwal models/APIs** | No backend functionality | 🔴 CRITICAL |
| **No production_orders model** | Core entity missing | 🔴 CRITICAL |
| **No SAP integration** | Cannot fetch material data | 🔴 CRITICAL |

### 4.2 Functional Gaps (Feature Incomplete)

| Gap | Impact | Severity |
|-----|--------|----------|
| **No cancel/abort preparation** | Cannot undo mistakes | 🟡 MEDIUM |
| **No manual barcode input** | Cannot work if scanner fails | 🟡 MEDIUM |
| **No offline capability** | Cannot work without network | 🟡 MEDIUM |
| **No supervisor monitoring view** | No oversight capability | 🟡 MEDIUM |

### 4.3 UX Gaps (Feature Works But Poor Experience)

| Gap | Impact | Severity |
|-----|--------|----------|
| **No empty state designs** | Confusing when no data | 🟢 LOW |
| **No pull-to-refresh** | Mobile users expect it | 🟢 LOW |
| **No numeric keypad hint** | Slower input on mobile | 🟢 LOW |
| **No color swatches for tinta** | Visual confusion | 🟢 LOW |

---

## 📊 PHASE 5: IMPLEMENTATION SEQUENCING

### 5.1 Dependency Graph

```
                     ┌─────────────────────────────────────┐
                     │   FOUNDATION (Must Build First)      │
                     │                                      │
                     │ 1. production_orders model           │
                     │ 2. khazwal_material_preparations     │
                     │ 3. Navigation & routing              │
                     │ 4. SAP service (mock first)          │
                     └───────────────┬─────────────────────┘
                                     │
         ┌───────────────────────────┼───────────────────────────┐
         │                           │                           │
         ▼                           ▼                           ▼
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│ US-KW-001       │       │ Database        │       │ Notification    │
│ Queue List      │       │ Migration       │       │ Integration     │
│                 │       │                 │       │                 │
│ - API endpoint  │       │ - Tables        │       │ - Categories    │
│ - Queue page    │       │ - Indexes       │       │ - Templates     │
│ - Card component│       │ - Seeder        │       │ - Routing       │
└────────┬────────┘       └─────────────────┘       └─────────────────┘
         │
         ▼
┌─────────────────┐
│ US-KW-002       │
│ Start Prep      │
│                 │
│ - Start API     │
│ - Start page    │
│ - Service logic │
└────────┬────────┘
         │
         ├────────────────┬────────────────┐
         │                │                │
         ▼                ▼                ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│ US-KW-003   │  │ US-KW-004   │  │ US-KW-005   │
│ Plat        │  │ Kertas      │  │ Tinta       │
│             │  │             │  │             │
│ (PARALLEL)  │  │ (PARALLEL)  │  │ (PARALLEL)  │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │                │                │
       └────────────────┼────────────────┘
                        │
                        ▼
              ┌─────────────────┐
              │ US-KW-006       │
              │ Finalize        │
              │                 │
              │ - All steps     │
              │ - Photo upload  │
              │ - Notification  │
              └────────┬────────┘
                       │
                       ▼
              ┌─────────────────┐
              │ CONSUMER SIDE   │
              │                 │
              │ - Cetak queue   │
              │ - Supervisor    │
              │ - Dashboard     │
              └─────────────────┘
```

### 5.2 Priority Matrix

| Priority | Feature | Reason |
|----------|---------|--------|
| **P0** | Foundation (models, routes, nav) | Nothing works without this |
| **P0** | US-KW-001: Queue List | Entry point for all workflow |
| **P0** | US-KW-002: Start Prep | Cannot begin any preparation |
| **P1** | US-KW-003: Confirm Plat | Core workflow step |
| **P1** | US-KW-004: Input Kertas | Core workflow step |
| **P1** | US-KW-005: Confirm Tinta | Core workflow step |
| **P1** | US-KW-006: Finalize | Complete workflow |
| **P1** | Cetak Queue View | Consumer must see results |
| **P2** | Supervisor Dashboard | Monitoring capability |
| **P2** | Offline support | Enhancement |
| **P2** | Print label | Nice to have |

---

## 📊 PHASE 6: DETAILED RECOMMENDATIONS

### 6.1 New Backend Files Needed

```
backend/
├── models/
│   ├── production_order.go       [P0] Core PO model
│   ├── po_stage_tracking.go      [P0] Stage tracking
│   └── khazwal_material_prep.go  [P0] Material preparation
├── handlers/
│   └── khazwal_handler.go        [P0] All khazwal endpoints
├── services/
│   ├── khazwal_service.go        [P0] Business logic
│   └── sap_service.go            [P1] SAP API integration
└── routes/
    └── routes.go                 [P0] Add khazwal routes
```

### 6.2 New Frontend Pages/Routes Needed

| Page | Path | Priority | Purpose |
|------|------|----------|---------|
| `MaterialPrepQueuePage.vue` | `/khazwal/material-prep` | P0 | Daftar PO untuk disiapkan |
| `MaterialPrepDetailPage.vue` | `/khazwal/material-prep/:id` | P0 | Detail & start preparation |
| `MaterialPrepProcessPage.vue` | `/khazwal/material-prep/:id/process` | P1 | Steps: plat, kertas, tinta |
| `MaterialPrepHistoryPage.vue` | `/khazwal/material-prep/history` | P2 | Riwayat persiapan |
| `CetakQueuePage.vue` | `/cetak/queue` | P1 | View for Unit Cetak |

### 6.3 New Frontend Components Needed

| Component | Location | Priority | Purpose |
|-----------|----------|----------|---------|
| `POQueueCard.vue` | `components/khazwal/` | P0 | Card untuk list PO |
| `PriorityBadge.vue` | `components/common/` | P0 | Badge prioritas (🔴🟡🟢) |
| `BarcodeScanner.vue` | `components/common/` | P1 | Barcode/QR scanner |
| `MaterialChecklist.vue` | `components/khazwal/` | P1 | Checklist steps |
| `TintaInputForm.vue` | `components/khazwal/` | P1 | Multi-color tinta form |
| `VarianceAlert.vue` | `components/khazwal/` | P1 | Warning untuk variance |
| `PhotoUploader.vue` | `components/common/` | P2 | Upload foto palet |

### 6.4 Navigation/Menu Changes

```javascript
// router/index.js - New routes
{
  path: '/khazwal',
  name: 'Khazwal',
  meta: { 
    requiresAuth: true,
    roles: ['STAFF_KHAZWAL', 'SUPERVISOR_KHAZWAL', 'MANAGER', 'ADMIN'],
  },
  children: [
    {
      path: 'material-prep',
      name: 'MaterialPrepQueue',
      component: () => import('../views/khazwal/MaterialPrepQueuePage.vue'),
    },
    {
      path: 'material-prep/:id',
      name: 'MaterialPrepDetail',
      component: () => import('../views/khazwal/MaterialPrepDetailPage.vue'),
    },
    {
      path: 'material-prep/:id/process',
      name: 'MaterialPrepProcess',
      component: () => import('../views/khazwal/MaterialPrepProcessPage.vue'),
    },
  ]
}
```

```javascript
// components/layout/Sidebar.vue - New menu item
{
  icon: 'Package',
  label: 'Khazanah Awal',
  path: '/khazwal/material-prep',
  roles: ['STAFF_KHAZWAL', 'SUPERVISOR_KHAZWAL', 'MANAGER', 'ADMIN'],
}
```

### 6.5 API Endpoints Needed

| Method | Endpoint | Priority | Handler Function |
|--------|----------|----------|------------------|
| GET | `/api/khazwal/material-prep/queue` | P0 | `GetMaterialPrepQueue` |
| GET | `/api/khazwal/material-prep/:id` | P0 | `GetMaterialPrepDetail` |
| POST | `/api/khazwal/material-prep/:id/start` | P0 | `StartMaterialPrep` |
| POST | `/api/khazwal/material-prep/:id/confirm-plat` | P1 | `ConfirmPlat` |
| PATCH | `/api/khazwal/material-prep/:id/kertas` | P1 | `UpdateKertas` |
| PATCH | `/api/khazwal/material-prep/:id/tinta` | P1 | `UpdateTinta` |
| POST | `/api/khazwal/material-prep/:id/finalize` | P1 | `FinalizeMaterialPrep` |
| GET | `/api/khazwal/material-prep/history` | P2 | `GetMaterialPrepHistory` |

### 6.6 Database Changes Needed

```sql
-- Priority: P0
-- File: backend/database/migrations/xxx_create_production_orders.sql

CREATE TABLE production_orders (
    id BIGSERIAL PRIMARY KEY,
    po_number BIGINT NOT NULL UNIQUE,
    obc_number VARCHAR(9) NOT NULL,
    sap_customer_code VARCHAR(50),
    sap_product_code VARCHAR(50),
    product_name VARCHAR(255),
    product_specifications JSONB,
    quantity_ordered INTEGER NOT NULL,
    quantity_target_lembar_besar INTEGER NOT NULL,
    estimated_rims INTEGER NOT NULL,
    order_date DATE NOT NULL,
    due_date DATE NOT NULL,
    priority VARCHAR(20) NOT NULL DEFAULT 'NORMAL',
    priority_score INTEGER,
    current_stage VARCHAR(50) NOT NULL,
    current_status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);

CREATE INDEX idx_po_status_priority ON production_orders(current_stage, current_status, priority, due_date);
```

```sql
-- Priority: P0
-- File: backend/database/migrations/xxx_create_khazwal_material_preparations.sql

CREATE TABLE khazwal_material_preparations (
    id BIGSERIAL PRIMARY KEY,
    production_order_id BIGINT NOT NULL UNIQUE,
    sap_plat_code VARCHAR(50) NOT NULL,
    kertas_blanko_quantity INTEGER NOT NULL,
    tinta_requirements JSONB,
    plat_retrieved_at TIMESTAMPTZ,
    kertas_blanko_actual INTEGER,
    kertas_blanko_variance INTEGER,
    kertas_blanko_variance_reason TEXT,
    tinta_actual JSONB,
    material_photos JSONB,
    status VARCHAR(50) DEFAULT 'PENDING',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    duration_minutes INTEGER,
    prepared_by BIGINT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_po_khazwal_prep FOREIGN KEY (production_order_id) REFERENCES production_orders(id),
    CONSTRAINT fk_prepared_by FOREIGN KEY (prepared_by) REFERENCES users(id)
);

CREATE INDEX idx_khazwal_prep_status ON khazwal_material_preparations(status);
CREATE INDEX idx_khazwal_prep_po ON khazwal_material_preparations(production_order_id);
```

---

## 📊 PHASE 7: USER JOURNEYS

### 7.1 Journey 1: Staff Melihat dan Memulai Persiapan (US-KW-001 + US-KW-002)

**Owner Journey (Staff Khazwal):**
1. Staff login dan masuk ke dashboard staff
2. Staff klik menu "Khazanah Awal" → "Persiapan Material" di sidebar
3. Sistem redirect ke `/khazwal/material-prep`
4. Sistem load daftar PO dengan status "Menunggu Penyiapan Material"
5. Staff melihat list PO dengan priority badges (🔴🟡🟢)
6. Staff klik salah satu PO card untuk melihat detail
7. Sistem redirect ke `/khazwal/material-prep/:id`
8. Staff melihat:
   - Info PO (nomor, OBC, quantity, due date)
   - Material yang dibutuhkan (plat, kertas, tinta)
9. Staff klik tombol "Mulai Persiapan"
10. Sistem tampil confirmation dialog
11. Staff konfirmasi → Sistem:
    - Create record di `khazwal_material_preparations`
    - Update `production_orders.current_status` = 'SEDANG_DISIAPKAN'
    - Log ke `activity_logs`
12. Staff diarahkan ke halaman process

**Consumer Journey (Supervisor):**
1. Supervisor melihat dashboard monitoring
2. Dashboard menampilkan:
   - Jumlah PO dalam queue
   - Jumlah PO sedang diproses
   - Staff yang sedang aktif
3. Supervisor bisa klik untuk detail

### 7.2 Journey 2: Staff Konfirmasi Material (US-KW-003 + US-KW-004 + US-KW-005)

**Owner Journey (Staff Khazwal):**
1. Staff berada di `/khazwal/material-prep/:id/process`
2. Sistem tampilkan stepper: Plat → Kertas → Tinta → Selesai

**Step 1 - Plat:**
3. Staff lihat kode plat yang diperlukan
4. Staff klik "Scan Plat" → Camera terbuka
5. Staff scan barcode pada plat fisik
6. Sistem validasi:
   - Match: ✅ checkmark, next step enabled
   - Mismatch: ❌ warning, tidak bisa lanjut
7. Jika match, sistem:
   - Update `plat_retrieved_at`
   - Log ke `activity_logs`

**Step 2 - Kertas:**
8. Sistem tampilkan target kertas: "[X] lembar besar"
9. Staff input jumlah aktual yang disiapkan
10. Sistem hitung variance real-time
11. Jika variance > 5%:
    - Tampil warning banner
    - Muncul input "Alasan Selisih"
12. Staff submit → Sistem:
    - Update `kertas_blanko_actual`, `kertas_blanko_variance`
    - (Future) Update SAP inventory
    - Log ke `activity_logs`

**Step 3 - Tinta:**
13. Sistem tampilkan list warna sesuai PO specs
14. Per warna, staff:
    - Centang checkbox
    - Input quantity (kg)
15. Jika stock warna < 10kg, tampil low stock alert
16. Staff submit → Sistem:
    - Update `tinta_actual` JSONB
    - (Future) Update SAP inventory
    - Log ke `activity_logs`

### 7.3 Journey 3: Staff Finalisasi dan Notifikasi ke Cetak (US-KW-006)

**Owner Journey (Staff Khazwal):**
1. Staff di step terakhir "Selesai"
2. Sistem tampilkan summary checklist:
   - ✅ Plat sudah diambil (dengan timestamp)
   - ✅ Kertas sudah disiapkan (dengan qty)
   - ✅ Tinta sudah disiapkan (dengan breakdown)
3. (Optional) Staff upload foto palet
4. Staff klik "Selesai - Kirim ke Unit Cetak"
5. Sistem tampil confirmation dialog dengan summary
6. Staff konfirmasi → Sistem:
   - Update `khazwal_material_preparations.status` = 'COMPLETED'
   - Calculate `duration_minutes`
   - Update `production_orders.current_stage` = 'CETAK'
   - Update `production_orders.current_status` = 'SIAP_CETAK'
   - Create notification ke semua user Unit Cetak
   - Log ke `activity_logs`
7. Staff melihat success screen dengan:
   - "Material untuk PO [X] berhasil disiapkan!"
   - Durasi persiapan: [X] menit
   - Tombol "Kembali ke Queue"

**Consumer Journey (Unit Cetak):**
1. Operator Cetak sedang bekerja
2. Notifikasi muncul: "Material Siap untuk PO [X]"
3. Operator klik notifikasi
4. Redirect ke `/cetak/queue` (atau detail page)
5. Operator melihat PO dengan status "Siap Cetak"
6. Operator bisa lihat detail material yang sudah disiapkan

---

## 📊 IMPLEMENTATION PLAN - Sprint Breakdown

### Sprint 1: Foundation (Week 1)

| Task | Priority | Effort | Assignee |
|------|----------|--------|----------|
| Create `production_orders` model + migration | P0 | 4h | Backend |
| Create `khazwal_material_preparations` model + migration | P0 | 3h | Backend |
| Create `po_stage_tracking` model + migration | P0 | 2h | Backend |
| Setup khazwal routes skeleton | P0 | 2h | Backend |
| Add STAFF_KHAZWAL role | P0 | 1h | Backend |
| Create khazwal menu in Sidebar | P0 | 1h | Frontend |
| Create khazwal routes in router | P0 | 2h | Frontend |
| Create page skeleton (empty) | P0 | 2h | Frontend |
| Seed sample PO data | P0 | 2h | Backend |

**Deliverable:** Navigation works, empty pages accessible

### Sprint 2: Queue & Detail (Week 2)

| Task | Priority | Effort | Assignee |
|------|----------|--------|----------|
| Implement `GET /api/khazwal/material-prep/queue` | P0 | 4h | Backend |
| Implement priority calculation logic | P0 | 2h | Backend |
| Implement filtering & search | P0 | 3h | Backend |
| Create `MaterialPrepQueuePage.vue` | P0 | 4h | Frontend |
| Create `POQueueCard.vue` | P0 | 3h | Frontend |
| Create `PriorityBadge.vue` | P0 | 1h | Frontend |
| Implement skeleton loading | P0 | 2h | Frontend |
| Create filter/search UI | P0 | 3h | Frontend |
| Mobile responsive | P0 | 2h | Frontend |

**Deliverable:** Staff can see and search PO queue

### Sprint 3: Start & Detail (Week 3)

| Task | Priority | Effort | Assignee |
|------|----------|--------|----------|
| Implement `GET /api/khazwal/material-prep/:id` | P0 | 3h | Backend |
| Implement `POST /api/khazwal/material-prep/:id/start` | P0 | 4h | Backend |
| Create MaterialPreparationService | P0 | 4h | Backend |
| Mock SAP service for materials | P1 | 3h | Backend |
| Create `MaterialPrepDetailPage.vue` | P0 | 4h | Frontend |
| Implement start preparation flow | P0 | 3h | Frontend |
| Add confirmation dialog | P0 | 2h | Frontend |
| Loading & feedback states | P0 | 2h | Frontend |

**Deliverable:** Staff can view detail and start preparation

### Sprint 4: Process Steps (Week 4)

| Task | Priority | Effort | Assignee |
|------|----------|--------|----------|
| Implement `POST /api/khazwal/material-prep/:id/confirm-plat` | P1 | 3h | Backend |
| Implement `PATCH /api/khazwal/material-prep/:id/kertas` | P1 | 3h | Backend |
| Implement `PATCH /api/khazwal/material-prep/:id/tinta` | P1 | 3h | Backend |
| Implement variance validation | P1 | 2h | Backend |
| Create `MaterialPrepProcessPage.vue` with stepper | P1 | 4h | Frontend |
| Create `BarcodeScanner.vue` | P1 | 4h | Frontend |
| Create kertas input with variance alert | P1 | 3h | Frontend |
| Create `TintaInputForm.vue` | P1 | 3h | Frontend |

**Deliverable:** Staff can complete all process steps

### Sprint 5: Finalize & Notification (Week 5)

| Task | Priority | Effort | Assignee |
|------|----------|--------|----------|
| Implement `POST /api/khazwal/material-prep/:id/finalize` | P1 | 4h | Backend |
| Implement notification to Cetak users | P1 | 3h | Backend |
| Implement activity logging for all steps | P1 | 2h | Backend |
| Create finalization summary UI | P1 | 3h | Frontend |
| Create `PhotoUploader.vue` | P2 | 3h | Frontend |
| Success screen with duration | P1 | 2h | Frontend |
| Cetak notification integration | P1 | 2h | Frontend |

**Deliverable:** Complete workflow end-to-end

### Sprint 6: Consumer Side & Polish (Week 6) ✅ COMPLETED

| Task | Priority | Effort | Status | Notes |
|------|----------|--------|--------|-------|
| Create basic Cetak queue API | P1 | 3h | ✅ Done | Enhanced dengan OBC Master integration |
| Create supervisor monitoring API | P2 | 3h | ✅ Done | Menggunakan existing endpoint |
| Create `CetakQueuePage.vue` | P1 | 4h | ✅ Done | With Cetak store integration |
| Add Cetak menu & routes | P1 | 2h | ✅ Done | Routes configured |
| Supervisor monitoring view | P2 | 4h | ✅ Done | With Khazwal store + OBC context |
| Empty states design | P2 | 2h | ✅ Done | Enhanced dengan animations |
| Pull-to-refresh mobile | P2 | 2h | ✅ Done | Touch gesture implemented |
| E2E testing | P1 | 4h | 📝 Planned | For next phase |

**Deliverable:** ✅ Complete feature with consumer views

**Implementation Summary:**
- Backend: OBC Master fields added to Cetak DTOs
- Frontend: Cetak Pinia store created for state management
- UI: CetakQueueCard enhanced with OBC Master display (Material, Seri, Warna, Personalization)
- UX: Pull-to-refresh for mobile with haptic feedback
- Monitoring: Enhanced StaffActivityCard with OBC context
- Empty States: Improved dengan Motion-V animations

**Documentation:**
- [Sprint 6 Completion Summary](./SPRINT6_COMPLETION_SUMMARY.md)

**Date Completed:** December 30, 2025

---

## 📝 ADDITIONAL NOTES

### Critical Reminders

1. **SAP Integration:** Untuk Phase 1, gunakan mock data. Real SAP integration di Phase 2.

2. **Mobile-First:** Semua UI harus didesain mobile-first. Staff Khazwal akan menggunakan tablet/handphone di lapangan.

3. **Offline Consideration:** Untuk Sprint 6+, pertimbangkan PWA dengan offline queue untuk scan dan input.

4. **Camera Permission:** BarcodeScanner component harus handle permission gracefully.

5. **Real-time Updates:** Untuk supervisor monitoring, consider WebSocket atau polling interval 30 detik.

### Technical Decisions

| Decision | Choice | Reason |
|----------|--------|--------|
| Scanner Library | `@aspect/barcode-scanner` atau `html5-qrcode` | Well-maintained, supports both barcode & QR |
| Stepper Component | Custom dengan Motion-V | Sesuai design standard |
| Photo Upload | S3/MinIO compatible | Scalable storage |
| State Management | Pinia store untuk khazwal | Consistent dengan existing pattern |

### Risk Mitigation

| Risk | Mitigation |
|------|------------|
| SAP API unavailable | Implement mock mode, graceful degradation |
| Scanner fails | Provide manual input fallback |
| Large photo uploads | Compress on client side, max 2MB |
| Concurrent edit same PO | Lock mechanism, show "sedang diproses oleh [nama]" |

---

**Document Status:** Ready for Development  
**Next Steps:** 
1. Review dengan tim
2. Create GitHub issues dari sprint breakdown
3. Setup development branch `feature/khazwal-material-prep`
4. Start Sprint 1
