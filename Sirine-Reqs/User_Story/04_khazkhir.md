# User Story: Khazanah Akhir (Khazkhir)

## Overview
Khazanah Akhir (Khazkhir) bertanggung jawab untuk penghitungan final HCS hasil verifikasi, pengemasan per rim (500 lembar), penyimpanan di gudang, dan pengiriman ke Dirjen Bea Cukai. Ini adalah tahap terakhir dalam proses produksi sebelum produk dikirim ke customer.

---

## 🎭 Personas

### 1. Staff Khazanah Akhir - Penghitungan & Pengemasan
**Nama:** Rina Kusuma  
**Posisi:** Staff Khazanah Akhir (Counting & Packaging)  
**Shift:** Pagi (07:00-15:00)  
**Pengalaman:** 4 tahun

**Tanggung Jawab:**
- Menghitung final HCS dari hasil verifikasi
- Pengemasan per 1 rim (500 lembar kirim)
- Labeling setiap rim dengan info PO/OBC
- Penyusunan rim di palet untuk storage
- Entry data hasil penghitungan dan pengemasan

**Pain Points:**
- Penghitungan manual memakan waktu dan rawan error
- Sulit tracking rim mana yang sudah dikemas untuk PO mana
- Labeling manual sering salah (nomor rim, tanggal, dll)
- Tidak tahu real-time progress pengemasan
- Laporan manual di akhir shift memakan waktu

**Goals:**
- Penghitungan lebih cepat dan akurat (digital counter)
- Auto-generate label untuk setiap rim
- Real-time tracking progress pengemasan
- Otomasi entry data
- Digital documentation

---

### 2. Staff Khazanah Akhir - Gudang & Inventory
**Nama:** Agung Prasetyo  
**Posisi:** Staff Khazanah Akhir (Warehouse & Inventory)  
**Posisi:** Shift rotasi  
**Pengalaman:** 6 tahun

**Tanggung Jawab:**
- Penyimpanan rim di gudang (rak/lokasi tertentu)
- Inventory management finished goods
- Stock opname
- Picking untuk pengiriman
- Koordinasi dengan tim pengiriman

**Pain Points:**
- Sulit cari lokasi rim tertentu di gudang
- Stock opname manual memakan waktu lama
- Tidak ada visibility real-time stock per PO/OBC
- FIFO tidak terjaga dengan baik
- Kesulitan tracking rim yang sudah dipick untuk pengiriman

**Goals:**
- Digital warehouse management system
- Barcode/QR untuk tracking lokasi
- Real-time inventory visibility
- FIFO automation
- Easy picking dengan digital guidance

---

### 3. Staff Khazanah Akhir - Pengiriman
**Nama:** Budi Hartono  
**Posisi:** Staff Pengiriman  
**Shift:** Pagi (07:00-15:00)  
**Pengalaman:** 5 tahun

**Tanggung Jawab:**
- Persiapan dokumen pengiriman (Surat Jalan, Packing List)
- Loading rim ke kendaraan
- Koordinasi dengan kurir/driver
- Tracking pengiriman
- Konfirmasi delivery ke customer

**Pain Points:**
- Dokumen pengiriman manual (rawan salah)
- Tidak ada visibility real-time status pengiriman
- Sulit tracking proof of delivery
- Customer sering komplain karena tidak tahu status pengiriman
- Laporan pengiriman manual

**Goals:**
- Auto-generate dokumen pengiriman
- Digital proof of delivery
- Real-time tracking pengiriman
- Customer portal untuk tracking
- Automated reporting

---

### 4. Supervisor Khazanah Akhir
**Nama:** Ibu Dewi Lestari  
**Posisi:** Supervisor Khazanah Akhir  
**Shift:** Rotasi  
**Pengalaman:** 10 tahun

**Tanggung Jawab:**
- Monitoring proses penghitungan, pengemasan, penyimpanan, dan pengiriman
- Koordinasi dengan Verifikasi dan customer (Dirjen Bea Cukai)
- Inventory management & stock control
- Delivery performance monitoring
- Reporting ke Logistics Manager

**Pain Points:**
- Tidak ada visibility real-time semua aktivitas
- Harus keliling untuk cek progress
- Sulit tracking on-time delivery performance
- Laporan manual memakan waktu
- Tidak ada alert jika ada issue (stock shortage, delay, dll)

**Goals:**
- Dashboard real-time untuk monitoring semua aktivitas
- Alert otomatis jika ada issue
- Automated reporting
- Data-driven decision making
- Improve on-time delivery rate

---

## 📱 User Stories

### Epic 1: Penghitungan & Pengemasan

#### US-KA-001: Melihat Daftar PO yang Perlu Dikemas
**Sebagai** Staff Khazanah Akhir  
**Saya ingin** melihat daftar PO yang sudah selesai verifikasi dan siap untuk dikemas  
**Sehingga** saya tahu prioritas pekerjaan hari ini

**Acceptance Criteria:**
- [ ] Tampil daftar PO dengan status "Selesai Verifikasi - Siap Pengemasan"
- [ ] Sorting berdasarkan prioritas dan due date
- [ ] Tampil info per PO:
  - Nomor PO & OBC
  - Total HCS (lembar kirim)
  - Estimasi jumlah rim (HCS ÷ 500)
  - Due date pengiriman
  - Prioritas (🔴 Urgent / 🟡 Normal / 🟢 Low)
  - Waktu selesai verifikasi
  - Durasi menunggu pengemasan
  - Lokasi penyimpanan HCS (dari verifikasi)
  - Customer: Dirjen Bea Cukai
- [ ] Filter berdasarkan: prioritas, tanggal, status
- [ ] Search berdasarkan nomor PO/OBC
- [ ] Alert jika PO menunggu > 4 jam
- [ ] Responsive untuk tablet/mobile

**Business Rules:**
- PO dengan due date < 2 hari = Urgent (🔴)
- PO dengan due date 2-5 hari = Normal (🟡)
- PO dengan due date > 5 hari = Low (🟢)
- FIFO: PO yang selesai verifikasi lebih dulu harus dikemas lebih dulu
- Alert jika waiting time > 4 jam

---

#### US-KA-002: Memulai Proses Penghitungan Final
**Sebagai** Staff Khazanah Akhir  
**Saya ingin** memulai proses penghitungan final HCS  
**Sehingga** sistem tracking progress pekerjaan saya

**Acceptance Criteria:**
- [ ] Button "Mulai Penghitungan" pada detail PO
- [ ] Sistem tampilkan:
  - Target HCS (dari verifikasi): [X] lembar
  - Lokasi penyimpanan HCS: [Rak/Palet]
  - Estimasi jumlah rim: [Y] rim ([X] ÷ 500)
  - Estimasi durasi: [N] jam
- [ ] Status PO berubah menjadi "Sedang Dihitung - Khazanah Akhir"
- [ ] Timestamp mulai penghitungan tercatat
- [ ] Nama staff yang handle tercatat
- [ ] Timer mulai berjalan (real-time)

**Business Rules:**
- 1 Staff bisa handle 1 PO dalam satu waktu
- Estimasi durasi = (Total HCS ÷ 500) × 15 menit per rim
- Target: 15 menit per rim (counting + packaging)

---

#### US-KA-003: Input Hasil Penghitungan Final
**Sebagai** Staff Khazanah Akhir  
**Saya ingin** input hasil penghitungan final HCS  
**Sehingga** ada record akurat sebelum pengemasan

**Acceptance Criteria:**
- [ ] Input field:
  - Jumlah HCS Aktual (lembar kirim)
  - Jumlah Rusak/Reject (jika ada kerusakan saat handling)
- [ ] Sistem otomatis hitung:
  - Target (dari verifikasi): [X] lembar
  - Aktual: [Y] lembar
  - Selisih: [Y-X] lembar
  - Persentase Akurasi: ([Y] ÷ [X]) × 100%
  - Estimasi Rim: [Y] ÷ 500
- [ ] Jika selisih > 1%, tampil warning dan wajib isi alasan
- [ ] Jika rusak > 0, wajib input:
  - Jenis kerusakan
  - Foto kerusakan
  - Alasan (handling damage, dll)
- [ ] Real-time validation input (tidak boleh negatif)
- [ ] Timestamp input tercatat

**Business Rules:**
- Toleransi selisih: ±1%
- Jika selisih > 1%, wajib isi alasan dan approval supervisor
- Jika rusak > 0.5%, escalate ke supervisor

---

#### US-KA-004: Pengemasan per Rim dengan Auto-Label
**Sebagai** Staff Khazanah Akhir  
**Saya ingin** sistem auto-generate label untuk setiap rim  
**Sehingga** labeling akurat dan cepat

**Acceptance Criteria:**
- [ ] Setelah input hasil penghitungan, sistem tampilkan:
  - Total rim yang akan dikemas: [X] rim
  - Sisa lembar (jika tidak genap 500): [Y] lembar
- [ ] Button "Mulai Pengemasan"
- [ ] Untuk setiap rim, sistem auto-generate label berisi:
  - **Nomor PO:** PO-XXXX
  - **Nomor OBC:** OBC-YYYY-ZZZ
  - **Nomor Rim:** 1/[Total], 2/[Total], ..., [Total]/[Total]
  - **Jumlah:** 500 lembar (atau sisa untuk rim terakhir)
  - **Tanggal Produksi:** DD/MM/YYYY
  - **Tanggal Verifikasi:** DD/MM/YYYY
  - **Tanggal Pengemasan:** DD/MM/YYYY
  - **Status:** HCS (Hasil Cetak Sempurna)
  - **QR Code:** Untuk tracking (encode: PO, OBC, Rim Number)
  - **Barcode:** Untuk scanning
  - **Customer:** Dirjen Bea Cukai
- [ ] Print label via thermal printer (atau export PDF untuk print)
- [ ] Checkbox "Rim [X] sudah dikemas" untuk setiap rim
- [ ] Upload foto rim yang sudah dikemas (opsional)
- [ ] Real-time progress: [X]/[Total] rim selesai dikemas

**Business Rules:**
- 1 Rim = 500 lembar kirim (standar)
- Rim terakhir bisa < 500 lembar (sisa)
- Label wajib ada QR Code dan Barcode untuk tracking
- Label format sesuai standar Dirjen Bea Cukai

---

#### US-KA-005: Penyusunan Rim di Palet
**Sebagai** Staff Khazanah Akhir  
**Saya ingin** tracking rim mana yang disusun di palet mana  
**Sehingga** mudah untuk storage dan picking

**Acceptance Criteria:**
- [ ] Setelah semua rim dikemas, button "Susun di Palet"
- [ ] Input/scan nomor palet (auto-generate atau manual)
- [ ] Sistem tampilkan:
  - Total rim: [X] rim
  - Estimasi jumlah palet (1 palet max 20 rim)
  - Palet yang akan digunakan
- [ ] Untuk setiap palet:
  - Nomor Palet: PAL-KA-XXX
  - Jumlah rim di palet ini: [Y] rim
  - Rim numbers: 1-20, 21-40, dst.
  - QR Code palet (untuk tracking)
- [ ] Print label palet berisi:
  - Nomor Palet
  - Nomor PO & OBC
  - Total rim di palet ini
  - Rim range (e.g., Rim 1-20 dari 59 rim)
  - Tanggal pengemasan
  - QR Code
- [ ] Upload foto palet yang sudah disusun
- [ ] Timestamp penyusunan tercatat

**Business Rules:**
- 1 Palet maksimal 20 rim (untuk safety & handling)
- 1 Palet hanya untuk 1 PO (tidak boleh mix)
- Palet harus di-label dengan jelas

---

#### US-KA-006: Finalisasi Pengemasan
**Sebagai** Staff Khazanah Akhir  
**Saya ingin** finalisasi proses pengemasan  
**Sehingga** rim siap untuk disimpan di gudang

**Acceptance Criteria:**
- [ ] Button "Selesai Pengemasan"
- [ ] Sistem validasi:
  - ✅ Semua rim sudah dikemas
  - ✅ Semua rim sudah di-label
  - ✅ Semua rim sudah disusun di palet
  - ✅ Foto palet sudah diupload (opsional tapi recommended)
- [ ] Sistem tampilkan summary:
  - Total HCS: [X] lembar
  - Total Rim: [Y] rim
  - Total Palet: [Z] palet
  - Durasi pengemasan: [N] jam [M] menit
  - Efficiency: [Target vs Actual]
- [ ] Konfirmasi finalisasi
- [ ] Status PO berubah "Selesai Pengemasan - Siap Simpan"
- [ ] Notifikasi ke Staff Gudang untuk penyimpanan
- [ ] Timestamp selesai pengemasan tercatat
- [ ] Sistem hitung durasi pengemasan
- [ ] Data tidak bisa diubah setelah finalisasi

**Business Rules:**
- Semua validasi harus pass sebelum bisa finalisasi
- Durasi Pengemasan = Timestamp Selesai - Timestamp Mulai
- Target efficiency: ≤ 15 menit per rim
- Data tidak bisa diubah setelah finalisasi (untuk audit)

---

### Epic 2: Warehouse & Inventory Management

#### US-KA-007: Penyimpanan Rim di Gudang
**Sebagai** Staff Gudang  
**Saya ingin** simpan rim di lokasi gudang yang tepat  
**Sehingga** mudah untuk picking nanti

**Acceptance Criteria:**
- [ ] Tampil daftar palet yang "Siap Simpan"
- [ ] Button "Simpan di Gudang" pada detail palet
- [ ] Sistem suggest lokasi optimal berdasarkan:
  - FIFO (due date paling dekat di lokasi paling accessible)
  - Kapasitas rak
  - Proximity (PO yang sama di lokasi berdekatan)
- [ ] Dropdown pilih lokasi penyimpanan:
  - Zona: A, B, C, D (berdasarkan area gudang)
  - Rak: 1-50
  - Level: 1-5 (tingkat rak)
  - Posisi: A-Z (posisi di rak)
  - Format: A-10-3-C (Zona A, Rak 10, Level 3, Posisi C)
- [ ] Scan QR Code palet untuk konfirmasi
- [ ] Scan QR Code lokasi untuk konfirmasi
- [ ] Upload foto palet di lokasi (opsional)
- [ ] Timestamp penyimpanan tercatat
- [ ] Status palet berubah "Tersimpan di Gudang"
- [ ] Lokasi palet terupdate di sistem
- [ ] Inventory stock terupdate

**Business Rules:**
- FIFO: PO dengan due date paling dekat disimpan di lokasi paling accessible
- 1 Lokasi bisa multiple palet (jika cukup space)
- Scan QR wajib untuk memastikan lokasi yang benar
- Real-time inventory update

---

#### US-KA-008: Real-Time Inventory Dashboard
**Sebagai** Staff Gudang / Supervisor  
**Saya ingin** melihat inventory real-time  
**Sehingga** tahu stock availability setiap saat

**Acceptance Criteria:**
- [ ] Dashboard Inventory dengan metrics:
  
  **Overall Stock:**
  - Total Rim di Gudang: [X] rim ([Y] lembar)
  - Total Palet: [Z] palet
  - Kapasitas Gudang: [W]% terpakai
  - Available Space: [N] lokasi kosong
  
  **Stock per PO/OBC:**
  - List PO/OBC dengan stock:
    - Nomor PO & OBC
    - Total rim
    - Total lembar
    - Jumlah palet
    - Lokasi penyimpanan
    - Due date
    - Status (Ready to Ship / Waiting / On Hold)
  
  **Stock Aging:**
  - < 7 hari: [X] rim
  - 7-14 hari: [Y] rim
  - 14-30 hari: [Z] rim
  - > 30 hari: [W] rim (alert jika ada)
  
  **Warehouse Utilization:**
  - Heatmap lokasi gudang (occupied vs available)
  - Kapasitas per zona
  - Optimal vs actual layout

- [ ] Search & Filter:
  - Search by PO/OBC
  - Filter by status, due date, zona
  - Sort by FIFO, quantity, location
- [ ] Export to Excel
- [ ] Auto-refresh setiap 1 menit

**Business Rules:**
- Real-time update berdasarkan penyimpanan dan picking
- Alert jika stock aging > 30 hari (slow-moving)
- Kapasitas gudang warning jika > 90%

---

#### US-KA-009: Stock Opname Digital
**Sebagai** Staff Gudang  
**Saya ingin** melakukan stock opname dengan digital  
**Sehingga** lebih cepat dan akurat

**Acceptance Criteria:**
- [ ] Button "Mulai Stock Opname"
- [ ] Pilih scope:
  - Full stock opname (seluruh gudang)
  - Partial stock opname (per zona/rak)
  - Cycle count (per PO/OBC)
- [ ] Sistem generate checklist:
  - List semua palet yang harus di-count
  - Expected quantity (dari sistem)
  - Lokasi penyimpanan
- [ ] Mobile app untuk counting:
  - Scan QR Code palet
  - Input actual quantity (jumlah rim)
  - Foto palet (jika ada discrepancy)
  - Notes (jika ada issue)
- [ ] Real-time progress: [X]/[Total] palet counted
- [ ] Sistem otomatis hitung discrepancy:
  - Expected vs Actual
  - Variance (quantity & %)
  - Variance value (Rp)
- [ ] Jika discrepancy > 1%, wajib investigation
- [ ] Finalisasi stock opname:
  - Summary report
  - Discrepancy list
  - Action items
  - Approval supervisor
- [ ] Inventory adjustment (jika approved)

**Business Rules:**
- Full stock opname: minimal 1× per bulan
- Cycle count: minimal 1× per minggu (rotating PO)
- Discrepancy > 1% wajib investigation & approval
- Inventory adjustment hanya setelah approval

---

#### US-KA-010: Warehouse Location Finder
**Sebagai** Staff Gudang  
**Saya ingin** mudah cari lokasi rim tertentu  
**Sehingga** picking lebih cepat

**Acceptance Criteria:**
- [ ] Search bar: input PO/OBC/Palet number
- [ ] Sistem tampilkan:
  - Lokasi penyimpanan (Zona-Rak-Level-Posisi)
  - Visual map gudang (highlight lokasi)
  - Navigasi ke lokasi (step-by-step)
  - Info rim:
    - Nomor PO & OBC
    - Jumlah rim di lokasi ini
    - Rim numbers
    - Tanggal simpan
    - Due date
- [ ] QR Code scanner untuk quick search
- [ ] Filter nearby locations (untuk batch picking)
- [ ] Export location list (untuk picking list)

**Business Rules:**
- Search harus < 2 detik
- Visual map untuk easy navigation
- Support batch picking (multiple PO sekaligus)

---

#### US-KA-011: FIFO Compliance Monitoring
**Sebagai** Supervisor Khazanah Akhir  
**Saya ingin** monitor FIFO compliance  
**Sehingga** tidak ada stock yang terlalu lama di gudang

**Acceptance Criteria:**
- [ ] FIFO Dashboard:
  
  **FIFO Compliance:**
  - FIFO Compliance Rate: [X]%
  - Target: ≥ 98%
  - Violations: [Y] cases
  
  **Aging Analysis:**
  - Stock aging distribution (chart)
  - Slow-moving items (> 30 hari)
  - At-risk items (approaching due date)
  
  **Alerts:**
  - 🔴 Stock > 30 hari (urgent action needed)
  - 🟡 Stock 21-30 hari (monitor closely)
  - 🔵 Stock < 7 hari dari due date (prioritize shipping)

- [ ] Drill-down per PO/OBC
- [ ] Action items tracking
- [ ] Export to Excel

**Business Rules:**
- FIFO: First In First Out (yang masuk dulu harus keluar dulu)
- Target FIFO compliance: ≥ 98%
- Alert jika ada stock > 30 hari
- Prioritize shipping untuk stock yang mendekati due date

---

### Epic 3: Pengiriman (Shipping & Delivery)

#### US-KA-012: Melihat Daftar PO yang Siap Kirim
**Sebagai** Staff Pengiriman  
**Saya ingin** melihat daftar PO yang siap untuk dikirim  
**Sehingga** bisa planning pengiriman

**Acceptance Criteria:**
- [ ] Tampil daftar PO dengan status "Tersimpan di Gudang - Siap Kirim"
- [ ] Sorting berdasarkan due date (FIFO)
- [ ] Tampil info per PO:
  - Nomor PO & OBC
  - Total rim
  - Total lembar
  - Jumlah palet
  - Lokasi penyimpanan (list)
  - Due date pengiriman
  - Prioritas (🔴 Urgent / 🟡 Normal / 🟢 Low)
  - Customer: Dirjen Bea Cukai
  - Alamat pengiriman
  - Contact person
- [ ] Filter berdasarkan: prioritas, due date, status
- [ ] Search berdasarkan nomor PO/OBC
- [ ] Alert jika due date < 2 hari
- [ ] Batch selection (untuk pengiriman multiple PO sekaligus)

**Business Rules:**
- PO dengan due date paling dekat harus dikirim dulu (FIFO)
- Alert jika due date < 2 hari
- Support batch shipping (multiple PO dalam 1 truk jika tujuan sama)

---

#### US-KA-013: Create Delivery Order
**Sebagai** Staff Pengiriman  
**Saya ingin** create delivery order  
**Sehingga** ada dokumen resmi untuk pengiriman

**Acceptance Criteria:**
- [ ] Button "Buat Delivery Order" pada detail PO
- [ ] Form Delivery Order:
  
  **Header:**
  - Nomor DO: DO-YYYY-MM-XXXX (auto-generate)
  - Tanggal DO: [Today]
  - Nomor PO & OBC
  - Customer: Dirjen Bea Cukai
  
  **Shipping Details:**
  - Alamat pengiriman (auto-fill dari master data)
  - Contact person (auto-fill)
  - Phone number (auto-fill)
  - Tanggal pengiriman (input)
  - Estimasi waktu tiba (input)
  
  **Items:**
  - Total rim: [X] rim
  - Total lembar: [Y] lembar
  - Total palet: [Z] palet
  - Lokasi penyimpanan (list)
  
  **Logistics:**
  - Jenis kendaraan (dropdown: Truk, Pickup, dll)
  - Nomor kendaraan (input)
  - Driver name (input)
  - Driver phone (input)
  - Ekspedisi (jika pakai 3rd party)
  
  **Documents:**
  - Surat Jalan (auto-generate)
  - Packing List (auto-generate)
  - Invoice (jika applicable)
  - BAST (Berita Acara Serah Terima) template

- [ ] Preview dokumen sebelum finalize
- [ ] Edit jika ada yang perlu diubah
- [ ] Button "Finalize DO"
- [ ] Status PO berubah "DO Created - Ready to Pick"
- [ ] Notifikasi ke Staff Gudang untuk picking

**Business Rules:**
- DO number auto-generate (format: DO-YYYY-MM-XXXX)
- Semua field wajib diisi
- Dokumen auto-generate berdasarkan template
- DO tidak bisa diubah setelah finalize (hanya bisa cancel & create new)

---

#### US-KA-014: Picking Process dengan Digital Guidance
**Sebagai** Staff Gudang  
**Saya ingin** digital guidance saat picking  
**Sehingga** picking lebih cepat dan akurat

**Acceptance Criteria:**
- [ ] Tampil daftar DO yang "Ready to Pick"
- [ ] Button "Mulai Picking" pada detail DO
- [ ] Sistem generate picking list:
  - List semua palet yang harus dipick
  - Lokasi penyimpanan (urut berdasarkan route optimal)
  - Jumlah rim per palet
  - Visual map gudang (route guidance)
- [ ] Mobile app untuk picking:
  - Step-by-step guidance:
    - "Pergi ke lokasi: A-10-3-C"
    - Visual map dengan highlight
    - "Ambil palet: PAL-KA-001"
    - "Scan QR Code palet untuk konfirmasi"
  - Checkbox setiap palet yang sudah dipick
  - Real-time progress: [X]/[Total] palet picked
  - Upload foto palet yang sudah dipick (opsional)
- [ ] Sistem validasi:
  - QR Code palet harus match dengan picking list
  - Jika salah palet, tampil error & guidance
- [ ] Setelah semua palet dipick:
  - Button "Selesai Picking"
  - Konfirmasi semua palet sudah di staging area
  - Timestamp selesai picking tercatat
  - Status DO berubah "Picked - Ready to Load"
  - Inventory stock terupdate (kurangi stock)

**Business Rules:**
- Route optimization untuk minimize walking distance
- QR Code scan wajib untuk setiap palet (untuk akurasi)
- Inventory update real-time saat picking
- Picking time target: ≤ 5 menit per palet

---

#### US-KA-015: Loading & Dispatch
**Sebagai** Staff Pengiriman  
**Saya ingin** track loading process dan dispatch  
**Sehingga** ada record lengkap

**Acceptance Criteria:**
- [ ] Tampil daftar DO yang "Ready to Load"
- [ ] Button "Mulai Loading" pada detail DO
- [ ] Loading checklist:
  - ✅ Semua palet sudah di staging area
  - ✅ Kendaraan sudah siap
  - ✅ Driver sudah hadir
  - ✅ Dokumen lengkap (Surat Jalan, Packing List, BAST)
- [ ] Untuk setiap palet:
  - Scan QR Code palet saat loading
  - Checkbox "Palet [X] sudah diload"
  - Real-time progress: [X]/[Total] palet loaded
- [ ] Upload foto:
  - Foto palet di kendaraan
  - Foto kendaraan (nomor polisi visible)
  - Foto driver
- [ ] Button "Selesai Loading"
- [ ] Konfirmasi dispatch:
  - Timestamp dispatch
  - Driver signature (digital)
  - Staff signature (digital)
  - Odometer reading (input)
- [ ] Print/Email dokumen ke driver:
  - Surat Jalan
  - Packing List
  - BAST template
  - Delivery instructions
- [ ] Status DO berubah "Dispatched - In Transit"
- [ ] Notifikasi ke customer (email/WhatsApp):
  - "Pesanan Anda sedang dalam perjalanan"
  - Nomor DO
  - Estimasi waktu tiba
  - Tracking link
- [ ] SMS/WhatsApp ke driver dengan delivery details

**Business Rules:**
- Semua palet wajib di-scan saat loading (untuk akurasi)
- Foto wajib untuk dokumentasi
- Digital signature untuk accountability
- Customer notification otomatis

---

#### US-KA-016: Delivery Tracking
**Sebagai** Staff Pengiriman / Supervisor  
**Saya ingin** track status pengiriman real-time  
**Sehingga** tahu posisi dan ETA

**Acceptance Criteria:**
- [ ] Delivery Tracking Dashboard:
  
  **Active Deliveries:**
  - List DO yang "In Transit"
  - Untuk setiap DO:
    - Nomor DO
    - Nomor PO & OBC
    - Customer
    - Driver name & phone
    - Kendaraan
    - Dispatch time
    - ETA
    - Status (On Time / Delayed)
    - GPS location (jika ada GPS tracker)
  
  **Delivery Status:**
  - 🟢 On Time: ETA sesuai schedule
  - 🟡 At Risk: Mendekati late
  - 🔴 Delayed: Sudah late
  
  **Map View:**
  - Visual map dengan pin lokasi kendaraan (jika ada GPS)
  - Route dari warehouse ke customer
  - Real-time location update

- [ ] Click DO untuk detail tracking:
  - Timeline:
    - DO Created: [Timestamp]
    - Picking Started: [Timestamp]
    - Picking Completed: [Timestamp]
    - Loading Started: [Timestamp]
    - Dispatched: [Timestamp]
    - In Transit: [Current]
    - Delivered: [Pending]
  - Driver contact (call/WhatsApp button)
  - Update ETA (jika ada perubahan)
  - Notes/Issues (jika ada)

- [ ] Alert:
  - 🔴 Delivery delayed > 2 jam
  - 🟡 Delivery at risk (approaching late)
  - 🔵 Delivery completed

**Business Rules:**
- Real-time tracking jika ada GPS tracker
- Manual update dari driver jika tidak ada GPS
- Alert otomatis untuk delayed delivery
- ETA update berdasarkan traffic & GPS data

---

#### US-KA-017: Proof of Delivery (POD)
**Sebagai** Driver / Staff Pengiriman  
**Saya ingin** digital proof of delivery  
**Sehingga** ada bukti resmi pengiriman

**Acceptance Criteria:**
- [ ] Mobile app untuk driver (saat sampai di customer):
  
  **Delivery Confirmation:**
  - Button "Sampai di Lokasi"
  - Timestamp arrival tercatat
  - GPS location captured
  
  **Unloading:**
  - Checkbox "Unloading selesai"
  - Upload foto:
    - Foto palet yang sudah diunload
    - Foto lokasi customer
  
  **BAST (Berita Acara Serah Terima):**
  - Digital BAST form:
    - Nomor DO
    - Nomor PO & OBC
    - Tanggal terima
    - Jumlah rim diterima
    - Kondisi barang (Baik/Rusak)
    - Jika rusak: foto & deskripsi
  - Customer signature (digital - touchscreen)
  - Customer name (input)
  - Customer ID/NIK (input)
  - Timestamp signature
  - Driver signature (digital)
  
  **Finalize Delivery:**
  - Button "Selesai Delivery"
  - Konfirmasi semua dokumen lengkap
  - Status DO berubah "Delivered"
  - Timestamp delivery tercatat
  - Auto-send POD ke office (PDF)

- [ ] Office staff bisa view POD:
  - Download PDF BAST
  - View photos
  - View signatures
  - Export untuk archive

- [ ] Customer notification:
  - Email/WhatsApp: "Pesanan Anda sudah diterima"
  - Attach POD (PDF)
  - Thank you message

**Business Rules:**
- Customer signature wajib untuk POD
- Foto wajib untuk dokumentasi
- POD tidak bisa diubah setelah finalize
- Auto-archive POD untuk audit (retention 5 tahun)

---

#### US-KA-018: Delivery Performance Dashboard
**Sebagai** Supervisor Khazanah Akhir  
**Saya ingin** monitor delivery performance  
**Sehingga** bisa improve on-time delivery rate

**Acceptance Criteria:**
- [ ] Delivery Performance Dashboard:
  
  **Key Metrics:**
  - **On-Time Delivery (OTD) Rate:** [X]%
    - Target: ≥ 95%
    - Trend: 📈/📉 (vs bulan lalu)
  - **Total Deliveries:** [Y] DO
    - On Time: [A] DO
    - Late: [B] DO
    - In Transit: [C] DO
  - **Average Delivery Time:** [Z] jam
    - From dispatch to delivery
    - Comparison vs target
  - **Perfect Delivery Rate:** [W]%
    - On time + No damage + Complete
  
  **Late Delivery Analysis:**
  - Late deliveries: [X] DO
  - Average late time: [Y] jam
  - Late reasons (Pareto chart):
    - Traffic
    - Vehicle breakdown
    - Customer not ready
    - Wrong address
    - Weather
    - Lainnya
  - Cost of late delivery (penalty, dll)
  
  **Delivery Trend:**
  - Line chart: OTD rate per minggu/bulan
  - Comparison: This month vs Last month
  - Seasonality pattern
  
  **Driver Performance:**
  - Top 5 drivers (by OTD rate)
  - Bottom 5 drivers (need improvement)
  - Average delivery time per driver

- [ ] Filter: periode, driver, customer
- [ ] Drill-down per DO
- [ ] Export to Excel/PDF

**Business Rules:**
- OTD = delivered on or before due date
- Target OTD: ≥ 95%
- Alert jika OTD < 90%
- Monthly review untuk continuous improvement

---

### Epic 4: Customer Portal & Communication

#### US-KA-019: Customer Portal untuk Dirjen Bea Cukai
**Sebagai** Customer (Dirjen Bea Cukai)  
**Saya ingin** portal untuk track order saya  
**Sehingga** tidak perlu telepon untuk tanya status

**Acceptance Criteria:**
- [ ] Customer Portal (web-based):
  
  **Login:**
  - Username & password (provided by company)
  - Role: Customer - Dirjen Bea Cukai
  
  **Dashboard:**
  - **Active Orders:**
    - List PO/OBC yang sedang dalam proses
    - Status per PO:
      - In Production
      - In Quality Control
      - In Packaging
      - Ready to Ship
      - In Transit
      - Delivered
    - Progress bar per PO
    - ETA per stage
  
  - **Order History:**
    - List PO/OBC yang sudah delivered
    - Delivery date
    - POD download
  
  - **Order Details:**
    - Click PO untuk detail:
      - Nomor PO & OBC
      - Spesifikasi produk
      - Jumlah order
      - Jumlah produced
      - Jumlah HCS
      - Jumlah HCTS
      - Quality metrics
      - Timeline (Gantt chart)
      - Current status
      - ETA delivery
  
  - **Delivery Tracking:**
    - Real-time tracking untuk order yang "In Transit"
    - Map view (jika ada GPS)
    - Driver contact
    - ETA
  
  - **Documents:**
    - Download documents:
      - Surat Jalan
      - Packing List
      - BAST (POD)
      - Invoice
      - Certificate of Quality

- [ ] Notifications:
  - Email notification untuk status changes
  - WhatsApp notification (optional)

- [ ] Support:
  - Contact form untuk inquiry
  - Live chat (jika available)

**Business Rules:**
- Customer hanya bisa lihat order mereka sendiri
- Real-time status update
- Document download available setelah delivery
- Secure access (HTTPS, authentication)

---

#### US-KA-020: Automated Customer Communication
**Sebagai** System  
**Saya ingin** otomatis kirim update ke customer  
**Sehingga** customer selalu informed

**Acceptance Criteria:**
- [ ] Auto-notification triggers:
  
  **Production Milestones:**
  - PO Created: "Order Anda telah diterima dan sedang diproses"
  - Production Started: "Produksi order Anda telah dimulai"
  - Production Completed: "Produksi order Anda telah selesai"
  - QC Completed: "Order Anda telah lulus Quality Control"
  - Packaging Completed: "Order Anda telah dikemas dan siap kirim"
  
  **Shipping Updates:**
  - Ready to Ship: "Order Anda siap untuk dikirim"
  - Dispatched: "Order Anda sedang dalam perjalanan. ETA: [Date/Time]"
  - Out for Delivery: "Order Anda akan tiba hari ini"
  - Delivered: "Order Anda telah diterima. Terima kasih!"
  
  **Issues/Delays:**
  - Quality Issue: "Ada issue pada order Anda. Tim kami sedang menangani."
  - Delivery Delay: "Pengiriman order Anda mengalami delay. ETA baru: [Date/Time]"

- [ ] Notification channels:
  - Email (primary)
  - WhatsApp (optional)
  - SMS (optional)
  - Portal notification (in-app)

- [ ] Notification content:
  - Subject/Title
  - Message body (clear & concise)
  - Order details (PO/OBC, quantity, ETA)
  - Action link (track order, view details)
  - Contact info (jika ada pertanyaan)

- [ ] Notification preferences:
  - Customer bisa set preference (email only, email+WhatsApp, dll)
  - Frequency (real-time, daily digest, dll)

**Business Rules:**
- Notification real-time (< 5 menit setelah event)
- Professional tone (Bahasa Indonesia formal)
- Clear & actionable
- Unsubscribe option (untuk non-critical notifications)

---

### Epic 5: Dashboard & Analytics

#### US-KA-021: Dashboard Overview Khazanah Akhir
**Sebagai** Supervisor Khazanah Akhir  
**Saya ingin** melihat overview semua aktivitas Khazanah Akhir  
**Sehingga** saya punya visibility real-time progress pekerjaan

**Acceptance Criteria:**
- [ ] Tampil di dashboard:
  
  **Pengemasan:**
  - Menunggu: [X] PO ([Y] lembar)
  - Sedang Dikerjakan: [Z] PO
  - Selesai Hari Ini: [W] PO ([A] rim)
  - Rata-rata Durasi: [N] menit per rim
  - Efficiency: [B]% (vs target 15 menit/rim)
  
  **Inventory:**
  - Total Stock: [X] rim ([Y] lembar)
  - Total Palet: [Z] palet
  - Kapasitas Gudang: [W]%
  - Stock Aging:
    - < 7 hari: [A] rim
    - 7-14 hari: [B] rim
    - 14-30 hari: [C] rim
    - > 30 hari: [D] rim (alert jika > 0)
  
  **Pengiriman:**
  - Ready to Ship: [X] PO
  - In Transit: [Y] DO
  - Delivered Today: [Z] DO
  - On-Time Delivery Rate: [W]%
  - Late Deliveries: [A] DO
  
  **Quality Metrics:**
  - Packaging Accuracy: [X]%
  - Picking Accuracy: [Y]%
  - Delivery Accuracy: [Z]% (no damage, complete)

- [ ] Chart: OTD trend (7 hari terakhir)
- [ ] Chart: Inventory level trend (30 hari terakhir)
- [ ] Chart: Packaging efficiency trend (7 hari terakhir)
- [ ] Alert section (jika ada issue)
- [ ] Auto-refresh setiap 30 detik

**Business Rules:**
- Dashboard harus load < 2 detik
- Real-time update berdasarkan aktivitas
- Alert prioritization: Red > Yellow > Blue

---

#### US-KA-022: Staff Performance Monitoring
**Sebagai** Supervisor Khazanah Akhir  
**Saya ingin** melihat performa individual staff  
**Sehingga** bisa evaluasi dan coaching

**Acceptance Criteria:**
- [ ] Tampil daftar staff Khazanah Akhir
- [ ] Per staff tampil metrics:
  
  **Pengemasan:**
  - Jumlah PO selesai hari ini
  - Jumlah rim dikemas hari ini
  - Rata-rata durasi per rim
  - Efficiency (vs target 15 menit/rim)
  - Accuracy (label, palet, dll)
  
  **Gudang:**
  - Jumlah palet disimpan hari ini
  - Jumlah palet dipick hari ini
  - Picking accuracy
  - Picking speed (menit per palet)
  
  **Pengiriman:**
  - Jumlah DO processed hari ini
  - On-time dispatch rate
  - Document accuracy
  - Customer satisfaction (jika ada feedback)

- [ ] Comparison dengan target dan rata-rata tim
- [ ] Filter berdasarkan: hari ini, minggu ini, bulan ini
- [ ] Export to Excel

**Business Rules:**
- Fair comparison (normalize by working hours & workload)
- Transparent metrics
- Focus on coaching & improvement (not punishment)

---

#### US-KA-023: Alert & Notification untuk Supervisor
**Sebagai** Supervisor Khazanah Akhir  
**Saya ingin** menerima alert jika ada issue  
**Sehingga** bisa cepat handling masalah

**Acceptance Criteria:**
- [ ] Alert otomatis untuk kondisi:
  
  **🔴 Critical:**
  - PO mendekati due date (< 2 hari) tapi belum dikemas
  - Stock aging > 30 hari
  - Delivery delayed > 4 jam
  - Inventory discrepancy > 5%
  - Kapasitas gudang > 95%
  
  **🟡 Warning:**
  - PO menunggu pengemasan > 4 jam
  - Stock aging 21-30 hari
  - Delivery at risk (approaching late)
  - Inventory discrepancy 2-5%
  - Kapasitas gudang 85-95%
  - OTD rate < 95%
  
  **🔵 Info:**
  - PO selesai dikemas
  - DO dispatched
  - Delivery completed
  - Stock opname scheduled

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

#### US-KA-024: Laporan Harian Khazanah Akhir
**Sebagai** Supervisor Khazanah Akhir  
**Saya ingin** generate laporan harian otomatis  
**Sehingga** tidak perlu buat laporan manual

**Acceptance Criteria:**
- [ ] Auto-generate laporan di akhir shift
- [ ] Isi laporan:
  
  **Summary Aktivitas:**
  - Pengemasan:
    - Jumlah PO selesai dikemas
    - Total rim dikemas
    - Rata-rata durasi per rim
    - Efficiency rate
  - Inventory:
    - Stock beginning of day
    - Stock in (dari pengemasan)
    - Stock out (dari pengiriman)
    - Stock end of day
    - Stock aging summary
  - Pengiriman:
    - Jumlah DO dispatched
    - Jumlah DO delivered
    - On-time delivery rate
    - Late deliveries & reasons
  
  **Performance Metrics:**
  - Packaging efficiency
  - Picking accuracy
  - Delivery performance
  - Staff productivity summary
  
  **Issues & Actions:**
  - Active issues
  - Resolved issues
  - Pending actions
  - Recommendations

- [ ] Format: PDF & Excel
- [ ] Auto-send via email ke Logistics Manager
- [ ] Bisa download manual kapan saja

**Business Rules:**
- Report generated automatically (no manual intervention)
- Data accurate & real-time
- Professional format (ready untuk management)

---

### Epic 6: Integration & Automation

#### US-KA-025: Integrasi dengan Verifikasi
**Sebagai** System  
**Saya ingin** seamless integration dengan Verifikasi  
**Sehingga** handover HCS smooth

**Acceptance Criteria:**
- [ ] **Verifikasi Selesai Notification:**
  - Tim Verifikasi finalisasi verifikasi
  - → Auto-update status PO di Khazanah Akhir "Selesai Verifikasi - Siap Pengemasan"
  - → Notifikasi ke Staff Khazanah Akhir
  - → HCS details available:
    - Total HCS (lembar)
    - Lokasi penyimpanan HCS (dari verifikasi)
    - Quality summary
    - Estimasi rim

- [ ] Real-time status sync
- [ ] Data consistency check
- [ ] Audit trail

**Business Rules:**
- Status sync real-time (< 5 detik)
- No manual intervention needed

---

#### US-KA-026: Integrasi dengan SAP - Goods Movement
**Sebagai** System  
**Saya ingin** otomatis update goods movement ke SAP  
**Sehingga** SAP inventory selalu accurate

**Acceptance Criteria:**
- [ ] **Setiap finalisasi pengemasan:**
  - Sistem otomatis kirim data ke SAP:
    - PO number
    - Goods receipt (finished goods)
    - Quantity (rim & lembar)
    - Storage location (gudang)
    - Timestamp
    - Batch number

- [ ] **Setiap picking (dispatch):**
  - Sistem otomatis kirim data ke SAP:
    - Goods issue (dari gudang)
    - Delivery order number
    - Quantity shipped
    - Customer
    - Timestamp

- [ ] **Setiap delivery completed:**
  - Sistem otomatis kirim data ke SAP:
    - Delivery confirmation
    - POD data
    - Billing trigger (jika applicable)

- [ ] Error handling:
  - Retry mechanism (3× retry)
  - Alert jika gagal
  - Manual sync option

- [ ] Log semua transaksi
- [ ] Audit trail

**Business Rules:**
- Real-time sync (immediately after event)
- Data validation sebelum kirim
- Rollback mechanism jika error

---

#### US-KA-027: Integrasi dengan Sistem Ekspedisi (3rd Party)
**Sebagai** System  
**Saya ingin** integrasi dengan sistem ekspedisi  
**Sehingga** tracking otomatis (jika pakai ekspedisi 3rd party)

**Acceptance Criteria:**
- [ ] **Create Shipment di Ekspedisi:**
  - Saat create DO, otomatis create shipment di sistem ekspedisi
  - Kirim data:
    - Shipment details
    - Pickup address (warehouse)
    - Delivery address (customer)
    - Package details (jumlah palet, berat, dimensi)
    - Required delivery date
  - Terima tracking number dari ekspedisi

- [ ] **Real-time Tracking:**
  - Pull tracking updates dari ekspedisi (API)
  - Update status di sistem internal
  - Display tracking info di dashboard & customer portal

- [ ] **POD dari Ekspedisi:**
  - Pull POD data dari ekspedisi
  - Sync dengan internal POD
  - Archive untuk audit

**Business Rules:**
- API integration dengan major ekspedisi (JNE, TIKI, dll)
- Real-time tracking update (setiap 30 menit)
- Fallback ke manual tracking jika API down

---

### Epic 7: Advanced Analytics

#### US-KA-028: Warehouse Optimization Analytics
**Sebagai** Logistics Manager  
**Saya ingin** warehouse optimization insights  
**Sehingga** bisa improve warehouse efficiency

**Acceptance Criteria:**
- [ ] Warehouse Analytics Dashboard:
  
  **Space Utilization:**
  - Current utilization: [X]%
  - Utilization trend (30 hari)
  - Optimal vs actual layout
  - Wasted space identification
  - Capacity forecast
  
  **Storage Efficiency:**
  - Average storage duration
  - FIFO compliance rate
  - Stock aging distribution
  - Slow-moving items
  - Fast-moving items
  
  **Picking Efficiency:**
  - Average picking time per palet
  - Picking route optimization
  - Walking distance per pick
  - Picking accuracy rate
  - Bottleneck identification
  
  **Layout Optimization:**
  - Heatmap: frequently picked items location
  - Recommendation: relocate fast-moving items to accessible locations
  - ABC analysis (A: fast-moving, B: medium, C: slow-moving)
  - Optimal slotting suggestion

- [ ] What-if scenario analysis
- [ ] ROI calculation untuk layout changes
- [ ] Export to PDF/Excel

**Business Rules:**
- Data-driven optimization
- Consider FIFO, accessibility, safety
- Continuous improvement

---

#### US-KA-029: Delivery Route Optimization
**Sebagai** Logistics Manager  
**Saya ingin** delivery route optimization  
**Sehingga** bisa reduce delivery time & cost

**Acceptance Criteria:**
- [ ] Route Optimization Tool:
  
  **Input:**
  - Multiple delivery addresses (untuk batch delivery)
  - Warehouse location (start point)
  - Delivery time windows
  - Vehicle capacity
  - Traffic data (real-time)
  
  **Optimization Algorithm:**
  - Calculate optimal route (shortest time/distance)
  - Consider:
    - Traffic conditions
    - Delivery time windows
    - Vehicle capacity
    - Driver working hours
  - Multiple route options (fastest, shortest, balanced)
  
  **Output:**
  - Optimized route (visual map)
  - Turn-by-turn directions
  - Estimated time per stop
  - Total distance & time
  - Fuel cost estimate
  - Comparison vs non-optimized route (savings)

- [ ] Save route untuk future reference
- [ ] Send route ke driver (mobile app)
- [ ] Track actual vs planned route
- [ ] Learn from historical data (ML)

**Business Rules:**
- Prioritize on-time delivery over shortest distance
- Consider traffic patterns (historical + real-time)
- Update route real-time jika ada perubahan (traffic, delay, dll)

---

#### US-KA-030: Cost Analysis & Profitability
**Sebagai** Finance Manager  
**Saya ingin** cost analysis untuk Khazanah Akhir operations  
**Sehingga** bisa optimize cost & improve profitability

**Acceptance Criteria:**
- [ ] Cost Analysis Dashboard:
  
  **Cost Breakdown:**
  - **Packaging Cost:**
    - Material cost (plastik wrapping, label, dll)
    - Labor cost (staff pengemasan)
    - Equipment cost (printer, dll)
    - Cost per rim
  - **Storage Cost:**
    - Warehouse rent (allocated)
    - Utilities (listrik, dll)
    - Equipment cost (rak, palet, dll)
    - Cost per rim per hari
  - **Shipping Cost:**
    - Transportation cost (fuel, toll, dll)
    - Labor cost (driver, staff pengiriman)
    - Ekspedisi cost (jika 3rd party)
    - Cost per delivery
  
  **Cost Trend:**
  - Cost per rim trend (30 hari)
  - Cost per delivery trend
  - Comparison: This month vs Last month
  - Budget vs Actual
  
  **Cost Optimization:**
  - High-cost areas identification
  - Cost reduction opportunities
  - ROI calculation untuk improvement initiatives
  - Benchmark vs industry standard

- [ ] Filter: periode, cost category
- [ ] Export to Excel
- [ ] Integration dengan finance system

**Business Rules:**
- Accurate cost allocation
- Regular cost review (monthly)
- Data-driven cost optimization

---

## 🎯 Key Performance Indicators (KPIs)

### KPI Pengemasan
- **Packaging Efficiency:** Target ≤ 15 menit per rim
- **Packaging Accuracy:** Target ≥ 99.5% (label, jumlah, dll)
- **Throughput:** Target [X] rim per hari
- **Cycle Time:** Target ≤ [Y] jam per PO (dari verifikasi selesai sampai siap simpan)

### KPI Inventory
- **Inventory Accuracy:** Target ≥ 99%
- **Stock Aging:** Target 0 rim > 30 hari
- **FIFO Compliance:** Target ≥ 98%
- **Warehouse Utilization:** Target 70-85% (optimal range)
- **Stock Opname Frequency:** Minimal 1× per bulan (full), 1× per minggu (cycle count)

### KPI Pengiriman
- **On-Time Delivery (OTD):** Target ≥ 95%
- **Perfect Delivery Rate:** Target ≥ 90% (on time + no damage + complete)
- **Average Delivery Time:** Target ≤ [X] jam (dari dispatch sampai delivery)
- **Picking Accuracy:** Target ≥ 99.5%
- **Document Accuracy:** Target 100%

### KPI Customer Satisfaction
- **Customer Satisfaction Score:** Target ≥ 4.5/5
- **Complaint Rate:** Target ≤ 1%
- **Response Time:** Target ≤ 2 jam untuk customer inquiry

---

## 📱 UI/UX Considerations

### Mobile-First Design
- **Prioritas:** Staff harus bisa input data di lapangan (gudang) dengan mudah
- **Large Touch Target:** Button minimal 44×44 px
- **Simple Interface:** Minimal distraction, focus pada task
- **Instant Feedback:** Visual & haptic feedback
- **Offline Capable:** Basic view & sync later
- **Fast:** Load time < 2 detik

### Warehouse Environment
- **High Contrast:** Mudah dibaca di lingkungan gudang (varying lighting)
- **Large Font:** Minimal 18px untuk body text
- **Glove-Friendly:** Touch targets besar (untuk staff pakai sarung tangan)
- **Dust/Water Resistant Device:** Recommend ruggedized tablet/phone

### Driver Mobile App
- **Simple Navigation:** Easy untuk driver yang tidak tech-savvy
- **Voice Guidance:** Turn-by-turn voice directions
- **Large Buttons:** Easy to tap while driving (saat parkir)
- **Offline Mode:** Work tanpa internet (sync later)
- **Battery Efficient:** Optimize untuk long trips

### Customer Portal
- **Professional Design:** Clean, modern, trustworthy
- **Easy Navigation:** Intuitive untuk non-technical users
- **Responsive:** Work di desktop, tablet, mobile
- **Fast Load:** < 3 detik
- **Secure:** HTTPS, strong authentication

---

## 🔐 Security & Access Control

### Role-Based Access

**Staff Khazanah Akhir - Pengemasan:**
- View: PO assigned to me, My performance
- Input: Penghitungan, Pengemasan data
- Cannot: Edit finalized data, View inventory detail, Delete data

**Staff Khazanah Akhir - Gudang:**
- View: Inventory, Lokasi, Picking list
- Input: Penyimpanan, Picking, Stock opname
- Cannot: Edit finalized data, View cost data, Delete data

**Staff Khazanah Akhir - Pengiriman:**
- View: DO, Delivery tracking
- Input: DO, Loading, Dispatch, POD
- Cannot: Edit inventory, View cost data, Delete data

**Driver:**
- View: DO assigned to me, Delivery instructions
- Input: POD, Delivery status
- Cannot: View other DO, Edit DO, View cost data

**Supervisor Khazanah Akhir:**
- View: All data (pengemasan, inventory, pengiriman)
- Input: Approve exceptions
- Edit: Data dengan audit trail
- Cannot: Delete finalized data

**Customer (Dirjen Bea Cukai):**
- View: Own orders only, Tracking, POD
- Cannot: View other customer data, Edit any data, View cost data

### Data Security
- **Encryption:** All data encrypted at rest & in transit (HTTPS)
- **Authentication:** Strong password policy, 2FA (optional)
- **Authorization:** Role-based access control (RBAC)
- **Audit Trail:** All actions logged (who, what, when, where)
- **Data Retention:** 5 tahun untuk compliance
- **Backup:** Daily backup, disaster recovery plan

---

## 🚀 Implementation Priority

### Phase 1 (MVP) - 2 Bulan
**Core Functionality:**
- [ ] US-KA-001 s/d US-KA-006: Penghitungan & Pengemasan
- [ ] US-KA-007 s/d US-KA-008: Penyimpanan & Inventory Dashboard
- [ ] US-KA-012 s/d US-KA-015: Delivery Order & Picking
- [ ] US-KA-016 s/d US-KA-017: Delivery Tracking & POD
- [ ] US-KA-021: Dashboard Overview

**Goal:** Paperless operations, real-time tracking, basic delivery management

---

### Phase 2 - 1.5 Bulan
**Advanced Features:**
- [ ] US-KA-009 s/d US-KA-011: Stock Opname, Location Finder, FIFO Monitoring
- [ ] US-KA-018: Delivery Performance Dashboard
- [ ] US-KA-022 s/d US-KA-024: Staff Performance, Alerts, Reporting
- [ ] US-KA-025 s/d US-KA-026: Integration (Verifikasi, SAP)

**Goal:** Advanced inventory management, performance analytics, integration

---

### Phase 3 - 1 Bulan
**Customer Experience:**
- [ ] US-KA-019 s/d US-KA-020: Customer Portal & Automated Communication
- [ ] US-KA-027: Integration dengan Ekspedisi 3rd Party

**Goal:** Excellent customer experience, transparency, proactive communication

---

### Phase 4 - 1 Bulan
**Optimization:**
- [ ] US-KA-028: Warehouse Optimization Analytics
- [ ] US-KA-029: Delivery Route Optimization
- [ ] US-KA-030: Cost Analysis & Profitability

**Goal:** Data-driven optimization, cost reduction, profitability improvement

---

## 📊 Success Metrics

### Adoption
- [ ] 100% staff Khazanah Akhir menggunakan sistem
- [ ] < 10 menit training time untuk basic operation
- [ ] User satisfaction score ≥ 4.5/5
- [ ] Customer portal adoption ≥ 80%

### Business Impact
- [ ] **Packaging Efficiency:** +30% (dengan digital support & auto-label)
- [ ] **Inventory Accuracy:** 95% → 99%+ (target +4%)
- [ ] **On-Time Delivery:** 85% → 95%+ (target +10%)
- [ ] **Stock Aging:** Reduce > 30 hari stock to 0
- [ ] **Picking Time:** -40% (dengan digital guidance)
- [ ] **Customer Satisfaction:** +20%
- [ ] **Cost per Delivery:** -15% (dengan route optimization)
- [ ] **Reporting Time:** 2 jam/hari → 0 jam (100% automated)

### Operational Excellence
- [ ] Paperless operation: 100%
- [ ] Real-time visibility: 100%
- [ ] FIFO compliance: ≥ 98%
- [ ] Perfect delivery rate: ≥ 90%
- [ ] Customer complaint rate: ≤ 1%

---

## 💡 Best Practices & Recommendations

### Untuk Staff Pengemasan
1. **Scan QR setiap rim** - untuk akurasi tracking
2. **Check label sebelum tempel** - untuk avoid error
3. **Foto palet sebelum kirim ke gudang** - untuk dokumentasi
4. **Update progress setiap 1 jam** - untuk real-time visibility
5. **Report issue immediately** - jangan tunggu sampai akhir shift

### Untuk Staff Gudang
1. **FIFO strictly** - yang masuk dulu harus keluar dulu
2. **Scan QR saat simpan & pick** - untuk inventory accuracy
3. **Organize gudang** - fast-moving items di lokasi accessible
4. **Regular cycle count** - untuk maintain inventory accuracy
5. **Report discrepancy immediately** - untuk quick resolution

### Untuk Staff Pengiriman
1. **Double-check dokumen** - sebelum dispatch
2. **Scan semua palet saat loading** - untuk akurasi
3. **Foto bukti loading** - untuk dokumentasi
4. **Communicate dengan driver** - clear instructions
5. **Follow up delivery** - sampai POD received

### Untuk Driver
1. **Check kendaraan sebelum berangkat** - untuk avoid breakdown
2. **Follow GPS route** - untuk optimal delivery time
3. **Communicate jika ada delay** - proactive communication
4. **Handle barang dengan hati-hati** - untuk avoid damage
5. **Complete POD immediately** - saat delivery selesai

### Untuk Supervisor
1. **Review dashboard setiap jam** - untuk proactive management
2. **Respond to alerts < 5 menit** - untuk minimize impact
3. **Daily coaching** - review performance dengan staff
4. **Weekly warehouse walk** - physical inspection
5. **Monthly performance review** - untuk continuous improvement

---

## 🎓 Training Requirements

### Staff Pengemasan Training (4 jam)
1. **Basic Operation (1 jam):**
   - Login & navigation
   - View PO & assignments
   - Input penghitungan
   - Pengemasan & labeling

2. **Auto-Label System (1 jam):**
   - Generate label
   - Print label
   - QR Code scanning
   - Troubleshooting

3. **Palet Management (1 jam):**
   - Penyusunan rim di palet
   - Palet labeling
   - Photo upload
   - Finalisasi

4. **Mobile App (1 jam):**
   - Mobile navigation
   - Quick actions
   - Offline mode

### Staff Gudang Training (6 jam)
1. **Warehouse Management (2 jam):**
   - Penyimpanan
   - Location system
   - FIFO principles
   - QR Code scanning

2. **Inventory Management (2 jam):**
   - Real-time inventory
   - Stock opname
   - Discrepancy handling
   - Cycle counting

3. **Picking Process (2 jam):**
   - Digital picking list
   - Route optimization
   - Picking accuracy
   - Mobile app

### Staff Pengiriman Training (4 jam)
1. **Delivery Order (1 jam):**
   - Create DO
   - Document generation
   - Editing & finalization

2. **Loading & Dispatch (1 jam):**
   - Loading checklist
   - QR scanning
   - Photo documentation
   - Dispatch process

3. **Tracking & POD (1 jam):**
   - Delivery tracking
   - Customer communication
   - POD process
   - Issue handling

4. **Customer Portal (1 jam):**
   - Portal overview
   - Customer support
   - Communication best practices

### Driver Training (2 jam)
1. **Mobile App Basics (1 jam):**
   - Login & navigation
   - View delivery instructions
   - GPS navigation
   - Contact support

2. **POD Process (1 jam):**
   - Arrival confirmation
   - BAST digital
   - Customer signature
   - Photo upload
   - Finalize delivery

---

## 📦 Label Specifications

### Rim Label (per 500 lembar)
```
┌─────────────────────────────────────────┐
│  PITA CUKAI - DIRJEN BEA CUKAI          │
├─────────────────────────────────────────┤
│  PO: PO-2025-001                        │
│  OBC: OBC-2024-001                      │
│  Rim: 1/59                              │
│  Qty: 500 lembar                        │
├─────────────────────────────────────────┤
│  Tgl Produksi: 27/12/2025               │
│  Tgl Verifikasi: 27/12/2025             │
│  Tgl Pengemasan: 27/12/2025             │
├─────────────────────────────────────────┤
│  Status: HCS (Hasil Cetak Sempurna)    │
├─────────────────────────────────────────┤
│  [QR CODE]     [BARCODE]                │
└─────────────────────────────────────────┘
```

### Palet Label
```
┌─────────────────────────────────────────┐
│  PALET - KHAZANAH AKHIR                 │
├─────────────────────────────────────────┤
│  Palet: PAL-KA-001                      │
│  PO: PO-2025-001                        │
│  OBC: OBC-2024-001                      │
├─────────────────────────────────────────┤
│  Rim: 1-20 dari 59 rim total            │
│  Total Qty: 10,000 lembar               │
├─────────────────────────────────────────┤
│  Tgl Pengemasan: 27/12/2025             │
│  Lokasi: A-10-3-C                       │
├─────────────────────────────────────────┤
│  [QR CODE PALET]                        │
└─────────────────────────────────────────┘
```

---

## 🚛 Delivery Documents

### Surat Jalan
- Nomor DO
- Tanggal pengiriman
- Pengirim (company details)
- Penerima (customer details)
- List items (PO, OBC, Qty)
- Kendaraan & driver details
- Signature (pengirim)

### Packing List
- Nomor DO
- Nomor PO & OBC
- Detail items:
  - Jumlah rim
  - Jumlah lembar
  - Jumlah palet
  - Palet numbers
- Total berat (estimasi)
- Total volume (estimasi)

### BAST (Berita Acara Serah Terima)
- Nomor DO
- Tanggal terima
- Pengirim details
- Penerima details
- List items dengan kondisi
- Signature penerima
- Signature pengirim (driver)
- Timestamp

---

**Catatan Akhir:**

User story ini comprehensive dan detail, mencakup semua aspek Khazanah Akhir dari penghitungan, pengemasan, warehouse management, sampai delivery ke customer. Fokus utama adalah:

1. **End-to-end visibility** - dari pengemasan sampai delivery
2. **Automation** - auto-label, auto-document, auto-notification
3. **Customer experience** - portal, tracking, proactive communication
4. **Inventory accuracy** - real-time tracking, FIFO, stock opname
5. **On-time delivery** - route optimization, real-time tracking, performance monitoring
6. **Integration** - seamless integration dengan Verifikasi, SAP, Ekspedisi 3rd party

Prioritas implementasi: **Phase 1 (MVP)** untuk establish core functionality, kemudian build advanced features di phase berikutnya.

**Prinsip desain:** Mobile-first (untuk staff di gudang), automation (untuk efficiency), customer-centric (untuk excellent customer experience), dan data-driven (untuk continuous improvement).
