# Database Architecture - Sistem Produksi Pita Cukai

**Dibuat oleh:** Zulfikar Hidayatullah  
**Tech Stack:** Laravel (Service Pattern), PostgreSQL  
**Timezone:** Asia/Jakarta (WIB)  
**Currency:** Rupiah (Rp)  
**Language:** Indonesian  

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Key Architectural Decisions](#key-architectural-decisions)
3. [Core Tables](#core-tables)
4. [Master Data Tables](#master-data-tables)
5. [Transaction Tables](#transaction-tables)
6. [Tracking & Audit Tables](#tracking--audit-tables)
7. [Relationships](#relationships)
8. [Indexes & Performance](#indexes--performance)
9. [ERD](#erd)
10. [Implementation Notes](#implementation-notes)
11. [Migration Strategy](#migration-strategy)

---

## Key Architectural Decisions

### ❌ What's NOT in This Database

1. **NO `customers` table** - Handled by SAP (fetch via API)
2. **NO `product_specifications` table** - Handled by SAP (fetch via API)
3. **NO `materials` table** - Handled by SAP (fetch via API)
4. **NO `palet` concept** - Assumption: **1 Palet = 1 PO** (simplification)

### ✅ Key Design Changes

1. **`po_number`** - BIGINT (numbers only, not VARCHAR)
2. **`obc_number`** - VARCHAR(9) (exactly 9 characters)
3. **`username`** - Using NIP VARCHAR(5) (max 5 characters)
4. **Verification Logic** - Redesigned completely:
   - **PER LABEL/RIM** (not per PO)
   - Multiple QC Inspectors can work on same PO simultaneously
   - Realistic workload: 1 employee cannot QC 80 rims/day
   - 1 Label = ~500 lembar (1 rim), QC-able by 1 person per day

### 📊 Total Tables: **31 Tables**

**Breakdown:**
- Core: 2 tables
- Master: 2 tables  
- Khazanah Awal: 3 tables
- Cetak: 6 tables
- Verifikasi: 4 tables (completely redesigned)
- Khazanah Akhir: 7 tables (no palet concept)
- Tracking & Audit: 6 tables
- Maintenance: 1 table

---

## Overview

Sistem ini mengelola end-to-end production flow dari penerimaan PO sampai delivery ke customer (Dirjen Bea Cukai). Flow utama:

```
PO → Khazwal (Material Prep) → Cetak → Khazwal (Counting & Cutting) 
   → Verifikasi (QC per Label/Rim) → Khazkhir (Packaging) → Khazkhir (Warehouse) 
   → Khazkhir (Shipping) → Delivery
```

**Design Principles:**
- **Everything centered around PO & OBC** (no palet concept - assume 1 palet = 1 PO)
- **SAP Integration**: Customer, Product Specs, Materials handled by SAP (via API)
- Single Source of Truth untuk setiap entity
- Normalized (3NF) dengan strategic denormalization untuk performance
- Audit trail untuk semua critical data
- Soft deletes untuk data integrity
- Timezone-aware timestamps (Asia/Jakarta)
- Support untuk future scalability

**Key Changes from Initial Design:**
- ❌ NO `customers` table (handled by SAP)
- ❌ NO `product_specifications` table (handled by SAP)
- ❌ NO `materials` table (handled by SAP)
- ✅ `po_number` is BIGINT (numbers only)
- ✅ `obc_number` is VARCHAR(9) (9 characters)
- ✅ `username` using NIP VARCHAR(5) (max 5 char)
- ✅ Verification redesigned: **per LABEL/RIM** (multiple QC can work on same PO)
- ✅ No palet concept (1 palet = 1 PO)

---

## Core Tables

### 1. `production_orders` (PO)
Tabel utama yang menjadi central reference untuk seluruh proses produksi.

```sql
CREATE TABLE production_orders (
    id BIGSERIAL PRIMARY KEY,
    
    -- Identifiers
    po_number BIGINT NOT NULL UNIQUE, -- Pure number only
    obc_number VARCHAR(9) NOT NULL, -- 9 characters
    
    -- SAP References (get via API)
    sap_customer_code VARCHAR(50), -- Customer code dari SAP
    sap_product_code VARCHAR(50), -- Product code dari SAP
    
    -- Product Info (cached from SAP for quick access)
    product_name VARCHAR(255),
    product_specifications JSONB, -- Cache dari SAP (plat_code, colors, size, dll)
    
    -- Quantities
    quantity_ordered INTEGER NOT NULL, -- Jumlah order (lembar kirim)
    quantity_target_lembar_besar INTEGER NOT NULL, -- Target lembar besar (÷ 2 dari ordered)
    estimated_rims INTEGER NOT NULL, -- Estimasi jumlah rim (quantity_ordered ÷ 500)
    
    -- Dates
    order_date DATE NOT NULL,
    due_date DATE NOT NULL,
    
    -- Priority
    priority VARCHAR(20) NOT NULL DEFAULT 'NORMAL', -- URGENT, NORMAL, LOW
    priority_score INTEGER, -- Calculated score untuk sorting
    
    -- Status Tracking
    current_stage VARCHAR(50) NOT NULL, -- KHAZWAL_MATERIAL_PREP, CETAK, KHAZWAL_COUNTING, dll
    current_status VARCHAR(50) NOT NULL, -- Status spesifik per stage
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Audit
    created_by BIGINT,
    updated_by BIGINT,
    
    -- Indexes
    INDEX idx_po_number (po_number),
    INDEX idx_obc_number (obc_number),
    INDEX idx_current_stage (current_stage),
    INDEX idx_current_status (current_status),
    INDEX idx_due_date (due_date),
    INDEX idx_priority (priority, priority_score),
    INDEX idx_sap_customer (sap_customer_code),
    INDEX idx_sap_product (sap_product_code)
);
```

### 2. `po_stage_tracking`
Tracking perjalanan PO melalui setiap stage dengan timestamps.

```sql
CREATE TABLE po_stage_tracking (
    id BIGSERIAL PRIMARY KEY,
    production_order_id BIGINT NOT NULL,
    
    -- Stage Info
    stage VARCHAR(50) NOT NULL, -- KHAZWAL_MATERIAL_PREP, CETAK, dll
    status VARCHAR(50) NOT NULL, -- PENDING, IN_PROGRESS, COMPLETED, dll
    
    -- Timestamps
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER, -- Auto-calculated
    
    -- Assigned Staff
    assigned_to BIGINT, -- FK ke users
    
    -- Notes
    notes TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_po_stage (production_order_id, stage),
    INDEX idx_status (status),
    
    CONSTRAINT fk_production_order FOREIGN KEY (production_order_id) REFERENCES production_orders(id) ON DELETE CASCADE
);
```

---

## Master Data Tables

### 3. `users`
Semua user (staff, supervisor, manager).

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    
    -- Auth (using NIP)
    nip VARCHAR(5) NOT NULL UNIQUE, -- NIP max 5 characters (username)
    password VARCHAR(255) NOT NULL,
    
    -- Profile
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    
    -- Role & Department
    role VARCHAR(50) NOT NULL, -- STAFF_KHAZWAL, OPERATOR_CETAK, QC_INSPECTOR, dll
    department VARCHAR(50), -- KHAZWAL, CETAK, VERIFIKASI, KHAZKHIR
    
    -- Shift
    default_shift VARCHAR(20), -- PAGI, SIANG, MALAM
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Gamification
    total_points INTEGER DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    INDEX idx_nip (nip),
    INDEX idx_role (role),
    INDEX idx_department (department),
    INDEX idx_is_active (is_active)
);
```

### 4. `machines`
Master data mesin cetak.

```sql
CREATE TABLE machines (
    id BIGSERIAL PRIMARY KEY,
    
    -- Identifiers
    machine_code VARCHAR(50) NOT NULL UNIQUE, -- MC-01, MC-02, dll
    machine_name VARCHAR(255) NOT NULL,
    
    -- Specifications
    machine_type VARCHAR(50), -- Jenis mesin
    capacity_per_hour INTEGER, -- Kapasitas produksi per jam
    
    -- Status
    current_status VARCHAR(50) DEFAULT 'AVAILABLE', -- AVAILABLE, RUNNING, DOWN, SETUP, MAINTENANCE
    
    -- Maintenance
    last_pm_date DATE, -- Preventive Maintenance terakhir
    next_pm_date DATE,
    operating_hours_since_pm INTEGER DEFAULT 0,
    
    -- OEE Baseline
    target_oee_percentage DECIMAL(5,2) DEFAULT 85.00,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

---

## Transaction Tables

**Note:** Master data untuk Customer, Product Specifications, dan Materials **TIDAK** disimpan di database lokal. 
Data ini di-fetch dari SAP via API saat dibutuhkan, dengan optional caching di Redis untuk performance.

### KHAZANAH AWAL (KHAZWAL)

### 5. `khazwal_material_preparations`
Proses penyiapan material di awal.

```sql
CREATE TABLE khazwal_material_preparations (
    id BIGSERIAL PRIMARY KEY,
    production_order_id BIGINT NOT NULL UNIQUE, -- 1-to-1 dengan PO
    
    -- Material Requirements (SAP reference)
    sap_plat_code VARCHAR(50) NOT NULL, -- Reference ke SAP
    kertas_blanko_quantity INTEGER NOT NULL, -- Jumlah lembar blanko (target ÷ 2)
    
    -- Tinta (JSONB karena multiple colors)
    tinta_requirements JSONB, -- [{"sap_material_code": "TINTA-RED", "color": "Merah", "quantity_kg": 5}, ...]
    
    -- Actual Preparation
    plat_retrieved_at TIMESTAMP WITH TIME ZONE,
    kertas_blanko_actual INTEGER, -- Actual yang disiapkan
    kertas_blanko_variance INTEGER, -- Selisih dari target
    kertas_blanko_variance_reason TEXT,
    
    tinta_actual JSONB, -- Actual tinta yang disiapkan
    
    -- Photos
    material_photos JSONB, -- Array of photo URLs
    
    -- Status & Timing
    status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, IN_PROGRESS, COMPLETED
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER,
    
    -- Staff
    prepared_by BIGINT, -- FK ke users
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_po_khazwal_prep FOREIGN KEY (production_order_id) REFERENCES production_orders(id),
    CONSTRAINT fk_prepared_by FOREIGN KEY (prepared_by) REFERENCES users(id)
);
```

### 6. `khazwal_counting_results`
Penghitungan hasil cetak dari Unit Cetak.

```sql
CREATE TABLE khazwal_counting_results (
    id BIGSERIAL PRIMARY KEY,
    production_order_id BIGINT NOT NULL UNIQUE,
    
    -- Counting Results (Lembar Besar)
    quantity_good INTEGER NOT NULL, -- Lembar besar yang baik
    quantity_defect INTEGER NOT NULL, -- Lembar besar yang rusak
    total_counted INTEGER, -- Auto-calculated
    variance_from_target INTEGER,
    
    -- Percentages
    percentage_good DECIMAL(5,2),
    percentage_defect DECIMAL(5,2),
    
    -- Defect Breakdown (jika rusak > 5%)
    defect_breakdown JSONB, -- [{"type": "Warna pudar", "quantity": 10}, ...]
    
    -- Status & Timing
    status VARCHAR(50) DEFAULT 'PENDING',
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER,
    
    -- Staff
    counted_by BIGINT,
    
    -- Notes
    variance_reason TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_po_khazwal_count FOREIGN KEY (production_order_id) REFERENCES production_orders(id),
    CONSTRAINT fk_counted_by FOREIGN KEY (counted_by) REFERENCES users(id)
);
```

### 7. `khazwal_cutting_results`
Pemotongan lembar besar → lembar kirim.

```sql
CREATE TABLE khazwal_cutting_results (
    id BIGSERIAL PRIMARY KEY,
    production_order_id BIGINT NOT NULL UNIQUE,
    
    -- Input
    input_lembar_besar INTEGER NOT NULL, -- Dari good count
    
    -- Output (Lembar Kirim)
    output_sisiran_kiri INTEGER NOT NULL,
    output_sisiran_kanan INTEGER NOT NULL,
    total_output INTEGER, -- Auto-calculated
    
    -- Expected vs Actual
    expected_output INTEGER, -- Input × 2
    waste_quantity INTEGER,
    waste_percentage DECIMAL(5,2),
    waste_reason TEXT,
    waste_photo_url VARCHAR(255),
    
    -- Machine & Staff
    cutting_machine VARCHAR(50),
    cut_by BIGINT,
    
    -- Status & Timing
    status VARCHAR(50) DEFAULT 'PENDING',
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_po_khazwal_cut FOREIGN KEY (production_order_id) REFERENCES production_orders(id),
    CONSTRAINT fk_cut_by FOREIGN KEY (cut_by) REFERENCES users(id)
);
```

### CETAK (PRINT)

### 8. `print_jobs`
Assignment PO ke mesin cetak.

```sql
CREATE TABLE print_jobs (
    id BIGSERIAL PRIMARY KEY,
    production_order_id BIGINT NOT NULL,
    machine_id BIGINT NOT NULL,
    
    -- Assignment
    assigned_to BIGINT NOT NULL, -- Operator
    assigned_by BIGINT, -- Supervisor
    assigned_at TIMESTAMP WITH TIME ZONE,
    
    -- Setup
    setup_started_at TIMESTAMP WITH TIME ZONE,
    setup_completed_at TIMESTAMP WITH TIME ZONE,
    setup_duration_minutes INTEGER,
    
    -- Setup Checklist
    setup_checklist JSONB, -- {"plat_installed": true, "tinta_ready": true, ...}
    
    -- Test Print
    test_print_quantity INTEGER,
    test_print_approved BOOLEAN,
    test_print_approved_at TIMESTAMP WITH TIME ZONE,
    test_print_photos JSONB, -- Array of URLs
    
    -- Production
    production_started_at TIMESTAMP WITH TIME ZONE,
    production_completed_at TIMESTAMP WITH TIME ZONE,
    production_duration_minutes INTEGER,
    
    -- Status
    status VARCHAR(50) DEFAULT 'ASSIGNED', -- ASSIGNED, SETUP, TEST_PRINT, RUNNING, COMPLETED
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_machine_status (machine_id, status),
    INDEX idx_operator (assigned_to),
    
    CONSTRAINT fk_po_print FOREIGN KEY (production_order_id) REFERENCES production_orders(id),
    CONSTRAINT fk_machine FOREIGN KEY (machine_id) REFERENCES machines(id),
    CONSTRAINT fk_operator FOREIGN KEY (assigned_to) REFERENCES users(id)
);
```

### 9. `print_production_logs`
Real-time production tracking (counter updates).

```sql
CREATE TABLE print_production_logs (
    id BIGSERIAL PRIMARY KEY,
    print_job_id BIGINT NOT NULL,
    
    -- Counter Update
    counter_value INTEGER NOT NULL, -- Cumulative output
    counter_increment INTEGER, -- Increment dari update sebelumnya
    
    -- Timestamp
    logged_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Auto-calculated metrics
    current_rate_per_hour INTEGER, -- Production rate saat ini
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_print_job_time (print_job_id, logged_at),
    
    CONSTRAINT fk_print_job_log FOREIGN KEY (print_job_id) REFERENCES print_jobs(id) ON DELETE CASCADE
);
```

### 10. `print_downtimes`
Recording downtime mesin.

```sql
CREATE TABLE print_downtimes (
    id BIGSERIAL PRIMARY KEY,
    print_job_id BIGINT NOT NULL,
    machine_id BIGINT NOT NULL,
    
    -- Downtime
    downtime_reason VARCHAR(50) NOT NULL, -- BREAKDOWN, WAITING_MATERIAL, MINOR_STOP, dll
    downtime_category VARCHAR(50), -- PLANNED, UNPLANNED
    downtime_description TEXT,
    
    -- Timing
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ended_at TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER,
    
    -- For Breakdown
    maintenance_ticket_id BIGINT, -- FK ke maintenance_tickets (jika breakdown)
    
    -- Logged by
    logged_by BIGINT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_print_job_downtime (print_job_id),
    INDEX idx_machine_downtime (machine_id, started_at),
    
    CONSTRAINT fk_print_job_downtime FOREIGN KEY (print_job_id) REFERENCES print_jobs(id),
    CONSTRAINT fk_machine_downtime FOREIGN KEY (machine_id) REFERENCES machines(id)
);
```

### 11. `print_quality_issues`
Recording quality issues saat produksi.

```sql
CREATE TABLE print_quality_issues (
    id BIGSERIAL PRIMARY KEY,
    print_job_id BIGINT NOT NULL,
    
    -- Issue Details
    defect_type VARCHAR(100) NOT NULL,
    defect_quantity INTEGER NOT NULL,
    defect_severity VARCHAR(20), -- CRITICAL, MAJOR, MINOR
    
    -- Root Cause (estimasi)
    estimated_cause VARCHAR(100),
    
    -- Documentation
    defect_photos JSONB, -- Array of URLs
    notes TEXT,
    
    -- Action Taken
    action_taken TEXT,
    
    -- Timing
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Logged by
    logged_by BIGINT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_print_job_quality (print_job_id),
    INDEX idx_defect_type (defect_type),
    
    CONSTRAINT fk_print_job_quality FOREIGN KEY (print_job_id) REFERENCES print_jobs(id)
);
```

### 12. `print_job_summaries`
Summary akhir setelah produksi selesai (untuk OEE calculation).

```sql
CREATE TABLE print_job_summaries (
    id BIGSERIAL PRIMARY KEY,
    print_job_id BIGINT NOT NULL UNIQUE,
    production_order_id BIGINT NOT NULL,
    machine_id BIGINT NOT NULL,
    operator_id BIGINT NOT NULL,
    
    -- Quantities
    target_quantity INTEGER NOT NULL,
    actual_quantity INTEGER NOT NULL,
    good_quantity INTEGER NOT NULL,
    defect_quantity INTEGER NOT NULL,
    
    -- Times (minutes)
    total_time_minutes INTEGER, -- Production end - start
    actual_production_time_minutes INTEGER, -- Exclude downtime
    downtime_minutes INTEGER,
    setup_time_minutes INTEGER,
    
    -- OEE Components
    availability_percentage DECIMAL(5,2),
    performance_percentage DECIMAL(5,2),
    quality_percentage DECIMAL(5,2),
    oee_percentage DECIMAL(5,2),
    
    -- Breakdown
    downtime_breakdown JSONB, -- {"breakdown": 45, "waiting": 10, ...}
    
    -- Photos (no palet concept - 1 palet = 1 PO)
    result_photos JSONB,
    
    -- Handover Notes (jika shift handover)
    handover_notes TEXT,
    handover_to BIGINT,
    
    -- Finalized
    finalized_at TIMESTAMP WITH TIME ZONE,
    finalized_by BIGINT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_machine_date (machine_id, finalized_at),
    INDEX idx_operator_date (operator_id, finalized_at),
    
    CONSTRAINT fk_print_job_summary FOREIGN KEY (print_job_id) REFERENCES print_jobs(id),
    CONSTRAINT fk_po_summary FOREIGN KEY (production_order_id) REFERENCES production_orders(id),
    CONSTRAINT fk_machine_summary FOREIGN KEY (machine_id) REFERENCES machines(id),
    CONSTRAINT fk_operator_summary FOREIGN KEY (operator_id) REFERENCES users(id)
);
```

### VERIFIKASI (QC)

**IMPORTANT REDESIGN:** Verification is per LABEL/RIM, NOT per PO. 
- 1 PO = ~80 rims (depending on quantity)
- 1 QC Inspector handles 1 LABEL/RIM at a time
- Multiple QC Inspectors can work on the same PO simultaneously (parallel processing)
- This makes it realistic: 1 employee cannot QC entire PO (80 rims) in a day

### 13. `verification_labels`
Master list of labels/rims yang perlu di-QC (auto-generated from cutting results).

```sql
CREATE TABLE verification_labels (
    id BIGSERIAL PRIMARY KEY,
    production_order_id BIGINT NOT NULL,
    
    -- Label Info
    label_number INTEGER NOT NULL, -- 1, 2, 3, ... (same as rim number)
    total_labels INTEGER NOT NULL, -- Total labels untuk PO ini
    
    -- Quantity (per label/rim)
    target_quantity INTEGER NOT NULL, -- ~500 lembar per label (or less untuk label terakhir)
    sisiran VARCHAR(10), -- KIRI or KANAN
    
    -- Source Info (dari cutting)
    cutting_result_id BIGINT,
    
    -- QC Status
    qc_status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, ASSIGNED, IN_PROGRESS, COMPLETED
    
    -- Assignment
    assigned_to BIGINT, -- QC Inspector yang currently handle label ini
    assigned_at TIMESTAMP WITH TIME ZONE,
    
    -- Timing
    qc_started_at TIMESTAMP WITH TIME ZONE,
    qc_completed_at TIMESTAMP WITH TIME ZONE,
    qc_duration_minutes INTEGER,
    
    -- Results (summary - detail di verification_inspection_results)
    hcs_quantity INTEGER,
    hcts_quantity INTEGER,
    hcs_percentage DECIMAL(5,2),
    
    -- Storage Location (after QC)
    storage_location VARCHAR(100),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_po_labels (production_order_id, label_number),
    INDEX idx_qc_status (qc_status),
    INDEX idx_assigned_qc (assigned_to, qc_status),
    
    CONSTRAINT fk_po_verification_label FOREIGN KEY (production_order_id) REFERENCES production_orders(id),
    CONSTRAINT fk_qc_inspector_label FOREIGN KEY (assigned_to) REFERENCES users(id),
    
    UNIQUE (production_order_id, label_number)
);
```

### 14. `verification_inspection_results`
Final inspection result per label (after QC completed).

```sql
CREATE TABLE verification_inspection_results (
    id BIGSERIAL PRIMARY KEY,
    verification_label_id BIGINT NOT NULL UNIQUE, -- 1-to-1 dengan label
    production_order_id BIGINT NOT NULL,
    
    -- QC Inspector
    qc_inspector_id BIGINT NOT NULL,
    
    -- Results
    target_quantity INTEGER NOT NULL,
    total_inspected INTEGER NOT NULL,
    hcs_quantity INTEGER NOT NULL,
    hcts_quantity INTEGER NOT NULL,
    
    -- Percentages
    hcs_percentage DECIMAL(5,2),
    hcts_percentage DECIMAL(5,2),
    
    -- Quality Score
    quality_score INTEGER, -- 0-100
    
    -- Performance Metrics
    inspection_duration_minutes INTEGER,
    inspection_rate_per_hour INTEGER, -- Calculated: total_inspected / (duration_minutes / 60)
    
    -- Defect Summary
    total_defects INTEGER,
    defect_breakdown JSONB, -- [{"type": "Warna pudar", "quantity": 5}, ...]
    
    -- Storage
    storage_location VARCHAR(100),
    
    -- Photos
    inspection_photos JSONB,
    
    -- Finalized
    finalized_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_po_inspection (production_order_id),
    INDEX idx_qc_inspector_result (qc_inspector_id, finalized_at),
    INDEX idx_hcs_percentage_result (hcs_percentage),
    
    CONSTRAINT fk_verification_label_result FOREIGN KEY (verification_label_id) REFERENCES verification_labels(id),
    CONSTRAINT fk_po_inspection FOREIGN KEY (production_order_id) REFERENCES production_orders(id),
    CONSTRAINT fk_qc_inspector_result FOREIGN KEY (qc_inspector_id) REFERENCES users(id)
);
```

### 15. `verification_defects`
Detail defect yang ditemukan saat verifikasi (per label).

```sql
CREATE TABLE verification_defects (
    id BIGSERIAL PRIMARY KEY,
    verification_label_id BIGINT NOT NULL,
    production_order_id BIGINT NOT NULL,
    
    -- Defect Details
    defect_type VARCHAR(100) NOT NULL,
    defect_quantity INTEGER NOT NULL,
    defect_severity VARCHAR(20), -- CRITICAL, MAJOR, MINOR
    
    -- Location
    defect_location VARCHAR(50), -- Visual map location
    
    -- Root Cause Analysis
    estimated_cause VARCHAR(100), -- Masalah plat, tinta, mesin, dll
    
    -- Documentation
    defect_photos JSONB,
    notes TEXT,
    
    -- Batch Info (range dalam label ini)
    batch_range_start INTEGER,
    batch_range_end INTEGER,
    
    -- Timing
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Logged by
    logged_by BIGINT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_verification_label_defect (verification_label_id),
    INDEX idx_po_defect (production_order_id),
    INDEX idx_defect_type_qc (defect_type),
    INDEX idx_severity (defect_severity),
    
    CONSTRAINT fk_verification_label_defect FOREIGN KEY (verification_label_id) REFERENCES verification_labels(id),
    CONSTRAINT fk_po_defect_verification FOREIGN KEY (production_order_id) REFERENCES production_orders(id)
);
```

### 16. `verification_po_summaries`
Aggregated summary per PO (all labels combined) - auto-calculated from label results.

```sql
CREATE TABLE verification_po_summaries (
    id BIGSERIAL PRIMARY KEY,
    production_order_id BIGINT NOT NULL UNIQUE,
    
    -- Overall Results (sum of all labels)
    total_labels INTEGER NOT NULL,
    labels_completed INTEGER NOT NULL,
    labels_pending INTEGER NOT NULL,
    
    target_quantity INTEGER NOT NULL,
    total_inspected INTEGER NOT NULL,
    hcs_quantity INTEGER NOT NULL,
    hcts_quantity INTEGER NOT NULL,
    
    -- Overall Percentages
    hcs_percentage DECIMAL(5,2),
    hcts_percentage DECIMAL(5,2),
    completion_percentage DECIMAL(5,2), -- (labels_completed / total_labels) × 100
    
    -- Average Quality Score (across all labels)
    avg_quality_score INTEGER,
    
    -- Performance (aggregate)
    total_duration_minutes INTEGER,
    avg_inspection_rate_per_hour INTEGER,
    
    -- Defect Summary (aggregate)
    total_defects INTEGER,
    defect_breakdown JSONB, -- Aggregated defect types
    
    -- QC Inspectors Involved
    qc_inspectors JSONB, -- [{"user_id": 1, "labels_count": 20}, ...]
    
    -- Status
    verification_status VARCHAR(50) DEFAULT 'IN_PROGRESS', -- IN_PROGRESS, COMPLETED
    
    -- Completion
    all_labels_completed_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_verification_status (verification_status),
    INDEX idx_completion_date (all_labels_completed_at),
    
    CONSTRAINT fk_po_verification_summary FOREIGN KEY (production_order_id) REFERENCES production_orders(id)
);
```

### KHAZANAH AKHIR (KHAZKHIR)

**Note:** No palet concept in Khazanah Akhir. Assumption: **1 Palet = 1 PO**. 
All tracking is centered around PO & OBC only.

### 17. `khazkhir_packaging_jobs`
Penghitungan final dan pengemasan.

```sql
CREATE TABLE khazkhir_packaging_jobs (
    id BIGSERIAL PRIMARY KEY,
    production_order_id BIGINT NOT NULL UNIQUE,
    
    -- Assignment
    assigned_to BIGINT NOT NULL, -- Staff Khazkhir
    assigned_at TIMESTAMP WITH TIME ZONE,
    
    -- HCS from Verification
    target_hcs_quantity INTEGER NOT NULL, -- Dari verification
    actual_hcs_counted INTEGER, -- Final count
    variance INTEGER,
    variance_reason TEXT,
    
    -- Damage during handling
    damaged_quantity INTEGER DEFAULT 0,
    damage_reason TEXT,
    damage_photos JSONB,
    
    -- Packaging
    total_rims INTEGER, -- Auto-calculated: actual ÷ 500
    remaining_sheets INTEGER, -- Sisa jika tidak genap 500
    
    -- Status & Timing
    status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, IN_PROGRESS, COMPLETED
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_po_packaging FOREIGN KEY (production_order_id) REFERENCES production_orders(id),
    CONSTRAINT fk_packaged_by FOREIGN KEY (assigned_to) REFERENCES users(id)
);
```

### 18. `khazkhir_rims`
Detail setiap rim yang dikemas (dengan label info).

```sql
CREATE TABLE khazkhir_rims (
    id BIGSERIAL PRIMARY KEY,
    packaging_job_id BIGINT NOT NULL,
    production_order_id BIGINT NOT NULL,
    
    -- Rim Info
    rim_number INTEGER NOT NULL, -- 1, 2, 3, ...
    total_rims INTEGER NOT NULL, -- Total untuk PO ini
    rim_quantity INTEGER NOT NULL, -- 500 or less (untuk rim terakhir)
    
    -- Label Info (auto-generated)
    rim_label_data JSONB, -- All label fields
    rim_qr_code VARCHAR(255), -- QR Code URL/data
    rim_barcode VARCHAR(255), -- Barcode URL/data
    
    -- Dates (untuk label)
    production_date DATE,
    verification_date DATE,
    packaging_date DATE,
    
    -- No palet concept - all rims belong to PO directly
    
    -- Photos
    rim_photo_url VARCHAR(255),
    
    -- Status
    is_packed BOOLEAN DEFAULT FALSE,
    packed_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_packaging_job_rim (packaging_job_id),
    
    CONSTRAINT fk_packaging_job_rim FOREIGN KEY (packaging_job_id) REFERENCES khazkhir_packaging_jobs(id),
    CONSTRAINT fk_po_rim FOREIGN KEY (production_order_id) REFERENCES production_orders(id),
    
    UNIQUE (packaging_job_id, rim_number)
);
```

### 19. `khazkhir_inventory_movements`
Tracking pergerakan inventory (stock in/out).

```sql
CREATE TABLE khazkhir_inventory_movements (
    id BIGSERIAL PRIMARY KEY,
    production_order_id BIGINT NOT NULL,
    
    -- Movement Type
    movement_type VARCHAR(50) NOT NULL, -- STOCK_IN, STOCK_OUT, ADJUSTMENT
    
    -- Quantity
    quantity_rims INTEGER NOT NULL,
    quantity_sheets INTEGER NOT NULL,
    
    -- Stock Balance (after this movement)
    balance_rims INTEGER,
    balance_sheets INTEGER,
    
    -- Location
    from_location VARCHAR(100),
    to_location VARCHAR(100),
    
    -- Reference
    reference_type VARCHAR(50), -- PACKAGING, PICKING, ADJUSTMENT
    reference_id BIGINT,
    
    -- User
    moved_by BIGINT,
    
    -- Timestamp
    movement_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Notes
    notes TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_po_inventory (production_order_id),
    INDEX idx_movement_type (movement_type),
    INDEX idx_movement_date (movement_date),
    
    CONSTRAINT fk_po_inventory FOREIGN KEY (production_order_id) REFERENCES production_orders(id)
);
```

### 20. `khazkhir_stock_opnames`
Stock opname (inventory count).

```sql
CREATE TABLE khazkhir_stock_opnames (
    id BIGSERIAL PRIMARY KEY,
    
    -- Opname Info
    opname_number VARCHAR(50) NOT NULL UNIQUE, -- SO-YYYY-MM-XXX
    opname_type VARCHAR(50) NOT NULL, -- FULL, PARTIAL, CYCLE_COUNT
    opname_scope VARCHAR(100), -- ALL, ZONA_A, PO-XXX, dll
    
    -- Status
    status VARCHAR(50) DEFAULT 'IN_PROGRESS', -- IN_PROGRESS, COMPLETED, APPROVED
    
    -- Timing
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    
    -- Performed by
    performed_by JSONB, -- Array of user IDs
    supervised_by BIGINT,
    
    -- Results
    total_pos_expected INTEGER,
    total_pos_counted INTEGER,
    total_discrepancies INTEGER,
    
    -- Approval
    approved_by BIGINT,
    approved_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_opname_status (status),
    INDEX idx_opname_date (started_at)
);
```

### 21. `khazkhir_stock_opname_details`
Detail stock opname per PO (no palet).

```sql
CREATE TABLE khazkhir_stock_opname_details (
    id BIGSERIAL PRIMARY KEY,
    stock_opname_id BIGINT NOT NULL,
    production_order_id BIGINT NOT NULL,
    
    -- Expected (from system)
    expected_quantity_rims INTEGER NOT NULL,
    expected_location VARCHAR(100),
    
    -- Actual (counted)
    actual_quantity_rims INTEGER,
    actual_location VARCHAR(100),
    
    -- Variance
    variance_rims INTEGER,
    variance_percentage DECIMAL(5,2),
    
    -- Status
    is_counted BOOLEAN DEFAULT FALSE,
    counted_at TIMESTAMP WITH TIME ZONE,
    counted_by BIGINT,
    
    -- Notes & Photos
    notes TEXT,
    photos JSONB,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_stock_opname_po (stock_opname_id, production_order_id),
    INDEX idx_variance (variance_percentage),
    
    CONSTRAINT fk_stock_opname FOREIGN KEY (stock_opname_id) REFERENCES khazkhir_stock_opnames(id),
    CONSTRAINT fk_po_opname FOREIGN KEY (production_order_id) REFERENCES production_orders(id)
);
```

### 22. `khazkhir_delivery_orders`
Delivery Order untuk pengiriman.

```sql
CREATE TABLE khazkhir_delivery_orders (
    id BIGSERIAL PRIMARY KEY,
    
    -- DO Info
    do_number VARCHAR(50) NOT NULL UNIQUE, -- DO-YYYY-MM-XXXX
    production_order_id BIGINT NOT NULL,
    sap_customer_code VARCHAR(50) NOT NULL, -- Customer from SAP
    
    -- Shipping Details
    shipping_address TEXT NOT NULL,
    shipping_city VARCHAR(100),
    shipping_province VARCHAR(100),
    contact_person VARCHAR(255),
    contact_phone VARCHAR(20),
    
    -- Dates
    delivery_date DATE NOT NULL,
    estimated_arrival TIMESTAMP WITH TIME ZONE,
    
    -- Items (no palet concept)
    total_rims INTEGER NOT NULL,
    total_sheets INTEGER NOT NULL,
    
    -- Logistics
    vehicle_type VARCHAR(50),
    vehicle_number VARCHAR(50),
    driver_name VARCHAR(255),
    driver_phone VARCHAR(20),
    ekspedisi VARCHAR(100), -- Jika 3rd party
    ekspedisi_tracking_number VARCHAR(100),
    
    -- Documents (auto-generated)
    surat_jalan_url VARCHAR(255),
    packing_list_url VARCHAR(255),
    bast_template_url VARCHAR(255),
    invoice_url VARCHAR(255),
    
    -- Status
    status VARCHAR(50) DEFAULT 'DRAFT', -- DRAFT, READY_TO_PICK, PICKED, LOADED, DISPATCHED, IN_TRANSIT, DELIVERED
    
    -- Timestamps per status
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    finalized_at TIMESTAMP WITH TIME ZONE,
    picking_started_at TIMESTAMP WITH TIME ZONE,
    picking_completed_at TIMESTAMP WITH TIME ZONE,
    loading_started_at TIMESTAMP WITH TIME ZONE,
    loading_completed_at TIMESTAMP WITH TIME ZONE,
    dispatched_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    
    -- Users
    created_by BIGINT,
    finalized_by BIGINT,
    dispatched_by BIGINT,
    
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_do_number (do_number),
    INDEX idx_po_do (production_order_id),
    INDEX idx_status_do (status),
    INDEX idx_delivery_date (delivery_date),
    INDEX idx_sap_customer_do (sap_customer_code),
    
    CONSTRAINT fk_po_do FOREIGN KEY (production_order_id) REFERENCES production_orders(id)
);
```

### 23. `khazkhir_delivery_tracking`
Tracking perjalanan delivery.

```sql
CREATE TABLE khazkhir_delivery_tracking (
    id BIGSERIAL PRIMARY KEY,
    delivery_order_id BIGINT NOT NULL,
    
    -- Status Update
    status VARCHAR(50) NOT NULL,
    status_description TEXT,
    
    -- Location (jika ada GPS)
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    location_name VARCHAR(255),
    
    -- Timestamp
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Updated by (driver or system)
    updated_by BIGINT,
    update_source VARCHAR(50), -- MANUAL, GPS, SYSTEM
    
    -- Notes
    notes TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_do_tracking (delivery_order_id, timestamp),
    
    CONSTRAINT fk_do_tracking FOREIGN KEY (delivery_order_id) REFERENCES khazkhir_delivery_orders(id)
);
```

### 24. `khazkhir_proof_of_deliveries`
Proof of Delivery (POD) digital.

```sql
CREATE TABLE khazkhir_proof_of_deliveries (
    id BIGSERIAL PRIMARY KEY,
    delivery_order_id BIGINT NOT NULL UNIQUE,
    
    -- Arrival
    arrival_timestamp TIMESTAMP WITH TIME ZONE,
    arrival_latitude DECIMAL(10,8),
    arrival_longitude DECIMAL(11,8),
    
    -- BAST Data
    bast_date DATE NOT NULL,
    rims_received INTEGER NOT NULL,
    sheets_received INTEGER NOT NULL,
    condition VARCHAR(50) DEFAULT 'BAIK', -- BAIK, RUSAK_SEBAGIAN, RUSAK
    condition_notes TEXT,
    damage_photos JSONB,
    
    -- Receiver Info
    receiver_name VARCHAR(255) NOT NULL,
    receiver_nik VARCHAR(50),
    receiver_phone VARCHAR(20),
    
    -- Signatures (base64 or URL)
    receiver_signature TEXT, -- Digital signature
    driver_signature TEXT,
    
    -- Timestamp Signatures
    receiver_signed_at TIMESTAMP WITH TIME ZONE,
    driver_signed_at TIMESTAMP WITH TIME ZONE,
    
    -- Photos
    delivery_photos JSONB, -- Photos of delivered goods
    location_photos JSONB, -- Photos of delivery location
    
    -- POD Document
    pod_document_url VARCHAR(255), -- Generated PDF
    
    -- Finalized
    finalized_at TIMESTAMP WITH TIME ZONE,
    finalized_by BIGINT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_do_pod FOREIGN KEY (delivery_order_id) REFERENCES khazkhir_delivery_orders(id)
);
```

---

## Tracking & Audit Tables

### 25. `activity_logs` (Audit Trail)
Comprehensive audit trail untuk semua critical actions.

```sql
CREATE TABLE activity_logs (
    id BIGSERIAL PRIMARY KEY,
    
    -- Who
    user_id BIGINT,
    user_role VARCHAR(50),
    user_name VARCHAR(255),
    
    -- What
    action VARCHAR(100) NOT NULL, -- CREATE, UPDATE, DELETE, APPROVE, FINALIZE, dll
    entity_type VARCHAR(100) NOT NULL, -- production_orders, print_jobs, dll
    entity_id BIGINT,
    
    -- Details
    description TEXT,
    changes JSONB, -- {"before": {...}, "after": {...}}
    
    -- Where
    ip_address VARCHAR(45),
    user_agent TEXT,
    device_type VARCHAR(50),
    
    -- When
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Reference
    reference_type VARCHAR(50), -- PO, DO, Print Job, dll
    reference_number VARCHAR(100),
    
    INDEX idx_user_activity (user_id, timestamp),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_timestamp (timestamp),
    INDEX idx_action (action)
);
```

### 26. `notifications`
Central notification system.

```sql
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    
    -- Recipient
    user_id BIGINT NOT NULL,
    
    -- Type
    notification_type VARCHAR(50) NOT NULL, -- ALERT, INFO, SUCCESS, WARNING
    notification_category VARCHAR(50), -- PRODUCTION, QUALITY, DELIVERY, dll
    
    -- Content
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    
    -- Link/Action
    action_url VARCHAR(255),
    action_text VARCHAR(100),
    
    -- Reference
    reference_type VARCHAR(50),
    reference_id BIGINT,
    
    -- Status
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,
    
    -- Priority
    priority VARCHAR(20) DEFAULT 'NORMAL', -- CRITICAL, HIGH, NORMAL, LOW
    
    -- Channels
    sent_via JSONB, -- ["in_app", "email", "whatsapp"]
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    
    INDEX idx_user_notifications (user_id, is_read, created_at),
    INDEX idx_priority (priority),
    
    CONSTRAINT fk_user_notification FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### 27. `gamification_achievements`
Master data achievements/badges.

```sql
CREATE TABLE gamification_achievements (
    id BIGSERIAL PRIMARY KEY,
    
    -- Achievement Info
    achievement_code VARCHAR(50) NOT NULL UNIQUE,
    achievement_name VARCHAR(255) NOT NULL,
    achievement_description TEXT,
    
    -- Icon/Badge
    icon_url VARCHAR(255),
    badge_rarity VARCHAR(20), -- BRONZE, SILVER, GOLD, DIAMOND
    
    -- Category
    category VARCHAR(50), -- QUALITY, PRODUCTIVITY, CONSISTENCY, dll
    department VARCHAR(50), -- CETAK, VERIFIKASI, dll
    
    -- Criteria
    criteria JSONB, -- Criteria untuk earn achievement
    
    -- Points
    points_reward INTEGER DEFAULT 0,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### 28. `user_achievements`
Achievements yang sudah di-earn user.

```sql
CREATE TABLE user_achievements (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    achievement_id BIGINT NOT NULL,
    
    -- Earned
    earned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Progress (untuk achievements dengan progress tracking)
    progress_current INTEGER DEFAULT 0,
    progress_target INTEGER,
    
    -- Reference (what triggered this achievement)
    reference_type VARCHAR(50),
    reference_id BIGINT,
    
    INDEX idx_user_achievements (user_id, earned_at),
    INDEX idx_achievement (achievement_id),
    
    CONSTRAINT fk_user_achieve FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_achievement FOREIGN KEY (achievement_id) REFERENCES gamification_achievements(id),
    
    UNIQUE (user_id, achievement_id)
);
```

### 29. `point_transactions`
History transaksi points (earn & redeem).

```sql
CREATE TABLE point_transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    
    -- Transaction Type
    transaction_type VARCHAR(50) NOT NULL, -- EARN, REDEEM, ADJUSTMENT
    
    -- Amount
    points_amount INTEGER NOT NULL, -- Positive untuk earn, negative untuk redeem
    points_balance_after INTEGER NOT NULL,
    
    -- Reason
    reason VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Reference
    reference_type VARCHAR(50),
    reference_id BIGINT,
    
    -- Approved by (untuk redemption)
    approved_by BIGINT,
    approved_at TIMESTAMP WITH TIME ZONE,
    
    -- Timestamp
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_user_points (user_id, transaction_date),
    INDEX idx_transaction_type (transaction_type),
    
    CONSTRAINT fk_user_points FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### 30. `alerts`
Central alert management.

```sql
CREATE TABLE alerts (
    id BIGSERIAL PRIMARY KEY,
    
    -- Alert Type
    alert_type VARCHAR(50) NOT NULL,
    alert_category VARCHAR(50),
    severity VARCHAR(20) NOT NULL, -- CRITICAL, WARNING, INFO
    
    -- Content
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    
    -- Reference
    reference_type VARCHAR(50),
    reference_id BIGINT,
    reference_number VARCHAR(100),
    
    -- Target (who should be alerted)
    target_roles JSONB, -- Array of roles
    target_users JSONB, -- Array of user IDs (jika specific users)
    
    -- Status
    status VARCHAR(50) DEFAULT 'ACTIVE', -- ACTIVE, ACKNOWLEDGED, RESOLVED, EXPIRED
    
    -- Acknowledgment
    acknowledged_by BIGINT,
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    
    -- Resolution
    resolved_by BIGINT,
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolution_notes TEXT,
    
    -- Response Time (auto-calculated)
    response_time_minutes INTEGER,
    resolution_time_minutes INTEGER,
    
    -- Timestamps
    triggered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_status_alert (status, triggered_at),
    INDEX idx_severity (severity),
    INDEX idx_reference (reference_type, reference_id)
);
```

---

## Maintenance Tables

### 31. `maintenance_tickets`
Breakdown & maintenance tickets.

```sql
CREATE TABLE maintenance_tickets (
    id BIGSERIAL PRIMARY KEY,
    
    -- Ticket Info
    ticket_number VARCHAR(50) NOT NULL UNIQUE, -- MT-YYYY-MM-XXXX
    ticket_type VARCHAR(50) NOT NULL, -- BREAKDOWN, PREVENTIVE, CORRECTIVE
    
    -- Machine
    machine_id BIGINT NOT NULL,
    
    -- Problem
    problem_description TEXT NOT NULL,
    problem_severity VARCHAR(20) NOT NULL, -- CRITICAL, HIGH, MEDIUM, LOW
    
    -- Reported
    reported_by BIGINT NOT NULL,
    reported_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Assignment
    assigned_to BIGINT, -- Technician
    assigned_at TIMESTAMP WITH TIME ZONE,
    
    -- Work
    work_started_at TIMESTAMP WITH TIME ZONE,
    work_completed_at TIMESTAMP WITH TIME ZONE,
    
    -- Root Cause & Action
    root_cause TEXT,
    action_taken TEXT,
    parts_used JSONB, -- [{"material_id": 1, "quantity": 2}, ...]
    
    -- Photos
    problem_photos JSONB,
    repair_photos JSONB,
    
    -- Status
    status VARCHAR(50) DEFAULT 'OPEN', -- OPEN, ASSIGNED, IN_PROGRESS, COMPLETED, CLOSED
    
    -- Downtime (link to print_downtimes)
    downtime_id BIGINT,
    
    -- Metrics
    response_time_minutes INTEGER,
    mttr_minutes INTEGER, -- Mean Time To Repair
    
    -- Approval
    approved_by BIGINT,
    approved_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_machine_maintenance (machine_id, reported_at),
    INDEX idx_status_maintenance (status),
    INDEX idx_assigned_technician (assigned_to),
    
    CONSTRAINT fk_machine_maintenance FOREIGN KEY (machine_id) REFERENCES machines(id),
    CONSTRAINT fk_reported_by FOREIGN KEY (reported_by) REFERENCES users(id),
    CONSTRAINT fk_assigned_technician FOREIGN KEY (assigned_to) REFERENCES users(id)
);
```

---

## Relationships

### Key Relationships Diagram (text-based)

```
production_orders (1) ─┬─ (1) khazwal_material_preparations
                       ├─ (1) khazwal_counting_results
                       ├─ (1) khazwal_cutting_results
                       ├─ (1-N) print_jobs
                       │
                       ├─ (N) verification_labels ← REDESIGNED!
                       ├─ (1) verification_po_summaries
                       │
                       ├─ (1) khazkhir_packaging_jobs
                       ├─ (1-N) khazkhir_rims
                       ├─ (N) khazkhir_inventory_movements
                       ├─ (1-N) khazkhir_delivery_orders
                       └─ (1-N) po_stage_tracking

print_jobs (1) ─────┬─ (N) print_production_logs
                    ├─ (N) print_downtimes
                    ├─ (N) print_quality_issues
                    └─ (1) print_job_summaries

verification_labels (1) ─┬─ (1) verification_inspection_results
                         └─ (N) verification_defects

khazkhir_delivery_orders (1) ─┬─ (N) khazkhir_delivery_tracking
                               └─ (1) khazkhir_proof_of_deliveries

machines (1) ─────┬─ (N) print_jobs
                  ├─ (N) print_downtimes
                  └─ (N) maintenance_tickets

users (1) ────────┬─ (N) print_jobs (operator)
                  ├─ (N) verification_labels (qc inspector - per label!)
                  ├─ (N) user_achievements
                  ├─ (N) point_transactions
                  └─ (N) notifications
```

**Key Changes in Relationships:**
- ❌ REMOVED: `verification_jobs` (per PO) 
- ✅ NEW: `verification_labels` (per label/rim) with 1-to-many from PO
- ❌ REMOVED: `khazkhir_palets` (no palet concept)
- ✅ SIMPLIFIED: Everything directly linked to PO

---

## Indexes & Performance

### Composite Indexes for Common Queries

```sql
-- Production Orders
CREATE INDEX idx_po_status_priority ON production_orders(current_stage, current_status, priority, due_date);
CREATE INDEX idx_po_due_date_status ON production_orders(due_date, current_status) WHERE deleted_at IS NULL;

-- Print Jobs - OEE Queries
CREATE INDEX idx_print_job_machine_date ON print_job_summaries(machine_id, finalized_at);
CREATE INDEX idx_print_job_operator_date ON print_job_summaries(operator_id, finalized_at);

-- Verification - Quality Queries
CREATE INDEX idx_verification_qc_date ON verification_summaries(qc_inspector_id, finalized_at);
CREATE INDEX idx_verification_hcs ON verification_summaries(hcs_percentage, finalized_at);

-- Warehouse - FIFO
CREATE INDEX idx_palets_fifo ON khazkhir_palets(production_order_id, stored_at) WHERE status = 'STORED';
CREATE INDEX idx_palets_location ON khazkhir_palets(storage_location, status);

-- Delivery - Performance
CREATE INDEX idx_do_delivery_date ON khazkhir_delivery_orders(delivery_date, status);
CREATE INDEX idx_do_customer ON khazkhir_delivery_orders(customer_id, delivery_date);

-- Activity Logs - Audit
CREATE INDEX idx_activity_user_time ON activity_logs(user_id, timestamp DESC);
CREATE INDEX idx_activity_entity ON activity_logs(entity_type, entity_id, timestamp DESC);

-- Alerts - Active Monitoring
CREATE INDEX idx_alerts_active ON alerts(status, severity, triggered_at) WHERE status = 'ACTIVE';
```

### Partial Indexes for Specific Scenarios

```sql
-- Only active POs
CREATE INDEX idx_active_pos ON production_orders(current_stage, current_status) 
    WHERE deleted_at IS NULL;

-- Only in-progress jobs
CREATE INDEX idx_active_print_jobs ON print_jobs(machine_id, status) 
    WHERE status IN ('SETUP', 'TEST_PRINT', 'RUNNING');

-- Only unread notifications
CREATE INDEX idx_unread_notifications ON notifications(user_id, created_at DESC) 
    WHERE is_read = FALSE;

-- Only stored palets
CREATE INDEX idx_stored_palets ON khazkhir_palets(storage_location) 
    WHERE status = 'STORED';
```

---

## ERD (Entity Relationship Diagram - Text Format)

```
┌─────────────────────────────────────────┐
│      production_orders (PO & OBC)       │ ◄── Everything centered here
│         (Central Entity)                │     NO palet concept
└──────────┬──────────────────────────────┘
           │
    ┌──────┴───────┬──────────┬────────────┬──────────┐
    │              │          │            │          │
    ▼              ▼          ▼            ▼          ▼
┌────────┐   ┌─────────┐ ┌─────────┐  ┌────────┐ ┌─────────┐
│KHAZWAL │   │  CETAK  │ │ VERIFY  │  │KHAZKHIR│ │  STAGE  │
│        │   │         │ │ LABELS  │  │        │ │ TRACKING│
│-prep   │   │-jobs    │ │ (NEW!)  │  │-package│ │         │
│-count  │   │-logs    │ │  PER    │  │-rims   │ │         │
│-cut    │   │-down    │ │ LABEL!  │  │-invent │ │         │
│        │   │-quality │ │-results │  │-deliver│ │         │
│        │   │-summary │ │-defects │  │-tracking│ │         │
└────────┘   └────┬────┘ │-po_sum  │  └────────┘ └─────────┘
                  │      └────┬────┘
                  │           │
                  ▼           ▼
             ┌─────────┐  ┌─────────┐
             │MACHINES │  │  USERS  │
             │ (no SAP)│  │  (NIP)  │
             │         │  │         │
             │-status  │  │-role    │
             │-pm      │  │-points  │
             │-oee     │  │-achieve │
             └────┬────┘  └────┬────┘
                  │            │
                  ▼            ▼
            ┌──────────┐  ┌──────────┐
            │MAINTEN-  │  │  NOTIF   │
            │  ANCE    │  │  POINTS  │
            │ TICKETS  │  │  ACHIEVE │
            └──────────┘  └──────────┘

┌─────────────────────────────────────────┐
│  DATA FROM SAP (via API - NOT in DB)   │
├─────────────────────────────────────────┤
│  • Customer (Dirjen Bea Cukai)         │
│  • Product Specifications              │
│  • Materials (Plat, Kertas, Tinta)     │
│  • Material Stock/Inventory            │
└─────────────────────────────────────────┘
```

---

## Implementation Notes

### 1. Timezone Handling
**CRITICAL:** Semua timestamps harus `TIMESTAMP WITH TIME ZONE` dengan default timezone Asia/Jakarta.

```php
// config/app.php
'timezone' => 'Asia/Jakarta',

// Migrations
$table->timestampTz('created_at');
$table->timestampTz('started_at')->nullable();
```

### 2. Soft Deletes
Gunakan soft deletes untuk semua master dan transaction tables (kecuali logs).

```php
// Migration
$table->timestampTz('deleted_at')->nullable();

// Model
use SoftDeletes;
```

### 3. Audit Trail Implementation
Semua perubahan critical data harus tercatat di `activity_logs`.

```php
// Eloquent Observer pattern
public function updated(Model $model)
{
    ActivityLog::create([
        'user_id' => auth()->id(),
        'action' => 'UPDATE',
        'entity_type' => get_class($model),
        'entity_id' => $model->id,
        'changes' => [
            'before' => $model->getOriginal(),
            'after' => $model->getAttributes()
        ],
        // ...
    ]);
}
```

### 4. JSONB Best Practices
Gunakan JSONB untuk data yang:
- Dynamic structure (e.g., defect breakdown, specifications)
- Array of objects (e.g., tinta requirements, photos)
- Metadata yang jarang di-query

**Indexing JSONB:**
```sql
CREATE INDEX idx_defect_breakdown ON print_quality_issues USING GIN (defect_breakdown);
```

### 5. Calculated Fields
Beberapa fields sebaiknya di-calculate di application layer, bukan database trigger:
- OEE components (availability, performance, quality)
- Percentages
- Durations
- Variances

**Alasan:** Easier to maintain, test, dan debug di application layer.

### 6. Status Tracking Pattern
Untuk setiap entity dengan status lifecycle, gunakan pattern:
1. Status column (VARCHAR)
2. Timestamp columns per status
3. Status tracking table (untuk history)

```php
// Model
public function transitionTo($newStatus)
{
    $this->status = $newStatus;
    $this->{$newStatus . '_at'} = now();
    $this->save();
    
    // Log to tracking table
    $this->statusTracking()->create([...]);
}
```

### 7. Performance Optimization
- Use `eager loading` untuk prevent N+1 queries
- Cache frequently accessed master data (users, machines, materials)
- Use queue untuk heavy operations (report generation, PDF creation)
- Implement pagination untuk all list views
- Use database views untuk complex aggregations

```php
// Eager loading
$jobs = PrintJob::with(['machine', 'operator', 'productionOrder'])
    ->where('status', 'running')
    ->get();
```

### 8. Real-time Updates
Untuk real-time dashboard, gunakan:
- **Laravel Echo + Pusher/Soketi** untuk WebSocket broadcasting
- **Redis** untuk caching real-time metrics
- **Database notifications** untuk in-app notifications

### 9. Data Retention
- **Transaction data:** Keep all (dengan archive strategy untuk data > 2 tahun)
- **Logs:** Keep 1 tahun di primary DB, archive older
- **Activity logs:** Keep 5 tahun (compliance)
- **Notifications:** Auto-delete after 30 hari untuk read notifications

### 10. Backup Strategy
- **Daily full backup** (automated)
- **Hourly incremental backup** (automated)
- **Pre-deployment backup** (manual)
- **Retention:** 30 hari daily, 12 bulan monthly
- **Test restore** quarterly

---

## Migration Strategy

### Phase 1: Core Infrastructure (Week 1-2)
```
1. users (NIP-based, no email required)
2. production_orders (with SAP references), po_stage_tracking
3. machines (only local - no SAP dependency)
4. activity_logs, notifications
```

### Phase 2: Khazanah Awal (Week 3)
```
5. khazwal_material_preparations
6. khazwal_counting_results
7. khazwal_cutting_results
```

### Phase 3: Cetak (Week 4-5)
```
8. print_jobs
9. print_production_logs
10. print_downtimes
11. print_quality_issues
12. print_job_summaries
13. maintenance_tickets
```

### Phase 4: Verifikasi (Week 6) - REDESIGNED!
```
13. verification_labels (per label/rim, not per PO!)
14. verification_inspection_results
15. verification_defects
16. verification_po_summaries (aggregated from labels)
```

### Phase 5: Khazanah Akhir - Packaging (Week 7)
```
17. khazkhir_packaging_jobs
18. khazkhir_rims
19. khazkhir_inventory_movements (no palet!)
```

### Phase 6: Khazanah Akhir - Warehouse (Week 8)
```
20. khazkhir_stock_opnames
21. khazkhir_stock_opname_details (track by PO, no palet)
```

### Phase 7: Khazanah Akhir - Delivery (Week 9)
```
22. khazkhir_delivery_orders
23. khazkhir_delivery_tracking
24. khazkhir_proof_of_deliveries
(NO picking_lists table - simplified, no palet concept)
```

### Phase 8: Gamification & Maintenance (Week 10)
```
27. gamification_achievements
28. user_achievements
29. point_transactions
30. alerts
31. maintenance_tickets
```

---

## Sample Queries

### 1. Get PO with Full Progress
```sql
SELECT 
    po.*,
    st.stage, st.status, st.started_at, st.completed_at,
    mp.status as material_prep_status,
    pj.status as print_status,
    vj.status as verification_status,
    pkg.status as packaging_status
FROM production_orders po
LEFT JOIN po_stage_tracking st ON st.production_order_id = po.id
LEFT JOIN khazwal_material_preparations mp ON mp.production_order_id = po.id
LEFT JOIN print_jobs pj ON pj.production_order_id = po.id
LEFT JOIN verification_jobs vj ON vj.production_order_id = po.id
LEFT JOIN khazkhir_packaging_jobs pkg ON pkg.production_order_id = po.id
WHERE po.id = ?;
```

### 2. Calculate Machine OEE (Daily)
```sql
SELECT 
    m.machine_code,
    DATE(pjs.finalized_at) as production_date,
    AVG(pjs.oee_percentage) as avg_oee,
    AVG(pjs.availability_percentage) as avg_availability,
    AVG(pjs.performance_percentage) as avg_performance,
    AVG(pjs.quality_percentage) as avg_quality,
    COUNT(*) as total_jobs,
    SUM(pjs.actual_quantity) as total_output
FROM print_job_summaries pjs
JOIN machines m ON m.id = pjs.machine_id
WHERE pjs.finalized_at >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY m.machine_code, DATE(pjs.finalized_at)
ORDER BY production_date DESC, m.machine_code;
```

### 3. Quality Trend (HCS Percentage) - NEW Query per Label
```sql
-- Aggregate from label-level results
SELECT 
    DATE(vir.finalized_at) as verification_date,
    COUNT(DISTINCT vir.production_order_id) as total_pos,
    COUNT(vir.id) as total_labels_completed,
    AVG(vir.hcs_percentage) as avg_hcs_percentage,
    SUM(vir.hcs_quantity) as total_hcs,
    SUM(vir.hcts_quantity) as total_hcts,
    ROUND(
        SUM(vir.hcs_quantity)::numeric / 
        (SUM(vir.hcs_quantity) + SUM(vir.hcts_quantity))::numeric * 100, 
        2
    ) as overall_hcs_percentage
FROM verification_inspection_results vir
WHERE vir.finalized_at >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY DATE(vir.finalized_at)
ORDER BY verification_date DESC;
```

### 3b. PO-Level Quality Summary
```sql
-- Get summary per PO (from verification_po_summaries)
SELECT 
    po.po_number,
    po.obc_number,
    vps.total_labels,
    vps.labels_completed,
    vps.completion_percentage,
    vps.hcs_percentage,
    vps.hcts_percentage,
    vps.avg_quality_score,
    vps.verification_status
FROM verification_po_summaries vps
JOIN production_orders po ON po.id = vps.production_order_id
WHERE vps.verification_status = 'IN_PROGRESS'
ORDER BY vps.completion_percentage ASC;
```

### 4. FIFO Compliance Check (NO PALET - track by PO)
```sql
SELECT 
    po.po_number,
    po.obc_number,
    po.due_date,
    pkg.completed_at as packaging_completed_at,
    im.to_location as storage_location,
    im.movement_date as stored_at,
    CURRENT_DATE - DATE(im.movement_date) as days_in_storage,
    im.balance_rims as current_stock_rims,
    CASE 
        WHEN CURRENT_DATE - DATE(im.movement_date) > 30 THEN 'CRITICAL'
        WHEN CURRENT_DATE - DATE(im.movement_date) > 21 THEN 'WARNING'
        WHEN CURRENT_DATE - DATE(im.movement_date) > 14 THEN 'MONITOR'
        ELSE 'OK'
    END as aging_status
FROM production_orders po
JOIN khazkhir_packaging_jobs pkg ON pkg.production_order_id = po.id
JOIN LATERAL (
    SELECT * FROM khazkhir_inventory_movements 
    WHERE production_order_id = po.id 
        AND movement_type = 'STOCK_IN'
    ORDER BY movement_date DESC 
    LIMIT 1
) im ON TRUE
WHERE pkg.status = 'COMPLETED'
    AND im.balance_rims > 0  -- Still have stock
ORDER BY days_in_storage DESC, po.due_date ASC;
```

### 5. On-Time Delivery Rate
```sql
SELECT 
    DATE_TRUNC('week', do.delivery_date) as week,
    COUNT(*) as total_deliveries,
    SUM(CASE WHEN do.delivered_at <= do.delivery_date THEN 1 ELSE 0 END) as on_time_deliveries,
    ROUND(
        SUM(CASE WHEN do.delivered_at <= do.delivery_date THEN 1 ELSE 0 END)::numeric 
        / COUNT(*)::numeric * 100, 
        2
    ) as otd_percentage
FROM khazkhir_delivery_orders do
WHERE do.status = 'DELIVERED'
    AND do.delivered_at >= CURRENT_DATE - INTERVAL '3 months'
GROUP BY DATE_TRUNC('week', do.delivery_date)
ORDER BY week DESC;
```

---

## Additional Considerations

### 1. Multi-Tenant Support (Future)
Jika ke depan ada multiple plants/factories:
```sql
-- Add to all tables
plant_id BIGINT NOT NULL,
CONSTRAINT fk_plant FOREIGN KEY (plant_id) REFERENCES plants(id)

-- Global filter di application layer
GlobalScope::where('plant_id', auth()->user()->plant_id)
```

### 2. Currency & Localization
- Currency: Simpan dalam integer (smallest unit - Rupiah)
- Format display di application layer
```php
// Database: 500000 (5 juta)
// Display: Rp 5.000.000
number_format($amount, 0, ',', '.')
```

### 3. File Storage
- **Photos/Documents:** Store di S3/MinIO
- **Database:** Store URL/path only
- **Naming convention:** `{type}/{po_number}/{timestamp}_{filename}`
- **Cleanup:** Periodic cleanup untuk orphaned files

### 4. API Design
Untuk customer portal dan mobile apps:
- RESTful API dengan Laravel Sanctum (authentication)
- API versioning (`/api/v1/...`)
- Rate limiting (60 requests/minute)
- Response caching untuk master data
- Pagination untuk list endpoints

### 5. Reporting Database (Future)
Untuk heavy analytics, consider **read replica** atau **separate reporting database**:
- ETL daily dari production DB
- Pre-aggregated tables untuk fast reporting
- No impact ke production performance

---

## Conclusion

Database architecture ini dirancang untuk:
✅ **Scalability** - Handle growth dalam volume produksi  
✅ **Performance** - Optimized indexes untuk common queries  
✅ **Maintainability** - Clear structure, well-documented  
✅ **Auditability** - Comprehensive audit trail  
✅ **Flexibility** - JSONB untuk dynamic data  
✅ **Data Integrity** - Foreign keys, constraints, validations  

**Next Steps:**
1. Review & approval dari tim
2. Setup development database (PostgreSQL 14+)
3. Create Laravel migrations (phase by phase)
4. Seed master data
5. Implement base models & relationships
6. Write repository/service layer
7. Unit tests untuk critical business logic

---

**Document Version:** 1.0  
**Last Updated:** 27 December 2025  
**Author:** Zulfikar Hidayatullah  
**Status:** Ready for Review
