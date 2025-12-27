# Alur Proses Produksi

## Flow Diagram
```
PPIC → Khazwal (Penyiapan Plat) → Khazwal (Penyiapan Bahan Baku) → Cetak → 
Khazwal (Penghitungan) → Khazwal (Pemotongan) → Verifikasi → 
Khazkhir (Penghitungan) → Khazkhir (Pengemasan) → Khazkhir (Pengiriman)
```

## Catatan Penting
- **1 Rim** = 500 Lembar Kirim
- **1 Rim Lembar Besar** = 2 Rim Lembar Kirim (karena belum dibagi menjadi sisiran kiri dan kanan)

---

## 1. PPIC (Production Planning and Inventory Control)

PPIC menerima order dari Dirjen Bea Cukai, order tersebut kemudian diberikan identitas **Nomor Order Bea Cukai (OBC)**. 

Dalam OBC ini mengandung:
- Spesifikasi produk
- Total jumlah order dari Dirjen Bea Cukai

Tiap OBC kemudian dibagi menjadi beberapa **Nomor Production Order (PO)** untuk produksi. Di dalam PO ini juga mengandung data spesifikasi produk seperti pada OBC.

**Konversi OBC → PO:**
```
Total Produksi = Jumlah Order + (Jumlah Order × 6%)
Jumlah PO = MOD(Total Produksi; 40.000)
```

### Contoh:
```
OBC-2024-001
├─ Spesifikasi: Pita Cukai Rokok Golongan A
├─ Jumlah Order: 500.000 lembar
├─ Kode Plat: PLT-RC-A-001
└─ Warna: Merah, Putih, Hitam

Perhitungan PO:
Total Produksi = 500.000 + (500.000 × 6%)
               = 500.000 + 30.000
               = 530.000 lembar

Pembagian PO (maks 40.000 lembar/PO):
530.000 ÷ 40.000 = 13,25 → 14 PO

Hasil:
├─ PO-001 s/d PO-013: masing-masing 40.000 lembar (13 × 40.000 = 520.000)
└─ PO-014: 10.000 lembar

Total: 14 PO dengan 530.000 lembar

---

Contoh lain - Order 100.000 lembar:
Total Produksi = 100.000 + (100.000 × 6%)
               = 100.000 + 6.000
               = 106.000 lembar

Pembagian PO:
106.000 ÷ 40.000 = 2,65 → 3 PO

Hasil:
├─ PO-001: 40.000 lembar
├─ PO-002: 40.000 lembar
└─ PO-003: 26.000 lembar

Total: 3 PO dengan 106.000 lembar
```

---

## 2. Khazanah Awal - Penyiapan Bahan Baku

Saat PO sudah turun, Khazanah Awal akan:

1. **Mengambil plat cetak** sesuai dengan kode plat yang tercantum pada OBC
2. **Menyiapkan bahan baku** untuk dicetak:
   - Kertas kosong (blanko)
   - Tinta untuk cetak sesuai spesifikasi warna
   
**Catatan Jumlah Kertas:**
Karena kertas yang digunakan adalah lembar besar yang nantinya akan dipotong menjadi 2 bagian (sisiran kiri dan sisiran kanan), maka:
```
Jumlah Kertas yang Disiapkan = Jumlah Cetak yang Tertera ÷ 2
```

3. **Penempatan di palet** - Setelah semua bahan baku siap dan sesuai dengan PO, semua akan ditaruh dalam satu palet yang kemudian dikirimkan ke unit cetak untuk proses cetak.

### Contoh:
```
PO-001 (30.000 lembar kirim)

Persiapan:
├─ Plat Cetak: PLT-RC-A-001
├─ Kertas Blanko: 30.000 ÷ 2 = 15.000 lembar besar
├─ Tinta:
│  ├─ Merah: 5 kg
│  ├─ Putih: 3 kg
│  └─ Hitam: 2 kg
└─ Palet: PAL-001

Status: Siap dikirim ke Unit Cetak
```

---

## 3. Cetak

Unit Cetak melakukan proses cetak sesuai dengan spesifikasi yang tertera pada kartu PO. Setelah selesai cetak, hasil akan dikirim ke Khazanah Awal untuk penghitungan.

### Contoh:
```
PO-001 - Proses Cetak

Input:
├─ Kertas Blanko: 15.000 lembar besar
├─ Plat: PLT-RC-A-001
└─ Tinta: Merah, Putih, Hitam

Proses:
├─ Waktu Cetak: 8 jam
├─ Mesin: Mesin Cetak #3
└─ Operator: Tim A

Output:
├─ Hasil Cetak: 15.000 lembar besar (tercetak)
└─ Status: Dikirim ke Khazanah Awal untuk penghitungan
```

---

## 4. Khazanah Awal - Penghitungan dan Pemotongan

### Penghitungan
Khazanah Awal akan menghitung hasil cetak untuk menyesuaikan datanya dengan yang tertera pada kartu PO. Jika ada kurang atau rusak, maka akan dilakukan pendataan perihal:
- Kurang
- Lebih
- Penggantian kerusakan cetak

### Pemotongan
Setelah dilakukan penghitungan, lembar besar akan dikirimkan ke proses pemotongan, yang kemudian akan dipotong menjadi lembar kirim yang dibagi menjadi:
- **Sisiran Kiri**
- **Sisiran Kanan**

### Contoh:
```
PO-001 - Penghitungan

Target: 15.000 lembar besar
Hasil Cetak: 15.000 lembar besar
Rusak/Cacat: 50 lembar besar

Hasil Penghitungan:
├─ Baik: 14.950 lembar besar
├─ Rusak: 50 lembar besar
└─ Status: Perlu cetak ulang 50 lembar

---

PO-001 - Pemotongan

Input: 14.950 lembar besar

Proses Pemotongan:
14.950 lembar besar × 2 = 29.900 lembar kirim

Output:
├─ Sisiran Kiri: 14.950 lembar kirim
├─ Sisiran Kanan: 14.950 lembar kirim
└─ Total: 29.900 lembar kirim

Status: Dikirim ke Verifikasi
```

---

## 5. Verifikasi

Setelah dilakukan pemotongan, akan dilakukan proses verifikasi/QC. Produk akan dibagi menjadi 2 kategori:

1. **HCS (Hasil Cetak Sempurna)**
2. **HCTS (Hasil Cetak Tidak Sempurna)**

Di akhir proses verifikasi akan dilakukan pendataan atau entry kedua kategori ini terkait hasil cetakan pada PO yang dikerjakan.

### Contoh:
```
PO-001 - Verifikasi/QC

Input: 29.900 lembar kirim

Proses Verifikasi:
├─ Tim QC: Tim Verifikasi A
├─ Kriteria: Warna, Ketajaman, Posisi Cetak, Kerusakan Fisik
└─ Waktu: 4 jam

Hasil Verifikasi:
├─ HCS (Hasil Cetak Sempurna): 29.500 lembar kirim
│  ├─ Sisiran Kiri: 14.750 lembar
│  └─ Sisiran Kanan: 14.750 lembar
│
└─ HCTS (Hasil Cetak Tidak Sempurna): 400 lembar kirim
   ├─ Warna Pudar: 150 lembar
   ├─ Posisi Cetak Miring: 100 lembar
   ├─ Sobek/Rusak: 100 lembar
   └─ Noda Tinta: 50 lembar

Persentase HCS: 98,66%
Status: HCS → Khazanah Akhir | HCTS → Unit Pengelolaan HCTS
```

---

## 6. Khazanah Akhir

Setelah selesai verifikasi:

### Pengemasan
- **HCS** akan dilakukan pengemasan per 1 Rim yang dikelompokkan per PO
- **HCTS** akan dikirimkan ke unit Pengelolaan HCTS

### Pengiriman
Setelah dilakukan pengemasan per 1 Rim yang dikelompokkan per PO, maka akan dilakukan pengiriman oleh unit pengiriman.

### Contoh:
```
PO-001 - Khazanah Akhir

Input HCS: 29.500 lembar kirim

Penghitungan Rim:
29.500 lembar ÷ 500 lembar/rim = 59 rim

Pengemasan:
├─ Rim 1-59: Masing-masing 500 lembar
├─ Kemasan: Plastik wrapping + Label PO
├─ Label berisi:
│  ├─ Nomor PO: PO-001
│  ├─ Nomor OBC: OBC-2024-001
│  ├─ Nomor Rim: 1/59, 2/59, ..., 59/59
│  ├─ Tanggal Produksi: 27 Desember 2025
│  └─ Status: HCS
└─ Disusun di Palet: PAL-KA-001

Pengiriman:
├─ Tujuan: Gudang Dirjen Bea Cukai
├─ Jumlah: 59 rim (29.500 lembar)
├─ Kendaraan: Truk T-001
├─ Tanggal Kirim: 28 Desember 2025
├─ Surat Jalan: SJ-2024-001
└─ Status: Siap Kirim

---

HCTS: 400 lembar → Dikirim ke Unit Pengelolaan HCTS
```
