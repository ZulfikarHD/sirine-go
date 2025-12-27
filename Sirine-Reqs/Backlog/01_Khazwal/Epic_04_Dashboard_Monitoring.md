# Epic 4: Dashboard & Monitoring (Supervisor)

**Epic ID:** KHAZWAL-EPIC-04  
**Priority:** 🟡 Medium (Phase 1 - MVP)  
**Estimated Duration:** 2-3 Minggu  

---

## 📋 Overview

Epic ini mencakup dashboard dan monitoring untuk Supervisor Khazanah Awal, termasuk overview aktivitas, monitoring performa staff, alert system, dan laporan harian otomatis.

---

## 🗄️ Database Reference

### Primary Tables (Aggregation Source)
- `khazwal_material_preparations`
- `khazwal_counting_results`
- `khazwal_cutting_results`

### Related Tables
- `production_orders` - Referensi PO
- `users` - Staff data & performance
- `alerts` - Alert management
- `notifications` - Notification delivery
- `activity_logs` - Audit trail

### Views (Recommended)
- `v_khazwal_daily_summary` - Daily aggregation view
- `v_khazwal_staff_performance` - Staff metrics view

---

## 📝 Backlog Items

### US-KW-015: Dashboard Overview Khazanah Awal

| Field | Value |
|-------|-------|
| **ID** | US-KW-015 |
| **Story Points** | 13 |
| **Priority** | 🟡 Medium |
| **Dependencies** | Epic 1, 2, 3 completed |

**User Story:**
> Sebagai Supervisor Khazanah Awal, saya ingin melihat overview semua aktivitas Khazanah Awal, sehingga saya punya visibility real-time progress pekerjaan.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-015-BE-01 | Create API `GET /api/khazwal/dashboard/overview` | 4h | Backend |
| KW-015-BE-02 | Create `KhazwalDashboardService.php` | 3h | Backend |
| KW-015-BE-03 | Query material prep stats (waiting, in_progress, completed) | 2h | Backend |
| KW-015-BE-04 | Query counting stats | 2h | Backend |
| KW-015-BE-05 | Query cutting stats | 2h | Backend |
| KW-015-BE-06 | Calculate average duration per stage | 2h | Backend |
| KW-015-BE-07 | Calculate defect & waste percentages | 2h | Backend |
| KW-015-BE-08 | Create defect trend (7 days) query | 2h | Backend |
| KW-015-BE-09 | Create waste trend (7 days) query | 2h | Backend |
| KW-015-BE-10 | Implement Redis caching (30 sec TTL) | 2h | Backend |
| KW-015-FE-01 | Create page `KhazwalDashboardPage.vue` | 4h | Frontend |
| KW-015-FE-02 | Create stats card components (3 sections) | 3h | Frontend |
| KW-015-FE-03 | Create defect trend chart (line chart) | 3h | Frontend |
| KW-015-FE-04 | Create waste trend chart (line chart) | 2h | Frontend |
| KW-015-FE-05 | Create alert section component | 2h | Frontend |
| KW-015-FE-06 | Implement auto-refresh (30 seconds) | 2h | Frontend |
| KW-015-FE-07 | Make responsive for tablet | 2h | Frontend |

#### Acceptance Criteria
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

#### Dashboard Layout
```
┌─────────────────────────────────────────────────────────────┐
│                    KHAZANAH AWAL DASHBOARD                  │
├─────────────────┬─────────────────┬─────────────────────────┤
│  📦 PENYIAPAN   │  📊 PENGHITUNG  │      ✂️ PEMOTONGAN      │
│  Menunggu: 5    │  Menunggu: 3    │      Menunggu: 2        │
│  Proses: 2      │  Proses: 1      │      Proses: 1          │
│  Selesai: 12    │  Selesai: 10    │      Selesai: 8         │
│  Avg: 42 min    │  Avg Rusak: 2%  │      Avg Waste: 1.2%    │
├─────────────────┴─────────────────┴─────────────────────────┤
│  📈 TREND RUSAK (7 Hari)        │  📉 TREND WASTE (7 Hari) │
│  [LINE CHART]                   │  [LINE CHART]            │
├─────────────────────────────────┴────────────────────────────┤
│  🚨 ALERTS                                                   │
│  • PO-12345 menunggu penghitungan > 2 jam (URGENT)          │
│  • Stok kertas blanko < minimum (WARNING)                    │
└──────────────────────────────────────────────────────────────┘
```

#### Database Query Reference
```sql
-- Daily Summary View
CREATE OR REPLACE VIEW v_khazwal_daily_summary AS
SELECT 
    -- Material Prep Stats
    COUNT(DISTINCT CASE WHEN mp.status = 'PENDING' THEN mp.id END) as prep_waiting,
    COUNT(DISTINCT CASE WHEN mp.status = 'IN_PROGRESS' THEN mp.id END) as prep_in_progress,
    COUNT(DISTINCT CASE WHEN mp.status = 'COMPLETED' 
        AND DATE(mp.completed_at) = CURRENT_DATE THEN mp.id END) as prep_completed_today,
    AVG(CASE WHEN mp.status = 'COMPLETED' 
        AND DATE(mp.completed_at) = CURRENT_DATE THEN mp.duration_minutes END) as prep_avg_duration,
    
    -- Counting Stats
    COUNT(DISTINCT CASE WHEN cr.status IS NULL 
        AND po.current_stage = 'KHAZWAL_COUNTING' THEN po.id END) as count_waiting,
    COUNT(DISTINCT CASE WHEN cr.status = 'IN_PROGRESS' THEN cr.id END) as count_in_progress,
    COUNT(DISTINCT CASE WHEN cr.status = 'COMPLETED' 
        AND DATE(cr.completed_at) = CURRENT_DATE THEN cr.id END) as count_completed_today,
    AVG(CASE WHEN cr.status = 'COMPLETED' 
        AND DATE(cr.completed_at) = CURRENT_DATE THEN cr.percentage_defect END) as count_avg_defect,
    
    -- Cutting Stats
    COUNT(DISTINCT CASE WHEN cut.status IS NULL 
        AND po.current_stage = 'KHAZWAL_CUTTING' THEN po.id END) as cut_waiting,
    COUNT(DISTINCT CASE WHEN cut.status = 'IN_PROGRESS' THEN cut.id END) as cut_in_progress,
    COUNT(DISTINCT CASE WHEN cut.status = 'COMPLETED' 
        AND DATE(cut.completed_at) = CURRENT_DATE THEN cut.id END) as cut_completed_today,
    AVG(CASE WHEN cut.status = 'COMPLETED' 
        AND DATE(cut.completed_at) = CURRENT_DATE THEN cut.waste_percentage END) as cut_avg_waste

FROM production_orders po
LEFT JOIN khazwal_material_preparations mp ON mp.production_order_id = po.id
LEFT JOIN khazwal_counting_results cr ON cr.production_order_id = po.id
LEFT JOIN khazwal_cutting_results cut ON cut.production_order_id = po.id
WHERE po.deleted_at IS NULL;

-- Trend Query (7 days)
SELECT 
    DATE(completed_at) as date,
    AVG(percentage_defect) as avg_defect,
    AVG(waste_percentage) as avg_waste
FROM khazwal_counting_results cr
JOIN khazwal_cutting_results cut ON cut.production_order_id = cr.production_order_id
WHERE cr.completed_at >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY DATE(completed_at)
ORDER BY date;
```

---

### US-KW-016: Monitoring Staff Performance

| Field | Value |
|-------|-------|
| **ID** | US-KW-016 |
| **Story Points** | 13 |
| **Priority** | 🟡 Medium |
| **Dependencies** | US-KW-015 |

**User Story:**
> Sebagai Supervisor Khazanah Awal, saya ingin melihat performa individual staff, sehingga saya bisa evaluasi dan coaching.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-016-BE-01 | Create API `GET /api/khazwal/dashboard/staff-performance` | 3h | Backend |
| KW-016-BE-02 | Query material prep metrics per staff | 2h | Backend |
| KW-016-BE-03 | Query counting metrics per staff | 2h | Backend |
| KW-016-BE-04 | Query cutting metrics per staff | 2h | Backend |
| KW-016-BE-05 | Calculate comparison vs target & team average | 2h | Backend |
| KW-016-BE-06 | Implement date filter (today, week, month) | 2h | Backend |
| KW-016-BE-07 | Create Excel export endpoint | 3h | Backend |
| KW-016-FE-01 | Create page `StaffPerformancePage.vue` | 4h | Frontend |
| KW-016-FE-02 | Create staff list with metrics table | 3h | Frontend |
| KW-016-FE-03 | Create comparison visualization (vs target) | 2h | Frontend |
| KW-016-FE-04 | Implement date filter UI | 2h | Frontend |
| KW-016-FE-05 | Add Excel export button | 1h | Frontend |
| KW-016-FE-06 | Make responsive | 2h | Frontend |

#### Acceptance Criteria
- [ ] Tampil daftar staff Khazanah Awal
- [ ] Per staff tampil metrics:
  - **Penyiapan Material:**
    - Jumlah PO selesai hari ini
    - Rata-rata durasi penyiapan
    - Tingkat akurasi (% tanpa selisih > 5%)
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

#### KPI Targets Reference
```
Penyiapan Material:
- Target Durasi: ≤ 45 menit/PO
- Target Akurasi: ≥ 98%

Penghitungan:
- Target Durasi: ≤ 30 menit/PO
- Target Rusak: ≤ 2%

Pemotongan:
- Target Durasi: ≤ 60 menit/PO
- Target Waste: ≤ 1%
```

#### Database Query Reference
```sql
-- Staff Performance View
SELECT 
    u.id as staff_id,
    u.nip,
    u.full_name,
    u.default_shift,
    
    -- Material Prep Metrics
    COUNT(DISTINCT CASE WHEN mp.status = 'COMPLETED' THEN mp.id END) as prep_completed,
    AVG(mp.duration_minutes) as prep_avg_duration,
    ROUND(
        COUNT(DISTINCT CASE WHEN mp.status = 'COMPLETED' 
            AND ABS(mp.kertas_blanko_variance) <= mp.kertas_blanko_quantity * 0.05 THEN mp.id END)::numeric
        / NULLIF(COUNT(DISTINCT CASE WHEN mp.status = 'COMPLETED' THEN mp.id END), 0) * 100,
        2
    ) as prep_accuracy,
    
    -- Counting Metrics
    COUNT(DISTINCT CASE WHEN cr.status = 'COMPLETED' THEN cr.id END) as count_completed,
    AVG(cr.duration_minutes) as count_avg_duration,
    AVG(cr.percentage_defect) as count_avg_defect,
    
    -- Cutting Metrics
    COUNT(DISTINCT CASE WHEN cut.status = 'COMPLETED' THEN cut.id END) as cut_completed,
    AVG(cut.duration_minutes) as cut_avg_duration,
    AVG(cut.waste_percentage) as cut_avg_waste
    
FROM users u
LEFT JOIN khazwal_material_preparations mp ON mp.prepared_by = u.id
    AND DATE(mp.completed_at) BETWEEN :start_date AND :end_date
LEFT JOIN khazwal_counting_results cr ON cr.counted_by = u.id
    AND DATE(cr.completed_at) BETWEEN :start_date AND :end_date
LEFT JOIN khazwal_cutting_results cut ON cut.cut_by = u.id
    AND DATE(cut.completed_at) BETWEEN :start_date AND :end_date
WHERE u.department = 'KHAZWAL'
    AND u.is_active = TRUE
GROUP BY u.id, u.nip, u.full_name, u.default_shift;
```

---

### US-KW-017: Alert & Notification untuk Supervisor

| Field | Value |
|-------|-------|
| **ID** | US-KW-017 |
| **Story Points** | 8 |
| **Priority** | 🟡 Medium |
| **Dependencies** | US-KW-015 |

**User Story:**
> Sebagai Supervisor Khazanah Awal, saya ingin menerima alert jika ada issue, sehingga saya bisa cepat handling masalah.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-017-BE-01 | Create `AlertService.php` for Khazwal | 3h | Backend |
| KW-017-BE-02 | Implement alert: PO waiting counting > 2 jam | 2h | Backend |
| KW-017-BE-03 | Implement alert: Persentase rusak > 10% | 2h | Backend |
| KW-017-BE-04 | Implement alert: Waste > 5% | 2h | Backend |
| KW-017-BE-05 | Implement alert: Stok material < minimum (SAP) | 2h | Backend |
| KW-017-BE-06 | Implement alert: Staff overtime > 2 jam | 2h | Backend |
| KW-017-BE-07 | Create scheduled job for alert checking | 2h | Backend |
| KW-017-BE-08 | Create API `GET /api/khazwal/alerts` | 2h | Backend |
| KW-017-BE-09 | Create API `POST /api/alerts/{id}/acknowledge` | 1h | Backend |
| KW-017-BE-10 | Implement in-app notification | 2h | Backend |
| KW-017-BE-11 | Implement WhatsApp notification (critical) | 3h | Backend |
| KW-017-FE-01 | Create alert list component | 3h | Frontend |
| KW-017-FE-02 | Create notification badge counter | 2h | Frontend |
| KW-017-FE-03 | Create acknowledge UI | 1h | Frontend |
| KW-017-FE-04 | Create alert history page | 2h | Frontend |

#### Acceptance Criteria
- [ ] Alert otomatis untuk kondisi:
  - 🔴 PO menunggu penghitungan > 2 jam
  - 🔴 Persentase rusak > 10%
  - 🔴 Waste pemotongan > 5%
  - 🟡 Stok kertas blanko < minimum
  - 🟡 Stok tinta < minimum
  - 🟡 Staff overtime > 2 jam
- [ ] Notification channel:
  - In-app notification (badge counter)
  - WhatsApp (untuk critical alert)
- [ ] Alert bisa di-acknowledge
- [ ] Alert history & resolution tracking

#### Alert Configuration
```php
// config/khazwal_alerts.php
return [
    'counting_overdue_minutes' => 120, // 2 hours
    'defect_critical_percentage' => 10,
    'waste_critical_percentage' => 5,
    'overtime_minutes' => 120,
];
```

#### Database Insert Reference
```sql
-- Create Alert
INSERT INTO alerts (
    alert_type,
    alert_category,
    severity,
    title,
    message,
    reference_type,
    reference_id,
    reference_number,
    target_roles,
    status,
    triggered_at,
    created_at
) VALUES (
    'COUNTING_OVERDUE',
    'KHAZWAL',
    'CRITICAL',
    'PO Menunggu Penghitungan Terlalu Lama',
    'PO ' || :po_number || ' sudah menunggu penghitungan selama ' || :waiting_hours || ' jam.',
    'PRODUCTION_ORDER',
    :po_id,
    :po_number,
    '["SUPERVISOR_KHAZWAL"]',
    'ACTIVE',
    NOW(),
    NOW()
);
```

---

### US-KW-018: Laporan Harian Khazanah Awal

| Field | Value |
|-------|-------|
| **ID** | US-KW-018 |
| **Story Points** | 8 |
| **Priority** | 🟡 Medium |
| **Dependencies** | US-KW-015, US-KW-016 |

**User Story:**
> Sebagai Supervisor Khazanah Awal, saya ingin generate laporan harian otomatis, sehingga saya tidak perlu buat laporan manual.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-018-BE-01 | Create `KhazwalReportService.php` | 3h | Backend |
| KW-018-BE-02 | Implement daily summary data collection | 2h | Backend |
| KW-018-BE-03 | Calculate all report metrics | 3h | Backend |
| KW-018-BE-04 | Create PDF generator (DomPDF/Snappy) | 4h | Backend |
| KW-018-BE-05 | Create Excel generator (Laravel Excel) | 3h | Backend |
| KW-018-BE-06 | Create scheduled job for end-of-shift | 2h | Backend |
| KW-018-BE-07 | Implement email delivery | 2h | Backend |
| KW-018-BE-08 | Create API `GET /api/khazwal/reports/daily` | 2h | Backend |
| KW-018-BE-09 | Create API `GET /api/khazwal/reports/download` | 2h | Backend |
| KW-018-FE-01 | Create report preview page | 3h | Frontend |
| KW-018-FE-02 | Add date range selector | 2h | Frontend |
| KW-018-FE-03 | Add download buttons (PDF/Excel) | 1h | Frontend |
| KW-018-FE-04 | Create report history list | 2h | Frontend |

#### Acceptance Criteria
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

#### Report Template Structure
```
╔════════════════════════════════════════════════════════════════╗
║           LAPORAN HARIAN KHAZANAH AWAL                         ║
║           Tanggal: [DATE] | Shift: [SHIFT]                     ║
╠════════════════════════════════════════════════════════════════╣
║ 1. RINGKASAN AKTIVITAS                                         ║
║    • Total PO Diproses: XX                                     ║
║    • Total PO Selesai: XX                                      ║
║    • Total PO Pending: XX                                      ║
╠════════════════════════════════════════════════════════════════╣
║ 2. PENYIAPAN MATERIAL                                          ║
║    • PO Selesai: XX                                            ║
║    • Rata-rata Durasi: XX menit (Target: 45 menit)             ║
║    • Akurasi Material: XX% (Target: 98%)                       ║
╠════════════════════════════════════════════════════════════════╣
║ 3. PENGHITUNGAN                                                ║
║    • PO Selesai: XX                                            ║
║    • Rata-rata Durasi: XX menit (Target: 30 menit)             ║
║    • Rata-rata Rusak: XX% (Target: ≤2%)                        ║
║    • Breakdown Kerusakan: [PIE CHART / TABLE]                  ║
╠════════════════════════════════════════════════════════════════╣
║ 4. PEMOTONGAN                                                  ║
║    • PO Selesai: XX                                            ║
║    • Rata-rata Durasi: XX menit (Target: 60 menit)             ║
║    • Rata-rata Waste: XX% (Target: ≤1%)                        ║
║    • Total Output: XX lembar kirim                             ║
╠════════════════════════════════════════════════════════════════╣
║ 5. PERFORMA STAFF                                              ║
║    [TABLE: Staff | PO Selesai | Durasi Avg | Score]            ║
╠════════════════════════════════════════════════════════════════╣
║ 6. ISSUE & RESOLUTION                                          ║
║    • [Issue description] - [Status]                            ║
╠════════════════════════════════════════════════════════════════╣
║ 7. REKOMENDASI                                                 ║
║    • [Action items based on data analysis]                     ║
╚════════════════════════════════════════════════════════════════╝
```

---

## 📊 Epic Summary

| User Story | Story Points | Priority | Phase |
|------------|--------------|----------|-------|
| US-KW-015 | 13 | Medium | MVP |
| US-KW-016 | 13 | Medium | Phase 2 |
| US-KW-017 | 8 | Medium | MVP |
| US-KW-018 | 8 | Medium | Phase 2 |
| **Total** | **42** | - | - |

---

## 🔗 Dependencies Graph

```
Epic 1, 2, 3 (Core Operations)
    │
    └── US-KW-015 (Dashboard Overview)
            │
            ├── US-KW-016 (Staff Performance)
            │       │
            │       └── US-KW-018 (Daily Report)
            │
            └── US-KW-017 (Alert System)
```

---

## 🧪 Testing Checklist

### Unit Tests
- [ ] KhazwalDashboardService - calculate stats
- [ ] KhazwalDashboardService - calculate trends
- [ ] AlertService - trigger conditions
- [ ] ReportService - data aggregation
- [ ] ReportService - PDF generation
- [ ] ReportService - Excel generation

### Integration Tests
- [ ] API endpoint: GET dashboard overview
- [ ] API endpoint: GET staff performance
- [ ] API endpoint: GET alerts
- [ ] API endpoint: POST acknowledge alert
- [ ] API endpoint: GET reports
- [ ] Scheduled job: alert checking
- [ ] Scheduled job: report generation

### E2E Tests
- [ ] Dashboard auto-refresh
- [ ] Alert notification flow
- [ ] Report download (PDF & Excel)

---

## 📱 UI/UX Notes

- **Dashboard:** Clean, scannable layout with clear visual hierarchy
- **Charts:** Use consistent color scheme (green = good, red = bad)
- **Alerts:** Non-intrusive but noticeable, priority-based sorting
- **Reports:** Easy to read, professional appearance
- **Auto-refresh:** Subtle indicator, no jarring page reloads

---

**Last Updated:** 27 December 2025  
**Status:** Ready for Development
