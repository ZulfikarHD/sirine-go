# Epic 3: Pemotongan (Cutting)

**Epic ID:** KHAZWAL-EPIC-03  
**Priority:** 🔴 High (Phase 1 - MVP)  
**Estimated Duration:** 2 Minggu  

---

## 📋 Overview

Epic ini mencakup proses pemotongan lembar besar menjadi lembar kirim (sisiran kiri dan kanan) serta tracking waste yang terjadi.

---

## 🗄️ Database Reference

### Primary Table
- `khazwal_cutting_results` - Menyimpan data hasil pemotongan

### Related Tables
- `production_orders` - Referensi PO utama
- `khazwal_counting_results` - Input dari penghitungan (good quantity)
- `po_stage_tracking` - Tracking perjalanan stage PO
- `users` - Staff yang menangani
- `activity_logs` - Audit trail
- `notifications` - Notifikasi ke Tim Verifikasi

---

## 📝 Backlog Items

### US-KW-011: Melihat Daftar PO yang Perlu Dipotong

| Field | Value |
|-------|-------|
| **ID** | US-KW-011 |
| **Story Points** | 5 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-010 (Finalisasi Penghitungan) |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin melihat daftar PO yang siap untuk dipotong, sehingga saya tahu prioritas pemotongan.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-011-BE-01 | Create API endpoint `GET /api/khazwal/cutting/queue` | 3h | Backend |
| KW-011-BE-02 | Query PO with status "Siap Potong" | 2h | Backend |
| KW-011-BE-03 | Calculate estimated output (input × 2) | 1h | Backend |
| KW-011-BE-04 | Implement sorting (priority + FIFO) | 1h | Backend |
| KW-011-BE-05 | Implement filter (tanggal, prioritas) | 2h | Backend |
| KW-011-FE-01 | Create page `CuttingQueuePage.vue` | 3h | Frontend |
| KW-011-FE-02 | Create component `CuttingQueueCard.vue` | 2h | Frontend |
| KW-011-FE-03 | Show estimated output calculation | 1h | Frontend |
| KW-011-FE-04 | Implement filter UI | 2h | Frontend |
| KW-011-FE-05 | Make responsive for mobile | 1h | Frontend |

#### Acceptance Criteria
- [ ] Tampil daftar PO dengan status "Siap Potong"
- [ ] Sorting berdasarkan prioritas dan FIFO
- [ ] Tampil info per PO:
  - Nomor PO & OBC
  - Jumlah lembar besar yang akan dipotong
  - Estimasi hasil: [X] lembar kirim (otomatis × 2)
  - Prioritas
- [ ] Filter berdasarkan tanggal, prioritas

#### Business Rules
```
Estimasi Hasil Pemotongan = Jumlah Baik (lembar besar) × 2
```

#### Database Query Reference
```sql
SELECT 
    po.id,
    po.po_number,
    po.obc_number,
    po.priority,
    cr.quantity_good as input_lembar_besar,
    cr.quantity_good * 2 as estimated_output,
    cr.completed_at as counting_completed_at
FROM production_orders po
JOIN khazwal_counting_results cr ON cr.production_order_id = po.id
WHERE po.current_stage = 'KHAZWAL_CUTTING'
    AND po.current_status = 'SIAP_POTONG'
    AND po.deleted_at IS NULL
ORDER BY 
    CASE po.priority 
        WHEN 'URGENT' THEN 1 
        WHEN 'NORMAL' THEN 2 
        ELSE 3 
    END,
    cr.completed_at ASC;
```

---

### US-KW-012: Memulai Proses Pemotongan

| Field | Value |
|-------|-------|
| **ID** | US-KW-012 |
| **Story Points** | 5 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-011 |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin memulai proses pemotongan, sehingga sistem tracking progress pemotongan.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-012-BE-01 | Create API `POST /api/khazwal/cutting/{po_id}/start` | 3h | Backend |
| KW-012-BE-02 | Create `CuttingService.php` | 2h | Backend |
| KW-012-BE-03 | Fetch counting result (good quantity) | 1h | Backend |
| KW-012-BE-04 | Insert record to `khazwal_cutting_results` | 2h | Backend |
| KW-012-BE-05 | Update `production_orders.current_status` = 'SEDANG_DIPOTONG' | 1h | Backend |
| KW-012-BE-06 | Update `po_stage_tracking` with started_at | 1h | Backend |
| KW-012-BE-07 | Log to `activity_logs` | 1h | Backend |
| KW-012-FE-01 | Create page `CuttingStartPage.vue` | 3h | Frontend |
| KW-012-FE-02 | Display input & estimated output info | 2h | Frontend |
| KW-012-FE-03 | Add cutting machine dropdown | 2h | Frontend |
| KW-012-FE-04 | Auto-fill operator from login user | 1h | Frontend |
| KW-012-FE-05 | Add "Mulai Pemotongan" button | 1h | Frontend |

#### Acceptance Criteria
- [ ] Button "Mulai Pemotongan" pada detail PO
- [ ] Sistem tampilkan:
  - Input: [X] lembar besar
  - Estimasi Output: [Y] lembar kirim (X × 2)
  - Mesin potong yang digunakan (dropdown)
  - Operator (auto-fill dari login user)
- [ ] Status PO berubah menjadi "Sedang Dipotong"
- [ ] Timestamp mulai pemotongan tercatat
- [ ] Nama staff & mesin tercatat

#### Database Insert Reference
```sql
INSERT INTO khazwal_cutting_results (
    production_order_id,
    input_lembar_besar,
    expected_output,
    cutting_machine,
    cut_by,
    status,
    started_at,
    created_at
) 
SELECT 
    :po_id,
    cr.quantity_good,
    cr.quantity_good * 2,
    :machine,
    :user_id,
    'IN_PROGRESS',
    NOW(),
    NOW()
FROM khazwal_counting_results cr
WHERE cr.production_order_id = :po_id;
```

---

### US-KW-013: Input Hasil Pemotongan

| Field | Value |
|-------|-------|
| **ID** | US-KW-013 |
| **Story Points** | 8 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-012 |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin input hasil pemotongan, sehingga ada record akurat konversi lembar besar → lembar kirim.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-013-BE-01 | Create API `PATCH /api/khazwal/cutting/{id}/result` | 3h | Backend |
| KW-013-BE-02 | Calculate total, waste, waste percentage | 2h | Backend |
| KW-013-BE-03 | Validate waste threshold (≤ 2%) | 2h | Backend |
| KW-013-BE-04 | Require reason & photo if waste > 2% | 2h | Backend |
| KW-013-BE-05 | Handle photo upload for waste evidence | 2h | Backend |
| KW-013-BE-06 | Log to `activity_logs` | 1h | Backend |
| KW-013-FE-01 | Create cutting result input form | 3h | Frontend |
| KW-013-FE-02 | Input field sisiran kiri & kanan | 2h | Frontend |
| KW-013-FE-03 | Real-time waste calculation | 2h | Frontend |
| KW-013-FE-04 | Conditional waste reason & photo upload | 3h | Frontend |
| KW-013-FE-05 | Input validation (tidak boleh negatif) | 1h | Frontend |

#### Acceptance Criteria
- [ ] Input field:
  - Sisiran Kiri (lembar kirim)
  - Sisiran Kanan (lembar kirim)
- [ ] Sistem otomatis hitung:
  - Total Hasil = Sisiran Kiri + Sisiran Kanan
  - Estimasi = Input (lembar besar) × 2
  - Selisih = Total Hasil - Estimasi
  - Waste = Estimasi - Total Hasil
  - Waste % = (Waste ÷ Estimasi) × 100%
- [ ] Jika waste > 2%, wajib isi alasan
- [ ] Real-time validation (tidak boleh negatif)

#### Business Rules
```
- Idealnya: Sisiran Kiri = Sisiran Kanan = Input (lembar besar)
- Toleransi waste: ≤ 2%
- Jika waste > 2%, wajib isi alasan dan foto bukti
```

#### Database Update Reference
```sql
UPDATE khazwal_cutting_results
SET 
    output_sisiran_kiri = :kiri,
    output_sisiran_kanan = :kanan,
    total_output = :kiri + :kanan,
    waste_quantity = expected_output - (:kiri + :kanan),
    waste_percentage = ROUND(
        ((expected_output - (:kiri + :kanan))::numeric / NULLIF(expected_output, 0)) * 100, 
        2
    ),
    waste_reason = :reason,
    waste_photo_url = :photo_url,
    updated_at = NOW()
WHERE id = :id;
```

---

### US-KW-014: Finalisasi Pemotongan

| Field | Value |
|-------|-------|
| **ID** | US-KW-014 |
| **Story Points** | 8 |
| **Priority** | 🔴 High |
| **Dependencies** | US-KW-013 |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin finalisasi hasil pemotongan, sehingga hasil bisa dikirim ke Verifikasi.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-014-BE-01 | Create API `POST /api/khazwal/cutting/{id}/finalize` | 3h | Backend |
| KW-014-BE-02 | Validate all required fields | 2h | Backend |
| KW-014-BE-03 | Validate waste reason if waste > 2% | 1h | Backend |
| KW-014-BE-04 | Update status to 'COMPLETED' | 1h | Backend |
| KW-014-BE-05 | Calculate duration_minutes | 1h | Backend |
| KW-014-BE-06 | Update `production_orders` to next stage | 1h | Backend |
| KW-014-BE-07 | Generate verification labels | 3h | Backend |
| KW-014-BE-08 | Update `po_stage_tracking` with completed_at | 1h | Backend |
| KW-014-BE-09 | Create notification to Tim Verifikasi | 2h | Backend |
| KW-014-BE-10 | Log to `activity_logs` | 1h | Backend |
| KW-014-FE-01 | Create finalization confirmation page | 3h | Frontend |
| KW-014-FE-02 | Show cutting summary | 2h | Frontend |
| KW-014-FE-03 | Validate before enable finalize button | 1h | Frontend |
| KW-014-FE-04 | Show success with next step info | 1h | Frontend |

#### Acceptance Criteria
- [ ] Button "Selesai - Kirim ke Verifikasi"
- [ ] Sistem validasi:
  - ✅ Hasil pemotongan sudah diinput
  - ✅ Jika waste > 2%, alasan sudah diisi
- [ ] Status PO berubah menjadi "Siap Verifikasi"
- [ ] Notifikasi ke Tim Verifikasi
- [ ] Timestamp selesai pemotongan tercatat
- [ ] Sistem hitung durasi pemotongan
- [ ] Data tersimpan dan tidak bisa diubah

#### Business Rules
```
- Durasi Pemotongan = Timestamp Selesai - Timestamp Mulai
- Data tidak bisa diubah setelah finalisasi
```

#### Database Operations
```sql
-- Update cutting result
UPDATE khazwal_cutting_results
SET 
    status = 'COMPLETED',
    completed_at = NOW(),
    duration_minutes = EXTRACT(EPOCH FROM (NOW() - started_at)) / 60,
    updated_at = NOW()
WHERE id = :id;

-- Update production order
UPDATE production_orders
SET 
    current_stage = 'VERIFIKASI',
    current_status = 'SIAP_VERIFIKASI',
    updated_at = NOW()
WHERE id = :po_id;

-- Generate verification labels (per rim)
-- Total labels = total_output / 500 (rounded up)
INSERT INTO verification_labels (
    production_order_id,
    label_number,
    total_labels,
    target_quantity,
    sisiran,
    cutting_result_id,
    qc_status,
    created_at
)
SELECT 
    :po_id,
    generate_series(1, CEIL(cr.total_output::numeric / 500)),
    CEIL(cr.total_output::numeric / 500),
    CASE 
        WHEN generate_series <= FLOOR(cr.total_output::numeric / 500) THEN 500
        ELSE cr.total_output % 500
    END,
    CASE WHEN generate_series % 2 = 1 THEN 'KIRI' ELSE 'KANAN' END,
    cr.id,
    'PENDING',
    NOW()
FROM khazwal_cutting_results cr
WHERE cr.production_order_id = :po_id;

-- Create notification to Verifikasi
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
    'QUALITY',
    'Hasil Potong Siap untuk PO ' || :po_number,
    'PO ' || :po_number || ' sudah selesai dipotong. Total ' || :total_labels || ' label siap verifikasi.',
    'PRODUCTION_ORDER',
    :po_id,
    'NORMAL',
    NOW()
FROM users u
WHERE u.department = 'VERIFIKASI' AND u.role IN ('QC_INSPECTOR', 'SUPERVISOR_VERIFIKASI');
```

---

## 📊 Epic Summary

| User Story | Story Points | Priority | Phase |
|------------|--------------|----------|-------|
| US-KW-011 | 5 | High | MVP |
| US-KW-012 | 5 | High | MVP |
| US-KW-013 | 8 | High | MVP |
| US-KW-014 | 8 | High | MVP |
| **Total** | **26** | - | - |

---

## 🔗 Dependencies Graph

```
US-KW-010 (Finalisasi Counting)
    │
    └── US-KW-011 (Cutting Queue)
            │
            └── US-KW-012 (Start Cutting)
                    │
                    └── US-KW-013 (Input Results)
                            │
                            └── US-KW-014 (Finalize)
                                    │
                                    └── [Verification Labels Generation]
                                            │
                                            └── [Verifikasi Queue]
```

---

## 🧪 Testing Checklist

### Unit Tests
- [ ] CuttingService - start cutting
- [ ] CuttingService - calculate waste percentage
- [ ] CuttingService - validate waste threshold
- [ ] CuttingService - finalize cutting
- [ ] Verification label generation logic

### Integration Tests
- [ ] API endpoint: GET cutting queue
- [ ] API endpoint: POST start cutting
- [ ] API endpoint: PATCH input result
- [ ] API endpoint: POST finalize
- [ ] Notification delivery to Verifikasi

### E2E Tests
- [ ] Complete flow: view queue → start → input → finalize
- [ ] Waste reason flow (when waste > 2%)
- [ ] Photo upload for waste evidence
- [ ] Verification label generation

---

## 📱 UI/UX Notes

- **Clear visual:** Sisiran kiri vs kanan dengan visual representation
- **Waste indicator:** Warna merah jika waste > 2%, hijau jika ≤ 2%
- **Photo upload:** Easy drag & drop atau camera capture
- **Summary card:** Tampilkan ringkasan sebelum finalisasi

---

## 🔄 Integration Notes

### Output to Verifikasi
Setelah finalisasi, sistem akan:
1. Generate verification labels (1 label per ~500 lembar)
2. Insert ke tabel `verification_labels`
3. Update PO stage ke VERIFIKASI
4. Send notification ke Tim Verifikasi

### Label Generation Logic
```
Total Labels = CEIL(total_output / 500)
Per Label Quantity = 500 (atau sisa untuk label terakhir)
Sisiran alternating: Label 1 = KIRI, Label 2 = KANAN, dst.
```

---

**Last Updated:** 27 December 2025  
**Status:** Ready for Development
