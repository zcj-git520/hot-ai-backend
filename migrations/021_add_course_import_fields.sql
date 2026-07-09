-- ======================================================
-- 课程动态导入与同步系统 - 数据库变更
-- 日期：2026-07-07
-- 说明：
--   1) learning_paths 加 source 关联字段
--   2) path_chapters 加源 key 字段
--   3) 新增 course_sync_log 同步日志表
-- 兼容性：
--   - 全部字段加 NULL DEFAULT，不破坏现有数据
--   - 索引加 IF NOT EXISTS 兼容 MariaDB 10.5
-- ======================================================

USE `hot_ai`;

-- ============================================
-- 1) learning_paths 加 source 关联字段
-- ============================================
-- 检查列是否存在后再 ALTER，避免 MySQL 1060 重复列错误
SET @col_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'learning_paths'
    AND COLUMN_NAME = 'source_id');
SET @sql = IF(@col_exists = 0,
  'ALTER TABLE learning_paths ADD COLUMN source_id VARCHAR(64) NULL COMMENT ''外部源标识, NULL=管理员手工创建''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'learning_paths'
    AND COLUMN_NAME = 'source_url');
SET @sql = IF(@col_exists = 0,
  'ALTER TABLE learning_paths ADD COLUMN source_url VARCHAR(500) NULL COMMENT ''外部源详情链接''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'learning_paths'
    AND COLUMN_NAME = 'source_meta');
SET @sql = IF(@col_exists = 0,
  'ALTER TABLE learning_paths ADD COLUMN source_meta JSON NULL COMMENT ''同步时的元信息快照''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'learning_paths'
    AND COLUMN_NAME = 'last_synced_at');
SET @sql = IF(@col_exists = 0,
  'ALTER TABLE learning_paths ADD COLUMN last_synced_at DATETIME NULL COMMENT ''最后一次同步成功时间''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 索引（先检查再加）
SET @idx_exists := (SELECT COUNT(*) FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'learning_paths'
    AND INDEX_NAME = 'idx_source');
SET @sql = IF(@idx_exists = 0,
  'ALTER TABLE learning_paths ADD INDEX idx_source (source_id)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ============================================
-- 2) path_chapters 加源 key 字段
-- ============================================
SET @col_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'path_chapters'
    AND COLUMN_NAME = 'source_chapter_key');
SET @sql = IF(@col_exists = 0,
  'ALTER TABLE path_chapters ADD COLUMN source_chapter_key VARCHAR(255) NULL COMMENT ''源内的章节标识(通常是文件相对路径)''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'path_chapters'
    AND COLUMN_NAME = 'content_hash');
SET @sql = IF(@col_exists = 0,
  'ALTER TABLE path_chapters ADD COLUMN content_hash CHAR(64) NULL COMMENT ''内容 sha256, 用于判断是否真有变化''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx_exists := (SELECT COUNT(*) FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'path_chapters'
    AND INDEX_NAME = 'idx_source_key');
SET @sql = IF(@idx_exists = 0,
  'ALTER TABLE path_chapters ADD INDEX idx_source_key (path_id, source_chapter_key)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ============================================
-- 3) course_sync_log 表 (新增)
-- ============================================
CREATE TABLE IF NOT EXISTS course_sync_log (
  id              INT AUTO_INCREMENT PRIMARY KEY,
  source_id       VARCHAR(64)    NOT NULL,
  trigger_type    ENUM('manual','scheduled') NOT NULL,
  status          ENUM('running','success','failed','partial') NOT NULL,
  started_at      DATETIME       NOT NULL,
  finished_at     DATETIME       NULL,
  duration_ms     INT            NULL,
  stats           JSON           NULL COMMENT '{created,updated,deleted,preserved,failed}',
  error_message   TEXT           NULL,
  operator        VARCHAR(64)    NULL COMMENT '管理员用户名 或 system',
  INDEX idx_source_time (source_id, started_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT '课程同步日志';

SELECT 'course import fields migration done' AS result;
