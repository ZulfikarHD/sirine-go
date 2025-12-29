-- Database Optimization untuk Sprint 6
-- Performance indexes, composite indexes, dan query optimization
-- Developer: Zulfikar Hidayatullah
-- Date: 2025-01-29

USE sirine_go;

-- ====================
-- COMPOSITE INDEXES
-- ====================

-- Composite index untuk notifications table (user_id + is_read)
-- untuk optimize query unread notifications count dan filter
CREATE INDEX IF NOT EXISTS idx_notifications_user_read 
ON notifications(user_id, is_read);

-- Composite index untuk notifications (user_id + created_at)
-- untuk optimize pagination dengan sorting
CREATE INDEX IF NOT EXISTS idx_notifications_user_created 
ON notifications(user_id, created_at DESC);

-- Composite index untuk user_sessions (user_id + is_revoked + expires_at)
-- untuk optimize active sessions query
CREATE INDEX IF NOT EXISTS idx_sessions_user_active 
ON user_sessions(user_id, is_revoked, expires_at);

-- Composite index untuk activity_logs (user_id + action + created_at)
-- untuk optimize filtering dan pagination di audit logs
CREATE INDEX IF NOT EXISTS idx_activity_user_action_created 
ON activity_logs(user_id, action, created_at DESC);

-- Composite index untuk activity_logs (entity_type + entity_id)
-- untuk optimize query by specific entity
CREATE INDEX IF NOT EXISTS idx_activity_entity 
ON activity_logs(entity_type, entity_id);

-- Composite index untuk password_reset_tokens (user_id + used_at + expires_at)
-- untuk optimize valid token lookup
CREATE INDEX IF NOT EXISTS idx_reset_tokens_valid 
ON password_reset_tokens(user_id, used_at, expires_at);

-- Composite index untuk achievements (user_id + earned_at)
-- untuk optimize user achievements display
CREATE INDEX IF NOT EXISTS idx_user_achievements 
ON user_achievements(user_id, earned_at DESC);

-- ====================
-- FULL TEXT SEARCH INDEXES
-- ====================

-- Full-text index untuk users (full_name, email) untuk fast search
ALTER TABLE users ADD FULLTEXT INDEX idx_users_search (full_name, email);

-- ====================
-- QUERY OPTIMIZATION VIEWS
-- ====================

-- View untuk active users dengan recent login
CREATE OR REPLACE VIEW v_active_users AS
SELECT 
    u.id,
    u.nip,
    u.full_name,
    u.email,
    u.role,
    u.department,
    u.last_login_at,
    u.total_points,
    u.level
FROM users u
WHERE u.status = 'ACTIVE' 
    AND u.deleted_at IS NULL
ORDER BY u.last_login_at DESC;

-- View untuk unread notifications per user
CREATE OR REPLACE VIEW v_unread_notifications AS
SELECT 
    user_id,
    COUNT(*) as unread_count,
    MAX(created_at) as latest_notification_at
FROM notifications
WHERE is_read = FALSE
GROUP BY user_id;

-- View untuk user sessions summary
CREATE OR REPLACE VIEW v_active_sessions AS
SELECT 
    s.user_id,
    u.nip,
    u.full_name,
    COUNT(*) as active_session_count,
    MAX(s.expires_at) as latest_session_expiry
FROM user_sessions s
JOIN users u ON s.user_id = u.id
WHERE s.is_revoked = FALSE 
    AND s.expires_at > NOW()
GROUP BY s.user_id, u.nip, u.full_name;

-- ====================
-- STORED PROCEDURES untuk Common Operations
-- ====================

-- Procedure untuk cleanup expired sessions
DELIMITER //
CREATE PROCEDURE IF NOT EXISTS sp_cleanup_expired_sessions()
BEGIN
    -- Mark expired sessions sebagai revoked
    UPDATE user_sessions 
    SET is_revoked = TRUE 
    WHERE expires_at < NOW() 
        AND is_revoked = FALSE;
    
    -- Delete old revoked sessions (older than 30 days)
    DELETE FROM user_sessions 
    WHERE is_revoked = TRUE 
        AND created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
    
    SELECT ROW_COUNT() as sessions_cleaned;
END //
DELIMITER ;

-- Procedure untuk cleanup expired password reset tokens
DELIMITER //
CREATE PROCEDURE IF NOT EXISTS sp_cleanup_expired_reset_tokens()
BEGIN
    -- Delete expired and used tokens (older than 7 days)
    DELETE FROM password_reset_tokens 
    WHERE (expires_at < NOW() OR used_at IS NOT NULL)
        AND created_at < DATE_SUB(NOW(), INTERVAL 7 DAY);
    
    SELECT ROW_COUNT() as tokens_cleaned;
END //
DELIMITER ;

-- Procedure untuk get user statistics
DELIMITER //
CREATE PROCEDURE IF NOT EXISTS sp_get_user_stats(IN p_user_id BIGINT)
BEGIN
    SELECT 
        u.id,
        u.nip,
        u.full_name,
        u.total_points,
        u.level,
        COUNT(DISTINCT s.id) as active_sessions,
        COUNT(DISTINCT n.id) as unread_notifications,
        COUNT(DISTINCT ua.achievement_id) as total_achievements,
        u.last_login_at
    FROM users u
    LEFT JOIN user_sessions s ON u.id = s.user_id 
        AND s.is_revoked = FALSE 
        AND s.expires_at > NOW()
    LEFT JOIN notifications n ON u.id = n.user_id 
        AND n.is_read = FALSE
    LEFT JOIN user_achievements ua ON u.id = ua.user_id
    WHERE u.id = p_user_id
    GROUP BY u.id, u.nip, u.full_name, u.total_points, u.level, u.last_login_at;
END //
DELIMITER ;

-- ====================
-- SCHEDULED CLEANUP EVENTS
-- ====================

-- Enable event scheduler
SET GLOBAL event_scheduler = ON;

-- Event untuk cleanup expired sessions (run setiap 1 jam)
CREATE EVENT IF NOT EXISTS evt_cleanup_sessions
ON SCHEDULE EVERY 1 HOUR
DO
    CALL sp_cleanup_expired_sessions();

-- Event untuk cleanup expired reset tokens (run setiap hari jam 2 pagi)
CREATE EVENT IF NOT EXISTS evt_cleanup_reset_tokens
ON SCHEDULE EVERY 1 DAY
STARTS TIMESTAMP(CURRENT_DATE) + INTERVAL 2 HOUR
DO
    CALL sp_cleanup_expired_reset_tokens();

-- Event untuk cleanup old activity logs (run setiap minggu, keep 90 days)
CREATE EVENT IF NOT EXISTS evt_cleanup_old_logs
ON SCHEDULE EVERY 1 WEEK
STARTS TIMESTAMP(CURRENT_DATE) + INTERVAL 3 HOUR
DO
    DELETE FROM activity_logs 
    WHERE created_at < DATE_SUB(NOW(), INTERVAL 90 DAY);

-- ====================
-- PERFORMANCE OPTIMIZATION SETTINGS
-- ====================

-- Optimize InnoDB settings untuk better performance
-- Note: These should be set in my.cnf for production

-- SHOW VARIABLES LIKE 'innodb_buffer_pool_size';
-- SET GLOBAL innodb_buffer_pool_size = 1073741824; -- 1GB (adjust based on server)

-- SHOW VARIABLES LIKE 'innodb_log_file_size';
-- SET GLOBAL innodb_log_file_size = 268435456; -- 256MB

-- Query cache settings (if available)
-- SET GLOBAL query_cache_size = 67108864; -- 64MB
-- SET GLOBAL query_cache_type = 1;

-- ====================
-- MONITORING QUERIES
-- ====================

-- Query untuk check index usage
-- SELECT * FROM sys.schema_unused_indexes WHERE object_schema = 'sirine_go';

-- Query untuk check table sizes
-- SELECT 
--     table_name,
--     ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb,
--     table_rows
-- FROM information_schema.TABLES
-- WHERE table_schema = 'sirine_go'
-- ORDER BY (data_length + index_length) DESC;

-- Query untuk check slow queries
-- SELECT * FROM mysql.slow_log ORDER BY start_time DESC LIMIT 10;

-- ====================
-- ANALYZE AND OPTIMIZE TABLES
-- ====================

-- Run ANALYZE untuk update statistics
ANALYZE TABLE users;
ANALYZE TABLE user_sessions;
ANALYZE TABLE password_reset_tokens;
ANALYZE TABLE activity_logs;
ANALYZE TABLE notifications;
ANALYZE TABLE achievements;
ANALYZE TABLE user_achievements;

-- Run OPTIMIZE untuk defragment tables (jalankan saat low traffic)
-- OPTIMIZE TABLE users;
-- OPTIMIZE TABLE user_sessions;
-- OPTIMIZE TABLE activity_logs;

-- ====================
-- VERIFICATION QUERIES
-- ====================

-- Verify indexes created
SELECT 
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    SEQ_IN_INDEX,
    INDEX_TYPE
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = 'sirine_go'
    AND TABLE_NAME IN ('users', 'user_sessions', 'notifications', 'activity_logs', 'password_reset_tokens')
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX;

-- Check database size
SELECT 
    table_schema AS 'Database',
    ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
FROM information_schema.TABLES
WHERE table_schema = 'sirine_go'
GROUP BY table_schema;

-- ====================
-- USAGE NOTES
-- ====================

/*
DEPLOYMENT CHECKLIST:
1. Backup database sebelum run optimization
   mysqldump -u root -p sirine_go > backup_before_optimization.sql

2. Run optimization script
   mysql -u root -p sirine_go < optimization.sql

3. Monitor query performance setelah optimization
   - Check slow query log
   - Monitor CPU dan memory usage
   - Verify index usage

4. Schedule regular maintenance
   - Run ANALYZE TABLE setiap minggu
   - Run OPTIMIZE TABLE setiap bulan (saat low traffic)
   - Monitor table sizes dan cleanup old data

5. Performance tuning berdasarkan workload
   - Adjust buffer pool size based on RAM
   - Monitor index usage dan drop unused indexes
   - Consider partitioning untuk large tables (activity_logs)

RECOMMENDED SERVER CONFIGURATION (my.cnf):
[mysqld]
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
innodb_flush_log_at_trx_commit = 2
innodb_flush_method = O_DIRECT
max_connections = 200
query_cache_size = 64M
query_cache_type = 1
slow_query_log = 1
long_query_time = 2
*/
