# Epic 6: Inventory & Material Management

**Epic ID:** KHAZWAL-EPIC-06  
**Priority:** 🟡 Medium (Phase 2)  
**Estimated Duration:** 2 Minggu  

---

## 📋 Overview

Epic ini mencakup monitoring stok material (kertas blanko, tinta, plat cetak) untuk mendukung operasional Khazanah Awal.

**Note:** Stok material di-manage oleh SAP. Sistem ini hanya untuk monitoring dan alert, bukan sebagai source of truth untuk inventory.

---

## 🗄️ Database Reference

### SAP Integration (via API)
- Material stock (kertas, tinta) - GET from SAP
- Material consumption - POST to SAP
- Plat tracking - Local + SAP sync

### Local Tables
- `production_orders` - Referensi PO
- `khazwal_material_preparations` - Consumption record
- `alerts` - Low stock alerts

---

## 📝 Backlog Items

### US-KW-021: Monitoring Stok Kertas Blanko

| Field | Value |
|-------|-------|
| **ID** | US-KW-021 |
| **Story Points** | 5 |
| **Priority** | 🟡 Medium |
| **Dependencies** | SAP API Integration |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin melihat stok kertas blanko real-time, sehingga tahu apakah stok cukup untuk PO yang akan dikerjakan.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-021-BE-01 | Create SAP API client for material stock | 4h | Backend |
| KW-021-BE-02 | Create API `GET /api/khazwal/inventory/kertas` | 2h | Backend |
| KW-021-BE-03 | Calculate reserved stock (allocated to pending POs) | 2h | Backend |
| KW-021-BE-04 | Calculate available stock | 1h | Backend |
| KW-021-BE-05 | Estimate days until stock depleted | 2h | Backend |
| KW-021-BE-06 | Create low stock alert trigger | 2h | Backend |
| KW-021-BE-07 | Cache SAP response (Redis, 5 min TTL) | 2h | Backend |
| KW-021-FE-01 | Create page `KertasStockPage.vue` | 3h | Frontend |
| KW-021-FE-02 | Create stock summary card | 2h | Frontend |
| KW-021-FE-03 | Create usage history chart (7 days) | 2h | Frontend |
| KW-021-FE-04 | Create consumption trend chart | 2h | Frontend |
| KW-021-FE-05 | Show low stock warning | 1h | Frontend |

#### Acceptance Criteria
- [ ] Tampil stok kertas blanko:
  - Stok saat ini (lembar)
  - Minimum stok
  - Stok reserved (untuk PO yang sudah dialokasikan)
  - Stok available
  - Estimasi habis dalam [X] hari
- [ ] Alert jika stok < minimum
- [ ] History penggunaan (7 hari terakhir)
- [ ] Trend consumption

#### SAP API Reference
```
GET /sap/api/materials/{material_code}/stock
Response: {
    "material_code": "KERTAS-BLANKO-001",
    "material_name": "Kertas Blanko A",
    "unit": "LEMBAR",
    "stock_quantity": 50000,
    "minimum_stock": 10000,
    "location": "WH-01"
}
```

---

### US-KW-022: Monitoring Stok Tinta

| Field | Value |
|-------|-------|
| **ID** | US-KW-022 |
| **Story Points** | 5 |
| **Priority** | 🟡 Medium |
| **Dependencies** | SAP API Integration |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin melihat stok tinta per warna, sehingga tahu apakah tinta cukup.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-022-BE-01 | Create API `GET /api/khazwal/inventory/tinta` | 2h | Backend |
| KW-022-BE-02 | Query multiple tinta colors from SAP | 2h | Backend |
| KW-022-BE-03 | Calculate reserved & available per color | 2h | Backend |
| KW-022-BE-04 | Estimate days until depleted per color | 2h | Backend |
| KW-022-BE-05 | Create low stock alert per color | 2h | Backend |
| KW-022-FE-01 | Create page `TintaStockPage.vue` | 3h | Frontend |
| KW-022-FE-02 | Create color-coded stock cards | 2h | Frontend |
| KW-022-FE-03 | Create usage history per color | 2h | Frontend |
| KW-022-FE-04 | Create consumption trend per color | 2h | Frontend |
| KW-022-FE-05 | Show low stock warning per color | 1h | Frontend |

#### Acceptance Criteria
- [ ] Tampil stok tinta per warna:
  - Stok saat ini (kg)
  - Minimum stok
  - Stok reserved
  - Stok available
  - Estimasi habis dalam [X] hari
- [ ] Alert jika stok < minimum per warna
- [ ] History penggunaan per warna
- [ ] Trend consumption per warna

#### Business Rules
```
Minimum stok tinta per warna: 10 kg
Alert jika stok < minimum
```

---

### US-KW-023: Tracking Plat Cetak

| Field | Value |
|-------|-------|
| **ID** | US-KW-023 |
| **Story Points** | 8 |
| **Priority** | 🟡 Medium |
| **Dependencies** | US-KW-003 (Confirm Plat) |

**User Story:**
> Sebagai Staff Khazanah Awal, saya ingin tracking lokasi dan status plat cetak, sehingga mudah cari plat yang dibutuhkan.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-023-BE-01 | Create API `GET /api/khazwal/inventory/plat` | 2h | Backend |
| KW-023-BE-02 | Fetch plat list from SAP | 2h | Backend |
| KW-023-BE-03 | Track plat usage status locally | 2h | Backend |
| KW-023-BE-04 | Sync plat location with SAP | 2h | Backend |
| KW-023-BE-05 | Create plat usage history | 2h | Backend |
| KW-023-BE-06 | Create maintenance alert trigger | 2h | Backend |
| KW-023-BE-07 | Implement search by plat code | 1h | Backend |
| KW-023-FE-01 | Create page `PlatTrackingPage.vue` | 3h | Frontend |
| KW-023-FE-02 | Create plat list with status badge | 2h | Frontend |
| KW-023-FE-03 | Create plat detail modal | 2h | Frontend |
| KW-023-FE-04 | Create usage history view | 2h | Frontend |
| KW-023-FE-05 | Implement search bar | 1h | Frontend |
| KW-023-FE-06 | Implement filter (status, lokasi) | 2h | Frontend |

#### Acceptance Criteria
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

#### Plat Status Flow
```
TERSEDIA → SEDANG_DIGUNAKAN → TERSEDIA
    │              │
    │              └── MAINTENANCE (if needed)
    │                      │
    └──────────────────────┘
```

---

## 📊 Epic Summary

| User Story | Story Points | Priority | Phase |
|------------|--------------|----------|-------|
| US-KW-021 | 5 | Medium | Phase 2 |
| US-KW-022 | 5 | Medium | Phase 2 |
| US-KW-023 | 8 | Medium | Phase 2 |
| **Total** | **18** | - | - |

---

## 🔗 Dependencies Graph

```
SAP API Integration
    │
    ├── US-KW-021 (Kertas Stock)
    │
    ├── US-KW-022 (Tinta Stock)
    │
    └── US-KW-023 (Plat Tracking)
            │
            └── Requires US-KW-003 (Confirm Plat) for usage tracking
```

---

## 🧪 Testing Checklist

### Unit Tests
- [ ] SAP API client - fetch stock
- [ ] Available stock calculation
- [ ] Days until depleted calculation
- [ ] Low stock alert trigger

### Integration Tests
- [ ] SAP API integration (mock)
- [ ] API endpoint: GET kertas stock
- [ ] API endpoint: GET tinta stock
- [ ] API endpoint: GET plat list

---

## 📱 UI/UX Notes

- **Color-coded cards:** Green = OK, Yellow = Warning, Red = Critical
- **Progress bars:** Visual representation of stock level
- **Trend charts:** Help identify consumption patterns
- **Search:** Fast lookup for plat codes

---

**Last Updated:** 27 December 2025  
**Status:** Ready for Development
