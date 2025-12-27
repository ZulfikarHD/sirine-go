# Epic 5: Analytics & Insights

**Epic ID:** KHAZWAL-EPIC-05  
**Priority:** 🟢 Low (Phase 3)  
**Estimated Duration:** 2 Minggu  

---

## 📋 Overview

Epic ini mencakup analisa mendalam untuk efisiensi pemotongan dan durasi proses, membantu identifikasi bottleneck dan best practices.

---

## 🗄️ Database Reference

### Primary Tables (Aggregation Source)
- `khazwal_cutting_results` - Waste data
- `khazwal_material_preparations` - Duration data
- `khazwal_counting_results` - Duration & defect data
- `print_job_summaries` - Machine & operator reference

### Related Tables
- `production_orders` - PO metadata
- `users` - Staff breakdown
- `machines` - Machine breakdown

---

## 📝 Backlog Items

### US-KW-019: Analisa Efisiensi Pemotongan

| Field | Value |
|-------|-------|
| **ID** | US-KW-019 |
| **Story Points** | 8 |
| **Priority** | 🟢 Low |
| **Dependencies** | Epic 3 (Cutting), US-KW-015 |

**User Story:**
> Sebagai Supervisor Khazanah Awal, saya ingin melihat analisa efisiensi pemotongan, sehingga bisa minimize waste.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-019-BE-01 | Create API `GET /api/khazwal/analytics/cutting-efficiency` | 3h | Backend |
| KW-019-BE-02 | Create waste trend query (time series) | 2h | Backend |
| KW-019-BE-03 | Create waste breakdown by machine | 2h | Backend |
| KW-019-BE-04 | Create waste breakdown by operator | 2h | Backend |
| KW-019-BE-05 | Create waste breakdown by shift | 2h | Backend |
| KW-019-BE-06 | Create waste breakdown by product type | 2h | Backend |
| KW-019-BE-07 | Calculate waste cost (Rp) | 2h | Backend |
| KW-019-BE-08 | Implement date filter | 1h | Backend |
| KW-019-FE-01 | Create page `CuttingEfficiencyPage.vue` | 4h | Frontend |
| KW-019-FE-02 | Create waste trend line chart | 2h | Frontend |
| KW-019-FE-03 | Create breakdown bar/pie charts | 3h | Frontend |
| KW-019-FE-04 | Create comparison table (operator vs operator) | 2h | Frontend |
| KW-019-FE-05 | Create waste cost summary card | 1h | Frontend |
| KW-019-FE-06 | Implement date range filter | 2h | Frontend |

#### Acceptance Criteria
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

#### Database Query Reference
```sql
-- Waste by Operator
SELECT 
    u.nip,
    u.full_name,
    COUNT(cut.id) as total_jobs,
    AVG(cut.waste_percentage) as avg_waste,
    SUM(cut.waste_quantity) as total_waste,
    SUM(cut.waste_quantity) * :cost_per_sheet as waste_cost
FROM khazwal_cutting_results cut
JOIN users u ON u.id = cut.cut_by
WHERE cut.completed_at BETWEEN :start_date AND :end_date
    AND cut.status = 'COMPLETED'
GROUP BY u.id, u.nip, u.full_name
ORDER BY avg_waste ASC;
```

---

### US-KW-020: Analisa Durasi Proses

| Field | Value |
|-------|-------|
| **ID** | US-KW-020 |
| **Story Points** | 8 |
| **Priority** | 🟢 Low |
| **Dependencies** | Epic 1, 2, 3, US-KW-015 |

**User Story:**
> Sebagai Supervisor Khazanah Awal, saya ingin melihat analisa durasi setiap tahap proses, sehingga bisa identifikasi bottleneck.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-020-BE-01 | Create API `GET /api/khazwal/analytics/duration` | 3h | Backend |
| KW-020-BE-02 | Calculate avg duration per stage | 2h | Backend |
| KW-020-BE-03 | Create duration breakdown by staff | 2h | Backend |
| KW-020-BE-04 | Create duration breakdown by shift | 2h | Backend |
| KW-020-BE-05 | Create duration breakdown by product type | 2h | Backend |
| KW-020-BE-06 | Create duration trend (time series) | 2h | Backend |
| KW-020-BE-07 | Calculate target vs actual comparison | 1h | Backend |
| KW-020-BE-08 | Identify bottleneck (longest stage) | 1h | Backend |
| KW-020-FE-01 | Create page `DurationAnalyticsPage.vue` | 4h | Frontend |
| KW-020-FE-02 | Create stacked bar chart (per stage) | 2h | Frontend |
| KW-020-FE-03 | Create breakdown charts | 3h | Frontend |
| KW-020-FE-04 | Create target vs actual comparison visual | 2h | Frontend |
| KW-020-FE-05 | Create bottleneck indicator | 1h | Frontend |
| KW-020-FE-06 | Implement date range filter | 2h | Frontend |

#### Acceptance Criteria
- [ ] Chart: Rata-rata durasi per tahap (Penyiapan, Penghitungan, Pemotongan)
- [ ] Breakdown durasi per: Staff, Shift, Jenis produk
- [ ] Identification bottleneck
- [ ] Trend durasi (time series)
- [ ] Target vs actual comparison
- [ ] Filter berdasarkan periode

#### KPI Targets Reference
```
Penyiapan Material: Target ≤ 45 menit
Penghitungan: Target ≤ 30 menit
Pemotongan: Target ≤ 60 menit
Total Cycle Time: Target ≤ 135 menit
```

#### Database Query Reference
```sql
-- Duration by Stage
SELECT 
    'Material Prep' as stage,
    AVG(duration_minutes) as avg_duration,
    MIN(duration_minutes) as min_duration,
    MAX(duration_minutes) as max_duration,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY duration_minutes) as median,
    45 as target
FROM khazwal_material_preparations
WHERE completed_at BETWEEN :start_date AND :end_date
    AND status = 'COMPLETED'

UNION ALL

SELECT 
    'Counting' as stage,
    AVG(duration_minutes),
    MIN(duration_minutes),
    MAX(duration_minutes),
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY duration_minutes),
    30 as target
FROM khazwal_counting_results
WHERE completed_at BETWEEN :start_date AND :end_date
    AND status = 'COMPLETED'

UNION ALL

SELECT 
    'Cutting' as stage,
    AVG(duration_minutes),
    MIN(duration_minutes),
    MAX(duration_minutes),
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY duration_minutes),
    60 as target
FROM khazwal_cutting_results
WHERE completed_at BETWEEN :start_date AND :end_date
    AND status = 'COMPLETED';
```

---

## 📊 Epic Summary

| User Story | Story Points | Priority | Phase |
|------------|--------------|----------|-------|
| US-KW-019 | 8 | Low | Phase 3 |
| US-KW-020 | 8 | Low | Phase 3 |
| **Total** | **16** | - | - |

---

## 🔗 Dependencies Graph

```
Epic 1, 2, 3 (Core Operations)
    │
    └── US-KW-015 (Dashboard)
            │
            ├── US-KW-019 (Cutting Efficiency)
            │
            └── US-KW-020 (Duration Analysis)
```

---

## 🧪 Testing Checklist

### Unit Tests
- [ ] Analytics calculations
- [ ] Waste cost calculation
- [ ] Bottleneck identification logic

### Integration Tests
- [ ] API endpoint: GET cutting efficiency
- [ ] API endpoint: GET duration analytics
- [ ] Date filter functionality

---

**Last Updated:** 27 December 2025  
**Status:** Ready for Development
