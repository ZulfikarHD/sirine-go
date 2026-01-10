# Sprint Planning - Epic 03: Pemotongan (Cutting)

**Epic ID:** KHAZWAL-EPIC-03  
**Total Story Points:** 26  
**Duration:** 3 Sprints (3 Minggu)  
**Team Setup:** Solo Fullstack Developer  
**Sprint Capacity:** 8-12 Story Points per Sprint  

---

## Overview

Epic 03 mencakup proses pemotongan lembar besar menjadi lembar kirim (sisiran kiri dan kanan) serta tracking waste yang terjadi. Epic ini berisi 4 user stories dengan total 26 story points yang akan dieksekusi dalam 3 sprint dengan dependencies yang sequential.

---

## Sprint Breakdown

### Sprint 1: Cutting Queue & Start Process (Week 1)

**Story Points:** 10  
**Goal:** Implement viewing cutting queue and starting the cutting process  
**Status:** 📋 Planned

#### User Stories
- **US-KW-011:** Melihat Daftar PO yang Perlu Dipotong (5 points)
- **US-KW-012:** Memulai Proses Pemotongan (5 points)

#### Key Deliverables
1. Backend API endpoint untuk cutting queue dengan filtering/sorting
2. Frontend page `CuttingQueuePage.vue`
3. Backend API untuk memulai proses cutting
4. Frontend page `CuttingStartPage.vue`
5. Database operations pada table `khazwal_cutting_results`

#### Task Sequence (Solo Developer)

**Day 1-2: US-KW-011 Backend (9 jam)**
- KW-011-BE-01: Create API endpoint `GET /api/khazwal/cutting/queue` (3h)
- KW-011-BE-02: Query PO with status "Siap Potong" (2h)
- KW-011-BE-03: Calculate estimated output (input × 2) (1h)
- KW-011-BE-04: Implement sorting (priority + FIFO) (1h)
- KW-011-BE-05: Implement filter (tanggal, prioritas) (2h)

**Day 3: US-KW-011 Frontend (9 jam)**
- KW-011-FE-01: Create page `CuttingQueuePage.vue` (3h)
- KW-011-FE-02: Create component `CuttingQueueCard.vue` (2h)
- KW-011-FE-03: Show estimated output calculation (1h)
- KW-011-FE-04: Implement filter UI (2h)
- KW-011-FE-05: Make responsive for mobile (1h)

**Day 4: US-KW-012 Backend (11 jam)**
- KW-012-BE-01: Create API `POST /api/khazwal/cutting/{po_id}/start` (3h)
- KW-012-BE-02: Create `CuttingService.php` (2h)
- KW-012-BE-03: Fetch counting result (good quantity) (1h)
- KW-012-BE-04: Insert record to `khazwal_cutting_results` (2h)
- KW-012-BE-05: Update `production_orders.current_status` = 'SEDANG_DIPOTONG' (1h)
- KW-012-BE-06: Update `po_stage_tracking` with started_at (1h)
- KW-012-BE-07: Log to `activity_logs` (1h)

**Day 5: US-KW-012 Frontend + Testing (9 jam)**
- KW-012-FE-01: Create page `CuttingStartPage.vue` (3h)
- KW-012-FE-02: Display input & estimated output info (2h)
- KW-012-FE-03: Add cutting machine dropdown (2h)
- KW-012-FE-04: Auto-fill operator from login user (1h)
- KW-012-FE-05: Add "Mulai Pemotongan" button (1h)
- Integration testing Sprint 1

#### Acceptance Criteria
- [ ] Tampil daftar PO dengan status "Siap Potong"
- [ ] Sorting berdasarkan prioritas dan FIFO berfungsi
- [ ] Filter berdasarkan tanggal dan prioritas berfungsi
- [ ] Estimasi hasil (input × 2) dihitung otomatis
- [ ] Button "Mulai Pemotongan" berfungsi
- [ ] Status PO berubah menjadi "Sedang Dipotong"
- [ ] Timestamp mulai pemotongan tercatat
- [ ] Nama staff & mesin tercatat

---

### Sprint 2: Input Cutting Results (Week 2)

**Story Points:** 8  
**Goal:** Implement result input with waste tracking and validation  
**Status:** 📋 Planned

#### User Stories
- **US-KW-013:** Input Hasil Pemotongan (8 points)

#### Key Deliverables
1. Backend API untuk input hasil cutting dengan waste calculation
2. Photo upload handling untuk bukti waste
3. Frontend form dengan real-time waste calculation
4. Waste threshold validation (2% rule)

#### Task Sequence (Solo Developer)

**Day 1-2: Backend (12 jam)**
- KW-013-BE-01: Create API `PATCH /api/khazwal/cutting/{id}/result` (3h)
- KW-013-BE-02: Calculate total, waste, waste percentage (2h)
- KW-013-BE-03: Validate waste threshold (≤ 2%) (2h)
- KW-013-BE-04: Require reason & photo if waste > 2% (2h)
- KW-013-BE-05: Handle photo upload for waste evidence (2h)
- KW-013-BE-06: Log to `activity_logs` (1h)

**Day 3-4: Frontend (11 jam)**
- KW-013-FE-01: Create cutting result input form (3h)
- KW-013-FE-02: Input field sisiran kiri & kanan (2h)
- KW-013-FE-03: Real-time waste calculation dengan visual indicators (2h)
- KW-013-FE-04: Conditional waste reason & photo upload (3h)
- KW-013-FE-05: Input validation (tidak boleh negatif) (1h)

**Day 5: Testing & Polish (8 jam)**
- Test waste calculation accuracy (2h)
- Test photo upload flow (2h)
- Test validation rules (1h)
- Mobile UX optimization (2h)
- Edge case testing (1h)

#### Acceptance Criteria
- [ ] Input field sisiran kiri dan kanan berfungsi
- [ ] Total hasil dihitung otomatis (kiri + kanan)
- [ ] Waste dihitung otomatis dan akurat
- [ ] Waste % dihitung dengan benar
- [ ] Jika waste > 2%, wajib isi alasan dan foto
- [ ] Photo upload berfungsi dengan baik
- [ ] Real-time validation mencegah input negatif
- [ ] Visual indicator untuk waste (hijau ≤ 2%, merah > 2%)

#### Business Rules
```
Idealnya: Sisiran Kiri = Sisiran Kanan = Input (lembar besar)
Estimasi Output = Input (lembar besar) × 2
Total Output = Sisiran Kiri + Sisiran Kanan
Waste = Estimasi Output - Total Output
Waste % = (Waste ÷ Estimasi Output) × 100%
Toleransi waste: ≤ 2%
Jika waste > 2%, wajib isi alasan dan foto bukti
```

---

### Sprint 3: Finalize Cutting (Week 3)

**Story Points:** 8  
**Goal:** Complete finalization flow with verification label generation  
**Status:** 📋 Planned

#### User Stories
- **US-KW-014:** Finalisasi Pemotongan (8 points)

#### Key Deliverables
1. Backend finalization logic dengan complex database operations
2. Verification label generation (1 label per 500 lembar)
3. Stage transition ke VERIFIKASI
4. Notification system ke Tim Verifikasi
5. Frontend finalization confirmation page

#### Task Sequence (Solo Developer)

**Day 1-2: Backend - Core Logic (10 jam)**
- KW-014-BE-01: Create API `POST /api/khazwal/cutting/{id}/finalize` (3h)
- KW-014-BE-02: Validate all required fields (2h)
- KW-014-BE-03: Validate waste reason if waste > 2% (1h)
- KW-014-BE-04: Update status to 'COMPLETED' (1h)
- KW-014-BE-05: Calculate duration_minutes (1h)
- KW-014-BE-06: Update `production_orders` to next stage (1h)
- KW-014-BE-07: Log to `activity_logs` (1h)

**Day 3: Backend - Label Generation (7 jam)**
- KW-014-BE-07: Generate verification labels logic (3h)
  - Total labels = CEIL(total_output / 500)
  - Per label quantity = 500 (atau sisa untuk label terakhir)
  - Sisiran alternating: KIRI, KANAN, KIRI, KANAN
- KW-014-BE-08: Update `po_stage_tracking` with completed_at (1h)
- KW-014-BE-09: Create notification to Tim Verifikasi (2h)
- KW-014-BE-10: Additional activity logs (1h)

**Day 4: Frontend (7 jam)**
- KW-014-FE-01: Create finalization confirmation page (3h)
- KW-014-FE-02: Show cutting summary dengan semua metrics (2h)
- KW-014-FE-03: Validate before enable finalize button (1h)
- KW-014-FE-04: Show success with next step info (1h)

**Day 5: Testing & Documentation (8 jam)**
- End-to-end testing: Queue → Start → Input → Finalize (3h)
- Test label generation accuracy dengan berbagai skenario (2h)
- Test notification delivery ke Tim Verifikasi (1h)
- Code cleanup & comments (Indonesian style per cursor rules) (1h)
- Update documentation (1h)

#### Acceptance Criteria
- [ ] Button "Selesai - Kirim ke Verifikasi" berfungsi
- [ ] Validasi sebelum finalisasi berfungsi dengan baik
- [ ] Jika waste > 2%, alasan dan foto harus sudah diisi
- [ ] Status PO berubah menjadi "Siap Verifikasi"
- [ ] Verification labels generated dengan benar
- [ ] Total labels = CEIL(total_output / 500)
- [ ] Sisiran alternating (KIRI, KANAN) sesuai aturan
- [ ] Notifikasi ke Tim Verifikasi terkirim
- [ ] Timestamp selesai pemotongan tercatat
- [ ] Durasi pemotongan terhitung otomatis
- [ ] Data tidak bisa diubah setelah finalisasi

#### Business Rules
```
Durasi Pemotongan = Timestamp Selesai - Timestamp Mulai
Total Labels = CEIL(total_output / 500)
Per Label Quantity = 500 (atau sisa untuk label terakhir)
Sisiran alternating: Label 1 = KIRI, Label 2 = KANAN, Label 3 = KIRI, dst.
Data tidak bisa diubah setelah finalisasi
```

---

## Technical Architecture

### Data Flow
```
Counting Result (lembar besar) 
    ↓
Cutting Queue (status: SIAP_POTONG)
    ↓
Start Cutting (status: SEDANG_DIPOTONG)
    ↓
Input Results (sisiran kiri + kanan)
    ↓
Finalize (status: COMPLETED)
    ↓
Generate Verification Labels (per 500 lembar)
    ↓
Notify Tim Verifikasi
    ↓
Stage: VERIFIKASI
```

### Database Tables Involved

#### Primary Table
- `khazwal_cutting_results`
  - Columns: id, production_order_id, input_lembar_besar, expected_output, 
    output_sisiran_kiri, output_sisiran_kanan, total_output, waste_quantity, 
    waste_percentage, waste_reason, waste_photo_url, cutting_machine, cut_by, 
    status, started_at, completed_at, duration_minutes

#### Related Tables
- `production_orders` - Status updates (KHAZWAL_CUTTING → VERIFIKASI)
- `khazwal_counting_results` - Input data source (quantity_good)
- `po_stage_tracking` - Progress tracking (started_at, completed_at)
- `verification_labels` - Generated in finalization
- `notifications` - Alert Tim Verifikasi
- `activity_logs` - Audit trail
- `users` - Staff information

### Critical Business Logic

```php
// Estimation
$estimatedOutput = $inputLembarBesar * 2;

// Waste Calculation
$totalOutput = $outputSisiranKiri + $outputSisiranKanan;
$waste = $estimatedOutput - $totalOutput;
$wastePercentage = round(($waste / $estimatedOutput) * 100, 2);

// Waste Threshold
if ($wastePercentage > 2) {
    // Require reason & photo
    $isWasteReasonRequired = true;
    $isWastePhotoRequired = true;
}

// Label Generation
$totalLabels = ceil($totalOutput / 500);
for ($i = 1; $i <= $totalLabels; $i++) {
    $targetQuantity = ($i < $totalLabels) ? 500 : ($totalOutput % 500);
    $sisiran = ($i % 2 == 1) ? 'KIRI' : 'KANAN';
    // Insert to verification_labels
}

// Duration Calculation
$durationMinutes = floor((strtotime($completedAt) - strtotime($startedAt)) / 60);
```

---

## Risk Mitigation

### Technical Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Photo upload size terlalu besar | High | Medium | Implement compression sebelum upload (max 2MB) |
| Label generation complexity | High | Low | Test dengan berbagai output quantities |
| Transaction integrity | Critical | Low | Use DB transactions untuk finalization |
| Waste calculation precision | High | Low | Use proper rounding, test edge cases |
| Database deadlock | Medium | Low | Optimize query order, use proper locking |

### Timeline Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Solo developer bottleneck | High | Prioritize backend completion before frontend |
| Testing time insufficient | Medium | Allocate full Day 5 each sprint for testing |
| Dependency on Epic 02 | Critical | Ensure US-KW-010 (Counting) is completed first |
| Photo upload implementation delay | Medium | Use existing storage service, test early |
| Label generation bugs | High | Create unit tests early, test multiple scenarios |

---

## Success Criteria

### Sprint 1 Success Criteria
- [ ] Can view list of PO ready for cutting with correct data
- [ ] Can filter PO by date and priority
- [ ] Can sort PO by priority and FIFO
- [ ] Estimated output calculated correctly (input × 2)
- [ ] Can start cutting process successfully
- [ ] Machine and operator assignment works
- [ ] Status updates reflect in database correctly
- [ ] Mobile responsive on all screens

### Sprint 2 Success Criteria
- [ ] Can input sisiran kiri & kanan quantities
- [ ] Total output calculated automatically
- [ ] Waste calculated automatically and accurately
- [ ] Waste percentage displayed with color indicator
- [ ] Photo upload works when waste > 2%
- [ ] Validations prevent invalid data entry
- [ ] Real-time calculation works smoothly
- [ ] Mobile-friendly input forms

### Sprint 3 Success Criteria
- [ ] Can finalize cutting process successfully
- [ ] All validations work before finalization
- [ ] Verification labels generated correctly
- [ ] Label count = CEIL(total_output / 500)
- [ ] Sisiran alternating pattern correct
- [ ] PO transitions to VERIFIKASI stage
- [ ] Notifications sent to Tim Verifikasi
- [ ] Complete audit trail in activity_logs
- [ ] Data immutable after finalization
- [ ] Duration calculated correctly

---

## Definition of Done (Per Sprint)

### Code Quality
- [ ] All tasks completed and functionally tested
- [ ] Code follows cursor rules (Indonesian comments dengan format yang proper)
- [ ] Code follows Laravel best practices
- [ ] Code follows Vue 3 Composition API best practices
- [ ] No console.log atau debug code yang tertinggal

### Testing
- [ ] Unit tests untuk business logic (jika applicable)
- [ ] Integration tests untuk API endpoints
- [ ] Manual testing untuk happy path
- [ ] Manual testing untuk edge cases
- [ ] Error handling tested

### UI/UX
- [ ] Mobile responsive (per design standard)
- [ ] Motion-V animations implemented
- [ ] Loading states implemented
- [ ] Error states implemented
- [ ] Success feedback implemented
- [ ] Apple-inspired design dengan Indigo & Fuchsia gradient

### Technical
- [ ] Linter errors resolved (yarn run lint)
- [ ] No breaking changes to existing features
- [ ] Integration tested with existing Counting module
- [ ] Database migrations created and tested
- [ ] API documentation updated (if applicable)

### Performance
- [ ] Page load < 2 detik
- [ ] No N+1 queries
- [ ] Efficient database queries
- [ ] Proper indexing on database tables

---

## Testing Strategy

### Unit Tests
- CuttingService - start cutting logic
- CuttingService - calculate waste percentage
- CuttingService - validate waste threshold
- CuttingService - finalize cutting logic
- Verification label generation logic

### Integration Tests
- API endpoint: GET cutting queue
- API endpoint: POST start cutting
- API endpoint: PATCH input result
- API endpoint: POST finalize
- Notification delivery to Verifikasi

### E2E Test Scenarios
1. **Happy Path:** View queue → Start → Input (waste ≤ 2%) → Finalize
2. **High Waste Path:** View queue → Start → Input (waste > 2% + reason + photo) → Finalize
3. **Edge Case:** Input dengan sisiran tidak seimbang
4. **Edge Case:** Output yang menghasilkan sisa label (bukan kelipatan 500)
5. **Error Case:** Finalisasi tanpa input hasil
6. **Error Case:** Finalisasi dengan waste > 2% tanpa alasan

---

## Dependencies

### External Dependencies
- Epic 02 (Counting) harus completed (US-KW-010: Finalisasi Penghitungan)
- Database tables sudah di-migrate
- User authentication & authorization system sudah ready
- File upload service sudah configured

### Internal Dependencies
- Sprint 1 harus completed sebelum Sprint 2
- Sprint 2 harus completed sebelum Sprint 3
- Sequential execution karena dependencies

---

## Sprint Retrospective Questions

### Sprint 1 Retro
- Apakah estimasi waktu untuk API development akurat?
- Apakah task sequence efisien untuk solo developer?
- Bottleneck apa yang ditemukan?
- Improvement apa yang bisa dilakukan di Sprint 2?

### Sprint 2 Retro
- Apakah waste calculation logic mudah dipahami?
- Apakah photo upload implementation smooth?
- Apakah testing time cukup?
- Improvement apa yang bisa dilakukan di Sprint 3?

### Sprint 3 Retro
- Apakah label generation logic robust?
- Apakah notification system berfungsi dengan baik?
- Apakah documentation memadai?
- Lessons learned untuk Epic berikutnya?

---

## Next Steps After Epic Completion

Setelah Epic 03 completed:
1. Deploy ke staging environment
2. User Acceptance Testing (UAT) dengan stakeholder
3. Fix bugs dari UAT
4. Deploy ke production
5. Monitor production usage
6. Gather feedback dari users (Staff Khazanah Awal)
7. Mulai Epic berikutnya (Epic 04: Verifikasi)

---

**Created:** {{ date }}  
**Last Updated:** {{ date }}  
**Status:** Ready for Execution  
**Owner:** Zulfikar Hidayatullah  
