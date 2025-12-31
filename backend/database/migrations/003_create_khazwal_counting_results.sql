-- Migration: Create khazwal_counting_results table untuk Epic 2 Penghitungan
-- Purpose: Tracking penghitungan hasil cetak dengan quantities, defects, dan variances

-- Create khazwal_counting_results table
CREATE TABLE IF NOT EXISTS khazwal_counting_results (
    id BIGSERIAL PRIMARY KEY,
    production_order_id BIGINT NOT NULL UNIQUE,
    
    -- Counting Results (Lembar Besar)
    quantity_good INTEGER NOT NULL DEFAULT 0 CHECK (quantity_good >= 0),
    quantity_defect INTEGER NOT NULL DEFAULT 0 CHECK (quantity_defect >= 0),
    total_counted INTEGER GENERATED ALWAYS AS (quantity_good + quantity_defect) STORED,
    variance_from_target INTEGER,
    
    -- Percentages (calculated in application, stored for reporting)
    percentage_good DECIMAL(5,2),
    percentage_defect DECIMAL(5,2),
    
    -- Defect Breakdown (JSONB: [{ "type": "...", "quantity": N }])
    defect_breakdown JSONB,
    
    -- Status & Timing
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED')),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER,
    
    -- Staff
    counted_by BIGINT,
    
    -- Notes
    variance_reason TEXT,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_po_counting FOREIGN KEY (production_order_id) 
        REFERENCES production_orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_counted_by FOREIGN KEY (counted_by) 
        REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes untuk performance optimization
CREATE INDEX idx_counting_status ON khazwal_counting_results(status);
CREATE INDEX idx_counting_po ON khazwal_counting_results(production_order_id);
CREATE INDEX idx_counting_staff ON khazwal_counting_results(counted_by, started_at);
CREATE INDEX idx_counting_completed ON khazwal_counting_results(completed_at) WHERE completed_at IS NOT NULL;

-- Add new statuses untuk PO tracking counting stage
-- Check if the enum type exists and has these values
DO $$ 
BEGIN
    -- Add KHAZWAL_COUNTING stage if not exists
    IF NOT EXISTS (
        SELECT 1 FROM pg_enum 
        WHERE enumlabel = 'KHAZWAL_COUNTING' 
        AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'po_stage')
    ) THEN
        ALTER TYPE po_stage ADD VALUE 'KHAZWAL_COUNTING' AFTER 'CETAK';
    END IF;

    -- Add KHAZWAL_CUTTING stage if not exists
    IF NOT EXISTS (
        SELECT 1 FROM pg_enum 
        WHERE enumlabel = 'KHAZWAL_CUTTING' 
        AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'po_stage')
    ) THEN
        ALTER TYPE po_stage ADD VALUE 'KHAZWAL_CUTTING' AFTER 'KHAZWAL_COUNTING';
    END IF;
END $$;

-- Add new statuses for counting stage
DO $$ 
BEGIN
    -- SELESAI_CETAK -> waiting for counting
    IF NOT EXISTS (
        SELECT 1 FROM pg_enum 
        WHERE enumlabel = 'SELESAI_CETAK' 
        AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'po_status')
    ) THEN
        ALTER TYPE po_status ADD VALUE 'SELESAI_CETAK';
    END IF;

    -- WAITING_COUNTING -> in counting queue
    IF NOT EXISTS (
        SELECT 1 FROM pg_enum 
        WHERE enumlabel = 'WAITING_COUNTING' 
        AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'po_status')
    ) THEN
        ALTER TYPE po_status ADD VALUE 'WAITING_COUNTING';
    END IF;

    -- SEDANG_DIHITUNG -> counting in progress
    IF NOT EXISTS (
        SELECT 1 FROM pg_enum 
        WHERE enumlabel = 'SEDANG_DIHITUNG' 
        AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'po_status')
    ) THEN
        ALTER TYPE po_status ADD VALUE 'SEDANG_DIHITUNG';
    END IF;

    -- SIAP_POTONG -> ready for cutting
    IF NOT EXISTS (
        SELECT 1 FROM pg_enum 
        WHERE enumlabel = 'SIAP_POTONG' 
        AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'po_status')
    ) THEN
        ALTER TYPE po_status ADD VALUE 'SIAP_POTONG';
    END IF;
END $$;

-- Create updated_at trigger untuk auto-update timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_counting_updated_at BEFORE UPDATE ON khazwal_counting_results
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments untuk documentation
COMMENT ON TABLE khazwal_counting_results IS 'Tabel untuk tracking hasil penghitungan setelah cetak';
COMMENT ON COLUMN khazwal_counting_results.quantity_good IS 'Jumlah lembar besar yang baik hasil penghitungan';
COMMENT ON COLUMN khazwal_counting_results.quantity_defect IS 'Jumlah lembar besar yang rusak hasil penghitungan';
COMMENT ON COLUMN khazwal_counting_results.total_counted IS 'Total hasil penghitungan (good + defect), generated column';
COMMENT ON COLUMN khazwal_counting_results.variance_from_target IS 'Selisih antara total counted dengan target quantity';
COMMENT ON COLUMN khazwal_counting_results.defect_breakdown IS 'Breakdown jenis kerusakan dalam format JSON';
COMMENT ON COLUMN khazwal_counting_results.status IS 'Status penghitungan: PENDING, IN_PROGRESS, COMPLETED';
COMMENT ON COLUMN khazwal_counting_results.duration_minutes IS 'Durasi penghitungan dalam menit';
