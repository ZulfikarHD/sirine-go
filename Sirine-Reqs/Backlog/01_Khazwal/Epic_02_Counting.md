# Epic 2: Penghitungan Hasil Cetak (Counting)

**Epic ID:** KHAZWAL-EPIC-02  
**Priority:** 🔴 High (Phase 1 - MVP)  
**Estimated Duration:** 2 Minggu  

---

## 📋 Overview

Epic ini mencakup proses penghitungan hasil cetak yang kembali dari Unit Cetak ke Khazanah Awal, termasuk kategorisasi hasil baik dan rusak.

---

## 🗄️ Database Reference

### Primary Table
- `khazwal_counting_results` - Menyimpan data hasil penghitungan

### Related Tables
- `production_orders` - Referensi PO utama
- `po_stage_tracking` - Tracking perjalanan stage PO
- `print_job_summaries` - Referensi hasil cetak dari Unit Cetak
- `users` - Staff yang menangani
- `activity_logs` - Audit trail
- `alerts` - Alert jika waiting time > 2 jam

---

## 📝 Backlog Items

### US-KW-007: Melihat Daftar PO yang Perlu Dihitung

| Field | Value |
|-------|-------|
| **ID** | US-KW-007 |
| **Story Points** | 5 |
| **Priority** | 🔴 High |
| **Dependencies** | Print Job Completion (Cetak) |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin melihat daftar PO yang hasil cetaknya sudah kembali dari Unit Cetak, sehingga saya tahu PO mana yang perlu dihitung.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-007-BE-01 | Create API endpoint `GET /api/khazwal/counting/queue` | 3h | Backend |
| KW-007-BE-02 | Query PO with status "Selesai Cetak - Menunggu Penghitungan" | 2h | Backend |
| KW-007-BE-03 | Implement FIFO sorting (waktu selesai cetak) | 1h | Backend |
| KW-007-BE-04 | Calculate waiting duration per PO | 1h | Backend |
| KW-007-BE-05 | Create alert if waiting > 2 hours | 2h | Backend |
| KW-007-BE-06 | Implement filter (tanggal, mesin) | 2h | Backend |
| KW-007-FE-01 | Create page `CountingQueuePage.vue` | 3h | Frontend |
| KW-007-FE-02 | Create component `CountingQueueCard.vue` | 2h | Frontend |
| KW-007-FE-03 | Show waiting duration with alert indicator | 2h | Frontend |
| KW-007-FE-04 | Implement filter UI | 2h | Frontend |
| KW-007-FE-05 | Make responsive for mobile | 1h | Frontend |

#### Acceptance Criteria
- [ ] Tampil daftar PO dengan status "Selesai Cetak - Menunggu Penghitungan"
- [ ] Sorting berdasarkan waktu selesai cetak (FIFO)
- [ ] Tampil info per PO:
  - Nomor PO & OBC
  - Target cetak (lembar besar)
  - Waktu selesai cetak
  - Durasi menunggu penghitungan
  - Mesin cetak yang digunakan
- [ ] Filter berdasarkan tanggal, mesin
- [ ] Alert jika ada PO yang menunggu > 2 jam

#### Business Rules
```
- PO yang selesai cetak lebih dulu harus dihitung lebih dulu (FIFO)
- Alert jika waiting time > 2 jam
```

#### Database Query Reference
```sql
SELECT 
    po.id,
    po.po_number,
    po.obc_number,
    po.quantity_target_lembar_besar as target,
    pjs.finalized_at as print_completed_at,
    m.machine_name,
    u.full_name as operator_name,
    EXTRACT(EPOCH FROM (NOW() - pjs.finalized_at)) / 60 as waiting_minutes,
    CASE 
        WHEN EXTRACT(EPOCH FROM (NOW() - pjs.finalized_at)) / 60 > 120 
        THEN true 
        ELSE false 
    END as is_overdue
FROM production_orders po
JOIN print_job_summaries pjs ON pjs.production_order_id = po.id
JOIN machines m ON m.id = pjs.machine_id
JOIN users u ON u.id = pjs.operator_id
WHERE po.current_stage = 'KHAZWAL_COUNTING'
    AND po.current_status = 'WAITING_COUNTING'
    AND po.deleted_at IS NULL
ORDER BY pjs.finalized_at ASC; -- FIFO
```

---

### US-KW-008: Memulai Proses Penghitungan

| Field | Value |
|-------|-------|
| **ID** | US-KW-008 |
| **Story Points** | 5 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-007 |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin memulai proses penghitungan hasil cetak, sehingga sistem tracking progress penghitungan.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-008-BE-01 | Create API `POST /api/khazwal/counting/{po_id}/start` | 3h | Backend |
| KW-008-BE-02 | Create `CountingService.php` | 2h | Backend |
| KW-008-BE-03 | Fetch print job summary data | 1h | Backend |
| KW-008-BE-04 | Insert record to `khazwal_counting_results` | 2h | Backend |
| KW-008-BE-05 | Update `production_orders.current_status` = 'SEDANG_DIHITUNG' | 1h | Backend |
| KW-008-BE-06 | Update `po_stage_tracking` with started_at | 1h | Backend |
| KW-008-BE-07 | Log to `activity_logs` | 1h | Backend |
| KW-008-FE-01 | Create page `CountingStartPage.vue` | 3h | Frontend |
| KW-008-FE-02 | Display print job info (target, mesin, operator, waktu) | 2h | Frontend |
| KW-008-FE-03 | Add "Mulai Penghitungan" button | 1h | Frontend |
| KW-008-FE-04 | Show loading state and redirect | 1h | Frontend |

#### Acceptance Criteria
- [ ] Button "Mulai Penghitungan" pada detail PO
- [ ] Sistem tampilkan:
  - Target: [X] lembar besar
  - Mesin cetak: [Nama Mesin]
  - Operator cetak: [Nama Operator]
  - Waktu cetak: [Tanggal & Jam]
- [ ] Status PO berubah menjadi "Sedang Dihitung"
- [ ] Timestamp mulai penghitungan tercatat
- [ ] Nama staff yang handle tercatat

#### Database Insert Reference
```sql
INSERT INTO khazwal_counting_results (
    production_order_id,
    status,
    started_at,
    counted_by,
    created_at
) VALUES (
    :po_id,
    'IN_PROGRESS',
    NOW(),
    :user_id,
    NOW()
);
```

---

### US-KW-009: Input Hasil Penghitungan

| Field | Value |
|-------|-------|
| **ID** | US-KW-009 |
| **Story Points** | 8 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-008 |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin input hasil penghitungan lembar besar, sehingga ada record akurat berapa yang baik dan berapa yang rusak.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-009-BE-01 | Create API `PATCH /api/khazwal/counting/{id}/result` | 3h | Backend |
| KW-009-BE-02 | Calculate total, variance, percentages | 2h | Backend |
| KW-009-BE-03 | Validate total vs target (±2% tolerance) | 2h | Backend |
| KW-009-BE-04 | Require defect breakdown if rusak > 5% | 2h | Backend |
| KW-009-BE-05 | Store defect_breakdown as JSONB | 1h | Backend |
| KW-009-BE-06 | Require variance reason if variance ≠ 0 | 1h | Backend |
| KW-009-BE-07 | Log to `activity_logs` | 1h | Backend |
| KW-009-FE-01 | Create counting input form | 3h | Frontend |
| KW-009-FE-02 | Real-time calculation display | 2h | Frontend |
| KW-009-FE-03 | Conditional defect breakdown form | 3h | Frontend |
| KW-009-FE-04 | Input validation (tidak boleh negatif) | 1h | Frontend |
| KW-009-FE-05 | Variance reason field (conditional) | 1h | Frontend |

#### Acceptance Criteria
- [ ] Input field:
  - Jumlah Baik (lembar besar)
  - Jumlah Rusak (lembar besar)
- [ ] Sistem otomatis hitung:
  - Total = Baik + Rusak
  - Selisih = Total - Target
  - Persentase Baik = (Baik ÷ Target) × 100%
  - Persentase Rusak = (Rusak ÷ Target) × 100%
- [ ] Jika selisih ≠ 0, tampil field "Keterangan Selisih"
- [ ] Jika rusak > 5%, wajib breakdown jenis kerusakan
- [ ] Real-time validation input (tidak boleh negatif)

#### Business Rules
```
- Total (Baik + Rusak) harus = Target ± toleransi
- Toleransi selisih: ±2%
- Jika rusak > 5%, wajib breakdown jenis kerusakan
```

#### JSONB Structure for defect_breakdown
```json
[
  {"type": "Warna pudar", "quantity": 15},
  {"type": "Tinta blobor", "quantity": 8},
  {"type": "Kertas sobek", "quantity": 3}
]
```

#### Database Update Reference
```sql
UPDATE khazwal_counting_results
SET 
    quantity_good = :good,
    quantity_defect = :defect,
    total_counted = :good + :defect,
    variance_from_target = (:good + :defect) - (
        SELECT quantity_target_lembar_besar 
        FROM production_orders 
        WHERE id = production_order_id
    ),
    percentage_good = ROUND((:good::numeric / NULLIF(:good + :defect, 0)) * 100, 2),
    percentage_defect = ROUND((:defect::numeric / NULLIF(:good + :defect, 0)) * 100, 2),
    defect_breakdown = :defect_json,
    variance_reason = :reason,
    updated_at = NOW()
WHERE id = :id;
```

---

### US-KW-010: Finalisasi Penghitungan

| Field | Value |
|-------|-------|
| **ID** | US-KW-010 |
| **Story Points** | 5 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-009 |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin finalisasi hasil penghitungan, sehingga proses bisa lanjut ke pemotongan.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-010-BE-01 | Create API `POST /api/khazwal/counting/{id}/finalize` | 3h | Backend |
| KW-010-BE-02 | Validate all required fields filled | 2h | Backend |
| KW-010-BE-03 | Validate defect breakdown if rusak > 5% | 1h | Backend |
| KW-010-BE-04 | Update status to 'COMPLETED' | 1h | Backend |
| KW-010-BE-05 | Calculate duration_minutes | 1h | Backend |
| KW-010-BE-06 | Update `production_orders.current_status` = 'SIAP_POTONG' | 1h | Backend |
| KW-010-BE-07 | Update `po_stage_tracking` with completed_at | 1h | Backend |
| KW-010-BE-08 | Log to `activity_logs` (immutable after finalize) | 1h | Backend |
| KW-010-FE-01 | Create finalization confirmation modal | 2h | Frontend |
| KW-010-FE-02 | Show summary before finalize | 2h | Frontend |
| KW-010-FE-03 | Disable edit after finalize | 1h | Frontend |
| KW-010-FE-04 | Redirect to cutting queue | 1h | Frontend |

#### Acceptance Criteria
- [ ] Button "Selesai Penghitungan - Lanjut Pemotongan"
- [ ] Sistem validasi:
  - ✅ Jumlah baik & rusak sudah diinput
  - ✅ Jika rusak > 5%, breakdown sudah diisi
- [ ] Status PO berubah menjadi "Siap Potong"
- [ ] Timestamp selesai penghitungan tercatat
- [ ] Sistem hitung durasi penghitungan
- [ ] Data tersimpan dan tidak bisa diubah (audit trail)

#### Business Rules
```
- Semua validasi harus pass sebelum bisa finalisasi
- Data penghitungan tidak bisa diubah setelah finalisasi (untuk audit)
```

#### Database Update Reference
```sql
-- Update counting result
UPDATE khazwal_counting_results
SET 
    status = 'COMPLETED',
    completed_at = NOW(),
    duration_minutes = EXTRACT(EPOCH FROM (NOW() - started_at)) / 60,
    updated_at = NOW()
WHERE id = :id;

-- Update production order
UPDATE production_orders
SET 
    current_stage = 'KHAZWAL_CUTTING',
    current_status = 'SIAP_POTONG',
    updated_at = NOW()
WHERE id = :po_id;
```

---

## 📊 Epic Summary

| User Story | Story Points | Priority | Phase |
|------------|--------------|----------|-------|
| US-KW-007 | 5 | High | MVP |
| US-KW-008 | 5 | High | MVP |
| US-KW-009 | 8 | High | MVP |
| US-KW-010 | 5 | High | MVP |
| **Total** | **23** | - | - |

---

## 🔗 Dependencies Graph

```
[Print Job Completion] (from Cetak)
    │
    └── US-KW-007 (Counting Queue)
            │
            └── US-KW-008 (Start Counting)
                    │
                    └── US-KW-009 (Input Results)
                            │
                            └── US-KW-010 (Finalize)
                                    │
                                    └── [Cutting Queue] (Epic 3)
```

---

## 🧪 Testing Checklist

### Unit Tests
- [ ] CountingService - start counting
- [ ] CountingService - calculate percentages
- [ ] CountingService - validate tolerance
- [ ] CountingService - finalize counting
- [ ] Defect breakdown validation

### Integration Tests
- [ ] API endpoint: GET counting queue
- [ ] API endpoint: POST start counting
- [ ] API endpoint: PATCH input result
- [ ] API endpoint: POST finalize
- [ ] Alert creation for overdue items

### E2E Tests
- [ ] Complete flow: view queue → start → input → finalize
- [ ] Defect breakdown flow (when rusak > 5%)
- [ ] Variance reason flow (when selisih ≠ 0)
- [ ] Immutability after finalize

---

## 📱 UI/UX Notes

- **Large input fields:** Untuk kemudahan input angka di tablet
- **Real-time calculation:** Tampilkan hasil hitung segera setelah input
- **Clear visual feedback:** Warna merah untuk rusak, hijau untuk baik
- **Confirmation dialog:** Sebelum finalisasi untuk mencegah kesalahan

---

**Last Updated:** 27 December 2025  
**Status:** Ready for Development
