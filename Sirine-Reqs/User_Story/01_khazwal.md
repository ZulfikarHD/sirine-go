# User Story: Khazanah Awal (Khazwal)

## Overview
Khazanah Awal bertanggung jawab untuk penyiapan material produksi (plat dan bahan baku) di awal proses, serta penghitungan dan pemotongan hasil cetak sebelum dikirim ke verifikasi.

---

## 🎭 Personas

### 1. Staff Khazanah Awal - Penyiapan Material
**Nama:** Dedi Kurniawan  
**Posisi:** Staff Khazanah Awal (Material Preparation)  
**Shift:** Pagi (07:00-15:00)  
**Pengalaman:** 3 tahun

**Tanggung Jawab:**
- Menyiapkan plat cetak sesuai kode plat di PO
- Menghitung dan menyiapkan kertas blanko
- Menyiapkan tinta sesuai spesifikasi warna
- Menyusun material di palet untuk dikirim ke Unit Cetak

**Pain Points:**
- Sering salah hitung jumlah kertas (lupa dibagi 2)
- Kesulitan tracking plat yang sedang dipakai
- Tidak tahu prioritas PO mana yang harus disiapkan dulu
- Manual input data ke form kertas

**Goals:**
- Tidak ada kesalahan penyiapan material
- Proses penyiapan lebih cepat dan efisien
- Tracking material lebih mudah
- Digital documentation

---

### 2. Staff Khazanah Awal - Penghitungan & Pemotongan
**Nama:** Siti Aminah  
**Posisi:** Staff Khazanah Awal (Counting & Cutting)  
**Shift:** Siang (15:00-23:00)  
**Pengalaman:** 5 tahun

**Tanggung Jawab:**
- Menghitung hasil cetak dari Unit Cetak
- Mencatat jumlah baik, rusak, kurang, lebih
- Koordinasi pemotongan lembar besar → lembar kirim
- Entry data hasil penghitungan dan pemotongan

**Pain Points:**
- Proses penghitungan manual memakan waktu lama
- Sering ada selisih antara target vs aktual
- Sulit tracking hasil cetak yang rusak per jenis kerusakan
- Laporan manual ke supervisor memakan waktu

**Goals:**
- Penghitungan lebih cepat dan akurat
- Real-time reporting ke supervisor
- Tracking kerusakan lebih detail
- Otomasi perhitungan konversi lembar besar → lembar kirim

---

### 3. Supervisor Khazanah Awal
**Nama:** Bambang Sutrisno  
**Posisi:** Supervisor Khazanah Awal  
**Shift:** Rotasi  
**Pengalaman:** 8 tahun

**Tanggung Jawab:**
- Monitoring proses penyiapan material
- Monitoring penghitungan dan pemotongan
- Koordinasi dengan Unit Cetak dan Verifikasi
- Reporting ke Production Manager
- Handling issue & escalation

**Pain Points:**
- Tidak ada visibility real-time progress pekerjaan
- Harus keliling cek satu-satu progress staff
- Laporan manual memakan waktu
- Sulit identifikasi bottleneck
- Tidak ada alert jika ada masalah

**Goals:**
- Dashboard real-time untuk monitoring semua aktivitas
- Alert otomatis jika ada issue
- Report otomatis
- Data-driven decision making

---

## 📱 User Stories

### Epic 1: Penyiapan Material (Material Preparation)

#### US-KW-001: Melihat Daftar PO yang Perlu Disiapkan
**Sebagai** Staff Khazanah Awal  
**Saya ingin** melihat daftar PO yang perlu disiapkan materialnya  
**Sehingga** saya tahu prioritas pekerjaan hari ini

**Acceptance Criteria:**
- [ ] Tampil daftar PO yang statusnya "Menunggu Penyiapan Material"
- [ ] Sorting berdasarkan prioritas (urgent, due date, sequence)
- [ ] Tampil info per PO:
  - Nomor PO & OBC
  - Jumlah cetak
  - Kode plat
  - Spesifikasi warna
  - Due date
  - Status prioritas (🔴 Urgent / 🟡 Normal / 🟢 Low)
- [ ] Filter berdasarkan: status, tanggal, prioritas
- [ ] Search berdasarkan nomor PO/OBC
- [ ] Responsive untuk tablet/mobile

**Business Rules:**
- PO dengan due date < 3 hari = Urgent (🔴)
- PO dengan due date 3-7 hari = Normal (🟡)
- PO dengan due date > 7 hari = Low (🟢)

---

#### US-KW-002: Memulai Proses Penyiapan Material
**Sebagai** Staff Khazanah Awal  
**Saya ingin** memulai proses penyiapan material untuk sebuah PO  
**Sehingga** sistem bisa tracking progress pekerjaan saya

**Acceptance Criteria:**
- [ ] Button "Mulai Persiapan" pada detail PO
- [ ] Sistem otomatis tampilkan:
  - Kode plat yang harus diambil
  - Jumlah kertas blanko yang harus disiapkan (otomatis dibagi 2)
  - Jenis dan estimasi jumlah tinta
  - Nomor palet yang akan digunakan
- [ ] Sistem ubah status PO menjadi "Sedang Disiapkan"
- [ ] Timestamp mulai penyiapan tercatat
- [ ] Nama staff yang handle tercatat

**Business Rules:**
- Jumlah Kertas Blanko = Jumlah Cetak PO ÷ 2
- 1 Palet maksimal untuk 2 PO (jika jumlah kecil)

---

#### US-KW-003: Konfirmasi Pengambilan Plat Cetak
**Sebagai** Staff Khazanah Awal  
**Saya ingin** konfirmasi bahwa plat cetak sudah diambil  
**Sehingga** sistem bisa tracking ketersediaan plat

**Acceptance Criteria:**
- [ ] Checkbox "Plat Sudah Diambil" dengan scan barcode/QR
- [ ] Sistem validasi kode plat sesuai dengan yang tertera di PO
- [ ] Jika salah plat, tampil warning dan tidak bisa lanjut
- [ ] Timestamp pengambilan plat tercatat
- [ ] Status plat berubah menjadi "Sedang Digunakan"
- [ ] Lokasi plat terupdate (di Unit Cetak)

**Business Rules:**
- Plat harus sesuai dengan kode plat di PO
- 1 Plat hanya bisa digunakan untuk 1 PO dalam waktu bersamaan

---

#### US-KW-004: Input Jumlah Kertas Blanko yang Disiapkan
**Sebagai** Staff Khazanah Awal  
**Saya ingin** input jumlah kertas blanko yang sudah disiapkan  
**Sehingga** ada record jika ada selisih dengan target

**Acceptance Criteria:**
- [ ] Input field "Jumlah Kertas Disiapkan" (dalam lembar besar)
- [ ] Tampil target: [Target] lembar besar
- [ ] Sistem otomatis hitung selisih (jika ada)
- [ ] Jika selisih > 5%, tampil warning dan wajib isi alasan
- [ ] Sistem kurangi stok kertas blanko di inventory
- [ ] Timestamp input tercatat

**Business Rules:**
- Target = Jumlah Cetak PO ÷ 2
- Toleransi selisih: ±5%
- Jika kurang > 5%, wajib isi alasan
- Jika lebih > 5%, wajib isi alasan

---

#### US-KW-005: Konfirmasi Penyiapan Tinta
**Sebagai** Staff Khazanah Awal  
**Saya ingin** konfirmasi jenis dan jumlah tinta yang sudah disiapkan  
**Sehingga** Unit Cetak tahu tinta apa saja yang tersedia

**Acceptance Criteria:**
- [ ] Checklist untuk setiap warna tinta sesuai spesifikasi PO
- [ ] Input jumlah tinta per warna (dalam kg)
- [ ] Sistem kurangi stok tinta di inventory
- [ ] Jika stok tinta < minimum, tampil alert
- [ ] Timestamp konfirmasi tercatat

**Business Rules:**
- Warna tinta sesuai spesifikasi di OBC/PO
- Minimum stok tinta per warna: 10 kg
- Alert jika stok < minimum

---

#### US-KW-006: Finalisasi Penyiapan Material
**Sebagai** Staff Khazanah Awal  
**Saya ingin** finalisasi bahwa semua material sudah siap di palet  
**Sehingga** Unit Cetak bisa mulai proses cetak

**Acceptance Criteria:**
- [ ] Button "Selesai - Kirim ke Unit Cetak"
- [ ] Sistem validasi semua checklist sudah complete:
  - ✅ Plat sudah diambil
  - ✅ Kertas sudah disiapkan
  - ✅ Tinta sudah disiapkan
- [ ] Input nomor palet
- [ ] Upload foto palet (opsional)
- [ ] Status PO berubah menjadi "Siap Cetak"
- [ ] Notifikasi ke Unit Cetak bahwa material sudah siap
- [ ] Timestamp selesai penyiapan tercatat
- [ ] Sistem hitung durasi penyiapan

**Business Rules:**
- Semua checklist harus complete sebelum bisa finalisasi
- Durasi penyiapan = Timestamp Selesai - Timestamp Mulai
---

### Epic 2: Penghitungan Hasil Cetak (Counting)

#### US-KW-007: Melihat Daftar PO yang Perlu Dihitung
**Sebagai** Staff Khazanah Awal  
**Saya ingin** melihat daftar PO yang hasil cetaknya sudah kembali dari Unit Cetak  
**Sehingga** saya tahu PO mana yang perlu dihitung

**Acceptance Criteria:**
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

**Business Rules:**
- PO yang selesai cetak lebih dulu harus dihitung lebih dulu (FIFO)
- Alert jika waiting time > 2 jam

---

#### US-KW-008: Memulai Proses Penghitungan
**Sebagai** Staff Khazanah Awal  
**Saya ingin** memulai proses penghitungan hasil cetak  
**Sehingga** sistem tracking progress penghitungan

**Acceptance Criteria:**
- [ ] Button "Mulai Penghitungan" pada detail PO
- [ ] Sistem tampilkan:
  - Target: [X] lembar besar
  - Mesin cetak: [Nama Mesin]
  - Operator cetak: [Nama Operator]
  - Waktu cetak: [Tanggal & Jam]
- [ ] Status PO berubah menjadi "Sedang Dihitung"
- [ ] Timestamp mulai penghitungan tercatat
- [ ] Nama staff yang handle tercatat

---

#### US-KW-009: Input Hasil Penghitungan
**Sebagai** Staff Khazanah Awal  
**Saya ingin** input hasil penghitungan lembar besar  
**Sehingga** ada record akurat berapa yang baik dan berapa yang rusak

**Acceptance Criteria:**
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

**Business Rules:**
- Total (Baik + Rusak) harus = Target ± toleransi
- Toleransi selisih: ±2%
- Jika rusak > 5%, wajib breakdown jenis kerusakan

---

#### US-KW-010: Finalisasi Penghitungan
**Sebagai** Staff Khazanah Awal  
**Saya ingin** finalisasi hasil penghitungan  
**Sehingga** proses bisa lanjut ke pemotongan

**Acceptance Criteria:**
- [ ] Button "Selesai Penghitungan - Lanjut Pemotongan"
- [ ] Sistem validasi:
  - ✅ Jumlah baik & rusak sudah diinput
  - ✅ Jika rusak > 5%, breakdown sudah diisi
  - ✅ Keputusan cetak ulang sudah dibuat (jika ada rusak)
- [ ] Status PO berubah menjadi "Siap Potong"
- [ ] Timestamp selesai penghitungan tercatat
- [ ] Sistem hitung durasi penghitungan
- [ ] Data tersimpan dan tidak bisa diubah (audit trail)

**Business Rules:**
- Semua validasi harus pass sebelum bisa finalisasi
- Data penghitungan tidak bisa diubah setelah finalisasi (untuk audit)

---

### Epic 3: Pemotongan (Cutting)

#### US-KW-011: Melihat Daftar PO yang Perlu Dipotong
**Sebagai** Staff Khazanah Awal  
**Saya ingin** melihat daftar PO yang siap untuk dipotong  
**Sehingga** saya tahu prioritas pemotongan

**Acceptance Criteria:**
- [ ] Tampil daftar PO dengan status "Siap Potong"
- [ ] Sorting berdasarkan prioritas dan FIFO
- [ ] Tampil info per PO:
  - Nomor PO & OBC
  - Jumlah lembar besar yang akan dipotong
  - Estimasi hasil: [X] lembar kirim (otomatis × 2)
  - Estimasi waktu pemotongan
  - Prioritas
- [ ] Filter berdasarkan tanggal, prioritas

**Business Rules:**
- Estimasi Hasil Pemotongan = Jumlah Baik (lembar besar) × 2

---

#### US-KW-012: Memulai Proses Pemotongan
**Sebagai** Staff Khazanah Awal  
**Saya ingin** memulai proses pemotongan  
**Sehingga** sistem tracking progress pemotongan

**Acceptance Criteria:**
- [ ] Button "Mulai Pemotongan" pada detail PO
- [ ] Sistem tampilkan:
  - Input: [X] lembar besar
  - Estimasi Output: [Y] lembar kirim (X × 2)
  - Mesin potong yang digunakan (dropdown)
  - Operator (auto-fill dari login user)
- [ ] Status PO berubah menjadi "Sedang Dipotong"
- [ ] Timestamp mulai pemotongan tercatat
- [ ] Nama staff & mesin tercatat

---

#### US-KW-013: Input Hasil Pemotongan
**Sebagai** Staff Khazanah Awal  
**Saya ingin** input hasil pemotongan  
**Sehingga** ada record akurat konversi lembar besar → lembar kirim

**Acceptance Criteria:**
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

**Business Rules:**
- Idealnya: Sisiran Kiri = Sisiran Kanan = Input (lembar besar)
- Toleransi waste: ≤ 2%
- Jika waste > 2%, wajib isi alasan dan foto bukti

---

#### US-KW-014: Finalisasi Pemotongan
**Sebagai** Staff Khazanah Awal  
**Saya ingin** finalisasi hasil pemotongan  
**Sehingga** hasil bisa dikirim ke Verifikasi

**Acceptance Criteria:**
- [ ] Button "Selesai - Kirim ke Verifikasi"
- [ ] Sistem validasi:
  - ✅ Hasil pemotongan sudah diinput
  - ✅ Jika waste > 2%, alasan sudah diisi
- [ ] Status PO berubah menjadi "Siap Verifikasi"
- [ ] Notifikasi ke Tim Verifikasi
- [ ] Timestamp selesai pemotongan tercatat
- [ ] Sistem hitung durasi pemotongan
- [ ] Data tersimpan dan tidak bisa diubah

**Business Rules:**
- Durasi Pemotongan = Timestamp Selesai - Timestamp Mulai
- Data tidak bisa diubah setelah finalisasi

---

### Epic 4: Dashboard & Monitoring (Untuk Supervisor)

#### US-KW-015: Dashboard Overview Khazanah Awal
**Sebagai** Supervisor Khazanah Awal  
**Saya ingin** melihat overview semua aktivitas Khazanah Awal  
**Sehingga** saya punya visibility real-time progress pekerjaan

**Acceptance Criteria:**
- [ ] Tampil di dashboard:
  - **Penyiapan Material:**
    - Menunggu: [X] PO
    - Sedang Dikerjakan: [Y] PO
    - Selesai Hari Ini: [Z] PO
    - Rata-rata Durasi: [N] menit
  - **Penghitungan:**
    - Menunggu: [X] PO
    - Sedang Dikerjakan: [Y] PO
    - Selesai Hari Ini: [Z] PO
    - Rata-rata Persentase Rusak: [N]%
  - **Pemotongan:**
    - Menunggu: [X] PO
    - Sedang Dikerjakan: [Y] PO
    - Selesai Hari Ini: [Z] PO
    - Rata-rata Waste: [N]%
- [ ] Chart: Trend persentase rusak (7 hari terakhir)
- [ ] Chart: Trend waste pemotongan (7 hari terakhir)
- [ ] Alert section (jika ada issue)
- [ ] Auto-refresh setiap 30 detik

---

#### US-KW-016: Monitoring Staff Performance
**Sebagai** Supervisor Khazanah Awal  
**Saya ingin** melihat performa individual staff  
**Sehingga** saya bisa evaluasi dan coaching

**Acceptance Criteria:**
- [ ] Tampil daftar staff Khazanah Awal
- [ ] Per staff tampil metrics:
  - **Penyiapan Material:**
    - Jumlah PO selesai hari ini
    - Rata-rata durasi penyiapan
    - Tingkat akurasi (% tanpa selisih)
  - **Penghitungan:**
    - Jumlah PO selesai hari ini
    - Rata-rata durasi penghitungan
    - Rata-rata persentase rusak yang ditemukan
  - **Pemotongan:**
    - Jumlah PO selesai hari ini
    - Rata-rata durasi pemotongan
    - Rata-rata waste
- [ ] Comparison dengan target dan rata-rata tim
- [ ] Filter berdasarkan: hari ini, minggu ini, bulan ini
- [ ] Export to Excel

---

#### US-KW-017: Alert & Notification untuk Supervisor
**Sebagai** Supervisor Khazanah Awal  
**Saya ingin** menerima alert jika ada issue  
**Sehingga** saya bisa cepat handling masalah

**Acceptance Criteria:**
- [ ] Alert otomatis untuk kondisi:
  - 🔴 PO menunggu penghitungan > 2 jam
  - 🔴 Persentase rusak > 10%
  - 🔴 Waste pemotongan > 5%
  - 🟡 Stok kertas blanko < minimum
  - 🟡 Stok tinta < minimum
  - 🟡 Plat tidak ditemukan
  - 🟡 Staff overtime > 2 jam
- [ ] Notification channel:
  - In-app notification (badge counter)
  - Push notification (mobile)
  - WhatsApp (untuk critical alert)
- [ ] Alert bisa di-acknowledge
- [ ] Alert history & resolution tracking

---

#### US-KW-018: Laporan Harian Khazanah Awal
**Sebagai** Supervisor Khazanah Awal  
**Saya ingin** generate laporan harian otomatis  
**Sehingga** saya tidak perlu buat laporan manual

**Acceptance Criteria:**
- [ ] Auto-generate laporan di akhir shift
- [ ] Isi laporan:
  - Summary aktivitas hari ini
  - Jumlah PO selesai per tahap
  - Rata-rata durasi per tahap
  - Total persentase rusak
  - Total waste pemotongan
  - Issue & resolution
  - Staff performance summary
  - Material usage
  - Recommendation/action items
- [ ] Format: PDF & Excel
- [ ] Auto-send via email ke Production Manager
- [ ] Bisa download manual kapan saja

---

### Epic 5: Analytics & Insights

#### US-KW-019: Analisa Efisiensi Pemotongan
**Sebagai** Supervisor Khazanah Awal  
**Saya ingin** melihat analisa efisiensi pemotongan  
**Sehingga** bisa minimize waste

**Acceptance Criteria:**
- [ ] Chart: Trend waste pemotongan
- [ ] Breakdown waste per:
  - Mesin potong
  - Operator
  - Shift
  - Jenis produk
- [ ] Comparison: Operator vs operator
- [ ] Best practice identification
- [ ] Waste cost calculation (Rp)
- [ ] Filter berdasarkan periode

---

#### US-KW-020: Analisa Durasi Proses
**Sebagai** Supervisor Khazanah Awal  
**Saya ingin** melihat analisa durasi setiap tahap proses  
**Sehingga** bisa identifikasi bottleneck

**Acceptance Criteria:**
- [ ] Chart: Rata-rata durasi per tahap
  - Penyiapan material
  - Penghitungan
  - Pemotongan
- [ ] Breakdown durasi per:
  - Staff
  - Shift
  - Jenis produk
- [ ] Identification bottleneck
- [ ] Trend durasi (time series)
- [ ] Target vs actual comparison
- [ ] Filter berdasarkan periode

---

### Epic 6: Inventory & Material Management

#### US-KW-021: Monitoring Stok Kertas Blanko
**Sebagai** Staff Khazanah Awal  
**Saya ingin** melihat stok kertas blanko real-time  
**Sehingga** tahu apakah stok cukup untuk PO yang akan dikerjakan

**Acceptance Criteria:**
- [ ] Tampil stok kertas blanko:
  - Stok saat ini (lembar)
  - Minimum stok
  - Stok reserved (untuk PO yang sudah dialokasikan)
  - Stok available
  - Estimasi habis dalam [X] hari
- [ ] Alert jika stok < minimum
- [ ] History penggunaan (7 hari terakhir)
- [ ] Trend consumption

---

#### US-KW-022: Monitoring Stok Tinta
**Sebagai** Staff Khazanah Awal  
**Saya ingin** melihat stok tinta per warna  
**Sehingga** tahu apakah tinta cukup

**Acceptance Criteria:**
- [ ] Tampil stok tinta per warna:
  - Stok saat ini (kg)
  - Minimum stok
  - Stok reserved
  - Stok available
  - Estimasi habis dalam [X] hari
- [ ] Alert jika stok < minimum per warna
- [ ] History penggunaan per warna
- [ ] Trend consumption per warna

---

#### US-KW-023: Tracking Plat Cetak
**Sebagai** Staff Khazanah Awal  
**Saya ingin** tracking lokasi dan status plat cetak  
**Sehingga** mudah cari plat yang dibutuhkan

**Acceptance Criteria:**
- [ ] List semua plat cetak dengan info:
  - Kode plat
  - Lokasi saat ini (Rak/Unit Cetak/Maintenance)
  - Status (Tersedia/Sedang Digunakan/Maintenance)
  - PO yang sedang menggunakan (jika ada)
  - History penggunaan
  - Kondisi plat (Baik/Perlu Maintenance)
- [ ] Search berdasarkan kode plat
- [ ] Filter berdasarkan status, lokasi
- [ ] Alert jika plat perlu maintenance

---

### Epic 7: Mobile App (Staff)

#### US-KW-024: Mobile App untuk Staff Khazanah Awal
**Sebagai** Staff Khazanah Awal  
**Saya ingin** akses aplikasi via mobile  
**Sehingga** bisa input data langsung di lapangan tanpa ke komputer

**Acceptance Criteria:**
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
- [ ] Fast & lightweight

---

### Epic 8: Integration & Automation

#### US-KW-025: Integrasi dengan SAP - Material Consumption
**Sebagai** System  
**Saya ingin** otomatis update consumption material ke SAP  
**Sehingga** inventory SAP selalu accurate

**Acceptance Criteria:**
- [ ] Setiap finalisasi penyiapan material:
  - Sistem otomatis kirim data consumption ke SAP:
    - Kertas blanko: [X] lembar
    - Tinta per warna: [Y] kg
    - Timestamp
    - PO number
- [ ] Setiap finalisasi pemotongan:
  - Sistem otomatis kirim data waste ke SAP
- [ ] Error handling jika SAP tidak respond
- [ ] Retry mechanism
- [ ] Log semua transaksi

---

#### US-KW-026: Integrasi dengan Unit Cetak
**Sebagai** System  
**Saya ingin** otomatis notifikasi ke Unit Cetak saat material siap  
**Sehingga** Unit Cetak langsung tahu dan bisa mulai cetak

**Acceptance Criteria:**
- [ ] Setiap finalisasi penyiapan material:
  - Notifikasi ke Unit Cetak (in-app)
  - Update status PO di sistem Unit Cetak
  - Kirim detail material yang sudah disiapkan
- [ ] Setiap ada work order cetak ulang:
  - Notifikasi urgent ke Unit Cetak
  - Prioritas di queue Unit Cetak

---

#### US-KW-027: Integrasi dengan Verifikasi
**Sebagai** System  
**Saya ingin** otomatis notifikasi ke Tim Verifikasi saat pemotongan selesai  
**Sehingga** Tim Verifikasi langsung tahu dan bisa mulai QC

**Acceptance Criteria:**
- [ ] Setiap finalisasi pemotongan:
  - Notifikasi ke Tim Verifikasi
  - Update status PO di sistem Verifikasi
  - Kirim detail hasil pemotongan:
    - Total lembar kirim
    - Sisiran kiri & kanan
    - Info kerusakan (jika ada dari penghitungan)

---

## 🎯 Key Performance Indicators (KPIs)

### KPI Penyiapan Material
- **Durasi Rata-rata Penyiapan:** Target ≤ 45 menit/PO
- **Akurasi Penyiapan:** Target ≥ 98% (tanpa selisih > 5%)
- **On-Time Delivery ke Unit Cetak:** Target ≥ 95%

### KPI Penghitungan
- **Durasi Rata-rata Penghitungan:** Target ≤ 30 menit/PO
- **Persentase Rusak Rata-rata:** Target ≤ 2%
- **Waiting Time:** Target ≤ 1 jam (dari selesai cetak sampai mulai hitung)

### KPI Pemotongan
- **Durasi Rata-rata Pemotongan:** Target ≤ 60 menit/PO
- **Waste Rate:** Target ≤ 1%
- **Akurasi Konversi:** Target ≥ 99% (sisiran kiri = sisiran kanan)

### KPI Overall
- **Throughput:** Jumlah PO selesai per hari
- **Cycle Time:** Durasi total (penyiapan + penghitungan + pemotongan)
- **First Time Right:** % PO tanpa issue

---

## 📱 UI/UX Considerations

### Mobile-First Design
- **Prioritas:** Input data di lapangan harus mudah
- **Large Touch Target:** Button minimal 44×44 px
- **Simple Form:** Max 3 field per screen
- **Instant Feedback:** Loading indicator, success/error message
- **Offline Capable:** Basic view & sync later

### Accessibility
- **Color Blind Friendly:** Jangan hanya pakai warna untuk status
- **Large Font:** Minimal 16px untuk body text
- **High Contrast:** Mudah dibaca di lingkungan pabrik
- **Icon + Text:** Jangan hanya icon

### Performance
- **Fast Load:** < 2 detik
- **Smooth Scroll:** 60 fps
- **Optimized Image:** Compress foto upload
- **Lazy Load:** Load data on demand

---

## 🔐 Security & Access Control

### Role-Based Access
- **Staff Khazanah Awal:**
  - View & input data sesuai tanggung jawabnya
  - Tidak bisa edit data yang sudah finalized
  - Tidak bisa delete data
- **Supervisor Khazanah Awal:**
  - Full access view semua data
  - Approval untuk keputusan tertentu
  - Edit data (dengan audit trail)
  - Generate report
- **Production Manager:**
  - View only (dashboard & report)
  - Export data

### Audit Trail
- Semua perubahan data tercatat:
  - Who (user)
  - What (action)
  - When (timestamp)
  - Before & after value

---

## 🚀 Implementation Priority

### Phase 1 (MVP) - 2 Bulan
- [ ] US-KW-001 s/d US-KW-006: Penyiapan Material
- [ ] US-KW-007 s/d US-KW-012: Penghitungan
- [ ] US-KW-013 s/d US-KW-016: Pemotongan
- [ ] US-KW-017: Dashboard Overview

### Phase 2 - 1 Bulan
- [ ] US-KW-018 s/d US-KW-020: Monitoring & Reporting
- [ ] US-KW-024 s/d US-KW-026: Inventory Management
- [ ] US-KW-027: Mobile App

### Phase 3 - 1 Bulan
- [ ] US-KW-021 s/d US-KW-023: Analytics
- [ ] US-KW-028 s/d US-KW-030: Integration

---

## 📊 Success Metrics

### Adoption
- [ ] 100% staff Khazanah Awal menggunakan sistem
- [ ] < 5 menit training time untuk basic operation
- [ ] User satisfaction score ≥ 4/5

### Business Impact
- [ ] Reduce durasi penyiapan material 30%
- [ ] Reduce waiting time 50%
- [ ] Reduce waste pemotongan 40%
- [ ] Increase data accuracy 95% → 99%
- [ ] Paperless operation 100%

---

**Catatan:**
User story ini comprehensive dan detail, siap untuk development. Prioritas implementasi bisa disesuaikan dengan kebutuhan bisnis dan resource yang tersedia.
