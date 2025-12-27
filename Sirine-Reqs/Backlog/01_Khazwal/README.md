# 📋 Backlog Khazanah Awal (Khazwal)

**Department:** Khazanah Awal  
**Total Epics:** 8  
**Total User Stories:** 27  
**Total Story Points:** 205  

---

## 🎯 Overview

Khazanah Awal bertanggung jawab untuk:
1. Penyiapan material produksi (plat, kertas blanko, tinta)
2. Penghitungan hasil cetak dari Unit Cetak
3. Pemotongan lembar besar → lembar kirim
4. Monitoring dan reporting

---

## 📊 Epic Summary

| Epic | Nama | Story Points | Priority | Phase | Status |
|------|------|--------------|----------|-------|--------|
| [Epic 1](Epic_01_Material_Preparation.md) | Penyiapan Material | 36 | 🔴 High | MVP | ⬜ Backlog |
| [Epic 2](Epic_02_Counting.md) | Penghitungan Hasil Cetak | 23 | 🔴 High | MVP | ⬜ Backlog |
| [Epic 3](Epic_03_Cutting.md) | Pemotongan | 26 | 🔴 High | MVP | ⬜ Backlog |
| [Epic 4](Epic_04_Dashboard_Monitoring.md) | Dashboard & Monitoring | 42 | 🟡 Medium | MVP/P2 | ⬜ Backlog |
| [Epic 5](Epic_05_Analytics.md) | Analytics & Insights | 16 | 🟢 Low | Phase 3 | ⬜ Backlog |
| [Epic 6](Epic_06_Inventory_Management.md) | Inventory Management | 18 | 🟡 Medium | Phase 2 | ⬜ Backlog |
| [Epic 7](Epic_07_Mobile_App.md) | Mobile App (PWA) | 21 | 🟡 Medium | Phase 2 | ⬜ Backlog |
| [Epic 8](Epic_08_Integration.md) | Integration & Automation | 23 | 🟢 Low | Phase 3 | ⬜ Backlog |
| **Total** | | **205** | | | |

---

## 📅 Implementation Phases

### Phase 1: MVP (Target: 2 Bulan)
**Story Points:** ~127

| User Story | Epic | Description |
|------------|------|-------------|
| US-KW-001 | Epic 1 | Melihat Daftar PO yang Perlu Disiapkan |
| US-KW-002 | Epic 1 | Memulai Proses Penyiapan Material |
| US-KW-003 | Epic 1 | Konfirmasi Pengambilan Plat Cetak |
| US-KW-004 | Epic 1 | Input Jumlah Kertas Blanko |
| US-KW-005 | Epic 1 | Konfirmasi Penyiapan Tinta |
| US-KW-006 | Epic 1 | Finalisasi Penyiapan Material |
| US-KW-007 | Epic 2 | Melihat Daftar PO yang Perlu Dihitung |
| US-KW-008 | Epic 2 | Memulai Proses Penghitungan |
| US-KW-009 | Epic 2 | Input Hasil Penghitungan |
| US-KW-010 | Epic 2 | Finalisasi Penghitungan |
| US-KW-011 | Epic 3 | Melihat Daftar PO yang Perlu Dipotong |
| US-KW-012 | Epic 3 | Memulai Proses Pemotongan |
| US-KW-013 | Epic 3 | Input Hasil Pemotongan |
| US-KW-014 | Epic 3 | Finalisasi Pemotongan |
| US-KW-015 | Epic 4 | Dashboard Overview |
| US-KW-017 | Epic 4 | Alert & Notification |

### Phase 2 (Target: 1 Bulan)
**Story Points:** ~57

| User Story | Epic | Description |
|------------|------|-------------|
| US-KW-016 | Epic 4 | Monitoring Staff Performance |
| US-KW-018 | Epic 4 | Laporan Harian Otomatis |
| US-KW-021 | Epic 6 | Monitoring Stok Kertas Blanko |
| US-KW-022 | Epic 6 | Monitoring Stok Tinta |
| US-KW-023 | Epic 6 | Tracking Plat Cetak |
| US-KW-024 | Epic 7 | Mobile App (PWA) |

### Phase 3 (Target: 1 Bulan)
**Story Points:** ~39

| User Story | Epic | Description |
|------------|------|-------------|
| US-KW-019 | Epic 5 | Analisa Efisiensi Pemotongan |
| US-KW-020 | Epic 5 | Analisa Durasi Proses |
| US-KW-025 | Epic 8 | Integrasi SAP - Material Consumption |
| US-KW-026 | Epic 8 | Integrasi dengan Unit Cetak |
| US-KW-027 | Epic 8 | Integrasi dengan Verifikasi |

---

## 🗄️ Database Tables Reference

### Primary Tables
| Table | Description |
|-------|-------------|
| `khazwal_material_preparations` | Data penyiapan material |
| `khazwal_counting_results` | Data hasil penghitungan |
| `khazwal_cutting_results` | Data hasil pemotongan |

### Related Tables
| Table | Description |
|-------|-------------|
| `production_orders` | PO utama |
| `po_stage_tracking` | Tracking stage PO |
| `users` | Staff data |
| `activity_logs` | Audit trail |
| `notifications` | Notifikasi |
| `alerts` | Alert system |
| `verification_labels` | Output ke Verifikasi |

### SAP Integration (via API)
- Material stock (kertas, tinta)
- Material consumption
- Plat data

---

## 🔑 Key Performance Indicators (KPIs)

### Penyiapan Material
- **Durasi Rata-rata:** Target ≤ 45 menit/PO
- **Akurasi:** Target ≥ 98%
- **On-Time Delivery:** Target ≥ 95%

### Penghitungan
- **Durasi Rata-rata:** Target ≤ 30 menit/PO
- **Persentase Rusak:** Target ≤ 2%
- **Waiting Time:** Target ≤ 1 jam

### Pemotongan
- **Durasi Rata-rata:** Target ≤ 60 menit/PO
- **Waste Rate:** Target ≤ 1%
- **Akurasi Konversi:** Target ≥ 99%

---

## 🔄 Process Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                         KHAZANAH AWAL                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  PO Masuk     Material Prep      Cetak         Counting         │
│     │              │               │              │              │
│     ▼              ▼               ▼              ▼              │
│  ┌─────┐      ┌─────────┐     ┌────────┐    ┌─────────┐         │
│  │Queue│─────►│Siapkan  │────►│ UNIT   │───►│Hitung   │         │
│  │ PO  │      │Material │     │ CETAK  │    │Hasil    │         │
│  └─────┘      └─────────┘     └────────┘    └─────────┘         │
│                    │                              │              │
│                    ▼                              ▼              │
│              ┌──────────┐                  ┌──────────┐         │
│              │• Plat    │                  │• Baik    │         │
│              │• Kertas  │                  │• Rusak   │         │
│              │• Tinta   │                  │• Breakdown│        │
│              └──────────┘                  └──────────┘         │
│                                                  │              │
│                                                  ▼              │
│                                            ┌─────────┐         │
│                                            │Potong   │         │
│                                            │Lembar   │         │
│                                            └─────────┘         │
│                                                  │              │
│                                                  ▼              │
│                                            ┌──────────┐        │
│                                            │• Sisir L │        │
│                                            │• Sisir R │        │
│                                            │• Waste   │        │
│                                            └──────────┘        │
│                                                  │              │
│                                                  ▼              │
│                                            ┌──────────┐        │
│                                            │ VERIFI-  │        │
│                                            │ KASI     │        │
│                                            └──────────┘        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 👥 Personas

| Persona | Role | Focus |
|---------|------|-------|
| Dedi Kurniawan | Staff - Material Prep | Epic 1 |
| Siti Aminah | Staff - Counting & Cutting | Epic 2, 3 |
| Bambang Sutrisno | Supervisor | Epic 4, 5 |

---

## 📱 Device Support

| Device | Support Level |
|--------|--------------|
| Desktop (Chrome, Firefox, Edge) | Full |
| Tablet (iPad, Android) | Full |
| Mobile (iOS, Android) | PWA (Phase 2) |

---

## 🔗 Cross-Department Integration

### Input From
- **PO Management:** Production Orders

### Output To
- **Unit Cetak:** Material ready notification
- **Verifikasi:** Cutting results + verification labels
- **SAP:** Material consumption

---

## 📝 Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 27 Dec 2025 | System | Initial backlog creation |

---

**Status:** ✅ Ready for Sprint Planning  
**Last Updated:** 27 December 2025
