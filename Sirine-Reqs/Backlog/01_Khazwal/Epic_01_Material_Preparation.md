# Epic 1: Penyiapan Material (Material Preparation)

**Epic ID:** KHAZWAL-EPIC-01  
**Priority:** 🔴 High (Phase 1 - MVP)  
**Estimated Duration:** 2-3 Minggu  

---

## 📋 Overview

Epic ini mencakup seluruh proses penyiapan material produksi di Khazanah Awal, mulai dari melihat daftar PO yang perlu disiapkan hingga finalisasi pengiriman material ke Unit Cetak.

---

## 🗄️ Database Reference

### Primary Table
- `khazwal_material_preparations` - Menyimpan data proses penyiapan material

### Related Tables
- `production_orders` - Referensi PO utama
- `po_stage_tracking` - Tracking perjalanan stage PO
- `users` - Staff yang menangani
- `activity_logs` - Audit trail
- `notifications` - Notifikasi ke Unit Cetak

### SAP Integration (via API)
- Material data (plat, kertas, tinta) dari SAP
- Stok material dari SAP

---

## 📝 Backlog Items

### US-KW-001: Melihat Daftar PO yang Perlu Disiapkan

| Field | Value |
|-------|-------|
| **ID** | US-KW-001 |
| **Story Points** | 5 |
| **Priority** | 🔴 High |
| **Dependencies** | - |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin melihat daftar PO yang perlu disiapkan materialnya, sehingga saya tahu prioritas pekerjaan hari ini.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-001-BE-01 | Create API endpoint `GET /api/khazwal/material-prep/queue` | 4h | Backend |
| KW-001-BE-02 | Implement query filter: status "Menunggu Penyiapan Material" | 2h | Backend |
| KW-001-BE-03 | Implement sorting by priority (urgent, due_date, sequence) | 2h | Backend |
| KW-001-BE-04 | Implement pagination & search (po_number, obc_number) | 2h | Backend |
| KW-001-BE-05 | Calculate priority status (🔴🟡🟢) based on due_date | 1h | Backend |
| KW-001-FE-01 | Create page `MaterialPrepQueuePage.vue` | 3h | Frontend |
| KW-001-FE-02 | Create component `POQueueCard.vue` with priority badge | 2h | Frontend |
| KW-001-FE-03 | Implement filter UI (status, tanggal, prioritas) | 2h | Frontend |
| KW-001-FE-04 | Implement search bar (nomor PO/OBC) | 1h | Frontend |
| KW-001-FE-05 | Make responsive for tablet/mobile | 2h | Frontend |
| KW-001-FE-06 | Implement skeleton loading state | 1h | Frontend |

#### Acceptance Criteria
- [ ] Tampil daftar PO dengan status "Menunggu Penyiapan Material"
- [ ] Sorting berdasarkan prioritas (urgent, due date, sequence)
- [ ] Tampil info: Nomor PO & OBC, Jumlah cetak, Kode plat, Spesifikasi warna, Due date, Status prioritas
- [ ] Filter berdasarkan: status, tanggal, prioritas
- [ ] Search berdasarkan nomor PO/OBC
- [ ] Responsive untuk tablet/mobile

#### Business Rules
```
Priority calculation:
- due_date < 3 hari = URGENT (🔴)
- due_date 3-7 hari = NORMAL (🟡)
- due_date > 7 hari = LOW (🟢)
```

#### Database Query Reference
```sql
SELECT 
    po.id,
    po.po_number,
    po.obc_number,
    po.quantity_ordered,
    po.quantity_target_lembar_besar,
    po.product_specifications->>'plat_code' as plat_code,
    po.product_specifications->>'colors' as colors,
    po.due_date,
    po.priority,
    mp.status as prep_status
FROM production_orders po
LEFT JOIN khazwal_material_preparations mp ON mp.production_order_id = po.id
WHERE po.current_stage = 'KHAZWAL_MATERIAL_PREP'
    AND po.current_status = 'WAITING_MATERIAL_PREP'
    AND po.deleted_at IS NULL
ORDER BY 
    CASE po.priority 
        WHEN 'URGENT' THEN 1 
        WHEN 'NORMAL' THEN 2 
        ELSE 3 
    END,
    po.due_date ASC;
```

---

### US-KW-002: Memulai Proses Penyiapan Material

| Field | Value |
|-------|-------|
| **ID** | US-KW-002 |
| **Story Points** | 8 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-001 |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin memulai proses penyiapan material untuk sebuah PO, sehingga sistem bisa tracking progress pekerjaan saya.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-002-BE-01 | Create API endpoint `POST /api/khazwal/material-prep/{po_id}/start` | 3h | Backend |
| KW-002-BE-02 | Create `MaterialPreparationService.php` | 2h | Backend |
| KW-002-BE-03 | Fetch plat code from SAP API (cache if needed) | 3h | Backend |
| KW-002-BE-04 | Calculate kertas blanko quantity (jumlah_cetak ÷ 2) | 1h | Backend |
| KW-002-BE-05 | Fetch tinta requirements from SAP API | 2h | Backend |
| KW-002-BE-06 | Insert record to `khazwal_material_preparations` | 2h | Backend |
| KW-002-BE-07 | Update `production_orders.current_status` = 'SEDANG_DISIAPKAN' | 1h | Backend |
| KW-002-BE-08 | Insert to `po_stage_tracking` with started_at | 1h | Backend |
| KW-002-BE-09 | Log to `activity_logs` | 1h | Backend |
| KW-002-FE-01 | Create page `MaterialPrepStartPage.vue` | 3h | Frontend |
| KW-002-FE-02 | Display material requirements (plat, kertas, tinta) | 2h | Frontend |
| KW-002-FE-03 | Add "Mulai Persiapan" button with confirmation | 1h | Frontend |
| KW-002-FE-04 | Show loading state and success feedback | 1h | Frontend |

#### Acceptance Criteria
- [ ] Button "Mulai Persiapan" pada detail PO
- [ ] Sistem otomatis tampilkan: Kode plat, Jumlah kertas blanko (÷2), Jenis dan estimasi tinta
- [ ] Status PO berubah menjadi "Sedang Disiapkan"
- [ ] Timestamp mulai penyiapan tercatat
- [ ] Nama staff yang handle tercatat

#### Business Rules
```
Jumlah Kertas Blanko = Jumlah Cetak PO ÷ 2
```

#### Database Insert Reference
```sql
INSERT INTO khazwal_material_preparations (
    production_order_id,
    sap_plat_code,
    kertas_blanko_quantity,
    tinta_requirements,
    status,
    started_at,
    prepared_by,
    created_at
) VALUES (
    :po_id,
    :plat_code,
    :quantity_target / 2,
    :tinta_json,
    'IN_PROGRESS',
    NOW(),
    :user_id,
    NOW()
);
```

---

### US-KW-003: Konfirmasi Pengambilan Plat Cetak

| Field | Value |
|-------|-------|
| **ID** | US-KW-003 |
| **Story Points** | 5 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-002 |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin konfirmasi bahwa plat cetak sudah diambil, sehingga sistem bisa tracking ketersediaan plat.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-003-BE-01 | Create API `POST /api/khazwal/material-prep/{id}/confirm-plat` | 2h | Backend |
| KW-003-BE-02 | Validate plat code via barcode/QR scan | 3h | Backend |
| KW-003-BE-03 | Update `plat_retrieved_at` timestamp | 1h | Backend |
| KW-003-BE-04 | Update plat status in SAP (Sedang Digunakan) - optional | 2h | Backend |
| KW-003-BE-05 | Log to `activity_logs` | 1h | Backend |
| KW-003-FE-01 | Create barcode/QR scanner component | 4h | Frontend |
| KW-003-FE-02 | Show validation result (match/mismatch) | 2h | Frontend |
| KW-003-FE-03 | Display warning if plat code mismatch | 1h | Frontend |
| KW-003-FE-04 | Update UI state after confirmation | 1h | Frontend |

#### Acceptance Criteria
- [ ] Checkbox "Plat Sudah Diambil" dengan scan barcode/QR
- [ ] Sistem validasi kode plat sesuai dengan yang tertera di PO
- [ ] Jika salah plat, tampil warning dan tidak bisa lanjut
- [ ] Timestamp pengambilan plat tercatat
- [ ] Status plat berubah menjadi "Sedang Digunakan"

#### Business Rules
```
- Plat harus sesuai dengan kode plat di PO
- 1 Plat hanya bisa digunakan untuk 1 PO dalam waktu bersamaan
```

---

### US-KW-004: Input Jumlah Kertas Blanko yang Disiapkan

| Field | Value |
|-------|-------|
| **ID** | US-KW-004 |
| **Story Points** | 5 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-002 |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin input jumlah kertas blanko yang sudah disiapkan, sehingga ada record jika ada selisih dengan target.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-004-BE-01 | Create API `PATCH /api/khazwal/material-prep/{id}/kertas` | 2h | Backend |
| KW-004-BE-02 | Calculate variance (actual - target) | 1h | Backend |
| KW-004-BE-03 | Validate variance threshold (±5%) | 1h | Backend |
| KW-004-BE-04 | Require reason if variance > 5% | 1h | Backend |
| KW-004-BE-05 | Update SAP stock (consumption) - via API | 3h | Backend |
| KW-004-BE-06 | Log to `activity_logs` | 1h | Backend |
| KW-004-FE-01 | Create input form for kertas blanko quantity | 2h | Frontend |
| KW-004-FE-02 | Show target vs actual comparison | 1h | Frontend |
| KW-004-FE-03 | Show variance percentage in real-time | 1h | Frontend |
| KW-004-FE-04 | Show warning & reason field if variance > 5% | 2h | Frontend |
| KW-004-FE-05 | Validation: tidak boleh negatif | 1h | Frontend |

#### Acceptance Criteria
- [ ] Input field "Jumlah Kertas Disiapkan" (dalam lembar besar)
- [ ] Tampil target: [Target] lembar besar
- [ ] Sistem otomatis hitung selisih (jika ada)
- [ ] Jika selisih > 5%, tampil warning dan wajib isi alasan
- [ ] Sistem kurangi stok kertas blanko di inventory (SAP)
- [ ] Timestamp input tercatat

#### Business Rules
```
Target = Jumlah Cetak PO ÷ 2
Toleransi selisih: ±5%
Jika kurang/lebih > 5%, wajib isi alasan
```

#### Database Update Reference
```sql
UPDATE khazwal_material_preparations
SET 
    kertas_blanko_actual = :actual_quantity,
    kertas_blanko_variance = :actual_quantity - kertas_blanko_quantity,
    kertas_blanko_variance_reason = :reason,
    updated_at = NOW()
WHERE id = :id;
```

---

### US-KW-005: Konfirmasi Penyiapan Tinta

| Field | Value |
|-------|-------|
| **ID** | US-KW-005 |
| **Story Points** | 5 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-002 |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin konfirmasi jenis dan jumlah tinta yang sudah disiapkan, sehingga Unit Cetak tahu tinta apa saja yang tersedia.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-005-BE-01 | Create API `PATCH /api/khazwal/material-prep/{id}/tinta` | 2h | Backend |
| KW-005-BE-02 | Validate tinta per warna (checklist) | 2h | Backend |
| KW-005-BE-03 | Update SAP stock per warna (consumption) | 3h | Backend |
| KW-005-BE-04 | Check minimum stock alert (< 10kg) | 1h | Backend |
| KW-005-BE-05 | Create alert if stock < minimum | 2h | Backend |
| KW-005-BE-06 | Log to `activity_logs` | 1h | Backend |
| KW-005-FE-01 | Create tinta checklist form (per warna) | 3h | Frontend |
| KW-005-FE-02 | Input quantity per warna (kg) | 2h | Frontend |
| KW-005-FE-03 | Show low stock warning per warna | 1h | Frontend |
| KW-005-FE-04 | Validation: semua warna harus diisi | 1h | Frontend |

#### Acceptance Criteria
- [ ] Checklist untuk setiap warna tinta sesuai spesifikasi PO
- [ ] Input jumlah tinta per warna (dalam kg)
- [ ] Sistem kurangi stok tinta di inventory (SAP)
- [ ] Jika stok tinta < minimum, tampil alert
- [ ] Timestamp konfirmasi tercatat

#### Business Rules
```
Warna tinta sesuai spesifikasi di OBC/PO
Minimum stok tinta per warna: 10 kg
Alert jika stok < minimum
```

#### JSONB Structure for tinta_actual
```json
[
  {"sap_material_code": "TINTA-RED", "color": "Merah", "quantity_kg": 5.5},
  {"sap_material_code": "TINTA-BLUE", "color": "Biru", "quantity_kg": 3.2}
]
```

---

### US-KW-006: Finalisasi Penyiapan Material

| Field | Value |
|-------|-------|
| **ID** | US-KW-006 |
| **Story Points** | 8 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-003, US-KW-004, US-KW-005 |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin finalisasi bahwa semua material sudah siap di palet, sehingga Unit Cetak bisa mulai proses cetak.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-006-BE-01 | Create API `POST /api/khazwal/material-prep/{id}/finalize` | 3h | Backend |
| KW-006-BE-02 | Validate all checklist complete (plat, kertas, tinta) | 2h | Backend |
| KW-006-BE-03 | Update status to 'COMPLETED' | 1h | Backend |
| KW-006-BE-04 | Calculate duration_minutes | 1h | Backend |
| KW-006-BE-05 | Update `production_orders.current_status` = 'SIAP_CETAK' | 1h | Backend |
| KW-006-BE-06 | Update `production_orders.current_stage` = 'CETAK' | 1h | Backend |
| KW-006-BE-07 | Update `po_stage_tracking` with completed_at | 1h | Backend |
| KW-006-BE-08 | Create notification to Unit Cetak | 2h | Backend |
| KW-006-BE-09 | Handle photo upload (optional) | 2h | Backend |
| KW-006-BE-10 | Log to `activity_logs` | 1h | Backend |
| KW-006-FE-01 | Create finalization page with checklist summary | 3h | Frontend |
| KW-006-FE-02 | Validate all steps completed before enable button | 2h | Frontend |
| KW-006-FE-03 | Add photo upload component (optional) | 2h | Frontend |
| KW-006-FE-04 | Show success message with duration summary | 1h | Frontend |
| KW-006-FE-05 | Redirect to queue after success | 1h | Frontend |

#### Acceptance Criteria
- [ ] Button "Selesai - Kirim ke Unit Cetak"
- [ ] Sistem validasi semua checklist sudah complete:
  - ✅ Plat sudah diambil
  - ✅ Kertas sudah disiapkan
  - ✅ Tinta sudah disiapkan
- [ ] Upload foto palet (opsional)
- [ ] Status PO berubah menjadi "Siap Cetak"
- [ ] Notifikasi ke Unit Cetak bahwa material sudah siap
- [ ] Timestamp selesai penyiapan tercatat
- [ ] Sistem hitung durasi penyiapan

#### Business Rules
```
Semua checklist harus complete sebelum bisa finalisasi
Durasi penyiapan = Timestamp Selesai - Timestamp Mulai
```

#### Database Update Reference
```sql
-- Update material preparation
UPDATE khazwal_material_preparations
SET 
    status = 'COMPLETED',
    completed_at = NOW(),
    duration_minutes = EXTRACT(EPOCH FROM (NOW() - started_at)) / 60,
    material_photos = :photos_json,
    updated_at = NOW()
WHERE id = :id;

-- Update production order
UPDATE production_orders
SET 
    current_stage = 'CETAK',
    current_status = 'SIAP_CETAK',
    updated_at = NOW()
WHERE id = :po_id;

-- Insert notification
INSERT INTO notifications (
    user_id,
    notification_type,
    notification_category,
    title,
    message,
    reference_type,
    reference_id,
    priority,
    created_at
) 
SELECT 
    u.id,
    'INFO',
    'PRODUCTION',
    'Material Siap untuk PO ' || :po_number,
    'Material untuk PO ' || :po_number || ' sudah siap di palet. Silakan mulai proses cetak.',
    'PRODUCTION_ORDER',
    :po_id,
    'NORMAL',
    NOW()
FROM users u
WHERE u.department = 'CETAK' AND u.role IN ('OPERATOR_CETAK', 'SUPERVISOR_CETAK');
```

---

## 📊 Epic Summary

| User Story | Story Points | Priority | Phase |
|------------|--------------|----------|-------|
| US-KW-001 | 5 | High | MVP |
| US-KW-002 | 8 | High | MVP |
| US-KW-003 | 5 | High | MVP |
| US-KW-004 | 5 | High | MVP |
| US-KW-005 | 5 | High | MVP |
| US-KW-006 | 8 | High | MVP |
| **Total** | **36** | - | - |

---

## 🔗 Dependencies Graph

```
US-KW-001 (Queue List)
    │
    └── US-KW-002 (Start Preparation)
            │
            ├── US-KW-003 (Confirm Plat) ──┐
            │                              │
            ├── US-KW-004 (Input Kertas) ──┼── US-KW-006 (Finalize)
            │                              │
            └── US-KW-005 (Confirm Tinta) ─┘
```

---

## 🧪 Testing Checklist

### Unit Tests
- [ ] MaterialPreparationService - start preparation
- [ ] MaterialPreparationService - validate plat code
- [ ] MaterialPreparationService - calculate variance
- [ ] MaterialPreparationService - finalize preparation
- [ ] Priority calculation logic

### Integration Tests
- [ ] API endpoint: GET queue list
- [ ] API endpoint: POST start preparation
- [ ] API endpoint: PATCH confirm plat
- [ ] API endpoint: PATCH input kertas
- [ ] API endpoint: PATCH confirm tinta
- [ ] API endpoint: POST finalize
- [ ] SAP API integration (mock)

### E2E Tests
- [ ] Complete flow: start → confirm plat → input kertas → confirm tinta → finalize
- [ ] Validation error scenarios
- [ ] Notification delivery to Cetak

---

## 📱 UI/UX Notes

- **Mobile-first:** Semua form harus mudah digunakan di tablet
- **Large touch targets:** Button minimal 44×44 px
- **Clear feedback:** Loading state, success/error toast
- **Offline capable:** Consider PWA for basic view

---

**Last Updated:** 27 December 2025  
**Status:** Ready for Development
