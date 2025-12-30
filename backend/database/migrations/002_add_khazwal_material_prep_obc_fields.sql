-- Migration: Add OBC Master Integration Fields to Khazwal Material Preparations
-- Date: 2025-12-30
-- Phase: Material Preparation Phase 6
-- Developer: Zulfikar Hidayatullah

-- Add new fields untuk enhanced OBC Master integration
ALTER TABLE khazwal_material_preparations 
ADD COLUMN IF NOT EXISTS plat_scanned_code VARCHAR(50) COMMENT 'Actual scanned plat barcode',
ADD COLUMN IF NOT EXISTS plat_match BOOLEAN DEFAULT FALSE COMMENT 'Whether scanned plat matches expected',
ADD COLUMN IF NOT EXISTS kertas_blanko_variance_percentage DECIMAL(5,2) COMMENT 'Variance percentage dari target',
ADD COLUMN IF NOT EXISTS tinta_low_stock_flags JSON COMMENT 'Low stock warnings per color';

-- Add indexes untuk performance
CREATE INDEX IF NOT EXISTS idx_khazwal_plat_match ON khazwal_material_preparations(plat_match);

-- Rollback script (jika diperlukan)
-- ALTER TABLE khazwal_material_preparations 
-- DROP COLUMN IF EXISTS plat_scanned_code,
-- DROP COLUMN IF EXISTS plat_match,
-- DROP COLUMN IF EXISTS kertas_blanko_variance_percentage,
-- DROP COLUMN IF EXISTS tinta_low_stock_flags;
-- DROP INDEX IF EXISTS idx_khazwal_plat_match;
