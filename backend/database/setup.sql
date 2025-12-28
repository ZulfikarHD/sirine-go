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

-- Tabel Users dengan enhanced fields untuk authentication
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nip VARCHAR(5) NOT NULL UNIQUE COMMENT 'Nomor Induk Pegawai (max 5 digit)',
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('ADMIN', 'MANAGER', 'STAFF_KHAZWAL', 'OPERATOR_CETAK', 'QC_INSPECTOR', 'VERIFIKATOR', 'STAFF_KHAZKHIR') NOT NULL,
    department ENUM('KHAZWAL', 'CETAK', 'VERIFIKASI', 'KHAZKHIR') NOT NULL,
    shift ENUM('PAGI', 'SIANG', 'MALAM') DEFAULT 'PAGI',
    profile_photo_url VARCHAR(500),
    total_points INT DEFAULT 0 COMMENT 'Total gamification points',
    level VARCHAR(20) DEFAULT 'Bronze' COMMENT 'User level based on points',
    status ENUM('ACTIVE', 'INACTIVE', 'SUSPENDED') DEFAULT 'ACTIVE',
    must_change_password BOOLEAN DEFAULT TRUE COMMENT 'Force password change on first login',
    last_login_at TIMESTAMP NULL,
    failed_login_attempts INT DEFAULT 0,
    locked_until TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_nip (nip),
    INDEX idx_email (email),
    INDEX idx_role (role),
    INDEX idx_department (department),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel User Sessions untuk JWT token tracking
CREATE TABLE IF NOT EXISTS user_sessions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    token_hash VARCHAR(255) NOT NULL COMMENT 'SHA256 hash dari JWT token',
    refresh_token_hash VARCHAR(255) COMMENT 'Hash dari refresh token',
    device_info VARCHAR(500),
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_token_hash (token_hash),
    INDEX idx_expires_at (expires_at),
    INDEX idx_is_revoked (is_revoked)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel Password Reset Tokens
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_token_hash (token_hash),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel Activity Logs untuk audit trail
CREATE TABLE IF NOT EXISTS activity_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    action ENUM('CREATE', 'UPDATE', 'DELETE', 'LOGIN', 'LOGOUT', 'PASSWORD_CHANGE') NOT NULL,
    entity_type VARCHAR(50) NOT NULL COMMENT 'Table name atau entity type',
    entity_id BIGINT UNSIGNED,
    changes JSON COMMENT 'Before/after values dalam JSON format',
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_entity_type (entity_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel Notifications untuk in-app notifications
CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type ENUM('INFO', 'SUCCESS', 'WARNING', 'ERROR', 'ACHIEVEMENT') DEFAULT 'INFO',
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    INDEX idx_user_read (user_id, is_read)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel Achievements untuk gamification system
CREATE TABLE IF NOT EXISTS achievements (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE COMMENT 'Unique code untuk achievement',
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    icon VARCHAR(100) COMMENT 'Icon identifier atau emoji',
    points INT NOT NULL DEFAULT 0 COMMENT 'Points yang didapat saat unlock',
    category ENUM('LOGIN', 'PRODUCTIVITY', 'QUALITY', 'MILESTONE') DEFAULT 'MILESTONE',
    criteria JSON COMMENT 'Criteria untuk unlock achievement',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_code (code),
    INDEX idx_category (category),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel User Achievements untuk tracking unlocked achievements
CREATE TABLE IF NOT EXISTS user_achievements (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    achievement_id BIGINT UNSIGNED NOT NULL,
    unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (achievement_id) REFERENCES achievements(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_achievement (user_id, achievement_id),
    INDEX idx_user_id (user_id),
    INDEX idx_achievement_id (achievement_id),
    INDEX idx_unlocked_at (unlocked_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Seed initial admin user
-- Password: Admin@123 (akan di-hash oleh aplikasi)
-- NIP: 99999
INSERT INTO users (nip, full_name, email, phone, password_hash, role, department, must_change_password, status) 
VALUES (
    '99999', 
    'Administrator', 
    'admin@sirine.local', 
    '081234567890',
    -- Bcrypt hash untuk "Admin@123" dengan cost 12
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIRQ8DrKIW',
    'ADMIN',
    'KHAZWAL',
    FALSE,
    'ACTIVE'
) ON DUPLICATE KEY UPDATE nip=nip;

-- Seed initial achievements
INSERT INTO achievements (code, name, description, icon, points, category, criteria) VALUES
('FIRST_LOGIN', 'Login Pertama', 'Berhasil login untuk pertama kalinya ke sistem', '🎉', 10, 'LOGIN', '{"type": "login_count", "value": 1}'),
('WEEK_STREAK', 'Semangat Seminggu', 'Login selama 7 hari berturut-turut', '🔥', 50, 'LOGIN', '{"type": "login_streak", "value": 7}'),
('MONTH_STREAK', 'Dedikasi Sebulan', 'Login selama 30 hari berturut-turut', '⭐', 100, 'LOGIN', '{"type": "login_streak", "value": 30}'),
('PROFILE_COMPLETE', 'Profil Lengkap', 'Melengkapi informasi profil dan foto', '✨', 20, 'MILESTONE', '{"type": "profile_complete"}'),
('EARLY_BIRD', 'Si Pagi Hari', 'Login sebelum jam 07:00 WIB sebanyak 10 kali', '🌅', 30, 'LOGIN', '{"type": "early_login", "value": 10}'),
('NIGHT_OWL', 'Burung Hantu', 'Login setelah jam 20:00 WIB sebanyak 10 kali', '🦉', 30, 'LOGIN', '{"type": "late_login", "value": 10}')
ON DUPLICATE KEY UPDATE code=code;

-- Catatan:
-- GORM Auto Migration akan handle kolom tambahan jika ada perubahan model
-- Jalankan aplikasi Go untuk sync dengan GORM models
