-- Setup Database untuk Sirine Go App
-- Developer: Zulfikar Hidayatullah
-- Timezone: Asia/Jakarta (WIB)

-- Buat database jika belum ada
CREATE DATABASE IF NOT EXISTS sirine_go 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

-- Gunakan database
USE sirine_go;

-- Set timezone
SET time_zone = '+07:00';

-- Catatan:
-- Tabel akan dibuat otomatis oleh GORM Auto Migration
-- Jalankan aplikasi Go untuk membuat tabel secara otomatis
