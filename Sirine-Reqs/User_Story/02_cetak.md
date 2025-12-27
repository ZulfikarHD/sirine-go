# User Story: Cetak (Print)

## Overview
Unit Cetak bertanggung jawab untuk mencetak pita cukai sesuai spesifikasi yang tertera pada PO. Proses ini merupakan tahap kritis yang mempengaruhi OEE (Overall Equipment Effectiveness), kualitas produk, dan efisiensi produksi secara keseluruhan.

---

## 🎭 Personas

### 1. Operator Cetak
**Nama:** Andi Wijaya  
**Posisi:** Operator Mesin Cetak  
**Shift:** Siang (15:00-23:00)  
**Pengalaman:** 4 tahun

**Tanggung Jawab:**
- Setup mesin cetak (plat, tinta, kertas)
- Menjalankan proses cetak sesuai spesifikasi PO
- Monitoring kualitas hasil cetak
- Melakukan adjustment jika ada masalah
- Entry data produksi real-time
- Maintenance harian mesin

**Pain Points:**
- Setup mesin memakan waktu lama (20-30 menit)
- Sulit tracking progress produksi real-time
- Tidak tahu apakah on-target atau behind schedule
- Harus turun dari mesin untuk input data ke komputer
- Tidak ada alert jika mesin bermasalah
- Laporan manual di akhir shift memakan waktu

**Goals:**
- Setup lebih cepat dan efisien
- Real-time visibility progress vs target
- Input data mudah (via mobile)
- Alert otomatis jika ada masalah
- Fokus ke produksi, bukan administrasi

---

### 2. Supervisor Unit Cetak
**Nama:** Hendra Gunawan  
**Posisi:** Supervisor Unit Cetak  
**Shift:** Rotasi  
**Pengalaman:** 10 tahun

**Tanggung Jawab:**
- Monitoring semua mesin cetak real-time
- Alokasi PO ke mesin dan operator
- Handling masalah produksi
- Koordinasi dengan Khazanah Awal dan Verifikasi
- Reporting ke Production Manager
- Improvement initiatives

**Pain Points:**
- Tidak ada visibility real-time semua mesin
- Harus keliling untuk cek progress satu-satu
- Sulit identifikasi bottleneck dengan cepat
- Laporan manual memakan waktu
- Tidak ada alert otomatis untuk masalah
- Sulit analisa OEE dan performance

**Goals:**
- Dashboard real-time untuk semua mesin
- Alert otomatis untuk masalah
- Report otomatis
- Data-driven decision making
- Improve OEE secara konsisten

---

### 3. Maintenance Technician
**Nama:** Joko Susilo  
**Posisi:** Teknisi Maintenance Unit Cetak  
**Shift:** On-call  
**Pengalaman:** 7 tahun

**Tanggung Jawab:**
- Preventive maintenance mesin cetak
- Repair breakdown mesin
- Tracking kondisi mesin
- Spare part management

**Pain Points:**
- Tidak tahu kondisi mesin real-time
- Breakdown sering tiba-tiba (tidak predictable)
- Sulit tracking history maintenance
- Tidak ada alert preventive maintenance
- Spare part sering tidak ready saat dibutuhkan

**Goals:**
- Predictive maintenance (tahu sebelum breakdown)
- Alert preventive maintenance schedule
- Easy tracking maintenance history
- Spare part inventory management

---

### 4. Production Manager
**Nama:** Ibu Sari Wulandari  
**Posisi:** Production Manager  
**Shift:** Office hours  
**Pengalaman:** 12 tahun

**Tanggung Jawab:**
- Strategic planning produksi
- Monitoring overall performance
- Decision making untuk improvement
- Reporting ke management

**Pain Points:**
- Tidak ada real-time visibility produksi
- Data terlambat (laporan manual di akhir hari)
- Sulit analisa root cause masalah
- Tidak ada predictive insights

**Goals:**
- Real-time dashboard di mobile
- Instant alert untuk critical issues
- Comprehensive analytics & insights
- Data-driven strategic decisions

---

## 📱 User Stories

### Epic 1: Setup & Persiapan Mesin

#### US-CT-001: Melihat Daftar PO yang Siap Cetak
**Sebagai** Operator Cetak  
**Saya ingin** melihat daftar PO yang material-nya sudah siap dari Khazanah Awal  
**Sehingga** saya tahu PO mana yang harus saya kerjakan

**Acceptance Criteria:**
- [ ] Tampil daftar PO dengan status "Siap Cetak"
- [ ] Sorting berdasarkan prioritas dan due date
- [ ] Tampil info per PO:
  - Nomor PO & OBC
  - Jumlah cetak (lembar besar)
  - Kode plat
  - Spesifikasi warna (jumlah warna & nama warna)
  - Due date
  - Prioritas (🔴 Urgent / 🟡 Normal / 🟢 Low)
  - Nomor palet material
  - Waktu material siap
- [ ] Filter berdasarkan: prioritas, tanggal, mesin yang cocok
- [ ] Search berdasarkan nomor PO/OBC
- [ ] Responsive untuk tablet/mobile
- [ ] Real-time update (auto-refresh)

**Business Rules:**
- PO dengan due date < 2 hari = Urgent (🔴)
- PO dengan due date 2-5 hari = Normal (🟡)
- PO dengan due date > 5 hari = Low (🟢)
- Prioritas juga mempertimbangkan sequence produksi (minimize changeover)

---

#### US-CT-002: Assign PO ke Mesin
**Sebagai** Supervisor Unit Cetak  
**Saya ingin** assign PO ke mesin cetak tertentu  
**Sehingga** ada alokasi kerja yang jelas dan terukur

**Acceptance Criteria:**
- [ ] Button "Assign ke Mesin" pada detail PO
- [ ] Dropdown pilih mesin cetak yang available
- [ ] Tampil info per mesin:
  - Nama mesin
  - Status (Available/Running/Down/Setup)
  - Current PO (jika ada)
  - Estimasi selesai current PO
  - Capability (cocok untuk jenis produk ini atau tidak)
  - Last maintenance date
- [ ] Dropdown pilih operator yang akan handle
- [ ] Sistem otomatis suggest mesin optimal berdasarkan:
  - Availability
  - Capability
  - Setup time (minimize changeover)
  - Performance history
- [ ] Konfirmasi assignment
- [ ] Status PO berubah menjadi "Assigned - Menunggu Setup"
- [ ] Notifikasi ke operator yang di-assign

**Business Rules:**
- 1 Mesin hanya bisa handle 1 PO dalam satu waktu
- Mesin dengan status "Down" tidak bisa di-assign
- Prioritas assign ke mesin yang minimize setup time

---

#### US-CT-003: Memulai Setup Mesin
**Sebagai** Operator Cetak  
**Saya ingin** memulai proses setup mesin untuk PO yang di-assign ke saya  
**Sehingga** sistem tracking waktu setup dan progress

**Acceptance Criteria:**
- [ ] Button "Mulai Setup" pada PO yang di-assign
- [ ] Sistem tampilkan checklist setup:
  - ✅ Plat cetak terpasang (scan QR plat untuk validasi)
  - ✅ Tinta warna [X] ready (checklist per warna)
  - ✅ Kertas blanko loaded
  - ✅ Mesin warming up
  - ✅ Test print OK
- [ ] Setiap checklist bisa di-tick saat selesai
- [ ] Timestamp mulai setup tercatat
- [ ] Status mesin berubah menjadi "Setup"
- [ ] Status PO berubah menjadi "Setup in Progress"
- [ ] Timer setup berjalan (real-time)
- [ ] Nama operator tercatat

**Business Rules:**
- Setup time target: ≤ 20 menit
- Semua checklist harus complete sebelum bisa mulai produksi
- Scan QR plat wajib untuk memastikan plat yang benar

---

#### US-CT-004: Test Print & Quality Check
**Sebagai** Operator Cetak  
**Saya ingin** melakukan test print dan validasi kualitas  
**Sehingga** memastikan hasil cetak sesuai spesifikasi sebelum produksi massal

**Acceptance Criteria:**
- [ ] Button "Test Print" setelah setup selesai
- [ ] Input jumlah test print (default: 10 lembar)
- [ ] Checklist validasi kualitas test print:
  - ✅ Warna sesuai spesifikasi
  - ✅ Posisi cetak tepat (tidak miring)
  - ✅ Ketajaman cetak baik
  - ✅ Tidak ada noda atau cacat
- [ ] Upload foto test print (opsional tapi recommended)
- [ ] Jika test print TIDAK OK:
  - Input masalah yang ditemukan
  - Button "Adjustment" → kembali ke setup
  - Timestamp adjustment tercatat
- [ ] Jika test print OK:
  - Button "Approve - Mulai Produksi"
  - Timestamp approval tercatat
  - Status berubah "Ready to Run"

**Business Rules:**
- Test print wajib dilakukan sebelum produksi massal
- Jika test print gagal > 3 kali, escalate ke supervisor
- Foto test print disimpan untuk audit trail

---

#### US-CT-005: Finalisasi Setup & Mulai Produksi
**Sebagai** Operator Cetak  
**Saya ingin** finalisasi setup dan mulai produksi  
**Sehingga** mesin bisa mulai cetak massal

**Acceptance Criteria:**
- [ ] Button "Mulai Produksi" setelah test print approved
- [ ] Sistem validasi:
  - ✅ Semua checklist setup complete
  - ✅ Test print approved
  - ✅ Plat tervalidasi
- [ ] Konfirmasi mulai produksi
- [ ] Timestamp mulai produksi tercatat
- [ ] Status mesin berubah menjadi "Running"
- [ ] Status PO berubah menjadi "In Production"
- [ ] Sistem hitung total setup time
- [ ] Sistem mulai tracking production counter
- [ ] Notifikasi ke supervisor bahwa produksi dimulai

**Business Rules:**
- Setup time = Timestamp Mulai Produksi - Timestamp Mulai Setup
- Setup time masuk ke OEE calculation (Availability loss)
- Operator tidak bisa mulai produksi jika validasi tidak pass

---

### Epic 2: Monitoring Produksi Real-Time

#### US-CT-006: Input Production Counter Real-Time
**Sebagai** Operator Cetak  
**Saya ingin** input jumlah hasil cetak secara real-time  
**Sehingga** progress produksi ter-track dengan akurat

**Acceptance Criteria:**
- [ ] Input field "Jumlah Cetak" (lembar besar)
- [ ] Bisa input via:
  - Manual input (mobile/tablet)
  - Auto-increment button (+100, +500, +1000)
  - Scan barcode (jika ada counter otomatis)
- [ ] Tampil real-time:
  - Target: [X] lembar
  - Actual: [Y] lembar
  - Progress: [Z]% (gauge chart)
  - Gap: [X-Y] lembar (warna merah jika behind, hijau jika on-track)
  - Estimasi selesai: [Waktu] (berdasarkan current rate)
- [ ] Auto-save setiap input (tidak perlu klik save)
- [ ] Timestamp setiap update tercatat
- [ ] History update tersimpan (audit trail)

**Business Rules:**
- Update minimal setiap 30 menit (reminder jika belum update)
- Counter tidak bisa dikurangi (hanya bisa naik)
- Jika gap > 10%, tampil warning ke operator

---

#### US-CT-007: Monitoring Production Rate
**Sebagai** Operator Cetak  
**Saya ingin** melihat production rate real-time  
**Sehingga** saya tahu apakah speed produksi sesuai target

**Acceptance Criteria:**
- [ ] Tampil di dashboard operator:
  - **Current Rate:** [X] lembar/jam
  - **Target Rate:** [Y] lembar/jam (dari standard time)
  - **Performance:** [Z]% (actual/target × 100%)
  - **Speed Loss:** [N]% (jika ada)
- [ ] Gauge chart untuk visualisasi performance
- [ ] Color coding:
  - 🟢 Green: Performance ≥ 95%
  - 🟡 Yellow: Performance 80-94%
  - 🔴 Red: Performance < 80%
- [ ] Chart trend production rate (per jam)
- [ ] Alert jika performance drop < 80%
- [ ] Auto-refresh setiap 1 menit

**Business Rules:**
- Performance = (Actual Output / Theoretical Output) × 100%
- Theoretical Output = Operating Time × Ideal Cycle Time
- Ideal cycle time per produk sudah di-define di master data

---

#### US-CT-008: Recording Downtime
**Sebagai** Operator Cetak  
**Saya ingin** record downtime saat mesin berhenti  
**Sehingga** ada data akurat untuk OEE calculation

**Acceptance Criteria:**
- [ ] Button "Mesin Stop" (prominent, easy access)
- [ ] Saat klik "Mesin Stop":
  - Timestamp stop tercatat
  - Status mesin berubah "Down"
  - Timer downtime mulai berjalan (real-time)
  - Dropdown pilih alasan stop:
    - 🔴 Breakdown (Equipment Failure)
    - 🟡 Waiting Material
    - 🟡 Waiting Operator
    - 🟡 Minor Stop (< 5 menit)
    - 🔵 Changeover/Setup
    - 🟢 Planned Maintenance
    - 🟢 Break Time
    - Lainnya (input manual)
  - Jika breakdown, wajib input deskripsi masalah
  - Jika breakdown, auto-notifikasi ke maintenance
- [ ] Button "Mesin Start" untuk resume produksi
- [ ] Saat klik "Mesin Start":
  - Timestamp start tercatat
  - Status mesin berubah "Running"
  - Sistem hitung durasi downtime
  - Downtime tercatat di log
- [ ] Alert jika downtime > 15 menit (notifikasi ke supervisor)
- [ ] Alert jika downtime > 30 menit (escalate ke manager)

**Business Rules:**
- Downtime classification:
  - Planned: Break, Planned Maintenance (tidak masuk OEE calculation)
  - Unplanned: Breakdown, Waiting, Minor Stop (masuk OEE calculation)
- Downtime > 5 menit wajib dicatat
- Downtime < 5 menit (minor stop) di-aggregate

---

#### US-CT-009: Recording Quality Issues During Production
**Sebagai** Operator Cetak  
**Saya ingin** record quality issues yang ditemukan saat produksi  
**Sehingga** ada tracking HCTS dan bisa improvement

**Acceptance Criteria:**
- [ ] Button "Catat Kerusakan" (easy access)
- [ ] Input:
  - Jumlah lembar rusak
  - Jenis kerusakan (dropdown):
    - Warna pudar
    - Warna tidak sesuai
    - Posisi cetak miring
    - Cetakan tidak tajam
    - Noda tinta
    - Sobek/rusak fisik
    - Lainnya (input manual)
  - Waktu kejadian (auto-fill)
  - Foto kerusakan (opsional)
  - Tindakan yang sudah dilakukan
- [ ] Sistem otomatis kurangi dari good output
- [ ] Sistem otomatis hitung defect rate real-time
- [ ] Jika defect rate > 5%, alert ke operator & supervisor
- [ ] History quality issues tersimpan

**Business Rules:**
- Defect rate = (Total Defect / Total Output) × 100%
- Target defect rate: < 2%
- Jika defect rate > 5%, wajib stop produksi dan investigasi

---

#### US-CT-010: Monitoring OEE Real-Time (Operator View)
**Sebagai** Operator Cetak  
**Saya ingin** melihat OEE mesin saya real-time  
**Sehingga** saya tahu performa saya dan bisa improve

**Acceptance Criteria:**
- [ ] Tampil di dashboard operator (prominent):
  - **OEE Score:** [X]% (gauge chart besar)
  - **Availability:** [Y]%
  - **Performance:** [Z]%
  - **Quality:** [W]%
- [ ] Color coding:
  - 🟢 Green: OEE ≥ 85% (World Class)
  - 🟡 Yellow: OEE 60-84% (Good)
  - 🔴 Red: OEE < 60% (Needs Improvement)
- [ ] Breakdown losses:
  - Downtime: [X] menit
  - Speed Loss: [Y]%
  - Quality Loss: [Z] lembar
- [ ] Comparison dengan:
  - Target OEE
  - Shift sebelumnya
  - Best performer hari ini
- [ ] Tips untuk improve OEE (contextual)
- [ ] Auto-refresh setiap 1 menit

**Business Rules:**
- OEE = Availability × Performance × Quality
- Availability = (Planned Production Time - Downtime) / Planned Production Time
- Performance = (Actual Output / Theoretical Output) × 100%
- Quality = (Good Output / Total Output) × 100%

---

### Epic 3: Completion & Handover

#### US-CT-011: Finalisasi Produksi
**Sebagai** Operator Cetak  
**Saya ingin** finalisasi produksi saat PO selesai  
**Sehingga** hasil bisa dikirim ke Khazanah Awal untuk penghitungan

**Acceptance Criteria:**
- [ ] Button "Selesai Produksi" (muncul saat actual ≥ target)
- [ ] Sistem tampilkan summary:
  - Target: [X] lembar
  - Actual: [Y] lembar
  - Selisih: [Y-X] lembar
  - Durasi produksi: [N] jam [M] menit
  - OEE: [Z]%
  - Downtime total: [N] menit
  - Quality issues: [N] lembar
- [ ] Konfirmasi finalisasi
- [ ] Input nomor palet hasil cetak
- [ ] Upload foto hasil cetak (opsional)
- [ ] Timestamp selesai produksi tercatat
- [ ] Status mesin berubah "Available"
- [ ] Status PO berubah "Selesai Cetak - Menunggu Penghitungan"
- [ ] Notifikasi ke Khazanah Awal untuk penghitungan
- [ ] Data tidak bisa diubah setelah finalisasi

**Business Rules:**
- Operator bisa finalisasi jika actual ≥ target
- Jika actual < target, wajib isi alasan dan approval supervisor
- Durasi produksi = Timestamp Selesai - Timestamp Mulai Produksi (exclude downtime)

---

#### US-CT-012: Handover ke Shift Berikutnya
**Sebagai** Operator Cetak  
**Saya ingin** handover PO yang belum selesai ke shift berikutnya  
**Sehingga** produksi bisa continue tanpa gap

**Acceptance Criteria:**
- [ ] Button "Handover ke Shift Berikutnya"
- [ ] Form handover:
  - Progress saat ini: [X]% ([Y] dari [Z] lembar)
  - Estimasi sisa waktu: [N] jam
  - Kondisi mesin: (dropdown: Baik/Perlu Perhatian/Bermasalah)
  - Issue yang perlu diperhatikan: (text area)
  - Tinta remaining: (per warna)
  - Kertas remaining: [X] lembar
  - Notes untuk shift berikutnya: (text area)
- [ ] Upload foto kondisi mesin (opsional)
- [ ] Konfirmasi handover
- [ ] Timestamp handover tercatat
- [ ] Status PO berubah "In Production - Handover"
- [ ] Notifikasi ke operator shift berikutnya
- [ ] Operator shift berikutnya wajib acknowledge handover

**Business Rules:**
- Handover hanya bisa dilakukan di akhir shift
- Operator shift berikutnya wajib review handover notes
- Jika ada issue critical, wajib involve supervisor

---

### Epic 4: Dashboard & Monitoring (Supervisor)

#### US-CT-013: Dashboard Overview Unit Cetak
**Sebagai** Supervisor Unit Cetak  
**Saya ingin** melihat overview semua mesin cetak real-time  
**Sehingga** saya punya visibility penuh dan bisa cepat respond jika ada masalah

**Acceptance Criteria:**
- [ ] Tampil di dashboard (layout seperti Andon Board):
  ```
  ┌─────────────────────────────────────────────────────────┐
  │  UNIT CETAK - REAL-TIME MONITORING                      │
  ├────────┬────────┬──────────┬─────────┬─────────┬────────┤
  │ Mesin  │ Status │ PO       │ Output  │ Target  │  OEE   │
  ├────────┼────────┼──────────┼─────────┼─────────┼────────┤
  │ MC-01  │  🟢    │ PO-001   │ 12,500  │ 15,000  │  87%   │
  │ MC-02  │  🟡    │ PO-002   │  8,200  │ 10,000  │  72%   │
  │ MC-03  │  🔴    │ PO-003   │  3,100  │  8,000  │  45%   │
  │ MC-04  │  🟢    │ PO-004   │ 18,900  │ 20,000  │  91%   │
  │ MC-05  │  ⚪    │ -        │    -    │    -    │   -    │
  └────────┴────────┴──────────┴─────────┴─────────┴────────┘
  
  🔴 ACTIVE ALERTS:
  - MC-03: BREAKDOWN - Downtime: 45 menit
  - MC-02: BEHIND TARGET - Gap: 1,800 lembar
  ```
- [ ] Status color coding:
  - 🟢 Green: Running, on-target
  - 🟡 Yellow: Running, behind target
  - 🔴 Red: Down/Problem
  - ⚪ White: Available/Idle
  - 🔵 Blue: Setup
- [ ] Click mesin untuk detail view
- [ ] Summary metrics:
  - Total mesin running: [X]/[Y]
  - Total output hari ini: [X] lembar
  - Average OEE: [Y]%
  - Total downtime: [Z] menit
  - Active alerts: [N]
- [ ] Auto-refresh setiap 10 detik
- [ ] Full-screen mode (untuk TV display)

**Business Rules:**
- Status mesin update real-time berdasarkan input operator
- Alert prioritization: Red > Yellow > Blue
- Dashboard harus load < 2 detik

---

#### US-CT-014: Monitoring OEE per Mesin (Supervisor)
**Sebagai** Supervisor Unit Cetak  
**Saya ingin** melihat OEE detail per mesin  
**Sehingga** bisa identifikasi mesin mana yang perlu improvement

**Acceptance Criteria:**
- [ ] Tampil list semua mesin dengan OEE metrics:
  - Mesin ID & Name
  - Current OEE (gauge chart)
  - Availability (%)
  - Performance (%)
  - Quality (%)
  - Breakdown losses:
    - Downtime: [X] menit
    - Speed loss: [Y]%
    - Quality loss: [Z] lembar
- [ ] Sorting berdasarkan: OEE (low to high), mesin ID
- [ ] Filter berdasarkan: status, shift, periode
- [ ] Comparison:
  - Mesin vs mesin
  - Hari ini vs kemarin
  - Shift ini vs shift sebelumnya
- [ ] Chart: OEE trend per mesin (7 hari terakhir)
- [ ] Export to Excel
- [ ] Drill-down ke detail losses

**Business Rules:**
- OEE calculation real-time berdasarkan data operator
- Target OEE: ≥ 85%
- Alert jika OEE < 60%

---

#### US-CT-015: Six Big Losses Analysis
**Sebagai** Supervisor Unit Cetak  
**Saya ingin** melihat analisa Six Big Losses  
**Sehingga** bisa fokus improvement di area yang paling impact

**Acceptance Criteria:**
- [ ] Dashboard Six Big Losses dengan breakdown:
  
  **Availability Losses:**
  1. **Breakdown Losses**
     - Total downtime: [X] menit
     - Frequency: [Y] kali
     - MTBF: [Z] jam
     - MTTR: [W] menit
     - Top 5 breakdown reasons (Pareto chart)
  
  2. **Setup & Adjustment Losses**
     - Total setup time: [X] menit
     - Average setup time: [Y] menit/changeover
     - Target setup time: [Z] menit
     - Gap: [Y-Z] menit
  
  **Performance Losses:**
  3. **Idling & Minor Stops**
     - Frequency: [X] kali
     - Total time: [Y] menit
     - Average duration: [Z] menit
  
  4. **Speed Losses**
     - Ideal speed: [X] lembar/jam
     - Actual speed: [Y] lembar/jam
     - Speed loss: [Z]%
  
  **Quality Losses:**
  5. **Quality Defects**
     - Total defects: [X] lembar
     - Defect rate: [Y]%
     - Top 5 defect types (Pareto chart)
  
  6. **Startup Losses**
     - Waste during startup: [X] lembar
     - Time to stable production: [Y] menit

- [ ] Waterfall chart untuk visualisasi losses
- [ ] Filter berdasarkan: mesin, periode, shift
- [ ] Drill-down untuk detail per loss category
- [ ] Export to PDF/Excel

**Business Rules:**
- Six Big Losses = penyebab utama OEE tidak optimal
- Focus improvement di top 2-3 losses (80/20 rule)

---

#### US-CT-016: Alert & Notification untuk Supervisor
**Sebagai** Supervisor Unit Cetak  
**Saya ingin** menerima alert otomatis jika ada masalah  
**Sehingga** bisa cepat respond dan minimize impact

**Acceptance Criteria:**
- [ ] Alert otomatis untuk kondisi:
  - 🔴 **Critical:**
    - Mesin breakdown
    - OEE drop > 20% dari average
    - Defect rate > 10%
    - Downtime > 30 menit
    - Production stop tanpa reason
  - 🟡 **Warning:**
    - Behind target > 10%
    - OEE < 70%
    - Defect rate > 5%
    - Downtime > 15 menit
    - Minor stops frequency > 5 kali/jam
  - 🔵 **Info:**
    - Setup time > target
    - Shift handover pending
    - Material running low
- [ ] Notification channels:
  - In-app notification (badge counter)
  - Push notification (mobile)
  - WhatsApp (untuk critical alert)
  - Email (untuk daily summary)
- [ ] Alert management:
  - Alert prioritization
  - Alert acknowledgment
  - Alert resolution tracking
  - Response time measurement
- [ ] Alert history & analytics

**Business Rules:**
- Critical alert wajib di-acknowledge dalam 5 menit
- Jika tidak di-acknowledge, escalate ke manager
- Alert resolution wajib di-document

---

#### US-CT-017: Production Planning & Scheduling
**Sebagai** Supervisor Unit Cetak  
**Saya ingin** optimize production schedule  
**Sehingga** maximize throughput dan minimize changeover

**Acceptance Criteria:**
- [ ] Tampil production queue:
  - PO yang menunggu (prioritized)
  - PO yang sedang dikerjakan
  - PO yang selesai hari ini
- [ ] Drag-and-drop untuk re-sequence PO
- [ ] Sistem suggest optimal sequence berdasarkan:
  - Due date
  - Setup time (minimize changeover)
  - Machine capability
  - Material availability
- [ ] Tampil estimasi:
  - Total production time
  - Completion time per PO
  - Machine utilization
- [ ] Gantt chart untuk visualisasi schedule
- [ ] What-if scenario analysis
- [ ] Save & publish schedule
- [ ] Notifikasi ke operator tentang schedule

**Business Rules:**
- Prioritas: Due date > Setup optimization > Machine utilization
- Minimize changeover untuk efisiensi
- Consider machine capability (tidak semua mesin bisa semua produk)

---

### Epic 5: Performance Analytics

#### US-CT-018: Operator Performance Dashboard
**Sebagai** Supervisor Unit Cetak  
**Saya ingin** melihat performa individual operator  
**Sehingga** bisa evaluasi, coaching, dan recognition

**Acceptance Criteria:**
- [ ] Tampil list operator dengan metrics:
  
  **Productivity:**
  - Output per shift (lembar)
  - Production rate (lembar/jam)
  - Target achievement (%)
  - Efficiency score
  
  **Quality:**
  - Defect rate (%)
  - First Pass Yield (%)
  - Quality score
  
  **Equipment Handling:**
  - Machine uptime saat bertugas (%)
  - Setup time (average)
  - Minor stops frequency
  - Equipment care score
  
  **OEE Contribution:**
  - Average OEE saat bertugas
  - Availability contribution
  - Performance contribution
  - Quality contribution

- [ ] Operator scorecard (detail per operator):
  ```
  ┌─────────────────────────────────────────────────┐
  │  OPERATOR CETAK: Andi Wijaya                    │
  │  Mesin: MC-02 | Shift: Siang | 27/12/25        │
  ├─────────────────────────────────────────────────┤
  │  📊 PRODUCTIVITY                                │
  │  ├─ Output: 18,500 lembar                      │
  │  ├─ Target: 16,000 lembar                      │
  │  ├─ Achievement: 115.6% ⭐                      │
  │  └─ Rate: 2,312 lembar/jam                     │
  ├─────────────────────────────────────────────────┤
  │  ✅ QUALITY                                     │
  │  ├─ Defect Rate: 1.8%                          │
  │  ├─ FPY: 98.2%                                 │
  │  └─ Quality Score: 95/100                      │
  ├─────────────────────────────────────────────────┤
  │  ⚙️ EQUIPMENT EFFICIENCY                        │
  │  ├─ Machine Uptime: 94%                        │
  │  ├─ Setup Time: 18 menit (Target: 20)         │
  │  ├─ Minor Stops: 3 kali                        │
  │  └─ Equipment Score: 92/100                    │
  ├─────────────────────────────────────────────────┤
  │  🏆 OVERALL SCORE: 94/100 (Very Good)          │
  └─────────────────────────────────────────────────┘
  ```

- [ ] Leaderboard (top 10 operators)
- [ ] Peer comparison (individual vs team average)
- [ ] Trend analysis (performance over time)
- [ ] Filter: shift, periode, mesin
- [ ] Export to Excel

**Business Rules:**
- Overall score = weighted average (Productivity 40%, Quality 30%, Equipment 30%)
- Ranking update real-time
- Recognition untuk top 3 performers

---

#### US-CT-019: Shift Performance Comparison
**Sebagai** Supervisor Unit Cetak  
**Saya ingin** compare performa antar shift  
**Sehingga** bisa identifikasi best practice dan improvement area

**Acceptance Criteria:**
- [ ] Tampil comparison 3 shift (Pagi/Siang/Malam):
  
  **Productivity:**
  - Total output per shift
  - Average production rate
  - Target achievement
  
  **Quality:**
  - Defect rate per shift
  - Quality score
  
  **OEE:**
  - Average OEE per shift
  - Availability per shift
  - Performance per shift
  - Quality per shift
  
  **Issues:**
  - Downtime per shift
  - Breakdown frequency
  - Quality issues frequency

- [ ] Chart: Trend comparison (7 hari terakhir)
- [ ] Best performing shift highlight
- [ ] Gap analysis (best vs worst)
- [ ] Root cause analysis untuk gap
- [ ] Recommendations

**Business Rules:**
- Fair comparison (normalize by working hours & machine availability)
- Consider shift team composition
- Identify systemic issues vs random variation

---

#### US-CT-020: Machine Performance Analytics
**Sebagai** Production Manager  
**Saya ingin** melihat analisa performa per mesin  
**Sehingga** bisa decision making untuk maintenance, replacement, atau improvement

**Acceptance Criteria:**
- [ ] Tampil analytics per mesin:
  
  **Performance Metrics:**
  - Average OEE (trend)
  - Utilization rate (%)
  - Throughput (lembar/hari)
  - Cycle time
  
  **Reliability:**
  - MTBF (Mean Time Between Failures)
  - MTTR (Mean Time To Repair)
  - Breakdown frequency
  - Availability rate
  
  **Quality:**
  - Defect rate per mesin
  - Quality consistency
  
  **Cost:**
  - Maintenance cost
  - Downtime cost
  - Cost per unit produced

- [ ] Benchmarking (mesin vs mesin)
- [ ] Trend analysis (performance over time)
- [ ] Predictive insights (maintenance recommendation)
- [ ] ROI analysis (improvement initiatives)
- [ ] Export to PDF/Excel

**Business Rules:**
- Data-driven decision untuk investment
- Identify underperforming machines
- Prioritize improvement initiatives by ROI

---

### Epic 6: Maintenance Management

#### US-CT-021: Preventive Maintenance Schedule
**Sebagai** Maintenance Technician  
**Saya ingin** melihat schedule preventive maintenance  
**Sehingga** bisa plan dan execute on-time

**Acceptance Criteria:**
- [ ] Tampil PM schedule per mesin:
  - Mesin ID & Name
  - Last PM date
  - Next PM due date
  - PM frequency (e.g., setiap 500 jam operasi)
  - Operating hours since last PM
  - Status (🟢 OK / 🟡 Due Soon / 🔴 Overdue)
- [ ] Calendar view untuk PM schedule
- [ ] Alert PM due dalam 7 hari
- [ ] Alert PM overdue
- [ ] Button "Execute PM" untuk mulai maintenance
- [ ] PM checklist (per mesin type)
- [ ] Record PM completion:
  - Timestamp
  - Technician
  - Work performed
  - Parts replaced
  - Findings
  - Next PM date
- [ ] PM history per mesin

**Business Rules:**
- PM schedule berdasarkan operating hours atau calendar
- PM overdue → mesin tidak boleh digunakan
- PM execution wajib follow checklist

---

#### US-CT-022: Breakdown Management
**Sebagai** Maintenance Technician  
**Saya ingin** manage breakdown repair dengan efisien  
**Sehingga** minimize downtime

**Acceptance Criteria:**
- [ ] Notifikasi otomatis saat operator report breakdown
- [ ] Breakdown ticket berisi:
  - Mesin ID
  - Problem description (dari operator)
  - Timestamp breakdown
  - Current downtime duration
  - Priority (based on impact)
  - PO yang terdampak
- [ ] Button "Acknowledge" (technician acknowledge ticket)
- [ ] Button "Start Repair" (mulai perbaikan)
- [ ] Input during repair:
  - Root cause
  - Action taken
  - Parts used
  - Photos (before/after)
- [ ] Button "Complete Repair" (selesai perbaikan)
- [ ] Sistem hitung MTTR (Mean Time To Repair)
- [ ] Feedback dari operator (mesin OK atau belum)
- [ ] Breakdown history & analytics

**Business Rules:**
- Response time target: < 10 menit
- MTTR target: < 30 menit (untuk minor issues)
- Escalate jika repair > 1 jam

---

#### US-CT-023: Spare Part Inventory Management
**Sebagai** Maintenance Technician  
**Saya ingin** manage spare part inventory  
**Sehingga** parts selalu ready saat dibutuhkan

**Acceptance Criteria:**
- [ ] Tampil list spare parts:
  - Part name & code
  - Stock quantity
  - Minimum stock
  - Location
  - Unit price
  - Last restock date
- [ ] Alert jika stock < minimum
- [ ] Record part usage (saat repair):
  - Part used
  - Quantity
  - Mesin
  - Ticket number
  - Technician
- [ ] Sistem auto-deduct stock
- [ ] Purchase request (jika stock low)
- [ ] Stock history & consumption trend
- [ ] ABC analysis (critical parts identification)

**Business Rules:**
- Critical parts (high usage, high impact) harus selalu ready
- Auto-generate purchase request jika stock < minimum
- Track part consumption per mesin (identify high-maintenance machines)

---

#### US-CT-024: Predictive Maintenance Insights
**Sebagai** Maintenance Manager  
**Saya ingin** predictive maintenance insights  
**Sehingga** bisa prevent breakdown sebelum terjadi

**Acceptance Criteria:**
- [ ] Machine health scoring per mesin:
  - Overall health score (0-100)
  - Breakdown risk (Low/Medium/High)
  - Factors contributing to risk:
    - Operating hours since last PM
    - Breakdown history
    - Performance degradation
    - Age of machine
- [ ] Predictive alerts:
  - "Mesin MC-03 berisiko breakdown dalam 7 hari" (berdasarkan pattern)
  - "Mesin MC-01 performa menurun 15% → perlu inspection"
- [ ] Maintenance recommendations:
  - Prioritized action items
  - Estimated impact (downtime prevention)
  - Estimated cost
- [ ] Trend analysis:
  - MTBF trend (improving or degrading?)
  - Breakdown pattern (time-based, usage-based?)
- [ ] ROI of predictive maintenance

**Business Rules:**
- Predictive model based on historical data & ML
- Prioritize high-risk machines
- Validate predictions dengan actual outcomes (continuous learning)

---

### Epic 7: Gamification & Motivation

#### US-CT-025: Achievement & Badges untuk Operator
**Sebagai** Operator Cetak  
**Saya ingin** earn achievements dan badges  
**Sehingga** termotivasi untuk perform better

**Acceptance Criteria:**
- [ ] Achievement system:
  - 🏅 **"Quality Champion"** - HCS 99%+ selama seminggu
  - 🏅 **"Production Hero"** - Exceed target 125%+
  - 🏅 **"Zero Defect"** - No HCTS selama shift
  - 🏅 **"Machine Master"** - OEE 95%+ selama seminggu
  - 🏅 **"Fast Setup"** - Setup time < 80% standard time
  - 🏅 **"Consistency King"** - Stable performance 30 hari
  - 🏅 **"Perfect Week"** - 100% target achievement 5 hari kerja
  - 🏅 **"Speed Demon"** - Production rate 120%+ dari standard

- [ ] Badge display:
  - Di profile operator
  - Di leaderboard
  - Di Andon board (saat operator login)
  - Notification saat earn badge

- [ ] Badge rarity:
  - 🥉 Bronze: Easy to achieve
  - 🥈 Silver: Moderate difficulty
  - 🥇 Gold: Hard to achieve
  - 💎 Diamond: Very rare

- [ ] Progress tracking untuk each badge
- [ ] Badge showcase (collection)

**Business Rules:**
- Badge earned berdasarkan actual data (tidak bisa dimanipulasi)
- Badge permanent (tidak bisa hilang)
- Rare badge = higher recognition

---

#### US-CT-026: Point System & Rewards
**Sebagai** Operator Cetak  
**Saya ingin** earn points dan redeem rewards  
**Sehingga** ada tangible benefit dari performa baik

**Acceptance Criteria:**
- [ ] Earn points for:
  - ✅ Achieve daily target (+10 points)
  - ✅ Exceed target 110% (+20 points)
  - ✅ Exceed target 120% (+30 points)
  - ✅ Zero defect shift (+25 points)
  - ✅ OEE ≥ 90% (+20 points)
  - ✅ Fast setup (< target) (+15 points)
  - ✅ No downtime shift (+30 points)
  - ✅ Help colleague (+10 points)
  - ✅ Submit improvement idea (+15 points)

- [ ] Lose points for:
  - ❌ Miss target (-10 points)
  - ❌ High defect rate > 5% (-15 points)
  - ❌ OEE < 60% (-10 points)
  - ❌ Late attendance (-5 points)
  - ❌ Safety violation (-30 points)

- [ ] Point balance display (real-time)
- [ ] Point history & transactions
- [ ] Redeem points:
  - 100 points: Extra 30 menit break
  - 200 points: Shift preference (pilih shift)
  - 300 points: Voucher makan Rp 50.000
  - 500 points: Voucher belanja Rp 100.000
  - 1000 points: Training voucher
  - 2000 points: Bonus Rp 500.000

- [ ] Redemption history
- [ ] Point expiry (1 tahun)

**Business Rules:**
- Points earned/lost based on actual performance data
- Redemption approval by supervisor
- Point balance tidak bisa negatif (minimum 0)

---

#### US-CT-027: Leaderboard & Competition
**Sebagai** Operator Cetak  
**Saya ingin** melihat ranking saya di leaderboard  
**Sehingga** termotivasi untuk compete dan improve

**Acceptance Criteria:**
- [ ] **Daily Leaderboard:**
  - Top 10 operators hari ini
  - Ranking berdasarkan overall score
  - Real-time update
  - Display di Andon board

- [ ] **Weekly Leaderboard:**
  - Top 10 operators minggu ini
  - Consistent performers
  - Most improved

- [ ] **Monthly Leaderboard:**
  - Top 10 operators bulan ini
  - Department champions
  - Hall of Fame

- [ ] **Team Leaderboard:**
  - Best performing shift
  - Inter-shift competition

- [ ] Leaderboard display:
  - Rank
  - Operator name
  - Score
  - Badges earned
  - Trend (🔺 naik / 🔻 turun / ➡️ stabil)

- [ ] My position highlight
- [ ] Filter: daily/weekly/monthly, by shift

**Business Rules:**
- Fair ranking (normalize by working hours & machine assigned)
- Transparent calculation
- Recognition untuk top 3 (di Andon board & WhatsApp group)

---

#### US-CT-028: Daily Challenges
**Sebagai** Operator Cetak  
**Saya ingin** participate in daily challenges  
**Sehingga** ada variety dan fun di pekerjaan

**Acceptance Criteria:**
- [ ] Daily challenge system:
  - **"Beat Your Best"** - Exceed personal record output
  - **"Zero Defect Day"** - No defects selama shift
  - **"Speed Run"** - Achieve target 2 jam lebih cepat
  - **"Perfect Setup"** - Setup time < 15 menit
  - **"Uptime Hero"** - No downtime selama shift
  - **"Team Challenge"** - Team total output > target

- [ ] Challenge display:
  - Challenge of the day (announced di pagi hari)
  - Challenge description & rules
  - Reward (points + badge)
  - Progress tracking (real-time)
  - Participants & standings

- [ ] Challenge completion:
  - Auto-detect completion
  - Instant notification
  - Reward auto-credited
  - Celebration (confetti animation 🎉)

- [ ] Challenge history
- [ ] Suggest challenges (operator bisa suggest)

**Business Rules:**
- 1 challenge per day
- All operators eligible
- Challenge designed untuk achievable tapi challenging
- Rotate challenge types untuk variety

---

### Epic 8: Mobile App untuk Operator

#### US-CT-029: Mobile App - Production Monitoring
**Sebagai** Operator Cetak  
**Saya ingin** monitor dan input data via mobile  
**Sehingga** tidak perlu bolak-balik ke komputer

**Acceptance Criteria:**
- [ ] Responsive mobile web (PWA)
- [ ] Mobile dashboard:
  - Current PO info
  - Progress (gauge chart)
  - Target vs actual
  - OEE score (real-time)
  - Timer (production time, downtime)

- [ ] Quick actions (prominent buttons):
  - ➕ Update Counter
  - ⏸️ Mesin Stop
  - ▶️ Mesin Start
  - ⚠️ Catat Kerusakan
  - ✅ Selesai Produksi

- [ ] Input forms optimized untuk mobile:
  - Large input fields
  - Number pad untuk input angka
  - Dropdown dengan search
  - Camera integration untuk foto

- [ ] Offline mode:
  - Cache data untuk view
  - Queue input saat offline
  - Auto-sync saat online

- [ ] Push notifications
- [ ] Fast & lightweight (< 1 MB)

**Business Rules:**
- Mobile-first design (prioritas UX mobile)
- Touch-friendly (button min 44×44 px)
- Work di lingkungan pabrik (high contrast, large font)

---

#### US-CT-030: Mobile App - Personal Dashboard
**Sebagai** Operator Cetak  
**Saya ingin** melihat personal performance via mobile  
**Sehingga** bisa track progress kapan saja

**Acceptance Criteria:**
- [ ] **My Performance:**
  - Today's stats (real-time)
  - This week summary
  - This month summary
  - Trend chart

- [ ] **My Achievements:**
  - Badges earned (showcase)
  - Points balance
  - Leaderboard position
  - Milestone progress

- [ ] **My Goals:**
  - Daily target
  - Weekly goals
  - Personal KPIs
  - Progress tracking

- [ ] **My Schedule:**
  - Today's shift
  - This week schedule
  - Shift swap request

- [ ] **Notifications:**
  - Performance alerts
  - Badge earned
  - Leaderboard update
  - Challenge completion

**Business Rules:**
- Personal data only (tidak bisa lihat data operator lain)
- Real-time sync
- Privacy-protected

---

### Epic 9: Reporting & Analytics

#### US-CT-031: Daily Production Report
**Sebagai** Supervisor Unit Cetak  
**Saya ingin** generate daily production report otomatis  
**Sehingga** tidak perlu buat laporan manual

**Acceptance Criteria:**
- [ ] Auto-generate report di akhir shift
- [ ] Isi report:
  
  **Production Summary:**
  - Total output per mesin
  - Total output unit cetak
  - Target vs actual
  - Achievement percentage
  
  **OEE Summary:**
  - Average OEE per mesin
  - Average OEE unit cetak
  - OEE breakdown (A × P × Q)
  
  **Downtime Summary:**
  - Total downtime per mesin
  - Downtime breakdown by reason
  - Top 5 downtime causes
  
  **Quality Summary:**
  - Total defects
  - Defect rate
  - Top 5 defect types
  
  **Operator Performance:**
  - Top 3 performers
  - Performance summary per operator
  
  **Issues & Actions:**
  - Active issues
  - Resolved issues
  - Pending actions
  
  **Recommendations:**
  - Improvement suggestions
  - Action items untuk besok

- [ ] Format: PDF & Excel
- [ ] Auto-send via email ke Production Manager
- [ ] Download manual kapan saja
- [ ] Customizable template

**Business Rules:**
- Report generated automatically (no manual intervention)
- Data accurate & real-time
- Professional format (ready untuk management)

---

#### US-CT-032: Weekly Performance Report
**Sebagai** Production Manager  
**Saya ingin** weekly performance report dengan insights  
**Sehingga** bisa review dan strategic planning

**Acceptance Criteria:**
- [ ] Auto-generate setiap Senin pagi (untuk minggu sebelumnya)
- [ ] Isi report:
  
  **Executive Summary:**
  - Key highlights
  - Key achievements
  - Key challenges
  - Action items
  
  **Production Performance:**
  - Total output minggu ini
  - Comparison vs target
  - Comparison vs minggu lalu
  - Trend analysis
  
  **OEE Analysis:**
  - Average OEE minggu ini
  - OEE trend (4 minggu terakhir)
  - Best performing machines
  - Underperforming machines
  - Six Big Losses breakdown
  
  **Quality Analysis:**
  - Defect rate trend
  - Top defect types
  - Quality improvement initiatives
  
  **People Performance:**
  - Top performers (recognition)
  - Performance distribution
  - Training needs identification
  
  **Maintenance:**
  - Breakdown summary
  - MTBF & MTTR
  - PM compliance
  
  **Strategic Insights:**
  - Bottleneck identification
  - Improvement opportunities
  - ROI of initiatives
  - Recommendations

- [ ] Format: PDF (executive-friendly)
- [ ] Auto-send via email
- [ ] Presentation-ready charts & graphs

**Business Rules:**
- Data-driven insights (not just raw data)
- Actionable recommendations
- Professional presentation

---

#### US-CT-033: Custom Report Builder
**Sebagai** Analyst  
**Saya ingin** build custom reports  
**Sehingga** bisa analisa data sesuai kebutuhan spesifik

**Acceptance Criteria:**
- [ ] Drag-and-drop report builder
- [ ] Select dimensions:
  - Mesin
  - Operator
  - Shift
  - PO/OBC
  - Tanggal/Periode
  - Produk

- [ ] Select metrics:
  - Output
  - OEE (A, P, Q)
  - Downtime
  - Defect rate
  - Setup time
  - Cycle time
  - Cost

- [ ] Select chart type:
  - Line chart (trend)
  - Bar chart (comparison)
  - Pie chart (composition)
  - Gauge chart (KPI)
  - Heatmap (pattern)
  - Pareto chart
  - Waterfall chart

- [ ] Filters & slicers
- [ ] Drill-down capability
- [ ] Save report template
- [ ] Schedule report delivery
- [ ] Export: Excel, PDF, CSV

**Business Rules:**
- User-friendly (no SQL knowledge required)
- Fast query (< 5 detik)
- Data access control (role-based)

---

### Epic 10: Integration & Automation

#### US-CT-034: Integrasi dengan SAP - Production Confirmation
**Sebagai** System  
**Saya ingin** otomatis kirim production confirmation ke SAP  
**Sehingga** SAP inventory & production data selalu accurate

**Acceptance Criteria:**
- [ ] Setiap finalisasi produksi:
  - Sistem otomatis kirim data ke SAP:
    - PO number
    - Quantity produced (good output)
    - Quantity scrapped (defects)
    - Production date & time
    - Mesin used
    - Operator
    - Actual hours
- [ ] SAP update:
  - Production order confirmation
  - Inventory update (finished goods)
  - Material consumption (kertas, tinta)
- [ ] Error handling:
  - Retry mechanism (3× retry)
  - Alert jika gagal
  - Manual sync option
- [ ] Log semua transaksi
- [ ] Audit trail

**Business Rules:**
- Real-time sync (immediately after finalization)
- Data validation sebelum kirim
- Rollback mechanism jika error

---

#### US-CT-035: Integrasi dengan Khazanah Awal
**Sebagai** System  
**Saya ingin** seamless integration dengan Khazanah Awal  
**Sehingga** handover material & hasil cetak smooth

**Acceptance Criteria:**
- [ ] **Material Ready Notification:**
  - Khazanah Awal finalisasi penyiapan material
  - → Auto-update status PO di Unit Cetak "Siap Cetak"
  - → Notifikasi ke Supervisor Unit Cetak
  - → Material details available (plat, kertas, tinta, palet)

- [ ] **Production Complete Notification:**
  - Operator finalisasi produksi
  - → Auto-update status PO di Khazanah Awal "Selesai Cetak - Menunggu Penghitungan"
  - → Notifikasi ke Staff Khazanah Awal
  - → Hasil cetak details available (jumlah, palet, issues)

- [ ] Real-time status sync
- [ ] Data consistency check

**Business Rules:**
- Status sync real-time (< 5 detik)
- No manual intervention needed
- Audit trail untuk tracking

---

#### US-CT-036: Integrasi dengan Maintenance System
**Sebagai** System  
**Saya ingin** otomatis create maintenance ticket saat breakdown  
**Sehingga** response time lebih cepat

**Acceptance Criteria:**
- [ ] Saat operator klik "Mesin Stop" dengan reason "Breakdown":
  - Auto-create maintenance ticket
  - Auto-assign ke technician on-duty
  - Push notification ke technician
  - WhatsApp notification (jika critical)

- [ ] Ticket berisi:
  - Mesin ID
  - Problem description
  - Timestamp
  - Priority
  - PO terdampak

- [ ] Tracking:
  - Acknowledgment time
  - Response time
  - Resolution time
  - MTTR

- [ ] Auto-close ticket saat repair complete
- [ ] Feedback loop (operator confirm mesin OK)

**Business Rules:**
- Critical breakdown = immediate notification (WhatsApp)
- Response time target: < 10 menit
- Auto-escalate jika no response dalam 15 menit

---

## 🎯 Key Performance Indicators (KPIs)

### KPI Produksi
- **OEE (Overall Equipment Effectiveness):** Target ≥ 85%
- **Availability:** Target ≥ 90%
- **Performance:** Target ≥ 95%
- **Quality:** Target ≥ 98%
- **Throughput:** Target [X] lembar/hari per mesin
- **Cycle Time:** Target ≤ [Y] jam per PO

### KPI Setup & Changeover
- **Setup Time:** Target ≤ 20 menit
- **Setup Efficiency:** Target ≥ 90%
- **Changeover Loss:** Target ≤ 5% dari production time

### KPI Downtime
- **Total Downtime:** Target ≤ 10% dari planned production time
- **MTBF (Mean Time Between Failures):** Target ≥ 100 jam
- **MTTR (Mean Time To Repair):** Target ≤ 30 menit
- **Breakdown Frequency:** Target ≤ 2 kali/bulan per mesin

### KPI Quality
- **Defect Rate:** Target ≤ 2%
- **First Pass Yield:** Target ≥ 98%
- **Cost of Poor Quality:** Target ≤ 1% dari production cost

### KPI Operator
- **Target Achievement:** Target ≥ 100%
- **Productivity:** Target [X] lembar/jam per operator
- **Quality Score:** Target ≥ 95/100
- **Equipment Care Score:** Target ≥ 90/100

---

## 📱 UI/UX Considerations

### Mobile-First Design
- **Prioritas:** Operator harus bisa input data di lapangan dengan mudah
- **Large Touch Target:** Button minimal 44×44 px
- **Simple Form:** Max 3 field per screen
- **Instant Feedback:** Loading indicator, success/error message
- **Offline Capable:** Basic view & sync later
- **Fast:** Load time < 2 detik

### Andon Board Design (untuk TV Display)
- **Large Display:** Optimized untuk TV/monitor besar (55"+)
- **High Contrast:** Readable dari jarak 10 meter
- **Color Coding:** Traffic light system (Red/Yellow/Green)
- **Live Updates:** Real-time tanpa flicker
- **Sound Alerts:** Audio notification untuk critical alerts
- **Full-screen Mode:** No distractions

### Accessibility
- **Color Blind Friendly:** Jangan hanya pakai warna untuk status (tambah icon/text)
- **Large Font:** Minimal 16px untuk body text, 24px untuk important info
- **High Contrast:** Mudah dibaca di lingkungan pabrik (bright lighting)
- **Icon + Text:** Jangan hanya icon (untuk clarity)

### Performance
- **Fast Load:** < 2 detik untuk dashboard
- **Smooth Scroll:** 60 fps
- **Optimized Image:** Compress foto upload (max 500 KB)
- **Lazy Load:** Load data on demand
- **Real-time Update:** WebSocket untuk live data (tidak polling)

---

## 🔐 Security & Access Control

### Role-Based Access

**Operator Cetak:**
- View: PO assigned to me, My performance, My achievements
- Input: Production data, Downtime, Quality issues
- Cannot: Edit finalized data, View other operator's detail data, Delete data

**Supervisor Unit Cetak:**
- View: All machines, All operators, All reports
- Input: Assign PO, Approve exceptions
- Edit: Data dengan audit trail
- Cannot: Delete finalized data

**Maintenance Technician:**
- View: Machine status, Maintenance tickets, Spare parts
- Input: Maintenance records, Part usage
- Cannot: View production data detail, Edit production data

**Production Manager:**
- View: All dashboards, All reports, All analytics
- Export: All data
- Cannot: Edit production data (view only)

### Audit Trail
- Semua perubahan data tercatat:
  - Who (user ID & name)
  - What (action: create/update/delete)
  - When (timestamp)
  - Where (IP address, device)
  - Before & after value (untuk update)
- Audit log tidak bisa dihapus
- Retention: 5 tahun

---

## 🚀 Implementation Priority

### Phase 1 (MVP) - 2 Bulan
**Core Production Functionality:**
- [ ] US-CT-001 s/d US-CT-005: Setup & Persiapan Mesin
- [ ] US-CT-006 s/d US-CT-010: Monitoring Produksi Real-Time
- [ ] US-CT-011 s/d US-CT-012: Completion & Handover
- [ ] US-CT-013 s/d US-CT-014: Dashboard Overview & OEE Monitoring
- [ ] US-CT-029: Mobile App - Production Monitoring

**Goal:** Paperless production, real-time tracking, basic OEE calculation

---

### Phase 2 - 1.5 Bulan
**Analytics & Insights:**
- [ ] US-CT-015: Six Big Losses Analysis
- [ ] US-CT-016: Alert & Notification
- [ ] US-CT-017: Production Planning & Scheduling
- [ ] US-CT-018 s/d US-CT-020: Performance Analytics
- [ ] US-CT-031 s/d US-CT-032: Reporting

**Goal:** Data-driven insights, automated alerts, comprehensive reporting

---

### Phase 3 - 1.5 Bulan
**Maintenance & Optimization:**
- [ ] US-CT-021 s/d US-CT-024: Maintenance Management
- [ ] US-CT-034 s/d US-CT-036: Integration (SAP, Khazanah Awal, Maintenance)
- [ ] US-CT-033: Custom Report Builder

**Goal:** Predictive maintenance, seamless integration, custom analytics

---

### Phase 4 - 1 Bulan
**Gamification & Engagement:**
- [ ] US-CT-025 s/d US-CT-028: Gamification (Badges, Points, Leaderboard, Challenges)
- [ ] US-CT-030: Mobile App - Personal Dashboard

**Goal:** Employee engagement, motivation, continuous improvement culture

---

## 📊 Success Metrics

### Adoption
- [ ] 100% operator menggunakan sistem untuk input data
- [ ] < 10 menit training time untuk basic operation
- [ ] User satisfaction score ≥ 4.5/5
- [ ] Mobile app usage ≥ 80% dari total input

### Business Impact
- [ ] **OEE Improvement:** 70% → 85%+ (target +15%)
- [ ] **Setup Time Reduction:** 30 menit → 20 menit (target -33%)
- [ ] **Downtime Reduction:** 15% → 10% (target -33%)
- [ ] **Defect Rate Reduction:** 5% → 2% (target -60%)
- [ ] **Throughput Increase:** +20%
- [ ] **Data Accuracy:** 95% → 99.5%
- [ ] **Reporting Time:** 2 jam/hari → 0 jam (100% automated)

### Operational Excellence
- [ ] Paperless operation: 100%
- [ ] Real-time visibility: 100%
- [ ] Alert response time: < 5 menit
- [ ] Report generation: Automated (0 manual effort)
- [ ] Maintenance response time: < 10 menit

---

## 💡 Best Practices & Recommendations

### Untuk Operator
1. **Update counter setiap 30 menit** - untuk tracking akurat
2. **Record downtime immediately** - jangan tunggu sampai akhir shift
3. **Foto quality issues** - untuk dokumentasi & learning
4. **Review OEE setiap jam** - untuk self-awareness & improvement
5. **Participate in challenges** - untuk fun & motivation

### Untuk Supervisor
1. **Review dashboard setiap jam** - untuk proactive management
2. **Respond to alerts < 5 menit** - untuk minimize impact
3. **Daily coaching** - review performance dengan operator
4. **Weekly improvement meeting** - discuss Six Big Losses
5. **Recognize top performers** - untuk motivation

### Untuk Maintenance
1. **Preventive maintenance on-time** - untuk prevent breakdown
2. **Track MTBF & MTTR** - untuk continuous improvement
3. **Stock critical spare parts** - untuk fast response
4. **Document root cause** - untuk learning & prevention
5. **Predictive maintenance** - shift dari reactive ke proactive

### Untuk Management
1. **Review weekly report** - untuk strategic insights
2. **Data-driven decisions** - use analytics untuk investment decisions
3. **Benchmark best practices** - internal & external
4. **Invest in training** - untuk skill development
5. **Celebrate wins** - recognize improvements & achievements

---

## 🎓 Training Requirements

### Operator Training (4 jam)
1. **Basic Operation (1 jam):**
   - Login & navigation
   - View PO & assignments
   - Setup checklist
   - Start/stop production

2. **Data Input (1 jam):**
   - Update production counter
   - Record downtime
   - Record quality issues
   - Finalize production

3. **Mobile App (1 jam):**
   - Mobile navigation
   - Quick actions
   - Photo upload
   - Offline mode

4. **Gamification (1 jam):**
   - Achievements & badges
   - Points & rewards
   - Leaderboard
   - Challenges

### Supervisor Training (8 jam)
1. **Dashboard & Monitoring (2 jam)**
2. **OEE & Analytics (2 jam)**
3. **Alert Management (1 jam)**
4. **Reporting (1 jam)**
5. **Production Planning (1 jam)**
6. **People Management (1 jam)**

### Maintenance Training (4 jam)
1. **Maintenance Management (2 jam)**
2. **Spare Parts Management (1 jam)**
3. **Predictive Maintenance (1 jam)**

---

**Catatan Akhir:**

User story ini comprehensive dan detail, mencakup semua aspek Unit Cetak dari operator level sampai management level. Fokus utama adalah:

1. **Real-time visibility** - untuk proactive management
2. **OEE optimization** - untuk efisiensi maksimal
3. **Data-driven insights** - untuk continuous improvement
4. **Employee engagement** - untuk motivation & performance
5. **Seamless integration** - untuk end-to-end process flow

Prioritas implementasi bisa disesuaikan dengan kebutuhan bisnis dan resource yang tersedia. Yang paling penting adalah **Phase 1 (MVP)** untuk establish foundation yang solid, baru kemudian build advanced features di phase berikutnya.

**Prinsip desain:** Mobile-first, user-friendly, fast, dan fokus pada **User Experience terbaik untuk operator** karena mereka adalah primary users yang akan menggunakan sistem setiap hari di lapangan.
